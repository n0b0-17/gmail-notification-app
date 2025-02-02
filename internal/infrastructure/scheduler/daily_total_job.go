package scheduler

import (
	"fmt"
	"net/http"
)

type DailyTotalJob struct {
	baseURL string
}

func NewDailyTotalJob(baseURL string) *DailyTotalJob {
	return &DailyTotalJob{baseURL: baseURL}
}

func (j *DailyTotalJob) Execute() error {
	url := fmt.Sprintf("%s/gmail/daily-totals", j.baseURL)
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to execute daily total job: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("daily total job failed with status: %d", resp.StatusCode)
	}
	return nil
}
