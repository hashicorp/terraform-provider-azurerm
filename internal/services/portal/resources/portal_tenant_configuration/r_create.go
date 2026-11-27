// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package portal_tenant_configuration

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/portal/2019-01-01-preview/tenantconfiguration"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	portalClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/portal/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/portal/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourcePortalTenantConfigurationCreateUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := portalClient.From(meta).TenantConfigurationsClient
	ctx, cancel := timeouts.ForCreateUpdate(account.StopContext(meta), d)
	defer cancel()

	// NOTE: we're using a Terraform-internal Resource ID here since the Go SDK no longer exposes one
	// since this is an operation on a Tenant (which doesn't expose any configurable values).
	id := parse.NewPortalTenantConfigurationID("default")
	if d.IsNewResource() {
		if !features.From(meta).SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
			existing, err := client.Get(ctx)
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for existing %s: %+v", id, err)
				}
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return tf.ImportAsExistsError("azurerm_portal_tenant_configuration", id.ID())
			}
		}
	}

	parameters := tenantconfiguration.Configuration{
		Properties: &tenantconfiguration.ConfigurationProperties{
			EnforcePrivateMarkdownStorage: pointer.To(d.Get("private_markdown_storage_enforced").(bool)),
		},
	}

	if _, err := client.Create(ctx, parameters); err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	if d.IsNewResource() {
		d.SetId(id.ID())
	}

	return resourcePortalTenantConfigurationRead(d, meta)
}
