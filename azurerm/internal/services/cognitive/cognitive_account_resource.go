package cognitive

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/cognitiveservices/mgmt/2017-04-18/cognitiveservices"
	"github.com/hashicorp/go-azure-helpers/response"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	commonValidate "github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/locks"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/cognitive/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/cognitive/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/network"
	networkParse "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/network/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/set"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceCognitiveAccount() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceCognitiveAccountCreate,
		Read:   resourceCognitiveAccountRead,
		Update: resourceCognitiveAccountUpdate,
		Delete: resourceCognitiveAccountDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.AccountID(id)
			return err
		}),

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.CognitiveServicesAccountName(),
			},

			"location": azure.SchemaLocation(),

			"resource_group_name": azure.SchemaResourceGroupName(),

			"kind": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"Academic",
					"AnomalyDetector",
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
					"Personalizer",
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
				Type:     pluginsdk.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"F0", "F1", "S0", "S", "S1", "S2", "S3", "S4", "S5", "S6", "P0", "P1", "P2",
				}, false),
			},

			"qna_runtime_endpoint": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: validation.IsURLWithHTTPorHTTPS,
			},

			"network_acls": {
				Type:         pluginsdk.TypeList,
				Optional:     true,
				MaxItems:     1,
				RequiredWith: []string{"custom_subdomain_name"},
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"default_action": {
							Type:     pluginsdk.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(cognitiveservices.Allow),
								string(cognitiveservices.Deny),
							}, false),
						},
						"ip_rules": {
							Type:     pluginsdk.TypeSet,
							Optional: true,
							Elem: &pluginsdk.Schema{
								Type: pluginsdk.TypeString,
								ValidateFunc: validation.Any(
									commonValidate.IPv4Address,
									commonValidate.CIDR,
								),
							},
							Set: set.HashIPv4AddressOrCIDR,
						},
						"virtual_network_subnet_ids": {
							Type:     pluginsdk.TypeSet,
							Optional: true,
							Elem:     &pluginsdk.Schema{Type: pluginsdk.TypeString},
						},
					},
				},
			},

			"custom_subdomain_name": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"tags": tags.Schema(),

			"endpoint": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"primary_access_key": {
				Type:      pluginsdk.TypeString,
				Computed:  true,
				Sensitive: true,
			},

			"secondary_access_key": {
				Type:      pluginsdk.TypeString,
				Computed:  true,
				Sensitive: true,
			},
		},
	}
}

func resourceCognitiveAccountCreate(d *pluginsdk.ResourceData, meta interface{}) error {
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
				return fmt.Errorf("checking for presence of existing Cognitive Account %q (Resource Group %q): %s", name, resourceGroup, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_cognitive_account", *existing.ID)
		}
	}

	sku, err := expandAccountSkuName(d.Get("sku_name").(string))
	if err != nil {
		return fmt.Errorf("expanding sku_name for Cognitive Account %s (Resource Group %q): %v", name, resourceGroup, err)
	}

	networkAcls, subnetIds := expandCognitiveAccountNetworkAcls(d.Get("network_acls").([]interface{}))

	// also lock on the Virtual Network ID's since modifications in the networking stack are exclusive
	virtualNetworkNames := make([]string, 0)
	for _, v := range subnetIds {
		id, err := networkParse.SubnetIDInsensitively(v)
		if err != nil {
			return err
		}
		if !utils.SliceContainsValue(virtualNetworkNames, id.VirtualNetworkName) {
			virtualNetworkNames = append(virtualNetworkNames, id.VirtualNetworkName)
		}
	}

	locks.MultipleByName(&virtualNetworkNames, network.VirtualNetworkResourceName)
	defer locks.UnlockMultipleByName(&virtualNetworkNames, network.VirtualNetworkResourceName)

	props := cognitiveservices.Account{
		Kind:     utils.String(kind),
		Location: utils.String(azure.NormalizeLocation(d.Get("location").(string))),
		Sku:      sku,
		Properties: &cognitiveservices.AccountProperties{
			APIProperties:       &cognitiveservices.AccountAPIProperties{},
			NetworkAcls:         networkAcls,
			CustomSubDomainName: utils.String(d.Get("custom_subdomain_name").(string)),
		},
		Tags: tags.Expand(d.Get("tags").(map[string]interface{})),
	}

	if kind == "QnAMaker" {
		if v, ok := d.GetOk("qna_runtime_endpoint"); ok && v != "" {
			props.Properties.APIProperties.QnaRuntimeEndpoint = utils.String(v.(string))
		} else {
			return fmt.Errorf("the QnAMaker runtime endpoint `qna_runtime_endpoint` is required when kind is set to `QnAMaker`")
		}
	}

	if _, err := client.Create(ctx, resourceGroup, name, props); err != nil {
		return fmt.Errorf("creating Cognitive Services Account %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	stateConf := &pluginsdk.StateChangeConf{
		Pending:    []string{"Creating"},
		Target:     []string{"Succeeded"},
		Refresh:    cognitiveAccountStateRefreshFunc(ctx, client, resourceGroup, name),
		MinTimeout: 15 * time.Second,
		Timeout:    d.Timeout(pluginsdk.TimeoutCreate),
	}

	if _, err = stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("waiting for Cognitive Account (%s) to become available: %s", d.Get("name"), err)
	}

	read, err := client.GetProperties(ctx, resourceGroup, name)
	if err != nil {
		return fmt.Errorf("retrieving Cognitive Services Account %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	d.SetId(*read.ID)

	return resourceCognitiveAccountRead(d, meta)
}

func resourceCognitiveAccountUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cognitive.AccountsClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.AccountID(d.Id())
	if err != nil {
		return err
	}

	sku, err := expandAccountSkuName(d.Get("sku_name").(string))
	if err != nil {
		return fmt.Errorf("error expanding sku_name for Cognitive Account %s (Resource Group %q): %v", id.Name, id.ResourceGroup, err)
	}

	networkAcls, subnetIds := expandCognitiveAccountNetworkAcls(d.Get("network_acls").([]interface{}))

	// also lock on the Virtual Network ID's since modifications in the networking stack are exclusive
	virtualNetworkNames := make([]string, 0)
	for _, v := range subnetIds {
		id, err := networkParse.SubnetIDInsensitively(v)
		if err != nil {
			return err
		}
		if !utils.SliceContainsValue(virtualNetworkNames, id.VirtualNetworkName) {
			virtualNetworkNames = append(virtualNetworkNames, id.VirtualNetworkName)
		}
	}

	locks.MultipleByName(&virtualNetworkNames, network.VirtualNetworkResourceName)
	defer locks.UnlockMultipleByName(&virtualNetworkNames, network.VirtualNetworkResourceName)

	props := cognitiveservices.Account{
		Sku: sku,
		Properties: &cognitiveservices.AccountProperties{
			APIProperties:       &cognitiveservices.AccountAPIProperties{},
			NetworkAcls:         networkAcls,
			CustomSubDomainName: utils.String(d.Get("custom_subdomain_name").(string)),
		},
		Tags: tags.Expand(d.Get("tags").(map[string]interface{})),
	}

	if kind := d.Get("kind"); kind == "QnAMaker" {
		if v, ok := d.GetOk("qna_runtime_endpoint"); ok && v != "" {
			props.Properties.APIProperties.QnaRuntimeEndpoint = utils.String(v.(string))
		} else {
			return fmt.Errorf("the QnAMaker runtime endpoint `qna_runtime_endpoint` is required when kind is set to `QnAMaker`")
		}
	}

	if _, err = client.Update(ctx, id.ResourceGroup, id.Name, props); err != nil {
		return fmt.Errorf("Error updating Cognitive Services Account %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	stateConf := &pluginsdk.StateChangeConf{
		Pending:    []string{"Accepted"},
		Target:     []string{"Succeeded"},
		Refresh:    cognitiveAccountStateRefreshFunc(ctx, client, id.ResourceGroup, id.Name),
		MinTimeout: 15 * time.Second,
		Timeout:    d.Timeout(pluginsdk.TimeoutCreate),
	}

	if _, err = stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("waiting for Cognitive Account (%s) to become available: %s", d.Get("name"), err)
	}

	return resourceCognitiveAccountRead(d, meta)
}

func resourceCognitiveAccountRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cognitive.AccountsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.AccountID(d.Id())
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
		d.Set("custom_subdomain_name", props.CustomSubDomainName)
		if err := d.Set("network_acls", flattenCognitiveAccountNetworkAcls(props.NetworkAcls)); err != nil {
			return fmt.Errorf("setting `network_acls` for Cognitive Account %q: %+v", *resp.Name, err)
		}
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

func resourceCognitiveAccountDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cognitive.AccountsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.AccountID(d.Id())
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

func cognitiveAccountStateRefreshFunc(ctx context.Context, client *cognitiveservices.AccountsClient, resourceGroupName string, cognitiveAccountName string) pluginsdk.StateRefreshFunc {
	return func() (interface{}, string, error) {
		res, err := client.GetProperties(ctx, resourceGroupName, cognitiveAccountName)
		if err != nil {
			return nil, "", fmt.Errorf("issuing read request in cognitiveAccountStateRefreshFunc to Azure ARM for Cognitive Account '%s' (RG: '%s'): %s", cognitiveAccountName, resourceGroupName, err)
		}

		return res, string(res.Properties.ProvisioningState), nil
	}
}

func expandCognitiveAccountNetworkAcls(input []interface{}) (*cognitiveservices.NetworkRuleSet, []string) {
	subnetIds := make([]string, 0)
	if len(input) == 0 || input[0] == nil {
		return nil, subnetIds
	}

	v := input[0].(map[string]interface{})

	defaultAction := v["default_action"].(string)

	ipRulesRaw := v["ip_rules"].(*pluginsdk.Set)
	ipRules := make([]cognitiveservices.IPRule, 0)

	for _, v := range ipRulesRaw.List() {
		rule := cognitiveservices.IPRule{
			Value: utils.String(v.(string)),
		}
		ipRules = append(ipRules, rule)
	}

	networkRulesRaw := v["virtual_network_subnet_ids"].(*pluginsdk.Set)
	networkRules := make([]cognitiveservices.VirtualNetworkRule, 0)
	for _, v := range networkRulesRaw.List() {
		rawId := v.(string)
		subnetIds = append(subnetIds, rawId)
		rule := cognitiveservices.VirtualNetworkRule{
			ID: utils.String(rawId),
		}
		networkRules = append(networkRules, rule)
	}

	ruleSet := cognitiveservices.NetworkRuleSet{
		DefaultAction:       cognitiveservices.NetworkRuleAction(defaultAction),
		IPRules:             &ipRules,
		VirtualNetworkRules: &networkRules,
	}
	return &ruleSet, subnetIds
}

func flattenCognitiveAccountNetworkAcls(input *cognitiveservices.NetworkRuleSet) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	output := make(map[string]interface{})

	output["default_action"] = string(input.DefaultAction)

	ipRules := make([]interface{}, 0)
	if input.IPRules != nil {
		for _, v := range *input.IPRules {
			if v.Value == nil {
				continue
			}

			ipRules = append(ipRules, *v.Value)
		}
	}
	output["ip_rules"] = pluginsdk.NewSet(pluginsdk.HashString, ipRules)

	virtualNetworkRules := make([]interface{}, 0)
	if input.VirtualNetworkRules != nil {
		for _, v := range *input.VirtualNetworkRules {
			if v.ID == nil {
				continue
			}

			id := *v.ID
			subnetId, err := networkParse.SubnetIDInsensitively(*v.ID)
			if err == nil {
				id = subnetId.ID()
			}

			virtualNetworkRules = append(virtualNetworkRules, id)
		}
	}
	output["virtual_network_subnet_ids"] = pluginsdk.NewSet(pluginsdk.HashString, virtualNetworkRules)

	return []interface{}{output}
}
