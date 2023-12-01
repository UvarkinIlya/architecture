package app_manager

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"time"

	"architecture/logger"
	"architecture/modellibrary"
)

type WatchdogChecker struct {
	startWatchdogReq modellibrary.WatchdogStartRequest
	startWatchdogURL string
	checkInterval    time.Duration
	maxWaitTime      time.Duration
}

func NewWatchdogChecker(startWatchdogReq modellibrary.WatchdogStartRequest, startWatchdogURL string, checkInterval, maxWaitTime int) *WatchdogChecker {
	return &WatchdogChecker{
		startWatchdogReq: startWatchdogReq,
		startWatchdogURL: startWatchdogURL,
		checkInterval:    time.Duration(checkInterval) * time.Second,
		maxWaitTime:      time.Duration(maxWaitTime) * time.Second,
	}
}

func (w WatchdogChecker) Check() (failed chan struct{}) {
	failed = make(chan struct{})

	go func() {
		err := w.startServiceTick()
		if err != nil {
			logger.Error("Failed start server tick due to err: %s", err)
			failed <- struct{}{}
		}

		_ = w.waitFailed()
		failed <- struct{}{}
	}()

	return failed
}

func (w WatchdogChecker) waitFailed() (err error) {
	ticker := time.NewTicker(w.checkInterval)
	for range ticker.C {
		file, err := os.Open(w.startWatchdogReq.FileName)
		if err != nil {
			return err
		}

		stat, err := file.Stat()
		if err != nil {
			return err
		}

		if stat.ModTime().Before(time.Now().Add(-w.maxWaitTime)) {
			return errors.New("file update timeout")
		}
	}

	return errors.New("watchdog ticker stopped")
}

func (w WatchdogChecker) startServiceTick() (err error) {
	bin, err := json.Marshal(w.startWatchdogReq)
	if err != nil {
		return err
	}

	resp, err := http.Post(w.startWatchdogURL, "application/json", bytes.NewBuffer(bin))
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return errors.New("Failed watchdog start, due to error statusCode != 200. ")
	}

	return nil
}
