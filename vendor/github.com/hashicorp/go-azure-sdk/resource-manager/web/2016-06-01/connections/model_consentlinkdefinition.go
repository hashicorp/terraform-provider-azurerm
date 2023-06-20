package connections

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ConsentLinkDefinition struct {
	DisplayName        *string    `json:"displayName,omitempty"`
	FirstPartyLoginUri *string    `json:"firstPartyLoginUri,omitempty"`
	Link               *string    `json:"link,omitempty"`
	Status             *LinkState `json:"status,omitempty"`
}
