package taskbytime
import (
	"testing"
	"time"
	"log"
	"github.com/ktship/testio"
)

func init() {
	log.Print("init taskbytime test")
	SetData(0, taskData{
		startNum:0,
		maxNum:3,
		interval:10,
		isRepeat:true,
	})
	SetData(1, taskData{
		startNum:5,
		maxNum:5,
		interval:13,
		isRepeat:true,
	})
	SetData(2, taskData{
		startNum:0,
		maxNum:1,
		interval:17,
		isRepeat:false,
	})
}

type clientTask struct {
	taskData
	taskId			int
	curNum			int
	remainedTime	int
	Ticker 			*time.Ticker
}

func newClient(t taskData) *clientTask {
	return &clientTask{
		taskData: t,
		Ticker			:time.NewTicker(1 * time.Second),
	}
}

func (c *clientTask)printSec() {
	log.Printf("Task: %d Current Num: %d Remained Time: %d", c.taskId, c.curNum, c.remainedTime)
}

func TestNew(t *testing.T) {
	// users
	var uid1 uint32
	uid1 = 123

	// a task
	var taskId_a uint32
	taskId_a = 0

	tio := testio.NewTestFileIO()
	cio := testio.NewTestCacheIO()
	taskm := NewTaskManager(tio, cio, uid1, taskId_a)

	curNum, interval, remainedTime, err := taskm.CreateTask()
	if err == nil {
		t.Errorf("Fail CreateTask %d, %d, %d, %s", curNum, interval, remainedTime, err)
	}

	if interval == 0 {
		t.Errorf("Fail interval is 0")
	}
	if remainedTime == 0 {
		t.Errorf("Fail remainedTime is 0")
	}

	client := newClient(taskData{
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

