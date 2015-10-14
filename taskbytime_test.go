package taskbytime
import (
	"testing"
	"time"
	"log"
)

func init() {
	log.Print("init taskbytime test")
	SetData(0,TaskData{
		startNum:0,
		maxNum:3,
		interval:10,
		isRepeat:true,
	})
	SetData(1,TaskData{
		startNum:5,
		maxNum:5,
		interval:13,
		isRepeat:true,
	})
	SetData(2,TaskData{
		startNum:0,
		maxNum:1,
		interval:17,
		isRepeat:false,
	})
}

type clientTask struct {
	TaskData
	taskId			int
	curNum			int
	remainedTime	int
	Ticker 			*time.Ticker
}

func newClient(t TaskData) *clientTask {
	return &clientTask{
		TaskData: t,
		Ticker			:time.NewTicker(1 * time.Second),
	}
}

func (c *clientTask)printSec() {
	log.Printf("Task: %d Current Num: %d Remained Time: %d", c.taskId, c.curNum, c.remainedTime)
}

func TestNew(t *testing.T) {
	// users
	var uid1 int
	uid1 = 123

	// a task
	var taskId_a int
	taskId_a = 0

	var curNum, interval, remainedTime int
	var err error
	if err = CreateTask(uid1, taskId_a); err == nil {
		t.Errorf("Fail CreateTask %s", err)
	}

	if curNum, interval, remainedTime, err = StartTask(uid1, taskId_a); err == nil {
		t.Errorf("Fail StartTask %d, %d, %d, %s", curNum, interval, remainedTime, err)
	}
	if interval == 0 {
		t.Errorf("Fail interval is 0")
	}
	if remainedTime == 0 {
		t.Errorf("Fail remainedTime is 0")
	}

	client := newClient(TaskData {
		startNum:	0,
		maxNum:		0,
		interval:	0,
		isRepeat:	true,
	})
	client.Ticker = time.NewTicker(1 * time.Second)
	go func(){
		for _ = range client.Ticker.C {
			client.printSec()
		}
	}()
	time.Sleep(100 * time.Second)

	t.Error("Fail")
}

