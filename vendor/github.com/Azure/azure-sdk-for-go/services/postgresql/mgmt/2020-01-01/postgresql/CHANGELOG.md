# Change History

## Additive Changes

### New Funcs

1. *RecoverableServerResource.UnmarshalJSON([]byte) error
1. *ServerParametersListUpdateConfigurationsFuture.UnmarshalJSON([]byte) error
1. *ServerSecurityAlertPolicyListResultIterator.Next() error
1. *ServerSecurityAlertPolicyListResultIterator.NextWithContext(context.Context) error
1. *ServerSecurityAlertPolicyListResultPage.Next() error
1. *ServerSecurityAlertPolicyListResultPage.NextWithContext(context.Context) error
1. NewRecoverableServersClient(string) RecoverableServersClient
1. NewRecoverableServersClientWithBaseURI(string, string) RecoverableServersClient
1. NewServerBasedPerformanceTierClient(string) ServerBasedPerformanceTierClient
1. NewServerBasedPerformanceTierClientWithBaseURI(string, string) ServerBasedPerformanceTierClient
1. NewServerParametersClient(string) ServerParametersClient
1. NewServerParametersClientWithBaseURI(string, string) ServerParametersClient
1. NewServerSecurityAlertPolicyListResultIterator(ServerSecurityAlertPolicyListResultPage) ServerSecurityAlertPolicyListResultIterator
1. NewServerSecurityAlertPolicyListResultPage(ServerSecurityAlertPolicyListResult, func(context.Context, ServerSecurityAlertPolicyListResult) (ServerSecurityAlertPolicyListResult, error)) ServerSecurityAlertPolicyListResultPage
1. RecoverableServerProperties.MarshalJSON() ([]byte, error)
1. RecoverableServerResource.MarshalJSON() ([]byte, error)
1. RecoverableServersClient.Get(context.Context, string, string) (RecoverableServerResource, error)
1. RecoverableServersClient.GetPreparer(context.Context, string, string) (*http.Request, error)
1. RecoverableServersClient.GetResponder(*http.Response) (RecoverableServerResource, error)
1. RecoverableServersClient.GetSender(*http.Request) (*http.Response, error)
1. ServerBasedPerformanceTierClient.List(context.Context, string, string) (PerformanceTierListResult, error)
1. ServerBasedPerformanceTierClient.ListPreparer(context.Context, string, string) (*http.Request, error)
1. ServerBasedPerformanceTierClient.ListResponder(*http.Response) (PerformanceTierListResult, error)
1. ServerBasedPerformanceTierClient.ListSender(*http.Request) (*http.Response, error)
1. ServerParametersClient.ListUpdateConfigurations(context.Context, string, string, ConfigurationListResult) (ServerParametersListUpdateConfigurationsFuture, error)
1. ServerParametersClient.ListUpdateConfigurationsPreparer(context.Context, string, string, ConfigurationListResult) (*http.Request, error)
1. ServerParametersClient.ListUpdateConfigurationsResponder(*http.Response) (ConfigurationListResult, error)
1. ServerParametersClient.ListUpdateConfigurationsSender(*http.Request) (ServerParametersListUpdateConfigurationsFuture, error)
1. ServerSecurityAlertPoliciesClient.ListByServer(context.Context, string, string) (ServerSecurityAlertPolicyListResultPage, error)
1. ServerSecurityAlertPoliciesClient.ListByServerComplete(context.Context, string, string) (ServerSecurityAlertPolicyListResultIterator, error)
1. ServerSecurityAlertPoliciesClient.ListByServerPreparer(context.Context, string, string) (*http.Request, error)
1. ServerSecurityAlertPoliciesClient.ListByServerResponder(*http.Response) (ServerSecurityAlertPolicyListResult, error)
1. ServerSecurityAlertPoliciesClient.ListByServerSender(*http.Request) (*http.Response, error)
1. ServerSecurityAlertPolicyListResult.IsEmpty() bool
1. ServerSecurityAlertPolicyListResult.MarshalJSON() ([]byte, error)
1. ServerSecurityAlertPolicyListResultIterator.NotDone() bool
1. ServerSecurityAlertPolicyListResultIterator.Response() ServerSecurityAlertPolicyListResult
1. ServerSecurityAlertPolicyListResultIterator.Value() ServerSecurityAlertPolicy
1. ServerSecurityAlertPolicyListResultPage.NotDone() bool
1. ServerSecurityAlertPolicyListResultPage.Response() ServerSecurityAlertPolicyListResult
1. ServerSecurityAlertPolicyListResultPage.Values() []ServerSecurityAlertPolicy

### Struct Changes

#### New Structs

1. RecoverableServerProperties
1. RecoverableServerResource
1. RecoverableServersClient
1. ServerBasedPerformanceTierClient
1. ServerParametersClient
1. ServerParametersListUpdateConfigurationsFuture
1. ServerSecurityAlertPolicyListResult
1. ServerSecurityAlertPolicyListResultIterator
1. ServerSecurityAlertPolicyListResultPage

#### New Struct Fields

1. PerformanceTierProperties.MaxBackupRetentionDays
1. PerformanceTierProperties.MaxLargeStorageMB
1. PerformanceTierProperties.MaxStorageMB
1. PerformanceTierProperties.MinBackupRetentionDays
1. PerformanceTierProperties.MinLargeStorageMB
1. PerformanceTierProperties.MinStorageMB
