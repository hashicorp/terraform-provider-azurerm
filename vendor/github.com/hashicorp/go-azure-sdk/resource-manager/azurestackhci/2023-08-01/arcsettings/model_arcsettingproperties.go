package arcsettings

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ArcSettingProperties struct {
	AggregateState              *ArcSettingAggregateState  `json:"aggregateState,omitempty"`
	ArcApplicationClientId      *string                    `json:"arcApplicationClientId,omitempty"`
	ArcApplicationObjectId      *string                    `json:"arcApplicationObjectId,omitempty"`
	ArcApplicationTenantId      *string                    `json:"arcApplicationTenantId,omitempty"`
	ArcInstanceResourceGroup    *string                    `json:"arcInstanceResourceGroup,omitempty"`
	ArcServicePrincipalObjectId *string                    `json:"arcServicePrincipalObjectId,omitempty"`
	ConnectivityProperties      *interface{}               `json:"connectivityProperties,omitempty"`
	DefaultExtensions           *[]DefaultExtensionDetails `json:"defaultExtensions,omitempty"`
	PerNodeDetails              *[]PerNodeState            `json:"perNodeDetails,omitempty"`
	ProvisioningState           *ProvisioningState         `json:"provisioningState,omitempty"`
}
