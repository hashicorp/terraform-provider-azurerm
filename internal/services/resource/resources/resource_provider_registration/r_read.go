// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package resource_provider_registration

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/resources/2021-07-01/features"
	"github.com/hashicorp/go-azure-sdk/resource-manager/resources/2022-09-01/providers"
	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/resourceproviders/custompollers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func (r Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Resource.ResourceProvidersClient
			featureClient := metadata.Client.Resource.FeaturesClient
			account := metadata.Client.Account

			id, err := providers.ParseSubscriptionProviderID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if err := r.checkIfManagedByTerraform(id.ProviderName, account); err != nil {
				return err
			}

			resp, err := client.Get(ctx, *id, providers.DefaultGetOperationOptions())
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}

				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			registrationState := ""
			if model := resp.Model; model != nil && model.RegistrationState != nil {
				registrationState = *model.RegistrationState
			}
			if !strings.EqualFold(registrationState, Registered) {
				// Account for API inconsistency before removing from state
				pollerType := custompollers.NewResourceProviderExistencePoller(client, *id, 10)
				poller := pollers.NewPoller(pollerType, 10*time.Second, pollers.DefaultNumberOfDroppedConnectionsToAllow)
				if err := poller.PollUntilDone(ctx); err != nil {
					if errors.As(err, &pollers.PollingCancelledError{}) {
						log.Printf("[WARN] %s was not registered - removing from state", id)
						return metadata.MarkAsGone(id)
					}
					return fmt.Errorf("polling %s: %+v", id, err)
				}
			}

			resourceProviderFeatureId := features.NewProviders2ID(id.SubscriptionId, id.ProviderName)
			result, err := featureClient.ListComplete(ctx, resourceProviderFeatureId)
			if err != nil {
				return fmt.Errorf("retrieving features for %s: %+v", *id, err)
			}
			features := make([]ResourceProviderRegistrationFeatureModel, 0)
			for _, item := range result.Items {
				if item.Properties != nil && item.Properties.State != nil && item.Name != nil {
					featureName := (*item.Name)[len(id.ProviderName)+1:]
					switch *item.Properties.State {
					case Registering, Registered:
						features = append(features, ResourceProviderRegistrationFeatureModel{Name: featureName, Registered: true})
					case Unregistering, Unregistered:
						features = append(features, ResourceProviderRegistrationFeatureModel{Name: featureName})
					}
				}
			}

			return metadata.Encode(&ResourceProviderRegistrationModel{
				Name:     id.ProviderName,
				Features: features,
			})
		},
		Timeout: 15 * time.Minute,
	}
}

const (
	Pending       = "Pending"
	Registering   = "Registering"
	Unregistering = "Unregistering"
	Registered    = "Registered"
	NotRegistered = "NotRegistered"
	Unregistered  = "Unregistered"
)

func (r Resource) checkIfManagedByTerraform(name string, account *clients.ResourceManagerAccount) error {
	for resourceProvider := range account.RegisteredResourceProviders {
		if resourceProvider == name {
			fmtStr := `
The Resource Provider %[1]q is automatically registered by Terraform.

To manage this Resource Provider registration with the "azurerm_resource_provider_registration" resource, you need to
prevent Terraform from managing this Resource Provider automatically by one of these methods:

1. Disable automatic Resource Provider registration by setting the following in the Provider block:

   resource_provider_registrations = "none"

2. Choose a set of Provider Registrations to automatically register, that do not include %[1]q. Refer to the
   provider documentation for more information:

   https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs#resource_provider_registrations`
			return fmt.Errorf(fmtStr, name)
		}
	}

	return nil
}

func (r Resource) applyFeatures(ctx context.Context, metadata sdk.ResourceMetaData, id providers.SubscriptionProviderId, oldFeatures []any, newFeatures []any) error {
	for _, v := range newFeatures {
		value := v.(map[string]any)
		name := value["name"].(string)
		featureId := features.NewFeatureID(id.SubscriptionId, id.ProviderName, name)
		if value["registered"].(bool) {
			if err := r.registerFeature(ctx, metadata, featureId); err != nil {
				return err
			}
		} else {
			if err := r.unregisterFeature(ctx, metadata, featureId); err != nil {
				return err
			}
		}
	}

	// unregister the features which block is removed now
	unmanagedRegisteredFeatures := make(map[string]bool)
	for _, v := range oldFeatures {
		value := v.(map[string]any)
		name := value["name"].(string)
		unmanagedRegisteredFeatures[name] = value["registered"].(bool)
	}
	for _, v := range newFeatures {
		value := v.(map[string]any)
		name := value["name"].(string)
		unmanagedRegisteredFeatures[name] = false
	}

	for featureName, registered := range unmanagedRegisteredFeatures {
		if registered {
			featureId := features.NewFeatureID(id.SubscriptionId, id.ProviderName, featureName)
			if err := r.unregisterFeature(ctx, metadata, featureId); err != nil {
				return err
			}
		}
	}
	return nil
}

func (r Resource) registerFeature(ctx context.Context, metadata sdk.ResourceMetaData, id features.FeatureId) error {
	client := metadata.Client.Resource.FeaturesClient
	existing, err := client.Get(ctx, id)
	if err != nil {
		return fmt.Errorf("error checking for existing feature %q: %+v", id, err)
	}

	if existing.Model != nil && existing.Model.Properties != nil && existing.Model.Properties.State != nil {
		if strings.EqualFold(*existing.Model.Properties.State, Pending) {
			return fmt.Errorf("%s which requires manual approval should not be managed by terraform", id)
		}
		if strings.EqualFold(*existing.Model.Properties.State, Registered) {
			return nil
		}
	}

	log.Printf("[INFO] registering feature %q.", id)
	resp, err := client.Register(ctx, id)
	if err != nil {
		return fmt.Errorf("error registering feature %q: %+v", id, err)
	}

	if resp.Model != nil && resp.Model.Properties != nil && resp.Model.Properties.State != nil {
		if strings.EqualFold(*resp.Model.Properties.State, Pending) {
			return fmt.Errorf("%s which requires manual approval can not be managed by terraform", id)
		}
	}

	deadline, ok := ctx.Deadline()
	if !ok {
		return fmt.Errorf("internal-error: context had no deadline")
	}
	stateConf := &pluginsdk.StateChangeConf{
		Pending:    []string{Registering},
		Target:     []string{Registered},
		Refresh:    r.featureRegisteringStateRefreshFunc(ctx, client, id),
		MinTimeout: 3 * time.Minute,
		Timeout:    time.Until(deadline),
	}

	if _, err = stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("waiting for %s registering to be completed: %+v", id, err)
	}
	return nil
}

func (r Resource) unregisterFeature(ctx context.Context, metadata sdk.ResourceMetaData, id features.FeatureId) error {
	client := metadata.Client.Resource.FeaturesClient
	existing, err := client.Get(ctx, id)
	if err != nil {
		return fmt.Errorf("error checking for existing feature %q: %+v", id, err)
	}

	if existing.Model != nil && existing.Model.Properties != nil && existing.Model.Properties.State != nil {
		if strings.EqualFold(*existing.Model.Properties.State, Pending) {
			return fmt.Errorf("%s which requires manual approval should not be managed by terraform", id)
		}
		if strings.EqualFold(*existing.Model.Properties.State, Unregistered) {
			return nil
		}
	}

	log.Printf("[INFO] unregistering feature %q.", id)
	resp, err := client.Unregister(ctx, id)
	if err != nil {
		return fmt.Errorf("unregistering feature %q: %+v", id, err)
	}

	if resp.Model != nil && resp.Model.Properties != nil && resp.Model.Properties.State != nil {
		if strings.EqualFold(*resp.Model.Properties.State, Pending) {
			return fmt.Errorf("%s requires manual registration approval and can not be managed by terraform", id)
		}
	}

	deadline, ok := ctx.Deadline()
	if !ok {
		return fmt.Errorf("internal-error: context had no deadline")
	}
	stateConf := &pluginsdk.StateChangeConf{
		Pending:    []string{Unregistering},
		Target:     []string{NotRegistered, Unregistered},
		Refresh:    r.featureRegisteringStateRefreshFunc(ctx, client, id),
		MinTimeout: 3 * time.Minute,
		Timeout:    time.Until(deadline),
	}

	if _, err = stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("waiting for %s to be complete unregistering: %+v", id, err)
	}

	return nil
}

func (r Resource) featureRegisteringStateRefreshFunc(ctx context.Context, client *features.FeaturesClient, id features.FeatureId) pluginsdk.StateRefreshFunc {
	return func() (any, string, error) {
		res, err := client.Get(ctx, id)
		if err != nil {
			return nil, "", fmt.Errorf("retrieving %s: %+v", id, err)
		}
		if res.Model == nil || res.Model.Properties == nil || res.Model.Properties.State == nil {
			return nil, "", fmt.Errorf("error reading %s registering status: %+v", id, err)
		}

		return res, *res.Model.Properties.State, nil
	}
}
