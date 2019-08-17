package automation

import (
	"github.com/Azure/azure-sdk-for-go/services/automation/mgmt/2015-10-31/automation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/common"
)

type Client struct {
	AccountClient               *automation.AccountClient
	AgentRegistrationInfoClient *automation.AgentRegistrationInformationClient
	CredentialClient            *automation.CredentialClient
	DscConfigurationClient      *automation.DscConfigurationClient
	DscNodeConfigurationClient  *automation.DscNodeConfigurationClient
	ModuleClient                *automation.ModuleClient
	RunbookClient               *automation.RunbookClient
	RunbookDraftClient          *automation.RunbookDraftClient
	ScheduleClient              *automation.ScheduleClient
	VariableClient              *automation.VariableClient
}

func BuildClient(o *common.ClientOptions) *Client {

	AccountClient := automation.NewAccountClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&AccountClient.Client, o.ResourceManagerAuthorizer)

	AgentRegistrationInfoClient := automation.NewAgentRegistrationInformationClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&AgentRegistrationInfoClient.Client, o.ResourceManagerAuthorizer)

	CredentialClient := automation.NewCredentialClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&CredentialClient.Client, o.ResourceManagerAuthorizer)

	DscConfigurationClient := automation.NewDscConfigurationClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&DscConfigurationClient.Client, o.ResourceManagerAuthorizer)

	DscNodeConfigurationClient := automation.NewDscNodeConfigurationClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&DscNodeConfigurationClient.Client, o.ResourceManagerAuthorizer)

	ModuleClient := automation.NewModuleClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&ModuleClient.Client, o.ResourceManagerAuthorizer)

	RunbookClient := automation.NewRunbookClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&RunbookClient.Client, o.ResourceManagerAuthorizer)

	ScheduleClient := automation.NewScheduleClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&ScheduleClient.Client, o.ResourceManagerAuthorizer)

	RunbookDraftClient := automation.NewRunbookDraftClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&RunbookDraftClient.Client, o.ResourceManagerAuthorizer)

	VariableClient := automation.NewVariableClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&VariableClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		AccountClient:               &AccountClient,
		AgentRegistrationInfoClient: &AgentRegistrationInfoClient,
		CredentialClient:            &CredentialClient,
		DscConfigurationClient:      &DscConfigurationClient,
		DscNodeConfigurationClient:  &DscNodeConfigurationClient,
		ModuleClient:                &ModuleClient,
		RunbookClient:               &RunbookClient,
		RunbookDraftClient:          &RunbookDraftClient,
		ScheduleClient:              &ScheduleClient,
		VariableClient:              &VariableClient,
	}
}
