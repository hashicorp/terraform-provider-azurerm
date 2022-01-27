package afdorigins

type AFDOriginProperties struct {
	AzureOrigin                 *ResourceReference    `json:"azureOrigin,omitempty"`
	DeploymentStatus            *DeploymentStatus     `json:"deploymentStatus,omitempty"`
	EnabledState                *EnabledState         `json:"enabledState,omitempty"`
	EnforceCertificateNameCheck *bool                 `json:"enforceCertificateNameCheck,omitempty"`
	HostName                    string                `json:"hostName"`
	HttpPort                    *int64                `json:"httpPort,omitempty"`
	HttpsPort                   *int64                `json:"httpsPort,omitempty"`
	OriginGroupName             *string               `json:"originGroupName,omitempty"`
	OriginHostHeader            *string               `json:"originHostHeader,omitempty"`
	Priority                    *int64                `json:"priority,omitempty"`
	ProvisioningState           *AfdProvisioningState `json:"provisioningState,omitempty"`
	SharedPrivateLinkResource   *interface{}          `json:"sharedPrivateLinkResource,omitempty"`
	Weight                      *int64                `json:"weight,omitempty"`
}
