package taskbytime

type taskIO interface {
	ReadUserTask(uid int, tid int) (map[string]interface{}, error)
	WriteUserTask(uid int, tid int, updateAttrs map[string]interface{}) error
	DelUserTask(uid int, tid int) error
}


