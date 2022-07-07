package servers

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AnalysisServicesServer struct {
	Id         *string                           `json:"id,omitempty"`
	Location   string                            `json:"location"`
	Name       *string                           `json:"name,omitempty"`
	Properties *AnalysisServicesServerProperties `json:"properties,omitempty"`
	Sku        ResourceSku                       `json:"sku"`
	Tags       *map[string]string                `json:"tags,omitempty"`
	Type       *string                           `json:"type,omitempty"`
}
