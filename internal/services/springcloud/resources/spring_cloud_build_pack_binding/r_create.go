// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package spring_cloud_build_pack_binding

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	appplatform_rm "github.com/hashicorp/go-azure-sdk/resource-manager/appplatform/2024-01-01-preview/appplatform"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	springcloudClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/springcloud/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/jackofallops/kermit/sdk/appplatform/2023-05-01-preview/appplatform"
)

func resourceSpringCloudBuildPackBindingCreateUpdate(d *pluginsdk.ResourceData, meta any) error {
	subscriptionId := account.From(meta).SubscriptionId
	client := springcloudClient.From(meta).BuildPackBindingClient
	ctx, cancel := timeouts.ForCreateUpdate(account.StopContext(meta), d)
	defer cancel()

	// the ID is parsed insensitively to preserve the behaviour of the legacy parser this replaced -
	// we are ok with this remaining insensitive as Azure Spring Apps is deprecated and will be removed
	builderId, err := appplatform_rm.ParseBuilderIDInsensitively(d.Get("spring_cloud_builder_id").(string))
	if err != nil {
		return err
	}
	id := appplatform_rm.NewBuildPackBindingID(subscriptionId, builderId.ResourceGroupName, builderId.SpringName, builderId.BuildServiceName, builderId.BuilderName, d.Get("name").(string))

	if d.IsNewResource() {
		if !features.From(meta).SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
			existing, err := client.Get(ctx, id.ResourceGroupName, id.SpringName, id.BuildServiceName, id.BuilderName, id.BuildPackBindingName)
			if err != nil {
				if !response.WasNotFound(existing.Response.Response) {
					return fmt.Errorf("checking for existing %s: %+v", id, err)
				}
			}
			if !response.WasNotFound(existing.Response.Response) {
				return tf.ImportAsExistsError("azurerm_spring_cloud_build_pack_binding", id.ID())
			}
		}
	}

	buildpackBinding := appplatform.BuildpackBindingResource{
		Properties: &appplatform.BuildpackBindingProperties{
			BindingType:      appplatform.BindingType(d.Get("binding_type").(string)),
			LaunchProperties: expandBuildPackBindingBuildPackBindingLaunchProperties(d.Get("launch").([]any)),
		},
	}
	future, err := client.CreateOrUpdate(ctx, id.ResourceGroupName, id.SpringName, id.BuildServiceName, id.BuilderName, id.BuildPackBindingName, buildpackBinding)
	if err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for creation/update of %s: %+v", id, err)
	}

	return resourceSpringCloudBuildPackBindingRead(d, meta)
}

func expandBuildPackBindingBuildPackBindingLaunchProperties(input []any) *appplatform.BuildpackBindingLaunchProperties {
	if len(input) == 0 {
		return nil
	}
	v := input[0].(map[string]any)
	var properties, secrets map[string]*string
	if valueRaw, ok := v["properties"]; ok && valueRaw != nil {
		properties = pluginsdk.ExpandMapStringPtrString(valueRaw.(map[string]any))
	}
	if valueRaw, ok := v["secrets"]; ok && valueRaw != nil {
		secrets = pluginsdk.ExpandMapStringPtrString(valueRaw.(map[string]any))
	}
	return &appplatform.BuildpackBindingLaunchProperties{
		Properties: properties,
		Secrets:    secrets,
	}
}
