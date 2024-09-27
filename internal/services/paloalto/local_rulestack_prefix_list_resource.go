// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package paloalto

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/paloaltonetworks/2022-08-29/localrulestacks"
	"github.com/hashicorp/go-azure-sdk/resource-manager/paloaltonetworks/2022-08-29/prefixlistlocalrulestack"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/paloalto/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type LocalRuleStackPrefixList struct{}

var _ sdk.ResourceWithUpdate = LocalRuleStackPrefixList{}

type LocalRuleStackPrefixListModel struct {
	Name         string   `tfschema:"name"`
	RuleStackID  string   `tfschema:"rulestack_id"`
	PrefixList   []string `tfschema:"prefix_list"`
	AuditComment string   `tfschema:"audit_comment"`
	Description  string   `tfschema:"description"`
}

func (r LocalRuleStackPrefixList) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return prefixlistlocalrulestack.ValidateLocalRulestackPrefixListID
}

func (r LocalRuleStackPrefixList) ResourceType() string {
	return "azurerm_palo_alto_local_rulestack_prefix_list"
}

func (r LocalRuleStackPrefixList) ModelObject() interface{} {
	return &LocalRuleStackPrefixListModel{}
}

func (r LocalRuleStackPrefixList) Arguments() map[string]*schema.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validate.LocalRuleStackRuleName, // TODO - Check this
		},

		"rulestack_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: prefixlistlocalrulestack.ValidateLocalRulestackID,
		},

		"prefix_list": {
			Type:     pluginsdk.TypeList,
			Required: true,
			MinItems: 1,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: validation.IsCIDR,
			},
		},

		"audit_comment": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},

		"description": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},
	}
}

func (r LocalRuleStackPrefixList) Attributes() map[string]*schema.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r LocalRuleStackPrefixList) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.PaloAlto.Client.PrefixListLocalRulestack
			rulestackClient := metadata.Client.PaloAlto.Client.LocalRulestacks
			model := LocalRuleStackPrefixListModel{}

			if err := metadata.Decode(&model); err != nil {
				return err
			}

			rulestackId, err := localrulestacks.ParseLocalRulestackID(model.RuleStackID)
			if err != nil {
				return err
			}
			locks.ByID(rulestackId.ID())
			defer locks.UnlockByID(rulestackId.ID())

			id := prefixlistlocalrulestack.NewLocalRulestackPrefixListID(rulestackId.SubscriptionId, rulestackId.ResourceGroupName, rulestackId.LocalRulestackName, model.Name)

			existing, err := client.Get(ctx, id)
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			props := prefixlistlocalrulestack.PrefixObject{
				PrefixList: model.PrefixList,
			}

			if model.AuditComment != "" {
				props.AuditComment = pointer.To(model.AuditComment)
			}

			if model.Description != "" {
				props.Description = pointer.To(model.Description)
			}

			prefixList := prefixlistlocalrulestack.PrefixListResource{
				Properties: props,
			}

			if err = client.CreateOrUpdateThenPoll(ctx, id, prefixList); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)

			if err = rulestackClient.CommitThenPoll(ctx, *rulestackId); err != nil {
				return fmt.Errorf("committing Local RuleStack config for %s: %+v", id, err)
			}

			return nil
		},
	}
}

func (r LocalRuleStackPrefixList) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.PaloAlto.Client.PrefixListLocalRulestack

			id, err := prefixlistlocalrulestack.ParseLocalRulestackPrefixListID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var state LocalRuleStackPrefixListModel

			existing, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(existing.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("reading %s: %+v", *id, err)
			}

			state.Name = id.PrefixListName
			state.RuleStackID = prefixlistlocalrulestack.NewLocalRulestackID(id.SubscriptionId, id.ResourceGroupName, id.LocalRulestackName).ID()
			if model := existing.Model; model != nil {
				props := model.Properties

				state.PrefixList = props.PrefixList
				state.AuditComment = pointer.From(props.AuditComment)
				state.Description = pointer.From(props.Description)
			}

			return metadata.Encode(&state)
		},
	}
}

func (r LocalRuleStackPrefixList) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.PaloAlto.Client.PrefixListLocalRulestack
			rulestackClient := metadata.Client.PaloAlto.Client.LocalRulestacks

			id, err := prefixlistlocalrulestack.ParseLocalRulestackPrefixListID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			rulestackId := localrulestacks.NewLocalRulestackID(id.SubscriptionId, id.ResourceGroupName, id.LocalRulestackName)
			locks.ByID(rulestackId.ID())
			defer locks.UnlockByID(rulestackId.ID())

			if err = client.DeleteThenPoll(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}

			if err = rulestackClient.CommitThenPoll(ctx, rulestackId); err != nil {
				return fmt.Errorf("committing Local Rulestack config for %s: %+v", id, err)
			}

			return nil
		},
	}
}

func (r LocalRuleStackPrefixList) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.PaloAlto.Client.PrefixListLocalRulestack
			rulestackClient := metadata.Client.PaloAlto.Client.LocalRulestacks

			id, err := prefixlistlocalrulestack.ParseLocalRulestackPrefixListID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			model := LocalRuleStackPrefixListModel{}

			if err = metadata.Decode(&model); err != nil {
				return err
			}

			rulestackId := localrulestacks.NewLocalRulestackID(id.SubscriptionId, id.ResourceGroupName, id.LocalRulestackName)
			locks.ByID(rulestackId.ID())
			defer locks.UnlockByID(rulestackId.ID())

			existing, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(existing.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("reading %s for update: %+v", *id, err)
			}

			prefixList := *existing.Model

			if metadata.ResourceData.HasChange("prefix_list") {
				prefixList.Properties.PrefixList = model.PrefixList
			}

			if metadata.ResourceData.HasChange("audit_comment") {
				prefixList.Properties.AuditComment = pointer.To(model.AuditComment)
			}

			if metadata.ResourceData.HasChange("description") {
				prefixList.Properties.Description = pointer.To(model.Description)
			}

			if _, err = client.CreateOrUpdate(ctx, *id, prefixList); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			if err = rulestackClient.CommitThenPoll(ctx, rulestackId); err != nil {
				return fmt.Errorf("committing Local Rulestack config for %s: %+v", id, err)
			}

			return nil
		},
	}
}
