package virtualmachine

type VirtualMachineConnectionProfile struct {
	AdminUsername    *string `json:"adminUsername,omitempty"`
	NonAdminUsername *string `json:"nonAdminUsername,omitempty"`
	PrivateIpAddress *string `json:"privateIpAddress,omitempty"`
	RdpAuthority     *string `json:"rdpAuthority,omitempty"`
	RdpInBrowserUrl  *string `json:"rdpInBrowserUrl,omitempty"`
	SshAuthority     *string `json:"sshAuthority,omitempty"`
	SshInBrowserUrl  *string `json:"sshInBrowserUrl,omitempty"`
}
