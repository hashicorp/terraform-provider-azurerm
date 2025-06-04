// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package springcloud

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/springcloud/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/springcloud/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/springcloud/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
	"github.com/jackofallops/kermit/sdk/appplatform/2023-05-01-preview/appplatform"
)

func resourceSpringCloudBuildServiceBuilder() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		DeprecationMessage: features.DeprecatedInFivePointOh("Azure Spring Apps is now deprecated and will be retired on 2028-05-31 - as such the `azurerm_spring_cloud_builder` resource is deprecated and will be removed in a future major version of the AzureRM Provider. See https://aka.ms/asaretirement for more information."),

		Create: resourceSpringCloudBuildServiceBuilderCreateUpdate,
		Read:   resourceSpringCloudBuildServiceBuilderRead,
		Update: resourceSpringCloudBuildServiceBuilderCreateUpdate,
		Delete: resourceSpringCloudBuildServiceBuilderDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		SchemaVersion: 1,
		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.SpringCloudBuildServiceBuilderV0ToV1{},
		}),

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.SpringCloudBuildServiceBuilderID(id)
			return err
		}),

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
			},

			"spring_cloud_service_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.SpringCloudServiceID,
			},

			"build_pack_group": {
				Type:     pluginsdk.TypeSet,
				Required: true,
				MinItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"name": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},

						"build_pack_ids": {
							Type:     pluginsdk.TypeList,
							Optional: true,
							MinItems: 1,
							Elem: &pluginsdk.Schema{
								Type:         pluginsdk.TypeString,
								ValidateFunc: validation.StringIsNotEmpty,
							},
						},
					},
				},
			},

			"stack": {
				Type:     pluginsdk.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"id": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},

						"version": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
					},
				},
			},
		},
	}
}

func resourceSpringCloudBuildServiceBuilderCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	client := meta.(*clients.Client).AppPlatform.BuildServiceBuilderClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	springId, err := parse.SpringCloudServiceID(d.Get("spring_cloud_service_id").(string))
	if err != nil {
		return err
	}
	id := parse.NewSpringCloudBuildServiceBuilderID(subscriptionId, springId.ResourceGroup, springId.SpringName, "default", d.Get("name").(string))

	if d.IsNewResource() {
		existing, err := client.Get(ctx, id.ResourceGroup, id.SpringName, id.BuildServiceName, id.BuilderName)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for existing %s: %+v", id, err)
			}
		}
		if !utils.ResponseWasNotFound(existing.Response) {
			return tf.ImportAsExistsError("azurerm_spring_cloud_builder", id.ID())
		}
	}

	builderResource := appplatform.BuilderResource{
		Properties: &appplatform.BuilderProperties{
			BuildpackGroups: expandBuildServiceBuilderBuildPacksGroupPropertiesArray(d.Get("build_pack_group").(*pluginsdk.Set).List()),
			Stack:           expandBuildServiceBuilderStackProperties(d.Get("stack").([]interface{})),
		},
	}
	future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.SpringName, id.BuildServiceName, id.BuilderName, builderResource)
	if err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for creation/update of %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceSpringCloudBuildServiceBuilderRead(d, meta)
}

func resourceSpringCloudBuildServiceBuilderRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).AppPlatform.BuildServiceBuilderClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.SpringCloudBuildServiceBuilderID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.SpringName, id.BuildServiceName, id.BuilderName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] appplatform %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.Set("name", id.BuilderName)
	d.Set("spring_cloud_service_id", parse.NewSpringCloudServiceID(id.SubscriptionId, id.ResourceGroup, id.SpringName).ID())
	if props := resp.Properties; props != nil {
		if err := d.Set("build_pack_group", flattenBuildServiceBuilderBuildPacksGroupPropertiesArray(props.BuildpackGroups)); err != nil {
			return fmt.Errorf("setting `build_pack_group`: %+v", err)
		}
		if err := d.Set("stack", flattenBuildServiceBuilderStackProperties(props.Stack)); err != nil {
			return fmt.Errorf("setting `stack`: %+v", err)
		}
	}
	return nil
}

func resourceSpringCloudBuildServiceBuilderDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).AppPlatform.BuildServiceBuilderClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.SpringCloudBuildServiceBuilderID(d.Id())
	if err != nil {
		return err
	}

	future, err := client.Delete(ctx, id.ResourceGroup, id.SpringName, id.BuildServiceName, id.BuilderName)
	if err != nil {
		return fmt.Errorf("deleting %s: %+v", id, err)
	}

	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for deletion of %s: %+v", id, err)
	}
	return nil
}

func expandBuildServiceBuilderBuildPacksGroupPropertiesArray(input []interface{}) *[]appplatform.BuildpacksGroupProperties {
	results := make([]appplatform.BuildpacksGroupProperties, 0)
	for _, item := range input {
		v := item.(map[string]interface{})
		results = append(results, appplatform.BuildpacksGroupProperties{
			Name:       utils.String(v["name"].(string)),
			Buildpacks: expandBuildServiceBuilderBuildPackPropertiesArray(v["build_pack_ids"].([]interface{})),
		})
	}
	return &results
}

func expandBuildServiceBuilderStackProperties(input []interface{}) *appplatform.StackProperties {
	if len(input) == 0 {
		return nil
	}
	v := input[0].(map[string]interface{})
	return &appplatform.StackProperties{
		ID:      utils.String(v["id"].(string)),
		Version: utils.String(v["version"].(string)),
	}
}

func expandBuildServiceBuilderBuildPackPropertiesArray(input []interface{}) *[]appplatform.BuildpackProperties {
	results := make([]appplatform.BuildpackProperties, 0)
	for _, item := range input {
		results = append(results, appplatform.BuildpackProperties{
			ID: utils.String(item.(string)),
		})
	}
	return &results
}

func flattenBuildServiceBuilderBuildPacksGroupPropertiesArray(input *[]appplatform.BuildpacksGroupProperties) []interface{} {
	results := make([]interface{}, 0)
	if input == nil {
		return results
	}

	for _, item := range *input {
		var name string
		if item.Name != nil {
			name = *item.Name
		}
		results = append(results, map[string]interface{}{
			"name":           name,
			"build_pack_ids": flattenBuildServiceBuilderBuildPackPropertiesArray(item.Buildpacks),
		})
	}
	return results
}

func flattenBuildServiceBuilderStackProperties(input *appplatform.StackProperties) []interface{} {
	if input == nil {
		return make([]interface{}, 0)
	}

	var id string
	if input.ID != nil {
		id = *input.ID
	}
	var version string
	if input.Version != nil {
		version = *input.Version
	}
	return []interface{}{
		map[string]interface{}{
			"id":      id,
			"version": version,
		},
	}
}

func flattenBuildServiceBuilderBuildPackPropertiesArray(input *[]appplatform.BuildpackProperties) []interface{} {
	results := make([]interface{}, 0)
	if input == nil {
		return results
	}

	for _, item := range *input {
		var id string
		if item.ID != nil {
			id = *item.ID
		}
		results = append(results, id)
	}
	return results
}
