package accounts

type HiveMetastoreProperties struct {
	DatabaseName                    *string                          `json:"databaseName,omitempty"`
	NestedResourceProvisioningState *NestedResourceProvisioningState `json:"nestedResourceProvisioningState,omitempty"`
	Password                        *string                          `json:"password,omitempty"`
	RuntimeVersion                  *string                          `json:"runtimeVersion,omitempty"`
	ServerUri                       *string                          `json:"serverUri,omitempty"`
	UserName                        *string                          `json:"userName,omitempty"`
}
