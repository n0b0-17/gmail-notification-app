package scheduler

import (
	"log"
	"time"

	"github.com/robfig/cron/v3"
)

type Scheduler struct  {
	cron *cron.Cron
	job *DailyTotalJob
}

func NewScheduler(baseURL string) *Scheduler {
	jst := cron.NewParser(
		cron.Second | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.DowOptional,
	)
	
	c := cron.New(cron.WithParser(jst), cron.WithLocation(time.FixedZone("Asia/Tokyo", 9*60*60)))

	return &Scheduler{
		cron: c,
		job: NewDailyTotalJob(baseURL),
	}
}


func (s *Scheduler) Start() error {
	_,err := s.cron.AddFunc("0 01 09 * * *", func() {
		if err := s.job.Execute(); err != nil {
			log.Printf("Error executing daily total job: %v", err)
		}
	})

	if err != nil {
		return err
	}

	s.cron.Start()
	return nil 
}

func (s *Scheduler) Stop() {
	s.cron.Stop()
}

