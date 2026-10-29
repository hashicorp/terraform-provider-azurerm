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

func (r Resource) Read() sdk.ResourceFunc {
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
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			var state ContainerAppEnvironmentManagedCertificateModel

			state.Name = id.ManagedCertificateName
			state.ContainerAppEnvironmentId = managedenvironments.NewManagedEnvironmentID(id.SubscriptionId, id.ResourceGroupName, id.ManagedEnvironmentName).ID()

			if model := existing.Model; model != nil {
				state.Tags = tags.Flatten(model.Tags)

				if props := model.Properties; props != nil {
					state.SubjectName = pointer.From(props.SubjectName)
					state.DomainControlValidation = pointer.FromEnum(props.DomainControlValidation)
					state.ValidationToken = pointer.From(props.ValidationToken)
				}
			}

			return metadata.Encode(&state)
		},
	}
}
