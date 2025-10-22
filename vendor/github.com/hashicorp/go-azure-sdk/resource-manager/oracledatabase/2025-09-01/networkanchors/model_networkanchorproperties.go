package networkanchors

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type NetworkAnchorProperties struct {
	CidrBlock                            *string                         `json:"cidrBlock,omitempty"`
	DnsForwardingEndpointIPAddress       *string                         `json:"dnsForwardingEndpointIpAddress,omitempty"`
	DnsForwardingEndpointNsgRulesURL     *string                         `json:"dnsForwardingEndpointNsgRulesUrl,omitempty"`
	DnsForwardingRules                   *[]DnsForwardingRule            `json:"dnsForwardingRules,omitempty"`
	DnsForwardingRulesURL                *string                         `json:"dnsForwardingRulesUrl,omitempty"`
	DnsListeningEndpointAllowedCidrs     *string                         `json:"dnsListeningEndpointAllowedCidrs,omitempty"`
	DnsListeningEndpointIPAddress        *string                         `json:"dnsListeningEndpointIpAddress,omitempty"`
	DnsListeningEndpointNsgRulesURL      *string                         `json:"dnsListeningEndpointNsgRulesUrl,omitempty"`
	IsOracleDnsForwardingEndpointEnabled *bool                           `json:"isOracleDnsForwardingEndpointEnabled,omitempty"`
	IsOracleDnsListeningEndpointEnabled  *bool                           `json:"isOracleDnsListeningEndpointEnabled,omitempty"`
	IsOracleToAzureDnsZoneSyncEnabled    *bool                           `json:"isOracleToAzureDnsZoneSyncEnabled,omitempty"`
	OciBackupCidrBlock                   *string                         `json:"ociBackupCidrBlock,omitempty"`
	OciSubnetId                          *string                         `json:"ociSubnetId,omitempty"`
	OciVcnDnsLabel                       *string                         `json:"ociVcnDnsLabel,omitempty"`
	OciVcnId                             *string                         `json:"ociVcnId,omitempty"`
	ProvisioningState                    *AzureResourceProvisioningState `json:"provisioningState,omitempty"`
	ResourceAnchorId                     string                          `json:"resourceAnchorId"`
	SubnetId                             string                          `json:"subnetId"`
	VnetId                               *string                         `json:"vnetId,omitempty"`
}
