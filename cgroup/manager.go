package cgroup

import (
	"fmt"
	"github.com/xujihui1985/mydocker/cgroup/subsystems"
)

type Manager struct {
	Path string
	Resource *subsystems.ResourceConfig
}

func NewManager(path string) *Manager {
	return &Manager{
		Path: path,
	}
}

func (c *Manager) Apply(pid int) error {
	for _, subSys := range subsystems.Instances {
		if err := subSys.Apply(c.Path, pid); err != nil {
			return err
		}
	}
	return nil
}

func (c *Manager) Set(cfg *subsystems.ResourceConfig) error {
	for _, subSys := range subsystems.Instances {
		if err := subSys.Set(c.Path, cfg); err != nil {
			return err
		}
	}
	return nil
}

func (c *Manager) Destroy() error {
	for _, subSys := range subsystems.Instances {
		if err := subSys.Remove(c.Path); err != nil {
			return fmt.Errorf("remove cgroup fail %v", err)
		}
	}
	return nil
}