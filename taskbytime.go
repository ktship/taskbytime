package taskbytime
import (
	"fmt"
)

//타임기반의 일처리 엔진
type Task struct{
	startNum	uint32
	maxNum		uint32
	interval	uint32
	isRepeat	bool
}

type TaskIO interface {
	Read(id uint32, id2 uint32) map[string]interface{}
	Write(id uint32, id2 uint32, data map[string]interface{})
}

type TaskManager struct {
	io TaskIO
}

func NewTaskManger(taskio TaskIO) *TaskManager {
	return &TaskManager{ io:taskio }
}

var taskData = make(map[uint32]Task)

// SetDatas : 파라메터 추가
// 인풋 : taskId uint32, startNum uint32, interval uint32, isRepeat bool
func SetData(taskId uint32, t Task) {
	taskData[taskId] = t
}

// validate : taskId 가 파라메터에 있는지 체크
// 인풋 : taskId uint32
// 리턴 : error
func (t *TaskManager)validate(taskId uint32) error {
	_, ok := taskData[taskId]
	if !ok {
		return fmt.Errorf("taskId : %d 는 taskDatas에 없음.", taskId)
	}
	return nil
}

// CreateTask : 태스크를 생성함
// 인풋 : user id, task index
// 리턴 : err
func (t *TaskManager)CreateTask(uid uint32, taskId uint32) (err error) {
	if err := t.validate(taskId); err != nil {
		return err
	}

	t.io.Write(uid, taskId, nil)

	return nil
}

// StartTask : 인터벌 시작. 현재수량이 최대수량보다 크거나 같을 경우에는 에러처리
// 인풋 : user id, task index
// 리턴 : 현재수량, 인터벌, remainTime( 0이면 스톱상태 ), err
func (t *TaskManager)StartTask(uid uint32, taskId uint32) (curNum uint32, interval uint32, remainTime uint32, err error) {
	if err := t.validate(taskId); err != nil {
		return 0, 0, 0, err
	}


	return 0, 0, 0, nil
}

// AddTask : 수량을 늘림 (친구의 하트선물등)
// 인풋 : user id, task index
// 리턴 : 현재수량, 인터벌, remainTime( 0이면 스톱상태 ), err
func (t *TaskManager)AddTask(uid uint32, taskId uint32) (curNum uint32, interval uint32, remainTime uint32, err error) {
	if err := t.validate(taskId); err != nil {
		return 0, 0, 0, err
	}

	return 0, 0, 0, nil
}

// ReduceTask : 수량을 줄임. 자동 삭제(옵션)
// 인풋 : user id, task index
// 리턴 : 현재수량, 인터벌, remainTime( 0이면 스톱상태 ), err
func (t *TaskManager)ReduceTask(uid uint32, taskId uint32) (curNum uint32, interval uint32, remainTime uint32, err error) {
	if err := t.validate(taskId); err != nil {
		return 0, 0, 0, err
	}


	return 0, 0, 0, nil
}

// DeleteTask : 태스크 삭제
// 인풋 : user id, task index
// 리턴 : err
func (t *TaskManager)DeleteTask(uid uint32, taskId uint32) (err error) {
	if err := t.validate(taskId); err != nil {
		return err
	}


	return nil
}