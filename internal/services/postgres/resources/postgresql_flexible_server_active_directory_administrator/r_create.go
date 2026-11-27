// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package postgresql_flexible_server_active_directory_administrator

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/postgresql/2026-10-01/administratormicrosoftentras"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	postgresClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/postgres/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourcePostgresqlFlexibleServerAdministratorCreate(d *pluginsdk.ResourceData, meta any) error {
	client := postgresClient.From(meta).FlexibleServerAdministratorsClient
	subscriptionId := account.From(meta).SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(account.StopContext(meta), d)
	defer cancel()

	id := administratormicrosoftentras.NewAdministratorID(subscriptionId, d.Get("resource_group_name").(string), d.Get("server_name").(string), d.Get("object_id").(string))

	locks.ByName(id.FlexibleServerName, postgresqlFlexibleServerResourceName)
	defer locks.UnlockByName(id.FlexibleServerName, postgresqlFlexibleServerResourceName)

	if !features.From(meta).SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		existing, err := client.AdministratorsMicrosoftEntraGet(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_postgresql_flexible_server_active_directory_administrator", id.ID())
		}
	}

	principalName := d.Get("principal_name").(string)
	principalType := administratormicrosoftentras.PrincipalType(d.Get("principal_type").(string))
	tenantId := d.Get("tenant_id").(string)
	parameters := administratormicrosoftentras.AdministratorMicrosoftEntraAdd{
		Properties: &administratormicrosoftentras.AdministratorMicrosoftEntraPropertiesForAdd{
			PrincipalName: &principalName,
			PrincipalType: &principalType,
			TenantId:      &tenantId,
		},
	}

	if err := client.AdministratorsMicrosoftEntraCreateOrUpdateCallbackThenPoll(ctx, id, parameters, sdk.SetIDCallback(meta, &id, d)); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}
	d.SetId(id.ID())

	return nil
}
