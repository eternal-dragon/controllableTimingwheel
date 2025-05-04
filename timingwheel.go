package controllabletimingwheel

import (
	"log"
	"sync"
	"time"
)

// 创建层级时间轮结构体
type ControllableTimingWheel struct {
	// 最小时间间隔
	minInterval time.Duration
	// 时间轮
	wheel *timingWheel
	// 时间缩放
	timeScale float64

	stopC   chan interface{}
	changeC chan interface{}
}

// 初始化层级时间轮
func NewControllableTimingWheel(minInterval time.Duration, slotCountPerLevel int) *ControllableTimingWheel {
	return &ControllableTimingWheel{
		minInterval: minInterval,
		wheel:       newTimingWheel(slotCountPerLevel, minInterval),
		timeScale:   1,
		stopC:       make(chan interface{}),
		changeC:     make(chan interface{}),
	}
}

// 添加定时任务
func (h *ControllableTimingWheel) AddTask(delay time.Duration, task func()) {
	// current slot的task会被直接执行，如果等待current触发，可能会迟到。
	if delay < h.wheel.interval {
		go task()
		return
	}
	h.wheel.addTask(delay, task)
}

// 启动层级时间轮
func (h *ControllableTimingWheel) Start() {
	// 从最底层开始运行
	realMillSecond := int(float64(h.minInterval.Milliseconds()) * h.timeScale)
	ticker := time.NewTicker(time.Millisecond * time.Duration(realMillSecond))
	defer ticker.Stop()
	for {
		select {
		case <-h.stopC:
			return
		case <-h.changeC:
			realMillSecond := int(float64(h.minInterval.Milliseconds()) * h.timeScale)
			ticker.Reset(time.Millisecond * time.Duration(realMillSecond))
		case <-ticker.C:
			tasks := h.wheel.advance()
			for _, task := range tasks {
				go task.run()
			}
		}
	}
}

func (h *ControllableTimingWheel) Stop() {
	h.stopC <- new(interface{})
}

func (h *ControllableTimingWheel) ChangeTimeScale(timeScale float64) {
	h.timeScale = timeScale
	h.changeC <- new(interface{})
}

// 基本时间轮结构体
type timingWheel struct {
	// 时间轮的槽位
	slots [][]task
	// 当前槽位指针
	currentSlot int
	// 槽位数
	slotCount int
	// 间隔时间
	interval time.Duration

	// 下一层时间轮
	nextWheel     *timingWheel
	nextWheelMark int
	sync.Mutex
}

type task struct {
	run   func()
	delay time.Duration
}

// 初始化基本时间轮
func newTimingWheel(slotCount int, interval time.Duration) *timingWheel {
	return &timingWheel{
		slots:     make([][]task, slotCount),
		slotCount: slotCount,
		interval:  interval,
	}
}

// 添加任务到时间轮
func (w *timingWheel) addTask(delay time.Duration, taskFunc func()) {
	// 小于本层时间轮间隔，不应当出现此情况
	if delay < w.interval {
		log.Default().Printf("WARN: add delay:%v is small than interval:%v", delay, w.interval)
	}
	// 超出本层时间轮范围，添加到下一层时间轮
	if delay > w.interval*time.Duration(w.slotCount) {
		if w.nextWheel == nil {
			w.Lock()
			if w.nextWheel == nil {
				w.nextWheel = newTimingWheel(w.slotCount, w.interval*time.Duration(w.slotCount))
				w.nextWheelMark = w.currentSlot
			}
			w.Unlock()
		}
		w.nextWheel.addTask(delay, taskFunc)
		return
	}

	delayInSlots := int(delay / w.interval)
	slot := (w.currentSlot + delayInSlots) % w.slotCount
	newDelay := delay % w.interval // should sub the time spend in this wheel
	w.slots[slot] = append(w.slots[slot], task{run: taskFunc, delay: newDelay})
}

// 推进时间轮
// 先整理本层时间轮，得到需要执行的tasks,并将本层置于下一个状态
// 然后触发下一层时间轮，并开始接收下一层时间轮的tasks
func (w *timingWheel) advance() []task {
	// 执行当前槽位的任务
	tasks := w.slots[w.currentSlot]
	w.slots[w.currentSlot] = nil
	// 更新当前槽位指针
	w.currentSlot = (w.currentSlot + 1) % w.slotCount
	// 如果当前层级不是最高层，则运行下一层级
	if w.currentSlot == w.nextWheelMark && w.nextWheel != nil {
		go func() {
			newTasks := w.nextWheel.advance()
			for _, newTask := range newTasks {
				w.addTask(newTask.delay, newTask.run)
			}
		}()
	}
	return tasks
}
