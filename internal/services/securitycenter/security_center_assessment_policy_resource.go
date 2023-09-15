// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package securitycenter

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/security/mgmt/v3.0/security" // nolint: staticcheck
	"github.com/google/uuid"
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/security/2021-06-01/assessmentsmetadata"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceArmSecurityCenterAssessmentPolicy() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceArmSecurityCenterAssessmentPolicyCreate,
		Read:   resourceArmSecurityCenterAssessmentPolicyRead,
		Update: resourceArmSecurityCenterAssessmentPolicyUpdate,
		Delete: resourceArmSecurityCenterAssessmentPolicyDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := assessmentsmetadata.ParseProviderAssessmentMetadataID(id)
			return err
		}),

		Schema: map[string]*pluginsdk.Schema{
			"description": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"display_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"severity": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				Default:  string(security.SeverityMedium),
				ValidateFunc: validation.StringInSlice([]string{
					string(security.SeverityLow),
					string(security.SeverityMedium),
					string(security.SeverityHigh),
				}, false),
			},

			// API would return `Unknown` when `categories` isn't set.
			// After synced with service team, they confirmed will add `Unknown` as possible value to this property and it will be published as a new version of this API.
			// https://github.com/Azure/azure-rest-api-specs/issues/14918
			"categories": {
				Type:     pluginsdk.TypeSet,
				Optional: true,
				Computed: true,
				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
					ValidateFunc: validation.StringInSlice([]string{
						"Unknown",
						string(security.Compute),
						string(security.Data),
						string(security.IdentityAndAccess),
						string(security.IoT),
						string(security.Networking),
					}, false),
				},
			},

			"implementation_effort": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(security.ImplementationEffortLow),
					string(security.ImplementationEffortModerate),
					string(security.ImplementationEffortHigh),
				}, false),
			},

			"remediation_description": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"threats": {
				Type:     pluginsdk.TypeSet,
				Optional: true,
				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
					ValidateFunc: validation.StringInSlice([]string{
						"AccountBreach",
						"DataExfiltration",
						"DataSpillage",
						"MaliciousInsider",
						"ElevationOfPrivilege",
						"ThreatResistance",
						"MissingCoverage",
						"DenialOfService",
					}, false),
				},
			},

			"user_impact": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(security.UserImpactLow),
					string(security.UserImpactModerate),
					string(security.UserImpactHigh),
				}, false),
			},

			"name": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceArmSecurityCenterAssessmentPolicyCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	client := meta.(*clients.Client).SecurityCenter.AssessmentsMetadataClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := uuid.New().String()

	id := assessmentsmetadata.NewProviderAssessmentMetadataID(subscriptionId, name)

	existing, err := client.GetInSubscription(ctx, id)
	if err != nil {
		if !response.WasNotFound(existing.HttpResponse) {
			return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
		}
	}

	if !response.WasNotFound(existing.HttpResponse) {
		return tf.ImportAsExistsError("azurerm_security_center_assessment_policy", id.ID())
	}

	params := assessmentsmetadata.SecurityAssessmentMetadataResponse{
		Properties: &assessmentsmetadata.SecurityAssessmentMetadataPropertiesResponse{
			AssessmentType: assessmentsmetadata.AssessmentTypeCustomerManaged,
			Description:    pointer.To(d.Get("description").(string)),
			DisplayName:    d.Get("display_name").(string),
			Severity:       assessmentsmetadata.Severity(d.Get("severity").(string)),
		},
	}

	if v, ok := d.GetOk("categories"); ok {
		categories := make([]assessmentsmetadata.Categories, 0)
		for _, item := range v.(*pluginsdk.Set).List() {
			categories = append(categories, (assessmentsmetadata.Categories)(item.(string)))
		}
		params.Properties.Categories = &categories
	}

	if v, ok := d.GetOk("threats"); ok {
		threats := make([]assessmentsmetadata.Threats, 0)
		for _, item := range v.(*pluginsdk.Set).List() {
			threats = append(threats, assessmentsmetadata.Threats(item.(string)))
		}
		params.Properties.Threats = &threats
	}

	if v, ok := d.GetOk("implementation_effort"); ok {
		params.Properties.ImplementationEffort = pointer.To(assessmentsmetadata.ImplementationEffort(v.(string)))
	}

	if v, ok := d.GetOk("remediation_description"); ok {
		params.Properties.RemediationDescription = pointer.To(v.(string))
	}

	if v, ok := d.GetOk("user_impact"); ok {
		params.Properties.UserImpact = pointer.To(assessmentsmetadata.UserImpact(v.(string)))
	}

	if _, err := client.CreateInSubscription(ctx, id, params); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceArmSecurityCenterAssessmentPolicyRead(d, meta)
}

func resourceArmSecurityCenterAssessmentPolicyRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).SecurityCenter.AssessmentsMetadataClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := assessmentsmetadata.ParseProviderAssessmentMetadataID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.GetInSubscription(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[INFO] %s does not exist - removing from state", *id)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.AssessmentMetadataName)

	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
			d.Set("description", pointer.From(props.Description))
			d.Set("display_name", props.DisplayName)
			d.Set("severity", string(props.Severity))
			d.Set("implementation_effort", string(pointer.From(props.ImplementationEffort)))
			d.Set("remediation_description", pointer.From(props.RemediationDescription))
			d.Set("user_impact", string(pointer.From(props.UserImpact)))

			categories := make([]string, 0)
			if props.Categories != nil {
				for _, item := range *props.Categories {
					categories = append(categories, string(item))
				}
			}
			d.Set("categories", utils.FlattenStringSlice(&categories))

			threats := make([]string, 0)
			if props.Threats != nil {
				for _, item := range *props.Threats {
					threats = append(threats, string(item))
				}
			}
			d.Set("threats", utils.FlattenStringSlice(&threats))
		}
	}

	return nil
}

func resourceArmSecurityCenterAssessmentPolicyUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).SecurityCenter.AssessmentsMetadataClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := assessmentsmetadata.ParseProviderAssessmentMetadataID(d.Id())
	if err != nil {
		return err
	}

	existing, err := client.GetInSubscription(ctx, *id)
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}
	if existing.Model == nil || existing.Model.Properties == nil {
		return fmt.Errorf("retrieving %s: `properties` was nil", *id)
	}

	if d.HasChange("description") {
		existing.Model.Properties.Description = pointer.To(d.Get("description").(string))
	}

	if d.HasChange("display_name") {
		existing.Model.Properties.DisplayName = d.Get("display_name").(string)
	}

	if d.HasChange("severity") {
		existing.Model.Properties.Severity = assessmentsmetadata.Severity(d.Get("severity").(string))
	}

	if d.HasChange("categories") {
		categories := make([]assessmentsmetadata.Categories, 0)
		for _, item := range d.Get("categories").(*pluginsdk.Set).List() {
			categories = append(categories, assessmentsmetadata.Categories(item.(string)))
		}
		existing.Model.Properties.Categories = &categories
	}

	if d.HasChange("threats") {
		threats := make([]assessmentsmetadata.Threats, 0)
		for _, item := range d.Get("threats").(*pluginsdk.Set).List() {
			threats = append(threats, (assessmentsmetadata.Threats)(item.(string)))
		}
		existing.Model.Properties.Threats = &threats
	}

	if d.HasChange("implementation_effort") {
		existing.Model.Properties.ImplementationEffort = pointer.To(assessmentsmetadata.ImplementationEffort(d.Get("implementation_effort").(string)))
	}

	if d.HasChange("remediation_description") {
		existing.Model.Properties.RemediationDescription = utils.String(d.Get("remediation_description").(string))
	}

	if d.HasChange("user_impact") {
		existing.Model.Properties.UserImpact = pointer.To(assessmentsmetadata.UserImpact(d.Get("user_impact").(string)))
	}

	if _, err := client.CreateInSubscription(ctx, *id, *existing.Model); err != nil {
		return fmt.Errorf("updating %s: %+v", *id, err)
	}

	return resourceArmSecurityCenterAssessmentPolicyRead(d, meta)
}

func resourceArmSecurityCenterAssessmentPolicyDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).SecurityCenter.AssessmentsMetadataClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := assessmentsmetadata.ParseProviderAssessmentMetadataID(d.Id())
	if err != nil {
		return err
	}

	if _, err := client.DeleteInSubscription(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return nil
}
