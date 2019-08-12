package automation

import (
	"github.com/Azure/azure-sdk-for-go/services/automation/mgmt/2015-10-31/automation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/common"
)

type Client struct {
	AccountClient               automation.AccountClient
	AgentRegistrationInfoClient automation.AgentRegistrationInformationClient
	CredentialClient            automation.CredentialClient
	DscConfigurationClient      automation.DscConfigurationClient
	DscNodeConfigurationClient  automation.DscNodeConfigurationClient
	ModuleClient                automation.ModuleClient
	RunbookClient               automation.RunbookClient
	RunbookDraftClient          automation.RunbookDraftClient
	ScheduleClient              automation.ScheduleClient
	VariableClient              automation.VariableClient
}

func BuildClient(o *common.ClientOptions) *Client {
	c := Client{}

	c.AccountClient = automation.NewAccountClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&c.AccountClient.Client, o.ResourceManagerAuthorizer)

	c.AgentRegistrationInfoClient = automation.NewAgentRegistrationInformationClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&c.AgentRegistrationInfoClient.Client, o.ResourceManagerAuthorizer)

	c.CredentialClient = automation.NewCredentialClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&c.CredentialClient.Client, o.ResourceManagerAuthorizer)

	c.DscConfigurationClient = automation.NewDscConfigurationClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&c.DscConfigurationClient.Client, o.ResourceManagerAuthorizer)

	c.DscNodeConfigurationClient = automation.NewDscNodeConfigurationClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&c.DscNodeConfigurationClient.Client, o.ResourceManagerAuthorizer)

	c.ModuleClient = automation.NewModuleClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&c.ModuleClient.Client, o.ResourceManagerAuthorizer)

	c.RunbookClient = automation.NewRunbookClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&c.RunbookClient.Client, o.ResourceManagerAuthorizer)

	c.ScheduleClient = automation.NewScheduleClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&c.ScheduleClient.Client, o.ResourceManagerAuthorizer)

	c.RunbookDraftClient = automation.NewRunbookDraftClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&c.RunbookDraftClient.Client, o.ResourceManagerAuthorizer)

	c.VariableClient = automation.NewVariableClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&c.VariableClient.Client, o.ResourceManagerAuthorizer)

	return &c
}
