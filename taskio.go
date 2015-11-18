package taskbytime

type taskio struct {

}

func NewTaskIO() *taskio {
	return &taskio{

	}
}

type taskIO interface {
	ReadUserTask(uid string, tid string) (map[string]interface{}, error)
	WriteUserTask(uid string, tid string, updateAttrs map[string]interface{}) error
	DelUserTask(uid string, tid string) error
}


