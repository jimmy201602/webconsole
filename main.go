package main

import (
	"webconsole/utils"
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"syscall"
	"time"

	"webconsole/website"
)

const (
	START  string = "start"
	STOP   string = "stop"
	STATUS string = "status"
)

type Apibox struct {
	PID int
}

func (a *Apibox) GetPID() (*Apibox, error) {
	b, err := ioutil.ReadFile(utils.PidPath)
	if nil != err {
		return nil, err
	}
	b = bytes.TrimSpace(b)
	pid, err := utils.StringUtils(string(b)).Int()
	if nil != err {
		return nil, err
	}
	a.PID = pid
	return a, nil
}

func (a *Apibox) Start() error {
	website.Run()
	return nil
}

func (a *Apibox) Stop() error {
	time.Sleep(time.Duration(1 * time.Second))
	if utils.Exists(utils.PidPath) {
		_, err := a.GetPID()
		if nil != err {
			return err
		}
		p, err := os.FindProcess(a.PID)
		if nil != err {
			return err
		}
		err = p.Kill()
		if nil != err {
			return err
		}
	} else {
		return fmt.Errorf("Unable to read the PID file.")
	}
	return nil
}

func (a *Apibox) Status() (bool, error) {
	if utils.Exists(utils.PidPath) {
		_, err := a.GetPID()
		if nil != err {
			return false, err
		}
		if err := syscall.Kill(a.PID, 0); nil != err {
			return false, nil
		} else {
			return true, nil
		}

	} else {
		return false, fmt.Errorf("Unable to read the PID file.")
	}
}

func main() {
	flag.Parse()
	var cmd string = flag.Arg(0)
	cmd = strings.ToLower(cmd)
	switch strings.TrimSpace(cmd) {
	case START:
		a := &Apibox{}
		err := a.Start()
		if nil != err {
			fmt.Fprintf(os.Stderr, err.Error()+"\n")
		}
	case STOP:
		a := &Apibox{}
		err := a.Stop()
		if nil != err {
			fmt.Fprintf(os.Stderr, err.Error()+"\n")
		}
	case STATUS:
		a := &Apibox{}
		t, err := a.Status()
		if nil != err {
			fmt.Fprintf(os.Stderr, err.Error()+"\n")
		}
		if !t {
			fmt.Fprintf(os.Stdout, "Stop.\n")
		} else {
			fmt.Fprintf(os.Stdout, "Running...\n")
		}
	default:
		fmt.Fprintf(os.Stderr, "Usage: %s command <start|stop|status>\n", os.Args[0])
	}
}
