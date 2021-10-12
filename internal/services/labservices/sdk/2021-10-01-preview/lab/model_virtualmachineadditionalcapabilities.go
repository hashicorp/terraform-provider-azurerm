package lab

type VirtualMachineAdditionalCapabilities struct {
	InstallGpuDrivers *EnableState `json:"installGpuDrivers,omitempty"`
}
