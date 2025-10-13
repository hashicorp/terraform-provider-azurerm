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

func resourceSpringCloudBuildPackBinding() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		DeprecationMessage: features.DeprecatedInFivePointOh("Azure Spring Apps is now deprecated and will be retired on 2028-05-31 - as such the `azurerm_spring_cloud_build_pack_binding` resource is deprecated and will be removed in a future major version of the AzureRM Provider. See https://aka.ms/asaretirement for more information."),

		Create: resourceSpringCloudBuildPackBindingCreateUpdate,
		Read:   resourceSpringCloudBuildPackBindingRead,
		Update: resourceSpringCloudBuildPackBindingCreateUpdate,
		Delete: resourceSpringCloudBuildPackBindingDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.SpringCloudBuildPackBindingID(id)
			return err
		}),

		SchemaVersion: 1,
		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.BuildPackBindingV0ToV1{},
		}),

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
			},

			"spring_cloud_builder_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.SpringCloudBuildServiceBuilderID,
			},

			"binding_type": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(appplatform.BindingTypeApplicationInsights),
					string(appplatform.BindingTypeApacheSkyWalking),
					string(appplatform.BindingTypeAppDynamics),
					string(appplatform.BindingTypeDynatrace),
					string(appplatform.BindingTypeNewRelic),
					string(appplatform.BindingTypeElasticAPM),
				}, false),
			},

			"launch": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"properties": {
							Type:     pluginsdk.TypeMap,
							Optional: true,
							Elem: &pluginsdk.Schema{
								Type:         pluginsdk.TypeString,
								ValidateFunc: validation.StringIsNotEmpty,
							},
						},

						"secrets": {
							Type:     pluginsdk.TypeMap,
							Optional: true,
							Elem: &pluginsdk.Schema{
								Type:         pluginsdk.TypeString,
								ValidateFunc: validation.StringIsNotEmpty,
							},
						},
					},
				},
			},
		},
	}
}

func resourceSpringCloudBuildPackBindingCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	client := meta.(*clients.Client).AppPlatform.BuildPackBindingClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	builderId, err := parse.SpringCloudBuildServiceBuilderID(d.Get("spring_cloud_builder_id").(string))
	if err != nil {
		return err
	}
	id := parse.NewSpringCloudBuildPackBindingID(subscriptionId, builderId.ResourceGroup, builderId.SpringName, builderId.BuildServiceName, builderId.BuilderName, d.Get("name").(string))

	if d.IsNewResource() {
		existing, err := client.Get(ctx, id.ResourceGroup, id.SpringName, id.BuildServiceName, id.BuilderName, id.BuildPackBindingName)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for existing %s: %+v", id, err)
			}
		}
		if !utils.ResponseWasNotFound(existing.Response) {
			return tf.ImportAsExistsError("azurerm_spring_cloud_build_pack_binding", id.ID())
		}
	}

	buildpackBinding := appplatform.BuildpackBindingResource{
		Properties: &appplatform.BuildpackBindingProperties{
			BindingType:      appplatform.BindingType(d.Get("binding_type").(string)),
			LaunchProperties: expandBuildPackBindingBuildPackBindingLaunchProperties(d.Get("launch").([]interface{})),
		},
	}
	future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.SpringName, id.BuildServiceName, id.BuilderName, id.BuildPackBindingName, buildpackBinding)
	if err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for creation/update of %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceSpringCloudBuildPackBindingRead(d, meta)
}

func resourceSpringCloudBuildPackBindingRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).AppPlatform.BuildPackBindingClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.SpringCloudBuildPackBindingID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.SpringName, id.BuildServiceName, id.BuilderName, id.BuildPackBindingName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] appplatform %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}
	d.Set("name", id.BuildPackBindingName)
	d.Set("spring_cloud_builder_id", parse.NewSpringCloudBuildServiceBuilderID(id.SubscriptionId, id.ResourceGroup, id.SpringName, id.BuildServiceName, id.BuilderName).ID())
	if props := resp.Properties; props != nil {
		d.Set("binding_type", props.BindingType)
		if err := d.Set("launch", flattenBuildPackBindingBuildPackBindingLaunchProperties(props.LaunchProperties, d.Get("launch").([]interface{}))); err != nil {
			return fmt.Errorf("setting `launch`: %+v", err)
		}
	}
	return nil
}

func resourceSpringCloudBuildPackBindingDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).AppPlatform.BuildPackBindingClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.SpringCloudBuildPackBindingID(d.Id())
	if err != nil {
		return err
	}

	future, err := client.Delete(ctx, id.ResourceGroup, id.SpringName, id.BuildServiceName, id.BuilderName, id.BuildPackBindingName)
	if err != nil {
		return fmt.Errorf("deleting %s: %+v", id, err)
	}

	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for deletion of %s: %+v", id, err)
	}
	return nil
}

func expandBuildPackBindingBuildPackBindingLaunchProperties(input []interface{}) *appplatform.BuildpackBindingLaunchProperties {
	if len(input) == 0 {
		return nil
	}
	v := input[0].(map[string]interface{})
	var properties, secrets map[string]*string
	if valueRaw, ok := v["properties"]; ok && valueRaw != nil {
		properties = utils.ExpandMapStringPtrString(valueRaw.(map[string]interface{}))
	}
	if valueRaw, ok := v["secrets"]; ok && valueRaw != nil {
		secrets = utils.ExpandMapStringPtrString(valueRaw.(map[string]interface{}))
	}
	return &appplatform.BuildpackBindingLaunchProperties{
		Properties: properties,
		Secrets:    secrets,
	}
}

func flattenBuildPackBindingBuildPackBindingLaunchProperties(input *appplatform.BuildpackBindingLaunchProperties, old []interface{}) []interface{} {
	if input == nil {
		return make([]interface{}, 0)
	}

	props := make(map[string]interface{})
	if input.Properties != nil {
		props = utils.FlattenMapStringPtrString(input.Properties)
	}
	secrets := make(map[string]interface{})
	if len(old) != 0 {
		v := old[0].(map[string]interface{})
		if secretsRaw, ok := v["secrets"]; ok && secretsRaw != nil {
			secrets = secretsRaw.(map[string]interface{})
		}
	}
	return []interface{}{
		map[string]interface{}{
			"properties": props,
			"secrets":    secrets,
		},
	}
}
