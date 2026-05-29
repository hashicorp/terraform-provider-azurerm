package cosmosdb

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DatabaseAccountKeysMetadata struct {
	PrimaryMasterKey           *AccountKeyMetadata `json:"primaryMasterKey,omitempty"`
	PrimaryReadonlyMasterKey   *AccountKeyMetadata `json:"primaryReadonlyMasterKey,omitempty"`
	SecondaryMasterKey         *AccountKeyMetadata `json:"secondaryMasterKey,omitempty"`
	SecondaryReadonlyMasterKey *AccountKeyMetadata `json:"secondaryReadonlyMasterKey,omitempty"`
}
