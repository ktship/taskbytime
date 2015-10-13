package taskbytime
import (
	"testing"
	"time"
	"log"
)

func init() {
	log.Print("init taskbytime test")
	SetData(0,Task{
		startNum:0,
		maxNum:3,
		interval:10,
		isRepeat:true,
	})
	SetData(1,Task{
		startNum:5,
		maxNum:5,
		interval:13,
		isRepeat:true,
	})
	SetData(2,Task{
		startNum:0,
		maxNum:1,
		interval:17,
		isRepeat:false,
	})
}

type clientTask struct {
	Task
	taskId			int32
	curNum			int32
	remainedTime	int32
	Ticker 			*time.Ticker
}

func newClient(t Task) *clientTask {
	return &clientTask{
		Task: t,
		Ticker			:time.NewTicker(1 * time.Second),
	}
}

func (c *clientTask)printSec() {
	log.Printf("Task: %d Current Num: %d Remained Time: %d", c.taskId, c.curNum, c.remainedTime)
}

func TestNew(t *testing.T) {
	// users
	var uid1 int32
	uid1 = 123

	// a task
	var taskId_a int32
	taskId_a = 0

	var curNum, interval, remainedTime int32
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

	client := newClient(Task {
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

