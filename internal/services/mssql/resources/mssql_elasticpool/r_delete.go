// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package mssql_elasticpool

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceMsSqlElasticPoolDelete(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).MSSQL.ElasticPoolsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := commonids.ParseSqlElasticPoolID(d.Id())
	if err != nil {
		return err
	}

	if err = client.DeleteThenPoll(ctx, *id); err != nil {
		return fmt.Errorf("deleting ElasticPool %s: %+v", id, err)
	}

	return nil
}
