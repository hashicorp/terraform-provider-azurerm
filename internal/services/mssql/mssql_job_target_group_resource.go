// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package mssql

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2023-08-01-preview/jobcredentials"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2023-08-01-preview/jobs"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2023-08-01-preview/jobtargetgroups"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/mssql/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type MsSqlJobTargetGroupResource struct{}

type MsSqlJobTargetGroupResourceModel struct {
	Name       string           `tfschema:"name"`
	JobAgentID string           `tfschema:"job_agent_id"`
	JobTargets []MsSqlJobTarget `tfschema:"job_target"`
}

type MsSqlJobTarget struct {
	ServerName      string `tfschema:"server_name"`
	Type            string `tfschema:"type"`
	DatabaseName    string `tfschema:"database_name"`
	ElasticPoolName string `tfschema:"elastic_pool_name"`
	JobCredentialId string `tfschema:"job_credential_id"`
	MembershipType  string `tfschema:"membership_type"`
}

var (
	_ sdk.ResourceWithUpdate        = MsSqlJobTargetGroupResource{}
	_ sdk.ResourceWithCustomizeDiff = MsSqlJobTargetGroupResource{}
)

func (r MsSqlJobTargetGroupResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
			ForceNew:     true,
		},
		"job_agent_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: jobs.ValidateJobAgentID,
			ForceNew:     true,
		},
		"job_target": {
			Type:     pluginsdk.TypeSet,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"server_name": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validate.ValidateMsSqlServerName,
					},
					"database_name": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validate.ValidateMsSqlDatabaseName,
					},
					"elastic_pool_name": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validate.ValidateMsSqlElasticPoolName,
					},
					"job_credential_id": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: jobcredentials.ValidateCredentialID,
					},
					"membership_type": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Default:      string(jobtargetgroups.JobTargetGroupMembershipTypeInclude),
						ValidateFunc: validation.StringInSlice(jobtargetgroups.PossibleValuesForJobTargetGroupMembershipType(), false),
					},
					"type": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
				},
			},
		},
	}
}

func (r MsSqlJobTargetGroupResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r MsSqlJobTargetGroupResource) CustomizeDiff() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 10 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var config MsSqlJobTargetGroupResourceModel
			if err := metadata.DecodeDiff(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			for _, v := range config.JobTargets {
				if v.DatabaseName != "" && v.ElasticPoolName != "" {
					return fmt.Errorf("`database_name` and `elastic_pool_name` are mutually exclusive")
				}

				targetType := determineJobTargetType(v)
				if isCredentialRequired(jobtargetgroups.JobTargetGroupMembershipType(v.MembershipType), targetType) {
					if v.JobCredentialId == "" {
						return fmt.Errorf("`job_credential_id` is required when `membership_type` is `%s` and `type` is `%s`", jobtargetgroups.JobTargetGroupMembershipTypeInclude, targetType)
					}
				}
			}

			return nil
		},
	}
}

func (r MsSqlJobTargetGroupResource) ModelObject() interface{} {
	return &MsSqlJobTargetGroupResourceModel{}
}

func (r MsSqlJobTargetGroupResource) ResourceType() string {
	return "azurerm_mssql_job_target_group"
}

func (r MsSqlJobTargetGroupResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.MSSQL.JobTargetGroupsClient

			var model MsSqlJobTargetGroupResourceModel
			if err := metadata.Decode(&model); err != nil {
				return err
			}

			jobAgent, err := jobtargetgroups.ParseJobAgentID(model.JobAgentID)
			if err != nil {
				return err
			}

			id := jobtargetgroups.NewTargetGroupID(jobAgent.SubscriptionId, jobAgent.ResourceGroupName, jobAgent.ServerName, jobAgent.JobAgentName, model.Name)

			existing, err := client.Get(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			parameters := jobtargetgroups.JobTargetGroup{
				Name: pointer.To(model.Name),
				Properties: pointer.To(jobtargetgroups.JobTargetGroupProperties{
					Members: expandJobTargets(model.JobTargets),
				}),
			}

			if _, err := client.CreateOrUpdate(ctx, id, parameters); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r MsSqlJobTargetGroupResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.MSSQL.JobTargetGroupsClient

			id, err := jobtargetgroups.ParseTargetGroupID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}

				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			state := MsSqlJobTargetGroupResourceModel{
				Name:       id.TargetGroupName,
				JobAgentID: jobtargetgroups.NewJobAgentID(id.SubscriptionId, id.ResourceGroupName, id.ServerName, id.JobAgentName).ID(),
			}

			if model := resp.Model; model != nil {
				if props := model.Properties; props != nil {
					state.JobTargets = flattenJobTargets(props.Members)
				}
			}

			return metadata.Encode(&state)
		},
	}
}

func (r MsSqlJobTargetGroupResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.MSSQL.JobTargetGroupsClient

			id, err := jobtargetgroups.ParseTargetGroupID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var config MsSqlJobTargetGroupResourceModel
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding %+v", err)
			}

			existing, err := client.Get(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			if existing.Model == nil {
				return fmt.Errorf("retrieving %s: `model` was nil", id)
			}

			if existing.Model.Properties == nil {
				return fmt.Errorf("retrieving %s: `model.Properties` was nil", id)
			}

			payload := existing.Model

			if metadata.ResourceData.HasChange("job_target") {
				payload.Properties.Members = expandJobTargets(config.JobTargets)
			}

			if _, err := client.CreateOrUpdate(ctx, *id, *payload); err != nil {
				return fmt.Errorf("updating: %s: %+v", id, err)
			}

			return nil
		},
	}
}

func (r MsSqlJobTargetGroupResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.MSSQL.JobTargetGroupsClient

			id, err := jobtargetgroups.ParseTargetGroupID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if _, err := client.Delete(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", id, err)
			}

			return nil
		},
	}
}

func (r MsSqlJobTargetGroupResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return jobtargetgroups.ValidateTargetGroupID
}

func expandJobTargets(input []MsSqlJobTarget) []jobtargetgroups.JobTarget {
	targets := make([]jobtargetgroups.JobTarget, 0)
	if len(input) == 0 {
		return targets
	}

	for _, v := range input {
		t := jobtargetgroups.JobTarget{
			MembershipType: pointer.To(jobtargetgroups.JobTargetGroupMembershipType(v.MembershipType)),
			ServerName:     pointer.To(v.ServerName),
		}

		targetType := determineJobTargetType(v)
		t.Type = targetType

		if isCredentialRequired(jobtargetgroups.JobTargetGroupMembershipType(v.MembershipType), targetType) {
			t.RefreshCredential = pointer.To(v.JobCredentialId)
		}

		if targetType == jobtargetgroups.JobTargetTypeSqlDatabase {
			t.DatabaseName = pointer.To(v.DatabaseName)
		}

		if targetType == jobtargetgroups.JobTargetTypeSqlElasticPool {
			t.ElasticPoolName = pointer.To(v.ElasticPoolName)
		}

		targets = append(targets, t)
	}

	return targets
}

func flattenJobTargets(input []jobtargetgroups.JobTarget) []MsSqlJobTarget {
	targets := make([]MsSqlJobTarget, 0)
	if len(input) == 0 {
		return targets
	}

	for _, v := range input {
		t := MsSqlJobTarget{
			DatabaseName:    pointer.From(v.DatabaseName),
			ElasticPoolName: pointer.From(v.ElasticPoolName),
			MembershipType:  string(pointer.From(v.MembershipType)),
			JobCredentialId: pointer.From(v.RefreshCredential),
			ServerName:      pointer.From(v.ServerName),
			Type:            string(v.Type),
		}

		targets = append(targets, t)
	}

	return targets
}

func determineJobTargetType(input MsSqlJobTarget) jobtargetgroups.JobTargetType {
	switch {
	case input.DatabaseName != "":
		return jobtargetgroups.JobTargetTypeSqlDatabase
	case input.ElasticPoolName != "":
		return jobtargetgroups.JobTargetTypeSqlElasticPool
	default:
		return jobtargetgroups.JobTargetTypeSqlServer
	}
}

func isCredentialRequired(membershipType jobtargetgroups.JobTargetGroupMembershipType, targetType jobtargetgroups.JobTargetType) bool {
	return membershipType == jobtargetgroups.JobTargetGroupMembershipTypeInclude && targetType != jobtargetgroups.JobTargetTypeSqlDatabase
}
