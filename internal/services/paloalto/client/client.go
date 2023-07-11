package client

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/paloaltonetworks/2022-08-29/certificateobjectlocalrulestack"
	"github.com/hashicorp/go-azure-sdk/resource-manager/paloaltonetworks/2022-08-29/firewalls"
	"github.com/hashicorp/go-azure-sdk/resource-manager/paloaltonetworks/2022-08-29/fqdnlistlocalrulestack"
	"github.com/hashicorp/go-azure-sdk/resource-manager/paloaltonetworks/2022-08-29/localrules"
	"github.com/hashicorp/go-azure-sdk/resource-manager/paloaltonetworks/2022-08-29/localrulestacks"
	"github.com/hashicorp/go-azure-sdk/resource-manager/paloaltonetworks/2022-08-29/prefixlistlocalrulestack"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	CertificatesClient    *certificateobjectlocalrulestack.CertificateObjectLocalRulestackClient
	FirewallClient        *firewalls.FirewallsClient
	FQDNListsClient       *fqdnlistlocalrulestack.FqdnListLocalRulestackClient
	LocalRulesClient      *localrules.LocalRulesClient
	LocalRuleStacksClient *localrulestacks.LocalRuleStacksClient
	PrefixListClient      *prefixlistlocalrulestack.PrefixListLocalRulestackClient
}

func NewClient(o *common.ClientOptions) (*Client, error) {
	certificatesClient, err := certificateobjectlocalrulestack.NewCertificateObjectLocalRulestackClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Firewall client: %+v", err)
	}
	o.Configure(certificatesClient.Client, o.Authorizers.ResourceManager)

	firewallClient, err := firewalls.NewFirewallsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Firewall client: %+v", err)
	}
	o.Configure(firewallClient.Client, o.Authorizers.ResourceManager)

	fqdnListsClient, err := fqdnlistlocalrulestack.NewFqdnListLocalRulestackClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Firewall client: %+v", err)
	}
	o.Configure(fqdnListsClient.Client, o.Authorizers.ResourceManager)

	localRuleStackClient, err := localrulestacks.NewLocalRuleStacksClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building LocalRuleStacks client: %+v", err)
	}
	o.Configure(localRuleStackClient.Client, o.Authorizers.ResourceManager)

	localRuleClient, err := localrules.NewLocalRulesClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building LocalRules client: %+v", err)
	}
	o.Configure(localRuleClient.Client, o.Authorizers.ResourceManager)

	prefixListClient, err := prefixlistlocalrulestack.NewPrefixListLocalRulestackClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building LocalRules client: %+v", err)
	}
	o.Configure(prefixListClient.Client, o.Authorizers.ResourceManager)

	return &Client{
		CertificatesClient:    certificatesClient,
		FirewallClient:        firewallClient,
		FQDNListsClient:       fqdnListsClient,
		LocalRulesClient:      localRuleClient,
		LocalRuleStacksClient: localRuleStackClient,
		PrefixListClient:      prefixListClient,
	}, nil
}
