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

## Preview Fields

Fields that are in preview should not be supported until they reach General Availability (GA) status, as they may change or be removed before becoming stable.

## Flattening nested properties

When designing schemas, consider flattening properties with `MaxItems: 1` that contain only a single nested property unless the service team has confirmed additional nested properties are imminent. In those cases, add an inline comment explaining why the block is left unflattened so reviewers understand the rationale.

**DO**
```go
"credential_certificate": {
    Type:     pluginsdk.TypeList,
    Optional: true,
    Elem:     &pluginsdk.Schema{
        Type:         pluginsdk.TypeString,
        // NOTE: validation is intentionally minimal since there is no stable RP contract for the certificate contents beyond a non-empty string.
        ValidateFunc: validation.StringIsNotEmpty,
    },
}
```

## Array fields with MinItems and MaxItems

If a field is an array, proper `MinItems` and `MaxItems` should be set based on the API constraints to provide clear validation feedback to users.

```go
"email_addresses": {
    Type:     pluginsdk.TypeList,
    Required: true,
    MinItems: 1,
    MaxItems: 20,
    Elem: &pluginsdk.Schema{
        Type:         pluginsdk.TypeString,
        ValidateFunc: validation.IsEmailAddress,
    },
},
```

## Required fields in Azure Portal vs API documentation

Fields marked as required in the Azure Portal (indicated by `*`) should be defined as `Required` in Terraform, unless the API accepts the request without them and still functions.

## Validation for TypeList fields with no Required fields

When a `pluginsdk.TypeList` block has no required nested fields, conditional validation such as `AtLeastOneOf` or `ExactlyOneOf` must be set on the optional fields to ensure the block is not empty and has at least one property configured.

```go
"setting": {
    Type:     pluginsdk.TypeList,
    Optional: true,
    MaxItems: 1,
    Elem: &pluginsdk.Resource{
        Schema: map[string]*pluginsdk.Schema{
            "linux": {
                Type:         pluginsdk.TypeList,
                Optional:     true,
                Elem:         osSchema(),
                AtLeastOneOf: []string{"setting.0.linux", "setting.0.windows"},
            },
            "windows": {
                Type:         pluginsdk.TypeList,
                Optional:     true,
                Elem:         osSchema(),
                AtLeastOneOf: []string{"setting.0.linux", "setting.0.windows"},
            },
        },
    },
}
```

## Validation

String arguments must be validated against the Resource Provider (RP) / API contract wherever possible.

As a rule of thumb, ensure common shapes have appropriate, specific validators:
- Validate `name`-like fields for length and allowed characters (use regex/patterns from the spec where available)
- Use `commonids` or service-specific ID validators for resource IDs
- Validate common formats like dates, IPs, ports, emails, and URIs

This means using constraints from the swagger/spec, SDK types/constants, or RP documentation:
- Allowed values (enums)
- Patterns/regex
- Length bounds
- Formats (dates, IPs, ports, emails, URIs)
- Resource ID shapes (prefer `commonids` or service-specific validators)

`validation.StringIsNotEmpty` is a LAST RESORT ONLY.
It is only acceptable when you have confirmed the RP truly accepts arbitrary free-form text and there are no stable rules to validate.
If you use `validation.StringIsNotEmpty`, you MUST add an inline comment explaining why stronger validation cannot be applied and what you checked.

Numeric arguments must specify a valid range wherever possible.
`validation.IntAtLeast(0)` is a LAST RESORT ONLY and must be justified the same way as above.

### Avoid overly-generic ValidateFunc (no "lazy validation")

When adding or modifying schema fields, do not default to minimal validators if the RP provides stronger constraints.

**DO** validate using the RP contract
```go
"allocation_strategy": {
    Type:     pluginsdk.TypeString,
    Optional: true,
    ValidateFunc: validation.StringInSlice([]string{
        string(compute.AllocationStrategyAutomatic),
        string(compute.AllocationStrategyPrioritized),
    }, false),
},

"name": {
    Type:     pluginsdk.TypeString,
    Required: true,
    ValidateFunc: validation.StringMatch(
        regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9.\-_]{0,79}$`),
        "...",
    ),
},
```

**DO** use ID/format validators for well-known shapes
```go
"subnet_id": {
    Type:         pluginsdk.TypeString,
    Required:     true,
    ValidateFunc: commonids.ValidateSubnetID,
},
```

**DO NOT** use minimal validation when stronger constraints exist
```go
// BAD: RP defines allowed values / pattern / bounds, but this accepts nearly anything.
ValidateFunc: validation.StringIsNotEmpty,

// BAD: RP bounds exist, but this only enforces non-negative.
ValidateFunc: validation.IntAtLeast(0),
```

**LAST RESORT** (only when RP truly has no constraints)
```go
"description": {
    Type:     pluginsdk.TypeString,
    Optional: true,
    // NOTE: validation is intentionally minimal because the RP accepts arbitrary free-form text for this field (no enum/pattern/length constraints found).
    ValidateFunc: validation.StringIsNotEmpty,
},
```

```go
"name": {
	Type:     pluginsdk.TypeString,
	Required: true,
	ForceNew: true,
	ValidateFunc: validation.StringMatch(
		regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9.\-_]{0,79}$`),
		"`name` must be between 1 and 80 characters. It must start with an alphanumeric character and can contain alphanumeric characters, dots (.), hyphens (-), and underscores (_).",
	),
},

"subnet_id": {
	Type:         pluginsdk.TypeString,
	Required:     true,
	ValidateFunc: commonids.ValidateSubnetID,
},

"description": {
	Type:         pluginsdk.TypeString,
	Optional:     true,
	// NOTE: validation is intentionally minimal because the RP accepts arbitrary free-form text for this field (no enum/pattern/length constraints found).
	ValidateFunc: validation.StringIsNotEmpty,
},

"extensions_time_budget": {
	Type:         pluginsdk.TypeString,
	Optional:     true,
	Default:      "PT1H30M",
	ValidateFunc: validate.ISO8601DurationBetween("PT15M", "PT2H"),
},

"filter_value_percentage": {
	Type:         pluginsdk.TypeFloat,
	Optional:     true,
	ValidateFunc: validation.FloatBetween(0, 100),
},

"ip_address": {
	Type:         pluginsdk.TypeString,
	Optional:     true,
	ValidateFunc: azValidate.IPv4Address,
},

"output_blob_uri": {
	Type:         pluginsdk.TypeString,
	Optional:     true,
	ValidateFunc: validation.IsURLWithHTTPS,
},

"sim_policy_id": {
	Type:         pluginsdk.TypeString,
	Optional:     true,
	ValidateFunc: simpolicy.ValidateSimPolicyID,
},

"storage_size_in_gb": {
	Type:         pluginsdk.TypeInt,
	Optional:     true,
	ValidateFunc: validation.IntBetween(32, 16384),
},
```
