// +build linux amd64

package container

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
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
    setUpMount()
	// defaultMountFlags := syscall.MS_NOEXEC | syscall.MS_NOSUID | syscall.MS_NODEV
	// syscall.Mount("proc", "/proc", "proc", uintptr(defaultMountFlags), "")
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

func pivotRoot(root string) error {
	if err := syscall.Mount(root, root, "bind", syscall.MS_BIND | syscall.MS_REC, ""); err != nil {
		return fmt.Errorf("failed to mount root %v", err)
	}
	pivotDir := filepath.Join(root, ".pivot_root")
	if err := os.MkdirAll(pivotDir, 0777); err != nil {
		return err
	}
	if err := syscall.PivotRoot(root, pivotDir); err != nil {
		return err
	}
	if err := syscall.Chdir("/"); err != nil {
		return err
	}
	pivotDir = filepath.Join("/", ".pivot_root")
	if err := syscall.Unmount(pivotDir, syscall.MNT_DETACH); err != nil {
		return err
	}
	return os.Remove(pivotDir)
}

func setUpMount() {
	pwd, err := os.Getwd()
	if err != nil {
		return
	}

	pivotRoot(pwd)
	defaultMountFlags := syscall.MS_NOEXEC | syscall.MS_NOSUID | syscall.MS_NODEV
	syscall.Mount("proc", "/proc", "proc", uintptr(defaultMountFlags), "")
	syscall.Mount("tmpfs", "/dev", "tmpfs", syscall.MS_NOSUID | syscall.MS_STRICTATIME, "mode=755")
}