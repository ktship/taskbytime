package taskbytime

// 태스크 기본 구성 파라메터. 변하지 않는 데이터
type taskData struct {
	startNum	int
	maxNum		int
	interval	int
	isRepeat	bool
}

// TaskData 로 구성된 맵. 파라메터처럼 사용함.
var taskDatas = make(map[int]taskData)

// SetDatas : 파라메터 추가
// 인풋 : taskId int, startNum int, interval int, isRepeat bool
func SetData(taskId int, t taskData) {
	taskDatas[taskId] = t
}
