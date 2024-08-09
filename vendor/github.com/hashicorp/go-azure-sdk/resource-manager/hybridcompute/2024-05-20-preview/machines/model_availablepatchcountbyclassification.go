package machines

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AvailablePatchCountByClassification struct {
	Critical     *int64 `json:"critical,omitempty"`
	Definition   *int64 `json:"definition,omitempty"`
	FeaturePack  *int64 `json:"featurePack,omitempty"`
	Other        *int64 `json:"other,omitempty"`
	Security     *int64 `json:"security,omitempty"`
	ServicePack  *int64 `json:"servicePack,omitempty"`
	Tools        *int64 `json:"tools,omitempty"`
	UpdateRollup *int64 `json:"updateRollup,omitempty"`
	Updates      *int64 `json:"updates,omitempty"`
}
