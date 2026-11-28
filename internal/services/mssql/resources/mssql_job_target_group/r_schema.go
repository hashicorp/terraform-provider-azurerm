// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package mssql_job_target_group

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2025-01-01/jobcredentials"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2025-01-01/jobs"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2025-01-01/jobtargetgroups"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/mssql/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func (r Resource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description:  "The name which should be used for this Job Target Group.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
		"job_agent_id": {
			Description:  "The ID of the Elastic Job Agent.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: jobs.ValidateJobAgentID,
		},
		"job_target": {
			Type:     pluginsdk.TypeSet,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"server_name": {
						Description:  "The name of the MS SQL Server.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validate.ValidateMsSqlServerName,
					},
					"database_name": {
						Description:  "The name of the MS SQL Database.",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validate.ValidateMsSqlDatabaseName,
					},
					"elastic_pool_name": {
						Description:  "The name of the MS SQL Elastic Pool.",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validate.ValidateMsSqlElasticPoolName,
					},
					"job_credential_id": {
						Description:  "The ID of the job credential to use during execution of jobs.",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: jobcredentials.ValidateCredentialID,
					},
					"membership_type": {
						Description:  "The membership type for this job target.",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Default:      string(jobtargetgroups.JobTargetGroupMembershipTypeInclude),
						ValidateFunc: validation.StringInSlice(jobtargetgroups.PossibleValuesForJobTargetGroupMembershipType(), false),
					},
					"type": {
						Description: "The job target type. This value is computed based on `server_name`, `database_name`, and `elastic_pool_name`.",
						Type:        pluginsdk.TypeString,
						Computed:    true,
					},
				},
			},
		},
	}
}

func (r Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r Resource) CustomizeDiff() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 10 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var config MsSqlJobTargetGroupResourceModel
			if err := metadata.DecodeDiff(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			for _, v := range config.JobTargets {
				if v.DatabaseName != "" && v.ElasticPoolName != "" {
					return errors.New("`database_name` and `elastic_pool_name` are mutually exclusive")
				}
			}

			return nil
		},
	}
}
