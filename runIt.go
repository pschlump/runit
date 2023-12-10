package runit

import (
	"bytes"
	"errors"
	"io"
	"os/exec"
	"syscall"
	"time"
)

func RunIt(timeout time.Duration, name string, args ...string) (result []byte, err error) {

	cmd := exec.Command(name, args...)
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true,
		Pgid:    0,
	}

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return
	}

	outBuf := bytes.NewBuffer(nil)
	errBuf := bytes.NewBuffer(nil)

	err = cmd.Start()
	if err != nil {
		return
	}

	stdin.Close()

	go io.Copy(outBuf, stdout)
	go io.Copy(errBuf, stderr)

	//	result = errBuf.String()
	//	if len(result) > 0 {
	//		err = fmt.Errorf("Non Empty Stderr")
	//		return
	//	}
	//	result = outBuf.String()

	ch := make(chan error)

	go func(cmd *exec.Cmd) {
		defer close(ch)
		ch <- cmd.Wait()
	}(cmd)

	select {
	case err = <-ch:
	case <-time.After(timeout):
		cmd.Process.Kill()
		err = errors.New("execute timeout")
		return
	}

	if err != nil {
		errStr := errBuf.String()
		return nil, errors.New(errStr)
	}

	if outBuf.Len() > 0 {
		return outBuf.Bytes(), nil
	}

	return
}
