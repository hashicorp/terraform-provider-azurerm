// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package container_app_environment_certificate

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerapps/2025-07-01/certificates"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	containerappsClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/containerapps/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/containerapps/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type DataSource struct{}

type ContainerAppEnvironmentCertificateDataSourceModel struct {
	Name                 string `tfschema:"name"`
	ManagedEnvironmentId string `tfschema:"container_app_environment_id"`

	// Read Only
	SubjectName    string         `tfschema:"subject_name"`
	Issuer         string         `tfschema:"issuer"`
	IssueDate      string         `tfschema:"issue_date"`
	ExpirationDate string         `tfschema:"expiration_date"`
	Thumbprint     string         `tfschema:"thumbprint"`
	Tags           map[string]any `tfschema:"tags"`
}

var _ sdk.DataSource = DataSource{}

func (r DataSource) ModelObject() any {
	return &ContainerAppEnvironmentCertificateDataSourceModel{}
}

func (r DataSource) ResourceType() string {
	return "azurerm_container_app_environment_certificate"
}

func (r DataSource) Description() string {
	return "Use this data source to access information about an existing Container App Environment Certificate."
}

func (r DataSource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description:  "The name of the Container Apps Certificate.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.CertificateName,
		},

		"container_app_environment_id": {
			Description:  "The ID of the Container App Environment to configure this Certificate on.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: certificates.ValidateManagedEnvironmentID,
		},
	}
}

func (r DataSource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"subject_name": {
			Description: "The Subject Name for the Certificate.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"issuer": {
			Description: "The Certificate Issuer.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"issue_date": {
			Description: "The date of issue for the Certificate.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"expiration_date": {
			Description: "The expiration date for the Certificate.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"thumbprint": {
			Description: "The Thumbprint of the Certificate.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"tags": commonschema.TagsDataSource(),
	}
}

func (r DataSource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := containerappsClient.From(metadata.Client).CertificatesClient

			var cert ContainerAppEnvironmentCertificateDataSourceModel
			if err := metadata.Decode(&cert); err != nil {
				return err
			}

			envId, err := certificates.ParseManagedEnvironmentID(cert.ManagedEnvironmentId)
			if err != nil {
				return err
			}

			id := certificates.NewCertificateID(envId.SubscriptionId, envId.ResourceGroupName, envId.ManagedEnvironmentName, cert.Name)

			existing, err := client.Get(ctx, id)
			if err != nil {
				if response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("%s was not found", id)
				}
				return fmt.Errorf("reading %s: %+v", id, err)
			}

			cert.Name = id.CertificateName
			cert.ManagedEnvironmentId = envId.ID()

			if model := existing.Model; model != nil {
				cert.Tags = tags.Flatten(model.Tags)

				if props := model.Properties; props != nil {
					cert.Issuer = pointer.From(props.Issuer)
					cert.IssueDate = pointer.From(props.IssueDate)
					cert.ExpirationDate = pointer.From(props.ExpirationDate)
					cert.Thumbprint = pointer.From(props.Thumbprint)
					cert.SubjectName = pointer.From(props.SubjectName)
				}
			}

			metadata.SetID(id)

			return metadata.Encode(&cert)
		},
	}
}
