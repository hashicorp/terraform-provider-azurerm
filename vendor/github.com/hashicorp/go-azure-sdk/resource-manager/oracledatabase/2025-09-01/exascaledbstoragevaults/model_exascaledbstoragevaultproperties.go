package exascaledbstoragevaults

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ExascaleDbStorageVaultProperties struct {
	AdditionalFlashCacheInPercent    *int64                                `json:"additionalFlashCacheInPercent,omitempty"`
	AttachedShapeAttributes          *[]ShapeAttribute                     `json:"attachedShapeAttributes,omitempty"`
	Description                      *string                               `json:"description,omitempty"`
	DisplayName                      string                                `json:"displayName"`
	ExadataInfrastructureId          *string                               `json:"exadataInfrastructureId,omitempty"`
	HighCapacityDatabaseStorage      *ExascaleDbStorageDetails             `json:"highCapacityDatabaseStorage,omitempty"`
	HighCapacityDatabaseStorageInput ExascaleDbStorageInputDetails         `json:"highCapacityDatabaseStorageInput"`
	LifecycleDetails                 *string                               `json:"lifecycleDetails,omitempty"`
	LifecycleState                   *ExascaleDbStorageVaultLifecycleState `json:"lifecycleState,omitempty"`
	OciURL                           *string                               `json:"ociUrl,omitempty"`
	Ocid                             *string                               `json:"ocid,omitempty"`
	ProvisioningState                *AzureResourceProvisioningState       `json:"provisioningState,omitempty"`
	TimeZone                         *string                               `json:"timeZone,omitempty"`
	VMClusterCount                   *int64                                `json:"vmClusterCount,omitempty"`
}
