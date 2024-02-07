// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package paloalto

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/paloaltonetworks/2022-08-29/fqdnlistlocalrulestack"
	"github.com/hashicorp/go-azure-sdk/resource-manager/paloaltonetworks/2022-08-29/localrulestacks"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/paloalto/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type LocalRulestackFQDNList struct{}

var _ sdk.ResourceWithUpdate = LocalRulestackFQDNList{}

type LocalRulestackFQDNListModel struct {
	Name         string   `tfschema:"name"`
	RuleStackID  string   `tfschema:"rulestack_id"`
	FQDNList     []string `tfschema:"fully_qualified_domain_names"`
	AuditComment string   `tfschema:"audit_comment"`
	Description  string   `tfschema:"description"`
}

func (r LocalRulestackFQDNList) ModelObject() interface{} {
	return &LocalRulestackFQDNListModel{}
}

func (r LocalRulestackFQDNList) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return fqdnlistlocalrulestack.ValidateLocalRulestackFqdnListID
}

func (r LocalRulestackFQDNList) ResourceType() string {
	return "azurerm_palo_alto_local_rulestack_fqdn_list"
}

func (r LocalRulestackFQDNList) Arguments() map[string]*schema.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validate.LocalRuleStackFQDNListName, // TODO - Check this
		},

		"rulestack_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: fqdnlistlocalrulestack.ValidateLocalRulestackID,
		},

		"fully_qualified_domain_names": {
			Type:     pluginsdk.TypeList,
			Required: true,
			MinItems: 1,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: validation.StringIsNotEmpty,
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

func (r LocalRulestackFQDNList) Attributes() map[string]*schema.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r LocalRulestackFQDNList) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.PaloAlto.Client.FqdnListLocalRulestack
			rulestackClient := metadata.Client.PaloAlto.Client.LocalRulestacks

			model := LocalRulestackFQDNListModel{}

			if err := metadata.Decode(&model); err != nil {
				return err
			}

			rulestackId, err := localrulestacks.ParseLocalRulestackID(model.RuleStackID)
			if err != nil {
				return err
			}
			locks.ByID(rulestackId.ID())
			defer locks.UnlockByID(rulestackId.ID())

			id := fqdnlistlocalrulestack.NewLocalRulestackFqdnListID(rulestackId.SubscriptionId, rulestackId.ResourceGroupName, rulestackId.LocalRulestackName, model.Name)

			existing, err := client.Get(ctx, id)
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			props := fqdnlistlocalrulestack.FqdnObject{
				FqdnList: model.FQDNList,
			}

			if model.AuditComment != "" {
				props.AuditComment = pointer.To(model.AuditComment)
			}
			if model.Description != "" {
				props.Description = pointer.To(model.Description)
			}

			fqdnList := fqdnlistlocalrulestack.FqdnListLocalRulestackResource{
				Properties: props,
			}

			if _, err = client.CreateOrUpdate(ctx, id, fqdnList); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)

			if err = rulestackClient.CommitThenPoll(ctx, *rulestackId); err != nil {
				return fmt.Errorf("committing Local Rulestack config for %s: %+v", id, err)
			}

			return nil
		},
	}
}

func (r LocalRulestackFQDNList) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.PaloAlto.Client.FqdnListLocalRulestack

			id, err := fqdnlistlocalrulestack.ParseLocalRulestackFqdnListID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var state LocalRulestackFQDNListModel

			existing, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(existing.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("reading %s: %+v", *id, err)
			}

			state.Name = id.FqdnListName
			state.RuleStackID = fqdnlistlocalrulestack.NewLocalRulestackID(id.SubscriptionId, id.ResourceGroupName, id.LocalRulestackName).ID()

			if model := existing.Model; model != nil {
				props := model.Properties

				state.FQDNList = props.FqdnList
				state.AuditComment = pointer.From(props.AuditComment)
				state.Description = pointer.From(props.Description)
			}

			return metadata.Encode(&state)
		},
	}
}

func (r LocalRulestackFQDNList) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.PaloAlto.Client.FqdnListLocalRulestack
			rulestackClient := metadata.Client.PaloAlto.Client.LocalRulestacks

			id, err := fqdnlistlocalrulestack.ParseLocalRulestackFqdnListID(metadata.ResourceData.Id())
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

func (r LocalRulestackFQDNList) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.PaloAlto.Client.FqdnListLocalRulestack
			rulestackClient := metadata.Client.PaloAlto.Client.LocalRulestacks

			model := LocalRulestackFQDNListModel{}

			if err := metadata.Decode(&model); err != nil {
				return err
			}

			id, err := fqdnlistlocalrulestack.ParseLocalRulestackFqdnListID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			rulestackId := localrulestacks.NewLocalRulestackID(id.SubscriptionId, id.ResourceGroupName, id.LocalRulestackName)
			locks.ByID(rulestackId.ID())
			defer locks.UnlockByID(rulestackId.ID())

			existing, err := client.Get(ctx, *id)
			if err != nil {
				return fmt.Errorf("retreiving %s: %+v", *id, err)
			}

			fqdnList := *existing.Model

			if metadata.ResourceData.HasChange("fully_qualified_domain_names") {
				fqdnList.Properties.FqdnList = model.FQDNList
			}

			if metadata.ResourceData.HasChange("audit_comment") {
				fqdnList.Properties.AuditComment = pointer.To(model.AuditComment)
			}

			if metadata.ResourceData.HasChange("description") {
				fqdnList.Properties.Description = pointer.To(model.Description)
			}

			if _, err = client.CreateOrUpdate(ctx, *id, fqdnList); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			if err = rulestackClient.CommitThenPoll(ctx, rulestackId); err != nil {
				return fmt.Errorf("committing Local Rulestack config for %s: %+v", id, err)
			}

			return nil
		},
	}
}
