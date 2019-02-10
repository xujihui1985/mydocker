package subsystems

import (
	"io/ioutil"
	"os"
	"path"
	"strconv"
)

type Memory struct {

}

func (s *Memory) Set(cgroupPath string, cfg *ResourceConfig) error {
	if cp, err := GetCgroupPath(s.Name(), cgroupPath, true); err == nil {
		if cfg.MemoryLimit != "" {
			if err := ioutil.WriteFile(path.Join(cp, "memory.limit_in_bytes"), []byte(cfg.MemoryLimit), 0644); err != nil {
				return err
			}
		}
		return nil
	} else {
		return err
	}
}

func (s *Memory) Remove(cgroupPath string) error {
	if cp, err := GetCgroupPath(s.Name(), cgroupPath, false); err == nil {
		return os.Remove(cp)
	} else {
		return err
	}
}

func (s *Memory) Apply(cgroupPath string, pid int) error {
	if cp, err := GetCgroupPath(s.Name(), cgroupPath, false); err == nil {
		if err := ioutil.WriteFile(path.Join(cp, "tasks"), []byte(strconv.Itoa(pid)), 0644); err != nil {
			return err
		}
		return nil
	} else {
		return err
	}
}

func (s *Memory) Name() string{
	return "memory"
}