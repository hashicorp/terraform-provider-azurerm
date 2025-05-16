// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package netapp

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/netapp/2025-01-01/volumequotarules"
	"github.com/hashicorp/go-azure-sdk/resource-manager/netapp/2025-01-01/volumes"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	netAppModels "github.com/hashicorp/terraform-provider-azurerm/internal/services/netapp/models"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type NetAppVolumeQuotaRuleDataSource struct{}

var _ sdk.DataSource = NetAppVolumeQuotaRuleDataSource{}

func (r NetAppVolumeQuotaRuleDataSource) ResourceType() string {
	return "azurerm_netapp_volume_quota_rule"
}

func (r NetAppVolumeQuotaRuleDataSource) ModelObject() interface{} {
	return &netAppModels.NetAppVolumeQuotaRuleDataSourceModel{}
}

func (r NetAppVolumeQuotaRuleDataSource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return volumequotarules.ValidateVolumeQuotaRuleID
}

func (r NetAppVolumeQuotaRuleDataSource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},

		"volume_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: volumequotarules.ValidateVolumeID,
		},
	}
}

func (r NetAppVolumeQuotaRuleDataSource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"location": commonschema.LocationComputed(),

		"quota_target": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"quota_size_in_kib": {
			Type:     pluginsdk.TypeInt,
			Computed: true,
		},

		"quota_type": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}

func (r NetAppVolumeQuotaRuleDataSource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.NetApp.VolumeQuotaRules

			var state netAppModels.NetAppVolumeQuotaRuleDataSourceModel
			if err := metadata.Decode(&state); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			volumeID, err := volumes.ParseVolumeID(state.VolumeID)
			if err != nil {
				return fmt.Errorf("error parsing volume id %s: %+v", state.VolumeID, err)
			}

			id := volumequotarules.NewVolumeQuotaRuleID(metadata.Client.Account.SubscriptionId, volumeID.ResourceGroupName, volumeID.NetAppAccountName, volumeID.CapacityPoolName, volumeID.VolumeName, state.Name)

			resp, err := client.Get(ctx, id)
			if err != nil {
				if resp.HttpResponse.StatusCode == http.StatusNotFound {
					return fmt.Errorf("%s was not found", id)
				}
				return fmt.Errorf("retrieving %s: %v", id, err)
			}

			model := resp.Model
			if model == nil || model.Properties == nil {
				return fmt.Errorf("retrieving %s: model was nil", id)
			}

			state.Location = location.Normalize(model.Location)
			state.QuotaSizeInKiB = pointer.From(model.Properties.QuotaSizeInKiBs)
			state.QuotaTarget = pointer.From(model.Properties.QuotaTarget)
			state.QuotaType = string(pointer.From(model.Properties.QuotaType))

			metadata.SetID(id)

			return metadata.Encode(&state)
		},
	}
}
