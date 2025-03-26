# Schema Design Considerations

Whilst it is acceptable in certain cases to map the schema of a new resource or feature when extending an existing resource one-to-one from the Azure API, in the majority of cases more consideration needs to be given how to expose the Azure API in Terraform so that the provider presents a consistent and intuitive experience to the end user.

Below are a list of common patterns found in the Azure API and how these typically get mapped within Terraform.

## Features that are toggled by the property `Enabled`

It is commonplace for features to be toggled on and off by an `Enabled` property within an object in the SDK used to interact with the Azure API. See the examples below.

Example A.
```go
type ManagedClusterStorageProfileBlobCSIDriver struct {
	Enabled *bool `json:"enabled,omitempty"`
}
```

Example B.
```go
type ManagedClusterWorkloadAutoScalerProfileVerticalPodAutoscaler struct {
	ControlledValues ControlledValues `json:"controlledValues"`
	Enabled          bool             `json:"enabled"`
	UpdateMode       UpdateMode       `json:"updateMode"`
}
```

This is handled in the provider one of two ways depending on if the `Enabled` field is by its self or with other fields in the object.

In the cases where `Enabled` is the only field within the object we opt to flatten the block into a single top level property (or higher level property if already nested inside a block). So in the case of Example A, this would become:

```go
"storage_blob_driver_enabled": {
    Type:     pluginsdk.TypeBool,
    Optional: true,
    Default:  false,
},
```

However, when there are multiple fields in addition to the `Enabled` field, and they are all required for the object/feature like in Example B, a block is created with all the fields including `Enabled`. The corresponding Terraform schema would be as follows:

```go
"vertical_pod_autoscaler": {
    Type:     pluginsdk.TypeList,
    Optional: true,
    MaxItems: 1,
    Elem: &pluginsdk.Resource{
        Schema: map[string]*pluginsdk.Schema{
            "enabled": {
                Type:     pluginsdk.TypeBool,
                Optional: true,
                Default:  false,
            },
            "update_mode": {
                Type:     pluginsdk.TypeString,
                Required: true,
                ValidateFunc: validation.StringInSlice([]string{
                    string(managedclusters.UpdateModeAuto),
                    string(managedclusters.UpdateModeInitial),
                    string(managedclusters.UpdateModeRecreate),
                }, false),
            },
            "controlled_values": {
                Type:     pluginsdk.TypeString,
                Required: true,
                ValidateFunc: validation.StringInSlice([]string{
                    string(managedclusters.ControlledValuesRequestsAndLimits),
                    string(managedclusters.ControlledValuesRequestsOnly),
                }, false),
            },
        },
    },
},
```

Finally, there are instances where the additional fields/properties for an object/feature are optional or few, as shown below.

Example C.
```go
type ManagedClusterStorageProfileDiskCSIDriver struct {
	Enabled *bool   `json:"enabled,omitempty"`
	Version *string `json:"version,omitempty"`
}
```

In cases like these one option is to flatten the block into two top level properties:

```go
"storage_disk_driver_enabled": {
    Type:     pluginsdk.TypeBool,
    Optional: true,
    Default:  false,
},

"storage_disk_driver_version": {
    Type:     pluginsdk.TypeString,
    Optional: true,
    Default:  "V1",
    ValidateFunc: validation.StringInSlice([]string{
        "V1",
        "V2",
    }, false),
},
```

A judgement call should be made based off the behaviour of the API and expectations of a user.

## The `None` value or similar

Many Azure APIs and services will accept the values like `None`, `Off`, or `Default` as a default value and expose it as a constant in the API specification. 

```
    "shutdownOnIdleMode": {
      "type": "string",
      "enum": [
        "None",
        "UserAbsence",
        "LowUsage"
      ],
```

Whilst it isn't uncommon to stumble across older resources in the provider that expose and accept these as a valid values, the provider is moving away from this pattern, since Terraform has its own null type i.e. by omitting the field. Existing `None`, `Off` or `Default` values within the provider are planned for removal in version 4.0.

This ultimately means that the end user doesn't need to bloat their configuration with superfluous information that is implied through the omission of information.

The resulting schema in Terraform would look as follows and also requires a conversion between the Terraform null value and `None` within the Create and Read functions.

```go
// How the property is exposed in the schema
"shutdown_on_idle": {
    Type:     pluginsdk.TypeString,
    Optional: true,
    ValidateFunc: validation.StringInSlice([]string{
        string(labplan.ShutdownOnIdleModeUserAbsence),
        string(labplan.ShutdownOnIdleModeLowUsage),
        // Note: Whilst the `None` value exists it's handled in the Create/Update and Read functions.
        // string(labplan.ShutdownOnIdleModeNone),
    }, false),
},

// Normalising in the create or expand function
func (r resource) Create() sdk.ResourceFunc {
	
	...
	
	var config resourceModel
	if err := metadata.Decode(&config); err != nil {
        return fmt.Errorf("decoding: %+v", err)
    }
	
	// The resource property shutdown_on_idle maps to the attribute shutdownOnIdle in the defined model for a typed resource in this example
	shutdownOnIdle := string(labplan.ShutdownOnIdleModeNone)
	if v := model.ShutdownOnIdle; v != "" {
		shutdownOnIdle = v
    }
	
	...
	
}

// Normalising in the read or flatten function
func (r resource) Read() sdk.ResourceFunc {
	
	...
	
	shutdownOnIdle := ""
	if v := props.ShutdownOnIdle; v != nil && v != string(labplan.ShutdownOnIdleModeNone) {
		shutdownOnIdle = string(*v)
    }
	
	state.ShutdownOnIdle = shutdownOnIdle
	
	...
	
}
```

## SKU fields

Because the Azure API implementation for SKU fields tends to vary we can't easily standardise on a single approach, however, we should try to stick to one of the following two implementations:

1. When the SKU can be set using a single argument (e.g. only the SKU name), use a top-level `sku` argument. 
2. When the SKU requires multiple arguments (e.g. `name` and `capacity`), use a `sku` block.

Example of a `sku` argument:
```go
"sku": {
	Type:     pluginsdk.TypeString,
	Optional: true,
	Default:  string(firewallpolicies.FirewallPolicySkuTierStandard),
	ForceNew: true,
	ValidateFunc: validation.StringInSlice([]string{
		string(firewallpolicies.FirewallPolicySkuTierPremium),
		string(firewallpolicies.FirewallPolicySkuTierStandard),
		string(firewallpolicies.FirewallPolicySkuTierBasic),
	}, false),
}
```

Example of a `sku` block:
```go
	"sku": {
		Type:     pluginsdk.TypeList,
		Required: true,
		MaxItems: 1,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"name": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ValidateFunc: validation.StringInSlice(helpers.PossibleValuesForSkuName(), false),
				},
				"capacity": {
					Type:     pluginsdk.TypeInt,
					Optional: true,
					Default:  1,
					ValidateFunc: validation.IntInSlice([]int{
						1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 20, 30, 40, 50, 60, 70, 80, 90, 100, 200,
					}),
				},
			},
		},
	},
```

While you may encounter arguments like `sku_name`, `sku_family`, and `capacity` in existing resources, new arguments should avoid this format and use one of the two options above.

## The `type` field

The Azure API makes use of classes and inheritance through discriminator types defined in the REST API specifications. A strong indicator that a resource is actually a discriminated type is through the definition of a `type` or `kind` property.

Rather than exposing a generic resource with all the possible fields for all the possible different `type`s, we intentionally opt to split these resources by the `type` to improve the user experience. This means we can only output the relevant fields for this `type` which in turn allows us to provide more granular validation etc.

Whilst there is a trade-off here, since this means that we have to maintain more Data Sources/Resources, this is a worthwhile trade-off since each of these resources only exposes the fields which are relevant for this resource, meaning the logic is far simpler than trying to maintain a generic resource and pushing the complexity onto end-users.

Taking the Data Factory Linked Service resources as an example which could have all of possible types defined below, each requiring a different set of inputs:

```go
"type": {
    Type:     pluginsdk.TypeString,
    Required: true,
    ValidateFunc: validation.StringInSlice([]string{
        string(datafactory.TypeBasicLinkedServiceTypeAzureBlobStorage),
        string(datafactory.TypeBasicLinkedServiceTypeAzureDatabricks),
        string(datafactory.TypeBasicLinkedServiceTypeAzureFileStorage),
        string(datafactory.TypeBasicLinkedServiceTypeAzureFunction),
        string(datafactory.TypeBasicLinkedServiceTypeAzureSearch),
		...
    }, false),
},
```

Would be better exposed as the following resources:

- `azurerm_data_factory_linked_service_azure_blob_storage`
- `azurerm_data_factory_linked_service_azure_databricks`
- `azurerm_data_factory_linked_service_azure_file_storage`
- `azurerm_data_factory_linked_service_azure_function`
- `azurerm_data_factory_linked_service_azure_search`

