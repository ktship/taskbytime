package taskbytime
import (
	"fmt"
	"time"
	"log"
)

// 타임기반의 일처리 엔진
// 모든 태스크 수량 처리는 update 함수를 호출해서 갱신한 후에 할 것.

// 모든 태스크는 TaskManager 에게 일을 맡김. 인터페이스 역할을 함.
type TaskManager struct {
	io *taskio
}

// TaskManager 생성
// 입력 : TaskIO, TaskCacheIO
// 리턴 : TaskManager 인스턴스
func NewTaskManager() *TaskManager {
	return &TaskManager{
		io:			NewTaskIO(),
	}
}

// DB에 저장될 내용. 변하는 데이터
type TaskVariable struct {
	checkTime	int
	curNum		int
}

// --------------------------------------------------------------------

// validate : taskId 가 파라메터에 있는지 체크
// 인풋 : taskId
// 리턴 : error
func (t *TaskManager)validate(tid int) error {
	_, ok := taskDatas[tid]
	if !ok {
		return fmt.Errorf("taskId : %d 는 taskDatas에 없음.", tid)
	}
	return nil
}

// CreateTask : 인터벌 시작. 시작수량과 최대수량이 같거나 크면 타임체크를 시작하지 않음.
// 인풋 : none
// 리턴 : 현재수량, 인터벌, remainTime( 0이면 스톱상태 ), err
func (t *TaskManager)CreateTask(uid int, tid int) (curNum int, interval int, remainTime int, err error) {
	if err := t.validate(tid); err != nil {
		return 0, 0, 0, err
	}
	// 메모리에서 태스크 관련 파라메터 get
	taskd := taskDatas[tid]

	// write 할 내용 편집
	checkTime := time.Now().Unix()

	// 디비에 씀
	data := make(map[string]interface{})
	data["ct"] = int(checkTime)
	data["num"] = taskd.startNum
	err = t.io.WriteTaskIO(uid, tid, data)
	if err != nil {
		return 0, 0, 0, err
	}

	return taskd.startNum, taskd.interval, taskd.interval, nil
}

// CalcTask : 수량을 더하고 뺌 (하트 사용, 하트 선물등. 시간에 따른 수량 변화는 update 함수로 처리)
// 인풋 : 추가/감소 수량
// 리턴 : 현재수량, 인터벌, remainTime( 0이면 스톱상태 ), err
func (t *TaskManager)CalcTask(uid int, tid int, num int) (curNum int, interval int, remainTime int, err error) {
	if err := t.validate(tid); err != nil {
		return 0, 0, 0, err
	}
	// 메모리에서 태스크 관련 파라메터 get
	taskd := taskDatas[uid]

	// 캐쉬에서 user 데이터 get
	dat, err := t.io.ReadTaskIO(uid, tid)
	if err != nil {
		return 0, 0, 0, err
	}

	// 현재 시간을 기준으로 업데이트 실시
	newNum, newRemainTime, NewCheckTime, err := t.update(uid, tid, dat)
	addedNum := newNum + num

	ctime := int(time.Now().Unix())
	// 계산된 수량이 음수이면 0으로 초기화, 그리고 체크타임등을 초기화.
	// 수량이 Max 치를 넘어서는 것에 대해서는 제한하지 않음. (update 함수내부에서는 제한함)
	if addedNum < 0 {
		addedNum = 0
		newRemainTime = taskd.interval
		NewCheckTime = ctime
	}

	// 최대치 이상이었다가 최대치보다 작아 졌을때는 체크타임과 remainTime 을 초기화해준다.
	// 예> 풀 하트였다가 하트 소모하면 그때부터 다시 재계산함.
	if newNum >= taskd.maxNum && addedNum < taskd.maxNum {
		newRemainTime = taskd.interval
		NewCheckTime = ctime
	}

	// 디비에 씀
	dat = make(map[string]interface{})
	dat["ct"] = NewCheckTime
	dat["num"] = addedNum
	err = t.io.WriteTaskIO(uid, tid, dat)
	if err != nil {
		return 0, 0, 0, err
	}

	return addedNum, taskd.interval, newRemainTime, nil
}

// DeleteTask : 태스크 삭제
// 인풋 : user id, task index
// 리턴 : err
func (t *TaskManager)DeleteTask(uid int, tid int) (err error) {
	if err := t.validate(tid); err != nil {
		return err
	}

	// 디비에서 제거
	err = t.io.DelTaskIO(uid, tid)
	if err != nil {
		return err
	}

	return nil
}

// calcNum : 체크시간을 기준으로 현재 수량과 남은 시간, 필요하다면(현재수량이 증가 했을때) 체크시간 업데이트
// 인풋 : user id, task index
// 리턴 : err
func (t *TaskManager)update(uid int, tid int, dat map[string]interface{}) (newNum int, newRemainTime int, newCheckTime int, err error) {
	if err := t.validate(tid); err != nil {
		return 0,0,0,err
	}

	taskd := taskDatas[tid]
	curTime := int(time.Now().Unix())

	// 총수량이 최대수량보다 많으면 더 볼것도 없이,
	// 총수량은 최대수량으로 고정시키고 남은 시간과 체크시간을 현재 시간 기준으로 바꿈.
	curNum := dat["num"].(int)
	if curNum >= taskd.maxNum {
		log.Printf("	- update curNum:%d interval: %d checktime : %d", curNum, taskd.interval, curTime)
		return curNum, taskd.interval, curTime, nil
	}

	// 체크시간과 현재시간의 차이에서 인터벌로 나눈 수만큼 갯수를 증가시킴.
	oldCheckTime := dat["ct"].(int)
	curInterval := curTime - oldCheckTime
	curInterval = Max(0, curInterval)
	if taskd.interval == 0 {
		return 0,0,0,fmt.Errorf("taskData.interval(%d) is 0", tid)
	}
	portion := curInterval / taskd.interval
	mod := curInterval % taskd.interval

	// 시간 계산후에 총수량이 최대수량보다 많아버리면,
	// 총수량은 최대수량으로 고정시키고 남은 시간과 체크시간을 현재 시간 기준으로 바꿈.
	rNum := curNum + portion
	rNum = Min(rNum, taskd.maxNum)

	// 새 체크시간 갱신
	rCheckTime := curTime - mod
	// 새 남은 시간 갱신
	remain := taskd.interval - mod
	log.Printf("	- update curInterval(%d) taskd.interval(%d) remain(%d)", curInterval, taskd.interval, remain)

	return rNum, remain, rCheckTime, nil
}
