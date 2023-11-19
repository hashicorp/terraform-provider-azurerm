package appplatform

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type Scale struct {
	MaxReplicas *int64       `json:"maxReplicas,omitempty"`
	MinReplicas *int64       `json:"minReplicas,omitempty"`
	Rules       *[]ScaleRule `json:"rules,omitempty"`
}
