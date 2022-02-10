package purview

import (
	"fmt"
	"regexp"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/purview/mgmt/2021-07-01/purview"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/location"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/purview/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourcePurviewAccount() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourcePurviewAccountCreateUpdate,
		Read:   resourcePurviewAccountRead,
		Update: resourcePurviewAccountCreateUpdate,
		Delete: resourcePurviewAccountDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.AccountID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: resourcePurviewSchema(),
	}
}

func resourcePurviewAccountCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Purview.AccountsClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	location := azure.NormalizeLocation(d.Get("location").(string))
	t := d.Get("tags").(map[string]interface{})

	id := parse.NewAccountID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))
	if d.IsNewResource() {
		existing, err := client.Get(ctx, id.ResourceGroup, id.Name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}

		if !utils.ResponseWasNotFound(existing.Response) {
			return tf.ImportAsExistsError("azurerm_purview_account", id.ID())
		}
	}

	account := purview.Account{
		AccountProperties: &purview.AccountProperties{},
		Location:          &location,
		Tags:              tags.Expand(t),
	}

	if features.ThreePointOhBeta() {
		expandedIdentity, err := expandIdentity(d.Get("identity").([]interface{}))
		if err != nil {
			return fmt.Errorf("expanding `identity`: %+v", err)
		}

		account.Identity = expandedIdentity
	} else {
		account.Identity = &purview.Identity{
			Type: purview.TypeSystemAssigned,
		}
	}

	if d.Get("public_network_enabled").(bool) {
		account.AccountProperties.PublicNetworkAccess = purview.PublicNetworkAccessEnabled
	} else {
		account.AccountProperties.PublicNetworkAccess = purview.PublicNetworkAccessDisabled
	}

	if v, ok := d.GetOk("managed_resource_group_name"); ok {
		account.AccountProperties.ManagedResourceGroupName = utils.String(v.(string))
	}

	future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.Name, account)
	if err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for create/update of %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourcePurviewAccountRead(d, meta)
}

func resourcePurviewAccountRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Purview.AccountsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.AccountID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.Name)
	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("location", location.NormalizeNilable(resp.Location))

	if err := d.Set("identity", flattenIdentity(resp.Identity)); err != nil {
		return fmt.Errorf("flattening `identity`: %+v", err)
	}

	if err := d.Set("managed_resources", flattenPurviewAccountManagedResources(resp.ManagedResources)); err != nil {
		return fmt.Errorf("flattening `managed_resources`: %+v", err)
	}

	if props := resp.AccountProperties; props != nil {
		d.Set("public_network_enabled", props.PublicNetworkAccess == purview.PublicNetworkAccessEnabled)

		managedResourceGroupName := ""
		if props.ManagedResourceGroupName != nil {
			managedResourceGroupName = *props.ManagedResourceGroupName
		}
		d.Set("managed_resource_group_name", managedResourceGroupName)

		if endpoints := resp.Endpoints; endpoints != nil {
			d.Set("catalog_endpoint", endpoints.Catalog)
			d.Set("guardian_endpoint", endpoints.Guardian)
			d.Set("scan_endpoint", endpoints.Scan)
		}
	}

	keys, err := client.ListKeys(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		return fmt.Errorf("retrieving Keys for %s: %+v", *id, err)
	}
	d.Set("atlas_kafka_endpoint_primary_connection_string", keys.AtlasKafkaPrimaryEndpoint)
	d.Set("atlas_kafka_endpoint_secondary_connection_string", keys.AtlasKafkaSecondaryEndpoint)

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourcePurviewAccountDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Purview.AccountsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.AccountID(d.Id())
	if err != nil {
		return err
	}

	future, err := client.Delete(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for deletion of %s: %+v", *id, err)
	}

	return nil
}

func expandIdentity(input []interface{}) (*purview.Identity, error) {
	expanded, err := identity.ExpandSystemAssigned(input)
	if err != nil {
		return nil, err
	}

	return &purview.Identity{
		Type: purview.Type(string(expanded.Type)),
	}, nil
}

func flattenIdentity(input *purview.Identity) interface{} {
	var transition *identity.SystemAssigned

	if input != nil {
		transition = &identity.SystemAssigned{
			Type: identity.Type(string(input.Type)),
		}
		if input.PrincipalID != nil {
			transition.PrincipalId = *input.PrincipalID
		}
		if input.TenantID != nil {
			transition.TenantId = *input.TenantID
		}
	}

	return identity.FlattenSystemAssigned(transition)
}

func flattenPurviewAccountManagedResources(managedResources *purview.AccountPropertiesManagedResources) interface{} {
	if managedResources == nil {
		return make([]interface{}, 0)
	}

	resourceGroup := ""
	if managedResources.ResourceGroup != nil {
		resourceGroup = *managedResources.ResourceGroup
	}
	storageAccount := ""
	if managedResources.StorageAccount != nil {
		storageAccount = *managedResources.StorageAccount
	}
	eventHubNamespace := ""
	if managedResources.EventHubNamespace != nil {
		eventHubNamespace = *managedResources.EventHubNamespace
	}
	return []interface{}{
		map[string]interface{}{
			"resource_group_id":      resourceGroup,
			"storage_account_id":     storageAccount,
			"event_hub_namespace_id": eventHubNamespace,
		},
	}
}

func resourcePurviewSchema() map[string]*pluginsdk.Schema {
	schema := map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.StringMatch(
				regexp.MustCompile(`^[a-zA-Z0-9][-a-zA-Z0-9]{1,61}[a-zA-Z0-9]$`),
				"The Purview account name must be between 3 and 63 characters long, it can contain only letters, numbers and hyphens, and the first and last characters must be a letter or number."),
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"public_network_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  true,
		},

		"managed_resource_group_name": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Computed:     true,
			ForceNew:     true,
			ValidateFunc: azure.ValidateResourceGroupName,
		},

		"identity": func() *schema.Schema {
			// TODO: document that this will become required in 3.0
			if features.ThreePointOhBeta() {
				return commonschema.SystemAssignedIdentityRequired()
			}

			return commonschema.SystemAssignedIdentityComputed()
		}(),

		"managed_resources": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"resource_group_id": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
					"storage_account_id": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
					"event_hub_namespace_id": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
				},
			},
		},

		"catalog_endpoint": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"guardian_endpoint": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"scan_endpoint": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"atlas_kafka_endpoint_primary_connection_string": {
			Type:      pluginsdk.TypeString,
			Computed:  true,
			Sensitive: true,
		},

		"atlas_kafka_endpoint_secondary_connection_string": {
			Type:      pluginsdk.TypeString,
			Computed:  true,
			Sensitive: true,
		},

		"tags": tags.Schema(),
	}

	if !features.ThreePointOhBeta() {

		schema["sku_name"] = &pluginsdk.Schema{
			Type:       pluginsdk.TypeString,
			Optional:   true,
			Deprecated: "This property can no longer be specified on create/update, it can only be updated by creating a support ticket at Azure",
		}
	}

	return schema
}
