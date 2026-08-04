// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package elastic

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/elastic/2025-06-01/elasticmonitorresources"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/elastic/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

//go:generate go run ../../tools/generator-tests resourceidentity -resource-name elastic_cloud_serverless -service-package-name elastic -properties "resource_group_name,name" -known-values "subscription_id:data.Subscriptions.Primary"

type ElasticCloudServerlessResource struct{}

var (
	_ sdk.ResourceWithIdentity = ElasticCloudServerlessResource{}
	_ sdk.ResourceWithUpdate   = ElasticCloudServerlessResource{}
)

type ElasticCloudServerlessResourceModel struct {
	Name                     string            `tfschema:"name"`
	ResourceGroupName        string            `tfschema:"resource_group_name"`
	Location                 string            `tfschema:"location"`
	Kind                     string            `tfschema:"kind"`
	SkuName                  string            `tfschema:"sku_name"`
	ProjectType              string            `tfschema:"project_type"`
	ConfigurationType        string            `tfschema:"configuration_type"`
	OfferID                  string            `tfschema:"offer_id"`
	TermID                   string            `tfschema:"term_id"`
	ElasticCloudEmailAddress string            `tfschema:"elastic_cloud_email_address"`
	GenerateAPIKey           bool              `tfschema:"generate_api_key"`
	MonitoringEnabled        bool              `tfschema:"monitoring_enabled"`
	PlanID                   string            `tfschema:"plan_id"`
	PublisherID              string            `tfschema:"publisher_id"`
	ElasticCloudDeploymentID string            `tfschema:"elastic_cloud_deployment_id"`
	ElasticsearchServiceURL  string            `tfschema:"elasticsearch_service_url"`
	KibanaServiceURL         string            `tfschema:"kibana_service_url"`
	KibanaSSOURI             string            `tfschema:"kibana_sso_uri"`
	Tags                     map[string]string `tfschema:"tags"`
}

func (r ElasticCloudServerlessResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.ElasticName,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"kind": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"sku_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"project_type": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.StringInSlice([]string{
				string(elasticmonitorresources.ProjectTypeElasticsearch),
				string(elasticmonitorresources.ProjectTypeObservability),
				string(elasticmonitorresources.ProjectTypeSecurity),
			}, false),
		},

		"configuration_type": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringInSlice(elasticmonitorresources.PossibleValuesForConfigurationType(), false),
		},

		"offer_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"term_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"elastic_cloud_email_address": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.IsEmailAddress,
		},

		"generate_api_key": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
			ForceNew: true,
		},

		"monitoring_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  true,
			ForceNew: true,
		},

		"plan_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Default:      "ess-consumption-2024",
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"publisher_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Default:      "elastic",
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"tags": commonschema.Tags(),
	}
}

func (r ElasticCloudServerlessResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"elastic_cloud_deployment_id": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"elasticsearch_service_url": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"kibana_service_url": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"kibana_sso_uri": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}

func (r ElasticCloudServerlessResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 60 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model ElasticCloudServerlessResourceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			client := metadata.Client.Elastic.ServerlessMonitorClient
			id := elasticmonitorresources.NewMonitorID(metadata.Client.Account.SubscriptionId, model.ResourceGroupName, model.Name)

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.MonitorsGet(ctx, id)
				if err != nil && !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for existing %s: %+v", id, err)
				}
				if !response.WasNotFound(existing.HttpResponse) {
					return metadata.ResourceRequiresImport(r.ResourceType(), id)
				}
			}

			monitoringStatus := elasticmonitorresources.MonitoringStatusDisabled
			if model.MonitoringEnabled {
				monitoringStatus = elasticmonitorresources.MonitoringStatusEnabled
			}

			body := elasticmonitorresources.ElasticMonitorResource{
				Kind:     pointer.To(model.Kind),
				Location: location.Normalize(model.Location),
				Properties: &elasticmonitorresources.MonitorProperties{
					GenerateApiKey:   pointer.To(model.GenerateAPIKey),
					HostingType:      pointer.To(elasticmonitorresources.HostingTypeServerless),
					MonitoringStatus: pointer.To(monitoringStatus),
					PlanDetails: &elasticmonitorresources.PlanDetails{
						OfferID:     pointer.To(model.OfferID),
						PlanID:      pointer.To(model.PlanID),
						PublisherID: pointer.To(model.PublisherID),
						TermID:      pointer.To(model.TermID),
					},
					ProjectDetails: &elasticmonitorresources.ProjectDetails{
						ConfigurationType: pointer.ToEnum[elasticmonitorresources.ConfigurationType](model.ConfigurationType),
						ProjectType:       pointer.ToEnum[elasticmonitorresources.ProjectType](model.ProjectType),
					},
					UserInfo: &elasticmonitorresources.UserInfo{
						EmailAddress: pointer.To(model.ElasticCloudEmailAddress),
					},
				},
				Sku: &elasticmonitorresources.ResourceSku{
					Name: model.SkuName,
				},
				Tags: pointer.To(model.Tags),
			}

			if err := client.MonitorsCreateCallbackThenPoll(ctx, id, body, metadata.SetIDAndIdentityCallback(&id)); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			if err := pluginsdk.SetResourceIdentityData(metadata.ResourceData, &id); err != nil {
				return err
			}

			return nil
		},
	}
}

func (r ElasticCloudServerlessResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Elastic.ServerlessMonitorClient
			id, err := elasticmonitorresources.ParseMonitorID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.MonitorsGet(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			return r.flatten(metadata, id, resp.Model)
		},
	}
}

func (r ElasticCloudServerlessResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 60 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model ElasticCloudServerlessResourceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			client := metadata.Client.Elastic.ServerlessMonitorClient
			id, err := elasticmonitorresources.ParseMonitorID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			existing, err := client.MonitorsGet(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}
			if err := validateElasticCloudServerlessMonitor(id, existing.Model); err != nil {
				return err
			}

			payload := elasticmonitorresources.ElasticMonitorResourceUpdateParameters{
				Tags: pointer.To(model.Tags),
			}
			if err := client.MonitorsUpdateThenPoll(ctx, *id, payload); err != nil {
				return fmt.Errorf("updating %s: %+v", id, err)
			}

			return nil
		},
	}
}

func (r ElasticCloudServerlessResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 60 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Elastic.ServerlessMonitorClient
			id, err := elasticmonitorresources.ParseMonitorID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			existing, err := client.MonitorsGet(ctx, *id)
			if err != nil {
				if response.WasNotFound(existing.HttpResponse) {
					return nil
				}
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}
			if err := validateElasticCloudServerlessMonitor(id, existing.Model); err != nil {
				return err
			}

			if err := client.MonitorsDeleteThenPoll(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", id, err)
			}

			return nil
		},
	}
}

func (r ElasticCloudServerlessResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return elasticmonitorresources.ValidateMonitorID
}

func (r ElasticCloudServerlessResource) Identity() resourceids.ResourceId {
	return &elasticmonitorresources.MonitorId{}
}

func (r ElasticCloudServerlessResource) ModelObject() interface{} {
	return &ElasticCloudServerlessResourceModel{}
}

func (r ElasticCloudServerlessResource) ResourceType() string {
	return "azurerm_elastic_cloud_serverless"
}

func (r ElasticCloudServerlessResource) flatten(metadata sdk.ResourceMetaData, id *elasticmonitorresources.MonitorId, model *elasticmonitorresources.ElasticMonitorResource) error {
	if err := validateElasticCloudServerlessMonitor(id, model); err != nil {
		return err
	}

	state := ElasticCloudServerlessResourceModel{
		Name:              id.MonitorName,
		ResourceGroupName: id.ResourceGroupName,
		Kind:              pointer.From(model.Kind),
		Location:          location.Normalize(model.Location),
		Tags:              pointer.From(model.Tags),
	}

	if model.Sku != nil {
		state.SkuName = model.Sku.Name
	}

	props := model.Properties
	state.GenerateAPIKey = pointer.From(props.GenerateApiKey)
	state.MonitoringEnabled = pointer.From(props.MonitoringStatus) == elasticmonitorresources.MonitoringStatusEnabled

	if plan := props.PlanDetails; plan != nil {
		state.OfferID = pointer.From(plan.OfferID)
		state.PlanID = pointer.From(plan.PlanID)
		state.PublisherID = pointer.From(plan.PublisherID)
		state.TermID = pointer.From(plan.TermID)
	}

	if project := props.ProjectDetails; project != nil {
		state.ConfigurationType = string(pointer.From(project.ConfigurationType))
		state.ProjectType = string(pointer.From(project.ProjectType))
	}

	if user := props.UserInfo; user != nil {
		state.ElasticCloudEmailAddress = pointer.From(user.EmailAddress)
	}

	if elasticProperties := props.ElasticProperties; elasticProperties != nil {
		if user := elasticProperties.ElasticCloudUser; user != nil {
			state.ElasticCloudEmailAddress = pointer.From(user.EmailAddress)
		}
		if deployment := elasticProperties.ElasticCloudDeployment; deployment != nil {
			state.ElasticCloudDeploymentID = pointer.From(deployment.DeploymentId)
			state.ElasticsearchServiceURL = pointer.From(deployment.ElasticsearchServiceURL)
			state.KibanaServiceURL = pointer.From(deployment.KibanaServiceURL)
			state.KibanaSSOURI = pointer.From(deployment.KibanaSsoURL)
		}
	}

	if err := pluginsdk.SetResourceIdentityData(metadata.ResourceData, id); err != nil {
		return err
	}

	return metadata.Encode(&state)
}

func validateElasticCloudServerlessMonitor(id *elasticmonitorresources.MonitorId, model *elasticmonitorresources.ElasticMonitorResource) error {
	if model == nil {
		return fmt.Errorf("retrieving %s: model was nil", id)
	}
	if model.Properties == nil {
		return fmt.Errorf("retrieving %s: `properties` was nil", id)
	}
	if model.Properties.HostingType == nil {
		return fmt.Errorf("retrieving %s: `properties.hostingType` was nil", id)
	}
	if *model.Properties.HostingType != elasticmonitorresources.HostingTypeServerless {
		return fmt.Errorf("expected %s to use `Serverless` hosting, got `%s`", id, *model.Properties.HostingType)
	}

	return nil
}
