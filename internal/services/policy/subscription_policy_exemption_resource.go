// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package policy

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/resources/mgmt/2021-06-01-preview/policy" // nolint: staticcheck
	"github.com/Azure/go-autorest/autorest/date"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/structure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	azValidate "github.com/hashicorp/terraform-provider-azurerm/helpers/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/policy/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/policy/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceArmSubscriptionPolicyExemption() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceArmSubscriptionPolicyExemptionCreateUpdate,
		Read:   resourceArmSubscriptionPolicyExemptionRead,
		Update: resourceArmSubscriptionPolicyExemptionCreateUpdate,
		Delete: resourceArmSubscriptionPolicyExemptionDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.SubscriptionPolicyExemptionID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"subscription_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: commonids.ValidateSubscriptionID,
			},

			"exemption_category": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(policy.ExemptionCategoryMitigated),
					string(policy.ExemptionCategoryWaiver),
				}, false),
			},

			"policy_assignment_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.PolicyAssignmentID,
			},

			"display_name": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(1, 128),
			},

			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(1, 512),
			},

			"policy_definition_reference_ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.StringIsNotEmpty,
				},
			},

			"expires_on": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: azValidate.ISO8601DateTime,
			},

			"metadata": metadataSchema(),
		},
	}
}

func resourceArmSubscriptionPolicyExemptionCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Policy.ExemptionsClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	subscriptionId, err := commonids.ParseSubscriptionID(d.Get("subscription_id").(string))
	if err != nil {
		return err
	}
	id := parse.NewSubscriptionPolicyExemptionID(subscriptionId.SubscriptionId, d.Get("name").(string))

	if d.IsNewResource() {
		existing, err := client.Get(ctx, subscriptionId.ID(), id.PolicyExemptionName)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id.ID(), err)
			}
		}
		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_resource_group_policy_exemption", *existing.ID)
		}
	}

	exemption := policy.Exemption{
		ExemptionProperties: &policy.ExemptionProperties{
			PolicyAssignmentID:           utils.String(d.Get("policy_assignment_id").(string)),
			PolicyDefinitionReferenceIds: utils.ExpandStringSlice(d.Get("policy_definition_reference_ids").([]interface{})),
			ExemptionCategory:            policy.ExemptionCategory(d.Get("exemption_category").(string)),
		},
	}

	if v, ok := d.GetOk("display_name"); ok {
		exemption.ExemptionProperties.DisplayName = utils.String(v.(string))
	}

	if v, ok := d.GetOk("description"); ok {
		exemption.ExemptionProperties.Description = utils.String(v.(string))
	}

	if v, ok := d.GetOk("expires_on"); ok {
		t, err := date.ParseTime(time.RFC3339, v.(string))
		if err != nil {
			return fmt.Errorf("expanding `expires_on`: %+v", err)
		}
		exemption.ExemptionProperties.ExpiresOn = &date.Time{Time: t}
	}

	if metaDataString := d.Get("metadata").(string); metaDataString != "" {
		metaData, err := structure.ExpandJsonFromString(metaDataString)
		if err != nil {
			return fmt.Errorf("unable to parse metadata: %+v", err)
		}
		exemption.ExemptionProperties.Metadata = &metaData
	}

	if _, err := client.CreateOrUpdate(ctx, subscriptionId.ID(), id.PolicyExemptionName, exemption); err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id.ID(), err)
	}

	d.SetId(id.ID())

	return resourceArmSubscriptionPolicyExemptionRead(d, meta)
}

func resourceArmSubscriptionPolicyExemptionRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Policy.ExemptionsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.SubscriptionPolicyExemptionID(d.Id())
	if err != nil {
		return fmt.Errorf("reading Policy Exemption: %+v", err)
	}

	subscriptionId := commonids.NewSubscriptionID(id.SubscriptionId)

	resp, err := client.Get(ctx, subscriptionId.ID(), id.PolicyExemptionName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] Error reading Policy Exemption %q - removing from state", d.Id())
			d.SetId("")
			return nil
		}

		return fmt.Errorf("reading %s: %+v", id.ID(), err)
	}

	d.Set("name", resp.Name)
	d.Set("subscription_id", subscriptionId.ID())
	if props := resp.ExemptionProperties; props != nil {
		d.Set("policy_assignment_id", props.PolicyAssignmentID)
		d.Set("display_name", props.DisplayName)
		d.Set("description", props.Description)
		d.Set("exemption_category", string(props.ExemptionCategory))

		if err := d.Set("policy_definition_reference_ids", utils.FlattenStringSlice(props.PolicyDefinitionReferenceIds)); err != nil {
			return fmt.Errorf("setting `policy_definition_reference_ids: %+v", err)
		}

		expiresOn := ""
		if expiresTime := props.ExpiresOn; expiresTime != nil {
			expiresOn = expiresTime.String()
		}
		d.Set("expires_on", expiresOn)

		if metadataStr := flattenJSON(props.Metadata); metadataStr != "" {
			d.Set("metadata", metadataStr)
		}
	}

	return nil
}

func resourceArmSubscriptionPolicyExemptionDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Policy.ExemptionsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.SubscriptionPolicyExemptionID(d.Id())
	if err != nil {
		return err
	}

	subscriptionId := commonids.NewSubscriptionID(id.SubscriptionId)

	if _, err := client.Delete(ctx, subscriptionId.ID(), id.PolicyExemptionName); err != nil {
		return fmt.Errorf("deleting %s: %+v", id.ID(), err)
	}

	return nil
}
