package client

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/paloaltonetworks/2022-08-29/firewalls"
	"github.com/hashicorp/go-azure-sdk/resource-manager/paloaltonetworks/2022-08-29/localrules"
	"github.com/hashicorp/go-azure-sdk/resource-manager/paloaltonetworks/2022-08-29/localrulestacks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	FirewallClient        *firewalls.FirewallsClient
	LocalRulesClient      *localrules.LocalRulesClient
	LocalRuleStacksClient *localrulestacks.LocalRuleStacksClient
}

func NewClient(o *common.ClientOptions) (*Client, error) {
	firewallClient, err := firewalls.NewFirewallsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Firewall client: %+v", err)
	}
	o.Configure(firewallClient.Client, o.Authorizers.ResourceManager)

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

	return &Client{
		FirewallClient:        firewallClient,
		LocalRulesClient:      localRuleClient,
		LocalRuleStacksClient: localRuleStackClient,
	}, nil
}
