// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package mssql_job_credential

import (
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2025-01-01/jobcredentials"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceMsSqlJobCredentialCreate,
		Read:   resourceMsSqlJobCredentialRead,
		Update: resourceMsSqlJobCredentialUpdate,
		Delete: resourceMsSqlJobCredentialDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := jobcredentials.ParseCredentialID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(60 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(60 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(60 * time.Minute),
		},

		Schema: mssqlJobCredentialSchema(),
	}
}
