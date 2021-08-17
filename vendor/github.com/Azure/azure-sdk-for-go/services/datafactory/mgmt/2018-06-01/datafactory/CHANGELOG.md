# Change History

## Additive Changes

### New Funcs

1. IntegrationRuntimesClient.ListOutboundNetworkDependenciesEndpoints(context.Context, string, string, string) (IntegrationRuntimeOutboundNetworkDependenciesEndpointsResponse, error)
1. IntegrationRuntimesClient.ListOutboundNetworkDependenciesEndpointsPreparer(context.Context, string, string, string) (*http.Request, error)
1. IntegrationRuntimesClient.ListOutboundNetworkDependenciesEndpointsResponder(*http.Response) (IntegrationRuntimeOutboundNetworkDependenciesEndpointsResponse, error)
1. IntegrationRuntimesClient.ListOutboundNetworkDependenciesEndpointsSender(*http.Request) (*http.Response, error)

### Struct Changes

#### New Structs

1. IntegrationRuntimeOutboundNetworkDependenciesCategoryEndpoint
1. IntegrationRuntimeOutboundNetworkDependenciesEndpoint
1. IntegrationRuntimeOutboundNetworkDependenciesEndpointDetails
1. IntegrationRuntimeOutboundNetworkDependenciesEndpointsResponse

#### New Struct Fields

1. CosmosDbMongoDbAPILinkedServiceTypeProperties.IsServerVersionAbove32
1. IntegrationRuntimeVNetProperties.SubnetID
