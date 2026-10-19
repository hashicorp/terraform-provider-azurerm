package virtualwans

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type EffectiveRoutesParameters struct {
	ResourceId             *string `json:"resourceId,omitempty"`
	VirtualWanResourceType *string `json:"virtualWanResourceType,omitempty"`
}
