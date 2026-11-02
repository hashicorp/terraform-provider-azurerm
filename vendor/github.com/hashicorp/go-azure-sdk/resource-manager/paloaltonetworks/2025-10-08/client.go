package v2025_10_08

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/paloaltonetworks/2025-10-08/certificateobjectglobalrulestackresources"
	"github.com/hashicorp/go-azure-sdk/resource-manager/paloaltonetworks/2025-10-08/certificateobjectlocalrulestackresources"
	"github.com/hashicorp/go-azure-sdk/resource-manager/paloaltonetworks/2025-10-08/firewallresources"
	"github.com/hashicorp/go-azure-sdk/resource-manager/paloaltonetworks/2025-10-08/firewallstatusresources"
	"github.com/hashicorp/go-azure-sdk/resource-manager/paloaltonetworks/2025-10-08/fqdnlistglobalrulestackresources"
	"github.com/hashicorp/go-azure-sdk/resource-manager/paloaltonetworks/2025-10-08/fqdnlistlocalrulestackresources"
	"github.com/hashicorp/go-azure-sdk/resource-manager/paloaltonetworks/2025-10-08/globalrulestackresources"
	"github.com/hashicorp/go-azure-sdk/resource-manager/paloaltonetworks/2025-10-08/localrulesresources"
	"github.com/hashicorp/go-azure-sdk/resource-manager/paloaltonetworks/2025-10-08/localrulestackresources"
	"github.com/hashicorp/go-azure-sdk/resource-manager/paloaltonetworks/2025-10-08/metricsobjectfirewallresources"
	"github.com/hashicorp/go-azure-sdk/resource-manager/paloaltonetworks/2025-10-08/paloaltonetworkscloudngfws"
	"github.com/hashicorp/go-azure-sdk/resource-manager/paloaltonetworks/2025-10-08/postrulesresources"
	"github.com/hashicorp/go-azure-sdk/resource-manager/paloaltonetworks/2025-10-08/prefixlistglobalrulestackresources"
	"github.com/hashicorp/go-azure-sdk/resource-manager/paloaltonetworks/2025-10-08/prefixlistresources"
	"github.com/hashicorp/go-azure-sdk/resource-manager/paloaltonetworks/2025-10-08/prerulesresources"
	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

type Client struct {
	CertificateObjectGlobalRulestackResources *certificateobjectglobalrulestackresources.CertificateObjectGlobalRulestackResourcesClient
	CertificateObjectLocalRulestackResources  *certificateobjectlocalrulestackresources.CertificateObjectLocalRulestackResourcesClient
	FirewallResources                         *firewallresources.FirewallResourcesClient
	FirewallStatusResources                   *firewallstatusresources.FirewallStatusResourcesClient
	FqdnListGlobalRulestackResources          *fqdnlistglobalrulestackresources.FqdnListGlobalRulestackResourcesClient
	FqdnListLocalRulestackResources           *fqdnlistlocalrulestackresources.FqdnListLocalRulestackResourcesClient
	GlobalRulestackResources                  *globalrulestackresources.GlobalRulestackResourcesClient
	LocalRulesResources                       *localrulesresources.LocalRulesResourcesClient
	LocalRulestackResources                   *localrulestackresources.LocalRulestackResourcesClient
	MetricsObjectFirewallResources            *metricsobjectfirewallresources.MetricsObjectFirewallResourcesClient
	PaloAltoNetworksCloudngfws                *paloaltonetworkscloudngfws.PaloAltoNetworksCloudngfwsClient
	PostRulesResources                        *postrulesresources.PostRulesResourcesClient
	PreRulesResources                         *prerulesresources.PreRulesResourcesClient
	PrefixListGlobalRulestackResources        *prefixlistglobalrulestackresources.PrefixListGlobalRulestackResourcesClient
	PrefixListResources                       *prefixlistresources.PrefixListResourcesClient
}

func NewClientWithBaseURI(sdkApi sdkEnv.Api, configureFunc func(c *resourcemanager.Client)) (*Client, error) {
	certificateObjectGlobalRulestackResourcesClient, err := certificateobjectglobalrulestackresources.NewCertificateObjectGlobalRulestackResourcesClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building CertificateObjectGlobalRulestackResources client: %+v", err)
	}
	configureFunc(certificateObjectGlobalRulestackResourcesClient.Client)

	certificateObjectLocalRulestackResourcesClient, err := certificateobjectlocalrulestackresources.NewCertificateObjectLocalRulestackResourcesClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building CertificateObjectLocalRulestackResources client: %+v", err)
	}
	configureFunc(certificateObjectLocalRulestackResourcesClient.Client)

	firewallResourcesClient, err := firewallresources.NewFirewallResourcesClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building FirewallResources client: %+v", err)
	}
	configureFunc(firewallResourcesClient.Client)

	firewallStatusResourcesClient, err := firewallstatusresources.NewFirewallStatusResourcesClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building FirewallStatusResources client: %+v", err)
	}
	configureFunc(firewallStatusResourcesClient.Client)

	fqdnListGlobalRulestackResourcesClient, err := fqdnlistglobalrulestackresources.NewFqdnListGlobalRulestackResourcesClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building FqdnListGlobalRulestackResources client: %+v", err)
	}
	configureFunc(fqdnListGlobalRulestackResourcesClient.Client)

	fqdnListLocalRulestackResourcesClient, err := fqdnlistlocalrulestackresources.NewFqdnListLocalRulestackResourcesClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building FqdnListLocalRulestackResources client: %+v", err)
	}
	configureFunc(fqdnListLocalRulestackResourcesClient.Client)

	globalRulestackResourcesClient, err := globalrulestackresources.NewGlobalRulestackResourcesClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building GlobalRulestackResources client: %+v", err)
	}
	configureFunc(globalRulestackResourcesClient.Client)

	localRulesResourcesClient, err := localrulesresources.NewLocalRulesResourcesClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building LocalRulesResources client: %+v", err)
	}
	configureFunc(localRulesResourcesClient.Client)

	localRulestackResourcesClient, err := localrulestackresources.NewLocalRulestackResourcesClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building LocalRulestackResources client: %+v", err)
	}
	configureFunc(localRulestackResourcesClient.Client)

	metricsObjectFirewallResourcesClient, err := metricsobjectfirewallresources.NewMetricsObjectFirewallResourcesClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building MetricsObjectFirewallResources client: %+v", err)
	}
	configureFunc(metricsObjectFirewallResourcesClient.Client)

	paloAltoNetworksCloudngfwsClient, err := paloaltonetworkscloudngfws.NewPaloAltoNetworksCloudngfwsClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building PaloAltoNetworksCloudngfws client: %+v", err)
	}
	configureFunc(paloAltoNetworksCloudngfwsClient.Client)

	postRulesResourcesClient, err := postrulesresources.NewPostRulesResourcesClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building PostRulesResources client: %+v", err)
	}
	configureFunc(postRulesResourcesClient.Client)

	preRulesResourcesClient, err := prerulesresources.NewPreRulesResourcesClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building PreRulesResources client: %+v", err)
	}
	configureFunc(preRulesResourcesClient.Client)

	prefixListGlobalRulestackResourcesClient, err := prefixlistglobalrulestackresources.NewPrefixListGlobalRulestackResourcesClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building PrefixListGlobalRulestackResources client: %+v", err)
	}
	configureFunc(prefixListGlobalRulestackResourcesClient.Client)

	prefixListResourcesClient, err := prefixlistresources.NewPrefixListResourcesClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building PrefixListResources client: %+v", err)
	}
	configureFunc(prefixListResourcesClient.Client)

	return &Client{
		CertificateObjectGlobalRulestackResources: certificateObjectGlobalRulestackResourcesClient,
		CertificateObjectLocalRulestackResources:  certificateObjectLocalRulestackResourcesClient,
		FirewallResources:                         firewallResourcesClient,
		FirewallStatusResources:                   firewallStatusResourcesClient,
		FqdnListGlobalRulestackResources:          fqdnListGlobalRulestackResourcesClient,
		FqdnListLocalRulestackResources:           fqdnListLocalRulestackResourcesClient,
		GlobalRulestackResources:                  globalRulestackResourcesClient,
		LocalRulesResources:                       localRulesResourcesClient,
		LocalRulestackResources:                   localRulestackResourcesClient,
		MetricsObjectFirewallResources:            metricsObjectFirewallResourcesClient,
		PaloAltoNetworksCloudngfws:                paloAltoNetworksCloudngfwsClient,
		PostRulesResources:                        postRulesResourcesClient,
		PreRulesResources:                         preRulesResourcesClient,
		PrefixListGlobalRulestackResources:        prefixListGlobalRulestackResourcesClient,
		PrefixListResources:                       prefixListResourcesClient,
	}, nil
}
