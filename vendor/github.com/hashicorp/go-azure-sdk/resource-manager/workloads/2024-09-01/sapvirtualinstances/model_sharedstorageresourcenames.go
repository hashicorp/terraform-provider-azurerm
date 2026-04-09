package sapvirtualinstances

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SharedStorageResourceNames struct {
	SharedStorageAccountName                *string `json:"sharedStorageAccountName,omitempty"`
	SharedStorageAccountPrivateEndPointName *string `json:"sharedStorageAccountPrivateEndPointName,omitempty"`
}
