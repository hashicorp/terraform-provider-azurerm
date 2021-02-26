Generated from https://github.com/Azure/azure-rest-api-specs/tree/a1eee0489c374782a934ec1f093abd16fa7718ca/specification/datafactory/resource-manager/readme.md tag: `package-2018-06`

Code generator @microsoft.azure/autorest.go@2.1.171


### New Constants

1. TypeBasicTrigger.TypeCustomEventsTrigger

### New Funcs

1. *CustomEventsTrigger.UnmarshalJSON([]byte) error
1. BlobEventsTrigger.AsCustomEventsTrigger() (*CustomEventsTrigger, bool)
1. BlobTrigger.AsCustomEventsTrigger() (*CustomEventsTrigger, bool)
1. ChainingTrigger.AsCustomEventsTrigger() (*CustomEventsTrigger, bool)
1. CustomEventsTrigger.AsBasicMultiplePipelineTrigger() (BasicMultiplePipelineTrigger, bool)
1. CustomEventsTrigger.AsBasicTrigger() (BasicTrigger, bool)
1. CustomEventsTrigger.AsBlobEventsTrigger() (*BlobEventsTrigger, bool)
1. CustomEventsTrigger.AsBlobTrigger() (*BlobTrigger, bool)
1. CustomEventsTrigger.AsChainingTrigger() (*ChainingTrigger, bool)
1. CustomEventsTrigger.AsCustomEventsTrigger() (*CustomEventsTrigger, bool)
1. CustomEventsTrigger.AsMultiplePipelineTrigger() (*MultiplePipelineTrigger, bool)
1. CustomEventsTrigger.AsRerunTumblingWindowTrigger() (*RerunTumblingWindowTrigger, bool)
1. CustomEventsTrigger.AsScheduleTrigger() (*ScheduleTrigger, bool)
1. CustomEventsTrigger.AsTrigger() (*Trigger, bool)
1. CustomEventsTrigger.AsTumblingWindowTrigger() (*TumblingWindowTrigger, bool)
1. CustomEventsTrigger.MarshalJSON() ([]byte, error)
1. MultiplePipelineTrigger.AsCustomEventsTrigger() (*CustomEventsTrigger, bool)
1. RerunTumblingWindowTrigger.AsCustomEventsTrigger() (*CustomEventsTrigger, bool)
1. ScheduleTrigger.AsCustomEventsTrigger() (*CustomEventsTrigger, bool)
1. Trigger.AsCustomEventsTrigger() (*CustomEventsTrigger, bool)
1. TumblingWindowTrigger.AsCustomEventsTrigger() (*CustomEventsTrigger, bool)

## Struct Changes

### New Structs

1. CustomEventsTrigger
1. CustomEventsTriggerTypeProperties
1. ManagedVirtualNetworkReference

### New Struct Fields

1. AzureDatabricksLinkedServiceTypeProperties.Authentication
1. AzureDatabricksLinkedServiceTypeProperties.PolicyID
1. AzureDatabricksLinkedServiceTypeProperties.WorkspaceResourceID
1. AzureMLExecutePipelineActivityTypeProperties.DataPathAssignments
1. AzureMLExecutePipelineActivityTypeProperties.MlPipelineEndpointID
1. AzureMLExecutePipelineActivityTypeProperties.Version
1. CustomActivityTypeProperties.AutoUserSpecification
1. IntegrationRuntimeSsisCatalogInfo.DualStandbyPairName
1. ManagedIntegrationRuntime.ManagedVirtualNetwork
1. WebActivityAuthentication.UserTenant
