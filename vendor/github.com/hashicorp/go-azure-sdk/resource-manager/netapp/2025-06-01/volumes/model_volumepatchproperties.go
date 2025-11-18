package volumes

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type VolumePatchProperties struct {
	CoolAccess                *bool                                `json:"coolAccess,omitempty"`
	CoolAccessRetrievalPolicy *CoolAccessRetrievalPolicy           `json:"coolAccessRetrievalPolicy,omitempty"`
	CoolAccessTieringPolicy   *CoolAccessTieringPolicy             `json:"coolAccessTieringPolicy,omitempty"`
	CoolnessPeriod            *int64                               `json:"coolnessPeriod,omitempty"`
	DataProtection            *VolumePatchPropertiesDataProtection `json:"dataProtection,omitempty"`
	DefaultGroupQuotaInKiBs   *int64                               `json:"defaultGroupQuotaInKiBs,omitempty"`
	DefaultUserQuotaInKiBs    *int64                               `json:"defaultUserQuotaInKiBs,omitempty"`
	ExportPolicy              *VolumePatchPropertiesExportPolicy   `json:"exportPolicy,omitempty"`
	IsDefaultQuotaEnabled     *bool                                `json:"isDefaultQuotaEnabled,omitempty"`
	ProtocolTypes             *[]string                            `json:"protocolTypes,omitempty"`
	ServiceLevel              *ServiceLevel                        `json:"serviceLevel,omitempty"`
	SmbAccessBasedEnumeration *SmbAccessBasedEnumeration           `json:"smbAccessBasedEnumeration,omitempty"`
	SmbNonBrowsable           *SmbNonBrowsable                     `json:"smbNonBrowsable,omitempty"`
	SnapshotDirectoryVisible  *bool                                `json:"snapshotDirectoryVisible,omitempty"`
	ThroughputMibps           *float64                             `json:"throughputMibps,omitempty"`
	UnixPermissions           *string                              `json:"unixPermissions,omitempty"`
	UsageThreshold            *int64                               `json:"usageThreshold,omitempty"`
}
