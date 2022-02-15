# Change History

## Additive Changes

### New Constants

1. NetworkFeatures.NetworkFeaturesBasic
1. NetworkFeatures.NetworkFeaturesStandard
1. VolumeStorageToNetworkProximity.VolumeStorageToNetworkProximityDefault
1. VolumeStorageToNetworkProximity.VolumeStorageToNetworkProximityT1
1. VolumeStorageToNetworkProximity.VolumeStorageToNetworkProximityT2

### New Funcs

1. *SubscriptionQuotaItem.UnmarshalJSON([]byte) error
1. AzureEntityResource.MarshalJSON() ([]byte, error)
1. NewResourceQuotaLimitsClient(string) ResourceQuotaLimitsClient
1. NewResourceQuotaLimitsClientWithBaseURI(string, string) ResourceQuotaLimitsClient
1. PossibleNetworkFeaturesValues() []NetworkFeatures
1. PossibleVolumeStorageToNetworkProximityValues() []VolumeStorageToNetworkProximity
1. ProxyResource.MarshalJSON() ([]byte, error)
1. Resource.MarshalJSON() ([]byte, error)
1. ResourceQuotaLimitsClient.Get(context.Context, string, string) (SubscriptionQuotaItem, error)
1. ResourceQuotaLimitsClient.GetPreparer(context.Context, string, string) (*http.Request, error)
1. ResourceQuotaLimitsClient.GetResponder(*http.Response) (SubscriptionQuotaItem, error)
1. ResourceQuotaLimitsClient.GetSender(*http.Request) (*http.Response, error)
1. ResourceQuotaLimitsClient.List(context.Context, string) (SubscriptionQuotaItemList, error)
1. ResourceQuotaLimitsClient.ListPreparer(context.Context, string) (*http.Request, error)
1. ResourceQuotaLimitsClient.ListResponder(*http.Response) (SubscriptionQuotaItemList, error)
1. ResourceQuotaLimitsClient.ListSender(*http.Request) (*http.Response, error)
1. SubscriptionQuotaItem.MarshalJSON() ([]byte, error)
1. SubscriptionQuotaItemProperties.MarshalJSON() ([]byte, error)
1. TrackedResource.MarshalJSON() ([]byte, error)

### Struct Changes

#### New Structs

1. AzureEntityResource
1. LogSpecification
1. ProxyResource
1. Resource
1. ResourceQuotaLimitsClient
1. SubscriptionQuotaItem
1. SubscriptionQuotaItemList
1. SubscriptionQuotaItemProperties
1. TrackedResource

#### New Struct Fields

1. MetricSpecification.EnableRegionalMdmAccount
1. MetricSpecification.IsInternal
1. ServiceSpecification.LogSpecifications
1. VolumeProperties.NetworkFeatures
1. VolumeProperties.NetworkSiblingSetID
1. VolumeProperties.StorageToNetworkProximity
