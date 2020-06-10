package advisor

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/advisor/mgmt/2020-01-01/advisor"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
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

			"duration_in_days": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      -1,
				ValidateFunc: validation.IntBetween(-1, 24855),
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
		SuppressionProperties: &advisor.SuppressionProperties{},
	}

	if v, ok := d.GetOk("duration_in_days"); ok && v.(int) != -1 {
		props.SuppressionProperties.TTL = utils.String(strconv.Itoa(v.(int)))
	}

	if _, err := client.Create(ctx, recommendation.ResourceUri, recommendation.Name, name, props); err != nil {
		return fmt.Errorf("failure in creating Advisor Suppressions %q: %+v", name, err)
	}

	resp, err := client.Get(ctx, recommendation.ResourceUri, recommendation.Name, name)
	if err != nil {
		return fmt.Errorf("failure in retrieving Advisor Suppressions %q: %+v", name, err)
	}
	if resp.ID == nil || *resp.ID == "" {
		return fmt.Errorf("nil or Empty ID of Advisor Suppressions %q ID", name)
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

	// ttl from api is in format dd.hh:mm:ss, we set only the day number into duration_in_days
	durationDays, err := strconv.Atoi(strings.Split(*resp.TTL, ".")[0])
	if err != nil {
		return fmt.Errorf("can't convert %s to int of field `duration_in_days` in Advisor Suppression %q: %+v", strings.Split(*resp.TTL, ".")[0], id.Name, err)
	}
	d.Set("duration_in_days", durationDays)

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
