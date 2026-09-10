package datashares

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type StorageDataShareProperties struct {
	AccessPolicies      []StorageDataShareAccessPolicy      `json:"accessPolicies"`
	Assets              []StorageDataShareAsset             `json:"assets"`
	DataShareIdentifier *string                             `json:"dataShareIdentifier,omitempty"`
	DataShareUri        *string                             `json:"dataShareUri,omitempty"`
	Description         *string                             `json:"description,omitempty"`
	ProvisioningState   *NativeDataSharingProvisioningState `json:"provisioningState,omitempty"`
}
