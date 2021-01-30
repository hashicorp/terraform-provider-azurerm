package purview

import (
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/purview/mgmt/2020-12-01-preview/purview"
	"github.com/hashicorp/go-azure-helpers/response"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/location"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/purview/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/purview/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourcePurviewAccount() *schema.Resource {
	return &schema.Resource{
		Create: resourcePurviewAccountCreateUpdate,
		Read:   resourcePurviewAccountRead,
		Update: resourcePurviewAccountCreateUpdate,
		Delete: resourcePurviewAccountDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.AccountName(),
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"location": azure.SchemaLocation(),

			"sku_capacity": {
				Type:     schema.TypeInt,
				Required: true,
				ValidateFunc: validation.IntInSlice([]int{
					4,
					16,
				}),
			},

			"public_network_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},

			"identity": {
				Type:     schema.TypeList,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"principal_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"tenant_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},

			"catalog_endpoint": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"guardian_endpoint": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"scan_endpoint": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"atlas_kafka_endpoint_primary_connection_string": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},

			"atlas_kafka_endpoint_secondary_connection_string": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},

			"tags": tags.Schema(),
		},
	}
}

func resourcePurviewAccountCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Purview.AccountsClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	location := azure.NormalizeLocation(d.Get("location").(string))
	resourceGroup := d.Get("resource_group_name").(string)
	t := d.Get("tags").(map[string]interface{})

	if d.IsNewResource() {
		existing, err := client.Get(ctx, resourceGroup, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of existing Purview Account %q (Resource Group %q): %s", name, resourceGroup, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_purview_account", *existing.ID)
		}
	}

	account := purview.Account{
		AccountProperties: &purview.AccountProperties{},
		Identity: &purview.Identity{
			Type: purview.SystemAssigned,
		},
		Location: &location,
		Sku: &purview.AccountSku{
			Capacity: utils.Int32(d.Get("sku_capacity").(int32)),
			Name:     purview.Standard,
		},
		Tags: tags.Expand(t),
	}

	if d.Get("public_network_enabled").(bool) {
		account.AccountProperties.PublicNetworkAccess = purview.Enabled
	} else {
		account.AccountProperties.PublicNetworkAccess = purview.Disabled
	}

	if _, err := client.CreateOrUpdate(ctx, resourceGroup, name, account); err != nil {
		return fmt.Errorf("Error creating/updating Purview Account %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	resp, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		return fmt.Errorf("Error retrieving Purview Account %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	if resp.ID == nil {
		return fmt.Errorf("Cannot read Purview Account %q (Resource Group %q) ID", name, resourceGroup)
	}

	d.SetId(*resp.ID)

	return resourcePurviewAccountRead(d, meta)
}

func resourcePurviewAccountRead(d *schema.ResourceData, meta interface{}) error {
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

		return fmt.Errorf("Error retrieving Purview Account %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	d.Set("name", id.Name)
	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("location", location.NormalizeNilable(resp.Location))

	if sku := resp.Sku; sku != nil {
		if err := d.Set("sku_capacity", *sku.Capacity); err != nil {
			return fmt.Errorf("Error setting `sku_capacity`: %+v", err)
		}
	}

	if err := d.Set("identity", flattenPurviewAccountIdentity(resp.Identity)); err != nil {
		return fmt.Errorf("Error flattening `identity`: %+v", err)
	}

	if props := resp.AccountProperties; props != nil {
		if err := d.Set("public_network_enabled", props.PublicNetworkAccess == purview.Enabled); err != nil {
			return fmt.Errorf("Error setting `public_network_enabled`: %+v", err)
		}

		if endpoints := resp.Endpoints; endpoints != nil {
			if err := d.Set("catalog_endpoint", *endpoints.Catalog); err != nil {
				return fmt.Errorf("Error setting `catalog_endpoint`: %+v", err)
			}
			if err := d.Set("guardian_endpoint", *endpoints.Guardian); err != nil {
				return fmt.Errorf("Error setting `guardian_endpoint`: %+v", err)
			}
			if err := d.Set("scan_endpoint", *endpoints.Scan); err != nil {
				return fmt.Errorf("Error setting `scan_endpoint`: %+v", err)
			}
		}
	}

	keys, err := client.ListKeys(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		return fmt.Errorf("Error retrieving Purview Account keys %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}
	if primary := keys.AtlasKafkaPrimaryEndpoint; primary != nil {
		if err := d.Set("atlas_kafka_endpoint_primary_connection_string", *primary); err != nil {
			return fmt.Errorf("Error setting `atlas_kafka_endpoint_primary_connection_string`: %+v", err)
		}
	}
	if secondary := keys.AtlasKafkaSecondaryEndpoint; secondary != nil {
		if err := d.Set("atlas_kafka_endpoint_secondary_connection_string", *secondary); err != nil {
			return fmt.Errorf("Error setting `atlas_kafka_endpoint_secondary_connection_string`: %+v", err)
		}
	}

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourcePurviewAccountDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Purview.AccountsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.AccountID(d.Id())
	if err != nil {
		return err
	}

	future, err := client.Delete(ctx, id.ResourceGroup, id.Name)

	if err != nil {
		return fmt.Errorf("Error deleting Purview Account %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		if !response.WasNotFound(future.Response()) {
			return fmt.Errorf("Error deleting Purview Account %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
		}
	}

	return nil
}

func flattenPurviewAccountIdentity(identity *purview.Identity) interface{} {
	if identity == nil {
		return make([]interface{}, 0)
	}

	result := make(map[string]interface{})
	if identity.Type != "" {
		result["type"] = identity.Type
	}
	if identity.PrincipalID != nil {
		result["principal_id"] = *identity.PrincipalID
	}
	if identity.TenantID != nil {
		result["tenant_id"] = *identity.TenantID
	}

	return []interface{}{result}
}
