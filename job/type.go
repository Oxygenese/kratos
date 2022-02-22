package job

import "github.com/robfig/cron"

type Job interface {
	Run()
	Stop()
	Start()
	AddFunc(spec string, cmd func())
	AddJob(spec string, job cron.Job)
	Schedule(spec string, job cron.Job)
}
