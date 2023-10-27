package watchdog

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"architecture/modellibrary"
)

type Watchdog interface {
	StartService(serverName string, interval int, maxWaitTime int) (err error)
}

type WatchdogImpl struct {
	service map[string]watchdogCheckModel
	pid     int
}

func NewWatchdog(service map[string]watchdogCheckModel) *WatchdogImpl {
	return &WatchdogImpl{service: service}
}

func (w WatchdogImpl) StartService(serverName string, interval int, maxWaitTime int) (err error) {
	err = w.start()
	if err != nil {
		return err
	}

	time.Sleep(1 * time.Second)

	watchdogReq := modellibrary.WatchdogStartRequest{
		FileName:        w.service[serverName].Url,
		IntervalSeconds: interval,
	}

	bin, err := json.Marshal(watchdogReq)
	if err != nil {
		return err
	}

	resp, err := http.Post(w.service[serverName].Url, "application/json", bytes.NewBuffer(bin))
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return errors.New("Failed watchdog start, due to error statusCode != 200. ")
	}

	reboot := make(chan struct{})
	go startWatchdog(reboot, w.service[serverName].FilePath, interval, maxWaitTime)

	<-reboot
	err = w.restart()
	if err != nil {
		return err
	}

	err = w.StartService(serverName, interval, maxWaitTime)
	if err != nil {
		return err
	}

	return nil
}

func (w WatchdogImpl) start() (err error) {
	cmd := exec.Command("monitorApp/scripts/start")

	var bufErr bytes.Buffer
	var bufOut bytes.Buffer
	cmd.Stderr = &bufErr
	cmd.Stdout = &bufOut

	err = cmd.Run()
	if err != nil {
		return err
	}

	if bufErr.String() != "" {
		return errors.New(bufErr.String())
	}

	w.pid, err = strconv.Atoi(strings.Replace(bufOut.String(), "\n", "", -1))
	log.Println("Service pid:", w.pid)
	if err != nil {
		return err
	}

	return nil
}

func (w WatchdogImpl) stop() (err error) {
	cmd := exec.Command("kill", "-9", strconv.Itoa(w.pid))
	return cmd.Run()
}

func (w WatchdogImpl) restart() (err error) {
	err = w.stop()
	if err != nil {
		log.Println("Failed stop err:", err.Error())
	}
	err = w.start()
	if err != nil {
		log.Println("Failed start err:", err.Error())
		return
	}

	return nil
}

func startWatchdog(reboot chan struct{}, filePath string, interval, maxWaitTime int) {
	var err error
	defer func() {
		if err != nil {
			log.Println("Failed watchdog, due to error", err.Error())
		}
	}()

	ticker := time.NewTicker(time.Duration(interval) * time.Second)
	for {
		select {
		case <-ticker.C:
			data, err := os.ReadFile(filePath)
			if err != nil {
				return
			}

			serverTime, err := time.Parse(string(data), time.RFC3339)
			if err != nil {
				return
			}

			if serverTime.Before(time.Now().Add(-time.Duration(maxWaitTime) * time.Second)) {
				reboot <- struct{}{}
				return
			}

		}
	}
}
