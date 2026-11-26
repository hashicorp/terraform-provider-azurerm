// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package mysql_flexible_database

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/mysql/2023-12-30/databases"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceMySqlFlexibleDatabaseRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).MySQL.FlexibleServers.Databases
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := databases.ParseDatabaseID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	return resourceMySqlFlexibleDatabaseFlatten(d, id, resp.Model)
}

// `utf8` is a deprecated alias for `utf8mb3`, see http://dev.mysql.com/doc/refman/8.4/en/charset-unicode-utf8.html
func resourceMysqlFlexibleDatabaseAliasDiffSuppress(firstAlias, secondAlias string) pluginsdk.SchemaDiffSuppressFunc {
	return func(key, oldValue, newValue string, d *pluginsdk.ResourceData) bool {
		return suppress.CaseDifference(key, oldValue, newValue, d) ||
			strings.EqualFold(oldValue, firstAlias) && strings.EqualFold(newValue, secondAlias) ||
			strings.EqualFold(oldValue, secondAlias) && strings.EqualFold(newValue, firstAlias)
	}
}
