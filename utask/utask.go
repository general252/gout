package utask

import (
	"fmt"
	"github.com/robfig/cron/v3"
	"log"
	"sync"
	"time"
)

var (
	cronTab *cron.Cron = nil
	mux     sync.Mutex
)

// AddTaskEvery 添加任务. id: 任务标识
func AddTaskEvery(everyDuration time.Duration, cmd func()) (int, error) {
	mux.Lock()
	defer mux.Unlock()

	if cronTab == nil {
		cronTab = cron.New(cron.WithSeconds())
		cronTab.Start()
	}

	var spec = fmt.Sprintf("@every %v", everyDuration.String())
	id, err := cronTab.AddFunc(spec, cmd)

	return int(id), err
}

// RemoveTask 删除任务
func RemoveTask(id int) {
	mux.Lock()
	defer mux.Unlock()

	if cronTab == nil {
		return
	}

	cronTab.Remove(cron.EntryID(id))
}

// StopCron 关闭任务调度
func StopCron() {
	mux.Lock()
	defer mux.Unlock()

	if cronTab != nil {
		ctx := cronTab.Stop()
		select {
		case <-ctx.Done():
		case <-time.After(time.Second * 10):
			log.Println("context was not done immediately")
		}
	}
}
