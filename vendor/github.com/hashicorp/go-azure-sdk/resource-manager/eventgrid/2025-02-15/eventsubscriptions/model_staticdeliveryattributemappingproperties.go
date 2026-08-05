package eventsubscriptions

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type StaticDeliveryAttributeMappingProperties struct {
	IsSecret *bool   `json:"isSecret,omitempty"`
	Value    *string `json:"value,omitempty"`
}
