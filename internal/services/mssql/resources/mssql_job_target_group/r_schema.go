// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package mssql_job_target_group

import (
	"context"
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
					return fmt.Errorf("`database_name` and `elastic_pool_name` are mutually exclusive")
				}
			}

			return nil
		},
	}
}
