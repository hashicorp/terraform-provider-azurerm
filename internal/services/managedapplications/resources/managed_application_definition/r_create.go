// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package managed_application_definition

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/managedapplications/2021-07-01/applicationdefinitions"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceManagedApplicationDefinitionCreate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).ManagedApplication.ApplicationDefinitionClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := applicationdefinitions.NewApplicationDefinitionID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		existing, err := client.Get(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("failed to check for presence of existing %s: %+v", id, err)
			}
		}
		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_managed_application_definition", id.ID())
		}
	}

	parameters := applicationdefinitions.ApplicationDefinition{
		Location: pointer.To(location.Normalize(d.Get("location").(string))),
		Properties: applicationdefinitions.ApplicationDefinitionProperties{
			Authorizations: expandManagedApplicationDefinitionAuthorization(d.Get("authorization").(*pluginsdk.Set).List()),
			Description:    pointer.To(d.Get("description").(string)),
			DisplayName:    pointer.To(d.Get("display_name").(string)),
			IsEnabled:      pointer.To(d.Get("package_enabled").(bool)),
			LockLevel:      applicationdefinitions.ApplicationLockLevel(d.Get("lock_level").(string)),
		},
		Tags: tags.Expand(d.Get("tags").(map[string]any)),
	}

	if v, ok := d.GetOk("create_ui_definition"); ok {
		parameters.Properties.CreateUiDefinition = &v
	}

	if v, ok := d.GetOk("main_template"); ok {
		parameters.Properties.MainTemplate = &v
	}

	if v, ok := d.GetOk("package_file_uri"); ok {
		parameters.Properties.PackageFileUri = pointer.To(v.(string))
	}

	if _, err := client.CreateOrUpdate(ctx, id, parameters); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceManagedApplicationDefinitionRead(d, meta)
}

func expandManagedApplicationDefinitionAuthorization(input []any) *[]applicationdefinitions.ApplicationAuthorization {
	results := make([]applicationdefinitions.ApplicationAuthorization, 0)
	for _, item := range input {
		v := item.(map[string]any)
		result := applicationdefinitions.ApplicationAuthorization{
			RoleDefinitionId: v["role_definition_id"].(string),
			PrincipalId:      v["service_principal_id"].(string),
		}

		results = append(results, result)
	}
	return &results
}
