package taskbytime
import (
	"fmt"
	"time"
)

// 타임기반의 일처리 엔진
// 모든 태스크 수량 처리는 update 함수를 호출해서 갱신한 후에 할 것.

// 모든 태스크는 TaskManager 에게 일을 맡김. 인터페이스 역할을 함.
type TaskManager struct {
	// DB 에 읽고 쓸 io
	io TaskIO
	// 레디스 캐시내용을 읽고 쓸 io : 메모리 처럼 사용가능
	cacheIo TaskCacheIO
	// user id
	uid 	uint32
	// 처리할 대상 태스크 아이디
	taskId	uint32
}

// TaskManager 생성
// 입력 : TaskIO, TaskCacheIO
// 리턴 : TaskManager 인스턴스
func NewTaskManager(taskIo TaskIO, cacheIo TaskCacheIO, uid uint32, taskId uint32) *TaskManager {
	return &TaskManager{
		io:			taskIo,
		cacheIo:	cacheIo,
		uid:		uid,
		taskId:		taskId,
	}
}

// DB에 저장될 내용. 변하는 데이터
type TaskVariable struct {
	checkTime	int64
	curNum		int32
}

// DB에 태스크 관련 데이터를 저장할 인터페이스
type TaskIO interface {
	Read(id uint32, id2 uint32) (taskVar map[string]interface{}, err error)
	Write(id uint32, id2 uint32, taskVar map[string]interface{}) error
	Del(id uint32, id2 uint32) (err error)
}

// 캐쉬에서 태스크 관련 정보들을 읽고 쓰는 인터페이스
type TaskCacheIO interface {
	GetCacheTask(id uint32, id2 uint32) (taskVar map[string]interface{}, err error)
	PutCacheTask(id uint32, id2 uint32, taskVar map[string]interface{}) error
	DelCacheTask(id uint32, id2 uint32) (err error)
}

// --------------------------------------------------------------------

// validate : taskId 가 파라메터에 있는지 체크
// 인풋 : taskId
// 리턴 : error
func (t *TaskManager)validate() error {
	_, ok := taskDatas[t.taskId]
	if !ok {
		return fmt.Errorf("taskId : %d 는 taskDatas에 없음.", t.taskId)
	}
	return nil
}

// CreateTask : 인터벌 시작. 시작수량과 최대수량이 같거나 크면 타임체크를 시작하지 않음.
// 인풋 : none
// 리턴 : 현재수량, 인터벌, remainTime( 0이면 스톱상태 ), err
func (t *TaskManager)CreateTask() (curNum int32, interval int32, remainTime int32, err error) {
	if err := t.validate(); err != nil {
		return 0, 0, 0, err
	}
	// 메모리에서 태스크 관련 파라메터 get
	taskd := taskDatas[t.taskId]

	// write 할 내용 편집
	checkTime := time.Now().Unix()

	// 디비에 씀
	m := make(map[string]interface{})
	m["ct"] = checkTime
	m["num"] = taskd.startNum
	err = t.io.Write(t.uid, t.taskId, m)
	if err != nil {
		return 0, 0, 0, err
	}
	// 캐쉬에 씀
	err = t.cacheIo.PutCacheTask(t.uid, t.taskId, m)
	if err != nil {
		return 0, 0, 0, err
	}

	return taskd.startNum, taskd.interval, taskd.interval, nil
}

// CalcTask : 수량을 더하고 뺌 (하트 사용, 하트 선물등. 시간에 따른 수량 변화는 update 함수로 처리)
// 인풋 : 추가/감소 수량
// 리턴 : 현재수량, 인터벌, remainTime( 0이면 스톱상태 ), err
func (t *TaskManager)CalcTask(num int32) (curNum int32, interval int32, remainTime int32, err error) {
	if err := t.validate(); err != nil {
		return 0, 0, 0, err
	}
	// 메모리에서 태스크 관련 파라메터 get
	taskd := taskDatas[t.taskId]

	// 캐쉬에서 user 데이터 get
	userTask, err := t.cacheIo.GetCacheTask(t.uid, t.taskId)
	if err != nil {
		return 0, 0, 0, err
	}

	// 현재 시간을 기준으로 업데이트 실시
	newNum, newRemainTime, NewCheckTime, err := t.update(userTask)
	addedNum := newNum + num
	// 계산된 수량이 음수이면 0으로 초기화, 그리고 체크타임등을 초기화.
	// 수량이 Max 치를 넘어서는 것에 대해서는 제한하지 않음. (update 함수내부에서는 제한함)
	if addedNum < 0 {
		addedNum = 0
		newRemainTime = taskd.interval
		NewCheckTime = time.Now().Unix()
	}

	// 최대치 이상이었다가 최대치보다 작아 졌을때는 체크타임과 remainTime 을 초기화해준다.
	// 예> 풀 하트였다가 하트 소모하면 그때부터 다시 재계산함.
	if newNum >= taskd.maxNum && addedNum < taskd.maxNum {
		newRemainTime = taskd.interval
		NewCheckTime = time.Now().Unix()
	}

	// 디비에 씀
	m := make(map[string]interface{})
	m["ct"] = NewCheckTime
	m["num"] = addedNum
	err = t.io.Write(t.uid, t.taskId, m)
	if err != nil {
		return 0, 0, 0, err
	}
	// 캐쉬에 씀
	err = t.cacheIo.PutCacheTask(t.uid, t.taskId, m)
	if err != nil {
		return 0, 0, 0, err
	}

	return addedNum, taskd.interval, newRemainTime, nil
}

// DeleteTask : 태스크 삭제
// 인풋 : user id, task index
// 리턴 : err
func (t *TaskManager)DeleteTask() (err error) {
	if err := t.validate(); err != nil {
		return err
	}

	// 디비에서 제거
	err = t.io.Del(t.uid, t.taskId)
	if err != nil {
		return err
	}
	// 캐쉬에서 제거
	err = t.cacheIo.DelCacheTask(t.uid, t.taskId)
	if err != nil {
		return err
	}

	return nil
}

// calcNum : 체크시간을 기준으로 현재 수량과 남은 시간, 필요하다면(현재수량이 증가 했을때) 체크시간 업데이트
// 인풋 : user id, task index
// 리턴 : err
func (t *TaskManager)update(userTask map[string]interface{}) (newNum int32, newRemainTime int32, newCheckTime int64, err error) {
	if err := t.validate(); err != nil {
		return 0,0,0,err
	}

	taskd := taskDatas[t.taskId]
	curTime := time.Now().Unix()

	// 총수량이 최대수량보다 많으면 더 볼것도 없이,
	// 총수량은 최대수량으로 고정시키고 남은 시간과 체크시간을 현재 시간 기준으로 바꿈.
	curNum := userTask["num"].(int32)
	if curNum >= taskd.maxNum {
		return curNum, taskd.interval, curTime, nil
	}

	// 체크시간과 현재시간의 차이에서 인터벌로 나눈 수만큼 갯수를 증가시킴.
	oldCheckTime := userTask["ct"].(int64)
	curInterval := curTime - oldCheckTime
	curInterval = Max64(0, curInterval)
	if taskd.interval == 0 {
		return 0,0,0,fmt.Errorf("taskData.interval(%d) is 0", t.taskId)
	}
	portion := curInterval / int64(taskd.interval)
	mod := curInterval % int64(taskd.interval)

	// 시간 계산후에 총수량이 최대수량보다 많아버리면,
	// 총수량은 최대수량으로 고정시키고 남은 시간과 체크시간을 현재 시간 기준으로 바꿈.
	rNum := curNum + int32(portion)
	rNum = Min(rNum, taskd.maxNum)

	// 새 체크시간 갱신
	rCheckTime := curTime - mod

	return rNum, int32(mod), rCheckTime, nil
}
