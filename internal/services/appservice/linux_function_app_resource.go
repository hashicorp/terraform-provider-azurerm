package appservice

import (
	"context"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/services/web/mgmt/2021-02-01/web"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appservice/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
	"time"

	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type LinuxFunctionAppResource struct{}

type LinuxFunctionAppModel struct {
	Name                  string `tfschema:"name"`
	ResourceGroup         string `tfschema:"resource_group_name"`
	Location              string `tfschema:"location"`
	ServicePlanId         string `tfschema:"service_plan_id"`
	ClientAffinityEnabled bool   `tfschema:"client_affinity_enabled"`
	ClientCertEnabled     bool   `tfschema:"client_cert_enabled"`
	ClientCertMode        string `tfschema:"client_cert_mode"`
	Enabled               bool   `tfschema:"enabled"`
	HttpsOnly             bool   `tfschema:"https_only"`

	Tags map[string]string `tfschema:"tags"`
}

var _ sdk.ResourceWithUpdate = LinuxFunctionAppResource{}

func (r LinuxFunctionAppResource) ModelObject() interface{} {
	return LinuxFunctionAppModel{}
}

func (r LinuxFunctionAppResource) ResourceType() string {
	return "azurerm_linux_function_app"
}

func (r LinuxFunctionAppResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	panic("Implement me") // TODO - Add Validation func return here
}

func (r LinuxFunctionAppResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		/*
		   TODO - This sections is for configurable items, `Required: true` items first, followed by `Optional: true`,
		   both in alphabetical order
		*/
	}
}

func (r LinuxFunctionAppResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		/*
		   TODO - This section is for `Computed: true` only items, i.e. useful values that are returned by the
		   datasource that can be used as outputs or passed programmatically to other resources or data sources.
		*/
	}
}

func (r LinuxFunctionAppResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var functionApp LinuxFunctionAppModel

			if err := metadata.Decode(&functionApp); err != nil {
				return err
			}

			client := metadata.Client.AppService.WebAppsClient
			aseClient := metadata.Client.AppService.AppServiceEnvironmentClient
			servicePlanClient := metadata.Client.AppService.ServicePlanClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			id := parse.NewWebAppID(subscriptionId, functionApp.ResourceGroup, functionApp.Name)

			existing, err := client.Get(ctx, id.ResourceGroup, id.SiteName)
			if err != nil && !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing Linux %s: %+v", id, err)
			}

			if !utils.ResponseWasNotFound(existing.Response) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			availabilityRequest := web.ResourceNameAvailabilityRequest{
				Name: utils.String(functionApp.Name),
				Type: web.CheckNameResourceTypesMicrosoftWebsites,
			}

			servicePlanId, err := parse.ServicePlanID(functionApp.ServicePlanId)
			if err != nil {
				return err
			}

			servicePlan, err := servicePlanClient.Get(ctx, servicePlanId.ResourceGroup, servicePlanId.ServerfarmName)
			if err != nil {
				return fmt.Errorf("reading %s: %+v", servicePlanId, err)
			}
			if ase := servicePlan.HostingEnvironmentProfile; ase != nil {
				// Attempt to check the ASE for the appropriate suffix for the name availability request.
				// This varies between internal and external ASE Types, and potentially has other names in other clouds
				// We use the "internal" as the fallback here, if we can read the ASE, we'll get the full one
				nameSuffix := "appserviceenvironment.net"
				if ase.ID != nil {
					aseId, err := parse.AppServiceEnvironmentID(*ase.ID)
					nameSuffix = fmt.Sprintf("%s.%s", aseId.HostingEnvironmentName, nameSuffix)
					if err != nil {
						metadata.Logger.Warnf("could not parse App Service Environment ID determine FQDN for name availability check, defaulting to `%s.%s.appserviceenvironment.net`", functionApp.Name, servicePlanId)
					} else {
						existingASE, err := aseClient.Get(ctx, aseId.ResourceGroup, aseId.HostingEnvironmentName)
						if err != nil {
							metadata.Logger.Warnf("could not read App Service Environment to determine FQDN for name availability check, defaulting to `%s.%s.appserviceenvironment.net`", functionApp.Name, servicePlanId)
						} else if props := existingASE.AppServiceEnvironment; props != nil && props.DNSSuffix != nil && *props.DNSSuffix != "" {
							nameSuffix = *props.DNSSuffix
						}
					}
				}

				availabilityRequest.Name = utils.String(fmt.Sprintf("%s.%s", functionApp.Name, nameSuffix))
				availabilityRequest.IsFqdn = utils.Bool(true)
			}

			checkName, err := client.CheckNameAvailability(ctx, availabilityRequest)
			if err != nil {
				return fmt.Errorf("checking name availability for Linux %s: %+v", id, err)
			}
			if !*checkName.NameAvailable {
				return fmt.Errorf("the Site Name %q failed the availability check: %+v", id.SiteName, *checkName.Message)
			}

			siteEnvelope := web.Site{
				Location: utils.String(functionApp.Location),
				Tags:     tags.FromTypedObject(functionApp.Tags),
				SiteProperties: &web.SiteProperties{
					ServerFarmID: utils.String(functionApp.ServicePlanId),
					Enabled:      utils.Bool(functionApp.Enabled),
					HTTPSOnly:    utils.Bool(functionApp.HttpsOnly),
					//SiteConfig:            siteConfig,
					ClientAffinityEnabled: utils.Bool(functionApp.ClientAffinityEnabled),
					ClientCertEnabled:     utils.Bool(functionApp.ClientCertEnabled),
					ClientCertMode:        web.ClientCertMode(functionApp.ClientCertMode),
				},
			}

			future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.SiteName, siteEnvelope)
			if err != nil {
				return fmt.Errorf("creating Linux %s: %+v", id, err)
			}

			if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
				return fmt.Errorf("waiting for creation of Linux %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r LinuxFunctionAppResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			// TODO - Read Func
			return nil
		},
	}
}

func (r LinuxFunctionAppResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			// TODO - Delete Func
			return nil
		},
	}
}

func (r LinuxFunctionAppResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			// TODO - Delete Func
			return nil
		},
	}
}
