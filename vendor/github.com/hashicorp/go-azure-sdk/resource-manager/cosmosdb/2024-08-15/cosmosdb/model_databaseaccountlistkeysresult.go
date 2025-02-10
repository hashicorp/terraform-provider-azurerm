package cosmosdb

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DatabaseAccountListKeysResult struct {
	PrimaryMasterKey           *string `json:"primaryMasterKey,omitempty"`
	PrimaryReadonlyMasterKey   *string `json:"primaryReadonlyMasterKey,omitempty"`
	SecondaryMasterKey         *string `json:"secondaryMasterKey,omitempty"`
	SecondaryReadonlyMasterKey *string `json:"secondaryReadonlyMasterKey,omitempty"`
}
