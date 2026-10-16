package oraclesubscriptions

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type OracleSubscriptionUpdateProperties struct {
	Intent      *Intent `json:"intent,omitempty"`
	ProductCode *string `json:"productCode,omitempty"`
}
