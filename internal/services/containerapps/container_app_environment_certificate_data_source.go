package containerapps

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerapps/2022-03-01/certificates"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/containerapps/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
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
	return "azurerm_container_app_environment_certificate_resource"
}

func (r ContainerAppEnvironmentCertificateDataSource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: helpers.ValidateCertificateName,
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
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"issuer": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"issue_date": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"expiration_date": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"thumbprint": {
			Type:     pluginsdk.TypeString,
			Computed: true,
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

			id := certificates.NewCertificateID(envId.SubscriptionId, envId.ResourceGroupName, envId.EnvironmentName, cert.Name)

			existing, err := client.Get(ctx, id)
			if err != nil {
				if response.WasNotFound(existing.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("reading %s: %+v", id, err)
			}

			cert.Name = id.CertificateName
			cert.ManagedEnvironmentId = envId.ID()

			if model := existing.Model; model != nil {
				cert.Tags = tags.Flatten(model.Tags)

				if props := model.Properties; props != nil {
					cert.Issuer = utils.NormalizeNilableString(props.Issuer)
					cert.IssueDate = utils.NormalizeNilableString(props.IssueDate)
					cert.ExpirationDate = utils.NormalizeNilableString(props.ExpirationDate)
					cert.Thumbprint = utils.NormalizeNilableString(props.Thumbprint)
				}
			}

			metadata.SetID(id)

			return metadata.Encode(&cert)
		},
	}
}
