package process

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
)

type PidCtrl struct {
	pidFile string
}

func New(pidFile string) *PidCtrl {
	return &PidCtrl{pidFile: pidFile}
}

func (p *PidCtrl) Save() error {
	_, err := os.Stat(p.pidFile) // os.Stat获取文件信息
	if (err != nil && os.IsExist(err)) || err == nil {
		return fmt.Errorf("gofbot already running.")
	}

	pid := strconv.Itoa(os.Getpid())
	return ioutil.WriteFile(p.pidFile, []byte(pid), 0600)
}

func (p *PidCtrl) Clean() error {
	return os.Remove(p.pidFile)
}

func (p *PidCtrl) Find() (*os.Process, error) {
	data, err := ioutil.ReadFile(p.pidFile)
	if err != nil {
		return nil, err
	}

	pidStr := string(data)
	pid, err := strconv.Atoi(pidStr)
	if err != nil {
		return nil, err
	}

	return os.FindProcess(pid)
}
