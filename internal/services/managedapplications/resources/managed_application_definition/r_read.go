// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package managed_application_definition

import (
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/managedapplications/2021-07-01/applicationdefinitions"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceManagedApplicationDefinitionRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).ManagedApplication.ApplicationDefinitionClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := applicationdefinitions.ParseApplicationDefinitionID(d.Id())
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
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.Set("name", id.ApplicationDefinitionName)
	d.Set("resource_group_name", id.ResourceGroupName) // missing from response?

	if model := resp.Model; model != nil {
		p := model.Properties

		d.Set("location", location.NormalizeNilable(model.Location))

		if err := d.Set("authorization", flattenManagedApplicationDefinitionAuthorization(p.Authorizations)); err != nil {
			return fmt.Errorf("setting `authorization`: %+v", err)
		}
		d.Set("description", p.Description)
		d.Set("display_name", p.DisplayName)
		d.Set("package_enabled", p.IsEnabled)
		d.Set("lock_level", string(p.LockLevel))

		// the following are not returned from the API so lets pull it from state
		if v, ok := d.GetOk("create_ui_definition"); ok {
			d.Set("create_ui_definition", v.(string))
		}
		if v, ok := d.GetOk("main_template"); ok {
			d.Set("main_template", v.(string))
		}
		if v, ok := d.GetOk("package_file_uri"); ok {
			d.Set("package_file_uri", v.(string))
		}

		if err := tags.FlattenAndSet(d, model.Tags); err != nil {
			return err
		}
	}

	return nil
}

func flattenManagedApplicationDefinitionAuthorization(input *[]applicationdefinitions.ApplicationAuthorization) []any {
	results := make([]any, 0)
	if input == nil {
		return results
	}

	for _, item := range *input {
		results = append(results, map[string]any{
			"role_definition_id":   item.RoleDefinitionId,
			"service_principal_id": item.PrincipalId,
		})
	}

	return results
}
