// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package managed_application

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/managedapplications/2021-07-01/applications"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceManagedApplicationCreate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).ManagedApplication.ApplicationClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := applications.NewApplicationID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		existing, err := client.Get(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("failed to check for presence of %s: %+v", id, err)
			}
		}
		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_managed_application", id.ID())
		}
	}

	parameters := applications.Application{
		Location: pointer.To(location.Normalize(d.Get("location").(string))),
		Kind:     d.Get("kind").(string),
		Tags:     tags.Expand(d.Get("tags").(map[string]any)),
	}

	if v, ok := d.GetOk("managed_resource_group_name"); ok {
		targetResourceGroupId := commonids.NewResourceGroupID(meta.(*clients.Client).Account.SubscriptionId, v.(string))
		parameters.Properties = applications.ApplicationProperties{
			ManagedResourceGroupId: pointer.To(targetResourceGroupId.ID()),
		}
	}

	if v, ok := d.GetOk("application_definition_id"); ok {
		parameters.Properties.ApplicationDefinitionId = pointer.To(v.(string))
	}

	if v, ok := d.GetOk("plan"); ok {
		parameters.Plan = expandManagedApplicationPlan(v.([]any))
	}

	if _, ok := d.GetOk("identity"); ok {
		managedApplicationIdentity, err := expandManagedApplicationIdentity(d.Get("identity").([]any))
		if err != nil {
			return fmt.Errorf("expanding `identity`: %+v", err)
		}
		parameters.Identity = managedApplicationIdentity
	}

	params, err := expandManagedApplicationParameters(d)
	if err != nil {
		return fmt.Errorf("expanding `parameter_values`: %+v", err)
	}
	parameters.Properties.Parameters = pointer.To(any(params))

	if err := client.CreateOrUpdateCallbackThenPoll(ctx, id, parameters, sdk.SetIDCallback(meta, &id, d)); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}
	d.SetId(id.ID())

	return resourceManagedApplicationRead(d, meta)
}
