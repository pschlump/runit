package runit

import (
	"context"
	"errors"
	"os/exec"
	"time"
)

func RunIt(timeout time.Duration, shellcmd string, args ...string) (rv []byte, err error) {

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, shellcmd, args...)

	out, err := cmd.CombinedOutput()

	if ctx.Err() == context.DeadlineExceeded {
		err = errors.New("execute timeout")
		return
	}

	if err != nil {
		return []byte(out), err
	}

	return []byte(out), nil

}
