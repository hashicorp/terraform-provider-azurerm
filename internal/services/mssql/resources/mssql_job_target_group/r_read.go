// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package mssql_job_target_group

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2025-01-01/jobtargetgroups"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Read() sdk.ResourceFunc {
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

func flattenJobTargets(input []jobtargetgroups.JobTarget) []MsSqlJobTarget {
	targets := make([]MsSqlJobTarget, 0)
	if len(input) == 0 {
		return targets
	}

	for _, v := range input {
		t := MsSqlJobTarget{
			DatabaseName:    pointer.From(v.DatabaseName),
			ElasticPoolName: pointer.From(v.ElasticPoolName),
			MembershipType:  pointer.FromEnum(v.MembershipType),
			JobCredentialId: pointer.From(v.RefreshCredential),
			ServerName:      pointer.From(v.ServerName),
			Type:            string(v.Type),
		}

		targets = append(targets, t)
	}

	return targets
}

type MsSqlJobTarget struct {
	ServerName      string `tfschema:"server_name"`
	Type            string `tfschema:"type"`
	DatabaseName    string `tfschema:"database_name"`
	ElasticPoolName string `tfschema:"elastic_pool_name"`
	JobCredentialId string `tfschema:"job_credential_id"`
	MembershipType  string `tfschema:"membership_type"`
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
