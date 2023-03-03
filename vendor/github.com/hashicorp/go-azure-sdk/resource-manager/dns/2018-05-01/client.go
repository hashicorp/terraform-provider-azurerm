package v2018_05_01

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/dns/2018-05-01/dns"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dns/2018-05-01/recordsets"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dns/2018-05-01/zones"
	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	"github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type Client struct {
	Dns        *dns.DnsClient
	RecordSets *recordsets.RecordSetsClient
	Zones      *zones.ZonesClient
}

func NewClientWithBaseURI(api environments.Api, configureFunc func(c *resourcemanager.Client)) (*Client, error) {
	dnsClient, err := dns.NewDnsClientWithBaseURI(api)
	if err != nil {
		return nil, fmt.Errorf("building Dns client: %+v", err)
	}
	configureFunc(dnsClient.Client)

	recordSetsClient, err := recordsets.NewRecordSetsClientWithBaseURI(api)
	if err != nil {
		return nil, fmt.Errorf("building RecordSets client: %+v", err)
	}
	configureFunc(recordSetsClient.Client)

	zonesClient, err := zones.NewZonesClientWithBaseURI(api)
	if err != nil {
		return nil, fmt.Errorf("building Zones client: %+v", err)
	}
	configureFunc(zonesClient.Client)

	return &Client{
		Dns:        dnsClient,
		RecordSets: recordSetsClient,
		Zones:      zonesClient,
	}, nil
}
