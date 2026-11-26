// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package data_factory_pipeline

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/datafactory/2018-06-01/pipelines"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceDataFactoryPipelineDelete(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).DataFactory.PipelinesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := pipelines.ParsePipelineID(d.Id())
	if err != nil {
		return err
	}

	if _, err = client.Delete(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return nil
}
