package v2023_07_01_preview

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/dns/2023-07-01-preview/dns"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dns/2023-07-01-preview/dnssecconfigs"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dns/2023-07-01-preview/recordsets"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dns/2023-07-01-preview/zones"
	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

type Client struct {
	Dns           *dns.DnsClient
	DnssecConfigs *dnssecconfigs.DnssecConfigsClient
	RecordSets    *recordsets.RecordSetsClient
	Zones         *zones.ZonesClient
}

func NewClientWithBaseURI(sdkApi sdkEnv.Api, configureFunc func(c *resourcemanager.Client)) (*Client, error) {
	dnsClient, err := dns.NewDnsClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building Dns client: %+v", err)
	}
	configureFunc(dnsClient.Client)

	dnssecConfigsClient, err := dnssecconfigs.NewDnssecConfigsClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building DnssecConfigs client: %+v", err)
	}
	configureFunc(dnssecConfigsClient.Client)

	recordSetsClient, err := recordsets.NewRecordSetsClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building RecordSets client: %+v", err)
	}
	configureFunc(recordSetsClient.Client)

	zonesClient, err := zones.NewZonesClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building Zones client: %+v", err)
	}
	configureFunc(zonesClient.Client)

	return &Client{
		Dns:           dnsClient,
		DnssecConfigs: dnssecConfigsClient,
		RecordSets:    recordSetsClient,
		Zones:         zonesClient,
	}, nil
}
