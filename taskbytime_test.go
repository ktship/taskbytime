package taskbytime
import (
	"testing"
	"time"
	"log"
	"fmt"
	"strconv"
	"runtime"
	"github.com/ktship/testio"
)

const TEST_KEY_USER = "test_user"
const TEST_KEY_TASK = "test_task"

func init() {
	fmt.Printf("Running On %s, %s, %s, %d-bit \n", runtime.Compiler, runtime.GOARCH, runtime.GOOS, strconv.IntSize)

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

type client struct {
	taskData
	uid				int
	tasks			map[int]*clientTask
}

type clientTask struct {
	curNum			int
	remainedTime	int
}

func newClient(t taskData, uid int) *client {
	return &client{
		taskData: 		t,
		uid:			uid,
		tasks: 			make(map[int]*clientTask),
	}
}

func (c *client)setTaskInfo(tid int, cNum int, rTime int) {
	c.tasks[tid] = &clientTask{
		curNum		:cNum,
		remainedTime:rTime,
	}
}

func (c *client)getTaskInfo(tid int) (curNum int, remainedTime int) {
	info := c.tasks[tid]
	return info.curNum, info.remainedTime
}

func (c *client)tick1sec() {
	for k, v := range c.tasks {
		taskd := taskDatas[k]
		// 꽉 차 있으면 패쓰
		if (v.curNum >= taskd.maxNum) {
			continue
		}

		// 시간이 다 되었으면 수량 1 증가
		v.remainedTime = v.remainedTime - 1
		if v.remainedTime <= 0 {
			v.curNum = Min(taskd.maxNum, v.curNum + 1)
			v.remainedTime = taskd.interval
		}
	}
}

func Test01_New(t *testing.T) {
	log.Printf(" ---- Test01_New ")
	testio := testio.New()

	// users
	var uid1 int
	uid1 = 111

	// a task
	var taskId int
	taskId = 0

	tData := taskDatas[taskId]

	client := newClient(tData, uid1)
	nTM1 := New(testio)
	curNum, interval, remainedTime, err := nTM1.CreateTask(uid1, taskId)
	client.setTaskInfo(taskId, curNum, remainedTime)
	if err != nil {
		t.Errorf("Fail CreateTask %d, %d, %d", curNum, interval, remainedTime)
	}
	if interval == 0 {
		t.Errorf("Fail interval is 0")
	}
	if remainedTime == 0 {
		t.Errorf("Fail remainedTime is 0")
	}

	ticker := time.NewTicker(1 * time.Second)
	stop := make(chan bool, 1)

	go func(){
		for {
			select {
			case <- ticker.C:
				client.tick1sec()
				retNum, retRemainTime := client.getTaskInfo(taskId)
				log.Printf("client time Num:(%d) rTime: (%d)", retNum, retRemainTime)
				if retNum == 0 && retRemainTime == 1 {
					nTM := New(testio)
					log.Printf("server check nTM.CalcTask Task:%d Num:%d rTime: %d", taskId, retNum, retRemainTime)
					curNum, interval, remainedTime, err := nTM.CalcTask(uid1, taskId, 0)
					if curNum != 0 || interval != 5 || err != nil {
						t.Errorf("Fail CalcTask(0) %d, %d, %d", curNum, interval, remainedTime)
					}
				}

				if retNum == 1 && retRemainTime == 5 {
					nTM := New(testio)
					log.Printf("server check nTM.CalcTask Task:%d Num:%d rTime: %d", taskId, retNum, retRemainTime)
					curNum, interval, remainedTime, err := nTM.CalcTask(uid1, taskId, 0)
					if curNum != 1 || interval != 5 || err != nil {
						t.Errorf("Fail CalcTask(0) %d, %d, %d", curNum, interval, remainedTime)
					}
				}

			case <- stop:
				retNum, retRemainTime := client.getTaskInfo(taskId)
				if retNum == 3 && retRemainTime == 5 {
					nTM := New(testio)
					log.Printf("server check nTM.CalcTask Task:%d Num:%d rTime: %d", taskId, retNum, retRemainTime)
					curNum, interval, remainedTime, err := nTM.CalcTask(uid1, taskId, 0)
					if curNum != 3 || interval != 5 || err != nil {
						t.Errorf("Fail CalcTask(0) %d, %d, %d", curNum, interval, remainedTime)
					}
				}
				return
			}
		}
	}()

	// 9초후에 1개가 되고 남은 시간이 1초인가???
	time.Sleep(9 * time.Second)
	log.Printf(" ---- 9 sec later ")
	if true {
		taskm2 := New(testio)
		curNum, _, remainedTime, err := taskm2.CalcTask(uid1, taskId, 0)
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
		taskm3 := New(testio)
		curNum, _, remainedTime, err := taskm3.CalcTask(uid1, taskId, 0)
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

func Test02_Calc(t *testing.T) {
	// users
	var uid1 int
	uid1 = 111

	// a task : 4, 4, 3
	var taskId int
	taskId = 1

	testio := testio.New()
	nTM := New(testio)
	curNum, interval, remainedTime, err := nTM.CreateTask(uid1, taskId)
	if err != nil {
		t.Errorf("Fail CreateTask %d, %d, %d", curNum, interval, remainedTime)
	}
	log.Printf("TestCalc ----------- start(%d), interval(%d), remainTime(%d)", curNum, interval, remainedTime)

	time.Sleep(2 * time.Second)
	if true {
		log.Printf(" - 하나 사용(-1)하고 task 시작")
		nTM := New(testio)
		curNum, interval, remainedTime, err := nTM.CalcTask(uid1, taskId, -1)
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
		nTM2 := New(testio)
		curNum, interval, remainedTime, err = nTM2.CalcTask(uid1, taskId, -1)
		if err != nil {
			t.Errorf("Fail CreateTask %d, %d, %d", curNum, interval, remainedTime)
		}
		if curNum != 2 || remainedTime != 2 {
			t.Errorf("Fail CreateTask %d, %d, %d", curNum, interval, remainedTime)
		}

		time.Sleep(3 * time.Second)
		log.Printf(" 3초후")
		log.Printf(" -- TestCalc : 하나 사용(-1)")
		nTM3 := New(testio)
		curNum, interval, remainedTime, err = nTM3.CalcTask(uid1, taskId, -1)
		if err != nil {
			t.Errorf("Fail CreateTask %d, %d, %d", curNum, interval, remainedTime)
		}
		if curNum != 2 || remainedTime != 2 {
			t.Errorf("Fail CreateTask %d, %d, %d", curNum, interval, remainedTime)
		}

		time.Sleep(6 * time.Second)
		log.Printf(" 6초후")
		log.Printf(" -- TestCalc : 하나 사용(-1)")
		nTM = New(testio)
		curNum, interval, remainedTime, err = nTM.CalcTask(uid1, taskId, -1)
		if err != nil {
			t.Errorf("Fail CreateTask %d, %d, %d", curNum, interval, remainedTime)
		}
		if curNum != 3 || remainedTime != 3 {
			t.Errorf("Fail CreateTask %d, %d, %d", curNum, interval, remainedTime)
		}

		time.Sleep(2 * time.Second)
		log.Printf(" 2초후")
		log.Printf(" -- TestCalc : 3개 더함(3)")
		nTM = New(testio)
		curNum, interval, remainedTime, err = nTM.CalcTask(uid1, taskId, 3)
		if err != nil {
			t.Errorf("Fail CreateTask %d, %d, %d", curNum, interval, remainedTime)
		}
		if curNum != 6 || remainedTime != 1 {
			t.Errorf("Fail CreateTask %d, %d, %d", curNum, interval, remainedTime)
		}

		time.Sleep(2 * time.Second)
		log.Printf(" 23초후")
		log.Printf(" -- TestCalc : 하나 사용(-1)")
		nTM = New(testio)
		curNum, interval, remainedTime, err = nTM3.CalcTask(uid1, taskId, -1)
		if err != nil {
			t.Errorf("Fail CreateTask %d, %d, %d, %s", curNum, interval, remainedTime, err)
		}
		if curNum != 5 || remainedTime != 3 {
			t.Errorf("Fail CreateTask %d, %d, %d", curNum, interval, remainedTime)
		}

		nTM = New(testio)
		err = nTM3.DeleteTask(uid1, taskId)
		if err != nil {
			t.Errorf("Fail DeleteTask %d, %d, %d, %s", curNum, interval, remainedTime, err)
		}
	}
}

func Test03_Finish(t *testing.T) {
	// users
	var uid3 int
	uid3 = 333

	// a task : 1, 1, 5
	var taskId int
	taskId = 2

	testio := testio.New()
	nTM := New(testio)
	curNum, interval, remainedTime, err := nTM.CreateTask(uid3, taskId)
	if err != nil {
		t.Errorf("Fail CreateTask %d, %d, %d", curNum, interval, remainedTime)
	}
	log.Printf("TestFinish ----------- start(%d), interval(%d), remainTime(%d)", curNum, interval, remainedTime)

	time.Sleep(2 * time.Second)
	log.Printf(" 2초후")
	if true {
		log.Printf(" - 확인")
		nTM := New(testio)
		curNum, interval, remainedTime, err := nTM.CalcTask(uid3, taskId, 0)
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
		nTM2 := New(testio)
		curNum, interval, remainedTime, err = nTM2.CalcTask(uid3, taskId, 0)
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
		nTM3 := New(testio)
		curNum, interval, remainedTime, err = nTM3.CalcTask(uid3, taskId, 0)
		if err != nil {
			t.Errorf("Fail CreateTask %d, %d, %d", curNum, interval, remainedTime)
		}
		log.Printf(" - num 은 1, remainTime은 4이라야 함")
		if curNum != 1 || remainedTime != 4 {
			t.Errorf("Fail CreateTask %d, %d, %d", curNum, interval, remainedTime)
		}

		nTM = New(testio)
		err = nTM3.DeleteTask(uid3, taskId)
		if err != nil {
			t.Errorf("Fail DeleteTask %d, %d, %d, %s", curNum, interval, remainedTime, err)
		}

		time.Sleep(2 * time.Second)
		log.Printf(" 2초후")
		log.Printf(" - 확인")
		nTM = New(testio)
		curNum, interval, remainedTime, err = nTM.CalcTask(uid3, taskId, 0)
		if err == nil {
			t.Errorf("Fail CreateTask %d, %d, %d", curNum, interval, remainedTime)
		} else {
			fmt.Println(err)
		}
	}
}