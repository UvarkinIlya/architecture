package monitor

import (
	"net/http"
	"time"

	"architecture/logger"
)

type Monitor interface {
	Check()
}

type MonitorImpl struct {
	serviceName      string
	checkURL         string
	timeBetweenCheck time.Duration
}

func NewMonitor() *MonitorImpl {
	return &MonitorImpl{}
}

func (m *MonitorImpl) Check() {
	for {
		_, err := http.Get(m.checkURL)
		if err != nil {
			logger.Error("Reset service: %s", m.serviceName)
			return
		}
		time.Sleep(m.timeBetweenCheck)
	}
}

func (m *MonitorImpl) reset() {

}
