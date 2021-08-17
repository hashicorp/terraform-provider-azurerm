package purview

import (
	"fmt"
	"regexp"
	"time"

	"github.com/hashicorp/go-azure-helpers/response"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/identity"
	"github.com/hashicorp/terraform-provider-azurerm/internal/location"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/purview/sdk/2020-12-01-preview/account"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type purviewAccountIdentity = identity.SystemAssigned

func resourcePurviewAccount() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourcePurviewAccountCreateUpdate,
		Read:   resourcePurviewAccountRead,
		Update: resourcePurviewAccountCreateUpdate,
		Delete: resourcePurviewAccountDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := account.ParseAccountID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringMatch(
					regexp.MustCompile(`^[a-zA-Z0-9][-a-zA-Z0-9]{1,61}[a-zA-Z0-9]$`),
					"The Purview account name must be between 3 and 63 characters long, it can contain only letters, numbers and hyphens, and the first and last characters must be a letter or number."),
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"location": azure.SchemaLocation(),

			"sku_name": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"Standard_4",
					"Standard_16",
				}, false),
			},

			"public_network_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  true,
			},

			"identity": purviewAccountIdentity{}.Schema(),

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
		},
	}
}

func resourcePurviewAccountCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Purview.AccountsClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := account.NewAccountID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))
	if d.IsNewResource() {
		existing, err := client.Get(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_purview_account", id.ID())
		}
	}

	publicNetworkAccess := account.PublicNetworkAccessEnabled
	if d.Get("public_network_enabled").(bool) {
		publicNetworkAccess = account.PublicNetworkAccessDisabled
	}

	identity, err := expandPurviewIdentity(d.Get("identity").([]interface{}))
	if err != nil {
		return err
	}
	model := account.Account{
		Properties: &account.AccountProperties{
			PublicNetworkAccess: &publicNetworkAccess,
		},
		Identity: identity,
		Location: utils.String(location.Normalize(d.Get("location").(string))),
		Sku:      expandPurviewSkuName(d),
		Tags:     expandTags(d.Get("tags").(map[string]interface{})),
	}

	if err := client.CreateOrUpdateThenPoll(ctx, id, model); err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourcePurviewAccountRead(d, meta)
}

func resourcePurviewAccountRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Purview.AccountsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := account.ParseAccountID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}
	keys, err := client.ListKeys(ctx, *id)
	if err != nil {
		return fmt.Errorf("retrieving Keys for %s: %+v", *id, err)
	}

	d.Set("name", id.Name)
	d.Set("resource_group_name", id.ResourceGroup)

	if model := resp.Model; model != nil {
		d.Set("location", location.NormalizeNilable(model.Location))
		d.Set("sku_name", flattenPurviewSkuName(model.Sku))

		if err := d.Set("identity", flattenPurviewAccountIdentity(model.Identity)); err != nil {
			return fmt.Errorf("flattening `identity`: %+v", err)
		}

		if props := model.Properties; props != nil {
			publicNetworkEnabled := true
			if props.PublicNetworkAccess != nil {
				publicNetworkEnabled = *props.PublicNetworkAccess == account.PublicNetworkAccessEnabled
			}
			d.Set("public_network_enabled", publicNetworkEnabled)

			if endpoints := props.Endpoints; endpoints != nil {
				d.Set("catalog_endpoint", endpoints.Catalog)
				d.Set("guardian_endpoint", endpoints.Guardian)
				d.Set("scan_endpoint", endpoints.Scan)
			}
		}

		if err := tags.FlattenAndSet(d, flattenTags(model.Tags)); err != nil {
			return err
		}
	}

	if model := keys.Model; model != nil {
		d.Set("atlas_kafka_endpoint_primary_connection_string", model.AtlasKafkaPrimaryEndpoint)
		d.Set("atlas_kafka_endpoint_secondary_connection_string", model.AtlasKafkaSecondaryEndpoint)
	}

	return nil
}

func resourcePurviewAccountDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Purview.AccountsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := account.ParseAccountID(d.Id())
	if err != nil {
		return err
	}

	if err := client.DeleteThenPoll(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return nil
}

func expandPurviewSkuName(d *pluginsdk.ResourceData) *account.AccountSku {
	sku := d.Get("sku_name").(string)

	if len(sku) == 0 {
		return nil
	}

	name, capacity, err := azure.SplitSku(sku)
	if err != nil {
		return nil
	}

	skuName := account.Name(name)
	return &account.AccountSku{
		Name:     &skuName,
		Capacity: utils.Int64(int64(capacity)),
	}
}

func flattenPurviewSkuName(input *account.AccountSku) string {
	if input == nil || input.Name == nil || input.Capacity == nil {
		return ""
	}

	return fmt.Sprintf("%s_%d", string(*input.Name), *input.Capacity)
}

func expandPurviewIdentity(input []interface{}) (*identity.SystemAssignedIdentity, error) {
	expanded, err := purviewAccountIdentity{}.Expand(input)
	if err != nil {
		return nil, err
	}

	out := identity.SystemAssignedIdentity{}
	out.FromExpandedConfig(*expanded)
	return &out, nil
}

func flattenPurviewAccountIdentity(identity *identity.SystemAssignedIdentity) interface{} {
	expandedConfig := identity.ToExpandedConfig()
	return purviewAccountIdentity{}.Flatten(&expandedConfig)
}
