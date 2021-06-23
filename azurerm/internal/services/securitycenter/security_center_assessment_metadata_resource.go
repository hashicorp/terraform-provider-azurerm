package securitycenter

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/security/mgmt/v3.0/security"
	"github.com/google/uuid"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/securitycenter/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmSecurityCenterAssessmentMetadata() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create:             resourceArmSecurityCenterAssessmentMetadataCreate,
		Read:               resourceArmSecurityCenterAssessmentMetadataRead,
		Update:             resourceArmSecurityCenterAssessmentMetadataUpdate,
		Delete:             resourceArmSecurityCenterAssessmentMetadataDelete,
		DeprecationMessage: "This resource has been renamed to `azurerm_security_center_assessment_policy` and will be removed in version 3.0 of the provider.",

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.AssessmentMetadataID(id)
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

func resourceArmSecurityCenterAssessmentMetadataCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	client := meta.(*clients.Client).SecurityCenter.AssessmentsMetadataClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := uuid.New().String()

	id := parse.NewAssessmentMetadataID(subscriptionId, name)

	existing, err := client.GetInSubscription(ctx, name)
	if err != nil {
		if !utils.ResponseWasNotFound(existing.Response) {
			return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
		}
	}

	if !utils.ResponseWasNotFound(existing.Response) {
		return tf.ImportAsExistsError("azurerm_security_center_assessment_metadata", id.ID())
	}

	params := security.AssessmentMetadata{
		AssessmentMetadataProperties: &security.AssessmentMetadataProperties{
			AssessmentType: security.CustomerManaged,
			Description:    utils.String(d.Get("description").(string)),
			DisplayName:    utils.String(d.Get("display_name").(string)),
			Severity:       security.Severity(d.Get("severity").(string)),
		},
	}

	if v, ok := d.GetOk("categories"); ok {
		categories := make([]security.Categories, 0)
		for _, item := range v.(*pluginsdk.Set).List() {
			categories = append(categories, (security.Categories)(item.(string)))
		}
		params.AssessmentMetadataProperties.Categories = &categories
	}

	if v, ok := d.GetOk("threats"); ok {
		threats := make([]security.Threats, 0)
		for _, item := range v.(*pluginsdk.Set).List() {
			threats = append(threats, (security.Threats)(item.(string)))
		}
		params.AssessmentMetadataProperties.Threats = &threats
	}

	if v, ok := d.GetOk("implementation_effort"); ok {
		params.AssessmentMetadataProperties.ImplementationEffort = security.ImplementationEffort(v.(string))
	}

	if v, ok := d.GetOk("remediation_description"); ok {
		params.AssessmentMetadataProperties.RemediationDescription = utils.String(v.(string))
	}

	if v, ok := d.GetOk("user_impact"); ok {
		params.AssessmentMetadataProperties.UserImpact = security.UserImpact(v.(string))
	}

	if _, err := client.CreateInSubscription(ctx, name, params); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceArmSecurityCenterAssessmentMetadataRead(d, meta)
}

func resourceArmSecurityCenterAssessmentMetadataRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).SecurityCenter.AssessmentsMetadataClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.AssessmentMetadataID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.GetInSubscription(ctx, id.AssessmentMetadataName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] %s does not exist - removing from state", id)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.Set("name", id.AssessmentMetadataName)

	if props := resp.AssessmentMetadataProperties; props != nil {
		d.Set("description", props.Description)
		d.Set("display_name", props.DisplayName)
		d.Set("severity", string(props.Severity))
		d.Set("implementation_effort", string(props.ImplementationEffort))
		d.Set("remediation_description", props.RemediationDescription)
		d.Set("user_impact", string(props.UserImpact))

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

	return nil
}

func resourceArmSecurityCenterAssessmentMetadataUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).SecurityCenter.AssessmentsMetadataClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.AssessmentMetadataID(d.Id())
	if err != nil {
		return err
	}

	existing, err := client.GetInSubscription(ctx, id.AssessmentMetadataName)
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}
	if existing.AssessmentMetadataProperties == nil {
		return fmt.Errorf("retrieving %s: `assessmentMetadataProperties` was nil", id)
	}

	if d.HasChange("description") {
		existing.AssessmentMetadataProperties.Description = utils.String(d.Get("description").(string))
	}

	if d.HasChange("display_name") {
		existing.AssessmentMetadataProperties.DisplayName = utils.String(d.Get("display_name").(string))
	}

	if d.HasChange("severity") {
		existing.AssessmentMetadataProperties.Severity = security.Severity(d.Get("severity").(string))
	}

	if d.HasChange("categories") {
		categories := make([]security.Categories, 0)
		for _, item := range d.Get("categories").(*pluginsdk.Set).List() {
			categories = append(categories, (security.Categories)(item.(string)))
		}
		existing.AssessmentMetadataProperties.Categories = &categories
	}

	if d.HasChange("threats") {
		threats := make([]security.Threats, 0)
		for _, item := range d.Get("threats").(*pluginsdk.Set).List() {
			threats = append(threats, (security.Threats)(item.(string)))
		}
		existing.AssessmentMetadataProperties.Threats = &threats
	}

	if d.HasChange("implementation_effort") {
		existing.AssessmentMetadataProperties.ImplementationEffort = security.ImplementationEffort(d.Get("implementation_effort").(string))
	}

	if d.HasChange("remediation_description") {
		existing.AssessmentMetadataProperties.RemediationDescription = utils.String(d.Get("remediation_description").(string))
	}

	if d.HasChange("user_impact") {
		existing.AssessmentMetadataProperties.UserImpact = security.UserImpact(d.Get("user_impact").(string))
	}

	if _, err := client.CreateInSubscription(ctx, id.AssessmentMetadataName, existing); err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}

	return resourceArmSecurityCenterAssessmentMetadataRead(d, meta)
}

func resourceArmSecurityCenterAssessmentMetadataDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).SecurityCenter.AssessmentsMetadataClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.AssessmentMetadataID(d.Id())
	if err != nil {
		return err
	}

	if _, err := client.DeleteInSubscription(ctx, id.AssessmentMetadataName); err != nil {
		return fmt.Errorf("deleting %s: %+v", id, err)
	}

	return nil
}
