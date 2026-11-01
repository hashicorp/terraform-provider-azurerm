// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package load_test

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/loadtestservice/2022-12-01/loadtests"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.LoadTestService.V20221201.LoadTests
			schema := LoadTestResourceSchema{}

			id, err := loadtests.ParseLoadTestID(metadata.ResourceData.Id())
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

			if model := resp.Model; model != nil {
				schema.Name = id.LoadTestName
				schema.ResourceGroupName = id.ResourceGroupName
				if err := r.mapLoadTestResourceToLoadTestResourceSchema(*model, &schema); err != nil {
					return fmt.Errorf("flattening model: %+v", err)
				}
				if property := model.Properties; property != nil {
					schema.Description = pointer.From(property.Description)
				}
			}

			return metadata.Encode(&schema)
		},
	}
}

type LoadTestResourceSchema struct {
	DataPlaneURI      string                                     `tfschema:"data_plane_uri"`
	Description       string                                     `tfschema:"description"`
	Encryption        []LoadTestEncryption                       `tfschema:"encryption"`
	Identity          []identity.ModelSystemAssignedUserAssigned `tfschema:"identity"`
	Location          string                                     `tfschema:"location"`
	Name              string                                     `tfschema:"name"`
	ResourceGroupName string                                     `tfschema:"resource_group_name"`
	Tags              map[string]any                             `tfschema:"tags"`
}

type LoadTestEncryption struct {
	KeyURL   string                       `tfschema:"key_url"`
	Identity []LoadTestEncryptionIdentity `tfschema:"identity"`
}

type LoadTestEncryptionIdentity struct {
	IdentityID string `tfschema:"identity_id"`
	Type       string `tfschema:"type"`
}

func (r Resource) EnsureEncryptionIdentityIDExistsInIdentity(model LoadTestResourceSchema) error {
	if len(model.Encryption) == 1 && len(model.Encryption[0].Identity) == 1 && model.Encryption[0].Identity[0].Type == string(loadtests.TypeUserAssigned) {
		msg := "when `encryption.identity.type` is set to `UserAssigned`, the `encryption.identity.identity_id` provided must also be specified in the `identity.identity_ids` list"
		if len(model.Identity) == 0 {
			return errors.New(msg)
		}

		existsInIdentity := slices.Contains(model.Identity[0].IdentityIds, model.Encryption[0].Identity[0].IdentityID)

		if !existsInIdentity {
			return errors.New(msg)
		}
	}

	return nil
}

//nolint:unparam
func (r Resource) mapLoadTestResourceSchemaToLoadTestProperties(input LoadTestResourceSchema, output *loadtests.LoadTestProperties) error {
	output.Description = &input.Description
	output.Encryption = r.mapLoadTestResourceSchemaToLoadTestEncryption(input.Encryption)

	return nil
}

func (r Resource) mapLoadTestResourceSchemaToLoadTestEncryption(input []LoadTestEncryption) *loadtests.EncryptionProperties {
	if len(input) == 0 || input[0].KeyURL == "" {
		return nil
	}

	attr := input[0]

	encryptionIdentity := &loadtests.EncryptionPropertiesIdentity{}
	if attrIdentity := attr.Identity; len(attrIdentity) > 0 {
		encryptionIdentity.ResourceId = pointer.To(attrIdentity[0].IdentityID)
		encryptionIdentity.Type = pointer.ToEnum[loadtests.Type](attrIdentity[0].Type)
	}

	return &loadtests.EncryptionProperties{
		KeyURL:   pointer.To(attr.KeyURL),
		Identity: encryptionIdentity,
	}
}

//nolint:unparam
func (r Resource) mapLoadTestPropertiesToLoadTestResourceSchema(input loadtests.LoadTestProperties, output *LoadTestResourceSchema) error {
	output.DataPlaneURI = pointer.From(input.DataPlaneURI)
	output.Description = pointer.From(input.Description)

	if encryption := input.Encryption; encryption != nil {
		output.Encryption = []LoadTestEncryption{{
			KeyURL:   pointer.From(encryption.KeyURL),
			Identity: make([]LoadTestEncryptionIdentity, 0),
		}}
		if encryptionIdentity := encryption.Identity; encryptionIdentity != nil {
			output.Encryption[0].Identity = append(output.Encryption[0].Identity, LoadTestEncryptionIdentity{
				IdentityID: pointer.From(encryptionIdentity.ResourceId),
			})

			if encryptionIdentity.Type != nil {
				output.Encryption[0].Identity[0].Type = pointer.FromEnum(encryptionIdentity.Type)
			}
		}
	}
	return nil
}

func (r Resource) mapLoadTestResourceToLoadTestResourceSchema(input loadtests.LoadTestResource, output *LoadTestResourceSchema) error {
	identity, err := identity.FlattenLegacySystemAndUserAssignedMapToModel(input.Identity)
	if err != nil {
		return fmt.Errorf("flattening Legacy SystemAndUserAssigned Identity: %+v", err)
	}
	output.Identity = identity

	output.Location = location.Normalize(input.Location)
	output.Tags = tags.Flatten(input.Tags)

	if input.Properties == nil {
		input.Properties = &loadtests.LoadTestProperties{}
	}
	if err := r.mapLoadTestPropertiesToLoadTestResourceSchema(*input.Properties, output); err != nil {
		return fmt.Errorf("mapping SDK Field %q / Model %q to Schema: %+v", "LoadTestProperties", "Properties", err)
	}

	return nil
}
