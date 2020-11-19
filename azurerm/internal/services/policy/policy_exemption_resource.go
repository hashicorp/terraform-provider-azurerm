package policy

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/resources/mgmt/2020-03-01-preview/policy"
	"github.com/Azure/go-autorest/autorest/date"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/structure"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	azValidate "github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/policy/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/policy/validate"
	azSchema "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmPolicyExemption() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmPolicyExemptionCreateUpdate,
		Read:   resourceArmPolicyExemptionRead,
		Update: resourceArmPolicyExemptionCreateUpdate,
		Delete: resourceArmPolicyExemptionDelete,

		Importer: azSchema.ValidateResourceIDPriorToImport(func(id string) error {
			_, err := parse.PolicyExemptionID(id)
			return err
		}),

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"scope": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.PolicyScopeID,
			},

			"exemption_category": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(policy.Mitigated),
					string(policy.Waiver),
				}, false),
			},

			"policy_assignment_id": {
				Type:         schema.TypeString,
				Required:     true,
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

			"metadata": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsJSON,
			},
		},
	}
}

func resourceArmPolicyExemptionCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Policy.ExemptionsClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	scope := d.Get("scope").(string)

	if d.IsNewResource() {
		existing, err := client.Get(ctx, scope, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing Policy Exemption %q (Scope %q): %+v", name, scope, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_policy_exemption", *existing.ID)
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

	if _, err := client.CreateOrUpdate(ctx, exemption, scope, name); err != nil {
		return fmt.Errorf("creating/updating Policy Exemption %q (Scope %q): %+v", name, scope, err)
	}

	resp, err := client.Get(ctx, scope, name)
	if err != nil {
		return fmt.Errorf("retrieving Policy Exemption %q (Scope %q): %+v", name, scope, err)
	}

	if resp.ID == nil || *resp.ID == "" {
		return fmt.Errorf("empty or nil ID returned for Policy Exemption %q (Scope %q)", name, scope)
	}
	d.SetId(*resp.ID)

	return resourceArmPolicyExemptionRead(d, meta)
}

func resourceArmPolicyExemptionRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Policy.ExemptionsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.PolicyExemptionID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ScopeId(), id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] Error reading Policy Exemption %q - removing from state", d.Id())
			d.SetId("")
			return nil
		}

		return fmt.Errorf("reading Policy Exemption %q (Scope %q): %+v", id.Name, id.ScopeId(), err)
	}

	d.Set("name", resp.Name)
	d.Set("scope", id.ScopeId())
	if props := resp.ExemptionProperties; props != nil {
		d.Set("policy_assignment_id", props.PolicyAssignmentID)
		d.Set("display_name", props.DisplayName)
		d.Set("description", props.Description)
		d.Set("exemption_category", string(props.ExemptionCategory))

		if err := d.Set("policy_definition_reference_ids", utils.FlattenStringSlice(props.PolicyDefinitionReferenceIds)); err != nil {
			return fmt.Errorf("setting `policy_definition_reference_ids: %+v", err)
		}

		if expiresTime := props.ExpiresOn; expiresTime != nil {
			d.Set("expires_on", expiresTime.String())
		}

		if metadataStr := flattenJSON(props.Metadata); metadataStr != "" {
			d.Set("metadata", metadataStr)
		}
	}

	return nil
}

func resourceArmPolicyExemptionDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Policy.ExemptionsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.PolicyExemptionID(d.Id())
	if err != nil {
		return err
	}

	if _, err := client.Delete(ctx, id.ScopeId(), id.Name); err != nil {
		return fmt.Errorf("deleting Policy Exemption %q (Scope %q): %+v", id.Name, id.ScopeId(), err)
	}

	return nil
}
