package subscriptions

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SBClientAffineProperties struct {
	ClientId  *string `json:"clientId,omitempty"`
	IsDurable *bool   `json:"isDurable,omitempty"`
	IsShared  *bool   `json:"isShared,omitempty"`
}
