package subsystems


type ResourceConfig struct {
	MemoryLimit string
	CpuShare string
	CpuSet string
}

type SubSystem interface {
	Name() string
	Set(path string, cfg *ResourceConfig) error
	Apply(path string, pid int) error
	Remove(path string) error
}

var Instances = []SubSystem{
	&Memory{},
}