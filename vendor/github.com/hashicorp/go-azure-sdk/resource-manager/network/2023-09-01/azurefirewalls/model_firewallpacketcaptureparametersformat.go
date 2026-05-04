package azurefirewalls

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type FirewallPacketCaptureParametersFormat struct {
	DurationInSeconds        *int64                             `json:"durationInSeconds,omitempty"`
	FileName                 *string                            `json:"fileName,omitempty"`
	Filters                  *[]AzureFirewallPacketCaptureRule  `json:"filters,omitempty"`
	Flags                    *[]AzureFirewallPacketCaptureFlags `json:"flags,omitempty"`
	NumberOfPacketsToCapture *int64                             `json:"numberOfPacketsToCapture,omitempty"`
	Protocol                 *AzureFirewallNetworkRuleProtocol  `json:"protocol,omitempty"`
	SasURL                   *string                            `json:"sasUrl,omitempty"`
}
