package machines

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type Feature struct {
	FeatureType *string `json:"featureType,omitempty"`
	Name        *string `json:"name,omitempty"`
	Parent      *string `json:"parent,omitempty"`
	Status      *string `json:"status,omitempty"`
}
