package hdinsight

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/Azure/azure-sdk-for-go/services/hdinsight/mgmt/2018-06-01/hdinsight"
	"github.com/hashicorp/go-getter/helper/url"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/location"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/hdinsight/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/suppress"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func SchemaHDInsightName() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:         pluginsdk.TypeString,
		Required:     true,
		ForceNew:     true,
		ValidateFunc: validate.HDInsightName,
	}
}

func SchemaHDInsightDataSourceName() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:         pluginsdk.TypeString,
		Required:     true,
		ValidateFunc: validate.HDInsightName,
	}
}

func SchemaHDInsightTier() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeString,
		Required: true,
		ForceNew: true,
		ValidateFunc: validation.StringInSlice([]string{
			string(hdinsight.TierStandard),
			string(hdinsight.TierPremium),
		}, true),
		// TODO: file a bug about this
		DiffSuppressFunc: location.DiffSuppressFunc,
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
				// TODO 3.0: remove this attribute
				"enabled": {
					Type:       pluginsdk.TypeBool,
					Optional:   true,
					Default:    true,
					Deprecated: "HDInsight doesn't support disabling gateway anymore",
					ValidateFunc: func(i interface{}, k string) (warnings []string, errors []error) {
						enabled := i.(bool)

						if !enabled {
							errors = append(errors, fmt.Errorf("Only true is supported, because HDInsight doesn't support disabling gateway anymore. Provided value %t", enabled))
						}
						return warnings, errors
					},
				},
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
					Default:  string(hdinsight.ResourceProviderConnectionInbound),
					ValidateFunc: validation.StringInSlice([]string{
						string(hdinsight.ResourceProviderConnectionInbound),
						string(hdinsight.ResourceProviderConnectionOutbound),
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

func ExpandHDInsightsMonitor(input []interface{}) hdinsight.ClusterMonitoringRequest {
	vs := input[0].(map[string]interface{})

	workspace := vs["log_analytics_workspace_id"].(string)
	key := vs["primary_key"].(string)

	return hdinsight.ClusterMonitoringRequest{
		WorkspaceID: utils.String(workspace),
		PrimaryKey:  utils.String(key),
	}
}

func ExpandHDInsightsNetwork(input []interface{}) *hdinsight.NetworkProperties {
	if len(input) == 0 {
		return nil
	}

	vs := input[0].(map[string]interface{})

	connDir := hdinsight.ResourceProviderConnectionOutbound
	if v, exists := vs["connection_direction"]; exists && v != string(hdinsight.ResourceProviderConnectionOutbound) {
		connDir = hdinsight.ResourceProviderConnectionInbound
	}

	privateLink := hdinsight.PrivateLinkDisabled
	if v, exists := vs["private_link_enabled"]; exists && v != false {
		privateLink = hdinsight.PrivateLinkEnabled
	}

	return &hdinsight.NetworkProperties{
		ResourceProviderConnection: connDir,
		PrivateLink:                privateLink,
	}
}

func FlattenHDInsightsNetwork(input *hdinsight.NetworkProperties) []interface{} {
	if input == nil {
		return nil
	}

	connDir := string(hdinsight.ResourceProviderConnectionOutbound)
	if v := input.ResourceProviderConnection; v != "" {
		connDir = string(v)
	}

	privateLink := false
	if v := input.PrivateLink; v != "" {
		privateLink = (v == hdinsight.PrivateLinkEnabled)
	}

	return []interface{}{
		map[string]interface{}{
			"connection_direction": connDir,
			"private_link_enabled": privateLink,
		},
	}
}

func FlattenHDInsightsConfigurations(input map[string]*string) []interface{} {
	enabled := true

	username := ""
	if v, exists := input["restAuthCredential.username"]; exists && v != nil {
		username = *v
	}

	password := ""
	if v, exists := input["restAuthCredential.password"]; exists && v != nil {
		password = *v
	}

	return []interface{}{
		map[string]interface{}{
			"enabled":  enabled,
			"username": username,
			"password": password,
		},
	}
}

func FlattenHDInsightsHiveMetastore(env map[string]*string, site map[string]*string) []interface{} {
	server := ""
	if v, exists := env["hive_hostname"]; exists && v != nil {
		server = *v
	}

	database := ""
	if v, exists := env["hive_database_name"]; exists && v != nil {
		database = *v
	}

	username := ""
	if v, exists := site["javax.jdo.option.ConnectionUserName"]; exists && v != nil {
		username = *v
	}

	password := ""
	if v, exists := site["javax.jdo.option.ConnectionPassword"]; exists && v != nil {
		password = *v
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

	return nil
}

func FlattenHDInsightsOozieMetastore(env map[string]*string, site map[string]*string) []interface{} {
	server := ""
	if v, exists := env["oozie_hostname"]; exists && v != nil {
		server = *v
	}

	database := ""
	if v, exists := env["oozie_database_name"]; exists && v != nil {
		database = *v
	}

	username := ""
	if v, exists := site["oozie.service.JPAService.jdbc.username"]; exists && v != nil {
		username = *v
	}

	password := ""
	if v, exists := site["oozie.service.JPAService.jdbc.password"]; exists && v != nil {
		password = *v
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

	return nil
}

func FlattenHDInsightsAmbariMetastore(conf map[string]*string) []interface{} {
	server := ""
	if v, exists := conf["database-server"]; exists && v != nil {
		server = *v
	}

	database := ""
	if v, exists := conf["database-name"]; exists && v != nil {
		database = *v
	}

	username := ""
	if v, exists := conf["database-user-name"]; exists && v != nil {
		username = *v
	}

	password := ""
	if v, exists := conf["database-user-password"]; exists && v != nil {
		password = *v
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

	return nil
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
				"storage_resource_id": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ForceNew:     true,
					ValidateFunc: azure.ValidateResourceID,
				},
				"filesystem_id": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ForceNew:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},
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

// ExpandHDInsightsStorageAccounts returns an array of StorageAccount structs, as well as a ClusterIdentity
// populated with any managed identities required for accessing Data Lake Gen2 storage.
func ExpandHDInsightsStorageAccounts(storageAccounts []interface{}, gen2storageAccounts []interface{}) (*[]hdinsight.StorageAccount, *hdinsight.ClusterIdentity, error) {
	results := make([]hdinsight.StorageAccount, 0)

	var clusterIndentity *hdinsight.ClusterIdentity

	for _, vs := range storageAccounts {
		v := vs.(map[string]interface{})

		storageAccountKey := v["storage_account_key"].(string)
		storageContainerID := v["storage_container_id"].(string)
		isDefault := v["is_default"].(bool)

		uri, err := url.Parse(storageContainerID)
		if err != nil {
			return nil, nil, fmt.Errorf("Error parsing %q: %s", storageContainerID, err)
		}

		result := hdinsight.StorageAccount{
			Name:      utils.String(uri.Host),
			Container: utils.String(strings.TrimPrefix(uri.Path, "/")),
			Key:       utils.String(storageAccountKey),
			IsDefault: utils.Bool(isDefault),
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
			return nil, nil, fmt.Errorf("Error parsing %q: %s", fileSystemID, err)
		}

		if clusterIndentity == nil {
			clusterIndentity = &hdinsight.ClusterIdentity{
				Type:                   hdinsight.ResourceIdentityTypeUserAssigned,
				UserAssignedIdentities: make(map[string]*hdinsight.ClusterIdentityUserAssignedIdentitiesValue),
			}
		}

		// ... API doesn't seem to require client_id or principal_id, so pass in an empty ClusterIdentityUserAssignedIdentitiesValue
		clusterIndentity.UserAssignedIdentities[managedIdentityResourceID] = &hdinsight.ClusterIdentityUserAssignedIdentitiesValue{}

		result := hdinsight.StorageAccount{
			Name:          utils.String(uri.Host), // https://storageaccountname.dfs.core.windows.net/filesystemname -> storageaccountname.dfs.core.windows.net
			ResourceID:    utils.String(storageResourceID),
			FileSystem:    utils.String(uri.Path[1:]), // https://storageaccountname.dfs.core.windows.net/filesystemname -> filesystemname
			MsiResourceID: utils.String(managedIdentityResourceID),
			IsDefault:     utils.Bool(isDefault),
		}
		results = append(results, result)
	}

	return &results, clusterIndentity, nil
}

type HDInsightNodeDefinition struct {
	CanSpecifyInstanceCount  bool
	MinInstanceCount         int
	MaxInstanceCount         *int
	CanSpecifyDisks          bool
	MaxNumberOfDisksPerNode  *int
	FixedMinInstanceCount    *int32
	FixedTargetInstanceCount *int32
	CanAutoScaleByCapacity   bool
	CanAutoScaleOnSchedule   bool
}

func SchemaHDInsightNodeDefinition(schemaLocation string, definition HDInsightNodeDefinition, required bool) *pluginsdk.Schema {
	result := map[string]*pluginsdk.Schema{
		"vm_size": {
			Type:             pluginsdk.TypeString,
			Required:         true,
			ForceNew:         true,
			DiffSuppressFunc: suppress.CaseDifference,
			ValidateFunc:     validate.NodeDefinitionVMSize(),
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
			ValidateFunc: azure.ValidateResourceIDOrEmpty,
		},

		"virtual_network_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: azure.ValidateResourceIDOrEmpty,
		},
	}

	if definition.CanSpecifyInstanceCount {
		countValidation := validation.IntAtLeast(definition.MinInstanceCount)
		if definition.MaxInstanceCount != nil {
			countValidation = validation.IntBetween(definition.MinInstanceCount, *definition.MaxInstanceCount)
		}

		// TODO 3.0: remove this property
		result["min_instance_count"] = &pluginsdk.Schema{
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			ForceNew:     true,
			Computed:     true,
			Deprecated:   "this has been deprecated from the API and will be removed in version 3.0 of the provider",
			ValidateFunc: countValidation,
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
													string(hdinsight.DaysOfWeekMonday),
													string(hdinsight.DaysOfWeekTuesday),
													string(hdinsight.DaysOfWeekWednesday),
													string(hdinsight.DaysOfWeekThursday),
													string(hdinsight.DaysOfWeekFriday),
													string(hdinsight.DaysOfWeekSaturday),
													string(hdinsight.DaysOfWeekSunday),
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

func ExpandHDInsightNodeDefinition(name string, input []interface{}, definition HDInsightNodeDefinition) (*hdinsight.Role, error) {
	v := input[0].(map[string]interface{})
	vmSize := v["vm_size"].(string)
	username := v["username"].(string)
	password := v["password"].(string)
	virtualNetworkId := v["virtual_network_id"].(string)
	subnetId := v["subnet_id"].(string)

	role := hdinsight.Role{
		Name: utils.String(name),
		HardwareProfile: &hdinsight.HardwareProfile{
			VMSize: utils.String(vmSize),
		},
		OsProfile: &hdinsight.OsProfile{
			LinuxOperatingSystemProfile: &hdinsight.LinuxOperatingSystemProfile{
				Username: utils.String(username),
			},
		},
	}

	virtualNetworkSpecified := virtualNetworkId != ""
	subnetSpecified := subnetId != ""
	if virtualNetworkSpecified && subnetSpecified {
		role.VirtualNetworkProfile = &hdinsight.VirtualNetworkProfile{
			ID:     utils.String(virtualNetworkId),
			Subnet: utils.String(subnetId),
		}
	} else if (virtualNetworkSpecified && !subnetSpecified) || (subnetSpecified && !virtualNetworkSpecified) {
		return nil, fmt.Errorf("`virtual_network_id` and `subnet_id` must both either be set or empty!")
	}

	if password != "" {
		role.OsProfile.LinuxOperatingSystemProfile.Password = utils.String(password)
	} else {
		sshKeysRaw := v["ssh_keys"].(*pluginsdk.Set).List()
		sshKeys := make([]hdinsight.SSHPublicKey, 0)
		for _, v := range sshKeysRaw {
			sshKeys = append(sshKeys, hdinsight.SSHPublicKey{
				CertificateData: utils.String(v.(string)),
			})
		}

		if len(sshKeys) == 0 {
			return nil, fmt.Errorf("Either a `password` or `ssh_key` must be specified!")
		}

		role.OsProfile.LinuxOperatingSystemProfile.SSHProfile = &hdinsight.SSHProfile{
			PublicKeys: &sshKeys,
		}
	}

	if definition.CanSpecifyInstanceCount {
		minInstanceCount := v["min_instance_count"].(int)
		if minInstanceCount > 0 {
			role.MinInstanceCount = utils.Int32(int32(minInstanceCount))
		}

		targetInstanceCount := v["target_instance_count"].(int)
		role.TargetInstanceCount = utils.Int32(int32(targetInstanceCount))

		if definition.CanAutoScaleByCapacity || definition.CanAutoScaleOnSchedule {
			autoscaleRaw := v["autoscale"].([]interface{})
			autoscale := ExpandHDInsightNodeAutoScaleDefinition(autoscaleRaw)
			if autoscale != nil {
				role.AutoscaleConfiguration = autoscale
			}
		}
	} else {
		role.MinInstanceCount = definition.FixedMinInstanceCount
		role.TargetInstanceCount = definition.FixedTargetInstanceCount
	}

	if definition.CanSpecifyDisks {
		numberOfDisksPerNode := v["number_of_disks_per_node"].(int)
		if numberOfDisksPerNode > 0 {
			role.DataDisksGroups = &[]hdinsight.DataDisksGroups{
				{
					DisksPerNode: utils.Int32(int32(numberOfDisksPerNode)),
				},
			}
		}
	}

	return &role, nil
}

func ExpandHDInsightNodeAutoScaleDefinition(input []interface{}) *hdinsight.Autoscale {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	vs := input[0].(map[string]interface{})

	if vs["capacity"] != nil {
		capacityRaw := vs["capacity"].([]interface{})

		capacity := ExpandHDInsightAutoscaleCapacityDefinition(capacityRaw)
		if capacity != nil {
			return &hdinsight.Autoscale{
				Capacity: capacity,
			}
		}
	}

	if vs["recurrence"] != nil {
		recurrenceRaw := vs["recurrence"].([]interface{})
		recurrence := ExpandHDInsightAutoscaleRecurrenceDefinition(recurrenceRaw)
		if recurrence != nil {
			return &hdinsight.Autoscale{
				Recurrence: recurrence,
			}
		}
	}

	return nil
}

func ExpandHDInsightAutoscaleCapacityDefinition(input []interface{}) *hdinsight.AutoscaleCapacity {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	vs := input[0].(map[string]interface{})

	return &hdinsight.AutoscaleCapacity{
		MinInstanceCount: utils.Int32(int32(vs["min_instance_count"].(int))),
		MaxInstanceCount: utils.Int32(int32(vs["max_instance_count"].(int))),
	}
}

func ExpandHDInsightAutoscaleRecurrenceDefinition(input []interface{}) *hdinsight.AutoscaleRecurrence {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	vs := input[0].(map[string]interface{})

	schedules := make([]hdinsight.AutoscaleSchedule, 0)

	for _, v := range vs["schedule"].([]interface{}) {
		val := v.(map[string]interface{})

		weekDays := val["days"].([]interface{})
		expandedWeekDays := make([]hdinsight.DaysOfWeek, len(weekDays))
		for i := range weekDays {
			expandedWeekDays[i] = hdinsight.DaysOfWeek(weekDays[i].(string))
		}

		schedules = append(schedules, hdinsight.AutoscaleSchedule{
			Days: &expandedWeekDays,
			TimeAndCapacity: &hdinsight.AutoscaleTimeAndCapacity{
				Time: utils.String(val["time"].(string)),
				// SDK supports min and max, but server side always overrides max to be equal to min
				MinInstanceCount: utils.Int32(int32(val["target_instance_count"].(int))),
				MaxInstanceCount: utils.Int32(int32(val["target_instance_count"].(int))),
			},
		})
	}

	result := &hdinsight.AutoscaleRecurrence{
		TimeZone: utils.String(vs["timezone"].(string)),
		Schedule: &schedules,
	}

	return result
}

func FlattenHDInsightNodeDefinition(input *hdinsight.Role, existing []interface{}, definition HDInsightNodeDefinition) []interface{} {
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
	}

	if profile := input.VirtualNetworkProfile; profile != nil {
		if profile.ID != nil {
			output["virtual_network_id"] = *profile.ID
		}
		if profile.Subnet != nil {
			output["subnet_id"] = *profile.Subnet
		}
	}

	if definition.CanSpecifyInstanceCount {
		output["min_instance_count"] = 0
		output["target_instance_count"] = 0

		if input.MinInstanceCount != nil {
			output["min_instance_count"] = int(*input.MinInstanceCount)
		}

		if input.TargetInstanceCount != nil {
			output["target_instance_count"] = int(*input.TargetInstanceCount)
		}

		if definition.CanAutoScaleByCapacity || definition.CanAutoScaleOnSchedule {
			autoscale := FlattenHDInsightNodeAutoscaleDefinition(input.AutoscaleConfiguration)
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

func FindHDInsightRole(input *[]hdinsight.Role, name string) *hdinsight.Role {
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

func FindHDInsightConnectivityEndpoint(name string, input *[]hdinsight.ConnectivityEndpoint) string {
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

func FlattenHDInsightNodeAutoscaleDefinition(input *hdinsight.Autoscale) []interface{} {
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

func FlattenHDInsightAutoscaleCapacityDefinition(input *hdinsight.AutoscaleCapacity) []interface{} {
	return []interface{}{
		map[string]interface{}{
			"min_instance_count": input.MinInstanceCount,
			"max_instance_count": input.MaxInstanceCount,
		},
	}
}

func FlattenHDInsightAutoscaleRecurrenceDefinition(input *hdinsight.AutoscaleRecurrence) []interface{} {
	if input.Schedule == nil {
		return []interface{}{}
	}

	schedules := make([]interface{}, 0)

	for _, schedule := range *input.Schedule {
		days := make([]hdinsight.DaysOfWeek, 0)
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
