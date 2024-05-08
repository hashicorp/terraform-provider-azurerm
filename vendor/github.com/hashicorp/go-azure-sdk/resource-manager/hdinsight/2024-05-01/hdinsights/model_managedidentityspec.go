package hdinsights

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ManagedIdentitySpec struct {
	ClientId   string              `json:"clientId"`
	ObjectId   string              `json:"objectId"`
	ResourceId string              `json:"resourceId"`
	Type       ManagedIdentityType `json:"type"`
}
