// Copyright IBM Corp. 2014, 2025
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
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
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
	Description     string `tfschema:"description"`
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

func (HorizondbFirewallRuleResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return firewallrules.ValidateFirewallRuleID
}

func (HorizondbFirewallRuleResource) Identity() resourceids.ResourceId {
	return &firewallrules.FirewallRuleId{}
}

func (HorizondbFirewallRuleResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type: pluginsdk.TypeString, Required: true, ForceNew: true,
			ValidateFunc: validation.All(
				validation.StringLenBetween(1, 80),
				validation.StringMatch(regexp.MustCompile(`^[a-zA-Z0-9][-a-zA-Z0-9]*[a-zA-Z0-9]$`), "must contain only alphanumeric characters and hyphens and start and end with an alphanumeric character"),
			),
		},
		"horizondb_pool_id": commonschema.ResourceIDReferenceRequiredForceNew(&firewallrules.PoolId{}),
		"end_ip_address": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.IsIPv4Address,
		},
		"start_ip_address": {
			Type:     pluginsdk.TypeString,
			Required: true, ValidateFunc: validation.IsIPv4Address,
		},
		"description": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringLenBetween(0, 1024),
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
			client := metadata.Client.HorizonDB.FirewallRulesClient
			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
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
			resp, err := metadata.Client.HorizonDB.FirewallRulesClient.Get(ctx, *id)
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
			if err := metadata.Client.HorizonDB.FirewallRulesClient.CreateOrReplaceThenPoll(ctx, *id, expandHorizondbFirewallRule(state)); err != nil {
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
			if err := metadata.Client.HorizonDB.FirewallRulesClient.DeleteThenPoll(ctx, *id); err != nil {
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
		Description:     pointer.From(model.Properties.Description),
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
			Description:    pointer.To(state.Description),
		},
	}
}
