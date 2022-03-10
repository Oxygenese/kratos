package cron

import (
	"github.com/projects-mars/kratos/v2/log"
	"github.com/robfig/cron"
)

func NewCron() *Job {
	return &Job{cron: cron.New()}
}

type Job struct {
	cron *cron.Cron
}

func (e Job) Stop() {
	e.cron.Stop()
}

func (e Job) Run() {
	e.cron.Run()
}

func (e Job) Start() {
	e.cron.Start()
}

func (e Job) AddFunc(spec string, cmd func()) {
	e.cron.AddFunc(spec, cmd)
}

func (e Job) AddJob(spec string, job cron.Job) {
	e.cron.AddJob(spec, job)
}

func (e Job) Schedule(spec string, job cron.Job) {
	s, err := cron.Parse(spec)
	if err != nil {
		log.Errorf("cron schedule parse spec err:%s", err)
	} else {
		e.cron.Schedule(s, job)
	}
}
