package networkinterfaces

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type EffectiveRouteListResult struct {
	NextLink *string           `json:"nextLink,omitempty"`
	Value    *[]EffectiveRoute `json:"value,omitempty"`
}
