package resource

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/profiles/2017-03-09/resources/mgmt/resources"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/resourceproviders"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/sdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/resource/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/resource/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

var _ sdk.Resource = ResourceProviderRegistrationResource{}
var _ sdk.ResourceWithCustomImporter = ResourceProviderRegistrationResource{}

type ResourceProviderRegistrationResource struct{}

type ResourceProviderRegistrationModel struct {
	Name string `tfschema:"name"`
}

func (r ResourceProviderRegistrationResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: resourceproviders.EnhancedValidate,
		},
	}
}

func (r ResourceProviderRegistrationResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r ResourceProviderRegistrationResource) ModelObject() interface{} {
	return ResourceProviderRegistrationModel{}
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

			log.Printf("[DEBUG] Registering Resource Provider %q..", resourceId.ResourceProvider)
			if _, err := client.Register(ctx, resourceId.ResourceProvider); err != nil {
				return fmt.Errorf("registering Resource Provider %q: %+v", resourceId.ResourceProvider, err)
			}

			// TODO: @tombuildsstuff - expose a nicer means of doing this in the SDK
			log.Printf("[DEBUG] Waiting for Resource Provider %q to finish registering..", resourceId.ResourceProvider)
			stateConf := &pluginsdk.StateChangeConf{
				Pending:      []string{"Processing"},
				Target:       []string{"Registered"},
				Refresh:      r.registerRefreshFunc(ctx, client, resourceId.ResourceProvider),
				MinTimeout:   15 * time.Second,
				PollInterval: 30 * time.Second,
				Timeout:      metadata.ResourceData.Timeout(pluginsdk.TimeoutCreate),
			}
			if _, err := stateConf.WaitForStateContext(ctx); err != nil {
				return fmt.Errorf("waiting for Resource Provider Namespace %q to be registered: %s", resourceId.ResourceProvider, err)
			}
			log.Printf("[DEBUG] Registered Resource Provider %q.", resourceId.ResourceProvider)

			metadata.SetID(resourceId)
			return nil
		},
		Timeout: 30 * time.Minute,
	}
}

func (r ResourceProviderRegistrationResource) Read() sdk.ResourceFunc {
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

			return metadata.Encode(&ResourceProviderRegistrationModel{
				Name: id.ResourceProvider,
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

			if _, err := client.Unregister(ctx, id.ResourceProvider); err != nil {
				return fmt.Errorf("unregistering Resource Provider %q: %+v", id.ResourceProvider, err)
			}

			// TODO: @tombuildsstuff - we should likely expose something in the SDK to make this easier

			stateConf := &pluginsdk.StateChangeConf{
				Pending:    []string{"Processing"},
				Target:     []string{"Unregistered"},
				Refresh:    r.unregisterRefreshFunc(ctx, client, id.ResourceProvider),
				MinTimeout: 15 * time.Second,
				Timeout:    metadata.ResourceData.Timeout(pluginsdk.TimeoutDelete),
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
