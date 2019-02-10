// +build linux amd64

package container

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
	"syscall"
)

func readUserCommand() []string {
	// fd3 is the read pair of the pipe that passed from parent process
	pipe := os.NewFile(uintptr(3), "pipe")
	msg, err := ioutil.ReadAll(pipe)
	if err != nil {
		log.Fatalf("failed to read pipe %v\n", err)
		return nil
	}
	msgStr := string(msg)
	return strings.Split(msgStr, " ")
}

func RunContainerInitProcess() error {
	cmdArray := readUserCommand()
	if cmdArray == nil || len(cmdArray) == 0 {
		return fmt.Errorf("get user command error, cmd is empty")
	}
	log.Printf("command %v\n", cmdArray)

	defaultMountFlags := syscall.MS_NOEXEC | syscall.MS_NOSUID | syscall.MS_NODEV
	syscall.Mount("proc", "/proc", "proc", uintptr(defaultMountFlags), "")
	cmd := cmdArray[0]
	path, err := exec.LookPath(cmd)
	if err != nil {
		log.Fatalf("failed to find cmd %s\n", cmd)
		return err
	}
	if err := syscall.Exec(path, cmdArray[0:], os.Environ()); err != nil {
		log.Fatal(err.Error())
	}
	return nil
}
