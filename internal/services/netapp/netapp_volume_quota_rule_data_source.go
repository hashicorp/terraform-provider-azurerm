// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package netapp

import (
	"context"
	"fmt"
	"math"
	"net/http"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/netapp/2022-05-01/volumequotarules"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	netAppModels "github.com/hashicorp/terraform-provider-azurerm/internal/services/netapp/models"
	netAppValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/netapp/validate"
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

		"resource_group_name": commonschema.ResourceGroupName(),

		"account_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: netAppValidate.AccountName,
		},

		"pool_name": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},

		"volume_name": {
			Type:     pluginsdk.TypeString,
			Required: true,
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

		"quota_size_in_mib": {
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

			id := volumequotarules.NewVolumeQuotaRuleID(metadata.Client.Account.SubscriptionId, state.ResourceGroupName, state.AccountName, state.CapacityPoolName, state.VolumeName, state.Name)

			resp, err := client.Get(ctx, id)
			if err != nil {
				if resp.HttpResponse.StatusCode == http.StatusNotFound {
					return fmt.Errorf("%s was not found", id)
				}
				return fmt.Errorf("retrieving %s: %v", id, err)
			}

			model := resp.Model
			if model == nil {
				return fmt.Errorf("retrieving %s: model was nil", id)
			}

			var quotaSizeInKiB int = 0
			var quotaSizeInMiB int = 0

			if pointer.From(model.Properties.QuotaSizeInKiBs) > math.MaxInt32 {
				quotaSizeInMiB = int(pointer.From(model.Properties.QuotaSizeInKiBs) / 1024)
			} else {
				quotaSizeInKiB = int(pointer.From(model.Properties.QuotaSizeInKiBs))
			}

			state.Location = location.Normalize(model.Location)
			state.QuotaSizeInKiB = quotaSizeInKiB
			state.QuotaSizeInMiB = quotaSizeInMiB
			state.QuotaTarget = pointer.From(model.Properties.QuotaTarget)
			state.QuotaType = string(pointer.From(model.Properties.QuotaType))

			metadata.SetID(id)

			return metadata.Encode(&model)
		},
	}
}
