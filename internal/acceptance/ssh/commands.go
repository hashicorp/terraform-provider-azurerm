package ssh

const (
	// LinuxAgentDeprovisionCommand is the command needed to "deprovision" the (Windows)Azure VM Agent
	// this clears out caches, resets the vm agent configuration and generally "generalizes" this VM
	// and is vaguely akin to `sysprep` on Windows.
	LinuxAgentDeprovisionCommand = "sudo waagent -verbose -deprovision+user -force"
)
