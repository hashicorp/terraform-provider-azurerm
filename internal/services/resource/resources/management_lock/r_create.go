// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package management_lock

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/resources/2020-05-01/managementlocks"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/custompollers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceManagementLockCreate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Resource.LocksClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := managementlocks.NewScopedLockID(d.Get("scope").(string), d.Get("name").(string))

	if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		existing, err := client.GetByScope(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_management_lock", id.ID())
		}
	}

	payload := managementlocks.ManagementLockObject{
		Properties: managementlocks.ManagementLockProperties{
			Level: managementlocks.LockLevel(d.Get("lock_level").(string)),
			Notes: pointer.To(d.Get("notes").(string)),
		},
	}

	if _, err := client.CreateOrUpdateByScope(ctx, id, payload); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	poller := custompollers.NewEventualConsistencyPoller(12, func(pollerCtx context.Context) (*http.Response, error) {
		resp, err := client.GetByScope(pollerCtx, id)
		return resp.HttpResponse, err
	}, custompollers.DefaultCreationEventualConsistencyPollerOptions())
	if err := poller.PollUntilDone(ctx); err != nil {
		return fmt.Errorf("waiting for %s to finish create replication", id)
	}

	d.SetId(id.ID())
	return resourceManagementLockRead(d, meta)
}
