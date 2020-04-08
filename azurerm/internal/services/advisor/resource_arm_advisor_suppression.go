package advisor

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/advisor/mgmt/2020-01-01/advisor"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/advisor/helper"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/advisor/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/advisor/validate"
	azSchema "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmAdvisorSuppression() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmAdvisorSuppressionCreateUpdate,
		Read:   resourceArmAdvisorSuppressionRead,
		Update: resourceArmAdvisorSuppressionCreateUpdate,
		Delete: resourceArmAdvisorSuppressionDelete,

		Importer: azSchema.ValidateResourceIDPriorToImport(func(id string) error {
			_, err := parse.AdvisorSuppressionID(id)
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
				ValidateFunc: validate.AdvisorSuppressionName(),
			},

			"recommendation_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.AdvisorRecommendationID,
			},

			"suppressed_duration": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      -1,
				ValidateFunc: validate.AdvisorSuppresionTTL,
			},
		},
	}
}

func resourceArmAdvisorSuppressionCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Advisor.SuppressionsClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	recommendation, _ := parse.AdvisorRecommendationID(d.Get("recommendation_id").(string))

	if d.IsNewResource() {
		existing, err := client.Get(ctx, recommendation.ResourceUri, recommendation.Name, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("failure in checking for present of existing Advisor Suppressions %q: %+v", name, err)
			}
		}
		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_advisor_suppression", *existing.ID)
		}
	}

	props := advisor.SuppressionContract{
		SuppressionProperties: &advisor.SuppressionProperties{
			TTL: utils.String(helper.ConvertToAdvisorSuppresionTTL(d.Get("suppressed_duration").(int))),
		},
	}

	if _, err := client.Create(ctx, recommendation.ResourceUri, recommendation.Name, name, props); err != nil {
		return fmt.Errorf("failure in creating Advisor Suppressions %q: %+v", name, err)
	}

	resp, err := client.Get(ctx, recommendation.ResourceUri, recommendation.Name, name)
	if err != nil {
		return fmt.Errorf("failure in retrieving Advisor Suppressions %q: %+v", name, err)
	}
	if resp.ID == nil {
		return fmt.Errorf("cannot read Advisor Suppressions %q ID", name)
	}
	d.SetId(*resp.ID)

	return resourceArmAdvisorSuppressionRead(d, meta)
}

func resourceArmAdvisorSuppressionRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Advisor.SuppressionsClient
	rclient := meta.(*clients.Client).Advisor.RecommendationsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.AdvisorSuppressionID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceUri, id.RecommendationName, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] Advisor Suppressions %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("failure in reading Advisor Suppressions %q: %+v", id.Name, err)
	}

	rResp, err := rclient.Get(ctx, id.ResourceUri, id.RecommendationName)
	if err != nil || rResp.ID == nil {
		return fmt.Errorf("failure in reading Advisor Recommendations %q: %+v", id.RecommendationName, err)
	}
	d.Set("name", id.Name)
	d.Set("recommendation_id", rResp.ID)
	if props := resp.SuppressionProperties; props != nil {
		if ttl := helper.ParseAdvisorSuppresionTTL(*props.TTL); ttl != 0 {
			d.Set("suppressed_duration", ttl)
		}
	}
	return nil
}

func resourceArmAdvisorSuppressionDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Advisor.SuppressionsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.AdvisorSuppressionID(d.Id())
	if err != nil {
		return err
	}

	if _, err = client.Delete(ctx, id.ResourceUri, id.RecommendationName, id.Name); err != nil {
		return fmt.Errorf("failure in deleting Advisor Suppressions %q: %+v", id.Name, err)
	}

	return nil
}
