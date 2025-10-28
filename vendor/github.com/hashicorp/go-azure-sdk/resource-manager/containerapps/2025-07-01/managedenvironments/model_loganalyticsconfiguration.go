package managedenvironments

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type LogAnalyticsConfiguration struct {
	CustomerId *string `json:"customerId,omitempty"`
	SharedKey  *string `json:"sharedKey,omitempty"`
}
