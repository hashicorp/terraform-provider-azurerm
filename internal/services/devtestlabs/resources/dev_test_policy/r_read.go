// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package dev_test_policy

import (
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/devtestlab/2018-09-15/policies"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/devtestlabs/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceArmDevTestPolicyRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).DevTestLabs.PoliciesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := policies.ParsePolicyID(d.Id())
	if err != nil {
		return err
	}

	read, err := client.Get(ctx, *id, policies.GetOperationOptions{})
	if err != nil {
		if response.WasNotFound(read.HttpResponse) {
			log.Printf("[DEBUG] %s was not found - removing from state!", *id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("making Read request on %s: %+v", *id, err)
	}

	d.Set("name", id.PolicyName)
	d.Set("policy_set_name", id.PolicySetName)
	d.Set("lab_name", id.LabName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := read.Model; model != nil {
		props := model.Properties
		d.Set("description", props.Description)
		d.Set("fact_data", props.FactData)
		d.Set("evaluator_type", pointer.FromEnum(props.EvaluatorType))
		d.Set("threshold", props.Threshold)

		if err = tags.FlattenAndSet(d, helpers.FlattenTags(model.Tags)); err != nil {
			return err
		}
	}
	return nil
}
