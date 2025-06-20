---
applyTo: "internal/**/*.go"
description: This document outlines the coding patterns for Go files in the Terraform AzureRM provider repository. It includes patterns for resource implementation, client management, schema design, and Azure SDK integration.
---

## Coding Patterns for Terraform AzureRM Provider
Given below are the coding patterns for the Terraform AzureRM provider which **MUST** be followed.

### Resource Implementation Pattern

#### Standard Resource Structure
```go
func resourceServiceNameResourceType() *pluginsdk.Resource {
    return &pluginsdk.Resource{
        Create: resourceServiceNameResourceTypeCreate,
        Read:   resourceServiceNameResourceTypeRead,
        Update: resourceServiceNameResourceTypeUpdate,
        Delete: resourceServiceNameResourceTypeDelete,

        Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
            _, err := parse.ServiceNameResourceTypeID(id)
            return err
        }),

        Timeouts: &pluginsdk.ResourceTimeout{
            Create: pluginsdk.DefaultTimeout(30 * time.Minute),
            Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
            Update: pluginsdk.DefaultTimeout(30 * time.Minute),
            Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
        },

        Schema: resourceServiceNameResourceTypeSchema(),
    }
}

func resourceServiceNameResourceTypeSchema() map[string]*pluginsdk.Schema {
    return map[string]*pluginsdk.Schema{
        "name": {
            Type:         pluginsdk.TypeString,
            Required:     true,
            ForceNew:     true,
            ValidateFunc: validation.StringIsNotEmpty,
        },

        "location": commonschema.Location(),

        "resource_group_name": commonschema.ResourceGroupName(),

        "tags": commonschema.Tags(),
    }
}
```

#### CRUD Operation Pattern
```go
func resourceServiceNameResourceTypeCreate(ctx context.Context, d *pluginsdk.ResourceData, meta interface{}) error {
    client := meta.(*clients.Client).ServiceName.ResourceTypeClient
    subscriptionId := meta.(*clients.Client).Account.SubscriptionId

    // Parse input parameters
    name := d.Get("name").(string)
    resourceGroupName := d.Get("resource_group_name").(string)
    location := azure.NormalizeLocation(d.Get("location").(string))

    // Create resource ID
    id := parse.NewServiceNameResourceTypeID(subscriptionId, resourceGroupName, name)

    // Check for existing resource
    existing, err := client.Get(ctx, id)
    if err != nil && !response.WasNotFound(existing.HttpResponse) {
        return fmt.Errorf("checking for existing %s: %w", id, err)
    }
    if !response.WasNotFound(existing.HttpResponse) {
        return tf.ImportAsExistsError("azurerm_service_name_resource_type", id.ID())
    }

    // Build parameters
    parameters := servicenametype.ResourceType{
        Location: location,
        Properties: &servicenametype.ResourceTypeProperties{
            // Add properties here
        },
    }

    // Handle tags
    if tagsRaw := d.Get("tags"); tagsRaw != nil {
        parameters.Tags = tags.Expand(tagsRaw.(map[string]interface{}))
    }

    // Create resource
    if err := client.CreateOrUpdateThenPoll(ctx, id, parameters); err != nil {
        return fmt.Errorf("creating %s: %w", id, err)
    }

    d.SetId(id.ID())
    return resourceServiceNameResourceTypeRead(ctx, d, meta)
}

func resourceServiceNameResourceTypeRead(ctx context.Context, d *pluginsdk.ResourceData, meta interface{}) error {
    client := meta.(*clients.Client).ServiceName.ResourceTypeClient

    id, err := parse.ServiceNameResourceTypeID(d.Id())
    if err != nil {
        return err
    }

    resp, err := client.Get(ctx, *id)
    if err != nil {
        if response.WasNotFound(resp.HttpResponse) {
            log.Printf("[DEBUG] %s was not found - removing from state!", *id)
            d.SetId("")
            return nil
        }
        return fmt.Errorf("retrieving %s: %w", *id, err)
    }

    d.Set("name", id.ResourceTypeName)
    d.Set("resource_group_name", id.ResourceGroupName)

    if model := resp.Model; model != nil {
        d.Set("location", azure.NormalizeLocation(model.Location))
        
        if props := model.Properties; props != nil {
            // Set properties
        }

        if err := tags.FlattenAndSet(d, model.Tags); err != nil {
            return err
        }
    }

    return nil
}
```

### Client Management Pattern

#### Client Registration
```go
type Client struct {
    ResourceTypeClient *servicenametype.ResourceTypeClient
}

func NewClient(o *common.ClientOptions) *Client {
    resourceTypeClient := servicenametype.NewResourceTypeClientWithBaseURI(o.ResourceManagerEndpoint)
    o.ConfigureClient(&resourceTypeClient.Client, o.ResourceManagerAuthorizer)

    return &Client{
        ResourceTypeClient: &resourceTypeClient,
    }
}
```

### Data Source Pattern

#### Standard Data Source Structure
```go
func dataSourceServiceNameResourceType() *pluginsdk.Resource {
    return &pluginsdk.Resource{
        Read: dataSourceManagedDiskRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:     pluginsdk.TypeString,
				Required: true,
			},

			"resource_group_name": commonschema.ResourceGroupNameForDataSource(),

            "disk_encryption_set_id": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"encryption_settings": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"enabled": {
							Type:     pluginsdk.TypeBool,
							Computed: true,
						},

						"disk_encryption_key": {
							Type:     pluginsdk.TypeList,
							Computed: true,
							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"secret_url": {
										Type:     pluginsdk.TypeString,
										Computed: true,
									},

									"source_vault_id": {
										Type:     pluginsdk.TypeString,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},

            "tags": commonschema.TagsDataSource(),
        },
    }
}

func dataSourceServiceNameResourceTypeRead(ctx context.Context, d *pluginsdk.ResourceData, meta interface{}) error {
    client := meta.(*clients.Client).Compute.DisksClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := commonids.NewManagedDiskID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))
	resp, err := client.Get(ctx, id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return fmt.Errorf("%s was not found", id)
		}
		return fmt.Errorf("making Read request on %s: %s", id, err)
	}

	d.SetId(id.ID())

	d.Set("name", id.DiskName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		d.Set("location", location.Normalize(model.Location))

		storageAccountType := ""
		if sku := model.Sku; sku != nil {
			storageAccountType = string(*sku.Name)
		}
		d.Set("storage_account_type", storageAccountType)

		if props := model.Properties; props != nil {
			creationData := props.CreationData

			diskEncryptionSetId := ""
			if props.Encryption != nil && props.Encryption.DiskEncryptionSetId != nil {
				diskEncryptionSetId = *props.Encryption.DiskEncryptionSetId
			}
			d.Set("disk_encryption_set_id", diskEncryptionSetId)

			if err := d.Set("encryption_settings", flattenManagedDiskEncryptionSettings(props.EncryptionSettingsCollection)); err != nil {
				return fmt.Errorf("setting `encryption_settings`: %+v", err)
			}
		}

		if err := tags.FlattenAndSet(d, model.Tags); err != nil {
			return err
		}
	}

	return nil
}
```

### Schema Design Patterns

#### Complex Schema with Validation
```go
func resourceComplexSchema() map[string]*pluginsdk.Schema {
    return map[string]*pluginsdk.Schema{
        "name": {
            Type:     pluginsdk.TypeString,
            Required: true,
            ForceNew: true,
            ValidateFunc: validation.All(
                validation.StringIsNotEmpty,
                validation.StringMatch(regexp.MustCompile(`^[a-zA-Z0-9-]+$`), "name can only contain alphanumeric characters and hyphens"),
                validation.StringLenBetween(1, 64),
            ),
        },

        "configuration": {
            Type:     pluginsdk.TypeList,
            Optional: true,
            MaxItems: 1,
            Elem: &pluginsdk.Resource{
                Schema: map[string]*pluginsdk.Schema{
                    "enabled": {
                        Type:     pluginsdk.TypeBool,
                        Optional: true,
                        Default:  true,
                    },
                    "settings": {
                        Type:     pluginsdk.TypeMap,
                        Optional: true,
                        Elem: &pluginsdk.Schema{
                            Type: pluginsdk.TypeString,
                        },
                    },
                },
            },
        },

        "network_configuration": {
            Type:     pluginsdk.TypeSet,
            Optional: true,
            Elem: &pluginsdk.Resource{
                Schema: map[string]*pluginsdk.Schema{
                    "subnet_id": {
                        Type:         pluginsdk.TypeString,
                        Required:     true,
                        ValidateFunc: commonids.ValidateSubnetID,
                    },
                    "private_ip_address": {
                        Type:         pluginsdk.TypeString,
                        Optional:     true,
                        ValidateFunc: validation.IsIPv4Address,
                    },
                },
            },
        },
    }
}
```

### Error Handling Patterns

#### Standard Error Handling
```go
// Check for existing resource
existing, err := client.Get(ctx, id)
if err != nil {
    if !response.WasNotFound(existing.HttpResponse) {
        return fmt.Errorf("checking for existing %s: %w", id, err)
    }
}

// Handle resource not found in Read operation
if response.WasNotFound(resp.HttpResponse) {
    log.Printf("[DEBUG] %s was not found - removing from state", id)
    d.SetId("")
    return nil
}

// Handle throttling
if response.WasThrottled(resp.HttpResponse) {
    return resource.RetryableError(fmt.Errorf("request was throttled, retrying"))
}

// Handle long-running operations
if err := client.CreateOrUpdateThenPoll(ctx, id, parameters); err != nil {
    return fmt.Errorf("creating %s: %w", id, err)
}
```

### Resource ID Parsing Pattern

#### ID Parser Implementation
```go
type ServiceNameResourceTypeId struct {
    SubscriptionId      string
    ResourceGroupName   string
    ResourceTypeName    string
}

func NewServiceNameResourceTypeID(subscriptionId, resourceGroupName, resourceTypeName string) ServiceNameResourceTypeId {
    return ServiceNameResourceTypeId{
        SubscriptionId:    subscriptionId,
        ResourceGroupName: resourceGroupName,
        ResourceTypeName:  resourceTypeName,
    }
}

func (id ServiceNameResourceTypeId) ID() string {
    fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.ServiceName/resourceTypes/%s"
    return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.ResourceTypeName)
}

func ServiceNameResourceTypeID(input string) (*ServiceNameResourceTypeId, error) {
    parser := resourceids.NewParserFromResourceIdType(ServiceNameResourceTypeId{})
    parsed, err := parser.Parse(input, false)
    if err != nil {
        return nil, fmt.Errorf("parsing %q: %w", input, err)
    }

    var ok bool
    id := ServiceNameResourceTypeId{}

    if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
        return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
    }

    if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
        return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
    }

    if id.ResourceTypeName, ok = parsed.Parsed["resourceTypeName"]; !ok {
        return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceTypeName", *parsed)
    }

    return &id, nil
}
```

### Testing Patterns

#### Acceptance Test Pattern
```go
func TestAccServiceNameResourceType_basic(t *testing.T) {
    data := acceptance.BuildTestData(t, "azurerm_service_name_resource_type", "test")
    r := ServiceNameResourceTypeTestResource{}

    data.ResourceTest(t, r, []acceptance.TestStep{
        {
            Config: r.basic(data),
            Check: acceptance.ComposeTestCheckFunc(
                check.That(data.ResourceName).ExistsInAzure(r),
                check.That(data.ResourceName).Key("name").HasValue(data.RandomString),
                check.That(data.ResourceName).Key("location").HasValue(data.Locations.Primary),
            ),
        },
        data.ImportStep(),
    })
}

func (r ServiceNameResourceTypeTestResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
    id, err := parse.ServiceNameResourceTypeID(state.ID)
    if err != nil {
        return nil, err
    }

    resp, err := clients.ServiceName.ResourceTypeClient.Get(ctx, *id)
    if err != nil {
        return nil, fmt.Errorf("retrieving %s: %w", *id, err)
    }

    return utils.Bool(resp.Model != nil), nil
}

func (r ServiceNameResourceTypeTestResource) basic(data acceptance.TestData) string {
    return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_service_name_resource_type" "test" {
  name                = "acctest-%s"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}
```
