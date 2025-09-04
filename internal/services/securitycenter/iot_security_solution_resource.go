// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package securitycenter

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/security/mgmt/v3.0/security" // nolint: staticcheck
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/operationalinsights/2020-08-01/workspaces"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	iothubValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/iothub/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/securitycenter/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/securitycenter/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/securitycenter/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/set"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceIotSecuritySolution() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceIotSecuritySolutionCreateUpdate,
		Read:   resourceIotSecuritySolutionRead,
		Update: resourceIotSecuritySolutionCreateUpdate,
		Delete: resourceIotSecuritySolutionDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.IotSecuritySolutionID(id)
			return err
		}),

		SchemaVersion: 1,
		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.SecurityCenterIotSecuritySolutionV0ToV1{},
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.IotSecuritySolutionName,
			},

			"resource_group_name": commonschema.ResourceGroupName(),

			"location": commonschema.Location(),

			"display_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"iothub_ids": {
				Type:     pluginsdk.TypeSet,
				Required: true,
				Elem: &pluginsdk.Schema{
					Type:         pluginsdk.TypeString,
					ValidateFunc: iothubValidate.IotHubID,
				},
				Set: set.HashStringIgnoreCase,
			},

			"additional_workspace": {
				Type:     pluginsdk.TypeSet,
				Optional: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"data_types": {
							Type:     pluginsdk.TypeSet,
							Required: true,
							Elem: &pluginsdk.Schema{
								Type: pluginsdk.TypeString,
								ValidateFunc: validation.StringInSlice([]string{
									string(security.Alerts),
									string(security.RawEvents),
								}, false),
							},
						},

						"workspace_id": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: workspaces.ValidateWorkspaceID,
						},
					},
				},
			},

			"disabled_data_sources": {
				Type:     pluginsdk.TypeSet,
				Optional: true,
				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
					ValidateFunc: validation.StringInSlice([]string{
						string(security.TwinData),
					}, false),
				},
			},

			"log_analytics_workspace_id": {
				Type:             pluginsdk.TypeString,
				Optional:         true,
				ValidateFunc:     workspaces.ValidateWorkspaceID,
				DiffSuppressFunc: suppress.CaseDifference,
			},

			"enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  true,
			},

			"log_unmasked_ips_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  false,
			},

			"events_to_export": {
				Type:     pluginsdk.TypeSet,
				Optional: true,
				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
					ValidateFunc: validation.StringInSlice([]string{
						string(security.RawEvents),
					}, false),
				},
			},

			"recommendations_enabled": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"acr_authentication": {
							Type:     pluginsdk.TypeBool,
							Optional: true,
							Default:  true,
						},

						"agent_send_unutilized_msg": {
							Type:     pluginsdk.TypeBool,
							Optional: true,
							Default:  true,
						},

						"baseline": {
							Type:     pluginsdk.TypeBool,
							Optional: true,
							Default:  true,
						},

						"edge_hub_mem_optimize": {
							Type:     pluginsdk.TypeBool,
							Optional: true,
							Default:  true,
						},

						"edge_logging_option": {
							Type:     pluginsdk.TypeBool,
							Optional: true,
							Default:  true,
						},

						"inconsistent_module_settings": {
							Type:     pluginsdk.TypeBool,
							Optional: true,
							Default:  true,
						},

						"install_agent": {
							Type:     pluginsdk.TypeBool,
							Optional: true,
							Default:  true,
						},

						"ip_filter_deny_all": {
							Type:     pluginsdk.TypeBool,
							Optional: true,
							Default:  true,
						},

						"ip_filter_permissive_rule": {
							Type:     pluginsdk.TypeBool,
							Optional: true,
							Default:  true,
						},

						"open_ports": {
							Type:     pluginsdk.TypeBool,
							Optional: true,
							Default:  true,
						},

						"permissive_firewall_policy": {
							Type:     pluginsdk.TypeBool,
							Optional: true,
							Default:  true,
						},

						"permissive_input_firewall_rules": {
							Type:     pluginsdk.TypeBool,
							Optional: true,
							Default:  true,
						},

						"permissive_output_firewall_rules": {
							Type:     pluginsdk.TypeBool,
							Optional: true,
							Default:  true,
						},

						"privileged_docker_options": {
							Type:     pluginsdk.TypeBool,
							Optional: true,
							Default:  true,
						},

						"shared_credentials": {
							Type:     pluginsdk.TypeBool,
							Optional: true,
							Default:  true,
						},

						"vulnerable_tls_cipher_suite": {
							Type:     pluginsdk.TypeBool,
							Optional: true,
							Default:  true,
						},
					},
				},
			},

			"query_for_resources": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				Computed: true,
			},

			"query_subscription_ids": {
				Type:     pluginsdk.TypeSet,
				Optional: true,
				Computed: true,
				Elem: &pluginsdk.Schema{
					Type:         pluginsdk.TypeString,
					ValidateFunc: validation.IsUUID,
				},
			},

			"tags": tags.Schema(),
		},
	}
}

func resourceIotSecuritySolutionCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
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
			Export:                       expandIotSecuritySolutionExport(d.Get("events_to_export").(*pluginsdk.Set).List()),
			IotHubs:                      utils.ExpandStringSlice(d.Get("iothub_ids").(*pluginsdk.Set).List()),
			RecommendationsConfiguration: expandIotSecuritySolutionRecommendation(d.Get("recommendations_enabled").([]interface{})),
			UnmaskedIPLoggingStatus:      unmaskedIPLoggingStatus,
		},
		Tags: tags.Expand(d.Get("tags").(map[string]interface{})),
	}

	if v, ok := d.GetOk("additional_workspace"); ok {
		solution.IoTSecuritySolutionProperties.AdditionalWorkspaces = expandIotSecuritySolutionAdditionalWorkspace(v.(*pluginsdk.Set).List())
	}

	if v, ok := d.GetOk("disabled_data_sources"); ok {
		solution.IoTSecuritySolutionProperties.DisabledDataSources = expandIotSecuritySolutionDisabledDataSources(v.(*pluginsdk.Set).List())
	}

	logAnalyticsWorkspaceId := d.Get("log_analytics_workspace_id").(string)
	if logAnalyticsWorkspaceId != "" {
		solution.IoTSecuritySolutionProperties.Workspace = utils.String(logAnalyticsWorkspaceId)
	}

	query := d.Get("query_for_resources").(string)
	querySubscriptions := d.Get("query_subscription_ids").(*pluginsdk.Set).List()
	if query != "" || len(querySubscriptions) > 0 {
		if query != "" && len(querySubscriptions) > 0 {
			solution.UserDefinedResources = &security.UserDefinedResourcesProperties{
				Query:              utils.String(query),
				QuerySubscriptions: utils.ExpandStringSlice(querySubscriptions),
			}
		} else {
			return fmt.Errorf("`query_for_resources` and `query_subscription_ids` must be set togetther")
		}
	}

	if _, err := client.CreateOrUpdate(ctx, resourceGroup, name, solution); err != nil {
		return fmt.Errorf("creating/updating Security Center Iot Security Solution %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	d.SetId(resourceId)
	return resourceIotSecuritySolutionRead(d, meta)
}

func resourceIotSecuritySolutionRead(d *pluginsdk.ResourceData, meta interface{}) error {
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
		if prop.UserDefinedResources != nil {
			d.Set("query_for_resources", prop.UserDefinedResources.Query)
			d.Set("query_subscription_ids", utils.FlattenStringSlice(prop.UserDefinedResources.QuerySubscriptions))
		}
		if err := d.Set("additional_workspace", flattenIotSecuritySolutionAdditionalWorkspace(prop.AdditionalWorkspaces)); err != nil {
			return fmt.Errorf("setting `additional_workspace`: %+v", err)
		}
		if err := d.Set("disabled_data_sources", flattenIotSecuritySolutionDisabledDataSources(prop.DisabledDataSources)); err != nil {
			return fmt.Errorf("setting `disabled_data_sources`: %s", err)
		}
	}

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceIotSecuritySolutionDelete(d *pluginsdk.ResourceData, meta interface{}) error {
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

func expandIotSecuritySolutionAdditionalWorkspace(input []interface{}) *[]security.AdditionalWorkspacesProperties {
	results := make([]security.AdditionalWorkspacesProperties, 0)

	for _, item := range input {
		v := item.(map[string]interface{})

		dataTypes := make([]security.AdditionalWorkspaceDataType, 0)
		for _, item := range *(utils.ExpandStringSlice(v["data_types"].(*pluginsdk.Set).List())) {
			dataTypes = append(dataTypes, (security.AdditionalWorkspaceDataType)(item))
		}

		results = append(results, security.AdditionalWorkspacesProperties{
			Workspace: utils.String(v["workspace_id"].(string)),
			Type:      security.Sentinel,
			DataTypes: &dataTypes,
		})
	}

	return &results
}

func expandIotSecuritySolutionDisabledDataSources(input []interface{}) *[]security.DataSource {
	if len(input) == 0 {
		return nil
	}

	disabledDataSources := make([]security.DataSource, 0)
	for _, item := range *(utils.ExpandStringSlice(input)) {
		disabledDataSources = append(disabledDataSources, (security.DataSource)(item))
	}

	return &disabledDataSources
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

func flattenIotSecuritySolutionAdditionalWorkspace(input *[]security.AdditionalWorkspacesProperties) []interface{} {
	results := make([]interface{}, 0)
	if input == nil {
		return results
	}

	for _, item := range *input {
		rawDataTypes := make([]string, 0)
		for _, item := range *item.DataTypes {
			rawDataTypes = append(rawDataTypes, (string)(item))
		}
		dataTypes := utils.FlattenStringSlice(&rawDataTypes)

		var workspaceId string
		if item.Workspace != nil {
			workspaceId = *item.Workspace
		}

		results = append(results, map[string]interface{}{
			"data_types":   dataTypes,
			"workspace_id": workspaceId,
		})
	}

	return results
}

func flattenIotSecuritySolutionDisabledDataSources(input *[]security.DataSource) []interface{} {
	if input == nil || len(*input) == 0 {
		return nil
	}

	results := make([]string, 0)
	for _, v := range *input {
		results = append(results, (string)(v))
	}

	return utils.FlattenStringSlice(&results)
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
