// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package containerapps

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerapps/2023-05-01/certificates"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerapps/2024-03-01/managedenvironments"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/containerapps/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type ContainerAppEnvironmentManagedCertificateDataSource struct{}

type ContainerAppEnvironmentManagedCertificateDataSourceModel struct {
	Name                 string `tfschema:"name"`
	ManagedEnvironmentId string `tfschema:"container_app_environment_id"`

	// Read Only
	SubjectName             string                 `tfschema:"subject_name"`
	DomainControlValidation string                 `tfschema:"domain_control_validation_type"`
	Tags                    map[string]interface{} `tfschema:"tags"`
}

var _ sdk.DataSource = ContainerAppEnvironmentManagedCertificateDataSource{}

func (r ContainerAppEnvironmentManagedCertificateDataSource) ModelObject() interface{} {
	return &ContainerAppEnvironmentManagedCertificateDataSourceModel{}
}

func (r ContainerAppEnvironmentManagedCertificateDataSource) ResourceType() string {
	return "azurerm_container_app_environment_managed_certificate"
}

func (r ContainerAppEnvironmentManagedCertificateDataSource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.CertificateName,
			Description:  "The name of the Container Apps Certificate.",
		},

		"container_app_environment_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: certificates.ValidateManagedEnvironmentID,
			Description:  "The Container App Managed Environment ID.",
		},
	}
}

func (r ContainerAppEnvironmentManagedCertificateDataSource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"subject_name": {
			Type:        pluginsdk.TypeString,
			Computed:    true,
			Description: "The Subject Name for the Certificate.",
		},

		"domain_control_validation_type": {
			Type:        pluginsdk.TypeString,
			Computed:    true,
			Description: "Type of domain control validation.",
		},

		"tags": commonschema.TagsDataSource(),
	}
}

func (r ContainerAppEnvironmentManagedCertificateDataSource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ContainerApps.ManagedEnvironmentClient

			var cert ContainerAppEnvironmentManagedCertificateDataSourceModel
			if err := metadata.Decode(&cert); err != nil {
				return err
			}

			envId, err := certificates.ParseManagedEnvironmentID(cert.ManagedEnvironmentId)
			if err != nil {
				return err
			}

			id := managedenvironments.NewManagedCertificateID(envId.SubscriptionId, envId.ResourceGroupName, envId.ManagedEnvironmentName, cert.Name)

			existing, err := client.ManagedCertificatesGet(ctx, id)
			if err != nil {
				if response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("%s was not found", id)
				}
				return fmt.Errorf("reading %s: %+v", id, err)
			}

			cert.Name = id.ManagedCertificateName
			cert.ManagedEnvironmentId = envId.ID()

			if model := existing.Model; model != nil {
				cert.Tags = tags.Flatten(model.Tags)

				if props := model.Properties; props != nil {
					cert.SubjectName = pointer.From(props.SubjectName)
					cert.DomainControlValidation = string(pointer.From(props.DomainControlValidation))
				}
			}

			metadata.SetID(id)

			return metadata.Encode(&cert)
		},
	}
}
