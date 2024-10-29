package subscriptiondiagnosticsettings

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SubscriptionLogSettings struct {
	Category      *string `json:"category,omitempty"`
	CategoryGroup *string `json:"categoryGroup,omitempty"`
	Enabled       bool    `json:"enabled"`
}
