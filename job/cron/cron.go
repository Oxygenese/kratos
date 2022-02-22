package cron

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/robfig/cron"
)

func NewCron() *CronJob {
	return &CronJob{cron: cron.New()}
}

type CronJob struct {
	cron *cron.Cron
}

func (e CronJob) Stop() {
	e.cron.Stop()
}

func (e CronJob) Run() {
	e.cron.Run()
}

func (e CronJob) Start() {
	e.cron.Start()
}

func (e CronJob) AddFunc(spec string, cmd func()) {
	e.cron.AddFunc(spec, cmd)
}

func (e CronJob) AddJob(spec string, job cron.Job) {
	e.cron.AddJob(spec, job)
}

func (e CronJob) Schedule(spec string, job cron.Job) {
	s, err := cron.Parse(spec)
	if err != nil {
		log.Errorf("cron schedule parse spec err:%s", err)
	} else {
		e.cron.Schedule(s, job)
	}
}
