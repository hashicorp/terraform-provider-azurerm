# Guide: Resource Identity

This guide covers adding Resource Identity to a new or existing resource. For more information on Resource Identity, see [Resources - Identity](https://developer.hashicorp.com/terraform/plugin/sdkv2/resources/identity).

> The provider's Resource Identity generator does not yet support all identity types. `commonids.CompositeResourceID` and any custom resource IDs (i.e. not one provided by `commonids` or `go-azure-sdk/resource-manager`) are not supported.

## Adding Resource Identity

### Typed Resources

To add Resource Identity to a typed resource, we will need to implement the `sdk.ResourceWithIdentity` interface and modify the `Read()` function.

1. Define a variable of type `sdk.ResourceWithIdentity` and assign it a value of the resource type struct.

    ```go
    package example

    import "github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
   
    type ExampleResource struct{}

    var _ sdk.ResourceWithIdentity = ExampleResource{}
    ```
   
2. Add the `Identity()` method, this method should return a pointer to the correct resource ID, if you are unsure, you can reference the `IDValidationFunc` method, the ID that is being validated here is the one you'll want to use.

    ```go
    package example
    
    import "github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
    import "github.com/hashicorp/go-azure-helpers/resourceids"
    
    type ExampleResource struct{}
    
    var _ sdk.ResourceWithIdentity = ExampleResource{}
    
    func (r ExampleResource) Identity() resourceids.ResourceId {
        return &examplepackage.ExampleResourceId{}
    }
    ```

3. Update the `Read()` function to include a step setting the Resource Identity data into state. Resource Identity data does not have to be set manually, we can make use of the `pluginsdk.SetResourceIdentityData` helper function.

    ```go
    func (r ExampleResource) Read() sdk.ResourceFunc {
        return sdk.ResourceFunc{
            Timeout: 5 * time.Minute,
            Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
                client := metadata.Client.Service.ExampleClient
                id, err := examplepackage.ParseExampleResourceID(metadata.ResourceData.Id())
                if err != nil {
                    return err
                }
                 
                ...
                 
                if err := pluginsdk.SetResourceIdentityData(metadata.ResourceData, id); err != nil {
                    return err
                }
                
                return metadata.Encode(&model)
            },
        }
    }
    ```

4. Add an acceptance test to ensure the identity data is accurately set into state, please reference [Resource Identity Tests](#resource-identity-tests).

### Untyped Resources

To add Resource Identity to an untyped resource, follow the steps below.

1. Add the `Identity` schema. Here, we make use of the `pluginsdk.GenerateIdentitySchema` function, which takes in a pointer to a `resourceids.ResourceId`. The ID provided here should be the same as the ID that is being parsed in the `Importer` field.

    ```go
    package example
    
    import (
        "github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
    )
    
    func resourceExample() *pluginsdk.Resource {
        return &pluginsdk.Resource{
            Create: resourceExampleCreate,
            Read: resourceExampleRead,
            Update: resourceExampleUpdate,
            Delete: resourceExampleDelete,
    
            Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
                _, err := examplepackage.ParseExampleID(id)
                return err
            }),
            
            // We will be including the new `Identity` field
            Identity: &schema.ResourceIdentity{
                SchemaFunc: pluginsdk.GenerateIdentitySchema(&examplepackage.ExampleId{}),
            },
            
            ...
        }
    }
    ```
   
2. Update the `Importer` field, we'll want to use the `pluginsdk.ImporterValidatingIdentity` function and provide it with the same resource ID as the `pluginsdk.GenerateIdentitySchema` function.

    ```go
        package example
        
        import (
            "github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
        )
        
        func resourceExample() *pluginsdk.Resource {
            return &pluginsdk.Resource{
                Create: resourceExampleCreate,
                Read: resourceExampleRead,
                Update: resourceExampleUpdate,
                Delete: resourceExampleDelete,
        
                Importer: pluginsdk.ImporterValidatingIdentity(&examplepackage.ExampleId{}),
                
                // We will be including the new `Identity` field
                Identity: &schema.ResourceIdentity{
                    SchemaFunc: pluginsdk.GenerateIdentitySchema(&examplepackage.ExampleId{}),
                },
                
                ...
            }
        }
    ```

3. Update the `resourceExampleRead` function to include a step setting the Resource Identity data into state. Resource Identity data does not have to be set manually, we can make use of the `pluginsdk.SetResourceIdentityData` helper function.

    ```go
        func resourceExampleRead(d *pluginsdk.ResourceData, meta interface{}) error {
            client := meta.(*clients.Client).Service.ExampleClient
            ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
            defer cancel()
    
            id, err := examplepackage.ParseExampleResourceID(d.Id())
            if err != nil {
                return err
            }
            
            ...
            
            // Most of the time we can simply replace the final `return nil` line with the return below.
            // Note: there are a number of resources that conditionally return earlier in the read function before reaching the final `return nil` line.
            // Keep an eye out for these, as they can cause test failures that are tedious to diagnose.
            return pluginsdk.SetResourceIdentityData(d, id)
        }
    ```

4. Add an acceptance test to ensure the identity data is accurately set into state, please reference [Resource Identity Tests](#resource-identity-tests).

## Resource Identity Tests

Just like the schema, Resource Identity tests are entirely generated. This is done by adding a `go:generate` comment. Both untyped and typed resources use the same format. To make this easy to find and modify, place it underneath the imports.

The schema is generated for us by taking different parts of the ID and converting them to snake_case. By default, if the last segment ends in `Name`, it will not be converted to snake case in the schema but rather set to `name`. 

For the tests to generate properly, you will need to specify a combination of `-properties`, `-known-values`, and `-compare-values` inputs. All fields in the ID struct must be mapped to one of these options.

To go through these in order:

- `-properties`: This flag specifies the 1:1 relationship between the Resource Schema and the Resource Identity Schema fields (i.e name, resource_group_name, etc), this would be specified as `name,resource_group_name`. If the schema property name does not match the Resource Identity schema name these should be mapped accordingly. This would be specified as `{id_field_name}:{schema_field_name}`, e.g. `api_management_id:api_management_name`.

- `-known-values`: This flag specifies values that are not exposed in the resource schema, but are present in the Resource Identity schema, e.g. a subscription ID. This would be specified as `{id_field_name}:{known_value}`, e.g. `subscription_id:data.Subscriptions.Primary`.

- `-compare-values`: This flag allows for comparing values that are exposed in the resource schema through another resource ID. This comes up when we use a parent resource ID in the schema but the Resource Identity Schema uses the individual parts of that parent ID. This would be specified as `{id_field_name}:{schema_field_id_name}`, e.g. `virtual_network_name:virtual_network_id`.

Please reference the [Resource Identity Test Generator](../../internal/tools/generator-tests/generators/resource_identity.go) for additional options that are used less frequently.

 ```go
package example
     
import (
    "time"
      
    "github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

// A basic example where the Resource Identity fields map directly to the resource schema
//go:generate go run ../../tools/generator-tests resourceidentity -resource-name example_resource -service-package-name example -properties "name,resource_group_name" -known-values "subscription_id:data.Subscriptions.Primary"

// An example where individual Resource Identity field values exist in a parent ID 
//go:generate go run ../../tools/generator-tests resourceidentity -resource-name example_resource -service-package-name example -properties "name" -compare-values "parent_name:parent_resource_id" -known-values "subscription_id:data.Subscriptions.Primary"
  
type ExampleResource struct{}
  
var _ sdk.ResourceWithIdentity = ExampleResource{}
  
func (r ExampleResource) Identity() resourceids.ResourceId {
    return &examplepackage.ExampleResourceId{}
}
  
func (r ExampleResource) Read() sdk.ResourceFunc {
    return sdk.ResourceFunc{
        Timeout: 5 * time.Minute,
        Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
            client := metadata.Client.Service.ExampleClient
            id, err := examplepackage.ParseExampleResourceID(metadata.ResourceData.Id())
            if err != nil {
                return err
            }
             
            ...
              
            if err := pluginsdk.SetResourceIdentityData(metadata.ResourceData, id); err != nil {
                return err
            }
             
            return metadata.Encode(&model)
        },
    }
}
 ```
