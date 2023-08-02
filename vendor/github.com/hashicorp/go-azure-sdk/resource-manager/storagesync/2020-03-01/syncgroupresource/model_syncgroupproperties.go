package syncgroupresource

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SyncGroupProperties struct {
	SyncGroupStatus *string `json:"syncGroupStatus,omitempty"`
	UniqueId        *string `json:"uniqueId,omitempty"`
}
