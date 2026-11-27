// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package servicebus_namespace

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/servicebus/2026-01-01/namespaces"
	servicebusClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/servicebus/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func waitForNamespaceStatusToBeReady(ctx context.Context, meta any, id namespaces.NamespaceId, timeout time.Duration) error {
	namespaceClient := servicebusClient.From(meta).NamespacesClient
	stateConf := &pluginsdk.StateChangeConf{
		Pending: []string{
			string(namespaces.EndPointProvisioningStateUpdating),
			string(namespaces.EndPointProvisioningStateCreating),
			string(namespaces.EndPointProvisioningStateDeleting),
		},
		Target:                    []string{string(namespaces.EndPointProvisioningStateSucceeded)},
		Refresh:                   serviceBusNamespaceProvisioningStateRefreshFunc(ctx, namespaceClient, id),
		Timeout:                   timeout,
		PollInterval:              10 * time.Second,
		ContinuousTargetOccurence: 5,
	}

	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return err
	}
	return nil
}

func serviceBusNamespaceProvisioningStateRefreshFunc(ctx context.Context, client *namespaces.NamespacesClient, id namespaces.NamespaceId) pluginsdk.StateRefreshFunc {
	return func() (any, string, error) {
		res, err := client.Get(ctx, id)
		if err != nil {
			return nil, "", fmt.Errorf("retrieving servicebus namespace error: %+v", err)
		}
		if res.Model == nil || res.Model.Properties == nil || res.Model.Properties.ProvisioningState == nil {
			return nil, "", fmt.Errorf("retrieving %s: model/provisioningState was nil", id)
		}

		return res, *res.Model.Properties.ProvisioningState, nil
	}
}
