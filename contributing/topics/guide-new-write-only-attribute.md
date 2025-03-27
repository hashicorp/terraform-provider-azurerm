# Guide: Adding a new Write-Only Attribute 

This guide covers how to add a new Write-Only (WO) Attribute to a resource. A WO Attribute can accept ephemeral values and is never persisted in state. 

> **Note:** Write-Only Attributes are only available in Terraform version 1.11 or higher.

Good candidates for WO Attributes are sensitive user supplied properties, e.g. passwords, certificates, and keys, can be added in addition to an existing sensitive property.

There are however limitations on what can be added as a WO Attribute, the original sensitive property:
* Cannot be `ForceNew`
* Cannot be `Computed`
* Cannot be within a set of nested blocks or set or nested attributes
* Cannot be a block (list or set) or a map
* Cannot be used in data sources or the provider schemas

Adding a new WO Attribute consists of the following steps:

1. [Updating the Resource Schema](#updating-the-resource-schema)
2. [Updating the Create, Read and Update functions](#updating-the-create-read-and-update-functions)
3. [Adding Validation to prefer the WO Attribute over the non-WO Sensitive Attribute](#adding-validation)
4. [Adding the Tests](#adding-tests)
5. [Updating the Documentation](#updating-the-documentation)

In the steps outlined above we're going to look at a fictional resource called `azurerm_some_database` that has an existing sensitive property called `password` for which we're going to add a WO attribute.

## Updating the Resource Schema

A new WO attribute must be accompanied by the addition of a regular attribute whose presence and value is used to determine when the value of a WO Attribute has changed and signals the provider to send the value of the WO attribute.

As a result we add two new properties to the schema, `password_wo` and `password_wo_version`.

```go
... // omitted for brevity

"password": {
	Type:          pluginsdk.TypeString,
	Optional:      true,
	Sensitive:     true,
	ConflictsWith: []string{"password_wo"} // this must be set to prevent both the sensitive `password` and the wo attribute `password_wo` from being set
},

"password_wo": {
	Type:          pluginsdk.TypeString, 
	Optional:      true,
	WriteOnly:     true, 
	RequiredWith:  []string{"password_wo_version"} // this must be set to ensure the "trigger" property is provided with the wo attribute 
	ConflictsWith: []string{"password_wo_version"} // this must be set to prevent both the sensitive `password` and the wo attribute `password_wo` from being set
}

"password_wo_version": {
	Type:         pluginsdk.TypeInt,
	Optional:     true,
	RequiredWith: []string{"password_wo"} // this must be set to ensure the "trigger" property is provided with the wo attribute
}

... // omitted for brevity
```

## Updating the Create, Read and Update functions

In the `Create` function we make use of the helper function `pluginsdk.GetWriteOnly` to retrieve the WO attribute.

```go
func (SomeDatabase) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			... // omitted for brevity
			
			// use the GetWriteOnly helper to retrieve the WO attribute
			woPassword, err := pluginsdk.GetWriteOnly(metadata.ResourceData, "password_wo", cty.String)
			if err != nil {
				return err
			}
			
			// set it in the payload if the WO attribute is not null
			if !woPassword.IsNull() {
			    payload.Properties.Password = woPassword.AsString()
			}
			
			... // omitted for brevity
		}
	}
}
```

The only update in the `Read` function is to retrieve the value for the WO attribute's trigger property `password_wo_version` and to set that into state.

```go
func (SomeDatabase) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			... // omitted for brevity

			// since WO attributes are not persisted in state we do not need to write it back to state
			// but we do need to retrieve the value for the trigger attribute from the config and set that into
			// state to prevent a perma diff
			state.PasswordWOVersion = metadata.ResourceData.Get("password_wo_version").(int)
			
			... // omitted for brevity
		}
	}
}
```

In the `Update` function we rely on changes to the trigger attribute `password_wo_version` to know when to update the WO attribute.

```go
func (SomeDatabase) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			... // omitted for brevity

			// check if the trigger attribute has any changes
			if metadata.ResourceData.HasChange("password_wo_version") {
				woPassword, err := pluginsdk.GetWriteOnly(metadata.ResourceData, "password_wo", cty.String)
				if err != nil {
					return err
				}
				
				// set it in the payload if the WO attribute is not null
				if !woPassword.IsNull() {
					payload.Properties.Password = woPassword.AsString()
				}
			}
			
			... // omitted for brevity
		}
	}
}
```

## Adding Validation

The `terraform-plugin-sdk@v2` provides a helpful validation for WO attributes that surfaces a warning to users if they are on a version of Terraform that supports WO attributes but are using the non-WO attribute version of a sensitive property.

> **Note:** We currently recommend not adding this validation to the resource since the only way to remove the warning diagnostic is to use to WO attribute.

```go
// update the interface that the resource should implement
var _ sdk.ResourceWithConfigValidation = SomeDatabase{}

// add the config validation method to satisfy the new interface
func (r SomeDatabase) ValidateRawResourceConfig() []schema.ValidateRawResourceConfigFunc {
    return []schema.ValidateRawResourceConfigFunc{
		pluginSdkValidation.PreferWriteOnlyAttribute(cty.GetAttrPath("password"), cty.GetAttrPath("password_wo")),
	}
}

... // omitted for brevity
```

## Adding Tests

To cover our bases we should test the following paths for a WO attribute:
* Creating a resource with the WO attribute
* Updating a resource with the WO attribute
* Updating a resource that uses the original sensitive property to the WO attribute
* Updating a resource that uses the WO attribute back to the original sensitive property

These paths can be tested by the addition of two test cases:

```go
func TestAccSomeDatabase_writeOnlyPassword(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_some_database", "test")
	r := SomeDatabaseResource{}

	resource.ParallelTest(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(version.Must(version.NewVersion("1.11.0"))),
		},
		ProtoV5ProviderFactories: framework.ProtoV5ProviderFactoriesInit(context.Background(), "azurerm"),
		Steps: []resource.TestStep{
			{
				Config: r.writeOnlyPassword(data, "a-secret-from-kv", 1),
				Check:  check.That(data.ResourceName).ExistsInAzure(r),
			},
			data.ImportStep("password_wo_version"),
			{
				Config: r.writeOnlyPassword(data, "a-secret-from-kv-updated", 2),
				Check:  check.That(data.ResourceName).ExistsInAzure(r),
			},
			data.ImportStep("password_wo_version"),
		},
	})
}
```

```go
func TestAccSomeDatabase_updateToWriteOnlyPassword(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_some_database", "test")
	r := SomeDatabaseResource{}

	resource.ParallelTest(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(version.Must(version.NewVersion("1.11.0"))),
		},
		ProtoV5ProviderFactories: framework.ProtoV5ProviderFactoriesInit(context.Background(), "azurerm"),
		Steps: []resource.TestStep{
			{
				Config: r.basic(data),
				Check:  check.That(data.ResourceName).ExistsInAzure(r),
			},
			data.ImportStep("password"),
			{
				Config: r.writeOnlyPassword(data, "a-secret-from-kv", 1),
				Check:  check.That(data.ResourceName).ExistsInAzure(r),
			},
			data.ImportStep("password", "password_wo_version"),
			{
				Config: r.basic(data),
				Check:  check.That(data.ResourceName).ExistsInAzure(r),
			},
			data.ImportStep("password"),
		},
	})
}
```

To reduce the amount of unnecessary test templating, we should make use of the `acceptance.WriteOnlyKeyVaultSecretTemplate` test config template which provisions all the necessary dependencies to reference a secret value using the `azurerm_key_vault_secret` ephemeral resource.

```go
func (r SomeDatabaseResource) writeOnlyPassword(data acceptance.TestData, secret string, version int) string {
	return fmt.Sprintf(`
%s

%s

resource "azurerm_some_database" "test" {
  name                = "acctest-db-%[3]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  login               = "some_admin_login"
  password_wo         = ephemeral.azurerm_key_vault_secret.test.value
  password_wo_version = %[4]d
}
`, r.template(data), acceptance.WriteOnlyKeyVaultSecretTemplate(data, secret), data.RandomInteger, version)
}
```

## Updating the Documentation

When documenting WO attributes we specify `Write-Only` in the parentheses that contains the `Required` and `Optional` information.

```markdown
...

* `password` - (Optional) The Password associated with the `login` for the Database.

* `password_wo` - (Optional, Write-Only) The Password associated with the `login` for the Database.

* `password_wo_version` - (Optional) An integer value used to trigger an update for `password_wo`. This property should be incremented when updating `password_wo`.

...

```