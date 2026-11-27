// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package horizondb

import (
	"context"
	"errors"
	"fmt"
	"net/netip"
	"regexp"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/horizondb/2026-10-01-preview/firewallrules"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	horizondbClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/horizondb/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

//go:generate go run ../../tools/generator-tests resourceidentity -properties "name" -compare-values "subscription_id:horizondb_pool_id,resource_group_name:horizondb_pool_id,cluster_name:horizondb_pool_id,pool_name:horizondb_pool_id"
type HorizondbFirewallRuleResource struct{}

type HorizondbFirewallRuleResourceModel struct {
	Name            string `tfschema:"name"`
	HorizondbPoolId string `tfschema:"horizondb_pool_id"`
	EndIpAddress    string `tfschema:"end_ip_address"`
	StartIpAddress  string `tfschema:"start_ip_address"`
}

var (
	_ sdk.ResourceWithUpdate        = HorizondbFirewallRuleResource{}
	_ sdk.ResourceWithIdentity      = HorizondbFirewallRuleResource{}
	_ sdk.ResourceWithCustomizeDiff = HorizondbFirewallRuleResource{}
)

func (HorizondbFirewallRuleResource) ModelObject() any {
	return &HorizondbFirewallRuleResourceModel{}
}

func (HorizondbFirewallRuleResource) ResourceType() string {
	return "azurerm_horizondb_firewall_rule"
}

func (HorizondbFirewallRuleResource) Description() string {
	return "Manages a HorizonDB Firewall Rule."
}

func (HorizondbFirewallRuleResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return firewallrules.ValidateFirewallRuleID
}

func (HorizondbFirewallRuleResource) Identity() resourceids.ResourceId {
	return &firewallrules.FirewallRuleId{}
}

func (HorizondbFirewallRuleResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description: "The name of the HorizonDB Firewall Rule.",
			Type:        pluginsdk.TypeString,
			Required:    true,
			ForceNew:    true,
			ValidateFunc: validation.All(
				validation.StringLenBetween(1, 80),
				validation.StringMatch(regexp.MustCompile(`^[a-zA-Z0-9][-a-zA-Z0-9]*[a-zA-Z0-9]$`), "must contain only alphanumeric characters and hyphens and start and end with an alphanumeric character"),
			),
		},

		"horizondb_pool_id": commonschema.ResourceIDReferenceRequiredForceNew(&firewallrules.PoolId{}, "HorizonDB Pool"),

		"end_ip_address": {
			Description:  "The last IPv4 address in the allowed range.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.IsIPv4Address,
		},

		"start_ip_address": {
			Description:  "The first IPv4 address in the allowed range.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.IsIPv4Address,
		},
	}
}

func (HorizondbFirewallRuleResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (HorizondbFirewallRuleResource) CustomizeDiff() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			diff := metadata.ResourceDiff
			if !diff.NewValueKnown("start_ip_address") || !diff.NewValueKnown("end_ip_address") {
				return nil
			}

			startAddress, _ := netip.ParseAddr(diff.Get("start_ip_address").(string))
			endAddress, _ := netip.ParseAddr(diff.Get("end_ip_address").(string))
			if endAddress.Less(startAddress) {
				return errors.New("`end_ip_address` must be greater than or equal to `start_ip_address`")
			}
			return nil
		},
	}
}

func (r HorizondbFirewallRuleResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var state HorizondbFirewallRuleResourceModel
			if err := metadata.Decode(&state); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			poolId, err := firewallrules.ParsePoolID(state.HorizondbPoolId)
			if err != nil {
				return err
			}

			id := firewallrules.NewFirewallRuleID(poolId.SubscriptionId, poolId.ResourceGroupName, poolId.ClusterName, poolId.PoolName, state.Name)
			client := horizondbClient.From(metadata.Client).FirewallRulesClient
			if !features.From(metadata.Client).SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.Get(ctx, id)
				if err != nil && !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}
				if !response.WasNotFound(existing.HttpResponse) {
					return metadata.ResourceRequiresImport(r.ResourceType(), id)
				}
			}

			if err := client.CreateOrReplaceCallbackThenPoll(ctx, id, expandHorizondbFirewallRule(state), metadata.SetIDAndIdentityCallback(&id)); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)

			return pluginsdk.SetResourceIdentityData(metadata.ResourceData, &id)
		},
	}
}

func (r HorizondbFirewallRuleResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			id, err := firewallrules.ParseFirewallRuleIDInsensitively(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := horizondbClient.From(metadata.Client).FirewallRulesClient.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(*id)
				}
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			return r.flatten(metadata, id, resp.Model)
		},
	}
}

func (HorizondbFirewallRuleResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			id, err := firewallrules.ParseFirewallRuleIDInsensitively(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var state HorizondbFirewallRuleResourceModel
			if err := metadata.Decode(&state); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			client := horizondbClient.From(metadata.Client).FirewallRulesClient
			existing, err := client.Get(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			if existing.Model == nil {
				return fmt.Errorf("retrieving %s: `model` was nil", *id)
			}

			if existing.Model.Properties == nil {
				return fmt.Errorf("retrieving %s: `properties` was nil", *id)
			}

			payload := pointer.From(existing.Model)

			if metadata.ResourceData.HasChange("start_ip_address") {
				payload.Properties.StartIPAddress = state.StartIpAddress
			}

			if metadata.ResourceData.HasChange("end_ip_address") {
				payload.Properties.EndIPAddress = state.EndIpAddress
			}

			if err := client.CreateOrReplaceThenPoll(ctx, *id, payload); err != nil {
				return fmt.Errorf("updating %s: %+v", id, err)
			}

			return nil
		},
	}
}

func (HorizondbFirewallRuleResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			id, err := firewallrules.ParseFirewallRuleIDInsensitively(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if err := horizondbClient.From(metadata.Client).FirewallRulesClient.DeleteThenPoll(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", id, err)
			}

			return nil
		},
	}
}

func (HorizondbFirewallRuleResource) flatten(metadata sdk.ResourceMetaData, id *firewallrules.FirewallRuleId, model *firewallrules.FirewallRule) error {
	if model == nil || model.Properties == nil {
		return fmt.Errorf("retrieving %s: firewall rule properties are missing", id)
	}

	state := HorizondbFirewallRuleResourceModel{
		Name:            id.FirewallRuleName,
		HorizondbPoolId: firewallrules.NewPoolID(id.SubscriptionId, id.ResourceGroupName, id.ClusterName, id.PoolName).ID(),
		StartIpAddress:  model.Properties.StartIPAddress,
		EndIpAddress:    model.Properties.EndIPAddress,
	}

	if err := pluginsdk.SetResourceIdentityData(metadata.ResourceData, id); err != nil {
		return err
	}

	return metadata.Encode(&state)
}

func expandHorizondbFirewallRule(state HorizondbFirewallRuleResourceModel) firewallrules.FirewallRule {
	return firewallrules.FirewallRule{
		Properties: &firewallrules.FirewallRuleProperties{
			StartIPAddress: state.StartIpAddress,
			EndIPAddress:   state.EndIpAddress,
		},
	}
}
