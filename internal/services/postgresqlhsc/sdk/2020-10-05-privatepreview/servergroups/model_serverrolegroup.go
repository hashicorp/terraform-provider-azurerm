package servergroups

type ServerRoleGroup struct {
	EnableHa         *bool             `json:"enableHa,omitempty"`
	EnablePublicIp   *bool             `json:"enablePublicIp,omitempty"`
	Name             *string           `json:"name,omitempty"`
	Role             *ServerRole       `json:"role,omitempty"`
	ServerCount      *int64            `json:"serverCount,omitempty"`
	ServerEdition    *ServerEdition    `json:"serverEdition,omitempty"`
	ServerNames      *[]ServerNameItem `json:"serverNames,omitempty"`
	StorageQuotaInMb *int64            `json:"storageQuotaInMb,omitempty"`
	VCores           *int64            `json:"vCores,omitempty"`
}
