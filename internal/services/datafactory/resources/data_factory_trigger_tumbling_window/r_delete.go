// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package data_factory_trigger_tumbling_window

import (
	"fmt"

	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datafactory/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceDataFactoryTriggerTumblingWindowDelete(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).DataFactory.TriggersClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.TriggerID(d.Id())
	if err != nil {
		return err
	}

	future, err := client.Stop(ctx, id.ResourceGroup, id.FactoryName, id.Name)
	if err != nil {
		return fmt.Errorf("stopping %s: %+v", id, err)
	}
	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting to stop %s: %+v", id, err)
	}

	if _, err = client.Delete(ctx, id.ResourceGroup, id.FactoryName, id.Name); err != nil {
		return fmt.Errorf("deleting %s: %+v", id, err)
	}

	return nil
}
