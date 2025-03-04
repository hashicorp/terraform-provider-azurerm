// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package loadtestservice

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/loadtestservice/2022-12-01/loadtests"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type LoadTestDataSource struct{}

var _ sdk.DataSource = LoadTestDataSource{}

type LoadTestDataSourceModel struct {
	DataPlaneURI      string                                     `tfschema:"data_plane_uri"`
	Description       string                                     `tfschema:"description"`
	Encryption        []LoadTestEncryption                       `tfschema:"encryption"`
	Identity          []identity.ModelSystemAssignedUserAssigned `tfschema:"identity"`
	Location          string                                     `tfschema:"location"`
	Name              string                                     `tfschema:"name"`
	ResourceGroupName string                                     `tfschema:"resource_group_name"`
	Tags              map[string]string                          `tfschema:"tags"`
}

func (r LoadTestDataSource) ModelObject() interface{} {
	return &LoadTestDataSourceModel{}
}

func (r LoadTestDataSource) ResourceType() string {
	return "azurerm_load_test"
}

func (r LoadTestDataSource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},

		"resource_group_name": commonschema.ResourceGroupNameForDataSource(),
	}
}

func (r LoadTestDataSource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"data_plane_uri": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"location": commonschema.LocationComputed(),

		"description": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"identity": commonschema.SystemAssignedUserAssignedIdentityComputed(),

		"encryption": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"key_url": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
					"identity": {
						Type:     pluginsdk.TypeList,
						Computed: true,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"type": {
									Type:     pluginsdk.TypeString,
									Computed: true,
								},
								"identity_id": {
									Type:     pluginsdk.TypeString,
									Computed: true,
								},
							},
						},
					},
				},
			},
		},

		"tags": tags.SchemaDataSource(),
	}
}

func (r LoadTestDataSource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.LoadTestService.V20221201.LoadTests
			subscriptionId := metadata.Client.Account.SubscriptionId

			var loadTest LoadTestDataSourceModel
			if err := metadata.Decode(&loadTest); err != nil {
				return err
			}

			id := loadtests.NewLoadTestID(subscriptionId, loadTest.ResourceGroupName, loadTest.Name)

			existing, err := client.Get(ctx, id)
			if err != nil {
				if response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("%s was not found", id)
				}
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			loadTest.Name = id.LoadTestName
			loadTest.ResourceGroupName = id.ResourceGroupName

			if model := existing.Model; model != nil {
				identity, err := identity.FlattenLegacySystemAndUserAssignedMapToModel(model.Identity)
				if err != nil {
					return fmt.Errorf("flattening Legacy SystemAndUserAssigned Identity: %+v", err)
				}
				loadTest.Identity = identity

				loadTest.Location = location.Normalize(model.Location)
				loadTest.Tags = pointer.From(model.Tags)

				if model.Properties == nil {
					model.Properties = &loadtests.LoadTestProperties{}
				}

				loadTest.DataPlaneURI = pointer.From(model.Properties.DataPlaneURI)
				loadTest.Description = pointer.From(model.Properties.Description)

				if encryption := model.Properties.Encryption; encryption != nil {
					outputEncryption := make([]LoadTestEncryption, 0)
					outputEncryption = append(outputEncryption, LoadTestEncryption{
						KeyURL:   pointer.From(encryption.KeyURL),
						Identity: []LoadTestEncryptionIdentity{},
					})
					loadTest.Encryption = outputEncryption

					if encryptionIdentity := encryption.Identity; encryptionIdentity != nil {
						loadTest.Encryption[0].Identity = append(loadTest.Encryption[0].Identity, LoadTestEncryptionIdentity{
							IdentityID: pointer.From(encryptionIdentity.ResourceId),
						})

						if encryptionIdentity.Type != nil {
							loadTest.Encryption[0].Identity[0].Type = string(pointer.From(encryptionIdentity.Type))
						}
					}
				}
			}

			metadata.SetID(id)

			return metadata.Encode(&loadTest)
		},
	}
}
