// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package purview

import (
	"fmt"
	"regexp"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourcegroups"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/purview/2021-07-01/account"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourcePurviewAccount() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourcePurviewAccountCreate,
		Read:   resourcePurviewAccountRead,
		Update: resourcePurviewAccountUpdate,
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
				ValidateFunc: resourcegroups.ValidateName,
			},

			"identity": commonschema.SystemOrUserAssignedIdentityRequired(),

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

			"tags": commonschema.Tags(),
		},
	}
}

func resourcePurviewAccountCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Purview.AccountsClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := account.NewAccountID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	existing, err := client.Get(ctx, id)
	if err != nil {
		if !response.WasNotFound(existing.HttpResponse) {
			return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
		}
	}

	if !response.WasNotFound(existing.HttpResponse) {
		return tf.ImportAsExistsError("azurerm_purview_account", id.ID())
	}

	purviewAccount := account.Account{
		Properties: &account.AccountProperties{},
		Location:   pointer.To(location.Normalize(d.Get("location").(string))),
		Tags:       tags.Expand(d.Get("tags").(map[string]interface{})),
	}

	expandedIdentity, err := identity.ExpandSystemOrUserAssignedMap(d.Get("identity").([]interface{}))
	if err != nil {
		return fmt.Errorf("expanding `identity`: %+v", err)
	}
	purviewAccount.Identity = expandedIdentity

	publicNetworkAccessEnabled := account.PublicNetworkAccessDisabled
	if d.Get("public_network_enabled").(bool) {
		publicNetworkAccessEnabled = account.PublicNetworkAccessEnabled
	}
	purviewAccount.Properties.PublicNetworkAccess = &publicNetworkAccessEnabled

	if v, ok := d.GetOk("managed_resource_group_name"); ok {
		purviewAccount.Properties.ManagedResourceGroupName = pointer.To(v.(string))
	}

	if err := client.CreateOrUpdateThenPoll(ctx, id, purviewAccount); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
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

	d.Set("name", id.AccountName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		d.Set("location", location.NormalizeNilable(model.Location))

		flattenedIdentity, err := identity.FlattenSystemOrUserAssignedMap(model.Identity)
		if err != nil {
			return fmt.Errorf("flattening `identity`: %+v", err)
		}
		if err := d.Set("identity", flattenedIdentity); err != nil {
			return fmt.Errorf("flattening `identity`: %+v", err)
		}

		if props := model.Properties; props != nil {
			if err := d.Set("managed_resources", flattenPurviewAccountManagedResources(props.ManagedResources)); err != nil {
				return fmt.Errorf("flattening `managed_resources`: %+v", err)
			}

			d.Set("public_network_enabled", pointer.From(props.PublicNetworkAccess) == account.PublicNetworkAccessEnabled)

			d.Set("managed_resource_group_name", pointer.From(props.ManagedResourceGroupName))

			if endpoints := props.Endpoints; endpoints != nil {
				d.Set("catalog_endpoint", endpoints.Catalog)
				d.Set("guardian_endpoint", endpoints.Guardian)
				d.Set("scan_endpoint", endpoints.Scan)
			}
		}

		if err = tags.FlattenAndSet(d, model.Tags); err != nil {
			return err
		}
	}

	if keys, err := client.ListKeys(ctx, *id); err == nil {
		if model := keys.Model; model != nil {
			d.Set("atlas_kafka_endpoint_primary_connection_string", model.AtlasKafkaPrimaryEndpoint)
			d.Set("atlas_kafka_endpoint_secondary_connection_string", model.AtlasKafkaSecondaryEndpoint)
		}
	} else {
		// if eventhubs have been disabled we will get a response was not found, so we can ignore that error
		if !response.WasNotFound(keys.HttpResponse) {
			return fmt.Errorf("retrieving Keys for %s: %+v", *id, err)
		}

		d.Set("atlas_kafka_endpoint_primary_connection_string", "")
		d.Set("atlas_kafka_endpoint_secondary_connection_string", "")
	}

	return nil
}

func resourcePurviewAccountUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Purview.AccountsClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := account.ParseAccountID(d.Id())
	if err != nil {
		return err
	}

	parameters := account.AccountUpdateParameters{
		Properties: &account.AccountProperties{},
	}

	if d.HasChange("public_network_enabled") {
		if d.Get("public_network_enabled").(bool) {
			parameters.Properties.PublicNetworkAccess = pointer.To(account.PublicNetworkAccessEnabled)
		} else {
			parameters.Properties.PublicNetworkAccess = pointer.To(account.PublicNetworkAccessDisabled)
		}
	}

	if d.HasChange("tags") {
		parameters.Tags = tags.Expand(d.Get("tags").(map[string]interface{}))
	}

	if d.HasChange("identity") {
		expandedIdentity, err := identity.ExpandSystemOrUserAssignedMap(d.Get("identity").([]interface{}))
		if err != nil {
			return fmt.Errorf("expanding `identity`: %+v", err)
		}
		parameters.Identity = expandedIdentity
	}

	if err := client.UpdateThenPoll(ctx, *id, parameters); err != nil {
		return fmt.Errorf("updating %s: %+v", *id, err)
	}

	return resourcePurviewAccountRead(d, meta)
}

func resourcePurviewAccountDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Purview.AccountsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := account.ParseAccountID(d.Id())
	if err != nil {
		return err
	}

	err = client.DeleteThenPoll(ctx, *id)
	if err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return nil
}

func flattenPurviewAccountManagedResources(managedResources *account.ManagedResources) interface{} {
	if managedResources == nil {
		return make([]interface{}, 0)
	}

	return []interface{}{
		map[string]interface{}{
			"resource_group_id":      pointer.From(managedResources.ResourceGroup),
			"storage_account_id":     pointer.From(managedResources.StorageAccount),
			"event_hub_namespace_id": pointer.From(managedResources.EventHubNamespace),
		},
	}
}
