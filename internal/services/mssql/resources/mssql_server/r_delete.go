// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package mssql_server

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	mssqlClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/mssql/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceMsSqlServerDelete(d *pluginsdk.ResourceData, meta any) error {
	client := mssqlClient.From(meta).ServersClient
	ctx, cancel := timeouts.ForDelete(account.StopContext(meta), d)
	defer cancel()

	id, err := commonids.ParseSqlServerID(d.Id())
	if err != nil {
		return err
	}

	if err = client.DeleteThenPoll(ctx, pointer.From(id)); err != nil {
		return fmt.Errorf("deleting SQL Server %s: %+v", id, err)
	}

	return nil
}
