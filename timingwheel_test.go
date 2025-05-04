package controllabletimingwheel_test

import (
	"fmt"
	"time"

	Timingwheel "github.com/eternal-dragon/controllableTimingwheel"
)

func Example_ControllableTimingWheel() {
	// 创建层级时间轮
	hierarchicalWheel := Timingwheel.NewControllableTimingWheel(
		time.Millisecond*10,
		10,
	)

	// 添加一些定时任务
	hierarchicalWheel.AddTask(500*time.Millisecond, func() {
		fmt.Println("Task executed after 500ms")
	})

	hierarchicalWheel.AddTask(2*time.Second, func() {
		fmt.Println("Task executed after 2s")
	})

	// 启动层级时间轮
	hierarchicalWheel.Start()
}
