package nodetype

type VMSSExtensionProperties struct {
	AutoUpgradeMinorVersion  *bool        `json:"autoUpgradeMinorVersion,omitempty"`
	ForceUpdateTag           *string      `json:"forceUpdateTag,omitempty"`
	ProtectedSettings        *interface{} `json:"protectedSettings,omitempty"`
	ProvisionAfterExtensions *[]string    `json:"provisionAfterExtensions,omitempty"`
	ProvisioningState        *string      `json:"provisioningState,omitempty"`
	Publisher                string       `json:"publisher"`
	Settings                 *interface{} `json:"settings,omitempty"`
	Type                     string       `json:"type"`
	TypeHandlerVersion       string       `json:"typeHandlerVersion"`
}
