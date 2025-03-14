# Guide: Adding a new Feature to the Features Block

This guide covers how to add a new Feature to the Features Block ([Terraform Docs](https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/guides/features-block)) that will change the default behaviour for how a resource or service works. Reasons for this can include:

* Purging a resource during delete

* Recovering a resource that has been soft deleted during create

* Detach a connected resource during deletion

Following are the steps needed to add a new Feature to the Feature Block:

> **Note:** The Azure Provider is in the process of moving towards a new Framework Plugin for the provider. Because of this, we must update the provider in a few areas when adding a new feature. We'll update the following areas `internal/features`, `internal/provider`, `internal/provider/framework`, and the resource file itself.

### Updating `internal/features`

1. Update `internal/features/user_flags.go` with either a new block for the service package or updating an existing service package with the new feature to add. Added struct names should represent the service package they affect, and feature names should concisely describe their effect.

```go
type UserFeatures struct {
    KeyVault KeyVaultFeatures
}

type KeyVaultFeatures struct {
    PurgeSoftDeleteOnDestroy bool
}
```

2. Update `internal/features/defaults.go` with what the default value for the new feature will be. This must represent the current default behaviour of the resource(s) to avoid this becoming a breaking change when the feature flagged behaviour is added to the target resource(s).

```go
func Default() UserFeatures {
    return UserFeatures{
    ...
	KeyVault: KeyVaultFeatures{
        PurgeSoftDeleteOnDestroy: true,
    }
    ...
}
```
### Updating `internal/provider`

1. Update `internal/provider/feature.go` with what the Terraform schema will look like and how to thread it into the features block

```go
func schemaFeatures(supportLegacyTestSuite bool) *pluginsdk.Schema {
    featuresMap := map[string]*pluginsdk.Schema{
        ...
        "key_vault": {
            Type:     pluginsdk.TypeList,
            Optional: true,
            MaxItems: 1,
            Elem: &pluginsdk.Resource{
            Schema: map[string]*pluginsdk.Schema{
                "purge_soft_delete_on_destroy": {
                    Description: "When enabled soft-deleted `azurerm_key_vault` resources will be permanently deleted (e.g purged), when destroyed",
                    Type:        pluginsdk.TypeBool,
                    Optional:    true,
                    Default:     true,
                },
            },
        },
        ...
    }
}

func expandFeatures(input []interface{}) features.UserFeatures {
    ...
    if raw, ok := val["key_vault"]; ok {
        items := raw.([]interface{})
        if len(items) > 0 && items[0] != nil {
            keyVaultRaw := items[0].(map[string]interface{})
            if v, ok := keyVaultRaw["purge_soft_delete_on_destroy"]; ok {
                featuresMap.KeyVault.PurgeSoftDeleteOnDestroy = v.(bool)
            }
        }
    }
    ...
}
```

2. Update `internal/provider/feature_test.go` to include a test for every permutation of the feature you are adding to the TestExpandFeatures test and a test dedicated to the service package of the feature.

```go
func TestExpandFeatures(t *testing.T) {
    testData := []struct {
        Name     string
        Input    []interface{}
        EnvVars  map[string]interface{}
        Expected features.UserFeatures
    }{
        {
            Name:  "Empty Block",
            Input: []interface{}{},
            Expected: features.UserFeatures{
                ...
                KeyVault: features.KeyVaultFeatures{
                    PurgeSoftDeleteOnDestroy:         true,
                },
                ...
            }
        },
        {
            Name: "Complete Enabled",
            Input: []interface{}{
                map[string]interface{}{
                    ...
                    "key_vault": []interface{}{
                        map[string]interface{}{
    	                    "purge_soft_delete_on_destroy": true,
    	                },
    	            },   
     	            ...
    	        },
            },
            Expected: features.UserFeatures{
                ...
                KeyVault: features.KeyVaultFeatures{
                    PurgeSoftDeleteOnDestroy: true,
                },
                ...
            },
        },
        {
            Name: "Complete Disabled",
            Input: []interface{}{
                map[string]interface{}{
                    ...
                    "key_vault": []interface{}{
                        map[string]interface{}{
                            "purge_soft_delete_on_destroy": false,
                        },
                    },
                    ...
                },
            },
            Expected: features.UserFeatures{
                ...
                KeyVault: features.KeyVaultFeatures{
                PurgeSoftDeleteOnDestroy: false,
                },
                ...
            },
        },
    },	
}


func TestExpandFeaturesKeyVault(t *testing.T) {
    testData := []struct {
        Name     string
        Input    []interface{}
        EnvVars  map[string]interface{}
        Expected features.UserFeatures
    }{
        {
            Name: "Empty Block",
            Input: []interface{}{
                map[string]interface{}{
                    "key_vault": []interface{}{},
                },
            },
            Expected: features.UserFeatures{
                KeyVault: features.KeyVaultFeatures{
                    PurgeSoftDeleteOnDestroy: true,
                },
            },
        },
        {
            Name: "Purge Soft Delete On Destroy",
            Input: []interface{}{
                map[string]interface{}{
                    "key_vault": []interface{}{
                        map[string]interface{}{
                            "purge_soft_delete_on_destroy": true,
                        },
                    },
                },
            },
            Expected: features.UserFeatures{
                KeyVault: features.KeyVaultFeatures{
                    PurgeSoftDeleteOnDestroy: true,
                },
            },
        },
        {
            Name: "Purge Soft Delete On Destroy Disabled",
            Input: []interface{}{
                map[string]interface{}{
                    "key_vault": []interface{}{
                        map[string]interface{}{
                            "purge_soft_delete_on_destroy": false,
                        },
                    },
                },
            },
            Expected: features.UserFeatures{
                KeyVault: features.KeyVaultFeatures{
                    PurgeSoftDeleteOnDestroy: false,
                },
            },
        },
    }

    for _, testCase := range testData {
        t.Logf("[DEBUG] Test Case: %q", testCase.Name)
        result := expandFeatures(testCase.Input)
        if !reflect.DeepEqual(result.KeyVault, testCase.Expected.KeyVault) {
            t.Fatalf("Expected %+v but got %+v", result.KeyVault, testCase.Expected.KeyVault)
        }
    }
}
```
### Updating `internal/provider/framework`

1. Update `internal/provider/framework/model.go`

```go
// For new services, add a List type for the new block with a `tfsdk` struct tag that matches the schema name for the block, for new features in an existing block/service, this can be skipped.
type Features struct {
    ...
    KeyVault types.List `tfsdk:"key_vault"`
    ...
}

// and an attribute map variable for the block, or add to the appropriate existing var
var FeaturesAttributes = map[string]attr.Type{
    ...
    "key_vault": types.ListType{}.WithElementType(types.ObjectType{}.WithAttributeTypes(KeyVaultAttributes)),
    ...
}

// Add a Go struct that matches the new block or add to the appropriate existing struct
type KeyVault struct {
    PurgeSoftDeleteOnDestroy types.Bool `tfsdk:"purge_soft_delete_on_destroy"`
}

// finally, create the attribute map variable for the new block, or add the feature to the appropriate existing map
var KeyVaultAttributes = map[string]attr.Type{
    "purge_soft_delete_on_destroy": types.BoolType
}
```

2. Update `internal/provider/framework/provider.go`

```go
func (p *azureRmFrameworkProvider) Schema(_ context.Context, _ provider.SchemaRequest, response *provider.SchemaResponse) {
    response.Schema = schema.Schema{
        ...
        Blocks: map[string]schema.Block{
            "features": schema.ListNestedBlock{
                Validators: []validator.List{
                    listvalidator.SizeBetween(1, 1),
                },
                NestedObject: schema.NestedBlockObject{
                    Blocks: map[string]schema.Block{
                        ...
                        // Add an attribute map variable for the new block or add to the existing map inside the Nested Object
                        "key_vault": schema.ListNestedBlock{
                            NestedObject: schema.NestedBlockObject{
                            	Attributes: map[string]schema.Attribute{
                                	"purge_soft_delete_on_destroy": schema.BoolAttribute{
                            	    	Description: "When enabled soft-deleted `azurerm_key_vault` resources will be permanently deleted (e.g purged), when destroyed",
                            	    	Optional:    true,
                                   },
                                },
                            },
                        },
                        ...
                    },
                },
            },
        },	
        ...
    }
}
```

3. Update `internal/provider/framework/config.go`

```go
// Add a new check that the feature has been specified in the config that then loads the feature into the provider or add the new feature to the existing block.
func (p *ProviderConfig) Load(ctx context.Context, data *ProviderModel, tfVersion string, diags *diag.Diagnostics) {
    ...
    if !features.KeyVault.IsNull() && !features.KeyVault.IsUnknown() {
        var feature []KeyVault
        d := features.KeyVault.ElementsAs(ctx, &feature, true)
        diags.Append(d...)
        if diags.HasError() {
            return
        }

        f.KeyVault.PurgeSoftDeleteOnDestroy = true
        if !feature[0].PurgeSoftDeleteOnDestroy.IsNull() && !feature[0].PurgeSoftDeleteOnDestroy.IsUnknown() {
            f.KeyVault.PurgeSoftDeleteOnDestroy = feature[0].PurgeSoftDeleteOnDestroy.ValueBool()
        }
    }
    ...
}
```

4. Update  `internal/provider/framework/config_test.go` with the Features Model and Attributes

```go
func defaultFeaturesList() types.List {
    ...
    // Add a NewObjectValueFrom that holds what type of feature you have or append to the existing ObjectValueFrom
    keyVault, _ := basetypes.NewObjectValueFrom(context.Background(), KeyVaultAttributes, map[string]attr.Value{
        "purge_soft_delete_on_destroy":                            basetypes.NewBoolNull(),
    })
    keyVaultList, _ := basetypes.NewListValue(types.ObjectType{}.WithAttributeTypes(KeyVaultAttributes), []attr.Value{keyVault})
    ...
    // If the added feature is supporting a new service, add it to the following list of services
    fData, d := basetypes.NewObjectValue(FeaturesAttributes, map[string]attr.Value{
        ...
        "key_vault": keyVaultList,
        ...
    }
}
```
### Update the resource

1. Update `internal/service/serviceName/resourceName.go` in this case `internal/service/keyvault/key_vault_resource.go` to include the functionality of the added feature.

```go
func resourceKeyVaultDelete(d *pluginsdk.ResourceData, meta interface{}) error {
    ...
    if meta.(*clients.Client).Features.KeyVault.PurgeSoftDeleteOnDestroy {
        // Purge the Keyvault
    }
    ...
}
```

2. Update `internal/service/serviceName/resourceName_test.go` in this case `internal/service/keyvault/key_vault_resource_test.go` to test the new feature.

```go
func TestAccKeyVault_softDeleteRecoveryDisabled(t *testing.T) {
    data := acceptance.BuildTestData(t, "azurerm_key_vault", "test")
    r := KeyVaultResource{}

    data.ResourceTest(t, r, []acceptance.TestStep{
        {
        	// create it regularly
        	Config: r.softDeleteRecoveryDisabled(data),
        	Check: acceptance.ComposeTestCheckFunc(
             	check.That(data.ResourceName).ExistsInAzure(r),
                check.That(data.ResourceName).Key("purge_protection_enabled").HasValue("false"),
        	),
        },
        data.ImportStep(),
        {
            // delete the key vault
            Config: r.softDeleteAbsent(data),
        },
        {
            // attempting to re-create it requires recovery, which is enabled by default
            Config:      r.softDeleteRecoveryDisabled(data),
            ExpectError: regexp.MustCompile("An existing soft-deleted Key Vault exists with the Name"),
        },
    })
}

func (KeyVaultResource) softDeleteRecoveryDisabled(data acceptance.TestData) string {
    return fmt.Sprintf(`
provider "azurerm" {
  features {
    key_vault {
      recover_soft_deleted_key_vaults = false
    }
  }
}
...
`)
}
```

At this point, if all tests have passed including the tests found in `internal/provider/function/normalise_resource_id_test.go` and `internal/provider/function/parse_resource_id_test.go`, the Feature should be implemented and ready for use. 