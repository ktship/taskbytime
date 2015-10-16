package taskbytime
import (
	"testing"
	"time"
	"log"
	"github.com/ktship/testio"
	"fmt"
)

func init() {
	SetData(0, taskData{
		startNum:0,
		maxNum:5,
		interval:5,
		isRepeat:true,
	})
	SetData(1, taskData{
		startNum:4,
		maxNum:4,
		interval:3,
		isRepeat:true,
	})
	SetData(2, taskData{
		startNum:0,
		maxNum:1,
		interval:5,
		isRepeat:false,
	})
}

type clientTask struct {
	taskData
	uid				uint32
	taskId			uint32
	curNum			int32
	remainedTime	int32
	Ticker 			*time.Ticker
}

func newClient(t taskData, uid uint32, tid uint32, cNum int32, rTime int32) *clientTask {
	return &clientTask{
		taskData: 		t,
		uid:			uid,
		taskId: 		tid,
		curNum: 		cNum,
		remainedTime: 	rTime,
		Ticker:			time.NewTicker(1 * time.Second),
	}
}

func (c *clientTask)newTM() *TaskManager {
	tio := testio.NewTestFileIO()
	cio := testio.NewTestCacheIO()
	taskm := NewTaskManager(tio, cio, c.uid, c.taskId)
	return taskm
}

func (c *clientTask)getClientInfo() (taskId uint32, curNum int32, remainedTime int32) {
	return c.taskId, c.curNum, c.remainedTime
}

func (c *clientTask)tick1sec() {
	taskd := taskDatas[c.taskId]

	// 꽉 차 있으면 패쓰
	if (c.curNum >= taskd.maxNum) {
		return
	}

	// 시간이 다 되었으면 수량 1 증가
	c.remainedTime = c.remainedTime - 1
	if c.remainedTime <= 0 {
		c.curNum = Min(taskd.maxNum, c.curNum + 1)
		c.remainedTime = taskd.interval
	}
}

func TestNew(t *testing.T) {
	// users
	var uid1, uid2, uid3 uint32
	uid1 = 111
	uid2 = 222
	uid3 = 333

	// a task
	var taskId_a uint32
	taskId_a = 0

	tio := testio.NewTestFileIO()
	cio := testio.NewTestCacheIO()

	tData := taskDatas[taskId_a]

	client := newClient(tData, uid1, taskId_a, tData.startNum, tData.interval)
	nTM1 := client.newTM()
	curNum, interval, remainedTime, err := nTM1.CreateTask()
	if err != nil {
		t.Errorf("Fail CreateTask %d, %d, %d", curNum, interval, remainedTime)
	}
	if interval == 0 {
		t.Errorf("Fail interval is 0")
	}
	if remainedTime == 0 {
		t.Errorf("Fail remainedTime is 0")
	}

	client2 := newClient(tData, uid2, taskId_a, tData.startNum, tData.interval)
	nTM2 := client2.newTM()
	curNum, interval, remainedTime, err = nTM2.CreateTask()
	client3 := newClient(tData, uid3, taskId_a, tData.startNum, tData.interval)
	nTM3 := client3.newTM()
	curNum, interval, remainedTime, err = nTM3.CreateTask()

	client.Ticker = time.NewTicker(1 * time.Second)
	stop := make(chan bool, 1)

	go func(){
		for {
			select {
			case <- client.Ticker.C:
				client.tick1sec()
				retTaskId, retNum, retRemainTime := client.getClientInfo()
				log.Printf("client time Num:(%d) rTime: (%d)", retNum, retRemainTime)
				if retNum == 0 && retRemainTime == 1 {
					nTM := client.newTM()
					log.Printf("server check nTM.CalcTask Task:%d Num:%d rTime: %d", retTaskId, retNum, retRemainTime)
					curNum, interval, remainedTime, err := nTM.CalcTask(0)
					if curNum != 0 || interval != 5 || err != nil {
						t.Errorf("Fail CalcTask(0) %d, %d, %d", curNum, interval, remainedTime)
					}
				}

				if retNum == 1 && retRemainTime == 5 {
					nTM := client.newTM()
					log.Printf("server check nTM.CalcTask Task:%d Num:%d rTime: %d", retTaskId, retNum, retRemainTime)
					curNum, interval, remainedTime, err := nTM.CalcTask(0)
					if curNum != 1 || interval != 5 || err != nil {
						t.Errorf("Fail CalcTask(0) %d, %d, %d", curNum, interval, remainedTime)
					}
				}

				if retNum == 3 && retRemainTime == 5 {
					nTM := client.newTM()
					log.Printf("server check nTM.CalcTask Task:%d Num:%d rTime: %d", retTaskId, retNum, retRemainTime)
					curNum, interval, remainedTime, err := nTM.CalcTask(0)
					if curNum != 3 || interval != 5 || err != nil {
						t.Errorf("Fail CalcTask(0) %d, %d, %d", curNum, interval, remainedTime)
					}
				}
			case <- stop:
				return
			}
		}
	}()

	// 9초후에 1개가 되고 남은 시간이 1초인가???
	time.Sleep(9 * time.Second)
	log.Printf(" ---- 9 sec later ")
	if true {
		taskm2 := NewTaskManager(tio, cio, uid2, taskId_a)
		curNum, _, remainedTime, err := taskm2.CalcTask(0)
		if err != nil {
			t.Error("taskm2.CalcTask")
		}
		log.Printf("--- 2 client curNum:%d remainedTime: %d", curNum, remainedTime)
		if curNum != 1 || remainedTime != 1 || interval != 5 || err != nil {
			t.Errorf("Fail CalcTask(0) %d, %d, %d", curNum, interval, remainedTime)
		}
	}
	time.Sleep(1 * time.Second)

	// 20초 후에 4개가 정확히 되는가???
	time.Sleep(10 * time.Second)
	log.Printf(" ---- 20 sec later ")
	if true {
		taskm3 := NewTaskManager(tio, cio, uid3, taskId_a)
		curNum, _, remainedTime, err := taskm3.CalcTask(0)
		if err != nil {
			t.Error("taskm3.CalcTask")
		}
		log.Printf("--- 3 client curNum:%d remainedTime: %d", curNum, remainedTime)
		if curNum != 4 || remainedTime != 5 || interval != 5 || err != nil {
			t.Errorf("Fail CalcTask(0) %d, %d, %d", curNum, interval, remainedTime)
		}
	}
	time.Sleep(1 * time.Second)
	stop <- true
}

func TestCalc(t *testing.T) {
	// users
	var uid1 uint32
	uid1 = 111

	// a task : 4, 4, 3
	var taskId uint32
	taskId = 1

	tData := taskDatas[taskId]

	client := newClient(tData, uid1, taskId, tData.startNum, tData.interval)
	nTM := client.newTM()
	curNum, interval, remainedTime, err := nTM.CreateTask()
	if err != nil {
		t.Errorf("Fail CreateTask %d, %d, %d", curNum, interval, remainedTime)
	}
	log.Printf("TestCalc ----------- start(%d), interval(%d), remainTime(%d)", curNum, interval, remainedTime)

	time.Sleep(2 * time.Second)
	if true {
		log.Printf(" - 하나 사용(-1)하고 task 시작")
		nTM := client.newTM()
		curNum, interval, remainedTime, err := nTM.CalcTask(-1)
		if err != nil {
			t.Errorf("Fail CreateTask %d, %d, %d", curNum, interval, remainedTime)
		}
		log.Printf(" - num 은 3, remainTime은 3이라야 함")
		if curNum != 3 || remainedTime != 3 {
			t.Errorf("Fail CreateTask %d, %d, %d", curNum, interval, remainedTime)
		}

		time.Sleep(1 * time.Second)
		log.Printf(" 1초후")
		log.Printf(" -- TestCalc : 하나 사용(-1)")
		nTM2 := client.newTM()
		curNum, interval, remainedTime, err = nTM2.CalcTask(-1)
		if err != nil {
			t.Errorf("Fail CreateTask %d, %d, %d", curNum, interval, remainedTime)
		}
		if curNum != 2 || remainedTime != 2 {
			t.Errorf("Fail CreateTask %d, %d, %d", curNum, interval, remainedTime)
		}

		time.Sleep(3 * time.Second)
		log.Printf(" 3초후")
		log.Printf(" -- TestCalc : 하나 사용(-1)")
		nTM3 := client.newTM()
		curNum, interval, remainedTime, err = nTM3.CalcTask(-1)
		if err != nil {
			t.Errorf("Fail CreateTask %d, %d, %d", curNum, interval, remainedTime)
		}
		if curNum != 2 || remainedTime != 2 {
			t.Errorf("Fail CreateTask %d, %d, %d", curNum, interval, remainedTime)
		}

		time.Sleep(6 * time.Second)
		log.Printf(" 6초후")
		log.Printf(" -- TestCalc : 하나 사용(-1)")
		nTM = client.newTM()
		curNum, interval, remainedTime, err = nTM.CalcTask(-1)
		if err != nil {
			t.Errorf("Fail CreateTask %d, %d, %d", curNum, interval, remainedTime)
		}
		if curNum != 3 || remainedTime != 3 {
			t.Errorf("Fail CreateTask %d, %d, %d", curNum, interval, remainedTime)
		}

		time.Sleep(2 * time.Second)
		log.Printf(" 2초후")
		log.Printf(" -- TestCalc : 세게 더함(3)")
		nTM = client.newTM()
		curNum, interval, remainedTime, err = nTM.CalcTask(3)
		if err != nil {
			t.Errorf("Fail CreateTask %d, %d, %d", curNum, interval, remainedTime)
		}
		if curNum != 6 || remainedTime != 1 {
			t.Errorf("Fail CreateTask %d, %d, %d", curNum, interval, remainedTime)
		}

		time.Sleep(2 * time.Second)
		log.Printf(" 23초후")
		log.Printf(" -- TestCalc : 하나 사용(-1)")
		nTM = client.newTM()
		curNum, interval, remainedTime, err = nTM3.CalcTask(-1)
		if err != nil {
			t.Errorf("Fail CreateTask %d, %d, %d, %s", curNum, interval, remainedTime, err)
		}
		if curNum != 5 || remainedTime != 3 {
			t.Errorf("Fail CreateTask %d, %d, %d", curNum, interval, remainedTime)
		}

		nTM = client.newTM()
		err = nTM3.DeleteTask()
		if err != nil {
			t.Errorf("Fail DeleteTask %d, %d, %d, %s", curNum, interval, remainedTime, err)
		}
	}
}

func TestFinish(t *testing.T) {
	// users
	var uid3 uint32
	uid3 = 333

	// a task : 1, 1, 5
	var taskId uint32
	taskId = 2

	tData := taskDatas[taskId]

	client := newClient(tData, uid3, taskId, tData.startNum, tData.interval)
	nTM := client.newTM()
	curNum, interval, remainedTime, err := nTM.CreateTask()
	if err != nil {
		t.Errorf("Fail CreateTask %d, %d, %d", curNum, interval, remainedTime)
	}
	log.Printf("TestFinish ----------- start(%d), interval(%d), remainTime(%d)", curNum, interval, remainedTime)

	time.Sleep(2 * time.Second)
	log.Printf(" 2초후")
	if true {
		log.Printf(" - 확인")
		nTM := client.newTM()
		curNum, interval, remainedTime, err := nTM.CalcTask(0)
		if err != nil {
			t.Errorf("Fail CreateTask %d, %d, %d", curNum, interval, remainedTime)
		}
		log.Printf(" - num 은 0, remainTime은 3이라야 함")
		if curNum != 0 || remainedTime != 3 {
			t.Errorf("Fail CreateTask %d, %d, %d", curNum, interval, remainedTime)
		}

		time.Sleep(2 * time.Second)
		log.Printf(" 2초후")
		log.Printf(" - 확인")
		nTM2 := client.newTM()
		curNum, interval, remainedTime, err = nTM2.CalcTask(0)
		if err != nil {
			t.Errorf("Fail CreateTask %d, %d, %d", curNum, interval, remainedTime)
		}
		log.Printf(" - num 은 0, remainTime은 1이라야 함")
		if curNum != 0 || remainedTime != 1 {
			t.Errorf("Fail CreateTask %d, %d, %d", curNum, interval, remainedTime)
		}

		time.Sleep(2 * time.Second)
		log.Printf(" 2초후")
		log.Printf(" - 확인")
		nTM3 := client.newTM()
		curNum, interval, remainedTime, err = nTM3.CalcTask(0)
		if err != nil {
			t.Errorf("Fail CreateTask %d, %d, %d", curNum, interval, remainedTime)
		}
		log.Printf(" - num 은 1, remainTime은 4이라야 함")
		if curNum != 1 || remainedTime != 4 {
			t.Errorf("Fail CreateTask %d, %d, %d", curNum, interval, remainedTime)
		}

		nTM = client.newTM()
		err = nTM3.DeleteTask()
		if err != nil {
			t.Errorf("Fail DeleteTask %d, %d, %d, %s", curNum, interval, remainedTime, err)
		}

		time.Sleep(2 * time.Second)
		log.Printf(" 2초후")
		log.Printf(" - 확인")
		nTM = client.newTM()
		curNum, interval, remainedTime, err = nTM.CalcTask(0)
		if err == nil {
			t.Errorf("Fail CreateTask %d, %d, %d", curNum, interval, remainedTime)
		} else {
			fmt.Println(err)
		}
	}
}