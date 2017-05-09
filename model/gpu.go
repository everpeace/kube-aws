package model

import (
	"errors"
	"fmt"
	"strings"
)

var GPUEnabledInstanceFamily = []string{"p2", "g2"}

type Gpu struct {
	Nvidia NvidiaSetting `yaml:"nvidia"`
}

type NvidiaSetting struct {
	Enabled bool   `yaml:"enabled,omitempty"`
	Version string `yaml:"version,omitempty"`
}

func isGpuEnabledInstanceType(instanceType string) bool {
	for _, family := range GPUEnabledInstanceFamily {
		if strings.HasPrefix(instanceType, family) {
			return true
		}
	}
	return false
}

func newDefaultGpu() Gpu {
	return Gpu{
		Nvidia: NvidiaSetting{
			Enabled: false,
			Version: "",
		},
	}
}

// This function is used when rendering cloud-config-worker
func (c NvidiaSetting) IsEnabledOn(instanceType string) bool {
	return isGpuEnabledInstanceType(instanceType) && c.Enabled
}

func (c Gpu) Valid(instanceType string) error {
	if c.Nvidia.Enabled && len(c.Nvidia.Version) == 0 {
		return errors.New(`gpu.nvidia.version must not be empty when gpu.nvidia is enabled.`)
	}
	if c.Nvidia.Enabled && !isGpuEnabledInstanceType(instanceType) {
		fmt.Printf("WARNING: instance type %v doesn't support GPU (only %v instance family do). GPU setting will be ignored.\n", instanceType, GPUEnabledInstanceFamily)

	}
	if !c.Nvidia.Enabled && isGpuEnabledInstanceType(instanceType) {
		fmt.Printf("WARNING: Nvidia GPU driver intallation is disabled although instance type %v does support GPU.  You have to install Nvidia GPU driver by yourself to schedule gpu resource.\n", instanceType)
	}

	return nil
}
