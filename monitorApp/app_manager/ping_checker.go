package app_manager

import (
	"net/http"
	"time"
)

type Checker interface {
	Check() (failed chan struct{}) //Change to chan error
}

type PingChecker struct {
	checkURL      string
	checkInterval time.Duration
	maxFailed     uint
}

func NewPingChecker(checkURL string, checkInterval time.Duration, maxFailed uint) *PingChecker {
	return &PingChecker{checkURL: checkURL, checkInterval: checkInterval, maxFailed: maxFailed}
}

func (p PingChecker) Check() (failed chan struct{}) {
	failed = make(chan struct{})
	var failedCount uint

	go func() {
		for {
			if failedCount > p.maxFailed {
				failed <- struct{}{}
			}

			select {
			case <-time.After(10 * time.Second):
				failedCount++
			case isFailed := <-p.check():
				if isFailed {
					failedCount++
					continue
				}

				failedCount = 0
			}
		}
	}()

	return failed
}

func (p PingChecker) check() (isFailed chan bool) {
	isFailed = make(chan bool)
	go func() {
		_, err := http.Get(p.checkURL)
		if err != nil {
			isFailed <- true
		}

		isFailed <- false
	}()

	return isFailed
}
