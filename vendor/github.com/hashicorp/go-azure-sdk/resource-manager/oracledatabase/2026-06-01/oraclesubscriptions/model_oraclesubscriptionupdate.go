package oraclesubscriptions

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type OracleSubscriptionUpdate struct {
	Plan       *PlanUpdate                         `json:"plan,omitempty"`
	Properties *OracleSubscriptionUpdateProperties `json:"properties,omitempty"`
}
