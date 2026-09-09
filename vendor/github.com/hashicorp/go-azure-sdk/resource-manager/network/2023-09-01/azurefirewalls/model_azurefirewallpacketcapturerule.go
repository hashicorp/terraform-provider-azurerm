package azurefirewalls

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AzureFirewallPacketCaptureRule struct {
	DestinationPorts *[]string `json:"destinationPorts,omitempty"`
	Destinations     *[]string `json:"destinations,omitempty"`
	Sources          *[]string `json:"sources,omitempty"`
}
