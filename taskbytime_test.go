package taskbytime
import (
	"testing"
	"time"
	"log"
)

func init() {
	log.Print("init taskbytime test")
	SetData(0, 0, 3, 10, true)
	SetData(1, 5, 5, 12, true)
	SetData(2, 0, 1, 15, false)
}

func TestNew(t *testing.T) {
	// users
	var uid1 int32
	uid1 = 123

	// a task
	var taskIndex_a int32
	taskIndex_a = 0

	if startNum, interval, err := CreateTask(uid1, taskIndex_a); err == nil {
		t.Errorf("Fail CreateTask %d, %d, %s", startNum, interval, err)
	}

	if curNum, interval, remainedTime, err := StartTask(uid1, taskIndex_a); err == nil {
		t.Errorf("Fail CreateTask %d, %d, %d, %s", curNum, interval, remainedTime, err)
	}

	ticker1Sec := time.NewTicker(1 * time.Second)
	go func(){
		for _ = range ticker1Sec.C {
			log.Println("tick")
		}
	}()
	time.Sleep(100 * time.Second)

	t.Error("Fail")
}

