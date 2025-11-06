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
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerapps/2025-07-01/managedenvironments"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/containerapps/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type ContainerAppEnvironmentManagedCertificateResource struct{}

type ContainerAppEnvironmentManagedCertificateModel struct {
	Name                      string                 `tfschema:"name"`
	ContainerAppEnvironmentId string                 `tfschema:"container_app_environment_id"`
	SubjectName               string                 `tfschema:"subject_name"`
	DomainControlValidation   string                 `tfschema:"domain_control_validation"`
	Tags                      map[string]interface{} `tfschema:"tags"`
	ValidationToken           string                 `tfschema:"validation_token"`
}

var _ sdk.ResourceWithUpdate = ContainerAppEnvironmentManagedCertificateResource{}

func (r ContainerAppEnvironmentManagedCertificateResource) ModelObject() interface{} {
	return &ContainerAppEnvironmentManagedCertificateModel{}
}

func (r ContainerAppEnvironmentManagedCertificateResource) ResourceType() string {
	return "azurerm_container_app_environment_managed_certificate"
}

func (r ContainerAppEnvironmentManagedCertificateResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return managedenvironments.ValidateManagedCertificateID
}

func (r ContainerAppEnvironmentManagedCertificateResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.CertificateName,
			Description:  "The name of the Container Apps Managed Certificate.",
		},

		"container_app_environment_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: managedenvironments.ValidateManagedEnvironmentID,
			Description:  "The Container App Managed Environment ID to configure this Managed Certificate on.",
		},

		"subject_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
			Description:  "The Subject Name of the Certificate. Must be a valid domain name.",
		},

		"domain_control_validation": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			ForceNew: true,
			Default:  string(managedenvironments.ManagedCertificateDomainControlValidationHTTP),
			ValidateFunc: validation.StringInSlice(
				managedenvironments.PossibleValuesForManagedCertificateDomainControlValidation(),
				false,
			),
			Description: "The domain control validation type for the managed certificate. Possible values are `CNAME`, `HTTP` and `TXT`. Defaults to `HTTP`.",
		},

		"tags": commonschema.Tags(),
	}
}

func (r ContainerAppEnvironmentManagedCertificateResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"validation_token": {
			Type:        pluginsdk.TypeString,
			Computed:    true,
			Description: "The validation token for the managed certificate.",
		},
	}
}

func (r ContainerAppEnvironmentManagedCertificateResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ContainerApps.ManagedEnvironmentClient

			var model ContainerAppEnvironmentManagedCertificateModel

			if err := metadata.Decode(&model); err != nil {
				return err
			}

			envId, err := managedenvironments.ParseManagedEnvironmentID(model.ContainerAppEnvironmentId)
			if err != nil {
				return err
			}

			id := managedenvironments.NewManagedCertificateID(envId.SubscriptionId, envId.ResourceGroupName, envId.ManagedEnvironmentName, model.Name)

			existing, err := client.ManagedCertificatesGet(ctx, id)
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			env, err := client.Get(ctx, *envId)
			if err != nil {
				return fmt.Errorf("reading %s for %s: %+v", *envId, id, err)
			}

			certificate := managedenvironments.ManagedCertificate{
				Location: env.Model.Location,
				Properties: &managedenvironments.ManagedCertificateProperties{
					SubjectName:             pointer.To(model.SubjectName),
					DomainControlValidation: pointer.To(managedenvironments.ManagedCertificateDomainControlValidation(model.DomainControlValidation)),
				},
				Tags: tags.Expand(model.Tags),
			}

			if err := client.ManagedCertificatesCreateOrUpdateThenPoll(ctx, id, certificate); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)

			return nil
		},
	}
}

func (r ContainerAppEnvironmentManagedCertificateResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ContainerApps.ManagedEnvironmentClient

			id, err := managedenvironments.ParseManagedCertificateID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			existing, err := client.ManagedCertificatesGet(ctx, *id)
			if err != nil {
				if response.WasNotFound(existing.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("reading %s: %+v", *id, err)
			}

			var state ContainerAppEnvironmentManagedCertificateModel

			state.Name = id.ManagedCertificateName
			state.ContainerAppEnvironmentId = managedenvironments.NewManagedEnvironmentID(id.SubscriptionId, id.ResourceGroupName, id.ManagedEnvironmentName).ID()

			if model := existing.Model; model != nil {
				state.Tags = tags.Flatten(model.Tags)

				if props := model.Properties; props != nil {
					state.SubjectName = pointer.From(props.SubjectName)
					state.DomainControlValidation = string(pointer.From(props.DomainControlValidation))
					state.ValidationToken = pointer.From(props.ValidationToken)
				}
			}

			return metadata.Encode(&state)
		},
	}
}

func (r ContainerAppEnvironmentManagedCertificateResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ContainerApps.ManagedEnvironmentClient

			id, err := managedenvironments.ParseManagedCertificateID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if _, err := client.ManagedCertificatesDelete(ctx, *id); err != nil {
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
			client := metadata.Client.ContainerApps.ManagedEnvironmentClient

			var model ContainerAppEnvironmentManagedCertificateModel

			if err := metadata.Decode(&model); err != nil {
				return err
			}

			id, err := managedenvironments.ParseManagedCertificateID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if metadata.ResourceData.HasChange("tags") {
				patch := managedenvironments.ManagedCertificatePatch{
					Tags: tags.Expand(model.Tags),
				}

				if _, err = client.ManagedCertificatesUpdate(ctx, *id, patch); err != nil {
					return fmt.Errorf("updating tags for %s: %+v", *id, err)
				}
			}

			return nil
		},
	}
}
