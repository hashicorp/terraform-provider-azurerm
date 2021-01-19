package hdinsight

import (
	"fmt"
	"strings"

	"github.com/Azure/azure-sdk-for-go/services/hdinsight/mgmt/2018-06-01/hdinsight"
	"github.com/hashicorp/go-getter/helper/url"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/suppress"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/location"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func SchemaHDInsightName() *schema.Schema {
	return &schema.Schema{
		Type:         schema.TypeString,
		Required:     true,
		ForceNew:     true,
		ValidateFunc: validate.HDInsightName,
	}
}

func SchemaHDInsightDataSourceName() *schema.Schema {
	return &schema.Schema{
		Type:         schema.TypeString,
		Required:     true,
		ValidateFunc: validate.HDInsightName,
	}
}

func SchemaHDInsightTier() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeString,
		Required: true,
		ForceNew: true,
		ValidateFunc: validation.StringInSlice([]string{
			string(hdinsight.Standard),
			string(hdinsight.Premium),
		}, true),
		// TODO: file a bug about this
		DiffSuppressFunc: location.DiffSuppressFunc,
	}
}

func SchemaHDInsightTls() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeString,
		Optional: true,
		ForceNew: true,
		ValidateFunc: validation.StringInSlice([]string{
			"1.0",
			"1.1",
			"1.2",
		}, false),
	}
}

func SchemaHDInsightClusterVersion() *schema.Schema {
	return &schema.Schema{
		Type:             schema.TypeString,
		Required:         true,
		ForceNew:         true,
		ValidateFunc:     validate.HDInsightClusterVersion,
		DiffSuppressFunc: hdinsightClusterVersionDiffSuppressFunc,
	}
}

func hdinsightClusterVersionDiffSuppressFunc(_, old, new string, _ *schema.ResourceData) bool {
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

func SchemaHDInsightsGateway() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Required: true,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				// TODO 3.0: remove this attribute
				"enabled": {
					Type:       schema.TypeBool,
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
					Type:     schema.TypeString,
					Required: true,
					ForceNew: true,
				},
				"password": {
					Type:      schema.TypeString,
					Required:  true,
					Sensitive: true,
					// Azure returns the key as *****. We'll suppress that here.
					DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
						return (new == d.Get(k).(string)) && (old == "*****")
					},
				},
			},
		},
	}
}

func SchemaHDInsightsExternalMetastore() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"server": {
					Type:     schema.TypeString,
					Required: true,
					ForceNew: true,
				},
				"database_name": {
					Type:     schema.TypeString,
					Required: true,
					ForceNew: true,
				},
				"username": {
					Type:     schema.TypeString,
					Required: true,
					ForceNew: true,
				},
				"password": {
					Type:      schema.TypeString,
					Required:  true,
					ForceNew:  true,
					Sensitive: true,
					// Azure returns the key as *****. We'll suppress that here.
					DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
						return (new == d.Get(k).(string)) && (old == "*****")
					},
				},
			},
		},
	}
}

func SchemaHDInsightsExternalMetastores() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"hive": SchemaHDInsightsExternalMetastore(),

				"oozie": SchemaHDInsightsExternalMetastore(),

				"ambari": SchemaHDInsightsExternalMetastore(),
			},
		},
	}
}

func SchemaHDInsightsMonitor() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"log_analytics_workspace_id": {
					Type:         schema.TypeString,
					Required:     true,
					ValidateFunc: validation.IsUUID,
				},
				"primary_key": {
					Type:         schema.TypeString,
					Required:     true,
					Sensitive:    true,
					ValidateFunc: validation.StringIsNotEmpty,
					// Azure doesn't return the key
					DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
						return (new == d.Get(k).(string)) && (old == "*****")
					},
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

func ExpandHDInsightsMonitor(input []interface{}) hdinsight.ClusterMonitoringRequest {
	vs := input[0].(map[string]interface{})

	workspace := vs["log_analytics_workspace_id"].(string)
	key := vs["primary_key"].(string)

	return hdinsight.ClusterMonitoringRequest{
		WorkspaceID: utils.String(workspace),
		PrimaryKey:  utils.String(key),
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

func SchemaHDInsightsStorageAccounts() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"storage_account_key": {
					Type:         schema.TypeString,
					Required:     true,
					ForceNew:     true,
					Sensitive:    true,
					ValidateFunc: validation.StringIsNotEmpty,
				},
				"storage_container_id": {
					Type:         schema.TypeString,
					Required:     true,
					ForceNew:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},
				"is_default": {
					Type:     schema.TypeBool,
					Required: true,
					ForceNew: true,
				},
			},
		},
	}
}

func SchemaHDInsightsGen2StorageAccounts() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		// HDInsight doesn't seem to allow adding more than one gen2 cluster right now.
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"storage_resource_id": {
					Type:         schema.TypeString,
					Required:     true,
					ForceNew:     true,
					ValidateFunc: azure.ValidateResourceID,
				},
				"filesystem_id": {
					Type:         schema.TypeString,
					Required:     true,
					ForceNew:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},
				"managed_identity_resource_id": {
					Type:         schema.TypeString,
					Required:     true,
					ForceNew:     true,
					ValidateFunc: azure.ValidateResourceID,
				},
				"is_default": {
					Type:     schema.TypeBool,
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
				Type:                   hdinsight.UserAssigned,
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
}

func ValidateSchemaHDInsightNodeDefinitionVMSize() schema.SchemaValidateFunc {
	return validation.StringInSlice([]string{
		// short of deploying every VM Sku for every node type for every HDInsight Cluster
		// this is the list I've (@tombuildsstuff) found for valid SKU's from an endpoint in the Portal
		// using another SKU causes a bad request from the API - as such this is a best effort UX
		"ExtraSmall",
		"Small",
		"Medium",
		"Large",
		"ExtraLarge",
		"A5",
		"A6",
		"A7",
		"A8",
		"A9",
		"A10",
		"A11",
		"Standard_A1_V2",
		"Standard_A2_V2",
		"Standard_A2m_V2",
		"Standard_A3",
		"Standard_A4_V2",
		"Standard_A4m_V2",
		"Standard_A8_V2",
		"Standard_A8m_V2",
		"Standard_D1",
		"Standard_D2",
		"Standard_D3",
		"Standard_D4",
		"Standard_D11",
		"Standard_D12",
		"Standard_D13",
		"Standard_D14",
		"Standard_D1_V2",
		"Standard_D2_V2",
		"Standard_D3_V2",
		"Standard_D4_V2",
		"Standard_D5_V2",
		"Standard_D11_V2",
		"Standard_D12_V2",
		"Standard_D13_V2",
		"Standard_D14_V2",
		"Standard_DS1_V2",
		"Standard_DS2_V2",
		"Standard_DS3_V2",
		"Standard_DS4_V2",
		"Standard_DS5_V2",
		"Standard_DS11_V2",
		"Standard_DS12_V2",
		"Standard_DS13_V2",
		"Standard_DS14_V2",
		"Standard_D4a_V4",
		"Standard_E2_V3",
		"Standard_E4_V3",
		"Standard_E8_V3",
		"Standard_E16_V3",
		"Standard_E20_V3",
		"Standard_E32_V3",
		"Standard_E64_V3",
		"Standard_E64i_V3",
		"Standard_E2s_V3",
		"Standard_E4s_V3",
		"Standard_E8s_V3",
		"Standard_E16s_V3",
		"Standard_E20s_V3",
		"Standard_E32s_V3",
		"Standard_E64s_V3",
		"Standard_E64is_V3",
		"Standard_G1",
		"Standard_G2",
		"Standard_G3",
		"Standard_G4",
		"Standard_G5",
		"Standard_F2s_V2",
		"Standard_F4s_V2",
		"Standard_F8s_V2",
		"Standard_F16s_V2",
		"Standard_F32s_V2",
		"Standard_F64s_V2",
		"Standard_F72s_V2",
		"Standard_GS1",
		"Standard_GS2",
		"Standard_GS3",
		"Standard_GS4",
		"Standard_GS5",
		"Standard_NC24",
	}, true)
}

func SchemaHDInsightNodeDefinition(schemaLocation string, definition HDInsightNodeDefinition, required bool) *schema.Schema {
	result := map[string]*schema.Schema{
		"vm_size": {
			Type:             schema.TypeString,
			Required:         true,
			ForceNew:         true,
			DiffSuppressFunc: suppress.CaseDifference,
			ValidateFunc:     ValidateSchemaHDInsightNodeDefinitionVMSize(),
		},
		"username": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},
		"password": {
			Type:      schema.TypeString,
			Optional:  true,
			ForceNew:  true,
			Sensitive: true,
		},
		"ssh_keys": {
			Type:     schema.TypeSet,
			Optional: true,
			ForceNew: true,
			Elem: &schema.Schema{
				Type:         schema.TypeString,
				ValidateFunc: validation.StringIsNotEmpty,
			},
			Set: schema.HashString,
			ConflictsWith: []string{
				fmt.Sprintf("%s.0.password", schemaLocation),
			},
		},

		"subnet_id": {
			Type:         schema.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: azure.ValidateResourceIDOrEmpty,
		},

		"virtual_network_id": {
			Type:         schema.TypeString,
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
		result["min_instance_count"] = &schema.Schema{
			Type:         schema.TypeInt,
			Optional:     true,
			ForceNew:     true,
			Computed:     true,
			Deprecated:   "this has been deprecated from the API and will be removed in version 3.0 of the provider",
			ValidateFunc: countValidation,
		}
		result["target_instance_count"] = &schema.Schema{
			Type:         schema.TypeInt,
			Required:     true,
			ValidateFunc: countValidation,
		}
	}

	if definition.CanSpecifyDisks {
		result["number_of_disks_per_node"] = &schema.Schema{
			Type:         schema.TypeInt,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.IntBetween(1, *definition.MaxNumberOfDisksPerNode),
		}
	}

	s := &schema.Schema{
		Type:     schema.TypeList,
		MaxItems: 1,
		Required: required,
		Optional: !required,
		Elem: &schema.Resource{
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
		sshKeysRaw := v["ssh_keys"].(*schema.Set).List()
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

func FlattenHDInsightNodeDefinition(input *hdinsight.Role, existing []interface{}, definition HDInsightNodeDefinition) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	output := map[string]interface{}{
		"vm_size":            "",
		"username":           "",
		"password":           "",
		"ssh_keys":           schema.NewSet(schema.HashString, []interface{}{}),
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

		sshKeys := existingV["ssh_keys"].(*schema.Set).List()
		output["ssh_keys"] = schema.NewSet(schema.HashString, sshKeys)

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
