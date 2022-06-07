package resource

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/profiles/2017-03-09/resources/mgmt/resources"
	"github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2015-12-01/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/resourceproviders"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/resource/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/resource/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
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
			client := metadata.Client.Resource.ProvidersClient
			account := metadata.Client.Account

			var obj ResourceProviderRegistrationModel
			if err := metadata.Decode(&obj); err != nil {
				return err
			}

			resourceId := parse.NewResourceProviderID(account.SubscriptionId, obj.Name)
			if err := r.checkIfManagedByTerraform(resourceId.ResourceProvider, account); err != nil {
				return err
			}

			provider, err := client.Get(ctx, resourceId.ResourceProvider, "")
			if err != nil {
				if utils.ResponseWasNotFound(provider.Response) {
					return fmt.Errorf("the Resource Provider %q was not found", resourceId.ResourceProvider)
				}

				return fmt.Errorf("retrieving Resource Provider %q: %+v", resourceId.ResourceProvider, err)
			}
			if provider.RegistrationState == nil {
				return fmt.Errorf("retrieving Resource Provider %q: `registrationState` was nil", resourceId.ResourceProvider)
			}

			if strings.EqualFold(*provider.RegistrationState, "Registered") {
				return metadata.ResourceRequiresImport(r.ResourceType(), resourceId)
			}

			if metadata.ResourceData.HasChange("feature") {
				oldFeaturesRaw, newFeaturesRaw := metadata.ResourceData.GetChange("feature")
				err := r.applyFeatures(ctx, metadata, resourceId, oldFeaturesRaw.(*pluginsdk.Set).List(), newFeaturesRaw.(*pluginsdk.Set).List())
				if err != nil {
					return fmt.Errorf("applying features for Resource Provider %q: %+v", resourceId.ResourceProvider, err)
				}
			}

			log.Printf("[DEBUG] Registering Resource Provider %q..", resourceId.ResourceProvider)
			if _, err := client.Register(ctx, resourceId.ResourceProvider); err != nil {
				return fmt.Errorf("registering Resource Provider %q: %+v", resourceId.ResourceProvider, err)
			}

			deadline, ok := ctx.Deadline()
			if !ok {
				return fmt.Errorf("could not retrieve context deadline for %s", resourceId.ID())
			}
			// TODO: @tombuildsstuff - expose a nicer means of doing this in the SDK
			log.Printf("[DEBUG] Waiting for Resource Provider %q to finish registering..", resourceId.ResourceProvider)
			stateConf := &pluginsdk.StateChangeConf{
				Pending:      []string{"Processing"},
				Target:       []string{"Registered"},
				Refresh:      r.registerRefreshFunc(ctx, client, resourceId.ResourceProvider),
				MinTimeout:   15 * time.Second,
				PollInterval: 30 * time.Second,
				Timeout:      time.Until(deadline),
			}
			if _, err := stateConf.WaitForStateContext(ctx); err != nil {
				return fmt.Errorf("waiting for Resource Provider Namespace %q to be registered: %s", resourceId.ResourceProvider, err)
			}
			log.Printf("[DEBUG] Registered Resource Provider %q.", resourceId.ResourceProvider)

			metadata.SetID(resourceId)
			return nil
		},

		Timeout: 120 * time.Minute,
	}
}

func (r ResourceProviderRegistrationResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Resource.ProvidersClient
			account := metadata.Client.Account

			var obj ResourceProviderRegistrationModel
			if err := metadata.Decode(&obj); err != nil {
				return err
			}

			resourceId := parse.NewResourceProviderID(account.SubscriptionId, obj.Name)
			if err := r.checkIfManagedByTerraform(resourceId.ResourceProvider, account); err != nil {
				return err
			}

			provider, err := client.Get(ctx, resourceId.ResourceProvider, "")
			if err != nil {
				if utils.ResponseWasNotFound(provider.Response) {
					return fmt.Errorf("the Resource Provider %q was not found", resourceId.ResourceProvider)
				}

				return fmt.Errorf("retrieving Resource Provider %q: %+v", resourceId.ResourceProvider, err)
			}
			if provider.RegistrationState == nil {
				return fmt.Errorf("retrieving Resource Provider %q: `registrationState` was nil", resourceId.ResourceProvider)
			}

			if !strings.EqualFold(*provider.RegistrationState, "Registered") {
				return fmt.Errorf("retrieving Resource Provider %q: `registrationState` was not `Registered`", resourceId.ResourceProvider)
			}

			if metadata.ResourceData.HasChange("feature") {
				oldFeaturesRaw, newFeaturesRaw := metadata.ResourceData.GetChange("feature")
				err := r.applyFeatures(ctx, metadata, resourceId, oldFeaturesRaw.(*pluginsdk.Set).List(), newFeaturesRaw.(*pluginsdk.Set).List())
				if err != nil {
					return fmt.Errorf("applying features for Resource Provider %q: %+v", resourceId.ResourceProvider, err)
				}
			}

			log.Printf("[DEBUG] Registering Resource Provider %q..", resourceId.ResourceProvider)
			if _, err := client.Register(ctx, resourceId.ResourceProvider); err != nil {
				return fmt.Errorf("registering Resource Provider %q: %+v", resourceId.ResourceProvider, err)
			}

			deadline, ok := ctx.Deadline()
			if !ok {
				return fmt.Errorf("could not retrieve context deadline for %s", resourceId.ID())
			}
			// TODO: @tombuildsstuff - expose a nicer means of doing this in the SDK
			log.Printf("[DEBUG] Waiting for Resource Provider %q to finish registering..", resourceId.ResourceProvider)
			stateConf := &pluginsdk.StateChangeConf{
				Pending:      []string{"Processing"},
				Target:       []string{"Registered"},
				Refresh:      r.registerRefreshFunc(ctx, client, resourceId.ResourceProvider),
				MinTimeout:   15 * time.Second,
				PollInterval: 30 * time.Second,
				Timeout:      time.Until(deadline),
			}
			if _, err := stateConf.WaitForStateContext(ctx); err != nil {
				return fmt.Errorf("waiting for Resource Provider Namespace %q to be registered: %s", resourceId.ResourceProvider, err)
			}
			log.Printf("[DEBUG] Registered Resource Provider %q.", resourceId.ResourceProvider)

			metadata.SetID(resourceId)
			return nil
		},
		Timeout: 120 * time.Minute,
	}
}

func (r ResourceProviderRegistrationResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Resource.ProvidersClient
			featureClient := metadata.Client.Resource.FeaturesClient
			account := metadata.Client.Account

			id, err := parse.ResourceProviderID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if err := r.checkIfManagedByTerraform(id.ResourceProvider, account); err != nil {
				return err
			}

			resp, err := client.Get(ctx, id.ResourceProvider, "")
			if err != nil {
				if utils.ResponseWasNotFound(resp.Response) {
					return metadata.MarkAsGone(id)
				}

				return fmt.Errorf("retrieving Resource Provider %q: %+v", id.ResourceProvider, err)
			}

			if resp.RegistrationState != nil && !strings.EqualFold(*resp.RegistrationState, "Registered") {
				log.Printf("[WARN] Resource Provider %q was not registered", id.ResourceProvider)
				return metadata.MarkAsGone(id)
			}

			result, err := featureClient.ListComplete(ctx, id.ResourceProvider)
			if err != nil {
				return fmt.Errorf("retrieving features for Resource Provider %q: %+v", id.ResourceProvider, err)
			}
			features := make([]ResourceProviderRegistrationFeatureModel, 0)
			for result.NotDone() {
				value := result.Value()
				if value.Properties != nil && value.Properties.State != nil && value.Name != nil {
					featureName := (*value.Name)[len(id.ResourceProvider)+1:]
					switch *value.Properties.State {
					case Registering, Registered:
						features = append(features, ResourceProviderRegistrationFeatureModel{Name: featureName, Registered: true})
					case Unregistering, Unregistered:
						features = append(features, ResourceProviderRegistrationFeatureModel{Name: featureName, Registered: false})
					}
				}
				if err := result.NextWithContext(ctx); err != nil {
					return fmt.Errorf("enumerating features: %+v", err)
				}
			}

			return metadata.Encode(&ResourceProviderRegistrationModel{
				Name:     id.ResourceProvider,
				Features: features,
			})
		},
		Timeout: 5 * time.Minute,
	}
}

func (r ResourceProviderRegistrationResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Resource.ProvidersClient
			account := metadata.Client.Account

			id, err := parse.ResourceProviderID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if err := r.checkIfManagedByTerraform(id.ResourceProvider, account); err != nil {
				return err
			}

			err = r.applyFeatures(ctx, metadata, *id, metadata.ResourceData.Get("feature").(*pluginsdk.Set).List(), make([]interface{}, 0))
			if err != nil {
				return fmt.Errorf("applying features for Resource Provider %q: %+v", id.ResourceProvider, err)
			}

			if _, err := client.Unregister(ctx, id.ResourceProvider); err != nil {
				return fmt.Errorf("unregistering Resource Provider %q: %+v", id.ResourceProvider, err)
			}

			deadline, ok := ctx.Deadline()
			if !ok {
				return fmt.Errorf("could not retrieve context deadline for %s", id.ID())
			}
			// TODO: @tombuildsstuff - we should likely expose something in the SDK to make this easier
			stateConf := &pluginsdk.StateChangeConf{
				Pending:    []string{"Processing"},
				Target:     []string{"Unregistered"},
				Refresh:    r.unregisterRefreshFunc(ctx, client, id.ResourceProvider),
				MinTimeout: 15 * time.Second,
				Timeout:    time.Until(deadline),
			}
			if _, err := stateConf.WaitForStateContext(ctx); err != nil {
				return fmt.Errorf("waiting for Resource Provider %q to become unregistered: %+v", id.ResourceProvider, err)
			}

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
		client := metadata.Client.Resource.ProvidersClient
		account := metadata.Client.Account

		id, err := parse.ResourceProviderID(metadata.ResourceData.Id())
		if err != nil {
			return err
		}

		provider, err := client.Get(ctx, id.ResourceProvider, "")
		if err != nil {
			return fmt.Errorf("retrieving Resource Provider %q: %+v", id.ResourceProvider, err)
		}

		if provider.Namespace == nil {
			return fmt.Errorf("retrieving Resource Provider %q: `namespace` was nil", id.ResourceProvider)
		}

		if *provider.Namespace != id.ResourceProvider {
			return fmt.Errorf("importing Resource Provider %q: expected %q", id.ResourceProvider, *provider.Namespace)
		}

		if provider.RegistrationState == nil || !strings.EqualFold(*provider.RegistrationState, "Registered") {
			return fmt.Errorf("importing Resource Provider %q: Resource Provider must be registered to be imported", id.ResourceProvider)
		}

		if err := r.checkIfManagedByTerraform(id.ResourceProvider, account); err != nil {
			return fmt.Errorf("importing Resource Provider %q: %+v", id.ResourceProvider, err)
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

func (r ResourceProviderRegistrationResource) registerRefreshFunc(ctx context.Context, client *resources.ProvidersClient, resourceProviderNamespace string) pluginsdk.StateRefreshFunc {
	return func() (interface{}, string, error) {
		resp, err := client.Get(ctx, resourceProviderNamespace, "")
		if err != nil {
			return resp, "Failed", err
		}

		if resp.RegistrationState != nil && strings.EqualFold(*resp.RegistrationState, "Registered") {
			return resp, "Registered", nil
		}

		return resp, "Processing", nil
	}
}

func (r ResourceProviderRegistrationResource) unregisterRefreshFunc(ctx context.Context, client *resources.ProvidersClient, resourceProvider string) pluginsdk.StateRefreshFunc {
	return func() (interface{}, string, error) {
		resp, err := client.Get(ctx, resourceProvider, "")
		if err != nil {
			return resp, "Failed", err
		}

		if resp.RegistrationState != nil && strings.EqualFold(*resp.RegistrationState, "Unregistered") {
			return resp, "Unregistered", nil
		}

		return resp, "Processing", nil
	}
}

func (r ResourceProviderRegistrationResource) applyFeatures(ctx context.Context, metadata sdk.ResourceMetaData, id parse.ResourceProviderId, oldFeatures []interface{}, newFeatures []interface{}) error {
	for _, v := range newFeatures {
		value := v.(map[string]interface{})
		name := value["name"].(string)
		if value["registered"].(bool) {
			if err := r.registerFeature(ctx, metadata, parse.NewFeatureID(id.SubscriptionId, id.ResourceProvider, name)); err != nil {
				return err
			}
		} else {
			if err := r.unregisterFeature(ctx, metadata, parse.NewFeatureID(id.SubscriptionId, id.ResourceProvider, name)); err != nil {
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
			if err := r.unregisterFeature(ctx, metadata, parse.NewFeatureID(id.SubscriptionId, id.ResourceProvider, featureName)); err != nil {
				return err
			}
		}
	}
	return nil
}

func (r ResourceProviderRegistrationResource) registerFeature(ctx context.Context, metadata sdk.ResourceMetaData, id parse.FeatureId) error {
	client := metadata.Client.Resource.FeaturesClient
	existing, err := client.Get(ctx, id.ProviderNamespace, id.Name)
	if err != nil {
		return fmt.Errorf("error checking for existing feature %q: %+v", id, err)
	}

	if existing.Properties != nil && existing.Properties.State != nil {
		if strings.EqualFold(*existing.Properties.State, Pending) {
			return fmt.Errorf("%s which requires manual approval should not be managed by terraform", id)
		}
		if strings.EqualFold(*existing.Properties.State, Registered) {
			return nil
		}
	}

	log.Printf("[INFO] registering feature %q.", id)
	resp, err := client.Register(ctx, id.ProviderNamespace, id.Name)
	if err != nil {
		return fmt.Errorf("error registering feature %q: %+v", id, err)
	}

	if resp.Properties != nil && resp.Properties.State != nil {
		if strings.EqualFold(*resp.Properties.State, Pending) {
			return fmt.Errorf("%s which requires manual approval can not be managed by terraform", id)
		}
	}

	deadline, ok := ctx.Deadline()
	if !ok {
		return fmt.Errorf("context had no deadline")
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

func (r ResourceProviderRegistrationResource) unregisterFeature(ctx context.Context, metadata sdk.ResourceMetaData, id parse.FeatureId) error {
	client := metadata.Client.Resource.FeaturesClient
	existing, err := client.Get(ctx, id.ProviderNamespace, id.Name)
	if err != nil {
		return fmt.Errorf("error checking for existing feature %q: %+v", id, err)
	}

	if existing.Properties != nil && existing.Properties.State != nil {
		if strings.EqualFold(*existing.Properties.State, Pending) {
			return fmt.Errorf("%s which requires manual approval should not be managed by terraform", id)
		}
		if strings.EqualFold(*existing.Properties.State, Unregistered) {
			return nil
		}
	}

	log.Printf("[INFO] unregistering feature %q.", id)
	resp, err := client.Unregister(ctx, id.ProviderNamespace, id.Name)
	if err != nil {
		return fmt.Errorf("unregistering feature %q: %+v", id, err)
	}

	if resp.Properties != nil && resp.Properties.State != nil {
		if strings.EqualFold(*resp.Properties.State, Pending) {
			return fmt.Errorf("%s requires manual registration approval and can not be managed by terraform", id)
		}
	}

	deadline, ok := ctx.Deadline()
	if !ok {
		return fmt.Errorf("context had no deadline")
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

func (r ResourceProviderRegistrationResource) featureRegisteringStateRefreshFunc(ctx context.Context, client *features.Client, id parse.FeatureId) pluginsdk.StateRefreshFunc {
	return func() (interface{}, string, error) {
		res, err := client.Get(ctx, id.ProviderNamespace, id.Name)
		if err != nil {
			return nil, "", fmt.Errorf("retrieving %s: %+v", id, err)
		}
		if res.Properties == nil || res.Properties.State == nil {
			return nil, "", fmt.Errorf("error reading %s registering status: %+v", id, err)
		}

		return res, *res.Properties.State, nil
	}
}
