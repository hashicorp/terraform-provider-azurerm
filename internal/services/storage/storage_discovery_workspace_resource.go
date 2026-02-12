// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package storage

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storagediscovery/2025-09-01/storagediscoveryworkspaces"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceStorageDiscoveryWorkspace() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceStorageDiscoveryWorkspaceCreate,
		Read:   resourceStorageDiscoveryWorkspaceRead,
		Update: resourceStorageDiscoveryWorkspaceUpdate,
		Delete: resourceStorageDiscoveryWorkspaceDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := storagediscoveryworkspaces.ParseProviderStorageDiscoveryWorkspaceID(id)
			return err
		}),

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"resource_group_name": commonschema.ResourceGroupName(),

			"location": commonschema.Location(),

			"workspace_roots": {
				Type:     pluginsdk.TypeList,
				Required: true,
				MinItems: 1,
				Elem: &pluginsdk.Schema{
					Type:         pluginsdk.TypeString,
					ValidateFunc: validation.StringIsNotEmpty,
				},
			},

			"scopes": {
				Type:     pluginsdk.TypeList,
				Required: true,
				MinItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"display_name": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},

						"resource_types": {
							Type:     pluginsdk.TypeList,
							Required: true,
							MinItems: 1,
							Elem: &pluginsdk.Schema{
								Type:         pluginsdk.TypeString,
								ValidateFunc: validation.StringIsNotEmpty,
							},
						},

						"tag_keys_only": {
							Type:     pluginsdk.TypeList,
							Optional: true,
							Elem: &pluginsdk.Schema{
								Type:         pluginsdk.TypeString,
								ValidateFunc: validation.StringIsNotEmpty,
							},
						},

						"tags": {
							Type:     pluginsdk.TypeMap,
							Optional: true,
							Elem: &pluginsdk.Schema{
								Type: pluginsdk.TypeString,
							},
						},
					},
				},
			},

			"description": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"sku": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				Default:  string(storagediscoveryworkspaces.StorageDiscoverySkuStandard),
				ValidateFunc: validation.StringInSlice(
					storagediscoveryworkspaces.PossibleValuesForStorageDiscoverySku(),
					false,
				),
			},

			"tags": commonschema.Tags(),
		},
	}
}

func resourceStorageDiscoveryWorkspaceCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Storage.StorageDiscoveryWorkspacesClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := storagediscoveryworkspaces.NewProviderStorageDiscoveryWorkspaceID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	existing, err := client.Get(ctx, id)
	if err != nil {
		if !response.WasNotFound(existing.HttpResponse) {
			return fmt.Errorf("checking for existing %s: %+v", id, err)
		}
	}

	if !response.WasNotFound(existing.HttpResponse) {
		return tf.ImportAsExistsError("azurerm_storage_discovery_workspace", id.ID())
	}

	payload := storagediscoveryworkspaces.StorageDiscoveryWorkspace{
		Location: location.Normalize(d.Get("location").(string)),
		Properties: &storagediscoveryworkspaces.StorageDiscoveryWorkspaceProperties{
			WorkspaceRoots: expandStorageDiscoveryWorkspaceRoots(d.Get("workspace_roots").([]interface{})),
			Scopes:         expandStorageDiscoveryScopes(d.Get("scopes").([]interface{})),
		},
		Tags: tags.Expand(d.Get("tags").(map[string]interface{})),
	}

	if v, ok := d.GetOk("description"); ok {
		payload.Properties.Description = pointer.To(v.(string))
	}

	if v, ok := d.GetOk("sku"); ok {
		sku := storagediscoveryworkspaces.StorageDiscoverySku(v.(string))
		payload.Properties.Sku = &sku
	}

	if _, err := client.CreateOrUpdate(ctx, id, payload); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceStorageDiscoveryWorkspaceRead(d, meta)
}

func resourceStorageDiscoveryWorkspaceRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Storage.StorageDiscoveryWorkspacesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := storagediscoveryworkspaces.ParseProviderStorageDiscoveryWorkspaceID(d.Id())
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

	d.Set("name", id.StorageDiscoveryWorkspaceName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		d.Set("location", location.Normalize(model.Location))

		if props := model.Properties; props != nil {
			d.Set("description", props.Description)

			sku := string(storagediscoveryworkspaces.StorageDiscoverySkuStandard)
			if props.Sku != nil {
				sku = string(*props.Sku)
			}
			d.Set("sku", sku)
			d.Set("workspace_roots", flattenStorageDiscoveryWorkspaceRoots(props.WorkspaceRoots))

			if err := d.Set("scopes", flattenStorageDiscoveryScopes(props.Scopes)); err != nil {
				return fmt.Errorf("setting `scopes`: %+v", err)
			}
		}

		if err := tags.FlattenAndSet(d, model.Tags); err != nil {
			return err
		}
	}

	return nil
}

func resourceStorageDiscoveryWorkspaceUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Storage.StorageDiscoveryWorkspacesClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := storagediscoveryworkspaces.ParseProviderStorageDiscoveryWorkspaceID(d.Id())
	if err != nil {
		return err
	}

	payload := storagediscoveryworkspaces.StorageDiscoveryWorkspaceUpdate{
		Properties: &storagediscoveryworkspaces.StorageDiscoveryWorkspacePropertiesUpdate{},
	}

	if d.HasChange("description") {
		payload.Properties.Description = pointer.To(d.Get("description").(string))
	}

	if d.HasChange("sku") {
		sku := storagediscoveryworkspaces.StorageDiscoverySku(d.Get("sku").(string))
		payload.Properties.Sku = &sku
	}

	if d.HasChange("workspace_roots") {
		workspaceRoots := expandStorageDiscoveryWorkspaceRoots(d.Get("workspace_roots").([]interface{}))
		payload.Properties.WorkspaceRoots = &workspaceRoots
	}

	if d.HasChange("scopes") {
		scopes := expandStorageDiscoveryScopes(d.Get("scopes").([]interface{}))
		payload.Properties.Scopes = &scopes
	}

	if d.HasChange("tags") {
		payload.Tags = tags.Expand(d.Get("tags").(map[string]interface{}))
	}

	if _, err := client.Update(ctx, *id, payload); err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}

	return resourceStorageDiscoveryWorkspaceRead(d, meta)
}

func resourceStorageDiscoveryWorkspaceDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Storage.StorageDiscoveryWorkspacesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := storagediscoveryworkspaces.ParseProviderStorageDiscoveryWorkspaceID(d.Id())
	if err != nil {
		return err
	}

	if _, err := client.Delete(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", id, err)
	}

	return nil
}

func expandStorageDiscoveryWorkspaceRoots(input []interface{}) []string {
	result := make([]string, 0)
	for _, item := range input {
		if item != nil {
			result = append(result, item.(string))
		}
	}
	return result
}

func expandStorageDiscoveryScopes(input []interface{}) []storagediscoveryworkspaces.StorageDiscoveryScope {
	result := make([]storagediscoveryworkspaces.StorageDiscoveryScope, 0)

	for _, item := range input {
		v := item.(map[string]interface{})
		scope := storagediscoveryworkspaces.StorageDiscoveryScope{
			DisplayName:   v["display_name"].(string),
			ResourceTypes: expandStorageDiscoveryResourceTypes(v["resource_types"].([]interface{})),
		}

		if tagKeysOnly := v["tag_keys_only"].([]interface{}); len(tagKeysOnly) > 0 {
			scope.TagKeysOnly = expandStorageDiscoveryTagKeysOnly(tagKeysOnly)
		}

		if scopeTags := v["tags"].(map[string]interface{}); len(scopeTags) > 0 {
			tagsMap := make(map[string]string)
			for k, val := range scopeTags {
				tagsMap[k] = val.(string)
			}
			scope.Tags = &tagsMap
		}

		result = append(result, scope)
	}

	return result
}

func expandStorageDiscoveryResourceTypes(input []interface{}) []storagediscoveryworkspaces.StorageDiscoveryResourceType {
	result := make([]storagediscoveryworkspaces.StorageDiscoveryResourceType, 0)
	for _, item := range input {
		if item != nil {
			result = append(result, storagediscoveryworkspaces.StorageDiscoveryResourceType(item.(string)))
		}
	}
	return result
}

func expandStorageDiscoveryTagKeysOnly(input []interface{}) *[]string {
	result := make([]string, 0)
	for _, item := range input {
		if item != nil {
			result = append(result, item.(string))
		}
	}
	return &result
}

func flattenStorageDiscoveryWorkspaceRoots(input []string) []interface{} {
	result := make([]interface{}, 0)
	for _, item := range input {
		result = append(result, item)
	}
	return result
}

func flattenStorageDiscoveryScopes(input []storagediscoveryworkspaces.StorageDiscoveryScope) []interface{} {
	result := make([]interface{}, 0)

	for _, scope := range input {
		scopeMap := map[string]interface{}{
			"display_name":   scope.DisplayName,
			"resource_types": flattenStorageDiscoveryResourceTypes(scope.ResourceTypes),
		}

		if scope.TagKeysOnly != nil {
			scopeMap["tag_keys_only"] = flattenStorageDiscoveryTagKeysOnly(scope.TagKeysOnly)
		}

		if scope.Tags != nil {
			scopeMap["tags"] = pointer.From(scope.Tags)
		}

		result = append(result, scopeMap)
	}

	return result
}

func flattenStorageDiscoveryResourceTypes(input []storagediscoveryworkspaces.StorageDiscoveryResourceType) []interface{} {
	result := make([]interface{}, 0)
	for _, item := range input {
		result = append(result, string(item))
	}
	return result
}

func flattenStorageDiscoveryTagKeysOnly(input *[]string) []interface{} {
	result := make([]interface{}, 0)
	if input != nil {
		for _, item := range *input {
			result = append(result, item)
		}
	}
	return result
}
