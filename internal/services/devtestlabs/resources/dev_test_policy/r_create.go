// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package dev_test_policy

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/devtestlab/2018-09-15/policies"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/devtestlabs/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceArmDevTestPolicyCreateUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).DevTestLabs.PoliciesClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := policies.NewPolicyID(subscriptionId, d.Get("resource_group_name").(string), d.Get("lab_name").(string), d.Get("policy_set_name").(string), d.Get("name").(string))

	if d.IsNewResource() {
		if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
			existing, err := client.Get(ctx, id, policies.GetOperationOptions{})
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %s", id, err)
				}
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return tf.ImportAsExistsError("azurerm_dev_test_policy", id.ID())
			}
		}
	}

	factData := d.Get("fact_data").(string)
	threshold := d.Get("threshold").(string)
	evaluatorType := policies.PolicyEvaluatorType(d.Get("evaluator_type").(string))

	description := d.Get("description").(string)

	parameters := policies.Policy{
		Tags: helpers.ExpandTags(d.Get("tags").(map[string]any)),
		Properties: policies.PolicyProperties{
			FactName:      pointer.ToEnum[policies.PolicyFactName](id.PolicyName),
			FactData:      pointer.To(factData),
			Description:   pointer.To(description),
			EvaluatorType: &evaluatorType,
			Threshold:     pointer.To(threshold),
		},
	}

	if _, err := client.CreateOrUpdate(ctx, id, parameters); err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceArmDevTestPolicyRead(d, meta)
}
