package cosmosdb

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SqlTriggerResource struct {
	Body             *string           `json:"body,omitempty"`
	Id               string            `json:"id"`
	TriggerOperation *TriggerOperation `json:"triggerOperation,omitempty"`
	TriggerType      *TriggerType      `json:"triggerType,omitempty"`
}
