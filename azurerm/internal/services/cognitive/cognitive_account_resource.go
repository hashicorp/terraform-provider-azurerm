package cognitive

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/cognitiveservices/mgmt/2017-04-18/cognitiveservices"
	"github.com/hashicorp/go-azure-helpers/response"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/cognitive/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	azSchema "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmCognitiveAccount() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmCognitiveAccountCreate,
		Read:   resourceArmCognitiveAccountRead,
		Update: resourceArmCognitiveAccountUpdate,
		Delete: resourceArmCognitiveAccountDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Importer: azSchema.ValidateResourceIDPriorToImport(func(id string) error {
			_, err := parse.CognitiveAccountID(id)
			return err
		}),

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.CognitiveServicesAccountName(),
			},

			"location": azure.SchemaLocation(),

			"resource_group_name": azure.SchemaResourceGroupName(),

			"kind": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"Academic",
					"Bing.Autosuggest",
					"Bing.Autosuggest.v7",
					"Bing.CustomSearch",
					"Bing.Search",
					"Bing.Search.v7",
					"Bing.Speech",
					"Bing.SpellCheck",
					"Bing.SpellCheck.v7",
					"CognitiveServices",
					"ComputerVision",
					"ContentModerator",
					"CustomSpeech",
					"CustomVision.Prediction",
					"CustomVision.Training",
					"Emotion",
					"Face",
					"FormRecognizer",
					"ImmersiveReader",
					"LUIS",
					"LUIS.Authoring",
					"QnAMaker",
					"Recommendations",
					"SpeakerRecognition",
					"Speech",
					"SpeechServices",
					"SpeechTranslation",
					"TextAnalytics",
					"TextTranslation",
					"WebLM",
				}, false),
			},

			"sku_name": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"F0", "F1", "S0", "S1", "S2", "S3", "S4", "S5", "S6", "P0", "P1", "P2",
				}, false),
			},

			"qna_runtime_endpoint": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.IsURLWithHTTPorHTTPS,
			},

			"tags": tags.Schema(),

			"endpoint": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"primary_access_key": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},

			"secondary_access_key": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},
		},
	}
}

func resourceArmCognitiveAccountCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cognitive.AccountsClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	kind := d.Get("kind").(string)

	if d.IsNewResource() {
		existing, err := client.GetProperties(ctx, resourceGroup, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of existing Cognitive Account %q (Resource Group %q): %s", name, resourceGroup, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_cognitive_account", *existing.ID)
		}
	}

	sku, err := expandAccountSkuName(d.Get("sku_name").(string))
	if err != nil {
		return fmt.Errorf("error expanding sku_name for Cognitive Account %s (Resource Group %q): %v", name, resourceGroup, err)
	}

	props := cognitiveservices.Account{
		Kind:     utils.String(kind),
		Location: utils.String(azure.NormalizeLocation(d.Get("location").(string))),
		Sku:      sku,
		Properties: &cognitiveservices.AccountProperties{
			APIProperties: &cognitiveservices.AccountAPIProperties{
				QnaRuntimeEndpoint: utils.String(d.Get("qna_runtime_endpoint").(string)),
			},
		},
		Tags: tags.Expand(d.Get("tags").(map[string]interface{})),
	}

	if kind == "QnAMaker" && *props.Properties.APIProperties.QnaRuntimeEndpoint == "" {
		return fmt.Errorf("the QnAMaker runtime endpoint `qna_runtime_endpoint` is required when kind is set to `QnAMaker`")
	}

	if _, err := client.Create(ctx, resourceGroup, name, props); err != nil {
		return fmt.Errorf("Error creating Cognitive Services Account %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	read, err := client.GetProperties(ctx, resourceGroup, name)
	if err != nil {
		return fmt.Errorf("Error retrieving Cognitive Services Account %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	d.SetId(*read.ID)

	return resourceArmCognitiveAccountRead(d, meta)
}

func resourceArmCognitiveAccountUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cognitive.AccountsClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.CognitiveAccountID(d.Id())
	if err != nil {
		return err
	}

	sku, err := expandAccountSkuName(d.Get("sku_name").(string))
	if err != nil {
		return fmt.Errorf("error expanding sku_name for Cognitive Account %s (Resource Group %q): %v", id.Name, id.ResourceGroup, err)
	}

	props := cognitiveservices.Account{
		Sku: sku,
		Properties: &cognitiveservices.AccountProperties{
			APIProperties: &cognitiveservices.AccountAPIProperties{
				QnaRuntimeEndpoint: utils.String(d.Get("qna_runtime_endpoint").(string)),
			},
		},
		Tags: tags.Expand(d.Get("tags").(map[string]interface{})),
	}

	if d.Get("kind").(string) == "QnAMaker" && *props.Properties.APIProperties.QnaRuntimeEndpoint == "" {
		return fmt.Errorf("the QnAMaker runtime endpoint `qna_runtime_endpoint` is required when kind is set to `QnAMaker`")
	}

	if _, err = client.Update(ctx, id.ResourceGroup, id.Name, props); err != nil {
		return fmt.Errorf("Error updating Cognitive Services Account %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	return resourceArmCognitiveAccountRead(d, meta)
}

func resourceArmCognitiveAccountRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cognitive.AccountsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.CognitiveAccountID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.GetProperties(ctx, id.ResourceGroup, id.Name)

	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] Cognitive Services Account %q was not found in Resource Group %q - removing from state!", id.Name, id.ResourceGroup)
			d.SetId("")
			return nil
		}
		return err
	}

	d.Set("name", id.Name)
	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("kind", resp.Kind)

	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

	if sku := resp.Sku; sku != nil {
		d.Set("sku_name", sku.Name)
	}

	if props := resp.Properties; props != nil {
		if apiProps := props.APIProperties; apiProps != nil {
			d.Set("qna_runtime_endpoint", apiProps.QnaRuntimeEndpoint)
		}
		d.Set("endpoint", props.Endpoint)
	}

	keys, err := client.ListKeys(ctx, id.ResourceGroup, id.Name)

	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] Not able to obtain keys for Cognitive Services Account %q in Resource Group %q - removing from state!", id.Name, id.ResourceGroup)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error obtaining keys for Cognitive Services Account %q in Resource Group %q: %v", id.Name, id.ResourceGroup, err)
	}

	d.Set("primary_access_key", keys.Key1)
	d.Set("secondary_access_key", keys.Key2)

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceArmCognitiveAccountDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cognitive.AccountsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.CognitiveAccountID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Delete(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		if !response.WasNotFound(resp.Response) {
			return fmt.Errorf("Error deleting Cognitive Services Account %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
		}
	}

	return nil
}

func expandAccountSkuName(skuName string) (*cognitiveservices.Sku, error) {
	var tier cognitiveservices.SkuTier
	switch skuName[0:1] {
	case "F":
		tier = cognitiveservices.Free
	case "S":
		tier = cognitiveservices.Standard
	case "P":
		tier = cognitiveservices.Premium
	default:
		return nil, fmt.Errorf("sku_name %s has unknown sku tier %s", skuName, skuName[0:1])
	}

	return &cognitiveservices.Sku{
		Name: utils.String(skuName),
		Tier: tier,
	}, nil
}
