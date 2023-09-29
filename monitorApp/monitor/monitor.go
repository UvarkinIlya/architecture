package monitor

import (
	"log"
	"net/http"
	"time"
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
			log.Println("Reset service:", m.serviceName)
			return
		}
		time.Sleep(m.timeBetweenCheck)
	}
}

func (m *MonitorImpl) reset() {

}
