package taskbytime

// 태스크 기본 구성 파라메터. 변하지 않는 데이터
type TaskData struct {
	startNum	int32
	maxNum		int32
	interval	int32
	isRepeat	bool
}

// TaskData 로 구성된 맵. 파라메터처럼 사용함.
var taskDatas = make(map[uint32]TaskData)

// SetDatas : 파라메터 추가
// 인풋 : taskId int, startNum int, interval int, isRepeat bool
func SetData(taskId uint32, t TaskData) {
	taskDatas[taskId] = t
}
