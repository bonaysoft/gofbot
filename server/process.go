package main

import (
	`fmt`
	`io/ioutil`
	`os`
	`strconv`
)

var pidFile = "gofbot.lock"

func savePid() error {
	_, err := os.Stat(pidFile) // os.Stat获取文件信息
	if (err != nil && os.IsExist(err)) || err == nil {
		return fmt.Errorf("gofbot already running.")
	}

	pid := strconv.Itoa(os.Getpid())
	return ioutil.WriteFile(pidFile, []byte(pid), 0600)
}

func cleanPid() error {
	return os.Remove(pidFile)
}

func findProcess() (*os.Process, error) {
	data, err := ioutil.ReadFile(pidFile)
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
