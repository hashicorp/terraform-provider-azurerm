package automation

import (
	"github.com/Azure/azure-sdk-for-go/services/automation/mgmt/2015-10-31/automation"
	"github.com/Azure/go-autorest/autorest"
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

func BuildClients(endpoint, subscriptionId, partnerId string, auth autorest.Authorizer) *Client {
	c := Client{}
	c.AccountClient = automation.NewAccountClientWithBaseURI(endpoint, subscriptionId)
	ar.ConfigureClient(&c.AccountClient.Client, auth, partnerId)

	c.AgentRegistrationInfoClient = automation.NewAgentRegistrationInformationClientWithBaseURI(endpoint, subscriptionId)
	ar.ConfigureClient(&c.AgentRegistrationInfoClient.Client, auth, partnerId)

	c.CredentialClient = automation.NewCredentialClientWithBaseURI(endpoint, subscriptionId)
	ar.ConfigureClient(&c.CredentialClient.Client, auth, partnerId)

	c.DscConfigurationClient = automation.NewDscConfigurationClientWithBaseURI(endpoint, subscriptionId)
	ar.ConfigureClient(&c.DscConfigurationClient.Client, auth, partnerId)

	c.DscNodeConfigurationClient = automation.NewDscNodeConfigurationClientWithBaseURI(endpoint, subscriptionId)
	ar.ConfigureClient(&c.DscNodeConfigurationClient.Client, auth, partnerId)

	c.ModuleClient = automation.NewModuleClientWithBaseURI(endpoint, subscriptionId)
	ar.ConfigureClient(&c.ModuleClient.Client, auth, partnerId)

	c.RunbookClient = automation.NewRunbookClientWithBaseURI(endpoint, subscriptionId)
	ar.ConfigureClient(&c.RunbookClient.Client, auth, partnerId)

	c.ScheduleClient = automation.NewScheduleClientWithBaseURI(endpoint, subscriptionId)
	ar.ConfigureClient(&c.ScheduleClient.Client, auth, partnerId)

	c.RunbookDraftClient = automation.NewRunbookDraftClientWithBaseURI(endpoint, subscriptionId)
	ar.ConfigureClient(&c.RunbookDraftClient.Client, auth, partnerId)

	c.VariableClient = automation.NewVariableClientWithBaseURI(endpoint, subscriptionId)
	ar.ConfigureClient(&c.VariableClient.Client, auth, partnerId)

	return &c
}
