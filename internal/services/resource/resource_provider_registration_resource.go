// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package resource

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/resources/2021-07-01/features"
	"github.com/hashicorp/go-azure-sdk/resource-manager/resources/2022-09-01/providers"
	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/resourceproviders"
	"github.com/hashicorp/terraform-provider-azurerm/internal/resourceproviders/custompollers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/resource/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

var (
	_ sdk.Resource                   = ResourceProviderRegistrationResource{}
	_ sdk.ResourceWithCustomImporter = ResourceProviderRegistrationResource{}
)

type ResourceProviderRegistrationResource struct{}

type ResourceProviderRegistrationModel struct {
	Name     string                                     `tfschema:"name"`
	Features []ResourceProviderRegistrationFeatureModel `tfschema:"feature"`
}

type ResourceProviderRegistrationFeatureModel struct {
	Name       string `tfschema:"name"`
	Registered bool   `tfschema:"registered"`
}

const (
	Pending       = "Pending"
	Registering   = "Registering"
	Unregistering = "Unregistering"
	Registered    = "Registered"
	NotRegistered = "NotRegistered"
	Unregistered  = "Unregistered"
)

func (r ResourceProviderRegistrationResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: resourceproviders.EnhancedValidate,
		},

		"feature": {
			Type:     pluginsdk.TypeSet,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"registered": {
						Type:     pluginsdk.TypeBool,
						Required: true,
					},
				},
			},
		},
	}
}

func (r ResourceProviderRegistrationResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r ResourceProviderRegistrationResource) ModelObject() interface{} {
	return &ResourceProviderRegistrationModel{}
}

func (r ResourceProviderRegistrationResource) ResourceType() string {
	return "azurerm_resource_provider_registration"
}

func (r ResourceProviderRegistrationResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Resource.ResourceProvidersClient
			account := metadata.Client.Account

			var obj ResourceProviderRegistrationModel
			if err := metadata.Decode(&obj); err != nil {
				return err
			}

			resourceId := providers.NewSubscriptionProviderID(account.SubscriptionId, obj.Name)
			if err := r.checkIfManagedByTerraform(resourceId.ProviderName, account); err != nil {
				return err
			}

			provider, err := client.Get(ctx, resourceId, providers.DefaultGetOperationOptions())
			if err != nil {
				if response.WasNotFound(provider.HttpResponse) {
					return fmt.Errorf("%s was not found", resourceId)
				}

				return fmt.Errorf("retrieving %q: %+v", resourceId, err)
			}
			registrationState := ""
			if model := provider.Model; model != nil && model.RegistrationState != nil {
				registrationState = *model.RegistrationState
			}
			if registrationState == "" {
				return fmt.Errorf("retrieving %s: `registrationState` was nil", resourceId)
			}
			if strings.EqualFold(registrationState, "Registered") {
				return metadata.ResourceRequiresImport(r.ResourceType(), resourceId)
			}

			if metadata.ResourceData.HasChange("feature") {
				oldFeaturesRaw, newFeaturesRaw := metadata.ResourceData.GetChange("feature")
				err := r.applyFeatures(ctx, metadata, resourceId, oldFeaturesRaw.(*pluginsdk.Set).List(), newFeaturesRaw.(*pluginsdk.Set).List())
				if err != nil {
					return fmt.Errorf("applying features for %q: %+v", resourceId, err)
				}
			}

			log.Printf("[DEBUG] Registering %s..", resourceId)
			payload := providers.ProviderRegistrationRequest{}
			if _, err := client.Register(ctx, resourceId, payload); err != nil {
				return fmt.Errorf("registering %s: %+v", resourceId, err)
			}

			log.Printf("[DEBUG] Waiting for %s to finish registering..", resourceId)
			pollerType := custompollers.NewResourceProviderRegistrationPoller(client, resourceId)
			poller := pollers.NewPoller(pollerType, 10*time.Second, pollers.DefaultNumberOfDroppedConnectionsToAllow)
			if err := poller.PollUntilDone(ctx); err != nil {
				return fmt.Errorf("waiting for %s to be registered: %s", resourceId, err)
			}
			log.Printf("[DEBUG] Registered Resource Provider %q.", resourceId)

			metadata.SetID(resourceId)
			return nil
		},

		Timeout: 120 * time.Minute,
	}
}

func (r ResourceProviderRegistrationResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Resource.ResourceProvidersClient
			account := metadata.Client.Account

			var obj ResourceProviderRegistrationModel
			if err := metadata.Decode(&obj); err != nil {
				return err
			}

			resourceId, err := providers.ParseSubscriptionProviderID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}
			if err := r.checkIfManagedByTerraform(resourceId.ProviderName, account); err != nil {
				return err
			}

			provider, err := client.Get(ctx, *resourceId, providers.DefaultGetOperationOptions())
			if err != nil {
				if response.WasNotFound(provider.HttpResponse) {
					return fmt.Errorf("the %s was not found", *resourceId)
				}

				return fmt.Errorf("retrieving %s: %+v", *resourceId, err)
			}
			registrationState := ""
			if model := provider.Model; model != nil && model.RegistrationState != nil {
				registrationState = *model.RegistrationState
			}
			if registrationState == "" {
				return fmt.Errorf("retrieving %s: `registrationState` was nil", *resourceId)
			}

			if !strings.EqualFold(registrationState, "Registered") {
				return fmt.Errorf("retrieving %s: `registrationState` was not `Registered` but %q", *resourceId, registrationState)
			}

			if metadata.ResourceData.HasChange("feature") {
				oldFeaturesRaw, newFeaturesRaw := metadata.ResourceData.GetChange("feature")
				err := r.applyFeatures(ctx, metadata, *resourceId, oldFeaturesRaw.(*pluginsdk.Set).List(), newFeaturesRaw.(*pluginsdk.Set).List())
				if err != nil {
					return fmt.Errorf("applying features for %s: %+v", *resourceId, err)
				}
			}

			log.Printf("[DEBUG] Registering %s..", *resourceId)
			payload := providers.ProviderRegistrationRequest{}
			if _, err := client.Register(ctx, *resourceId, payload); err != nil {
				return fmt.Errorf("registering %s: %+v", *resourceId, err)
			}

			log.Printf("[DEBUG] Waiting for %s to finish registering..", resourceId)
			pollerType := custompollers.NewResourceProviderRegistrationPoller(client, *resourceId)
			poller := pollers.NewPoller(pollerType, 10*time.Second, pollers.DefaultNumberOfDroppedConnectionsToAllow)
			if err := poller.PollUntilDone(ctx); err != nil {
				return fmt.Errorf("waiting for %s to be registered: %s", resourceId, err)
			}
			log.Printf("[DEBUG] Registered Resource Provider %q.", resourceId)

			metadata.SetID(resourceId)
			return nil
		},
		Timeout: 120 * time.Minute,
	}
}

func (r ResourceProviderRegistrationResource) Read() sdk.ResourceFunc {
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
			if !strings.EqualFold(registrationState, "Registered") {
				log.Printf("[WARN] %s was not registered - removing from state", id)
				return metadata.MarkAsGone(id)
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
						features = append(features, ResourceProviderRegistrationFeatureModel{Name: featureName, Registered: false})
					}
				}
			}

			return metadata.Encode(&ResourceProviderRegistrationModel{
				Name:     id.ProviderName,
				Features: features,
			})
		},
		Timeout: 5 * time.Minute,
	}
}

func (r ResourceProviderRegistrationResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Resource.ResourceProvidersClient
			account := metadata.Client.Account

			id, err := providers.ParseSubscriptionProviderID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if err := r.checkIfManagedByTerraform(id.ProviderName, account); err != nil {
				return err
			}

			err = r.applyFeatures(ctx, metadata, *id, metadata.ResourceData.Get("feature").(*pluginsdk.Set).List(), make([]interface{}, 0))
			if err != nil {
				return fmt.Errorf("applying features for %s: %+v", *id, err)
			}

			if _, err := client.Unregister(ctx, *id); err != nil {
				return fmt.Errorf("unregistering Resource Provider %q: %+v", *id, err)
			}

			log.Printf("[DEBUG] Waiting for %s to finish unregistering..", *id)
			pollerType := custompollers.NewResourceProviderUnregistrationPoller(client, *id)
			poller := pollers.NewPoller(pollerType, 10*time.Second, pollers.DefaultNumberOfDroppedConnectionsToAllow)
			if err := poller.PollUntilDone(ctx); err != nil {
				return fmt.Errorf("waiting for %s to become unregistered: %+v", *id, err)
			}
			log.Printf("[DEBUG] Unregistered Resource Provider %q.", *id)

			return nil
		},
		Timeout: 30 * time.Minute,
	}
}

func (r ResourceProviderRegistrationResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return validate.ResourceProviderID
}

func (r ResourceProviderRegistrationResource) CustomImporter() sdk.ResourceRunFunc {
	return func(ctx context.Context, metadata sdk.ResourceMetaData) error {
		client := metadata.Client.Resource.ResourceProvidersClient
		account := metadata.Client.Account

		id, err := providers.ParseSubscriptionProviderID(metadata.ResourceData.Id())
		if err != nil {
			return err
		}

		provider, err := client.Get(ctx, *id, providers.DefaultGetOperationOptions())
		if err != nil {
			return fmt.Errorf("retrieving %s: %+v", *id, err)
		}

		namespace := ""
		registrationState := ""
		if model := provider.Model; model != nil {
			if model.Namespace != nil {
				namespace = *model.Namespace
			}
			if model.RegistrationState != nil {
				registrationState = *model.RegistrationState
			}
		}
		if namespace != id.ProviderName {
			return fmt.Errorf("importing %s: expected %q but got %q", *id, id.ProviderName, namespace)
		}

		if !strings.EqualFold(registrationState, "Registered") {
			return fmt.Errorf("importing %s: Resource Provider must be registered to be imported", id.ProviderName)
		}

		if err := r.checkIfManagedByTerraform(id.ProviderName, account); err != nil {
			return fmt.Errorf("importing %s: %+v", *id, err)
		}

		return nil
	}
}

func (r ResourceProviderRegistrationResource) checkIfManagedByTerraform(name string, account *clients.ResourceManagerAccount) error {
	if account.SkipResourceProviderRegistration {
		return nil
	}

	for resourceProvider := range resourceproviders.Required() {
		if resourceProvider == name {
			fmtStr := `The Resource Provider %q is automatically registered by Terraform.

To manage this Resource Provider Registration with Terraform you need to opt-out
of Automatic Resource Provider Registration (by setting 'skip_provider_registration'
to 'true' in the Provider block) to avoid conflicting with Terraform.`
			return fmt.Errorf(fmtStr, name)
		}
	}

	return nil
}

func (r ResourceProviderRegistrationResource) applyFeatures(ctx context.Context, metadata sdk.ResourceMetaData, id providers.SubscriptionProviderId, oldFeatures []interface{}, newFeatures []interface{}) error {
	for _, v := range newFeatures {
		value := v.(map[string]interface{})
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
		value := v.(map[string]interface{})
		name := value["name"].(string)
		unmanagedRegisteredFeatures[name] = value["registered"].(bool)
	}
	for _, v := range newFeatures {
		value := v.(map[string]interface{})
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

func (r ResourceProviderRegistrationResource) registerFeature(ctx context.Context, metadata sdk.ResourceMetaData, id features.FeatureId) error {
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

func (r ResourceProviderRegistrationResource) unregisterFeature(ctx context.Context, metadata sdk.ResourceMetaData, id features.FeatureId) error {
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

func (r ResourceProviderRegistrationResource) featureRegisteringStateRefreshFunc(ctx context.Context, client *features.FeaturesClient, id features.FeatureId) pluginsdk.StateRefreshFunc {
	return func() (interface{}, string, error) {
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
