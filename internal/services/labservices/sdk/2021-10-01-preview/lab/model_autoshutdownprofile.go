package lab

type AutoShutdownProfile struct {
	DisconnectDelay          *string             `json:"disconnectDelay,omitempty"`
	IdleDelay                *string             `json:"idleDelay,omitempty"`
	NoConnectDelay           *string             `json:"noConnectDelay,omitempty"`
	ShutdownOnDisconnect     *EnableState        `json:"shutdownOnDisconnect,omitempty"`
	ShutdownOnIdle           *ShutdownOnIdleMode `json:"shutdownOnIdle,omitempty"`
	ShutdownWhenNotConnected *EnableState        `json:"shutdownWhenNotConnected,omitempty"`
}
