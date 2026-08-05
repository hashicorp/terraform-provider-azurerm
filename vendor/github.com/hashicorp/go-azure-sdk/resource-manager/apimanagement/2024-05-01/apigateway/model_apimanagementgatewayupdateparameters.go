package apigateway

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ApiManagementGatewayUpdateParameters struct {
	Etag       *string                                    `json:"etag,omitempty"`
	Id         *string                                    `json:"id,omitempty"`
	Name       *string                                    `json:"name,omitempty"`
	Properties *ApiManagementGatewayBaseProperties        `json:"properties,omitempty"`
	Sku        *ApiManagementGatewaySkuPropertiesForPatch `json:"sku,omitempty"`
	Tags       *map[string]string                         `json:"tags,omitempty"`
	Type       *string                                    `json:"type,omitempty"`
}
