package runit

import (
	"fmt"
	"testing"
	"time"
)

func TestRunIt(t *testing.T) {
	// t.Fatal("not implemented")
	rv, err := RunIt(30*time.Second, "ls", "-a")
	_, _ = rv, err
	fmt.Printf("rv ->%s<- err : %s\n", rv, err)
}
