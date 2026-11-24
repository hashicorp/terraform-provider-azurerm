// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package mssql_job_credential

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2025-01-01/jobagents"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func mssqlJobCredentialSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"job_agent_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.AsGeneratedID(jobagents.ParseJobAgentIDInsensitively),
		},

		"username": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},

		"password": {
			Type:          pluginsdk.TypeString,
			Optional:      true,
			Sensitive:     true,
			ConflictsWith: []string{"password_wo"},
			ExactlyOneOf:  []string{"password", "password_wo"},
		},
		"password_wo": {
			Type:          pluginsdk.TypeString,
			Optional:      true,
			WriteOnly:     true,
			RequiredWith:  []string{"password_wo_version"},
			ConflictsWith: []string{"password"},
			ExactlyOneOf:  []string{"password_wo", "password"},
		},
		"password_wo_version": {
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			RequiredWith: []string{"password_wo"},
		},
	}
}
