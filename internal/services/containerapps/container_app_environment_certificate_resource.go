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
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerapps/2023-05-01/managedenvironments"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerapps/2024-03-01/managedcertificates"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/containerapps/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type ContainerAppEnvironmentCertificateResource struct{}

type ContainerAppCertificateModel struct {
	Name                 string                 `tfschema:"name"`
	ManagedEnvironmentId string                 `tfschema:"container_app_environment_id"`
	Tags                 map[string]interface{} `tfschema:"tags"`

	// Fields for Bring Your Own Certificate
	CertificatePassword string `tfschema:"certificate_password"`
	CertificateBlob     string `tfschema:"certificate_blob_base64"`

	// Fields for Managed Certificate
	SubjectName             string `tfschema:"subject_name"`
	DomainControlValidation string `tfschema:"domain_control_validation"`

	// Read Only BYO
	Issuer         string `tfschema:"issuer"`
	IssueDate      string `tfschema:"issue_date"`
	ExpirationDate string `tfschema:"expiration_date"`
	Thumbprint     string `tfschema:"thumbprint"`

	// Read Only Managed Certificate
	ValidationToken string `tfschema:"validation_token"`
}

var _ sdk.ResourceWithUpdate = ContainerAppEnvironmentCertificateResource{}

func (r ContainerAppEnvironmentCertificateResource) ModelObject() interface{} {
	return &ContainerAppCertificateModel{}
}

func (r ContainerAppEnvironmentCertificateResource) ResourceType() string {
	return "azurerm_container_app_environment_certificate"
}

func (r ContainerAppEnvironmentCertificateResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return certificates.ValidateCertificateID
}

func (r ContainerAppEnvironmentCertificateResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.CertificateName,
			Description:  "The name of the Container Apps Environment Certificate.",
		},

		"container_app_environment_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: certificates.ValidateManagedEnvironmentID,
			Description:  "The Container App Managed Environment ID to configure this Certificate on.",
		},

		"certificate_blob_base64": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsBase64,
			Description:  "The Certificate Private Key as a base64 encoded PFX or PEM. Required for BYO Certificate.",
		},

		"certificate_password": {
			Type:        pluginsdk.TypeString,
			Optional:    true,
			ForceNew:    true,
			Sensitive:   true,
			Description: "The password for the Certificate. Required for BYO Certificate.",
		},

		"subject_name": {
			Type:        pluginsdk.TypeString,
			Optional:    true,
			ForceNew:    true,
			Description: "The subject name of the managed certificate.",
		},

		"domain_control_validation": {
			Type:        pluginsdk.TypeString,
			Optional:    true,
			ForceNew:    true,
			Description: "The domain control validation method for the managed certificate.",
		},

		"tags": commonschema.Tags(),
	}
}

func (r ContainerAppEnvironmentCertificateResource) Attributes() map[string]*pluginsdk.Schema {
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
	}
}

func (r ContainerAppEnvironmentCertificateResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			certClient := metadata.Client.ContainerApps.CertificatesClient
			envClient := metadata.Client.ContainerApps.ManagedEnvironmentClient
			managedCertClient := metadata.Client.ContainerApps.ManagedCertificatesClient

			var cert ContainerAppCertificateModel

			if err := metadata.Decode(&cert); err != nil {
				return err
			}

			envId, err := managedenvironments.ParseManagedEnvironmentID(cert.ManagedEnvironmentId)
			if err != nil {
				return err
			}

			if cert.CertificateBlob != "" && cert.CertificatePassword != "" {
				// Create BYO Certificate
				id := certificates.NewCertificateID(metadata.Client.Account.SubscriptionId, envId.ResourceGroupName, envId.ManagedEnvironmentName, cert.Name)
				env, err := envClient.Get(ctx, *envId)
				if err != nil {
					return fmt.Errorf("reading %s for %s: %+v", *envId, id, err)
				}

				model := certificates.Certificate{
					Location: env.Model.Location,
					Name:     pointer.To(id.CertificateName),
					Properties: &certificates.CertificateProperties{
						Password: pointer.To(cert.CertificatePassword),
						Value:    pointer.To(cert.CertificateBlob),
					},
					Tags: tags.Expand(cert.Tags),
				}

				if _, err := certClient.CreateOrUpdate(ctx, id, model); err != nil {
					return fmt.Errorf("creating %s: %+v", id, err)
				}

				metadata.SetID(id)
			} else if cert.SubjectName != "" && cert.DomainControlValidation != "" {
				// Create Managed Certificate
				domainControlValidation, err := parseDomainControlValidation(cert.DomainControlValidation)
				if err != nil {
					return err
				}

				id := managedcertificates.NewManagedCertificateID(metadata.Client.Account.SubscriptionId, envId.ResourceGroupName, envId.ManagedEnvironmentName, cert.Name)
				env, err := envClient.Get(ctx, *envId)
				if err != nil {
					return fmt.Errorf("reading %s for %s: %+v", *envId, id, err)
				}

				model := managedcertificates.ManagedCertificate{
					Location: env.Model.Location,
					Properties: &managedcertificates.ManagedCertificateProperties{
						SubjectName:             pointer.To(cert.SubjectName),
						DomainControlValidation: domainControlValidation,
					},
					Tags: tags.Expand(cert.Tags),
				}

				if err := managedCertClient.CreateOrUpdateThenPoll(ctx, id, model); err != nil {
					return fmt.Errorf("creating %s: %+v", id, err)
				}

				metadata.SetID(id)
			} else {
				return fmt.Errorf("either certificate_blob_base64 and certificate_password or subject_name and domain_control_validation must be provided")
			}

			return nil
		},
	}
}
func (r ContainerAppEnvironmentCertificateResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			certClient := metadata.Client.ContainerApps.CertificatesClient
			managedCertClient := metadata.Client.ContainerApps.ManagedCertificatesClient

			id, err := certificates.ParseCertificateID(metadata.ResourceData.Id())
			if err != nil {
				id, err := managedcertificates.ParseManagedCertificateID(metadata.ResourceData.Id())
				if err != nil {
					return err
				}

				existing, err := managedCertClient.Get(ctx, *id)
				if err != nil {
					if response.WasNotFound(existing.HttpResponse) {
						return metadata.MarkAsGone(id)
					}
					return fmt.Errorf("reading %s: %+v", *id, err)
				}

				var state ContainerAppCertificateModel

				state.Name = id.ManagedCertificateName
				state.ManagedEnvironmentId = managedcertificates.NewManagedEnvironmentID(id.SubscriptionId, id.ResourceGroupName, id.ManagedEnvironmentName).ID()

				if model := existing.Model; model != nil {
					state.Tags = tags.Flatten(model.Tags)

					if props := model.Properties; props != nil {
						state.SubjectName = pointer.From(props.SubjectName)
						state.ValidationToken = pointer.From(props.ValidationToken)

					}
				}

				return metadata.Encode(&state)
			}

			existing, err := certClient.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(existing.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("reading %s: %+v", *id, err)
			}

			var state ContainerAppCertificateModel

			state.Name = id.CertificateName
			state.ManagedEnvironmentId = certificates.NewManagedEnvironmentID(id.SubscriptionId, id.ResourceGroupName, id.ManagedEnvironmentName).ID()

			if model := existing.Model; model != nil {
				state.Tags = tags.Flatten(model.Tags)

				if certBlob, ok := metadata.ResourceData.GetOk("certificate_blob_base64"); ok {
					state.CertificateBlob = certBlob.(string)
				}

				if certPassword, ok := metadata.ResourceData.GetOk("certificate_password"); ok {
					state.CertificatePassword = certPassword.(string)
				}

				if props := model.Properties; props != nil {
					state.SubjectName = pointer.From(props.SubjectName)
					state.Issuer = pointer.From(props.Issuer)
					state.IssueDate = pointer.From(props.IssueDate)
					state.ExpirationDate = pointer.From(props.ExpirationDate)
					state.Thumbprint = pointer.From(props.Thumbprint)
				}
			}

			return metadata.Encode(&state)
		},
	}
}
func (r ContainerAppEnvironmentCertificateResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			certClient := metadata.Client.ContainerApps.CertificatesClient
			managedCertClient := metadata.Client.ContainerApps.ManagedCertificatesClient

			id, err := certificates.ParseCertificateID(metadata.ResourceData.Id())
			if err != nil {
				id, err := managedcertificates.ParseManagedCertificateID(metadata.ResourceData.Id())
				if err != nil {
					return err
				}

				if _, err := managedCertClient.Delete(ctx, *id); err != nil {
					return fmt.Errorf("deleting %s: %+v", *id, err)
				}

				return nil
			}

			if _, err := certClient.Delete(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (r ContainerAppEnvironmentCertificateResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			certClient := metadata.Client.ContainerApps.CertificatesClient
			managedCertClient := metadata.Client.ContainerApps.ManagedCertificatesClient

			var cert ContainerAppCertificateModel

			if err := metadata.Decode(&cert); err != nil {
				return err
			}

			id, err := certificates.ParseCertificateID(metadata.ResourceData.Id())
			if err != nil {
				id, err := managedcertificates.ParseManagedCertificateID(metadata.ResourceData.Id())
				if err != nil {
					return err
				}

				if metadata.ResourceData.HasChange("tags") {
					patch := managedcertificates.ManagedCertificatePatch{
						Tags: tags.Expand(cert.Tags),
					}

					if _, err = managedCertClient.Update(ctx, *id, patch); err != nil {
						return fmt.Errorf("updating tags for %s: %+v", *id, err)
					}
				}

				return nil
			}

			if metadata.ResourceData.HasChange("tags") {
				patch := certificates.CertificatePatch{
					Tags: tags.Expand(cert.Tags),
				}

				if _, err = certClient.Update(ctx, *id, patch); err != nil {
					return fmt.Errorf("updating tags for %s: %+v", *id, err)
				}
			}

			return nil
		},
	}
}

func parseDomainControlValidation(input string) (*managedcertificates.ManagedCertificateDomainControlValidation, error) {
	vals := map[string]managedcertificates.ManagedCertificateDomainControlValidation{
		"CNAME": managedcertificates.ManagedCertificateDomainControlValidationCNAME,
		"HTTP":  managedcertificates.ManagedCertificateDomainControlValidationHTTP,
		"TXT":   managedcertificates.ManagedCertificateDomainControlValidationTXT,
	}
	if v, ok := vals[input]; ok {
		return &v, nil
	}
	return nil, fmt.Errorf("invalid DomainControlValidation value: %s", input)
}
