package servers

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AnalysisServicesServerUpdateParameters struct {
	Properties *AnalysisServicesServerMutableProperties `json:"properties"`
	Sku        *ResourceSku                             `json:"sku"`
	Tags       *map[string]string                       `json:"tags,omitempty"`
}
