package videoanalyzers

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type StorageAccount struct {
	Id       *string           `json:"id,omitempty"`
	Identity *ResourceIdentity `json:"identity,omitempty"`
	Status   *string           `json:"status,omitempty"`
}
