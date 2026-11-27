// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package storage_mover_job_definition

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/storagemover/2025-07-01/jobdefinitions"
)

type StorageMoverJobDefinitionResourceModel struct {
	Name                  string                  `tfschema:"name"`
	StorageMoverProjectId string                  `tfschema:"storage_mover_project_id"`
	SourceName            string                  `tfschema:"source_name"`
	TargetName            string                  `tfschema:"target_name"`
	CopyMode              jobdefinitions.CopyMode `tfschema:"copy_mode"`
	SourceSubpath         string                  `tfschema:"source_sub_path"`
	TargetSubpath         string                  `tfschema:"target_sub_path"`
	AgentName             string                  `tfschema:"agent_name"`
	Description           string                  `tfschema:"description"`
}

func (r Resource) ModelObject() any {
	return &StorageMoverJobDefinitionResourceModel{}
}
