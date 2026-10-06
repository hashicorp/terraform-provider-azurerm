package oraclesubscriptions

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type Plan struct {
	Name          string  `json:"name"`
	Product       string  `json:"product"`
	PromotionCode *string `json:"promotionCode,omitempty"`
	Publisher     string  `json:"publisher"`
	Version       *string `json:"version,omitempty"`
}
