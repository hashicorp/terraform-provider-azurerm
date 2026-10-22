// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package mssql_job_agent

import (
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2025-01-01/jobagents"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

//go:generate go run ../../../../tools/generator-tests resourceidentity -resource-name mssql_job_agent -service-package-name mssql -properties "name" -compare-values "resource_group_name:database_id,server_name:database_id"

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceMsSqlJobAgentCreate,
		Read:   resourceMsSqlJobAgentRead,
		Update: resourceMsSqlJobAgentUpdate,
		Delete: resourceMsSqlJobAgentDelete,

		Importer: pluginsdk.ImporterValidatingIdentity(&jobagents.JobAgentId{}),

		Identity: &schema.ResourceIdentity{
			SchemaFunc: pluginsdk.GenerateIdentitySchema(&jobagents.JobAgentId{}),
		},

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(60 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(60 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(60 * time.Minute),
		},

		Schema: mssqlJobAgentSchema(),
	}
}
