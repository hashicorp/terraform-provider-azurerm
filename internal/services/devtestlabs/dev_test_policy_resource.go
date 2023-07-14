// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package devtestlabs

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/devtestlab/2018-09-15/policies"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/devtestlabs/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/devtestlabs/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceArmDevTestPolicy() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceArmDevTestPolicyCreateUpdate,
		Read:   resourceArmDevTestPolicyRead,
		Update: resourceArmDevTestPolicyCreateUpdate,
		Delete: resourceArmDevTestPolicyDelete,
		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := policies.ParsePolicyID(id)
			return err
		}),

		SchemaVersion: 1,
		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.DevTestLabPolicyUpgradeV0ToV1{},
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(policies.PolicyFactNameGalleryImage),
					string(policies.PolicyFactNameLabPremiumVMCount),
					string(policies.PolicyFactNameLabTargetCost),
					string(policies.PolicyFactNameLabVMCount),
					string(policies.PolicyFactNameLabVMSize),
					string(policies.PolicyFactNameUserOwnedLabPremiumVMCount),
					string(policies.PolicyFactNameUserOwnedLabVMCount),
					string(policies.PolicyFactNameUserOwnedLabVMCountInSubnet),
				}, false),
			},

			"policy_set_name": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
			},

			"lab_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.DevTestLabName(),
			},

			// There's a bug in the Azure API where this is returned in lower-case
			// BUG: https://github.com/Azure/azure-rest-api-specs/issues/3964
			"resource_group_name": azure.SchemaResourceGroupNameDiffSuppress(),

			"threshold": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"evaluator_type": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(policies.PolicyEvaluatorTypeAllowedValuesPolicy),
					string(policies.PolicyEvaluatorTypeMaxValuePolicy),
				}, false),
			},

			"description": {
				Type:     pluginsdk.TypeString,
				Optional: true,
			},

			"fact_data": {
				Type:     pluginsdk.TypeString,
				Optional: true,
			},

			"tags": tags.Schema(),
		},
	}
}

func resourceArmDevTestPolicyCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DevTestLabs.PoliciesClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for DevTest Policy creation")

	id := policies.NewPolicyID(subscriptionId, d.Get("resource_group_name").(string), d.Get("lab_name").(string), d.Get("policy_set_name").(string), d.Get("name").(string))

	if d.IsNewResource() {
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

	factData := d.Get("fact_data").(string)
	threshold := d.Get("threshold").(string)
	evaluatorType := policies.PolicyEvaluatorType(d.Get("evaluator_type").(string))

	description := d.Get("description").(string)
	factName := policies.PolicyFactName(id.PolicyName)

	parameters := policies.Policy{
		Tags: expandTags(d.Get("tags").(map[string]interface{})),
		Properties: policies.PolicyProperties{
			FactName:      &factName,
			FactData:      utils.String(factData),
			Description:   utils.String(description),
			EvaluatorType: &evaluatorType,
			Threshold:     utils.String(threshold),
		},
	}

	if _, err := client.CreateOrUpdate(ctx, id, parameters); err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceArmDevTestPolicyRead(d, meta)
}

func resourceArmDevTestPolicyRead(d *pluginsdk.ResourceData, meta interface{}) error {
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
		d.Set("evaluator_type", string(pointer.From(props.EvaluatorType)))
		d.Set("threshold", props.Threshold)

		if err = tags.FlattenAndSet(d, flattenTags(model.Tags)); err != nil {
			return err
		}
	}
	return nil
}

func resourceArmDevTestPolicyDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DevTestLabs.PoliciesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := policies.ParsePolicyID(d.Id())
	if err != nil {
		return err
	}

	read, err := client.Get(ctx, *id, policies.GetOperationOptions{})
	if err != nil {
		if response.WasNotFound(read.HttpResponse) {
			// deleted outside of TF
			log.Printf("[DEBUG] %s was not found  - assuming removed!", *id)
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	if _, err = client.Delete(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return err
}
