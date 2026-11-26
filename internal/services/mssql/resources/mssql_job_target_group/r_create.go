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

func (r Resource) Create() sdk.ResourceFunc {
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

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.Get(ctx, id)
				if err != nil && !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}

				if !response.WasNotFound(existing.HttpResponse) {
					return metadata.ResourceRequiresImport(r.ResourceType(), id)
				}
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

func expandJobTargets(input []MsSqlJobTarget) []jobtargetgroups.JobTarget {
	targets := make([]jobtargetgroups.JobTarget, 0)
	if len(input) == 0 {
		return targets
	}

	for _, v := range input {
		t := jobtargetgroups.JobTarget{
			MembershipType: pointer.ToEnum[jobtargetgroups.JobTargetGroupMembershipType](v.MembershipType),
			ServerName:     pointer.To(v.ServerName),
		}

		targetType := determineJobTargetType(v)
		t.Type = targetType

		if v.JobCredentialId != "" {
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
