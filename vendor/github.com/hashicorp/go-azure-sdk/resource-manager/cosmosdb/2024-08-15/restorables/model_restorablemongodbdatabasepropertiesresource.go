package restorables

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RestorableMongodbDatabasePropertiesResource struct {
	CanUndelete       *string        `json:"canUndelete,omitempty"`
	CanUndeleteReason *string        `json:"canUndeleteReason,omitempty"`
	EventTimestamp    *string        `json:"eventTimestamp,omitempty"`
	OperationType     *OperationType `json:"operationType,omitempty"`
	OwnerId           *string        `json:"ownerId,omitempty"`
	OwnerResourceId   *string        `json:"ownerResourceId,omitempty"`
	Rid               *string        `json:"_rid,omitempty"`
}
