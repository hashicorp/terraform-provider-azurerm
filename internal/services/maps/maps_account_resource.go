// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package maps

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/maps/2023-06-01/accounts"
	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/maps/custompollers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/maps/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceMapsAccount() *pluginsdk.Resource {
	resource := &pluginsdk.Resource{
		Create: resourceMapsAccountCreate,
		Read:   resourceMapsAccountRead,
		Update: resourceMapsAccountUpdate,
		Delete: resourceMapsAccountDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := accounts.ParseAccountID(id)
			return err
		}),

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.AccountName(),
			},

			"resource_group_name": commonschema.ResourceGroupName(),

			"location": commonschema.Location(),

			"sku_name": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(accounts.NameSZero),
					string(accounts.NameSOne),
					string(accounts.NameGTwo),
				}, false),
			},

			"cors": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"allowed_origins": {
							Type:     pluginsdk.TypeList,
							Required: true,
							MinItems: 1,
							Elem: &pluginsdk.Schema{
								Type: pluginsdk.TypeString,
							},
						},
					},
				},
			},

			"identity": commonschema.SystemAssignedUserAssignedIdentityOptional(),

			"local_authentication_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  true,
			},

			"data_store": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"unique_name": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
						"storage_account_id": commonschema.ResourceIDReferenceOptional(&commonids.StorageAccountId{}),
					},
				},
			},

			"tags": commonschema.Tags(),

			"x_ms_client_id": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"primary_access_key": {
				Type:      pluginsdk.TypeString,
				Computed:  true,
				Sensitive: true,
			},

			"secondary_access_key": {
				Type:      pluginsdk.TypeString,
				Computed:  true,
				Sensitive: true,
			},
		},
	}

	if !features.FourPointOhBeta() {
		resource.Schema["location"] = &pluginsdk.Schema{
			Type:             schema.TypeString,
			Optional:         true,
			Computed:         true,
			ForceNew:         true,
			StateFunc:        location.StateFunc,
			DiffSuppressFunc: location.DiffSuppressFunc,
		}
	}
	return resource
}

func resourceMapsAccountCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Maps.AccountsClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for AzureRM Maps Account creation.")

	id := accounts.NewAccountID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	existing, err := client.Get(ctx, id)
	if err != nil {
		if !response.WasNotFound(existing.HttpResponse) {
			return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
		}
	}

	if !response.WasNotFound(existing.HttpResponse) {
		return tf.ImportAsExistsError("azurerm_maps_account", id.ID())
	}

	dataStores, err := expandDataStore(d.Get("data_store").([]interface{}))
	if err != nil {
		return fmt.Errorf("expanding `data_store`: %+v", err)
	}

	loc := "global"
	if v, ok := d.GetOk("location"); ok {
		loc = location.Normalize(v.(string))
	}

	parameters := accounts.MapsAccount{
		Location: loc,
		Sku: accounts.Sku{
			Name: accounts.Name(d.Get("sku_name").(string)),
		},
		Tags: tags.Expand(d.Get("tags").(map[string]interface{})),
		Properties: &accounts.MapsAccountProperties{
			DisableLocalAuth: pointer.To(!d.Get("local_authentication_enabled").(bool)),
			Cors:             expandCors(d.Get("cors").([]interface{})),
			LinkedResources:  dataStores,
		},
	}

	// setting anything into identity returns a 400 Bad Request error if the location of the maps account is `global` which is
	// what we were defaulting to previously - when `location` becomes Required in 4.0 we can remove this check and set
	// identity in the payload like we do elsewhere
	if v, ok := d.GetOk("identity"); ok {
		identityExpanded, err := identity.ExpandSystemAndUserAssignedMap(v.([]interface{}))
		if err != nil {
			return fmt.Errorf("expanding `identity`: %+v", err)
		}
		parameters.Identity = identityExpanded
	}

	if _, err := client.CreateOrUpdate(ctx, id, parameters); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	// These should actually be LROs, but they're not, custom poller is required until https://github.com/Azure/azure-rest-api-specs/issues/29501 is resolved
	pollerType := custompollers.NewMapsAccountPoller(client, id)
	poller := pollers.NewPoller(pollerType, 10*time.Second, pollers.DefaultNumberOfDroppedConnectionsToAllow)
	if err := poller.PollUntilDone(ctx); err != nil {
		return err
	}

	d.SetId(id.ID())

	return resourceMapsAccountRead(d, meta)
}

func resourceMapsAccountUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Maps.AccountsClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for AzureRM Maps Account creation.")

	id, err := accounts.ParseAccountID(d.Id())
	if err != nil {
		return err
	}

	existing, err := client.Get(ctx, *id)
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	if existing.Model == nil {
		return fmt.Errorf("retrieving %s: `model` was nil", id)
	}
	if existing.Model.Properties == nil {
		return fmt.Errorf("retrieving %s: `properties` was nil", id)
	}

	payload := existing.Model

	if d.HasChange("local_authentication_enabled") {
		payload.Properties.DisableLocalAuth = pointer.To(!d.Get("local_authentication_enabled").(bool))
	}

	if d.HasChange("cors") {
		payload.Properties.Cors = expandCors(d.Get("cors").([]interface{}))
	}

	if d.HasChange("data_store") {
		dataStores, err := expandDataStore(d.Get("data_store").([]interface{}))
		if err != nil {
			return fmt.Errorf("expanding `data_store`: %+v", err)
		}
		payload.Properties.LinkedResources = dataStores
	}

	if d.HasChange("tags") {
		payload.Tags = tags.Expand(d.Get("tags").(map[string]interface{}))
	}

	if _, err := client.CreateOrUpdate(ctx, *id, *payload); err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}

	// These should actually be LROs, but they're not, custom poller is required until https://github.com/Azure/azure-rest-api-specs/issues/29501 is resolved
	pollerType := custompollers.NewMapsAccountPoller(client, *id)
	poller := pollers.NewPoller(pollerType, 10*time.Second, pollers.DefaultNumberOfDroppedConnectionsToAllow)
	if err := poller.PollUntilDone(ctx); err != nil {
		return err
	}

	return resourceMapsAccountRead(d, meta)
}

func resourceMapsAccountRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Maps.AccountsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := accounts.ParseAccountID(d.Id())
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
		d.Set("location", location.Normalize(model.Location))
		d.Set("sku_name", string(model.Sku.Name))

		identityFlattened, err := identity.FlattenSystemAndUserAssignedMap(model.Identity)
		if err != nil {
			return fmt.Errorf("flattening `identity`: %+v", err)
		}

		d.Set("identity", identityFlattened)

		if props := model.Properties; props != nil {
			d.Set("x_ms_client_id", props.UniqueId)
			d.Set("cors", flattenCors(props.Cors))

			dataStore, err := flattenDataStore(props.LinkedResources)
			if err != nil {
				return fmt.Errorf("flattening `data_store`: %+v", err)
			}
			d.Set("data_store", dataStore)

			localAuthenticationEnabled := true
			if props.DisableLocalAuth != nil {
				localAuthenticationEnabled = !*props.DisableLocalAuth
			}
			d.Set("local_authentication_enabled", localAuthenticationEnabled)
		}

		if err := tags.FlattenAndSet(d, model.Tags); err != nil {
			return err
		}
	}

	keysResp, err := client.ListKeys(ctx, *id)
	if err != nil {
		return fmt.Errorf("retrieving Access Keys for %s: %+v", *id, err)
	}
	if model := keysResp.Model; model != nil {
		d.Set("primary_access_key", model.PrimaryKey)
		d.Set("secondary_access_key", model.SecondaryKey)
	}

	return nil
}

func resourceMapsAccountDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Maps.AccountsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := accounts.ParseAccountID(d.Id())
	if err != nil {
		return err
	}

	if _, err := client.Delete(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return nil
}

func expandCors(input []interface{}) *accounts.CorsRules {
	if len(input) == 0 {
		return nil
	}

	cors := input[0].(map[string]interface{})

	corsRule := make([]accounts.CorsRule, 0)

	corsRule = append(corsRule, accounts.CorsRule{
		AllowedOrigins: pointer.From(utils.ExpandStringSlice(cors["allowed_origins"].([]interface{}))),
	})

	return &accounts.CorsRules{
		CorsRules: &corsRule,
	}
}

func expandDataStore(input []interface{}) (*[]accounts.LinkedResource, error) {
	if len(input) == 0 {
		return nil, nil
	}

	linkedResources := make([]accounts.LinkedResource, 0)

	for _, i := range input {
		dataStore := i.(map[string]interface{})

		storageAccountId, err := commonids.ParseStorageAccountID(dataStore["storage_account_id"].(string))
		if err != nil {
			return nil, err
		}

		linkedResources = append(linkedResources, accounts.LinkedResource{
			Id:         storageAccountId.ID(),
			UniqueName: dataStore["unique_name"].(string),
		})
	}

	return &linkedResources, nil
}

func flattenCors(input *accounts.CorsRules) []interface{} {
	output := make([]interface{}, 0)

	if input == nil || input.CorsRules == nil || len(*input.CorsRules) == 0 {
		return output
	}

	// although this is a slice, only one element can be supplied/is present
	allowedOrigins := (*input.CorsRules)[0].AllowedOrigins

	output = append(output, map[string]interface{}{
		"allowed_origins": allowedOrigins,
	})

	return output
}

func flattenDataStore(input *[]accounts.LinkedResource) ([]interface{}, error) {
	output := make([]interface{}, 0)

	if input == nil || len(*input) == 0 {
		return output, nil
	}

	for _, resource := range *input {
		storageAccountId, err := commonids.ParseStorageAccountID(resource.Id)
		if err != nil {
			return nil, err
		}

		output = append(output, map[string]interface{}{
			"storage_account_id": storageAccountId.ID(),
			"unique_name":        resource.UniqueName,
		})
	}

	return output, nil
}
