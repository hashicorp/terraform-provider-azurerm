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
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerapps/2025-01-01/certificates"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/containerapps/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type ContainerAppEnvironmentCertificateDataSource struct{}

type ContainerAppEnvironmentCertificateDataSourceModel struct {
	Name                 string `tfschema:"name"`
	ManagedEnvironmentId string `tfschema:"container_app_environment_id"`

	// Read Only
	SubjectName    string                 `tfschema:"subject_name"`
	Issuer         string                 `tfschema:"issuer"`
	IssueDate      string                 `tfschema:"issue_date"`
	ExpirationDate string                 `tfschema:"expiration_date"`
	Thumbprint     string                 `tfschema:"thumbprint"`
	Tags           map[string]interface{} `tfschema:"tags"`
}

var _ sdk.DataSource = ContainerAppEnvironmentCertificateDataSource{}

func (r ContainerAppEnvironmentCertificateDataSource) ModelObject() interface{} {
	return &ContainerAppEnvironmentCertificateDataSourceModel{}
}

func (r ContainerAppEnvironmentCertificateDataSource) ResourceType() string {
	return "azurerm_container_app_environment_certificate"
}

func (r ContainerAppEnvironmentCertificateDataSource) Arguments() map[string]*pluginsdk.Schema {
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
			Description:  "The Container App Managed Environment ID to configure this Certificate on.",
		},
	}
}

func (r ContainerAppEnvironmentCertificateDataSource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"subject_name": {
			Type:        pluginsdk.TypeString,
			Computed:    true,
			Description: "The Subject Name for the Certificate.",
		},

		"issuer": {
			Type:        pluginsdk.TypeString,
			Computed:    true,
			Description: "The Certificate Issuer.",
		},

		"issue_date": {
			Type:        pluginsdk.TypeString,
			Computed:    true,
			Description: "The date of issue for the Certificate.",
		},

		"expiration_date": {
			Type:        pluginsdk.TypeString,
			Computed:    true,
			Description: "The expiration date for the Certificate.",
		},

		"thumbprint": {
			Type:        pluginsdk.TypeString,
			Computed:    true,
			Description: "The Thumbprint of the Certificate.",
		},

		"tags": commonschema.TagsDataSource(),
	}
}

func (r ContainerAppEnvironmentCertificateDataSource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ContainerApps.CertificatesClient

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
