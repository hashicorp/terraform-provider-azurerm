# Change History

## Additive Changes

### New Funcs

1. *ExtensionsDisableAzureMonitorFuture.UnmarshalJSON([]byte) error
1. *ExtensionsEnableAzureMonitorFuture.UnmarshalJSON([]byte) error
1. AzureMonitorSelectedConfigurations.MarshalJSON() ([]byte, error)
1. ExtensionsClient.DisableAzureMonitor(context.Context, string, string) (ExtensionsDisableAzureMonitorFuture, error)
1. ExtensionsClient.DisableAzureMonitorPreparer(context.Context, string, string) (*http.Request, error)
1. ExtensionsClient.DisableAzureMonitorResponder(*http.Response) (autorest.Response, error)
1. ExtensionsClient.DisableAzureMonitorSender(*http.Request) (ExtensionsDisableAzureMonitorFuture, error)
1. ExtensionsClient.EnableAzureMonitor(context.Context, string, string, AzureMonitorRequest) (ExtensionsEnableAzureMonitorFuture, error)
1. ExtensionsClient.EnableAzureMonitorPreparer(context.Context, string, string, AzureMonitorRequest) (*http.Request, error)
1. ExtensionsClient.EnableAzureMonitorResponder(*http.Response) (autorest.Response, error)
1. ExtensionsClient.EnableAzureMonitorSender(*http.Request) (ExtensionsEnableAzureMonitorFuture, error)
1. ExtensionsClient.GetAzureMonitorStatus(context.Context, string, string) (AzureMonitorResponse, error)
1. ExtensionsClient.GetAzureMonitorStatusPreparer(context.Context, string, string) (*http.Request, error)
1. ExtensionsClient.GetAzureMonitorStatusResponder(*http.Response) (AzureMonitorResponse, error)
1. ExtensionsClient.GetAzureMonitorStatusSender(*http.Request) (*http.Response, error)

### Struct Changes

#### New Structs

1. AzureMonitorRequest
1. AzureMonitorResponse
1. AzureMonitorSelectedConfigurations
1. AzureMonitorTableConfiguration
1. ExtensionsDisableAzureMonitorFuture
1. ExtensionsEnableAzureMonitorFuture
