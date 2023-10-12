package app_manager

import (
	"bytes"
	"errors"
	"log"
	"net/http"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type AppManager interface {
	Start() (err error)
	Stop() (err error)
	Restart() (err error)
	Check()
}

type ManagerImpl struct {
	serviceName      string
	checkURL         string
	timeBetweenCheck time.Duration
	pid              int
}

func NewManager(serviceName string, checkURL string, timeBetweenCheck time.Duration) *ManagerImpl {
	return &ManagerImpl{
		serviceName:      serviceName,
		checkURL:         checkURL,
		timeBetweenCheck: timeBetweenCheck,
	}
}

func (m *ManagerImpl) Start() (err error) {
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

	m.pid, err = strconv.Atoi(strings.Replace(bufOut.String(), "\n", "", -1))
	log.Println("Service pid:", m.pid)
	if err != nil {
		return err
	}

	return nil
}

func (m *ManagerImpl) Stop() (err error) {
	cmd := exec.Command("kill", "-9", strconv.Itoa(m.pid))
	return cmd.Run()
}

func (m *ManagerImpl) Restart() (err error) {
	err = m.Stop()
	if err != nil {
		log.Println("Failed stop err:", err.Error())
	}
	err = m.Start()
	if err != nil {
		log.Println("Failed start err:", err.Error())
		return
	}

	return nil
}

func (m *ManagerImpl) Check() {
	for {
		select {
		case <-time.After(10 * time.Second):
			log.Println("Restart service:", m.serviceName)
			m.Restart()
		case <-m.check():
		}

		time.Sleep(m.timeBetweenCheck)
	}
}

func (m *ManagerImpl) check() (ch chan struct{}) {
	ch = make(chan struct{})
	go func() {
		_, err := http.Get(m.checkURL)
		if err != nil {
			log.Println("Restart service:", m.serviceName)
			m.Restart()
		}

		ch <- struct{}{}
	}()

	return ch
}
