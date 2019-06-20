package automation

import (
	"github.com/Azure/azure-sdk-for-go/services/automation/mgmt/2015-10-31/automation"
	"github.com/Azure/go-autorest/autorest"
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

func BuildClient(endpoint string, authorizer autorest.Authorizer, o *common.ClientOptions) *Client {
	c := Client{}

	c.AccountClient = automation.NewAccountClientWithBaseURI(endpoint, o.SubscriptionId)
	o.ConfigureClient(&c.AccountClient.Client, authorizer)

	c.AgentRegistrationInfoClient = automation.NewAgentRegistrationInformationClientWithBaseURI(endpoint, o.SubscriptionId)
	o.ConfigureClient(&c.AgentRegistrationInfoClient.Client, authorizer)

	c.CredentialClient = automation.NewCredentialClientWithBaseURI(endpoint, o.SubscriptionId)
	o.ConfigureClient(&c.CredentialClient.Client, authorizer)

	c.DscConfigurationClient = automation.NewDscConfigurationClientWithBaseURI(endpoint, o.SubscriptionId)
	o.ConfigureClient(&c.DscConfigurationClient.Client, authorizer)

	c.DscNodeConfigurationClient = automation.NewDscNodeConfigurationClientWithBaseURI(endpoint, o.SubscriptionId)
	o.ConfigureClient(&c.DscNodeConfigurationClient.Client, authorizer)

	c.ModuleClient = automation.NewModuleClientWithBaseURI(endpoint, o.SubscriptionId)
	o.ConfigureClient(&c.ModuleClient.Client, authorizer)

	c.RunbookClient = automation.NewRunbookClientWithBaseURI(endpoint, o.SubscriptionId)
	o.ConfigureClient(&c.RunbookClient.Client, authorizer)

	c.ScheduleClient = automation.NewScheduleClientWithBaseURI(endpoint, o.SubscriptionId)
	o.ConfigureClient(&c.ScheduleClient.Client, authorizer)

	c.RunbookDraftClient = automation.NewRunbookDraftClientWithBaseURI(endpoint, o.SubscriptionId)
	o.ConfigureClient(&c.RunbookDraftClient.Client, authorizer)

	c.VariableClient = automation.NewVariableClientWithBaseURI(endpoint, o.SubscriptionId)
	o.ConfigureClient(&c.VariableClient.Client, authorizer)

	return &c
}
