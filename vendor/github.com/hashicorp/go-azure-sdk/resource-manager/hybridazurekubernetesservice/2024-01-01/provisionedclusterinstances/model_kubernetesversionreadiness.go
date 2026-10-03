package provisionedclusterinstances

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type KubernetesVersionReadiness struct {
	ErrorMessage *string `json:"errorMessage,omitempty"`
	OsSku        *OSSKU  `json:"osSku,omitempty"`
	OsType       *OsType `json:"osType,omitempty"`
	Ready        *bool   `json:"ready,omitempty"`
}
