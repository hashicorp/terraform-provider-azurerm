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
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerapps/2024-03-01/managedenvironments"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type ContainerAppEnvironmentManagedCertificateResource struct{}

var _ sdk.Resource = ContainerAppEnvironmentManagedCertificateResource{}

type ContainerAppEnvironmentManagedCertificateModel struct {
	Name                    string                 `tfschema:"name"`
	ManagedEnvironmentId    string                 `tfschema:"container_app_environment_id"`
	DomainControlValidation string                 `tfschema:"domain_control_validation_type"`
	SubjectName             string                 `tfschema:"subject_name"`
	Tags                    map[string]interface{} `tfschema:"tags"`
}

func (a ContainerAppEnvironmentManagedCertificateResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"container_app_environment_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: managedenvironments.ValidateManagedEnvironmentID,
			Description:  "The ID of the Container App Environment to host this Container App.",
		},

		"domain_control_validation_type": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringInSlice(managedenvironments.PossibleValuesForManagedCertificateDomainControlValidation(), false),
			Description:  "Type of domain control validation. Possible values include `CNAME`, `HTTP` and `TXT`.",
		},

		"subject_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
			Description:  "Subject name of the certificate.",
		},

		"tags": commonschema.Tags(),
	}
}

func (a ContainerAppEnvironmentManagedCertificateResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (a ContainerAppEnvironmentManagedCertificateResource) ModelObject() interface{} {
	return &ContainerAppEnvironmentManagedCertificateModel{}
}

func (a ContainerAppEnvironmentManagedCertificateResource) ResourceType() string {
	return "azurerm_container_app_environment_managed_certificate"
}

func (a ContainerAppEnvironmentManagedCertificateResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return managedenvironments.ValidateManagedCertificateID
}

func (a ContainerAppEnvironmentManagedCertificateResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			managedEnvironmentsClient := metadata.Client.ContainerApps.ManagedEnvironmentClient

			var cert ContainerAppEnvironmentManagedCertificateModel

			if err := metadata.Decode(&cert); err != nil {
				return err
			}

			managedEnvironmentId, err := managedenvironments.ParseManagedEnvironmentID(cert.ManagedEnvironmentId)
			if err != nil {
				return err
			}

			managedEnvironment, err := managedEnvironmentsClient.Get(ctx, *managedEnvironmentId)
			if err != nil {
				return fmt.Errorf("reading %s: %+v", *managedEnvironmentId, err)
			}

			managedCertId := managedenvironments.NewManagedCertificateID(metadata.Client.Account.SubscriptionId, managedEnvironmentId.ResourceGroupName, managedEnvironmentId.ManagedEnvironmentName, cert.Name)

			existingManagedCert, err := managedEnvironmentsClient.ManagedCertificatesGet(ctx, managedCertId)
			if err != nil {
				if !response.WasNotFound(existingManagedCert.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %+v", managedCertId, err)
				}
			}
			if !response.WasNotFound(existingManagedCert.HttpResponse) {
				return metadata.ResourceRequiresImport(a.ResourceType(), managedCertId)
			}

			domainControlValidation := managedenvironments.ManagedCertificateDomainControlValidation(cert.DomainControlValidation)

			model := managedenvironments.ManagedCertificate{
				Location: managedEnvironment.Model.Location,
				Name:     pointer.To(managedCertId.ManagedCertificateName),
				Properties: &managedenvironments.ManagedCertificateProperties{
					DomainControlValidation: &domainControlValidation,
					SubjectName:             pointer.To(cert.SubjectName),
				},
				Tags: tags.Expand(cert.Tags),
			}

			if err := managedEnvironmentsClient.ManagedCertificatesCreateOrUpdateThenPoll(ctx, managedCertId, model); err != nil {
				return fmt.Errorf("creating %s: %+v", managedCertId, err)
			}

			metadata.SetID(managedCertId)

			return nil
		},
	}
}

func (a ContainerAppEnvironmentManagedCertificateResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			managedEnvironmentsClient := metadata.Client.ContainerApps.ManagedEnvironmentClient

			id, err := managedenvironments.ParseManagedCertificateID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			existing, err := managedEnvironmentsClient.ManagedCertificatesGet(ctx, *id)
			if err != nil {
				if response.WasNotFound(existing.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("reading %s: %+v", *id, err)
			}

			var state ContainerAppEnvironmentManagedCertificateModel

			state.Name = id.ManagedCertificateName
			state.ManagedEnvironmentId = managedenvironments.NewManagedEnvironmentID(id.SubscriptionId, id.ResourceGroupName, id.ManagedEnvironmentName).ID()

			if model := existing.Model; model != nil {
				state.Tags = tags.Flatten(model.Tags)

				if props := model.Properties; props != nil {
					state.SubjectName = pointer.From(props.SubjectName)
					state.DomainControlValidation = string(pointer.From(props.DomainControlValidation))
				}
			}

			return metadata.Encode(&state)
		},
	}
}

func (a ContainerAppEnvironmentManagedCertificateResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			managedEnvironmentsClient := metadata.Client.ContainerApps.ManagedEnvironmentClient

			id, err := managedenvironments.ParseManagedCertificateID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if _, err := managedEnvironmentsClient.ManagedCertificatesDelete(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (r ContainerAppEnvironmentManagedCertificateResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			managedEnvironmentsClient := metadata.Client.ContainerApps.ManagedEnvironmentClient

			var cert ContainerAppEnvironmentManagedCertificateModel

			if err := metadata.Decode(&cert); err != nil {
				return err
			}

			id, err := managedenvironments.ParseManagedCertificateID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if metadata.ResourceData.HasChange("tags") {
				patch := managedenvironments.ManagedCertificatePatch{
					Tags: tags.Expand(cert.Tags),
				}

				if _, err = managedEnvironmentsClient.ManagedCertificatesUpdate(ctx, *id, patch); err != nil {
					return fmt.Errorf("updating tags for %s: %+v", *id, err)
				}
			}

			return nil
		},
	}
}
