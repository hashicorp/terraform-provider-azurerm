package helper

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/sqlvirtualmachine/2022-02-01/sqlvirtualmachinegroups"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type WsfcDomainProfile struct {
	Fqdn                        string `tfschema:"fqdn"`
	OuPath                      string `tfschema:"ou_path"`
	ClusterBootstrapAccountName string `tfschema:"cluster_bootstrap_account_name"`
	ClusterOperatorAccountName  string `tfschema:"cluster_operator_account_name"`
	SqlServiceAccountName       string `tfschema:"sql_service_account_name"`
	FileShareWitnessPath        string `tfschema:"file_share_witness_path"`
	StorageAccountUrl           string `tfschema:"storage_account_url"`
	StorageAccountPrimaryKey    string `tfschema:"storage_account_primary_key"`
	ClusterSubnetType           string `tfschema:"cluster_subnet_type"`
}

func WsfcDomainProfileSchemaMsSqlVirtualMachineAvailabilityGroup() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Required: true,
		MaxItems: 1,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"fqdn": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ForceNew:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},

				"ou_path": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ForceNew:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},

				"cluster_bootstrap_account_name": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ForceNew:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},

				"cluster_operator_account_name": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ForceNew:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},

				"sql_service_account_name": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ForceNew:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},

				"file_share_witness_path": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ForceNew:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},

				"storage_account_url": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ForceNew:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},

				"storage_account_primary_key": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					Sensitive:    true,
					ValidateFunc: validation.StringIsNotEmpty,
				},

				"cluster_subnet_type": {
					Type:     pluginsdk.TypeString,
					Required: true,
					ValidateFunc: validation.StringInSlice([]string{
						string(sqlvirtualmachinegroups.ClusterSubnetTypeMultiSubnet),
						string(sqlvirtualmachinegroups.ClusterSubnetTypeSingleSubnet),
					}, false),
				},
			},
		},
	}
}
