// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2015-10-31/webhook"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2019-06-01/agentregistrationinformation"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2019-06-01/softwareupdateconfiguration"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2020-01-13-preview/watcher"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2023-11-01/module"
	automation_2024_10_23 "github.com/hashicorp/go-azure-sdk/resource-manager/automation/2024-10-23"
	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	*automation_2024_10_23.Client

	// Legacy module client for PowerShell72Module resources, which is deprecated in the latest API version
	// Should look into deprecate PowerShell72Module resource
	ModuleClientV2023 *module.ModuleClient

	AgentRegistrationInfoClient *agentregistrationinformation.AgentRegistrationInformationClient
	SoftwareUpdateConfigClient  *softwareupdateconfiguration.SoftwareUpdateConfigurationClient
	WebhookClient               *webhook.WebhookClient
	WatcherClient               *watcher.WatcherClient
}

func NewClient(o *common.ClientOptions) (*Client, error) {
	// Latest version client for most resources
	metaClient, err := automation_2024_10_23.NewClientWithBaseURI(o.Environment.ResourceManager, func(c *resourcemanager.Client) {
		o.Configure(c, o.Authorizers.ResourceManager)
	})
	if err != nil {
		return nil, fmt.Errorf("building Automation client: %+v", err)
	}

	// Legacy module client for PowerShell72Module resources
	moduleClientV2023, err := module.NewModuleClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Module v2023 client: %+v", err)
	}
	o.Configure(moduleClientV2023.Client, o.Authorizers.ResourceManager)

	agentRegistrationInfoClient, err := agentregistrationinformation.NewAgentRegistrationInformationClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Agent Registration Info client : %+v", err)
	}
	o.Configure(agentRegistrationInfoClient.Client, o.Authorizers.ResourceManager)

	softUpClient, err := softwareupdateconfiguration.NewSoftwareUpdateConfigurationClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Soft Up client : %+v", err)
	}
	o.Configure(softUpClient.Client, o.Authorizers.ResourceManager)

	watcherClient, err := watcher.NewWatcherClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Watcher client : %+v", err)
	}
	o.Configure(watcherClient.Client, o.Authorizers.ResourceManager)

	webhookClient, err := webhook.NewWebhookClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Webhook client : %+v", err)
	}
	o.Configure(webhookClient.Client, o.Authorizers.ResourceManager)

	return &Client{
		Client:            metaClient,
		ModuleClientV2023: moduleClientV2023,

		AgentRegistrationInfoClient: agentRegistrationInfoClient,
		SoftwareUpdateConfigClient:  softUpClient,
		WatcherClient:               watcherClient,
		WebhookClient:               webhookClient,
	}, nil
}
