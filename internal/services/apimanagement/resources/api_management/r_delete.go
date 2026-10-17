// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package api_management

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/deletedservice"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/apimanagementservice"
	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/custompollers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceApiManagementServiceDelete(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).ApiManagement.ServiceClient
	deletedServicesClient := meta.(*clients.Client).ApiManagement.DeletedServicesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := apimanagementservice.ParseServiceID(d.Id())
	if err != nil {
		return err
	}

	existing, err := client.Get(ctx, *id)
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	resp, err := client.Delete(ctx, *id)
	if err != nil {
		return fmt.Errorf("deleting %s: %v", *id, err)
	}

	isLroStatus := resp.HttpResponse.StatusCode == http.StatusCreated || resp.HttpResponse.StatusCode == http.StatusAccepted
	if isLroStatus {
		pollerType, err := custompollers.NewAPIManagementPoller(client, resp.HttpResponse)
		if err != nil {
			return fmt.Errorf("polling deleting %s: %+v", id, err)
		}
		if pollerType != nil {
			poller := pollers.NewPoller(pollerType, 20*time.Second, pollers.DefaultNumberOfDroppedConnectionsToAllow)
			if err := poller.PollUntilDone(ctx); err != nil {
				return fmt.Errorf("polling deleting %s: %+v", id, err)
			}
		}
	} else {
		if err := resp.Poller.PollUntilDone(ctx); err != nil {
			return fmt.Errorf("deleting %s: %v", *id, err)
		}
	}

	if model := existing.Model; model != nil {
		locationName := location.NormalizeNilable(pointer.To(model.Location))

		// Purge the soft deleted Api Management permanently if the feature flag is enabled
		if meta.(*clients.Client).Features.ApiManagement.PurgeSoftDeleteOnDestroy {
			log.Printf("[DEBUG] %s marked for purge - executing purge", *id)
			deletedServiceId := deletedservice.NewDeletedServiceID(id.SubscriptionId, locationName, id.ServiceName)
			if _, err := deletedServicesClient.GetByName(ctx, deletedServiceId); err != nil {
				return fmt.Errorf("retrieving the deleted %s to be able to purge it: %+v", *id, err)
			}
			resp, err := deletedServicesClient.Purge(ctx, deletedServiceId)
			if err != nil && !response.WasNotFound(resp.HttpResponse) {
				return fmt.Errorf("purging the deleted %s: %+v", *id, err)
			}

			isLroStatus = resp.HttpResponse.StatusCode == http.StatusCreated || resp.HttpResponse.StatusCode == http.StatusAccepted
			if isLroStatus {
				pollerType, err := custompollers.NewAPIManagementPoller(client, resp.HttpResponse)
				if err != nil {
					return fmt.Errorf("polling deleting %s: %+v", id, err)
				}
				if pollerType != nil {
					poller := pollers.NewPoller(pollerType, 20*time.Second, pollers.DefaultNumberOfDroppedConnectionsToAllow)
					if err := poller.PollUntilDone(ctx); err != nil {
						return fmt.Errorf("polling purging the deleting %s: %+v", id, err)
					}
				}
			} else {
				if err := resp.Poller.PollUntilDone(ctx); err != nil {
					return fmt.Errorf("purging the deleted %s: %+v", *id, err)
				}
			}
			log.Printf("[DEBUG] Purged %s.", *id)
			return nil
		}
	}

	return nil
}
