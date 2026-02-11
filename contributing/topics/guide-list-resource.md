# Guide: List Resource

This guide covers how to add a List Resource for an existing resource, using `azurerm_network_profile` as an example. For more information on Lists, see [Resources - List](https://developer.hashicorp.com/terraform/plugin/sdkv2/resources/list).

## Prerequisites

Before adding a List Resource, the resource must have Resource Identity implemented. For more information on implementing Resource Identity see [Guide: Resource Identity](guide-resource-identity.md).

## Adding List Resource

> **Note:** There are some minor differences between the implementation of a List Resource for an untyped or typed resource. These differences are highlighted in separated code snippets.

1. In the resource, refactor the Read function to have a separate flatten function containing only the logic to set the attributes into state. This will be used by both the Read function and later in the List Resource.<br><br>

    For untyped resources:
    ```go
    func resourceNetworkProfileFlatten(d *pluginsdk.ResourceData, id *networkprofiles.NetworkProfileId, profile *networkprofiles.NetworkProfile) error {
        d.Set("name", id.NetworkProfileName)
        d.Set("resource_group_name", id.ResourceGroupName)
    
        if profile != nil {
            if props := profile.Properties; props != nil {
                cniConfigs := flattenNetworkProfileContainerNetworkInterface(props.ContainerNetworkInterfaceConfigurations)
                if err := d.Set("container_network_interface", cniConfigs); err != nil {
                    return fmt.Errorf("setting `container_network_interface`: %+v", err)
                }
    
                cniIDs := flattenNetworkProfileContainerNetworkInterfaceIDs(props.ContainerNetworkInterfaces)
                if err := d.Set("container_network_interface_ids", cniIDs); err != nil {
                    return fmt.Errorf("setting `container_network_interface_ids`: %+v", err)
                }
            }
            d.Set("location", location.NormalizeNilable(profile.Location))
            if err := tags.FlattenAndSet(d, profile.Tags); err != nil {
                return err
            }
        }
        return pluginsdk.SetResourceIdentityData(d, id)
    }
    ```

    For typed resources:
    ```go
    func (ExampleResource) flatten(metadata sdk.ResourceMetaData, id *example.ExampleId, model *example.ExampleModel) error {
        // Instantiate state, set any fields with known values (e.g. ones we can derive from the ID)
        state := ExampleResourceModel{
            Name: id.ExampleResourceName
            ResourceGroupName: id.ResourceGroupName
        }
   
        if model != nil {
            state.Location = location.Normalize(model.Location)
            
            if props := model.Properties; props != nil {
                // Set remaining properties into the Resource Model (`state`)   
            }
        }
   
        // Set the Resource Identity Data
        if err := pluginsdk.SetResourceIdentityData(metadata.ResourceData, id); err != nil {
            return err
        }
   
        return metadata.Encode(&state)
    }
    ```

2. Create a new file for the List Resource (for example, `network_profile_resource_list.go`) and scaffold the empty resource:<br><br>

    For untyped resources:
    ```go
    type NetworkProfileListResource struct{}
    
    var _ sdk.FrameworkListWrappedResource = new(NetworkProfileListResource)
    
    func (NetworkProfileListResource) ResourceFunc() *pluginsdk.Resource {
        return resourceNetworkProfile()
    }

    // set this with a const from the resource containing the resource name eg, `azurerm_network_profile`
    func (r NetworkProfileListResource) Metadata(_ context.Context, _ resource.MetadataRequest, response *resource.MetadataResponse) {
        response.TypeName = azureNetworkProfileResourceName
    }
    ```

    For typed resources:
    ```go
    type ExampleListResource struct{}

    var _ sdk.FrameworkListWrappedResource = new(ExampleListResource)
    
    func (ExampleListResource) ResourceFunc() *pluginsdk.Resource {
        // Use the `sdk.WrappedResource` helper to convert a typed resource into `*pluginsdk.Resource`
        return sdk.WrappedResource(ExampleResource{})
    }

    // Set the name using the `ResourceType()` function
    func (ExampleListResource) Metadata(_ context.Context, _ resource.MetadataRequest, response *resource.MetadataResponse) {
        response.TypeName = ExampleResource{}.ResourceType()
    }
    ```

3. Define any List Resource specific configuration options. This step can be omitted if using the DefaultListModel (which includes `subscription_id` and `resource_group_name`). However, other resources may have different configuration options that need to be defined here and would look something like this:

    ```go
    type NetworkProfileListModel struct {
        SubscriptionId    types.String `tfsdk:"subscription_id"`
        ResourceGroupName types.String `tfsdk:"resource_group_name"`
    }
    
    func (NetworkProfileListResource) ListResourceConfigSchema(_ context.Context, _ list.ListResourceSchemaRequest, response *list.ListResourceSchemaResponse) {
        response.Schema = schema.Schema{
            Attributes: map[string]schema.Attribute{
                "subscription_id": schema.StringAttribute{
                    Optional: true,
                    Validators: []validator.String{
                        typehelpers.WrappedStringValidator{
                            Func: commonids.ValidateSubscriptionID,
                        },
                    },
                },
                "resource_group_name": schema.StringAttribute{
                    Optional: true,
                    Validators: []validator.String{
                        typehelpers.WrappedStringValidator{
                            Func: resourcegroups.ValidateName,
                        },
                    },
                },
            },
        }
    }
    ```

4. Implement the List function.<br><br>

    For untyped resources:

    ```go
    func (NetworkProfileListResource) List(ctx context.Context, request list.ListRequest, stream *list.ListResultsStream, metadata sdk.ResourceMetadata) {
    
        client := metadata.Client.Network.NetworkProfiles
    
        // Read the list config data into the model
        var data sdk.DefaultListModel
        diags := request.Config.Get(ctx, &data)
        if diags.HasError() {
            stream.Results = list.ListResultsStreamDiagnostics(diags)
            return
        }
    
        // Initialize a list for the results of the API request
        results := make([]networkprofiles.NetworkProfile, 0)
    
        subscriptionID := metadata.SubscriptionId
        if !data.SubscriptionId.IsNull() {
            subscriptionID = data.SubscriptionId.ValueString()
        }
    
        // Make the request based on which list parameters have been set in the config
        switch {
        case !data.ResourceGroupName.IsNull():
            resp, err := client.ListComplete(ctx, commonids.NewResourceGroupID(subscriptionID, data.ResourceGroupName.ValueString()))
            if err != nil {
                sdk.SetResponseErrorDiagnostic(stream, fmt.Sprintf("listing `%s`", azureNetworkProfileResourceName), err)
                return
            }
    
            results = resp.Items
        default:
            resp, err := client.ListAllComplete(ctx, commonids.NewSubscriptionID(subscriptionID))
            if err != nil {
                sdk.SetResponseErrorDiagnostic(stream, fmt.Sprintf("listing `%s`", azureNetworkProfileResourceName), err)
                return
            }
    
            results = resp.Items
        }
    
        // Define the function that will push results into the stream 
        stream.Results = func(push func(list.ListResult) bool) {
            for _, profile := range results {
            
                // Initialize a new result object for each resource in the list
                result := request.NewListResult(ctx)
                
                // Set the display name of the item as the resource name
                result.DisplayName = pointer.From(profile.Name)
    
                // Create a new ResourceData object to hold the state of the resource
                rd := resourceNetworkProfile().Data(&terraform.InstanceState{})
                
                // Set the ID of the resource for the ResourceData object
                id, err := networkprofiles.ParseNetworkProfileID(pointer.From(profile.Id))
                if err != nil {
                    sdk.SetErrorDiagnosticAndPushListResult(result, push, "parsing Network Profile ID", err)
                    return
                }
                rd.SetId(id.ID())
    
                // Use the resource flatten function to set the attributes into the resource state
                if err := resourceNetworkProfileFlatten(rd, id, &profile); err != nil {
                    sdk.SetErrorDiagnosticAndPushListResult(result, push, fmt.Sprintf("encoding `%s` resource data", azureNetworkProfileResourceName), err)
                    return
                }
    
               // Convert and set the identity and resource state into the result
               sdk.EncodeListResult(ctx, rd, &result)
               if result.Diagnostics.HasError() {
                   push(result)
                   return
               }

                if !push(result) {
                    return
                }
            }
        }
    }
    ```
   
    For typed resources:
    ```go
    func (ExampleListResource) List(ctx context.Context, request list.ListRequest, stream *list.ListResultsStream, metadata sdk.ResourceMetadata) {
        client := metadata.Client.Example.ExampleResourceClient

        var data sdk.DefaultListModel
        diags := request.Config.Get(ctx, &data)
        if diags.HasError() {
            stream.Results = list.ListResultsStreamDiagnostics(diags)
            return
        }

        var results []example.ExampleModel

        subscriptionID := metadata.SubscriptionId
        if !data.SubscriptionId.IsNull() {
            subscriptionID = data.SubscriptionId.ValueString()
        }

        r := ExampleResource{}
   
        switch {
        case !data.ResourceGroupName.IsNull():
            resp, err := client.ListByResourceGroupComplete(ctx, commonids.NewResourceGroupID(subscriptionID, data.ResourceGroupName.ValueString()))
            if err != nil {
                sdk.SetResponseErrorDiagnostic(stream, fmt.Sprintf("listing `%s`", r.ResourceType()), err)
                return
            }

            results = resp.Items
        default:
            resp, err := client.ListComplete(ctx, commonids.NewSubscriptionID(subscriptionID))
            if err != nil {
                sdk.SetResponseErrorDiagnostic(stream, fmt.Sprintf("listing `%s`", r.ResourceType()), err)
                return
            }

        results = resp.Items
    }

    stream.Results = func(push func(list.ListResult) bool) {
        for _, exampleResult := range results {
            result := request.NewListResult(ctx)
            result.DisplayName = pointer.From(exampleResult.Name)

            id, err := example.ParseExampleID(pointer.From(exampleResult.Id))
            if err != nil {
                sdk.SetErrorDiagnosticAndPushListResult(result, push, "parsing Example ID", err)
                return
            }

            // Instantiate a new ResourceMetaData object to leverage the resource's `flatten` function
            // which uses the `(ResourceMetaData).Encode()` function to populate the resource state.
            rmd := sdk.NewResourceMetaData(metadata.Client, r)
            rmd.SetID(id)

            if err := r.flatten(rmd, id, &exampleResult); err != nil {
                sdk.SetErrorDiagnosticAndPushListResult(result, push, fmt.Sprintf("encoding `%s` resource data", r.ResourceType()), err)
                return
            }

            sdk.EncodeListResult(ctx, rmd.ResourceData, &result)
            if result.Diagnostics.HasError() {
                push(result)
                return
            }

            if !push(result) {
                return
            }
        }
    }
    ```

5. Register the new List Resource

   List Resources are registered within the `registration.go` within each Service Package - and should look something like this:

    ```
    package network
    
    import "github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
    
    type Registration struct{}
    
    var _ sdk.FrameworkServiceRegistration = Registration{}
    
    // ...
    
    // Resources returns a list of List Resources supported by this Service
    func (r Registration) ListResources() []sdk.FrameworkListWrappedResource {
        return []sdk.FrameworkListWrappedResource{
            NetworkProfileListResource{},
            }
    }
    ```

6. Add Acceptance Tests for this List Resource

    Create a new acceptance test file for the List Resource (for example, `network_profile_resource_list_test.go`) and add tests to cover the List Resource functionality. The test should provision any prerequisite resources and multiple resources of the type of List Resource we want to test.

    The test should look something like this:

    ```
    package network_test
    
    import (
        "context"
        "fmt"
        "testing"
    
        "github.com/hashicorp/terraform-plugin-testing/helper/resource"
        "github.com/hashicorp/terraform-plugin-testing/querycheck"
        "github.com/hashicorp/terraform-plugin-testing/tfversion"
        "github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
        "github.com/hashicorp/terraform-provider-azurerm/internal/provider/framework"
    )
    
    func TestAccNetworkProfile_list_basic(t *testing.T) {
        r := NetworkProfileResource{}
        listResourceAddress := "azurerm_network_profile.list"
    
        data := acceptance.BuildTestData(t, "azurerm_network_profile", "test1")
    
        resource.Test(t, resource.TestCase{
            TerraformVersionChecks: []tfversion.TerraformVersionCheck{
                tfversion.SkipBelow(tfversion.Version1_14_0),
            },
            ProtoV5ProviderFactories: framework.ProtoV5ProviderFactoriesInit(context.Background(), "azurerm"),
            Steps: []resource.TestStep{
                {
                    Config: r.basicList(data), // provision multiple resources
                },
                {
                    Query:  true,
                    Config: r.basicQuery(),
                    QueryResultChecks: []querycheck.QueryResultCheck{
                        querycheck.ExpectLengthAtLeast(listResourceAddress, 3), // expect at least the 3 we created
                    },
                },
                {
                    Query:  true,
                    Config: r.basicQueryByResourceGroupName(data),
                    QueryResultChecks: []querycheck.QueryResultCheck{
                        querycheck.ExpectLength(listResourceAddress, 3), // expect exactly the 3 we created in that resource group
                    },
                },
            },
        })
    }
    
    // provision multiple Network Profile resources for testing
    func (r NetworkProfileResource) basicList(data acceptance.TestData) string {
        return fmt.Sprintf(`
    provider "azurerm" {
      features {}
    }
    
    // Prerequisite Resources ....
    
    resource "azurerm_network_profile" "test" {
      // Where possible, use the `count` meta argument to provision multiple resources to query
      count = 3
   
      name                = "acctestnetprofile${count.index}-%[1]d"
      location            = azurerm_resource_group.test.location
      resource_group_name = azurerm_resource_group.test.name
    
      container_network_interface {
        name = "acctesteth-%[1]d"
    
        ip_configuration {
          name      = "acctestipconfig-%[1]d"
          subnet_id = azurerm_subnet.test.id
        }
      }
    }
    `, data.RandomInteger, data.Locations.Primary)
    }
    
    // define the basic list query for testing
    func (r NetworkProfileResource) basicQuery() string {
        return `
    list "azurerm_network_profile" "list" {
      provider = azurerm
      config {}
    }
    `
    }
    
    // define the list query for testing by resource group name
    func (r NetworkProfileResource) basicQueryByResourceGroupName(data acceptance.TestData) string {
        return fmt.Sprintf(`
    list "azurerm_network_profile" "list" {
      provider = azurerm
      config {
        resource_group_name = "acctestRG-%[1]d"
      }
    }
    `, data.RandomInteger)
    }
    
    ```

7. Add documentation for this List Resource

    Documentation should be written manually and added to the `./website/docs/list-resources/` folder.

    It should include an example, arguments reference, and look something like this:

    ````markdown
    ---
    subcategory: "Network"
    layout: "azurerm"
    page_title: "Azure Resource Manager: azurerm_network_profile"
    description: |-
    Lists Network Profile resources.
    ---
    
    # List resource: azurerm_network_profile
    
    Lists Network Profile resources.
    
    ## Example Usage
    
    ### List all Network Profiles in the subscription
    
    ```hcl
    list "azurerm_network_profile" "example" {
      provider = azurerm
      config {}
    }
    ```
    
    ### List all Network Profiles in a specific resource group
    
    ```hcl
    list "azurerm_network_profile" "example" {
      provider = azurerm
      config {
        resource_group_name = "example-rg"
      }
    }
    ```
    
    ## Argument Reference
    
    This list resource supports the following arguments:
    
    * `resource_group_name` - (Optional) The name of the resource group to query.
    
    * `subscription_id` - (Optional) The Subscription ID to query. Defaults to the value specified in the Provider Configuration.
    ````

## Known Issues and Considerations

### Cancelled Context

Some resources need to send additional API requests in the flatten function, these API requests require a valid context (i.e. not cancelled or done). However, due to the way the List resources function, the context provided will be cancelled by the time Terraform calls the iterator (`stream.Results`).

In this scenario, you must instantiate a new context within the iterator using the deadline from the provided context, this should look like the below:

```go
func (ExampleListResource) List(ctx context.Context, request list.ListRequest, stream *list.ListResultsStream, metadata sdk.ResourceMetadata) {
    ...

    // retrieve the deadline from the supplied context
    deadline, ok := ctx.Deadline()
    if !ok {
        // This *should* never happen given the List Wrapper instantiates a context with a timeout
        sdk.SetResponseErrorDiagnostic(stream, "internal-error", "context had no deadline")
        return
    }
    
    stream.Result = func(push func(list.ListResult) bool) {
        // Instantiate a new context based on the deadline retrieved earlier
        ctx, cancel := context.WithDeadline(context.Background(), deadline)
        defer cancel()
        
        for _, example := range results {
            // Remaining logic to retrieve and set the resource data
        }
    }
}
```