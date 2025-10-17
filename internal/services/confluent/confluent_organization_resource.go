// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package confluent

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/confluent/2024-07-01/organizationresources"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/confluent/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type ConfluentOrganizationResource struct{}

type ConfluentOrganizationResourceModel struct {
	Name              string                       `tfschema:"name"`
	ResourceGroupName string                       `tfschema:"resource_group_name"`
	Location          string                       `tfschema:"location"`
	OfferDetail       []ConfluentOfferDetailModel  `tfschema:"offer_detail"`
	UserDetail        []ConfluentUserDetailModel   `tfschema:"user_detail"`
	LinkOrganization  []ConfluentLinkOrgModel      `tfschema:"link_organization"`
	Tags              map[string]string            `tfschema:"tags"`

	OrganizationId    string `tfschema:"organization_id"`
	SsoUrl            string `tfschema:"sso_url"`
	CreatedTime       string `tfschema:"created_time"`
	ProvisioningState string `tfschema:"provisioning_state"`
}

type ConfluentOfferDetailModel struct {
	Id              string   `tfschema:"id"`
	PlanId          string   `tfschema:"plan_id"`
	PlanName        string   `tfschema:"plan_name"`
	PublisherId     string   `tfschema:"publisher_id"`
	TermUnit        string   `tfschema:"term_unit"`
	PrivateOfferId  string   `tfschema:"private_offer_id"`
	PrivateOfferIds []string `tfschema:"private_offer_ids"`
	TermId          string   `tfschema:"term_id"`
	Status          string   `tfschema:"status"`
}

type ConfluentUserDetailModel struct {
	EmailAddress      string `tfschema:"email_address"`
	FirstName         string `tfschema:"first_name"`
	LastName          string `tfschema:"last_name"`
	AadEmail          string `tfschema:"aad_email"`
	UserPrincipalName string `tfschema:"user_principal_name"`
}

type ConfluentLinkOrgModel struct {
	Token string `tfschema:"token"`
}

func (r ConfluentOrganizationResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"offer_detail": {
			Type:     pluginsdk.TypeList,
			Required: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"id": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"plan_id": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"plan_name": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"publisher_id": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"term_unit": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"private_offer_id": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"private_offer_ids": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: validation.StringIsNotEmpty,
						},
					},

					"term_id": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"status": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
				},
			},
		},

		"user_detail": {
			Type:     pluginsdk.TypeList,
			Required: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"email_address": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"first_name": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"last_name": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"aad_email": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"user_principal_name": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
				},
			},
		},

		"link_organization": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"token": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						Sensitive:    true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},
			},
		},

		"tags": commonschema.Tags(),
	}
}

func (r ConfluentOrganizationResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"organization_id": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"sso_url": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"created_time": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"provisioning_state": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}

func (r ConfluentOrganizationResource) ModelObject() interface{} {
	return &ConfluentOrganizationResourceModel{}
}

func (r ConfluentOrganizationResource) ResourceType() string {
	return "azurerm_confluent_organization"
}

func (r ConfluentOrganizationResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Confluent.OrganizationResourcesClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			var model ConfluentOrganizationResourceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			id := organizationresources.NewOrganizationID(subscriptionId, model.ResourceGroupName, model.Name)

			existing, err := client.OrganizationGet(ctx, id)
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for existing %s: %+v", id, err)
				}
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return tf.ImportAsExistsError("azurerm_confluent_organization", id.ID())
			}

			payload := organizationresources.OrganizationResource{
				Location:   location.Normalize(model.Location),
				Properties: expandConfluentOrganizationProperties(model),
				Tags:       &model.Tags,
			}

			if err := client.OrganizationCreateThenPoll(ctx, id, payload); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r ConfluentOrganizationResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Confluent.OrganizationResourcesClient

			id, err := organizationresources.ParseOrganizationID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.OrganizationGet(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			var state ConfluentOrganizationResourceModel
			state.Name = id.OrganizationName
			state.ResourceGroupName = id.ResourceGroupName

			if model := resp.Model; model != nil {
				state.Location = location.Normalize(model.Location)
				state.Tags = pointer.From(model.Tags)

				props := model.Properties
				state.OrganizationId = pointer.From(props.OrganizationId)
				state.SsoUrl = pointer.From(props.SsoURL)
				state.CreatedTime = pointer.From(props.CreatedTime)
				state.ProvisioningState = string(pointer.From(props.ProvisioningState))
				state.OfferDetail = flattenConfluentOfferDetail(&props.OfferDetail)
				state.UserDetail = flattenConfluentUserDetail(&props.UserDetail)
				state.LinkOrganization = flattenConfluentLinkOrganization(props.LinkOrganization)
			}

			return metadata.Encode(&state)
		},
	}
}

func (r ConfluentOrganizationResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Confluent.OrganizationResourcesClient

			id, err := organizationresources.ParseOrganizationID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model ConfluentOrganizationResourceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			// The Azure API only supports updating tags via PATCH
			// offer_detail, user_detail, and link_organization cannot be updated after creation
			payload := organizationresources.OrganizationResourceUpdate{
				Tags: &model.Tags,
			}

			if _, err := client.OrganizationUpdate(ctx, *id, payload); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (r ConfluentOrganizationResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Confluent.OrganizationResourcesClient

			id, err := organizationresources.ParseOrganizationID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if err := client.OrganizationDeleteThenPoll(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (r ConfluentOrganizationResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return validate.OrganizationID
}

var _ sdk.ResourceWithUpdate = ConfluentOrganizationResource{}

func expandConfluentOrganizationProperties(model ConfluentOrganizationResourceModel) organizationresources.OrganizationResourceProperties {
	props := organizationresources.OrganizationResourceProperties{
		OfferDetail: expandConfluentOfferDetail(model.OfferDetail),
		UserDetail:  expandConfluentUserDetail(model.UserDetail),
	}

	if len(model.LinkOrganization) > 0 {
		props.LinkOrganization = expandConfluentLinkOrganization(model.LinkOrganization)
	}

	return props
}

func expandConfluentOfferDetail(input []ConfluentOfferDetailModel) organizationresources.OfferDetail {
	if len(input) == 0 {
		return organizationresources.OfferDetail{}
	}

	v := input[0]

	result := organizationresources.OfferDetail{
		Id:          v.Id,
		PlanId:      v.PlanId,
		PlanName:    v.PlanName,
		PublisherId: v.PublisherId,
		TermUnit:    v.TermUnit,
	}

	if v.PrivateOfferId != "" {
		result.PrivateOfferId = pointer.To(v.PrivateOfferId)
	}

	if len(v.PrivateOfferIds) > 0 {
		result.PrivateOfferIds = &v.PrivateOfferIds
	}

	if v.TermId != "" {
		result.TermId = pointer.To(v.TermId)
	}

	return result
}

func expandConfluentUserDetail(input []ConfluentUserDetailModel) organizationresources.UserDetail {
	if len(input) == 0 {
		return organizationresources.UserDetail{}
	}

	v := input[0]

	result := organizationresources.UserDetail{
		EmailAddress: v.EmailAddress,
	}

	if v.FirstName != "" {
		result.FirstName = pointer.To(v.FirstName)
	}

	if v.LastName != "" {
		result.LastName = pointer.To(v.LastName)
	}

	return result
}

func expandConfluentLinkOrganization(input []ConfluentLinkOrgModel) *organizationresources.LinkOrganization {
	if len(input) == 0 {
		return nil
	}

	v := input[0]

	return &organizationresources.LinkOrganization{
		Token: v.Token,
	}
}

func flattenConfluentOfferDetail(input *organizationresources.OfferDetail) []ConfluentOfferDetailModel {
	if input == nil {
		return []ConfluentOfferDetailModel{}
	}

	privateOfferIds := make([]string, 0)
	if input.PrivateOfferIds != nil {
		privateOfferIds = *input.PrivateOfferIds
	}

	return []ConfluentOfferDetailModel{
		{
			Id:              input.Id,
			PlanId:          input.PlanId,
			PlanName:        input.PlanName,
			PublisherId:     input.PublisherId,
			TermUnit:        input.TermUnit,
			PrivateOfferId:  pointer.From(input.PrivateOfferId),
			PrivateOfferIds: privateOfferIds,
			TermId:          pointer.From(input.TermId),
			Status:          string(pointer.From(input.Status)),
		},
	}
}

func flattenConfluentUserDetail(input *organizationresources.UserDetail) []ConfluentUserDetailModel {
	if input == nil {
		return []ConfluentUserDetailModel{}
	}

	return []ConfluentUserDetailModel{
		{
			EmailAddress:      input.EmailAddress,
			FirstName:         pointer.From(input.FirstName),
			LastName:          pointer.From(input.LastName),
			AadEmail:          pointer.From(input.AadEmail),
			UserPrincipalName: pointer.From(input.UserPrincipalName),
		},
	}
}

func flattenConfluentLinkOrganization(input *organizationresources.LinkOrganization) []ConfluentLinkOrgModel {
	if input == nil {
		return []ConfluentLinkOrgModel{}
	}

	// Don't expose the token in state for security reasons
	return []ConfluentLinkOrgModel{
		{
			Token: "",
		},
	}
}
