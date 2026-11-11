// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package container_app_environment_managed_certificate

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerapps/2025-07-01/managedenvironments"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Create() sdk.ResourceFunc {
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

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.ManagedCertificatesGet(ctx, id)
				if err != nil {
					if !response.WasNotFound(existing.HttpResponse) {
						return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
					}
				}

				if !response.WasNotFound(existing.HttpResponse) {
					return metadata.ResourceRequiresImport(r.ResourceType(), id)
				}
			}

			env, err := client.Get(ctx, *envId)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *envId, err)
			}

			if env.Model == nil {
				return fmt.Errorf("retrieving %s: `model` was nil", envId)
			}

			certificate := managedenvironments.ManagedCertificate{
				Location: env.Model.Location,
				Properties: &managedenvironments.ManagedCertificateProperties{
					SubjectName:             pointer.To(model.SubjectName),
					DomainControlValidation: pointer.ToEnum[managedenvironments.ManagedCertificateDomainControlValidation](model.DomainControlValidation),
				},
				Tags: tags.Expand(model.Tags),
			}

			if err := client.ManagedCertificatesCreateOrUpdateCallbackThenPoll(ctx, id, certificate, metadata.SetIDCallback(&id)); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)

			return nil
		},
	}
}
