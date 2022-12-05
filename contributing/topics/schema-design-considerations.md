# Schema Design Considerations

Whilst it is acceptable in certain cases to map the schema of a new resource, or feature when extending an existing resource, one-to-one from the Azure API, in the majority of cases more consideration needs to be given on how to expose the Azure API in Terraform so that the provider presents a consistent and intuitive experience to the end user.

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

Although there are still examples within the provider where this property is exposed in the Terraform schema, the provider is moving away from this and instead translates this behaviour in one of two ways.

In cases where `Enabled` is the only field required to turn a feature on and off we opt to flatten the block into a single top level property (or higher level property if already nested inside a block). So in the case of Example A, this would become:

```go
"storage_blob_driver_enabled": {
    Type:     pluginsdk.TypeBool,
    Optional: true,
    Default:  false,
},
```

For features that can accept or require configuration, i.e. the object contains additional properties other than `Enabled` like in Example B, the behaviour should be such that when the block is present the feature is enabled, and when it is absent it is disabled.

```go
"vertical_pod_autoscaler": {
    Type:     pluginsdk.TypeList,
    Optional: true,
    MaxItems: 1,
    Elem: &pluginsdk.Resource{
        Schema: map[string]*pluginsdk.Schema{
            "update_mode": {
                Type:     pluginsdk.TypeString,
                Required: true,
                ValidateFunc: validation.StringInSlice([]string{
                    string(managedclusters.UpdateModeAuto),
                    string(managedclusters.UpdateModeInitial),
                    string(managedclusters.UpdateModeOff),
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

There are instances where configuration properties for a feature is optional, as shown below.

Example C.
```go
type ManagedClusterStorageProfileDiskCSIDriver struct {
	Enabled *bool   `json:"enabled,omitempty"`
	Version *string `json:"version,omitempty"`
}
```

In cases like these we can choose to either flatten the block into two top level properties.

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

Depending on the behaviour of the Azure API and the default set by it, it's often worthwhile turning the optional properties into required ones in the Terraform schema, to avoid having to set empty blocks to enable features.

```go
"storage_disk_driver": {
    Type:     pluginsdk.TypeList,
    Optional: true,
    MaxItems: 1,
    Elem: &pluginsdk.Resource{
        Schema: map[string]*pluginsdk.Schema{
            "version": {
                Type:     pluginsdk.TypeString,
                Required: true,
                ValidateFunc: validation.StringInSlice([]string{
                    "V1",
                    "V2",
                }, false),
            },
        },
    },
},
```

## The `None` value

Many Azure APIs and services will accept the value `None` and expose it as an enum in the API specification. 

```
    "shutdownOnIdleMode": {
      "type": "string",
      "enum": [
        "None",
        "UserAbsence",
        "LowUsage"
      ],
```

Whilst it isn't uncommon to stumble across older resources in the provider that accept this as a valid value, the provider is moving away from the pattern of exposing `None`, since Terraform has its own null type i.e. by omitting the field this.
This ultimately means that the end user doesn't need to bloat their configuration with superfluous information.

The resulting schema in Terraform would look as follows and also requires a conversion between the Terraform null value and `None` within the Create and Read functions.

```go
"shutdown_on_idle": {
    Type:     pluginsdk.TypeString,
    Optional: true,
    ValidateFunc: validation.StringInSlice([]string{
        string(labplan.ShutdownOnIdleModeUserAbsence),
        string(labplan.ShutdownOnIdleModeLowUsage),
    }, false),
},
```





