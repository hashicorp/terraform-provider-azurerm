package v2022_08_29

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/paloaltonetworks/2022-08-29/certificateobjectglobalrulestack"
	"github.com/hashicorp/go-azure-sdk/resource-manager/paloaltonetworks/2022-08-29/certificateobjectlocalrulestack"
	"github.com/hashicorp/go-azure-sdk/resource-manager/paloaltonetworks/2022-08-29/firewalls"
	"github.com/hashicorp/go-azure-sdk/resource-manager/paloaltonetworks/2022-08-29/firewallstatus"
	"github.com/hashicorp/go-azure-sdk/resource-manager/paloaltonetworks/2022-08-29/fqdnlistglobalrulestack"
	"github.com/hashicorp/go-azure-sdk/resource-manager/paloaltonetworks/2022-08-29/fqdnlistlocalrulestack"
	"github.com/hashicorp/go-azure-sdk/resource-manager/paloaltonetworks/2022-08-29/globalrulestack"
	"github.com/hashicorp/go-azure-sdk/resource-manager/paloaltonetworks/2022-08-29/localrules"
	"github.com/hashicorp/go-azure-sdk/resource-manager/paloaltonetworks/2022-08-29/localrulestacks"
	"github.com/hashicorp/go-azure-sdk/resource-manager/paloaltonetworks/2022-08-29/postrules"
	"github.com/hashicorp/go-azure-sdk/resource-manager/paloaltonetworks/2022-08-29/prefixlistglobalrulestack"
	"github.com/hashicorp/go-azure-sdk/resource-manager/paloaltonetworks/2022-08-29/prefixlistlocalrulestack"
	"github.com/hashicorp/go-azure-sdk/resource-manager/paloaltonetworks/2022-08-29/prerules"
	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	"github.com/hashicorp/go-azure-sdk/sdk/environments"
)

type Client struct {
	CertificateObjectGlobalRulestack *certificateobjectglobalrulestack.CertificateObjectGlobalRulestackClient
	CertificateObjectLocalRulestack  *certificateobjectlocalrulestack.CertificateObjectLocalRulestackClient
	FirewallStatus                   *firewallstatus.FirewallStatusClient
	Firewalls                        *firewalls.FirewallsClient
	FqdnListGlobalRulestack          *fqdnlistglobalrulestack.FqdnListGlobalRulestackClient
	FqdnListLocalRulestack           *fqdnlistlocalrulestack.FqdnListLocalRulestackClient
	GlobalRulestack                  *globalrulestack.GlobalRulestackClient
	LocalRules                       *localrules.LocalRulesClient
	LocalRulestacks                  *localrulestacks.LocalRulestacksClient
	PostRules                        *postrules.PostRulesClient
	PreRules                         *prerules.PreRulesClient
	PrefixListGlobalRulestack        *prefixlistglobalrulestack.PrefixListGlobalRulestackClient
	PrefixListLocalRulestack         *prefixlistlocalrulestack.PrefixListLocalRulestackClient
}

func NewClientWithBaseURI(api environments.Api, configureFunc func(c *resourcemanager.Client)) (*Client, error) {
	certificateObjectGlobalRulestackClient, err := certificateobjectglobalrulestack.NewCertificateObjectGlobalRulestackClientWithBaseURI(api)
	if err != nil {
		return nil, fmt.Errorf("building CertificateObjectGlobalRulestack client: %+v", err)
	}
	configureFunc(certificateObjectGlobalRulestackClient.Client)

	certificateObjectLocalRulestackClient, err := certificateobjectlocalrulestack.NewCertificateObjectLocalRulestackClientWithBaseURI(api)
	if err != nil {
		return nil, fmt.Errorf("building CertificateObjectLocalRulestack client: %+v", err)
	}
	configureFunc(certificateObjectLocalRulestackClient.Client)

	firewallStatusClient, err := firewallstatus.NewFirewallStatusClientWithBaseURI(api)
	if err != nil {
		return nil, fmt.Errorf("building FirewallStatus client: %+v", err)
	}
	configureFunc(firewallStatusClient.Client)

	firewallsClient, err := firewalls.NewFirewallsClientWithBaseURI(api)
	if err != nil {
		return nil, fmt.Errorf("building Firewalls client: %+v", err)
	}
	configureFunc(firewallsClient.Client)

	fqdnListGlobalRulestackClient, err := fqdnlistglobalrulestack.NewFqdnListGlobalRulestackClientWithBaseURI(api)
	if err != nil {
		return nil, fmt.Errorf("building FqdnListGlobalRulestack client: %+v", err)
	}
	configureFunc(fqdnListGlobalRulestackClient.Client)

	fqdnListLocalRulestackClient, err := fqdnlistlocalrulestack.NewFqdnListLocalRulestackClientWithBaseURI(api)
	if err != nil {
		return nil, fmt.Errorf("building FqdnListLocalRulestack client: %+v", err)
	}
	configureFunc(fqdnListLocalRulestackClient.Client)

	globalRulestackClient, err := globalrulestack.NewGlobalRulestackClientWithBaseURI(api)
	if err != nil {
		return nil, fmt.Errorf("building GlobalRulestack client: %+v", err)
	}
	configureFunc(globalRulestackClient.Client)

	localRulesClient, err := localrules.NewLocalRulesClientWithBaseURI(api)
	if err != nil {
		return nil, fmt.Errorf("building LocalRules client: %+v", err)
	}
	configureFunc(localRulesClient.Client)

	localRulestacksClient, err := localrulestacks.NewLocalRulestacksClientWithBaseURI(api)
	if err != nil {
		return nil, fmt.Errorf("building LocalRulestacks client: %+v", err)
	}
	configureFunc(localRulestacksClient.Client)

	postRulesClient, err := postrules.NewPostRulesClientWithBaseURI(api)
	if err != nil {
		return nil, fmt.Errorf("building PostRules client: %+v", err)
	}
	configureFunc(postRulesClient.Client)

	preRulesClient, err := prerules.NewPreRulesClientWithBaseURI(api)
	if err != nil {
		return nil, fmt.Errorf("building PreRules client: %+v", err)
	}
	configureFunc(preRulesClient.Client)

	prefixListGlobalRulestackClient, err := prefixlistglobalrulestack.NewPrefixListGlobalRulestackClientWithBaseURI(api)
	if err != nil {
		return nil, fmt.Errorf("building PrefixListGlobalRulestack client: %+v", err)
	}
	configureFunc(prefixListGlobalRulestackClient.Client)

	prefixListLocalRulestackClient, err := prefixlistlocalrulestack.NewPrefixListLocalRulestackClientWithBaseURI(api)
	if err != nil {
		return nil, fmt.Errorf("building PrefixListLocalRulestack client: %+v", err)
	}
	configureFunc(prefixListLocalRulestackClient.Client)

	return &Client{
		CertificateObjectGlobalRulestack: certificateObjectGlobalRulestackClient,
		CertificateObjectLocalRulestack:  certificateObjectLocalRulestackClient,
		FirewallStatus:                   firewallStatusClient,
		Firewalls:                        firewallsClient,
		FqdnListGlobalRulestack:          fqdnListGlobalRulestackClient,
		FqdnListLocalRulestack:           fqdnListLocalRulestackClient,
		GlobalRulestack:                  globalRulestackClient,
		LocalRules:                       localRulesClient,
		LocalRulestacks:                  localRulestacksClient,
		PostRules:                        postRulesClient,
		PreRules:                         preRulesClient,
		PrefixListGlobalRulestack:        prefixListGlobalRulestackClient,
		PrefixListLocalRulestack:         prefixListLocalRulestackClient,
	}, nil
}
