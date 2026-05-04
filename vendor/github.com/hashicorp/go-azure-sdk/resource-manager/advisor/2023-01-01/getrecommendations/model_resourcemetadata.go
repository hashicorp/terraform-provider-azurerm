package getrecommendations

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ResourceMetadata struct {
	Action     *map[string]interface{} `json:"action,omitempty"`
	Plural     *string                 `json:"plural,omitempty"`
	ResourceId *string                 `json:"resourceId,omitempty"`
	Singular   *string                 `json:"singular,omitempty"`
	Source     *string                 `json:"source,omitempty"`
}
