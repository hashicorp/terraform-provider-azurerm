# Change History

## Additive Changes

### New Constants

1. StagingEnvironmentPolicy.StagingEnvironmentPolicyDisabled
1. StagingEnvironmentPolicy.StagingEnvironmentPolicyEnabled

### New Funcs

1. *RemotePrivateEndpointConnection.UnmarshalJSON([]byte) error
1. AppsClient.CreateOrUpdateSwiftVirtualNetworkConnectionWithCheck(context.Context, string, string, SwiftVirtualNetwork) (SwiftVirtualNetwork, error)
1. AppsClient.CreateOrUpdateSwiftVirtualNetworkConnectionWithCheckPreparer(context.Context, string, string, SwiftVirtualNetwork) (*http.Request, error)
1. AppsClient.CreateOrUpdateSwiftVirtualNetworkConnectionWithCheckResponder(*http.Response) (SwiftVirtualNetwork, error)
1. AppsClient.CreateOrUpdateSwiftVirtualNetworkConnectionWithCheckSender(*http.Request) (*http.Response, error)
1. PossibleStagingEnvironmentPolicyValues() []StagingEnvironmentPolicy
1. RemotePrivateEndpointConnection.MarshalJSON() ([]byte, error)
1. RemotePrivateEndpointConnectionProperties.MarshalJSON() ([]byte, error)
1. ResponseMessageEnvelopeRemotePrivateEndpointConnection.MarshalJSON() ([]byte, error)

### Struct Changes

#### New Structs

1. ArmPlan
1. RemotePrivateEndpointConnection
1. RemotePrivateEndpointConnectionProperties
1. ResponseMessageEnvelopeRemotePrivateEndpointConnection

#### New Struct Fields

1. SiteConfig.AcrUseManagedIdentityCreds
1. SiteConfig.AcrUserManagedIdentityID
1. SiteConfig.PublicNetworkAccess
1. SitePatchResourceProperties.VirtualNetworkSubnetID
1. SiteProperties.VirtualNetworkSubnetID
1. StaticSite.AllowConfigFileUpdates
1. StaticSite.PrivateEndpointConnections
1. StaticSite.StagingEnvironmentPolicy
