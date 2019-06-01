package automation

import "github.com/Azure/azure-sdk-for-go/services/automation/mgmt/2015-10-31/automation"

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
