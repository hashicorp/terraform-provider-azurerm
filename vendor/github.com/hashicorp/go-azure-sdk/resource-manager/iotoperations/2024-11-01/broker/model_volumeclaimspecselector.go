package broker

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type VolumeClaimSpecSelector struct {
	MatchExpressions *[]VolumeClaimSpecSelectorMatchExpressions `json:"matchExpressions,omitempty"`
	MatchLabels      *map[string]string                         `json:"matchLabels,omitempty"`
}
