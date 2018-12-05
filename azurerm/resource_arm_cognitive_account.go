package azurerm

import (
	"fmt"
	"log"

	"github.com/Azure/azure-sdk-for-go/services/cognitiveservices/mgmt/2017-04-18/cognitiveservices"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/response"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type cognitiveServicesPropertiesStruct struct{}

func resourceArmCognitiveAccount() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmCognitiveAccountCreate,
		Read:   resourceArmCognitiveAccountRead,
		Update: resourceArmCognitiveAccountUpdate,
		Delete: resourceArmCognitiveAccountDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.CognitiveServicesAccountName(),
			},

			"location": locationSchema(),

			"resource_group_name": resourceGroupNameSchema(),

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
					"ComputerVision",
					"ContentModerator",
					"CustomSpeech",
					"Emotion",
					"Face",
					"LUIS",
					"Recommendations",
					"SpeakerRecognition",
					"Speech",
					"SpeechTranslation",
					"TextAnalytics",
					"TextTranslation",
					"WebLM",
				}, false),
			},

			"sku": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(cognitiveservices.F0),
								string(cognitiveservices.S0),
								string(cognitiveservices.S1),
								string(cognitiveservices.S2),
								string(cognitiveservices.S3),
								string(cognitiveservices.S4),
								string(cognitiveservices.S5),
								string(cognitiveservices.S6),
								string(cognitiveservices.P0),
								string(cognitiveservices.P1),
								string(cognitiveservices.P2),
							}, false),
						},

						"tier": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(cognitiveservices.Free),
								string(cognitiveservices.Standard),
								string(cognitiveservices.Premium),
							}, false),
						},
					},
				},
			},

			"tags": tagsSchema(),

			"endpoint": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceArmCognitiveAccountCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).cognitiveAccountsClient
	ctx := meta.(*ArmClient).StopContext

	name := d.Get("name").(string)
	location := azure.NormalizeLocation(d.Get("location"))
	resourceGroup := d.Get("resource_group_name").(string)
	kind := d.Get("kind").(string)
	tags := d.Get("tags").(map[string]interface{})
	sku := expandCognitiveAccountSku(d)

	properties := cognitiveservices.AccountCreateParameters{
		Kind:       cognitiveservices.Kind(kind),
		Location:   utils.String(location),
		Sku:        sku,
		Properties: &cognitiveServicesPropertiesStruct{},
		Tags:       expandTags(tags),
	}

	_, err := client.Create(ctx, resourceGroup, name, properties)
	if err != nil {
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
	client := meta.(*ArmClient).cognitiveAccountsClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	resourceGroup := id.ResourceGroup
	name := id.Path["accounts"]

	tags := d.Get("tags").(map[string]interface{})
	sku := expandCognitiveAccountSku(d)

	properties := cognitiveservices.AccountUpdateParameters{
		Sku:  sku,
		Tags: expandTags(tags),
	}

	_, err = client.Update(ctx, resourceGroup, name, properties)
	if err != nil {
		return fmt.Errorf("Error updating Cognitive Services Account %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	return resourceArmCognitiveAccountRead(d, meta)
}

func resourceArmCognitiveAccountRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).cognitiveAccountsClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	resourceGroup := id.ResourceGroup
	name := id.Path["accounts"]

	resp, err := client.GetProperties(ctx, resourceGroup, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] Cognitive Services Account %q was not found in Resource Group %q - removing from state!", name, resourceGroup)
			d.SetId("")
			return nil
		}
		return err
	}

	d.Set("name", resp.Name)
	d.Set("resource_group_name", resourceGroup)
	d.Set("kind", resp.Kind)

	if err := azure.FlattenAndSetLocation(d, resp.Location); err != nil {
		return err
	}

	if err := d.Set("sku", flattenCognitiveAccountSku(resp.Sku)); err != nil {
		return fmt.Errorf("Error setting `sku`: %+v", err)
	}

	if props := resp.AccountProperties; props != nil {
		d.Set("endpoint", props.Endpoint)
	}

	flattenAndSetTags(d, resp.Tags)

	return nil
}

func resourceArmCognitiveAccountDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).cognitiveAccountsClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	resourceGroup := id.ResourceGroup
	name := id.Path["accounts"]

	resp, err := client.Delete(ctx, resourceGroup, name)
	if err != nil {
		if !response.WasNotFound(resp.Response) {
			return fmt.Errorf("Error deleting Cognitive Services Account %q (Resource Group %q): %+v", name, resourceGroup, err)
		}
	}

	return nil
}

func expandCognitiveAccountSku(d *schema.ResourceData) *cognitiveservices.Sku {
	skus := d.Get("sku").([]interface{})
	sku := skus[0].(map[string]interface{})

	return &cognitiveservices.Sku{
		Name: cognitiveservices.SkuName(sku["name"].(string)),
		Tier: cognitiveservices.SkuTier(sku["tier"].(string)),
	}
}

func flattenCognitiveAccountSku(input *cognitiveservices.Sku) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	return []interface{}{
		map[string]interface{}{
			"name": string(input.Name),
			"tier": string(input.Tier),
		},
	}
}
