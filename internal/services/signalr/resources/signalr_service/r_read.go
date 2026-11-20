// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package signalr_service

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/signalr/2024-03-01/signalr"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceArmSignalRServiceRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).SignalR.SignalRClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := signalr.ParseSignalRID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[DEBUG] %s was not found - removing from state!", *id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	keys, err := client.ListKeys(ctx, *id)
	if err != nil {
		return fmt.Errorf("listing keys for %s: %+v", *id, err)
	}

	return resourceArmSignalRServiceFlatten(d, id, resp.Model, keys.Model)
}

func signalRIsInServerlessMode(features *[]signalr.SignalRFeature) bool {
	if features == nil {
		return false
	}

	for _, feature := range *features {
		if feature.Flag == signalr.FeatureFlagsServiceMode {
			return strings.EqualFold(feature.Value, "Serverless")
		}
	}

	return false
}

func signalRFeature(featureFlag signalr.FeatureFlags, value string) signalr.SignalRFeature {
	return signalr.SignalRFeature{
		Flag:  featureFlag,
		Value: value,
	}
}

func signalrServiceProvisioningStateRefreshFunc(ctx context.Context, client *signalr.SignalRClient, id signalr.SignalRId) pluginsdk.StateRefreshFunc {
	return func() (any, string, error) {
		res, err := client.Get(ctx, id)
		if err != nil {
			return nil, "", fmt.Errorf("polling for the provisioning state of %s: %+v", id, err)
		}

		if model := res.Model; model != nil {
			if model.Properties != nil && model.Properties.ProvisioningState != nil {
				return res, string(*model.Properties.ProvisioningState), nil
			}
		}

		return nil, "", fmt.Errorf("unable to read provisioning state")
	}
}
