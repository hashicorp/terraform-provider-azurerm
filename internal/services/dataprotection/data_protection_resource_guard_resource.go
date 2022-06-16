package dataprotection

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/dataprotection/sdk/2022-04-01/resourceguards"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceDataProtectionResourceGuard() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceDataProtectionResourceGuardCreate,
		Read:   resourceDataProtectionResourceGuardRead,
		Update: resourceDataProtectionResourceGuardUpdate,
		Delete: resourceDataProtectionResourceGuardDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := resourceguards.ParseResourceGuardID(id)
			return err
		}),

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
			},

			"resource_group_name": commonschema.ResourceGroupName(),

			"location": commonschema.Location(),

			"identity": commonschema.SystemAssignedIdentityOptional(),

			"vault_critical_operation_exclusion_list": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				ForceNew: true,
				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
				},
			},

			"tags": commonschema.Tags(),
		},
	}
}

func resourceDataProtectionResourceGuardCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	client := meta.(*clients.Client).DataProtection.ResourceGuardClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := resourceguards.NewResourceGuardID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	existing, err := client.Get(ctx, id)
	if err != nil {
		if !response.WasNotFound(existing.HttpResponse) {
			return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
		}
	}
	if !response.WasNotFound(existing.HttpResponse) {
		return tf.ImportAsExistsError("azurerm_data_protection_resource_guard", id.ID())
	}

	expandedIdentity, err := expandResourceGuardIdentity(d.Get("identity").([]interface{}))
	if err != nil {
		return fmt.Errorf("expanding `identity`: %+v", err)
	}

	parameters := resourceguards.ResourceGuardResource{
		Identity: expandedIdentity,
		Location: utils.String(location.Normalize(d.Get("location").(string))),
		Properties: &resourceguards.ResourceGuard{
			VaultCriticalOperationExclusionList: utils.ExpandStringSlice(d.Get("vault_critical_operation_exclusion_list").([]interface{})),
		},
		Tags: tags.Expand(d.Get("tags").(map[string]interface{})),
	}

	if _, err := client.Put(ctx, id, parameters); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceDataProtectionResourceGuardRead(d, meta)
}

func resourceDataProtectionResourceGuardRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DataProtection.ResourceGuardClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := resourceguards.ParseResourceGuardID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[INFO] %s does not exist - removing from state", *id)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.ResourceGuardsName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		d.Set("location", location.Normalize(*model.Location))

		if v := model.Identity; v != nil {
			if err := d.Set("identity", flattenResourceGuardIdentity(v)); err != nil {
				return fmt.Errorf("setting `identity`: %+v", err)
			}
		}

		props := model.Properties
		d.Set("vault_critical_operation_exclusion_list", utils.FlattenStringSlice(props.VaultCriticalOperationExclusionList))

		if err := tags.FlattenAndSet(d, model.Tags); err != nil {
			return err
		}
	}

	return nil
}

func resourceDataProtectionResourceGuardUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DataProtection.ResourceGuardClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := resourceguards.ParseResourceGuardID(d.Id())
	if err != nil {
		return err
	}

	parameters := resourceguards.PatchResourceRequestInput{
		Properties: &resourceguards.PatchBackupVaultInput{},
	}

	if d.HasChange("identity") {
		expandedIdentity, err := expandResourceGuardIdentity(d.Get("identity").([]interface{}))
		if err != nil {
			return fmt.Errorf("expanding `identity`: %+v", err)
		}
		parameters.Identity = expandedIdentity
	}

	// https://github.com/Azure/azure-rest-api-specs/issues/19453
	// As there is a bug in rest api, so API doesn't support to update `vault_critical_operation_exclusion_list` for now
	if d.HasChange("vault_critical_operation_exclusion_list") {
		return fmt.Errorf("the service API doesn't support to update `vault_critical_operation_exclusion_list` for now")
	}

	if d.HasChange("tags") {
		parameters.Tags = tags.Expand(d.Get("tags").(map[string]interface{}))
	}

	if _, err := client.Patch(ctx, *id, parameters); err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}

	return resourceDataProtectionResourceGuardRead(d, meta)
}

func resourceDataProtectionResourceGuardDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DataProtection.ResourceGuardClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := resourceguards.ParseResourceGuardID(d.Id())
	if err != nil {
		return err
	}

	if _, err := client.Delete(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", id, err)
	}

	return nil
}

func expandResourceGuardIdentity(input []interface{}) (*resourceguards.DppIdentityDetails, error) {
	expanded, err := identity.ExpandSystemAssigned(input)
	if err != nil {
		return nil, err
	}

	return &resourceguards.DppIdentityDetails{
		Type: utils.String(string(expanded.Type)),
	}, nil
}

func flattenResourceGuardIdentity(input *resourceguards.DppIdentityDetails) []interface{} {
	var transform *identity.SystemAssigned

	if input != nil {
		transform = &identity.SystemAssigned{
			Type: identity.Type(*input.Type),
		}
		if input.PrincipalId != nil {
			transform.PrincipalId = *input.PrincipalId
		}
		if input.TenantId != nil {
			transform.TenantId = *input.TenantId
		}
	}

	return identity.FlattenSystemAssigned(transform)
}
