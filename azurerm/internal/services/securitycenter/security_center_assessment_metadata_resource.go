package securitycenter

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/security/mgmt/v3.0/security"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/securitycenter/parse"
	azSchema "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmSecurityCenterAssessmentMetadata() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmSecurityCenterAssessmentMetadataCreateUpdate,
		Read:   resourceArmSecurityCenterAssessmentMetadataRead,
		Update: resourceArmSecurityCenterAssessmentMetadataCreateUpdate,
		Delete: resourceArmSecurityCenterAssessmentMetadataDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Importer: azSchema.ValidateResourceIDPriorToImport(func(id string) error {
			_, err := parse.AssessmentMetadataID(id)
			return err
		}),

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsUUID,
			},

			"description": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"display_name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"assessment_type": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  string(security.CustomerManaged),
				ValidateFunc: validation.StringInSlice([]string{
					string(security.CustomerManaged),
					string(security.VerifiedPartner),
				}, false),
			},

			"severity": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  string(security.SeverityMedium),
				ValidateFunc: validation.StringInSlice([]string{
					string(security.SeverityLow),
					string(security.SeverityMedium),
					string(security.SeverityHigh),
				}, false),
			},

			// The Azure API always returns `nil` for this property after set it.
			// BUG: https://github.com/Azure/azure-rest-api-specs/issues/12297
			"categories": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
					ValidateFunc: validation.StringInSlice([]string{
						string(security.Compute),
						string(security.Networking),
						string(security.Data),
						string(security.IdentityAndAccess),
						string(security.IoT),
					}, false),
				},
			},

			"implementation_effort": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(security.ImplementationEffortLow),
					string(security.ImplementationEffortModerate),
					string(security.ImplementationEffortHigh),
				}, false),
			},

			"is_preview": {
				Type:     schema.TypeBool,
				Optional: true,
			},

			"remediation_description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"threats": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
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
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(security.UserImpactLow),
					string(security.UserImpactModerate),
					string(security.UserImpactHigh),
				}, false),
			},
		},
	}
}

func resourceArmSecurityCenterAssessmentMetadataCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	client := meta.(*clients.Client).SecurityCenter.AssessmentsMetadataClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)

	id := parse.NewAssessmentMetadataID(subscriptionId, name)

	if d.IsNewResource() {
		existing, err := client.GetInSubscription(ctx, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for present of existing Security Center Assessment Metadata %q : %+v", name, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_security_center_assessment_metadata", id.ID())
		}
	}

	params := security.AssessmentMetadata{
		AssessmentMetadataProperties: &security.AssessmentMetadataProperties{
			AssessmentType: security.AssessmentType(d.Get("assessment_type").(string)),
			Description:    utils.String(d.Get("description").(string)),
			DisplayName:    utils.String(d.Get("display_name").(string)),
			Severity:       security.Severity(d.Get("severity").(string)),
		},
	}

	if v, ok := d.GetOk("categories"); ok {
		category := make([]security.Category, 0)
		for _, item := range *(utils.ExpandStringSlice(v.(*schema.Set).List())) {
			category = append(category, (security.Category)(item))
		}
		params.AssessmentMetadataProperties.Category = &category
	}

	if v, ok := d.GetOk("threats"); ok {
		threats := make([]security.Threats, 0)
		for _, item := range *(utils.ExpandStringSlice(v.(*schema.Set).List())) {
			threats = append(threats, (security.Threats)(item))
		}
		params.AssessmentMetadataProperties.Threats = &threats
	}

	if v, ok := d.GetOk("implementation_effort"); ok {
		params.AssessmentMetadataProperties.ImplementationEffort = security.ImplementationEffort(v.(string))
	}

	if v, ok := d.GetOk("is_preview"); ok {
		params.AssessmentMetadataProperties.Preview = utils.Bool(v.(bool))
	}

	if v, ok := d.GetOk("remediation_description"); ok {
		params.AssessmentMetadataProperties.RemediationDescription = utils.String(v.(string))
	}

	if v, ok := d.GetOk("user_impact"); ok {
		params.AssessmentMetadataProperties.UserImpact = security.UserImpact(v.(string))
	}

	if _, err := client.CreateInSubscription(ctx, name, params); err != nil {
		return fmt.Errorf("creating/updating Security Center Assessment Metadata %q : %+v", name, err)
	}

	d.SetId(id.ID())

	return resourceArmSecurityCenterAssessmentMetadataRead(d, meta)
}

func resourceArmSecurityCenterAssessmentMetadataRead(d *schema.ResourceData, meta interface{}) error {
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
			log.Printf("[INFO] security Center Assessment Metadata %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving Security Center Assessment Metadata %q : %+v", id.AssessmentMetadataName, err)
	}

	d.Set("name", id.AssessmentMetadataName)

	if props := resp.AssessmentMetadataProperties; props != nil {
		d.Set("assessment_type", props.AssessmentType)
		d.Set("description", props.Description)
		d.Set("display_name", props.DisplayName)
		d.Set("severity", props.Severity)
		d.Set("implementation_effort", props.ImplementationEffort)
		d.Set("is_preview", props.Preview)
		d.Set("remediation_description", props.RemediationDescription)
		d.Set("user_impact", props.UserImpact)

		threats := make([]string, 0)
		if props.Threats != nil {
			for _, item := range *props.Threats {
				threats = append(threats, (string)(item))
			}
		}
		d.Set("threats", utils.FlattenStringSlice(&threats))
	}

	return nil
}

func resourceArmSecurityCenterAssessmentMetadataDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).SecurityCenter.AssessmentsMetadataClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.AssessmentMetadataID(d.Id())
	if err != nil {
		return err
	}

	if _, err := client.DeleteInSubscription(ctx, id.AssessmentMetadataName); err != nil {
		return fmt.Errorf("deleting Security Center Assessment Metadata %q : %+v", id.AssessmentMetadataName, err)
	}

	return nil
}
