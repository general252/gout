package utask

import (
	"fmt"
	"github.com/robfig/cron/v3"
	"log"
	"sync"
	"time"
)

var (
	cronTab      *cron.Cron = nil
	cronLocation            = time.Local
	mux          sync.Mutex
)

// SetLocation 设置时区
func SetLocation(loc *time.Location) {
	cronLocation = loc
}

// AddTaskEvery 添加任务. id: 任务标识
func AddTaskEvery(everyDuration time.Duration, cmd func()) (int, error) {
	mux.Lock()
	defer mux.Unlock()

	var spec = fmt.Sprintf("@every %v", everyDuration.String())
	return AddCron(spec, cmd)
}

// AddCron AddCron("* * * * * ?", func() { })
func AddCron(spec string, cmd func()) (int, error) {
	mux.Lock()
	defer mux.Unlock()

	if cronTab == nil {
		cronTab = cron.New(cron.WithSeconds(), cron.WithLocation(cronLocation))
		cronTab.Start()
	}

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
