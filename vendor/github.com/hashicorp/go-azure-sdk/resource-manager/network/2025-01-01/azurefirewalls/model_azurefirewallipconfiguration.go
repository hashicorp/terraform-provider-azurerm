package azurefirewalls

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AzureFirewallIPConfiguration struct {
	Etag       *string                                       `json:"etag,omitempty"`
	Id         *string                                       `json:"id,omitempty"`
	Name       *string                                       `json:"name,omitempty"`
	Properties *AzureFirewallIPConfigurationPropertiesFormat `json:"properties,omitempty"`
	Type       *string                                       `json:"type,omitempty"`
}
