package app_manager

import (
	"bytes"
	"errors"
	"fmt"
	"os/exec"
	"strconv"
	"time"

	"architecture/logger"
)

type AppManager interface {
	Start() (err error)
	Stop() (err error)
	Restart() (err error)
}

type ManagerImpl struct {
	checker      Checker
	serverConfig string
	serverPath   string
	pid          int
}

func NewManager(checker Checker, serverConfig, serverPath string) *ManagerImpl {
	return &ManagerImpl{
		checker:      checker,
		serverConfig: serverConfig,
		serverPath:   serverPath,
	}
}

func (m *ManagerImpl) Start() (err error) {
	cmd := exec.Command(fmt.Sprintf("./utils/%s", m.serverPath), "--config", m.serverConfig)

	var bufErr bytes.Buffer
	var bufOut bytes.Buffer
	cmd.Stderr = &bufErr
	cmd.Stdout = &bufOut

	err = cmd.Start()
	if err != nil {
		logger.Info("Failed start server err:%s", err)
		return err
	}

	if bufErr.String() != "" {
		return errors.New(bufErr.String())
	}

	m.pid = cmd.Process.Pid
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
