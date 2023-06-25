package hostpool

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type HostPoolPatchProperties struct {
	CustomRdpProperty             *string                        `json:"customRdpProperty,omitempty"`
	Description                   *string                        `json:"description,omitempty"`
	FriendlyName                  *string                        `json:"friendlyName,omitempty"`
	LoadBalancerType              *LoadBalancerType              `json:"loadBalancerType,omitempty"`
	MaxSessionLimit               *int64                         `json:"maxSessionLimit,omitempty"`
	PersonalDesktopAssignmentType *PersonalDesktopAssignmentType `json:"personalDesktopAssignmentType,omitempty"`
	PreferredAppGroupType         *PreferredAppGroupType         `json:"preferredAppGroupType,omitempty"`
	PublicNetworkAccess           *PublicNetworkAccess           `json:"publicNetworkAccess,omitempty"`
	RegistrationInfo              *RegistrationInfoPatch         `json:"registrationInfo,omitempty"`
	Ring                          *int64                         `json:"ring,omitempty"`
	SsoClientId                   *string                        `json:"ssoClientId,omitempty"`
	SsoClientSecretKeyVaultPath   *string                        `json:"ssoClientSecretKeyVaultPath,omitempty"`
	SsoSecretType                 *SSOSecretType                 `json:"ssoSecretType,omitempty"`
	SsoadfsAuthority              *string                        `json:"ssoadfsAuthority,omitempty"`
	StartVMOnConnect              *bool                          `json:"startVMOnConnect,omitempty"`
	VMTemplate                    *string                        `json:"vmTemplate,omitempty"`
	ValidationEnvironment         *bool                          `json:"validationEnvironment,omitempty"`
}
