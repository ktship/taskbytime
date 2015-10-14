package taskbytime
import (
	"fmt"
	"time"
)

//타임기반의 일처리 엔진

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
	Read(id uint32, id2 uint32) (taskVar TaskVariable, err error)
	Write(id uint32, id2 uint32, taskVar TaskVariable) error
}

// 캐쉬에서 태스크 관련 정보들을 읽고 쓰는 인터페이스
type TaskCacheIO interface {
	GetCacheTask(id uint32, id2 uint32) (taskVar TaskVariable, err error)
	PutCacheTask(id uint32, id2 uint32, taskVar TaskVariable) error
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

// StartTask : 인터벌 시작. 시작수량과 최대수량이 같거나 크면 타임체크를 시작하지 않음.
// 인풋 : none
// 리턴 : 현재수량, 인터벌, remainTime( 0이면 스톱상태 ), err
func (t *TaskManager)StartTask() (curNum int32, interval int32, remainTime int32, err error) {
	if err := t.validate(); err != nil {
		return 0, 0, 0, err
	}

	// write 할 내용 편집
	taskData := taskDatas[t.taskId]
	// 시작수량이 최대수량보다 작을때만 체크타임
	var checkTime int64
	checkTime = 0
	if taskData.startNum < taskData.maxNum {
		checkTime = time.Now().Unix()
	}
	taskVar := TaskVariable {
		checkTime:	checkTime,
		curNum:		taskData.startNum,
	}

	// 디비에 씀
	err = t.io.Write(t.uid, t.taskId, taskVar)
	if err != nil {
		return 0, 0, 0, err
	}

	return taskData.startNum, taskData.interval, taskData.interval, nil
}

// AddTask : 수량을 늘림 (친구의 하트선물등)
// 인풋 : add number
// 리턴 : 현재수량, 인터벌, remainTime( 0이면 스톱상태 ), err
func (t *TaskManager)AddTask(addNum int32) (curNum int32, interval int32, remainTime int32, err error) {
	if err := t.validate(); err != nil {
		return 0, 0, 0, err
	}

	taskData := taskDatas[t.taskId]

	// user 데이터 get
	userTask, err := t.cacheIo.GetCacheTask(t.uid, t.taskId)
	if err != nil {
		return 0, 0, 0, err
	}

	// 현재 시간을 기준으로 업데이트 실시
	t.update(&userTask)

	// 수량 추가: 수량은 음수가 될 수 없음
	addedNum := userTask.curNum + addNum
	addedNum = Max(addedNum, 0)

	taskVar := TaskVariable {
		checkTime:	userTask.checkTime,
		curNum:		addedNum,
	}

	// 디비에 씀
	err = t.io.Write(t.uid, t.taskId, taskVar)
	if err != nil {
		return 0, 0, 0, err
	}

	return taskData.startNum, taskData.interval, taskData.interval, nil
}

// ReduceTask : 수량을 줄임. 자동 삭제(옵션)
// 인풋 : reduce num
// 리턴 : 현재수량, 인터벌, remainTime( 0이면 스톱상태 ), err
func (t *TaskManager)ReduceTask(reduceNum int32) (curNum int32, interval int32, remainTime int32, err error) {
	if err := t.validate(); err != nil {
		return 0, 0, 0, err
	}


	return 0, 0, 0, nil
}

// DeleteTask : 태스크 삭제
// 인풋 : user id, task index
// 리턴 : err
func (t *TaskManager)DeleteTask() (err error) {
	if err := t.validate(); err != nil {
		return err
	}


	return nil
}

// calcNum : 체크시간을 기준으로 현재 수량과 남은 시간, 필요하다면(현재수량이 증가 했을때) 체크시간 업데이트
// 인풋 : user id, task index
// 리턴 : err
func (t *TaskManager)update(userTask *TaskVariable) (newNum int32, newRemainTime int32, newCheckTime int32, err error) {
	if err := t.validate(); err != nil {
		return 0,0,0,err
	}
	taskData := taskDatas[t.taskId]

	oldCheckTime := userTask.checkTime
	curTime := time.Now().Unix()
	curInterval := curTime - oldCheckTime
	curInterval = Max64(0, curInterval)
	if taskData.interval == 0 {
		return 0,0,0,fmt.Errorf("taskData.interval is 0")
	}
	portion := curInterval / int64(taskData.interval)
	mod := curInterval % int64(taskData.interval)

	rNum := userTask.curNum + int32(portion)
	rCheckTime := curTime - mod

	return rNum, int32(mod), int32(rCheckTime), nil
}
