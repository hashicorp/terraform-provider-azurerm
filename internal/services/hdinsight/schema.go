// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package hdinsight

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-sdk/resource-manager/hdinsight/2021-06-01/clusters"
	"github.com/hashicorp/go-azure-sdk/resource-manager/hdinsight/2021-06-01/extensions"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	azValidate "github.com/hashicorp/terraform-provider-azurerm/helpers/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/hdinsight/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/parse"
	keyVault "github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func SchemaHDInsightName() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:         pluginsdk.TypeString,
		Required:     true,
		ForceNew:     true,
		ValidateFunc: validate.HDInsightName,
	}
}

func SchemaHDInsightTier() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeString,
		Required: true,
		ForceNew: true,
		ValidateFunc: validation.StringInSlice([]string{
			string(clusters.TierStandard),
			string(clusters.TierPremium),
		}, false),
	}
}

func SchemaHDInsightTls() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeString,
		Optional: true,
		ForceNew: true,
		ValidateFunc: validation.StringInSlice([]string{
			"1.0",
			"1.1",
			"1.2",
		}, false),
	}
}

func SchemaHDInsightClusterVersion() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:             pluginsdk.TypeString,
		Required:         true,
		ForceNew:         true,
		ValidateFunc:     validate.HDInsightClusterVersion,
		DiffSuppressFunc: hdinsightClusterVersionDiffSuppressFunc,
	}
}

func hdinsightClusterVersionDiffSuppressFunc(_, old, new string, _ *pluginsdk.ResourceData) bool {
	// `3.6` gets converted to `3.6.1000.67`; so let's just compare major/minor if possible
	o := strings.Split(old, ".")
	n := strings.Split(new, ".")

	if len(o) >= 2 && len(n) >= 2 {
		oldMajor := o[0]
		oldMinor := o[1]
		newMajor := n[0]
		newMinor := n[1]

		return oldMajor == newMajor && oldMinor == newMinor
	}

	return false
}

func SchemaHDInsightsGateway() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Required: true,
		MaxItems: 1,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				// NOTE: these are Required since if these aren't present you get a `500 bad request`
				"username": {
					Type:     pluginsdk.TypeString,
					Required: true,
					ForceNew: true,
				},
				"password": {
					Type:      pluginsdk.TypeString,
					Required:  true,
					Sensitive: true,
					// Azure returns the key as *****. We'll suppress that here.
					DiffSuppressFunc: func(k, old, new string, d *pluginsdk.ResourceData) bool {
						return (new == d.Get(k).(string)) && (old == "*****")
					},
				},
			},
		},
	}
}

func SchemaHDInsightsComputeIsolation() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		MaxItems: 1,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*schema.Schema{
				"compute_isolation_enabled": {
					Type:     pluginsdk.TypeBool,
					Optional: true,
					Default:  false,
				},

				"host_sku": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},
			},
		},
	}
}

func SchemaHDInsightsExternalMetastore() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		MaxItems: 1,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"server": {
					Type:     pluginsdk.TypeString,
					Required: true,
					ForceNew: true,
				},
				"database_name": {
					Type:     pluginsdk.TypeString,
					Required: true,
					ForceNew: true,
				},
				"username": {
					Type:     pluginsdk.TypeString,
					Required: true,
					ForceNew: true,
				},
				"password": {
					Type:      pluginsdk.TypeString,
					Required:  true,
					ForceNew:  true,
					Sensitive: true,
					// Azure returns the key as *****. We'll suppress that here.
					DiffSuppressFunc: func(k, old, new string, d *pluginsdk.ResourceData) bool {
						return (new == d.Get(k).(string)) && (old == "*****")
					},
				},
			},
		},
	}
}

func SchemaHDInsightsExternalMetastores() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		MaxItems: 1,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"hive": SchemaHDInsightsExternalMetastore(),

				"oozie": SchemaHDInsightsExternalMetastore(),

				"ambari": SchemaHDInsightsExternalMetastore(),
			},
		},
	}
}

func SchemaHDInsightsMonitor() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		MaxItems: 1,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"log_analytics_workspace_id": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ValidateFunc: validation.IsUUID,
				},
				"primary_key": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					Sensitive:    true,
					ValidateFunc: validation.StringIsNotEmpty,
					// Azure doesn't return the key
					DiffSuppressFunc: func(k, old, new string, d *pluginsdk.ResourceData) bool {
						return (new == d.Get(k).(string)) && (old == "*****")
					},
				},
			},
		},
	}
}

func SchemaHDInsightsExtension() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		MaxItems: 1,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"log_analytics_workspace_id": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ValidateFunc: validation.IsUUID,
				},
				"primary_key": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					Sensitive:    true,
					ValidateFunc: validation.StringIsNotEmpty,
					DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
						return (new == d.Get(k).(string)) && (old == "*****")
					},
				},
			},
		},
	}
}

func SchemaHDInsightsNetwork() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		MaxItems: 1,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"connection_direction": {
					Type:     pluginsdk.TypeString,
					Optional: true,
					ForceNew: true,
					Default:  string(clusters.ResourceProviderConnectionInbound),
					ValidateFunc: validation.StringInSlice([]string{
						string(clusters.ResourceProviderConnectionInbound),
						string(clusters.ResourceProviderConnectionOutbound),
					}, false),
				},

				"private_link_enabled": {
					Type:     pluginsdk.TypeBool,
					Optional: true,
					ForceNew: true,
					Default:  false,
				},
			},
		},
	}
}

func SchemaHDInsightsSecurityProfile() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		ForceNew: true,
		MaxItems: 1,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"aadds_resource_id": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ForceNew:     true,
					ValidateFunc: azure.ValidateResourceID,
				},

				"domain_name": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ForceNew:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},

				"domain_username": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ForceNew:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},

				"domain_user_password": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ForceNew:     true,
					Sensitive:    true,
					ValidateFunc: validation.StringIsNotEmpty,
				},

				"ldaps_urls": {
					Type:     pluginsdk.TypeSet,
					Required: true,
					ForceNew: true,
					Elem: &pluginsdk.Schema{
						Type:         pluginsdk.TypeString,
						ValidateFunc: validate.HDInsightClusterLdapsUrls,
					},
				},

				"msi_resource_id": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ForceNew:     true,
					ValidateFunc: commonids.ValidateUserAssignedIdentityID,
				},

				"cluster_users_group_dns": {
					Type:     pluginsdk.TypeSet,
					Optional: true,
					ForceNew: true,
					Elem: &pluginsdk.Schema{
						Type:         pluginsdk.TypeString,
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},
			},
		},
	}
}

func SchemaHDInsightsScriptActions() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		MinItems: 1,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"name": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},

				"uri": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ValidateFunc: validation.IsURLWithHTTPorHTTPS,
				},

				"parameters": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},
			},
		},
	}
}

func SchemaHDInsightsHttpsEndpoints() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"access_modes": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					Elem: &pluginsdk.Schema{
						Type:         pluginsdk.TypeString,
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},

				"destination_port": {
					Type:         pluginsdk.TypeInt,
					Optional:     true,
					ValidateFunc: azValidate.PortNumber,
				},

				"disable_gateway_auth": {
					Type:     pluginsdk.TypeBool,
					Optional: true,
				},

				"private_ip_address": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: validation.IsIPAddress,
				},

				"sub_domain_suffix": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},
			},
		},
	}
}

type HttpEndpointModel struct {
	AccessModes        []string `tfschema:"access_modes"`
	DestinationPort    int32    `tfschema:"destination_port"`
	DisableGatewayAuth bool     `tfschema:"disable_gateway_auth"`
	PrivateIpAddress   string   `tfschema:"private_ip_address"`
	SubDomainSuffix    string   `tfschema:"sub_domain_suffix"`
}

func ExpandHDInsightsRolesScriptActions(input []interface{}) *[]clusters.ScriptAction {
	if len(input) == 0 {
		return nil
	}

	scriptActions := make([]clusters.ScriptAction, 0)

	for _, vs := range input {
		v := vs.(map[string]interface{})

		scriptActions = append(scriptActions, clusters.ScriptAction{
			Name:       v["name"].(string),
			Uri:        v["uri"].(string),
			Parameters: v["parameters"].(string),
		})
	}

	return &scriptActions
}

func ExpandHDInsightComputeIsolationProperties(input []interface{}) *clusters.ComputeIsolationProperties {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	v := input[0].(map[string]interface{})
	enableComputeIsolation := v["compute_isolation_enabled"].(bool)
	hostSku := v["host_sku"].(string)

	return &clusters.ComputeIsolationProperties{
		EnableComputeIsolation: &enableComputeIsolation,
		HostSku:                &hostSku,
	}
}

func ExpandHDInsightsConfigurations(input []interface{}) map[string]interface{} {
	vs := input[0].(map[string]interface{})

	// NOTE: Admin username must be different from SSH Username
	enabled := true
	username := vs["username"].(string)
	password := vs["password"].(string)

	return map[string]interface{}{
		"gateway": map[string]interface{}{
			"restAuthCredential.isEnabled": enabled,
			"restAuthCredential.username":  username,
			"restAuthCredential.password":  password,
		},
	}
}

func ExpandHDInsightsHiveMetastore(input []interface{}) map[string]interface{} {
	if len(input) == 0 {
		return nil
	}
	vs := input[0].(map[string]interface{})

	server := vs["server"].(string)
	database := vs["database_name"].(string)
	username := vs["username"].(string)
	password := vs["password"].(string)

	return map[string]interface{}{
		"hive-site": map[string]interface{}{
			"javax.jdo.option.ConnectionDriverName": "com.microsoft.sqlserver.jdbc.SQLServerDriver",
			"javax.jdo.option.ConnectionURL":        fmt.Sprintf("jdbc:sqlserver://%s;database=%s;encrypt=true;trustServerCertificate=true;create=false;loginTimeout=300", server, database),
			"javax.jdo.option.ConnectionUserName":   username,
			"javax.jdo.option.ConnectionPassword":   password,
		},
		"hive-env": map[string]interface{}{
			"hive_database":                       "Existing MSSQL Server database with SQL authentication",
			"hive_database_name":                  database,
			"hive_database_type":                  "mssql",
			"hive_existing_mssql_server_database": database,
			"hive_existing_mssql_server_host":     server,
			"hive_hostname":                       server,
		},
	}
}

func ExpandHDInsightsOozieMetastore(input []interface{}) map[string]interface{} {
	if len(input) == 0 {
		return nil
	}
	vs := input[0].(map[string]interface{})

	server := vs["server"].(string)
	database := vs["database_name"].(string)
	username := vs["username"].(string)
	password := vs["password"].(string)

	return map[string]interface{}{
		"oozie-site": map[string]interface{}{
			"oozie.service.JPAService.jdbc.driver":   "com.microsoft.sqlserver.jdbc.SQLServerDriver",
			"oozie.service.JPAService.jdbc.url":      fmt.Sprintf("jdbc:sqlserver://%s;database=%s;encrypt=true;trustServerCertificate=true;create=false;loginTimeout=300", server, database),
			"oozie.service.JPAService.jdbc.username": username,
			"oozie.service.JPAService.jdbc.password": password,
			"oozie.db.pluginsdk.name":                "oozie",
			"oozie.db.schema.name":                   "oozie",
		},
		"oozie-env": map[string]interface{}{
			"oozie_database":                       "Existing MSSQL Server database with SQL authentication",
			"oozie_database_name":                  database,
			"oozie_database_type":                  "mssql",
			"oozie_existing_mssql_server_database": database,
			"oozie_existing_mssql_server_host":     server,
			"oozie_hostname":                       server,
		},
	}
}

func ExpandHDInsightsAmbariMetastore(input []interface{}) map[string]interface{} {
	if len(input) == 0 {
		return nil
	}
	vs := input[0].(map[string]interface{})

	server := vs["server"].(string)
	database := vs["database_name"].(string)
	username := vs["username"].(string)
	password := vs["password"].(string)

	return map[string]interface{}{
		"ambari-conf": map[string]interface{}{
			"database-server":        server,
			"database-name":          database,
			"database-user-name":     username,
			"database-user-password": password,
		},
	}
}

func ExpandHDInsightsMonitor(input []interface{}) extensions.ClusterMonitoringRequest {
	vs := input[0].(map[string]interface{})

	return extensions.ClusterMonitoringRequest{
		WorkspaceId: pointer.To(vs["log_analytics_workspace_id"].(string)),
		PrimaryKey:  pointer.To(vs["primary_key"].(string)),
	}
}

func ExpandHDInsightsNetwork(input []interface{}) *clusters.NetworkProperties {
	if len(input) == 0 {
		return nil
	}

	vs := input[0].(map[string]interface{})

	connDir := clusters.ResourceProviderConnectionOutbound
	if v, exists := vs["connection_direction"]; exists && v != string(clusters.ResourceProviderConnectionOutbound) {
		connDir = clusters.ResourceProviderConnectionInbound
	}

	privateLink := clusters.PrivateLinkDisabled
	if v, exists := vs["private_link_enabled"]; exists && v != false {
		privateLink = clusters.PrivateLinkEnabled
	}

	return &clusters.NetworkProperties{
		ResourceProviderConnection: pointer.To(connDir),
		PrivateLink:                pointer.To(privateLink),
	}
}

func ExpandHDInsightPrivateLinkConfigurations(input []interface{}) *[]clusters.PrivateLinkConfiguration {
	if len(input) == 0 {
		return nil
	}

	configs := make([]clusters.PrivateLinkConfiguration, 0)

	for _, vs := range input {
		v := vs.(map[string]interface{})

		configs = append(configs, clusters.PrivateLinkConfiguration{
			Name:       v["name"].(string),
			Properties: ExpandHDInsightPrivateLinkConfigurationProperties(input),
		})
	}

	return pointer.To(configs)
}

func ExpandHDInsightPrivateLinkConfigurationProperties(input []interface{}) clusters.PrivateLinkConfigurationProperties {
	v := input[0].(map[string]interface{})

	return clusters.PrivateLinkConfigurationProperties{
		GroupId:          v["group_id"].(string),
		IPConfigurations: ExpandHDInsightPrivateLinkConfigurationIpConfiguration(v["ip_configuration"].([]interface{})),
	}
}

func ExpandHDInsightPrivateLinkConfigurationIpConfiguration(input []interface{}) []clusters.IPConfiguration {
	ipConfigs := make([]clusters.IPConfiguration, 0)

	for _, vs := range input {
		v := vs.(map[string]interface{})

		ipConfigs = append(ipConfigs, clusters.IPConfiguration{
			Name:       v["name"].(string),
			Properties: ExpandHDInsightPrivateLinkConfigurationIpConfigurationProperties(input),
		})
	}

	return ipConfigs
}

func ExpandHDInsightPrivateLinkConfigurationIpConfigurationProperties(input []interface{}) *clusters.IPConfigurationProperties {
	v := input[0].(map[string]interface{})

	props := clusters.IPConfigurationProperties{
		Primary:                   pointer.To(v["primary"].(bool)),
		PrivateIPAllocationMethod: pointer.To(clusters.PrivateIPAllocationMethod(v["private_ip_allocation_method"].(string))),
	}
	if v["private_ip_address"] != nil && v["private_ip_address"].(string) != "" {
		props.PrivateIPAddress = pointer.To(v["private_ip_address"].(string))
	}
	if v["subnet_id"] != nil && v["subnet_id"].(string) != "" {
		props.Subnet = pointer.To(clusters.ResourceId{Id: pointer.To(v["subnet_id"].(string))})
	}

	return pointer.To(props)
}

func flattenHDInsightComputeIsolationProperties(input *clusters.ComputeIsolationProperties) []interface{} {
	hostSku := ""
	enableComputeIsolation := false
	if input != nil {
		enableComputeIsolation = pointer.From(input.EnableComputeIsolation)
		hostSku = pointer.From(input.HostSku)
	}

	if !enableComputeIsolation {
		return make([]interface{}, 0)
	}

	return []interface{}{
		map[string]interface{}{
			"compute_isolation_enabled": enableComputeIsolation,
			"host_sku":                  hostSku,
		},
	}
}

func flattenHDInsightsNetwork(input *clusters.NetworkProperties) []interface{} {
	if input == nil {
		return make([]interface{}, 0)
	}

	connDir := string(clusters.ResourceProviderConnectionOutbound)
	if v := input.ResourceProviderConnection; v != nil {
		connDir = string(*v)
	}

	privateLink := false
	if v := input.PrivateLink; v != nil {
		privateLink = *v == clusters.PrivateLinkEnabled
	}

	return []interface{}{
		map[string]interface{}{
			"connection_direction": connDir,
			"private_link_enabled": privateLink,
		},
	}
}

func flattenHDInsightPrivateLinkConfigurations(input *[]clusters.PrivateLinkConfiguration) []interface{} {
	if input == nil {
		return make([]interface{}, 0)
	}

	v := pointer.From(input)[0]
	ipConfig := v.Properties.IPConfigurations[0]
	return []interface{}{
		map[string]interface{}{
			"name":             v.Name,
			"group_id":         v.Properties.GroupId,
			"ip_configuration": flattenHDInsightPrivateLinkConfigurationIpConfigurationProperties(&ipConfig),
		},
	}
}
func flattenHDInsightPrivateLinkConfigurationIpConfigurationProperties(input *clusters.IPConfiguration) []interface{} {
	if input == nil {
		return make([]interface{}, 0)
	}

	if input.Properties != nil {
		return []interface{}{
			map[string]interface{}{
				"name":                         input.Name,
				"primary":                      pointer.From(input.Properties.Primary),
				"private_ip_allocation_method": pointer.From(input.Properties.PrivateIPAllocationMethod),
				"private_ip_address":           pointer.From(input.Properties.PrivateIPAddress),
				"subnet_id":                    pointer.From(input.Properties.Subnet.Id),
			},
		}
	}

	return []interface{}{
		map[string]interface{}{
			"name": input.Name,
		},
	}
}

func FlattenHDInsightsConfigurations(input map[string]string, d *pluginsdk.ResourceData) []interface{} {
	username := ""
	if v, exists := input["restAuthCredential.username"]; exists {
		username = v
	}

	password := ""
	if v, exists := input["restAuthCredential.password"]; exists {
		password = v
	} else {
		password = d.Get("gateway.0.password").(string)
	}

	out := map[string]interface{}{
		"username": username,
		"password": password,
	}

	return []interface{}{out}
}

func FlattenHDInsightsHiveMetastore(env map[string]string, site map[string]string) []interface{} {
	server := ""
	if v, exists := env["hive_hostname"]; exists {
		server = v
	}

	database := ""
	if v, exists := env["hive_database_name"]; exists {
		database = v
	}

	username := ""
	if v, exists := site["javax.jdo.option.ConnectionUserName"]; exists {
		username = v
	}

	password := ""
	if v, exists := site["javax.jdo.option.ConnectionPassword"]; exists {
		password = v
	}

	if server != "" && database != "" {
		return []interface{}{
			map[string]interface{}{
				"server":        server,
				"database_name": database,
				"username":      username,
				"password":      password,
			},
		}
	}

	return []interface{}{}
}

func FlattenHDInsightsOozieMetastore(env map[string]string, site map[string]string) []interface{} {
	server := ""
	if v, exists := env["oozie_hostname"]; exists {
		server = v
	}

	database := ""
	if v, exists := env["oozie_database_name"]; exists {
		database = v
	}

	username := ""
	if v, exists := site["oozie.service.JPAService.jdbc.username"]; exists {
		username = v
	}

	password := ""
	if v, exists := site["oozie.service.JPAService.jdbc.password"]; exists {
		password = v
	}

	if server != "" && database != "" {
		return []interface{}{
			map[string]interface{}{
				"server":        server,
				"database_name": database,
				"username":      username,
				"password":      password,
			},
		}
	}

	return []interface{}{}
}

func FlattenHDInsightsAmbariMetastore(conf map[string]string) []interface{} {
	server := ""
	if v, exists := conf["database-server"]; exists {
		server = v
	}

	database := ""
	if v, exists := conf["database-name"]; exists {
		database = v
	}

	username := ""
	if v, exists := conf["database-user-name"]; exists {
		username = v
	}

	password := ""
	if v, exists := conf["database-user-password"]; exists {
		password = v
	}

	if server != "" && database != "" {
		return []interface{}{
			map[string]interface{}{
				"server":        server,
				"database_name": database,
				"username":      username,
				"password":      password,
			},
		}
	}

	return []interface{}{}
}

func SchemaHDInsightsStorageAccounts() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"storage_account_key": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ForceNew:     true,
					Sensitive:    true,
					ValidateFunc: validation.StringIsNotEmpty,
				},
				"storage_container_id": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ForceNew:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},
				// TODO: this should become `storage_account_id` in 4.0
				"storage_resource_id": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ForceNew:     true,
					ValidateFunc: commonids.ValidateStorageAccountID,
				},
				"is_default": {
					Type:     pluginsdk.TypeBool,
					Required: true,
					ForceNew: true,
				},
			},
		},
	}
}

func SchemaHDInsightsGen2StorageAccounts() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		// HDInsight doesn't seem to allow adding more than one gen2 cluster right now.
		MaxItems: 1,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				// TODO: this should become `storage_account_id` in 4.0
				"storage_resource_id": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ForceNew:     true,
					ValidateFunc: commonids.ValidateStorageAccountID,
				},
				"filesystem_id": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ForceNew:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},
				// TODO: this should become `user_assigned_identity_id` in 4.0
				"managed_identity_resource_id": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ForceNew:     true,
					ValidateFunc: azure.ValidateResourceID,
				},
				"is_default": {
					Type:     pluginsdk.TypeBool,
					Required: true,
					ForceNew: true,
				},
			},
		},
	}
}

func SchemaHDInsightsDiskEncryptionProperties() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"encryption_algorithm": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: validation.StringInSlice(clusters.PossibleValuesForJsonWebKeyEncryptionAlgorithm(), false),
				},

				"encryption_at_host_enabled": {
					Type:     pluginsdk.TypeBool,
					Optional: true,
				},

				"key_vault_managed_identity_id": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: commonids.ValidateUserAssignedIdentityID,
				},

				"key_vault_key_id": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: keyVault.NestedItemId,
				},
			},
		},
	}
}

func SchemaHDInsightPrivateLinkConfigurations() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		MaxItems: 1,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"name": {
					Type:     pluginsdk.TypeString,
					Required: true,
				},

				"group_id": {
					Type:     pluginsdk.TypeString,
					Required: true,
				},

				"ip_configuration": SchemaHDInsightPrivateLinkConfigurationIpConfiguration(),
			},
		},
	}
}

func SchemaHDInsightPrivateLinkConfigurationIpConfiguration() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Required: true,
		MaxItems: 1,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"name": {
					Type:     pluginsdk.TypeString,
					Required: true,
				},

				"primary": {
					Type:     pluginsdk.TypeBool,
					Optional: true,
				},

				"private_ip_address": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: validation.IsIPAddress,
				},

				"private_ip_allocation_method": {
					Type:     pluginsdk.TypeString,
					Optional: true,
					ValidateFunc: validation.StringInSlice([]string{
						string(clusters.PrivateIPAllocationMethodDynamic),
						string(clusters.PrivateIPAllocationMethodStatic),
					}, false),
				},

				"subnet_id": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: commonids.ValidateSubnetID,
				},
			},
		},
	}
}

func ExpandHDInsightsDiskEncryptionProperties(input []interface{}) (*clusters.DiskEncryptionProperties, error) {
	v := input[0].(map[string]interface{})

	encryptionAlgorithm := v["encryption_algorithm"].(string)
	encryptionAtHost := v["encryption_at_host_enabled"].(bool)
	keyVaultManagedIdentityId := v["key_vault_managed_identity_id"].(string)

	diskEncryptionProps := &clusters.DiskEncryptionProperties{
		EncryptionAlgorithm: pointer.To(clusters.JsonWebKeyEncryptionAlgorithm(encryptionAlgorithm)),
		EncryptionAtHost:    &encryptionAtHost,
		MsiResourceId:       &keyVaultManagedIdentityId,
	}

	if id, ok := v["key_vault_key_id"]; ok && id.(string) != "" {
		keyVaultKeyId, err := parse.ParseNestedItemID(id.(string))
		if err != nil {
			return nil, err
		}
		diskEncryptionProps.KeyName = &keyVaultKeyId.Name
		diskEncryptionProps.KeyVersion = &keyVaultKeyId.Version
		diskEncryptionProps.VaultUri = &keyVaultKeyId.KeyVaultBaseUrl
	}

	return diskEncryptionProps, nil
}

func flattenHDInsightsDiskEncryptionProperties(input *clusters.DiskEncryptionProperties) (*[]interface{}, error) {
	if input == nil {
		return pointer.To(make([]interface{}, 0)), nil
	}
	encryptionAlgorithm := ""
	if input.EncryptionAlgorithm != nil {
		encryptionAlgorithm = string(*input.EncryptionAlgorithm)
	}

	keyName := pointer.From(input.KeyName)
	keyVersion := pointer.From(input.KeyVersion)
	keyVaultKeyId := ""
	if (keyName != "" || keyVersion != "") && input.VaultUri != nil {
		keyVaultKeyIdRaw, err := parse.NewNestedItemID(*input.VaultUri, parse.NestedItemTypeKey, keyName, keyVersion)
		if err != nil {
			return nil, err
		}
		keyVaultKeyId = keyVaultKeyIdRaw.ID()
	}

	return &[]interface{}{
		map[string]interface{}{
			"encryption_algorithm":          encryptionAlgorithm,
			"encryption_at_host_enabled":    pointer.From(input.EncryptionAtHost),
			"key_vault_key_id":              keyVaultKeyId,
			"key_vault_managed_identity_id": pointer.From(input.MsiResourceId),
		},
	}, nil
}

// ExpandHDInsightsStorageAccounts returns an array of StorageAccount structs, as well as a ClusterIdentity
// populated with any managed identities required for accessing Data Lake Gen2 storage.
func ExpandHDInsightsStorageAccounts(storageAccounts []interface{}, gen2storageAccounts []interface{}) (*[]clusters.StorageAccount, *identity.SystemAndUserAssignedMap, error) {
	results := make([]clusters.StorageAccount, 0)

	var clusterIdentity *identity.SystemAndUserAssignedMap

	for _, vs := range storageAccounts {
		v := vs.(map[string]interface{})

		storageAccountKey := v["storage_account_key"].(string)
		storageContainerID := v["storage_container_id"].(string)
		storageResourceID := v["storage_resource_id"].(string)
		isDefault := v["is_default"].(bool)

		uri, err := url.Parse(storageContainerID)
		if err != nil {
			return nil, nil, fmt.Errorf("parsing %q: %s", storageContainerID, err)
		}

		result := clusters.StorageAccount{
			Name:       utils.String(uri.Host),
			ResourceId: utils.String(storageResourceID),
			Container:  utils.String(strings.TrimPrefix(uri.Path, "/")),
			Key:        utils.String(storageAccountKey),
			IsDefault:  utils.Bool(isDefault),
		}
		results = append(results, result)
	}

	for _, vs := range gen2storageAccounts {
		v := vs.(map[string]interface{})

		fileSystemID := v["filesystem_id"].(string)
		storageResourceID := v["storage_resource_id"].(string)
		managedIdentityResourceID := v["managed_identity_resource_id"].(string)

		isDefault := v["is_default"].(bool)

		uri, err := url.Parse(fileSystemID)
		if err != nil {
			return nil, nil, fmt.Errorf("parsing %q: %s", fileSystemID, err)
		}

		if clusterIdentity == nil {
			clusterIdentity = &identity.SystemAndUserAssignedMap{
				Type:        identity.TypeUserAssigned,
				IdentityIds: make(map[string]identity.UserAssignedIdentityDetails),
			}
		}

		clusterIdentity.IdentityIds[managedIdentityResourceID] = identity.UserAssignedIdentityDetails{
			// intentionally empty
		}

		result := clusters.StorageAccount{
			Name:          utils.String(uri.Host), // https://storageaccountname.dfs.core.windows.net/filesystemname -> storageaccountname.dfs.core.windows.net
			ResourceId:    utils.String(storageResourceID),
			FileSystem:    utils.String(uri.Path[1:]), // https://storageaccountname.dfs.core.windows.net/filesystemname -> filesystemname
			MsiResourceId: utils.String(managedIdentityResourceID),
			IsDefault:     utils.Bool(isDefault),
		}
		results = append(results, result)
	}

	return &results, clusterIdentity, nil
}

type HDInsightNodeDefinition struct {
	CanSpecifyInstanceCount  bool
	MinInstanceCount         int
	MaxInstanceCount         *int
	CanSpecifyDisks          bool
	MaxNumberOfDisksPerNode  *int
	FixedMinInstanceCount    *int64
	FixedTargetInstanceCount *int64
	CanAutoScaleByCapacity   bool
	CanAutoScaleOnSchedule   bool
	// todo remove in 4.0
	CanAutoScaleByCapacityDeprecated4PointOh bool
}

func SchemaHDInsightNodeDefinition(schemaLocation string, definition HDInsightNodeDefinition, required bool) *pluginsdk.Schema {
	result := map[string]*pluginsdk.Schema{
		"vm_size": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringInSlice(validate.NodeDefinitionVMSize, false),
		},
		"username": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},
		"password": {
			Type:      pluginsdk.TypeString,
			Optional:  true,
			ForceNew:  true,
			Sensitive: true,
		},
		"ssh_keys": {
			Type:     pluginsdk.TypeSet,
			Optional: true,
			ForceNew: true,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: validation.StringIsNotEmpty,
			},
			Set: pluginsdk.HashString,
			ConflictsWith: []string{
				fmt.Sprintf("%s.0.password", schemaLocation),
			},
		},

		"subnet_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: commonids.ValidateSubnetID,
		},

		"virtual_network_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: commonids.ValidateVirtualNetworkID,
		},

		"script_actions": SchemaHDInsightsScriptActions(),
	}

	if definition.CanSpecifyInstanceCount {
		countValidation := validation.IntAtLeast(definition.MinInstanceCount)
		if definition.MaxInstanceCount != nil {
			countValidation = validation.IntBetween(definition.MinInstanceCount, *definition.MaxInstanceCount)
		}

		result["target_instance_count"] = &pluginsdk.Schema{
			Type:         pluginsdk.TypeInt,
			Required:     true,
			ValidateFunc: countValidation,
		}

		if definition.CanAutoScaleByCapacity || definition.CanAutoScaleOnSchedule {
			autoScales := map[string]*pluginsdk.Schema{}

			if definition.CanAutoScaleByCapacity {
				autoScales["capacity"] = &pluginsdk.Schema{
					Type:     pluginsdk.TypeList,
					Optional: true,
					MaxItems: 1,
					ConflictsWith: []string{
						fmt.Sprintf("%s.0.autoscale.0.recurrence", schemaLocation),
					},
					Elem: &pluginsdk.Resource{
						Schema: map[string]*pluginsdk.Schema{
							"min_instance_count": {
								Type:         pluginsdk.TypeInt,
								Required:     true,
								ValidateFunc: countValidation,
							},
							"max_instance_count": {
								Type:         pluginsdk.TypeInt,
								Required:     true,
								ValidateFunc: countValidation,
							},
						},
					},
				}
				if definition.CanAutoScaleOnSchedule {
					autoScales["capacity"].ConflictsWith = []string{
						fmt.Sprintf("%s.0.autoscale.0.recurrence", schemaLocation),
					}
				}
			}
			// managing `azurerm_hdinsight_interactive_query_cluster` autoscaling through `capacity` doesn't work so we'll deprecate this portion of the schema for 4.0
			if definition.CanAutoScaleByCapacityDeprecated4PointOh {
				autoScales["capacity"] = &pluginsdk.Schema{
					Type:     pluginsdk.TypeList,
					Optional: true,
					MaxItems: 1,
					ConflictsWith: []string{
						fmt.Sprintf("%s.0.autoscale.0.recurrence", schemaLocation),
					},
					Deprecated: "HDInsight interactive query clusters can no longer be configured through `autoscale.0.capacity`. Use `autoscale.0.recurrence` instead.",
					Elem: &pluginsdk.Resource{
						Schema: map[string]*pluginsdk.Schema{
							"min_instance_count": {
								Type:         pluginsdk.TypeInt,
								Required:     true,
								ValidateFunc: countValidation,
							},
							"max_instance_count": {
								Type:         pluginsdk.TypeInt,
								Required:     true,
								ValidateFunc: countValidation,
							},
						},
					},
				}
				if definition.CanAutoScaleOnSchedule {
					autoScales["capacity"].ConflictsWith = []string{
						fmt.Sprintf("%s.0.autoscale.0.recurrence", schemaLocation),
					}
				}
			}
			if definition.CanAutoScaleOnSchedule {
				autoScales["recurrence"] = &pluginsdk.Schema{
					Type:     pluginsdk.TypeList,
					Optional: true,
					MaxItems: 1,
					Elem: &pluginsdk.Resource{
						Schema: map[string]*pluginsdk.Schema{
							"timezone": {
								Type:     pluginsdk.TypeString,
								Required: true,
							},
							"schedule": {
								Type:     pluginsdk.TypeList,
								Required: true,
								MinItems: 1,
								Elem: &pluginsdk.Resource{
									Schema: map[string]*pluginsdk.Schema{
										"time": {
											Type:     pluginsdk.TypeString,
											Required: true,
											ValidateFunc: validation.StringMatch(
												regexp.MustCompile("^([01][0-9]|[2][0-3]):([03][0])$"), // time must be on the hour or half past
												"Time of day must match the format HH:mm where HH is 00-23 and mm is 00 or 30",
											),
										},
										"days": {
											Type:     pluginsdk.TypeList,
											Required: true,
											Elem: &pluginsdk.Schema{
												Type: pluginsdk.TypeString,
												ValidateFunc: validation.StringInSlice([]string{
													string(clusters.DaysOfWeekMonday),
													string(clusters.DaysOfWeekTuesday),
													string(clusters.DaysOfWeekWednesday),
													string(clusters.DaysOfWeekThursday),
													string(clusters.DaysOfWeekFriday),
													string(clusters.DaysOfWeekSaturday),
													string(clusters.DaysOfWeekSunday),
												}, false),
											},
										},

										"target_instance_count": {
											Type:         pluginsdk.TypeInt,
											Required:     true,
											ValidateFunc: countValidation,
										},
									},
								},
							},
						},
					},
				}
				if definition.CanAutoScaleByCapacity {
					autoScales["recurrence"].ConflictsWith = []string{
						fmt.Sprintf("%s.0.autoscale.0.capacity", schemaLocation),
					}
				}
			}

			result["autoscale"] = &pluginsdk.Schema{
				Type:     pluginsdk.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: autoScales,
				},
			}
		}
	}

	if definition.CanSpecifyDisks {
		result["number_of_disks_per_node"] = &pluginsdk.Schema{
			Type:         pluginsdk.TypeInt,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.IntBetween(1, *definition.MaxNumberOfDisksPerNode),
		}
	}

	s := &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		MaxItems: 1,
		Required: required,
		Optional: !required,
		Elem: &pluginsdk.Resource{
			Schema: result,
		},
	}

	return s
}

func SchemaHDInsightNodeDefinitionKafka(schemaLocation string, definition HDInsightNodeDefinition, required bool) *pluginsdk.Schema {
	result := map[string]*pluginsdk.Schema{
		"vm_size": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringInSlice(validate.NodeDefinitionVMSize, false),
		},
		"username": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
		"password": {
			Type:      pluginsdk.TypeString,
			Optional:  true,
			ForceNew:  true,
			Sensitive: true,
		},
		"ssh_keys": {
			Type:     pluginsdk.TypeSet,
			Optional: true,
			ForceNew: true,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: validation.StringIsNotEmpty,
			},
			Set: pluginsdk.HashString,
			ConflictsWith: []string{
				fmt.Sprintf("%s.0.password", schemaLocation),
			},
		},

		"subnet_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: commonids.ValidateSubnetID,
		},

		"virtual_network_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: commonids.ValidateVirtualNetworkID,
		},

		"script_actions": SchemaHDInsightsScriptActions(),
	}

	if definition.CanSpecifyInstanceCount {
		countValidation := validation.IntAtLeast(definition.MinInstanceCount)
		if definition.MaxInstanceCount != nil {
			countValidation = validation.IntBetween(definition.MinInstanceCount, *definition.MaxInstanceCount)
		}

		result["target_instance_count"] = &pluginsdk.Schema{
			Type:         pluginsdk.TypeInt,
			Required:     true,
			ValidateFunc: countValidation,
		}

		if definition.CanAutoScaleByCapacity || definition.CanAutoScaleOnSchedule {
			autoScales := map[string]*pluginsdk.Schema{}

			if definition.CanAutoScaleByCapacity {
				autoScales["capacity"] = &pluginsdk.Schema{
					Type:     pluginsdk.TypeList,
					Optional: true,
					MaxItems: 1,
					ConflictsWith: []string{
						fmt.Sprintf("%s.0.autoscale.0.recurrence", schemaLocation),
					},
					Elem: &pluginsdk.Resource{
						Schema: map[string]*pluginsdk.Schema{
							"min_instance_count": {
								Type:         pluginsdk.TypeInt,
								Required:     true,
								ValidateFunc: countValidation,
							},
							"max_instance_count": {
								Type:         pluginsdk.TypeInt,
								Required:     true,
								ValidateFunc: countValidation,
							},
						},
					},
				}
				if definition.CanAutoScaleOnSchedule {
					autoScales["capacity"].ConflictsWith = []string{
						fmt.Sprintf("%s.0.autoscale.0.recurrence", schemaLocation),
					}
				}
			}
			// managing `azurerm_hdinsight_interactive_query_cluster` autoscaling through `capacity` doesn't work so we'll deprecate this portion of the schema for 4.0
			if definition.CanAutoScaleByCapacityDeprecated4PointOh {
				autoScales["capacity"] = &pluginsdk.Schema{
					Type:     pluginsdk.TypeList,
					Optional: true,
					MaxItems: 1,
					ConflictsWith: []string{
						fmt.Sprintf("%s.0.autoscale.0.recurrence", schemaLocation),
					},
					Deprecated: "HDInsight interactive query clusters can no longer be configured through `autoscale.0.capacity`. Use `autoscale.0.recurrence` instead.",
					Elem: &pluginsdk.Resource{
						Schema: map[string]*pluginsdk.Schema{
							"min_instance_count": {
								Type:         pluginsdk.TypeInt,
								Required:     true,
								ValidateFunc: countValidation,
							},
							"max_instance_count": {
								Type:         pluginsdk.TypeInt,
								Required:     true,
								ValidateFunc: countValidation,
							},
						},
					},
				}
				if definition.CanAutoScaleOnSchedule {
					autoScales["capacity"].ConflictsWith = []string{
						fmt.Sprintf("%s.0.autoscale.0.recurrence", schemaLocation),
					}
				}
			}
			if definition.CanAutoScaleOnSchedule {
				autoScales["recurrence"] = &pluginsdk.Schema{
					Type:     pluginsdk.TypeList,
					Optional: true,
					MaxItems: 1,
					Elem: &pluginsdk.Resource{
						Schema: map[string]*pluginsdk.Schema{
							"timezone": {
								Type:     pluginsdk.TypeString,
								Required: true,
							},
							"schedule": {
								Type:     pluginsdk.TypeList,
								Required: true,
								MinItems: 1,
								Elem: &pluginsdk.Resource{
									Schema: map[string]*pluginsdk.Schema{
										"time": {
											Type:     pluginsdk.TypeString,
											Required: true,
											ValidateFunc: validation.StringMatch(
												regexp.MustCompile("^([01][0-9]|[2][0-3]):([03][0])$"), // time must be on the hour or half past
												"Time of day must match the format HH:mm where HH is 00-23 and mm is 00 or 30",
											),
										},
										"days": {
											Type:     pluginsdk.TypeList,
											Required: true,
											Elem: &pluginsdk.Schema{
												Type: pluginsdk.TypeString,
												ValidateFunc: validation.StringInSlice([]string{
													string(clusters.DaysOfWeekMonday),
													string(clusters.DaysOfWeekTuesday),
													string(clusters.DaysOfWeekWednesday),
													string(clusters.DaysOfWeekThursday),
													string(clusters.DaysOfWeekFriday),
													string(clusters.DaysOfWeekSaturday),
													string(clusters.DaysOfWeekSunday),
												}, false),
											},
										},

										"target_instance_count": {
											Type:         pluginsdk.TypeInt,
											Required:     true,
											ValidateFunc: countValidation,
										},
									},
								},
							},
						},
					},
				}
				if definition.CanAutoScaleByCapacity {
					autoScales["recurrence"].ConflictsWith = []string{
						fmt.Sprintf("%s.0.autoscale.0.capacity", schemaLocation),
					}
				}
			}

			result["autoscale"] = &pluginsdk.Schema{
				Type:     pluginsdk.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: autoScales,
				},
			}
		}
	}

	if definition.CanSpecifyDisks {
		result["number_of_disks_per_node"] = &pluginsdk.Schema{
			Type:         pluginsdk.TypeInt,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.IntBetween(1, *definition.MaxNumberOfDisksPerNode),
		}
	}

	s := &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		MaxItems: 1,
		Required: required,
		Optional: !required,
		Elem: &pluginsdk.Resource{
			Schema: result,
		},
	}

	return s
}

func ExpandHDInsightNodeDefinition(name string, input []interface{}, definition HDInsightNodeDefinition) (*clusters.Role, error) {
	v := input[0].(map[string]interface{})
	vmSize := v["vm_size"].(string)
	username := v["username"].(string)
	password := v["password"].(string)
	virtualNetworkId := v["virtual_network_id"].(string)
	subnetId := v["subnet_id"].(string)
	scriptActions := v["script_actions"].([]interface{})

	role := clusters.Role{
		Name: utils.String(name),
		HardwareProfile: &clusters.HardwareProfile{
			VMSize: utils.String(vmSize),
		},
		OsProfile: &clusters.OsProfile{
			LinuxOperatingSystemProfile: &clusters.LinuxOperatingSystemProfile{},
		},
		ScriptActions: ExpandHDInsightsRolesScriptActions(scriptActions),
	}

	if name != "kafkamanagementnode" {
		role.OsProfile.LinuxOperatingSystemProfile.Username = utils.String(username)
	} else {
		// kafkamanagementnode generates a username and discards the value sent, however, the API has `Username` marked
		// as required non-empty, so we'll send a dummy one avoiding the Portal's default value, which is reserved/invalid.
		role.OsProfile.LinuxOperatingSystemProfile.Username = utils.String("sshadmin")
	}

	virtualNetworkSpecified := virtualNetworkId != ""
	subnetSpecified := subnetId != ""
	if virtualNetworkSpecified && subnetSpecified {
		role.VirtualNetworkProfile = &clusters.VirtualNetworkProfile{
			Id:     utils.String(virtualNetworkId),
			Subnet: utils.String(subnetId),
		}
	} else if (virtualNetworkSpecified && !subnetSpecified) || (subnetSpecified && !virtualNetworkSpecified) {
		return nil, fmt.Errorf("`virtual_network_id` and `subnet_id` must both either be set or empty")
	}

	if password != "" {
		role.OsProfile.LinuxOperatingSystemProfile.Password = utils.String(password)
	} else {
		sshKeysRaw := v["ssh_keys"].(*pluginsdk.Set).List()
		sshKeys := make([]clusters.SshPublicKey, 0)
		for _, v := range sshKeysRaw {
			sshKeys = append(sshKeys, clusters.SshPublicKey{
				CertificateData: utils.String(v.(string)),
			})
		}

		if len(sshKeys) == 0 {
			return nil, fmt.Errorf("either a `password` or `ssh_key` must be specified")
		}

		role.OsProfile.LinuxOperatingSystemProfile.SshProfile = &clusters.SshProfile{
			PublicKeys: &sshKeys,
		}
	}

	if definition.CanSpecifyInstanceCount {
		targetInstanceCount := v["target_instance_count"].(int)
		role.TargetInstanceCount = pointer.To(int64(targetInstanceCount))

		if definition.CanAutoScaleByCapacity || definition.CanAutoScaleOnSchedule {
			autoscaleRaw := v["autoscale"].([]interface{})
			autoscale := ExpandHDInsightNodeAutoScaleDefinition(autoscaleRaw)
			if autoscale != nil {
				role.Autoscale = autoscale
			}
		}
	} else {
		role.MinInstanceCount = definition.FixedMinInstanceCount
		role.TargetInstanceCount = definition.FixedTargetInstanceCount
	}

	if definition.CanSpecifyDisks {
		numberOfDisksPerNode := v["number_of_disks_per_node"].(int)
		if numberOfDisksPerNode > 0 {
			role.DataDisksGroups = &[]clusters.DataDisksGroups{
				{
					DisksPerNode: pointer.To(int64(numberOfDisksPerNode)),
				},
			}
		}
	}

	return &role, nil
}

func ExpandHDInsightNodeAutoScaleDefinition(input []interface{}) *clusters.Autoscale {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	vs := input[0].(map[string]interface{})

	if vs["capacity"] != nil {
		capacityRaw := vs["capacity"].([]interface{})

		capacity := ExpandHDInsightAutoscaleCapacityDefinition(capacityRaw)
		if capacity != nil {
			return &clusters.Autoscale{
				Capacity: capacity,
			}
		}
	}

	if vs["recurrence"] != nil {
		recurrenceRaw := vs["recurrence"].([]interface{})
		recurrence := ExpandHDInsightAutoscaleRecurrenceDefinition(recurrenceRaw)
		if recurrence != nil {
			return &clusters.Autoscale{
				Recurrence: recurrence,
			}
		}
	}

	return nil
}

func ExpandHDInsightAutoscaleCapacityDefinition(input []interface{}) *clusters.AutoscaleCapacity {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	vs := input[0].(map[string]interface{})

	return &clusters.AutoscaleCapacity{
		MinInstanceCount: pointer.To(int64(vs["min_instance_count"].(int))),
		MaxInstanceCount: pointer.To(int64(vs["max_instance_count"].(int))),
	}
}

func ExpandHDInsightAutoscaleRecurrenceDefinition(input []interface{}) *clusters.AutoscaleRecurrence {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	vs := input[0].(map[string]interface{})

	schedules := make([]clusters.AutoscaleSchedule, 0)

	for _, v := range vs["schedule"].([]interface{}) {
		val := v.(map[string]interface{})

		weekDays := val["days"].([]interface{})
		expandedWeekDays := make([]clusters.DaysOfWeek, len(weekDays))
		for i := range weekDays {
			expandedWeekDays[i] = clusters.DaysOfWeek(weekDays[i].(string))
		}

		schedules = append(schedules, clusters.AutoscaleSchedule{
			Days: &expandedWeekDays,
			TimeAndCapacity: &clusters.AutoscaleTimeAndCapacity{
				Time: utils.String(val["time"].(string)),
				// SDK supports min and max, but server side always overrides max to be equal to min
				MinInstanceCount: pointer.To(int64(val["target_instance_count"].(int))),
				MaxInstanceCount: pointer.To(int64(val["target_instance_count"].(int))),
			},
		})
	}

	result := &clusters.AutoscaleRecurrence{
		TimeZone: utils.String(vs["timezone"].(string)),
		Schedule: &schedules,
	}

	return result
}

func ExpandHDInsightSecurityProfile(input []interface{}) *clusters.SecurityProfile {
	if len(input) == 0 {
		return nil
	}

	v := input[0].(map[string]interface{})

	result := clusters.SecurityProfile{
		DirectoryType:      pointer.To(clusters.DirectoryTypeActiveDirectory),
		Domain:             utils.String(v["domain_name"].(string)),
		LdapsURLs:          utils.ExpandStringSlice(v["ldaps_urls"].(*pluginsdk.Set).List()),
		DomainUsername:     utils.String(v["domain_username"].(string)),
		DomainUserPassword: utils.String(v["domain_user_password"].(string)),
		AaddsResourceId:    utils.String(v["aadds_resource_id"].(string)),
		MsiResourceId:      utils.String(v["msi_resource_id"].(string)),
	}

	if clusterUsersGroupDNS := v["cluster_users_group_dns"].(*pluginsdk.Set).List(); len(clusterUsersGroupDNS) != 0 {
		result.ClusterUsersGroupDNs = utils.ExpandStringSlice(clusterUsersGroupDNS)
	}

	return &result
}

func FlattenHDInsightNodeDefinition(input *clusters.Role, existing []interface{}, definition HDInsightNodeDefinition) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	output := map[string]interface{}{
		"vm_size":            "",
		"username":           "",
		"password":           "",
		"ssh_keys":           pluginsdk.NewSet(pluginsdk.HashString, []interface{}{}),
		"subnet_id":          "",
		"virtual_network_id": "",
		"script_actions":     make([]interface{}, 0),
	}

	if profile := input.OsProfile; profile != nil {
		if osProfile := profile.LinuxOperatingSystemProfile; osProfile != nil {
			if username := osProfile.Username; username != nil {
				output["username"] = *username
			}
		}
	}

	// neither Password / SSH Keys are returned from the API, so we need to look them up to not force a diff
	if len(existing) > 0 {
		existingV := existing[0].(map[string]interface{})
		output["password"] = existingV["password"].(string)

		sshKeys := existingV["ssh_keys"].(*pluginsdk.Set).List()
		output["ssh_keys"] = pluginsdk.NewSet(pluginsdk.HashString, sshKeys)

		// whilst the VMSize can be returned from `input.HardwareProfile.VMSize` - it can be malformed
		// for example, `small`, `medium`, `large` and `extralarge` can be returned inside of actual VM Size
		// after extensive experimentation it appears multiple instance sizes fit `extralarge`, as such
		// unfortunately we can't transform these; since it can't be changed
		// we should be "safe" to try and pull it from the state instead, but clearly this isn't ideal
		vmSize := existingV["vm_size"].(string)
		output["vm_size"] = vmSize

		scriptActions := existingV["script_actions"].([]interface{})
		output["script_actions"] = scriptActions
	}

	if profile := input.VirtualNetworkProfile; profile != nil {
		if profile.Id != nil {
			output["virtual_network_id"] = *profile.Id
		}
		if profile.Subnet != nil {
			output["subnet_id"] = *profile.Subnet
		}
	}

	if definition.CanSpecifyInstanceCount {
		output["target_instance_count"] = 0

		if input.TargetInstanceCount != nil {
			output["target_instance_count"] = int(*input.TargetInstanceCount)
		}

		if definition.CanAutoScaleByCapacity || definition.CanAutoScaleOnSchedule {
			autoscale := FlattenHDInsightNodeAutoscaleDefinition(input.Autoscale)
			if autoscale != nil {
				output["autoscale"] = autoscale
			}
		}
	}

	if definition.CanSpecifyDisks {
		output["number_of_disks_per_node"] = 0
		if input.DataDisksGroups != nil && len(*input.DataDisksGroups) > 0 {
			group := (*input.DataDisksGroups)[0]
			if group.DisksPerNode != nil {
				output["number_of_disks_per_node"] = int(*group.DisksPerNode)
			}
		}
	}

	return []interface{}{output}
}

func FindHDInsightRole(input *[]clusters.Role, name string) *clusters.Role {
	if input == nil {
		return nil
	}

	for _, v := range *input {
		if v.Name == nil {
			continue
		}

		actualName := *v.Name
		if strings.EqualFold(name, actualName) {
			return &v
		}
	}

	return nil
}

func findHDInsightConnectivityEndpoint(name string, input *[]clusters.ConnectivityEndpoint) string {
	if input == nil {
		return ""
	}

	for _, v := range *input {
		if v.Name == nil || v.Location == nil {
			continue
		}

		if strings.EqualFold(*v.Name, name) {
			return *v.Location
		}
	}

	return ""
}

func FlattenHDInsightNodeAutoscaleDefinition(input *clusters.Autoscale) []interface{} {
	if input == nil {
		return nil
	}

	result := map[string]interface{}{}

	if input.Capacity != nil {
		result["capacity"] = FlattenHDInsightAutoscaleCapacityDefinition(input.Capacity)
	}

	if input.Recurrence != nil {
		result["recurrence"] = FlattenHDInsightAutoscaleRecurrenceDefinition(input.Recurrence)
	}

	if len(result) > 0 {
		return []interface{}{result}
	}
	return nil
}

func FlattenHDInsightAutoscaleCapacityDefinition(input *clusters.AutoscaleCapacity) []interface{} {
	return []interface{}{
		map[string]interface{}{
			"min_instance_count": input.MinInstanceCount,
			"max_instance_count": input.MaxInstanceCount,
		},
	}
}

func FlattenHDInsightAutoscaleRecurrenceDefinition(input *clusters.AutoscaleRecurrence) []interface{} {
	if input.Schedule == nil {
		return []interface{}{}
	}

	schedules := make([]interface{}, 0)

	for _, schedule := range *input.Schedule {
		days := make([]clusters.DaysOfWeek, 0)
		if schedule.Days != nil {
			days = *schedule.Days
		}

		targetInstanceCount := 0
		time := ""
		if schedule.TimeAndCapacity != nil {
			if schedule.TimeAndCapacity.MinInstanceCount != nil {
				// note: min / max are the same
				targetInstanceCount = int(*schedule.TimeAndCapacity.MinInstanceCount)
			}
			if *schedule.TimeAndCapacity.Time != "" {
				time = *schedule.TimeAndCapacity.Time
			}
		}
		schedules = append(schedules, map[string]interface{}{
			"days":                  days,
			"target_instance_count": targetInstanceCount,
			"time":                  time,
		})
	}

	return []interface{}{
		map[string]interface{}{
			"timezone": input.TimeZone,
			"schedule": &schedules,
		},
	}
}

func flattenHDInsightSecurityProfile(input *clusters.SecurityProfile, d *pluginsdk.ResourceData) []interface{} {
	if input == nil {
		return make([]interface{}, 0)
	}

	var aaddsResourceId string
	if input.AaddsResourceId != nil {
		aaddsResourceId = *input.AaddsResourceId
	}

	var domain string
	if input.Domain != nil {
		domain = *input.Domain
	}

	var domainUsername string
	if input.DomainUsername != nil {
		domainUsername = *input.DomainUsername
	}

	var msiResourceId string
	if input.MsiResourceId != nil {
		msiResourceId = *input.MsiResourceId
	}

	return []interface{}{
		map[string]interface{}{
			"aadds_resource_id":       aaddsResourceId,
			"cluster_users_group_dns": utils.FlattenStringSlice(input.ClusterUsersGroupDNs),
			"domain_name":             domain,
			"domain_username":         domainUsername,
			"domain_user_password":    d.Get("security_profile.0.domain_user_password"),
			"ldaps_urls":              utils.FlattenStringSlice(input.LdapsURLs),
			"msi_resource_id":         msiResourceId,
		},
	}
}
