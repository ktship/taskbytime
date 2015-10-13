package taskbytime
import (
	"fmt"
)

//타임기반의 일처리 엔진
type Task struct{
	startNum	int32
	maxNum		int32
	interval	int32
	isRepeat	bool
}

var taskData = make(map[int32]Task)

// SetDatas : 파라메터 추가
// 인풋 : taskId int32, startNum int32, interval int32, isRepeat bool
func SetData(taskId int32, t Task) {
	taskData[taskId] = t
}

// validate : taskId 가 파라메터에 있는지 체크
// 인풋 : taskId int32
// 리턴 : error
func validate(taskId int32) error {
	_, ok := taskData[taskId]
	if !ok {
		return fmt.Errorf("taskId : %d 는 taskDatas에 없음.", taskId)
	}
	return nil
}

// CreateTask : 태스크를 생성함
// 인풋 : user id, task index
// 리턴 : err
func CreateTask(uid int32, taskId int32) (err error) {
	if err := validate(taskId); err != nil {
		return err
	}

	return nil
}

// StartTask : 인터벌 시작. 현재수량이 최대수량보다 크거나 같을 경우에는 에러처리
// 인풋 : user id, task index
// 리턴 : 현재수량, 인터벌, remainTime( 0이면 스톱상태 ), err
func StartTask(uid int32, taskId int32) (curNum int32, interval int32, remainTime int32, err error) {
	if err := validate(taskId); err != nil {
		return 0, 0, 0, err
	}


	return 0, 0, 0, nil
}

// AddTask : 수량을 늘림 (친구의 하트선물등)
// 인풋 : user id, task index
// 리턴 : 현재수량, 인터벌, remainTime( 0이면 스톱상태 ), err
func AddTask(uid int32, taskId int32) (curNum int32, interval int32, remainTime int32, err error) {
	if err := validate(taskId); err != nil {
		return 0, 0, 0, err
	}

	return 0, 0, 0, nil
}

// ReduceTask : 수량을 줄임. 자동 삭제(옵션)
// 인풋 : user id, task index
// 리턴 : 현재수량, 인터벌, remainTime( 0이면 스톱상태 ), err
func ReduceTask(uid int32, taskId int32) (curNum int32, interval int32, remainTime int32, err error) {
	if err := validate(taskId); err != nil {
		return 0, 0, 0, err
	}


	return 0, 0, 0, nil
}

// DeleteTask : 태스크 삭제
// 인풋 : user id, task index
// 리턴 : err
func DeleteTask(uid int32, taskId int32) (err error) {
	if err := validate(taskId); err != nil {
		return err
	}


	return nil
}