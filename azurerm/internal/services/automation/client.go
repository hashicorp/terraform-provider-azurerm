package automation

import (
	"github.com/Azure/azure-sdk-for-go/services/automation/mgmt/2015-10-31/automation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/ar"
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

func BuildClient(endpoint, subscriptionId string, o *ar.ClientOptions) *Client {
	c := Client{}

	c.AccountClient = automation.NewAccountClientWithBaseURI(endpoint, subscriptionId)
	ar.ConfigureClient(&c.AccountClient.Client, o)

	c.AgentRegistrationInfoClient = automation.NewAgentRegistrationInformationClientWithBaseURI(endpoint, subscriptionId)
	ar.ConfigureClient(&c.AgentRegistrationInfoClient.Client, o)

	c.CredentialClient = automation.NewCredentialClientWithBaseURI(endpoint, subscriptionId)
	ar.ConfigureClient(&c.CredentialClient.Client, o)

	c.DscConfigurationClient = automation.NewDscConfigurationClientWithBaseURI(endpoint, subscriptionId)
	ar.ConfigureClient(&c.DscConfigurationClient.Client, o)

	c.DscNodeConfigurationClient = automation.NewDscNodeConfigurationClientWithBaseURI(endpoint, subscriptionId)
	ar.ConfigureClient(&c.DscNodeConfigurationClient.Client, o)

	c.ModuleClient = automation.NewModuleClientWithBaseURI(endpoint, subscriptionId)
	ar.ConfigureClient(&c.ModuleClient.Client, o)

	c.RunbookClient = automation.NewRunbookClientWithBaseURI(endpoint, subscriptionId)
	ar.ConfigureClient(&c.RunbookClient.Client, o)

	c.ScheduleClient = automation.NewScheduleClientWithBaseURI(endpoint, subscriptionId)
	ar.ConfigureClient(&c.ScheduleClient.Client, o)

	c.RunbookDraftClient = automation.NewRunbookDraftClientWithBaseURI(endpoint, subscriptionId)
	ar.ConfigureClient(&c.RunbookDraftClient.Client, o)

	c.VariableClient = automation.NewVariableClientWithBaseURI(endpoint, subscriptionId)
	ar.ConfigureClient(&c.VariableClient.Client, o)

	return &c
}
