// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package federated_identity_credential

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/managedidentity/2024-11-30/federatedidentitycredentials"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ManagedIdentity.V20241130.FederatedIdentityCredentials

			id, err := federatedidentitycredentials.ParseFederatedIdentityCredentialID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(*id)
				}
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			return r.flatten(metadata, id, resp.Model)
		},
	}
}

type FederatedIdentityCredentialResourceSchema struct {
	Audience               []string `tfschema:"audience"`
	Issuer                 string   `tfschema:"issuer"`
	Name                   string   `tfschema:"name"`
	UserAssignedIdentityId string   `tfschema:"user_assigned_identity_id"`
	Subject                string   `tfschema:"subject"`
}

func (r Resource) mapFederatedIdentityCredentialResourceSchemaToFederatedIdentityCredentialProperties(input FederatedIdentityCredentialResourceSchema, output *federatedidentitycredentials.FederatedIdentityCredentialProperties) {
	output.Audiences = input.Audience
	output.Issuer = input.Issuer
	output.Subject = input.Subject
}

func (r Resource) mapFederatedIdentityCredentialPropertiesToFederatedIdentityCredentialResourceSchema(input federatedidentitycredentials.FederatedIdentityCredentialProperties, output *FederatedIdentityCredentialResourceSchema) {
	output.Audience = input.Audiences
	output.Issuer = input.Issuer
	output.Subject = input.Subject
}

func (r Resource) mapFederatedIdentityCredentialResourceSchemaToFederatedIdentityCredential(input FederatedIdentityCredentialResourceSchema, output *federatedidentitycredentials.FederatedIdentityCredential) {
	if output.Properties == nil {
		output.Properties = &federatedidentitycredentials.FederatedIdentityCredentialProperties{}
	}
	r.mapFederatedIdentityCredentialResourceSchemaToFederatedIdentityCredentialProperties(input, output.Properties)
}

func (r Resource) mapFederatedIdentityCredentialToFederatedIdentityCredentialResourceSchema(input federatedidentitycredentials.FederatedIdentityCredential, output *FederatedIdentityCredentialResourceSchema) {
	if input.Properties == nil {
		input.Properties = &federatedidentitycredentials.FederatedIdentityCredentialProperties{}
	}
	r.mapFederatedIdentityCredentialPropertiesToFederatedIdentityCredentialResourceSchema(*input.Properties, output)
}
