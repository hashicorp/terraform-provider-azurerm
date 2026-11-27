// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package elastic_cloud_serverless

import (
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/elastic/2025-06-01/elasticmonitorresources"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func (r Resource) flatten(metadata sdk.ResourceMetaData, id *elasticmonitorresources.MonitorId, model *elasticmonitorresources.ElasticMonitorResource) error {
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
		state.Sku = model.Sku.Name
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
		state.ConfigurationType = pointer.FromEnum(project.ConfigurationType)
		state.ProjectType = pointer.FromEnum(project.ProjectType)
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
