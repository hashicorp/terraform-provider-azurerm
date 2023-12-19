package azurestackhci

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/azurestackhci/2023-08-01/arcsettings"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/azurestackhci/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

var _ sdk.DataSource = StackHCIClusterArcSettingDataSource{}

type StackHCIClusterArcSettingDataSource struct{}

func (r StackHCIClusterArcSettingDataSource) ResourceType() string {
	return "azurerm_stack_hci_cluster_arc_setting"
}

func (r StackHCIClusterArcSettingDataSource) ModelObject() interface{} {
	return &StackHCIClusterArcSettingDataSourceModel{}
}

type StackHCIClusterArcSettingDataSourceModel struct {
	Name              string `tfschema:"name"`
	StackHciClusterId string `tfschema:"stack_hci_cluster_id"`
}

func (r StackHCIClusterArcSettingDataSource) Arguments() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validate.ClusterName,
		},

		"stack_hci_cluster_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: arcsettings.ValidateClusterID,
		},
	}
}

func (r StackHCIClusterArcSettingDataSource) Attributes() map[string]*schema.Schema {
	return map[string]*schema.Schema{}
}

func (r StackHCIClusterArcSettingDataSource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.AzureStackHCI.ArcSettings

			var setting StackHCIClusterArcSettingDataSourceModel
			if err := metadata.Decode(&setting); err != nil {
				return err
			}

			clusterId, err := arcsettings.ParseClusterID(setting.StackHciClusterId)
			if err != nil {
				return err
			}

			id := arcsettings.NewArcSettingID(clusterId.SubscriptionId, clusterId.ResourceGroupName, clusterId.ClusterName, setting.Name)

			existing, err := client.ArcSettingsGet(ctx, id)
			if err != nil {
				if response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("%s was not found", id)
				}
				return fmt.Errorf("reading %s: %+v", id, err)
			}

			metadata.SetID(id)

			return metadata.Encode(&setting)
		},
	}
}
