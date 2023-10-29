package app_manager

import (
	"bytes"
	"errors"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"architecture/logger"
)

type AppManager interface {
	Start() (err error)
	Stop() (err error)
	Restart() (err error)
}

type ManagerImpl struct {
	checker Checker
	pid     int
}

func NewManager(checker Checker) *ManagerImpl {
	return &ManagerImpl{
		checker: checker,
	}
}

func (m *ManagerImpl) Start() (err error) {
	cmd := exec.Command("scripts/start")

	var bufErr bytes.Buffer
	var bufOut bytes.Buffer
	cmd.Stderr = &bufErr
	cmd.Stdout = &bufOut

	err = cmd.Start()
	if err != nil {
		return err
	}

	if bufErr.String() != "" {
		return errors.New(bufErr.String())
	}

	time.Sleep(500 * time.Millisecond)
	m.pid, err = strconv.Atoi(strings.Replace(bufOut.String(), "\n", "", -1))
	logger.Info("Service pid: %d", m.pid)
	if err != nil {
		return err
	}

	time.Sleep(3 * time.Second)
	for {
		err = m.check()
		logger.Error("Check failed: %s", err)
		err = m.Restart()
		if err != nil {
			return err
		}
	}
}

func (m *ManagerImpl) Stop() (err error) {
	cmd := exec.Command("kill", "-9", strconv.Itoa(m.pid))
	return cmd.Run()
}

func (m *ManagerImpl) Restart() (err error) {
	err = m.Stop()
	if err != nil {
		logger.Error("Failed stop err: %s", err.Error())
	}
	err = m.Start()
	if err != nil {
		logger.Error("Failed start err:", err.Error())
		return
	}

	return nil
}

func (m *ManagerImpl) check() error {
	<-m.checker.Check()
	return errors.New("check failed")
}
