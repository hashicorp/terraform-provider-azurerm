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
			Description: "The name which should be used for this Elastic Job Credential.",
			Type:        pluginsdk.TypeString,
			Required:    true,
			ForceNew:    true,
		},

		"job_agent_id": {
			Description:  "The ID of the Elastic Job Agent.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.AsGeneratedID(jobagents.ParseJobAgentIDInsensitively),
		},

		"username": {
			Description: "The username to use for this Elastic Job credential.",
			Type:        pluginsdk.TypeString,
			Required:    true,
		},

		"password": {
			Description:   "The password to use for this Elastic Job credential.",
			Type:          pluginsdk.TypeString,
			Optional:      true,
			Sensitive:     true,
			ConflictsWith: []string{"password_wo"},
			ExactlyOneOf:  []string{"password", "password_wo"},
		},
		"password_wo": {
			Description:   "The password to use for this Elastic Job credential.",
			Type:          pluginsdk.TypeString,
			Optional:      true,
			WriteOnly:     true,
			RequiredWith:  []string{"password_wo_version"},
			ConflictsWith: []string{"password"},
			ExactlyOneOf:  []string{"password_wo", "password"},
		},
		"password_wo_version": {
			Description:  "An integer value used to trigger an update for `password_wo`. This property should be incremented when updating `password_wo`.",
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			ValidateFunc: validation.IntAtLeast(1),
			RequiredWith: []string{"password_wo"},
		},
	}
}
