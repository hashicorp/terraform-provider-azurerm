// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package confluent

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/confluent/2024-07-01/scclusterrecords"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/confluent/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type ConfluentClusterResource struct{}

type ConfluentClusterResourceModel struct {
	ClusterId         string                        `tfschema:"cluster_id"`
	EnvironmentId     string                        `tfschema:"environment_id"`
	OrganizationId    string                        `tfschema:"organization_id"`
	ResourceGroupName string                        `tfschema:"resource_group_name"`
	DisplayName       string                        `tfschema:"display_name"`
	Availability      string                        `tfschema:"availability"`
	Cloud             string                        `tfschema:"cloud"`
	Region            string                        `tfschema:"region"`
	Package           string                        `tfschema:"package"`
	Spec              []ConfluentClusterSpecModel   `tfschema:"spec"`

	// Computed
	Id                      string                         `tfschema:"id"`
	Kind                    string                         `tfschema:"kind"`
	ApiEndpoint             string                         `tfschema:"api_endpoint"`
	HttpEndpoint            string                         `tfschema:"http_endpoint"`
	KafkaBootstrapEndpoint  string                         `tfschema:"kafka_bootstrap_endpoint"`
	Metadata                []ConfluentClusterMetadataModel `tfschema:"metadata"`
	Status                  []ConfluentClusterStatusModel   `tfschema:"status"`
}

type ConfluentClusterSpecModel struct {
	Zone            string                                      `tfschema:"zone"`
	Config          []ConfluentClusterConfigModel               `tfschema:"config"`
	Environment     []ConfluentClusterNetworkEnvironmentModel   `tfschema:"environment"`
	Network         []ConfluentClusterNetworkEnvironmentModel   `tfschema:"network"`
	Byok            []ConfluentClusterByokModel                 `tfschema:"byok"`
}

type ConfluentClusterConfigModel struct {
	Kind string `tfschema:"kind"`
}

type ConfluentClusterNetworkEnvironmentModel struct {
	Id           string `tfschema:"id"`
	Environment  string `tfschema:"environment"`
	Related      string `tfschema:"related"`
	ResourceName string `tfschema:"resource_name"`
}

type ConfluentClusterByokModel struct {
	Id           string `tfschema:"id"`
	Related      string `tfschema:"related"`
	ResourceName string `tfschema:"resource_name"`
}

type ConfluentClusterMetadataModel struct {
	Self             string `tfschema:"self"`
	ResourceName     string `tfschema:"resource_name"`
	CreatedTimestamp string `tfschema:"created_timestamp"`
	UpdatedTimestamp string `tfschema:"updated_timestamp"`
	DeletedTimestamp string `tfschema:"deleted_timestamp"`
}

type ConfluentClusterStatusModel struct {
	Phase string `tfschema:"phase"`
	Cku   int64  `tfschema:"cku"`
}

func (r ConfluentClusterResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"cluster_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"environment_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"organization_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"resource_group_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"display_name": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"availability": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"cloud": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"region": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"package": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringInSlice(scclusterrecords.PossibleValuesForPackage(), false),
		},

		"spec": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			ForceNew: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"zone": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"config": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"kind": {
									Type:         pluginsdk.TypeString,
									Optional:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},
							},
						},
					},

					"environment": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"id": {
									Type:         pluginsdk.TypeString,
									Optional:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},

								"environment": {
									Type:         pluginsdk.TypeString,
									Optional:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},

								"related": {
									Type:     pluginsdk.TypeString,
									Computed: true,
								},

								"resource_name": {
									Type:     pluginsdk.TypeString,
									Computed: true,
								},
							},
						},
					},

					"network": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"id": {
									Type:         pluginsdk.TypeString,
									Optional:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},

								"environment": {
									Type:         pluginsdk.TypeString,
									Optional:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},

								"related": {
									Type:     pluginsdk.TypeString,
									Computed: true,
								},

								"resource_name": {
									Type:     pluginsdk.TypeString,
									Computed: true,
								},
							},
						},
					},

					"byok": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"id": {
									Type:         pluginsdk.TypeString,
									Optional:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},

								"related": {
									Type:     pluginsdk.TypeString,
									Computed: true,
								},

								"resource_name": {
									Type:     pluginsdk.TypeString,
									Computed: true,
								},
							},
						},
					},
				},
			},
		},
	}
}

func (r ConfluentClusterResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"id": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"kind": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"api_endpoint": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"http_endpoint": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"kafka_bootstrap_endpoint": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"metadata": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"self": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"resource_name": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"created_timestamp": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"updated_timestamp": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"deleted_timestamp": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
				},
			},
		},

		"status": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"phase": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"cku": {
						Type:     pluginsdk.TypeInt,
						Computed: true,
					},
				},
			},
		},
	}
}

func (r ConfluentClusterResource) ModelObject() interface{} {
	return &ConfluentClusterResourceModel{}
}

func (r ConfluentClusterResource) ResourceType() string {
	return "azurerm_confluent_cluster"
}

func (r ConfluentClusterResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Confluent.ClusterClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			var model ConfluentClusterResourceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			id := scclusterrecords.NewClusterID(subscriptionId, model.ResourceGroupName, model.OrganizationId, model.EnvironmentId, model.ClusterId)

			existing, err := client.OrganizationGetClusterById(ctx, id)
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for existing %s: %+v", id, err)
				}
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return tf.ImportAsExistsError("azurerm_confluent_cluster", id.ID())
			}

			payload := scclusterrecords.SCClusterRecord{
				Kind:       pointer.To("Cluster"),
				Properties: expandConfluentClusterProperties(model),
			}

			if model.DisplayName != "" {
				payload.Name = pointer.To(model.DisplayName)
			}

			if _, err := client.ClusterCreateOrUpdate(ctx, id, payload); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r ConfluentClusterResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Confluent.ClusterClient

			id, err := scclusterrecords.ParseClusterID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.OrganizationGetClusterById(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			var state ConfluentClusterResourceModel
			state.ClusterId = id.ClusterId
			state.EnvironmentId = id.EnvironmentId
			state.OrganizationId = id.OrganizationName
			state.ResourceGroupName = id.ResourceGroupName

			if model := resp.Model; model != nil {
				state.Id = pointer.From(model.Id)
				state.Kind = pointer.From(model.Kind)
				state.DisplayName = pointer.From(model.Name)

				if props := model.Properties; props != nil {
					state.Metadata = flattenConfluentClusterMetadata(props.Metadata)
					state.Status = flattenConfluentClusterStatus(props.Status)

					if spec := props.Spec; spec != nil {
						state.ApiEndpoint = pointer.From(spec.ApiEndpoint)
						state.HttpEndpoint = pointer.From(spec.HTTPEndpoint)
						state.KafkaBootstrapEndpoint = pointer.From(spec.KafkaBootstrapEndpoint)
						state.Availability = pointer.From(spec.Availability)
						state.Cloud = pointer.From(spec.Cloud)
						state.Region = pointer.From(spec.Region)
						state.Package = string(pointer.From(spec.Package))
						state.Spec = flattenConfluentClusterSpec(spec)
					}
				}
			}

			return metadata.Encode(&state)
		},
	}
}

func (r ConfluentClusterResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Confluent.ClusterClient

			id, err := scclusterrecords.ParseClusterID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if _, err := client.ClusterDelete(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (r ConfluentClusterResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return validate.ClusterID
}

func expandConfluentClusterProperties(model ConfluentClusterResourceModel) *scclusterrecords.ClusterProperties {
	props := &scclusterrecords.ClusterProperties{}

	spec := &scclusterrecords.SCClusterSpecEntity{}

	if model.DisplayName != "" {
		spec.Name = pointer.To(model.DisplayName)
	}

	if model.Availability != "" {
		spec.Availability = pointer.To(model.Availability)
	}

	if model.Cloud != "" {
		spec.Cloud = pointer.To(model.Cloud)
	}

	if model.Region != "" {
		spec.Region = pointer.To(model.Region)
	}

	if model.Package != "" {
		spec.Package = pointer.To(scclusterrecords.Package(model.Package))
	}

	if len(model.Spec) > 0 {
		specModel := model.Spec[0]

		if specModel.Zone != "" {
			spec.Zone = pointer.To(specModel.Zone)
		}

		if len(specModel.Config) > 0 {
			spec.Config = &scclusterrecords.ClusterConfigEntity{
				Kind: pointer.To(specModel.Config[0].Kind),
			}
		}

		if len(specModel.Environment) > 0 {
			env := specModel.Environment[0]
			spec.Environment = &scclusterrecords.SCClusterNetworkEnvironmentEntity{
				Id:          pointer.To(env.Id),
				Environment: pointer.To(env.Environment),
			}
		}

		if len(specModel.Network) > 0 {
			net := specModel.Network[0]
			spec.Network = &scclusterrecords.SCClusterNetworkEnvironmentEntity{
				Id:          pointer.To(net.Id),
				Environment: pointer.To(net.Environment),
			}
		}

		if len(specModel.Byok) > 0 {
			byok := specModel.Byok[0]
			spec.Byok = &scclusterrecords.SCClusterByokEntity{
				Id: pointer.To(byok.Id),
			}
		}
	}

	props.Spec = spec

	return props
}

func flattenConfluentClusterMetadata(input *scclusterrecords.SCMetadataEntity) []ConfluentClusterMetadataModel {
	if input == nil {
		return []ConfluentClusterMetadataModel{}
	}

	return []ConfluentClusterMetadataModel{
		{
			Self:             pointer.From(input.Self),
			ResourceName:     pointer.From(input.ResourceName),
			CreatedTimestamp: pointer.From(input.CreatedTimestamp),
			UpdatedTimestamp: pointer.From(input.UpdatedTimestamp),
			DeletedTimestamp: pointer.From(input.DeletedTimestamp),
		},
	}
}

func flattenConfluentClusterStatus(input *scclusterrecords.ClusterStatusEntity) []ConfluentClusterStatusModel {
	if input == nil {
		return []ConfluentClusterStatusModel{}
	}

	return []ConfluentClusterStatusModel{
		{
			Phase: pointer.From(input.Phase),
			Cku:   pointer.From(input.Cku),
		},
	}
}

func flattenConfluentClusterSpec(input *scclusterrecords.SCClusterSpecEntity) []ConfluentClusterSpecModel {
	if input == nil {
		return []ConfluentClusterSpecModel{}
	}

	result := ConfluentClusterSpecModel{
		Zone: pointer.From(input.Zone),
	}

	if input.Config != nil {
		result.Config = []ConfluentClusterConfigModel{
			{
				Kind: pointer.From(input.Config.Kind),
			},
		}
	}

	if input.Environment != nil {
		result.Environment = []ConfluentClusterNetworkEnvironmentModel{
			{
				Id:           pointer.From(input.Environment.Id),
				Environment:  pointer.From(input.Environment.Environment),
				Related:      pointer.From(input.Environment.Related),
				ResourceName: pointer.From(input.Environment.ResourceName),
			},
		}
	}

	if input.Network != nil {
		result.Network = []ConfluentClusterNetworkEnvironmentModel{
			{
				Id:           pointer.From(input.Network.Id),
				Environment:  pointer.From(input.Network.Environment),
				Related:      pointer.From(input.Network.Related),
				ResourceName: pointer.From(input.Network.ResourceName),
			},
		}
	}

	if input.Byok != nil {
		result.Byok = []ConfluentClusterByokModel{
			{
				Id:           pointer.From(input.Byok.Id),
				Related:      pointer.From(input.Byok.Related),
				ResourceName: pointer.From(input.Byok.ResourceName),
			},
		}
	}

	return []ConfluentClusterSpecModel{result}
}
