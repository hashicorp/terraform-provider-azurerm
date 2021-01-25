package securitycenter

import (
	"fmt"
	"log"
	"time"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/suppress"

	"github.com/Azure/azure-sdk-for-go/services/preview/security/mgmt/v3.0/security"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/location"
	iothubValidate "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/iothub/validate"
	loganalyticsValidate "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/loganalytics/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/securitycenter/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/securitycenter/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	azSchema "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/set"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceIotSecuritySolution() *schema.Resource {
	return &schema.Resource{
		Create: resourceIotSecuritySolutionCreateUpdate,
		Read:   resourceIotSecuritySolutionRead,
		Update: resourceIotSecuritySolutionCreateUpdate,
		Delete: resourceIotSecuritySolutionDelete,

		Importer: azSchema.ValidateResourceIDPriorToImport(func(id string) error {
			_, err := parse.IotSecuritySolutionID(id)
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
				ValidateFunc: validate.IotSecuritySolutionName,
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"location": location.Schema(),

			"display_name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"iothub_ids": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: iothubValidate.IotHubID,
				},
				Set: set.HashStringIgnoreCase,
			},

			"log_analytics_workspace_id": {
				Type:             schema.TypeString,
				Optional:         true,
				ValidateFunc:     loganalyticsValidate.LogAnalyticsWorkspaceID,
				DiffSuppressFunc: suppress.CaseDifference,
			},

			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},

			"log_unmasked_ips_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},

			"events_to_export": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
					ValidateFunc: validation.StringInSlice([]string{
						string(security.RawEvents),
					}, false),
				},
			},

			"recommendations_enabled": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"acr_authentication": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  true,
						},

						"agent_send_unutilized_msg": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  true,
						},

						"baseline": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  true,
						},

						"edge_hub_mem_optimize": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  true,
						},

						"edge_logging_option": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  true,
						},

						"inconsistent_module_settings": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  true,
						},

						"install_agent": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  true,
						},

						"ip_filter_deny_all": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  true,
						},

						"ip_filter_permissive_rule": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  true,
						},

						"open_ports": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  true,
						},

						"permissive_firewall_policy": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  true,
						},

						"permissive_input_firewall_rules": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  true,
						},

						"permissive_output_firewall_rules": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  true,
						},

						"privileged_docker_options": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  true,
						},

						"shared_credentials": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  true,
						},

						"vulnerable_tls_cipher_suite": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  true,
						},
					},
				},
			},

			"user_defined_resource": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"query_for_resources": {
							Type:     schema.TypeString,
							Required: true,
						},

						"query_subscription_ids": {
							Type:     schema.TypeSet,
							Required: true,
							Elem: &schema.Schema{
								Type:         schema.TypeString,
								ValidateFunc: validation.IsUUID,
							},
						},
					},
				},
			},

			"tags": tags.Schema(),
		},
	}
}

func resourceIotSecuritySolutionCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).SecurityCenter.IotSecuritySolutionClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	location := location.Normalize(d.Get("location").(string))

	resourceId := parse.NewIotSecuritySolutionID(subscriptionId, resourceGroup, name).ID()
	if d.IsNewResource() {
		existing, err := client.Get(ctx, resourceGroup, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing Security Center Iot Security Solution %q (Resource Group %q): %+v", name, resourceGroup, err)
			}
		}
		if !utils.ResponseWasNotFound(existing.Response) {
			return tf.ImportAsExistsError("azurerm_iot_security_solution", resourceId)
		}
	}

	status := security.SolutionStatusDisabled
	if d.Get("enabled").(bool) {
		status = security.SolutionStatusEnabled
	}

	unmaskedIPLoggingStatus := security.UnmaskedIPLoggingStatusDisabled
	if d.Get("log_unmasked_ips_enabled").(bool) {
		unmaskedIPLoggingStatus = security.UnmaskedIPLoggingStatusEnabled
	}
	solution := security.IoTSecuritySolutionModel{
		Location: utils.String(location),
		IoTSecuritySolutionProperties: &security.IoTSecuritySolutionProperties{
			DisplayName:                  utils.String(d.Get("display_name").(string)),
			Status:                       status,
			Export:                       expandIotSecuritySolutionExport(d.Get("events_to_export").(*schema.Set).List()),
			IotHubs:                      utils.ExpandStringSlice(d.Get("iothub_ids").(*schema.Set).List()),
			UserDefinedResources:         expandIotSecuritySolutionUserDefinedResources(d.Get("user_defined_resource").([]interface{})),
			RecommendationsConfiguration: expandIotSecuritySolutionRecommendation(d.Get("recommendations_enabled").([]interface{})),
			UnmaskedIPLoggingStatus:      unmaskedIPLoggingStatus,
		},
		Tags: tags.Expand(d.Get("tags").(map[string]interface{})),
	}

	logAnalyticsWorkspaceId := d.Get("log_analytics_workspace_id").(string)
	if logAnalyticsWorkspaceId != "" {
		solution.IoTSecuritySolutionProperties.Workspace = utils.String(logAnalyticsWorkspaceId)
	}

	if _, err := client.CreateOrUpdate(ctx, resourceGroup, name, solution); err != nil {
		return fmt.Errorf("creating/updating Security Center Iot Security Solution %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	d.SetId(resourceId)
	return resourceIotSecuritySolutionRead(d, meta)
}

func resourceIotSecuritySolutionRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).SecurityCenter.IotSecuritySolutionClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.IotSecuritySolutionID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] Security Center Iot Security Solution %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("reading Security Center Iot Security Solution %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	d.Set("name", id.Name)
	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("location", location.NormalizeNilable(resp.Location))
	if prop := resp.IoTSecuritySolutionProperties; prop != nil {
		d.Set("enabled", prop.Status == security.SolutionStatusEnabled)
		d.Set("display_name", prop.DisplayName)
		d.Set("iothub_ids", utils.FlattenStringSlice(prop.IotHubs))
		d.Set("log_analytics_workspace_id", prop.Workspace)
		d.Set("log_unmasked_ips_enabled", prop.UnmaskedIPLoggingStatus == security.UnmaskedIPLoggingStatusEnabled)
		if err := d.Set("events_to_export", flattenIotSecuritySolutionExport(prop.Export)); err != nil {
			return fmt.Errorf("setting `events_to_export`: %s", err)
		}
		if err := d.Set("recommendations_enabled", flattenIotSecuritySolutionRecommendation(prop.RecommendationsConfiguration)); err != nil {
			return fmt.Errorf("setting `recommendations_enabled`: %s", err)
		}
		if err := d.Set("user_defined_resource", flattenIotSecuritySolutionUserDefinedResources(prop.UserDefinedResources)); err != nil {
			return fmt.Errorf("setting `user_defined_resource`: %s", err)
		}
	}

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceIotSecuritySolutionDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).SecurityCenter.IotSecuritySolutionClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.IotSecuritySolutionID(d.Id())
	if err != nil {
		return err
	}

	if _, err := client.Delete(ctx, id.ResourceGroup, id.Name); err != nil {
		return fmt.Errorf("deleting Security Center Iot Security Solution %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	return nil
}

func expandIotSecuritySolutionExport(input []interface{}) *[]security.ExportData {
	if len(input) == 0 || input[0] == nil {
		return nil
	}
	result := make([]security.ExportData, 0)
	for _, item := range input {
		result = append(result, security.ExportData(item.(string)))
	}
	return &result
}

func expandIotSecuritySolutionUserDefinedResources(input []interface{}) *security.UserDefinedResourcesProperties {
	if len(input) == 0 || input[0] == nil {
		return nil
	}
	v := input[0].(map[string]interface{})
	return &security.UserDefinedResourcesProperties{
		Query:              utils.String(v["query_for_resources"].(string)),
		QuerySubscriptions: utils.ExpandStringSlice(v["query_subscription_ids"].(*schema.Set).List()),
	}
}

func expandIotSecuritySolutionRecommendation(input []interface{}) *[]security.RecommendationConfigurationProperties {
	if len(input) == 0 || input[0] == nil {
		return nil
	}
	result := make([]security.RecommendationConfigurationProperties, 0)
	v := input[0].(map[string]interface{})
	for k, item := range getRecommendationSchemaMap() {
		status := security.Disabled
		if v[item].(bool) {
			status = security.Enabled
		}
		result = append(result, security.RecommendationConfigurationProperties{
			RecommendationType: k,
			Status:             status,
		})
	}
	return &result
}

func flattenIotSecuritySolutionExport(input *[]security.ExportData) []interface{} {
	result := make([]interface{}, 0)
	if input != nil {
		for _, item := range *input {
			result = append(result, string(item))
		}
	}
	return result
}

func flattenIotSecuritySolutionUserDefinedResources(input *security.UserDefinedResourcesProperties) []interface{} {
	if input == nil {
		return []interface{}{}
	}
	if input.QuerySubscriptions == nil && input.Query == nil {
		return []interface{}{}
	}
	query := ""
	if input.Query != nil {
		query = *input.Query
	}
	return []interface{}{
		map[string]interface{}{
			"query_for_resources":    query,
			"query_subscription_ids": utils.FlattenStringSlice(input.QuerySubscriptions),
		},
	}
}

func flattenIotSecuritySolutionRecommendation(input *[]security.RecommendationConfigurationProperties) []interface{} {
	if input == nil || len(*input) == 0 {
		return []interface{}{}
	}
	result := make(map[string]interface{})
	schemaMap := getRecommendationSchemaMap()
	for _, item := range *input {
		if v, ok := schemaMap[item.RecommendationType]; ok {
			result[v] = item.Status == security.Enabled
		}
	}
	return []interface{}{result}
}

func getRecommendationSchemaMap() map[security.RecommendationType]string {
	return map[security.RecommendationType]string{
		security.IoTACRAuthentication:             "acr_authentication",
		security.IoTAgentSendsUnutilizedMessages:  "agent_send_unutilized_msg",
		security.IoTBaseline:                      "baseline",
		security.IoTEdgeHubMemOptimize:            "edge_hub_mem_optimize",
		security.IoTEdgeLoggingOptions:            "edge_logging_option",
		security.IoTInconsistentModuleSettings:    "inconsistent_module_settings",
		security.IoTInstallAgent:                  "install_agent",
		security.IoTIPFilterDenyAll:               "ip_filter_deny_all",
		security.IoTIPFilterPermissiveRule:        "ip_filter_permissive_rule",
		security.IoTOpenPorts:                     "open_ports",
		security.IoTPermissiveFirewallPolicy:      "permissive_firewall_policy",
		security.IoTPermissiveInputFirewallRules:  "permissive_input_firewall_rules",
		security.IoTPermissiveOutputFirewallRules: "permissive_output_firewall_rules",
		security.IoTPrivilegedDockerOptions:       "privileged_docker_options",
		security.IoTSharedCredentials:             "shared_credentials",
		security.IoTVulnerableTLSCipherSuite:      "vulnerable_tls_cipher_suite",
	}
}
