// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package dynatrace

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dynatrace/2023-04-27/monitors"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type MonitorsResource struct{}

var _ sdk.ResourceWithUpdate = MonitorsResource{}

type MonitorsResourceModel struct {
	Name                          string                         `tfschema:"name"`
	ResourceGroup                 string                         `tfschema:"resource_group_name"`
	Location                      string                         `tfschema:"location"`
	MonitoringStatus              bool                           `tfschema:"monitoring_enabled"`
	MarketplaceSubscriptionStatus string                         `tfschema:"marketplace_subscription"`
	Identity                      []identity.ModelSystemAssigned `tfschema:"identity"`
	PlanData                      []PlanData                     `tfschema:"plan"`
	UserInfo                      []UserInfo                     `tfschema:"user"`
	Tags                          map[string]string              `tfschema:"tags"`
}

type PlanData struct {
	BillingCycle  string `tfschema:"billing_cycle"`
	EffectiveDate string `tfschema:"effective_date"`
	PlanDetails   string `tfschema:"plan"`
	UsageType     string `tfschema:"usage_type"`
}

type UserInfo struct {
	Country      string `tfschema:"country"`
	EmailAddress string `tfschema:"email"`
	FirstName    string `tfschema:"first_name"`
	LastName     string `tfschema:"last_name"`
	PhoneNumber  string `tfschema:"phone_number"`
}

func (r MonitorsResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"monitoring_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			ForceNew: true,
			Default:  true,
		},

		"marketplace_subscription": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.StringInSlice([]string{
				string(monitors.MarketplaceSubscriptionStatusActive),
				string(monitors.MarketplaceSubscriptionStatusSuspended),
			}, false),
		},

		"identity": commonschema.SystemAssignedIdentityRequired(),

		"plan": {
			Type:     pluginsdk.TypeList,
			Required: true,
			ForceNew: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					// Issue: https://github.com/Azure/azure-rest-api-specs/issues/31284
					// This should be an enum.
					"billing_cycle": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						ValidateFunc: validation.StringInSlice([]string{
							"MONTHLY",
							"WEEKLY",
						}, false),
					},

					"plan": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					// Issue: https://github.com/Azure/azure-rest-api-specs/issues/31284
					// This should be an enum.
					"usage_type": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						ValidateFunc: validation.StringInSlice([]string{
							"PAYG",
							"COMMITTED",
						}, false),
					},

					"effective_date": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
				},
			},
		},

		"user": {
			Type:     pluginsdk.TypeList,
			Required: true,
			ForceNew: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"country": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"email": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validate.EmailAddress,
					},

					"first_name": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"last_name": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"phone_number": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},
			},
		},

		"tags": tags.Schema(),
	}
}

func (r MonitorsResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r MonitorsResource) ModelObject() interface{} {
	return &MonitorsResourceModel{}
}

func (r MonitorsResource) ResourceType() string {
	return "azurerm_dynatrace_monitor"
}

func (r MonitorsResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Dynatrace.MonitorsClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			var model MonitorsResourceModel
			if err := metadata.Decode(&model); err != nil {
				return err
			}

			id := monitors.NewMonitorID(subscriptionId, model.ResourceGroup, model.Name)

			existing, err := client.Get(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			monitoringStatus := monitors.MonitoringStatusEnabled
			if !model.MonitoringStatus {
				monitoringStatus = monitors.MonitoringStatusDisabled
			}
			monitorsProps := monitors.MonitorProperties{
				MarketplaceSubscriptionStatus: pointer.To(monitors.MarketplaceSubscriptionStatus(model.MarketplaceSubscriptionStatus)),
				MonitoringStatus:              pointer.To(monitoringStatus),
				PlanData:                      ExpandDynatracePlanData(model.PlanData),
				UserInfo:                      ExpandDynatraceUserInfo(model.UserInfo),
			}

			dynatraceIdentity, err := expandDynatraceIdentity(model.Identity)
			if err != nil {
				return fmt.Errorf("expanding identity: %+v", err)
			}

			monitor := monitors.MonitorResource{
				Identity:   dynatraceIdentity,
				Location:   model.Location,
				Name:       pointer.To(model.Name),
				Properties: monitorsProps,
				Tags:       pointer.To(model.Tags),
			}

			if err := client.CreateOrUpdateThenPoll(ctx, id, monitor); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)

			return nil
		},
	}
}

func (r MonitorsResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Dynatrace.MonitorsClient
			id, err := monitors.ParseMonitorID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("reading %s: %+v", id, err)
			}

			if model := resp.Model; model != nil {
				props := model.Properties
				identityProps, err := flattenDynatraceIdentity(model.Identity)
				if err != nil {
					return fmt.Errorf("flattening identity: %+v", err)
				}
				monitoringStatus := true
				if *props.MonitoringStatus == monitors.MonitoringStatusDisabled {
					monitoringStatus = false
				}

				state := MonitorsResourceModel{
					Name:                          id.MonitorName,
					ResourceGroup:                 id.ResourceGroupName,
					Location:                      model.Location,
					MonitoringStatus:              monitoringStatus,
					MarketplaceSubscriptionStatus: string(*props.MarketplaceSubscriptionStatus),
					Identity:                      identityProps,
					PlanData:                      FlattenDynatracePlanData(props.PlanData),
					UserInfo:                      FlattenDynatraceUserInfo(props.UserInfo),
				}

				if model.Tags != nil {
					state.Tags = *model.Tags
				}

				return metadata.Encode(&state)
			}
			return nil
		},
	}
}

func (r MonitorsResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Dynatrace.MonitorsClient
			id, err := monitors.ParseMonitorID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			metadata.Logger.Infof("deleting %s", *id)

			if resp, err := client.Delete(ctx, *id); err != nil {
				if !response.WasNotFound(resp.HttpResponse) {
					return fmt.Errorf("deleting %s: %+v", *id, err)
				}
			}
			return nil
		},
	}
}

func (r MonitorsResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return monitors.ValidateMonitorID
}

func (r MonitorsResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Dynatrace.MonitorsClient
			id, err := monitors.ParseMonitorID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var state MonitorsResourceModel
			if err := metadata.Decode(&state); err != nil {
				return fmt.Errorf("decoding %+v", err)
			}

			if metadata.ResourceData.HasChange("tags") {
				props := monitors.MonitorResourceUpdate{
					Tags: pointer.To(state.Tags),
				}

				if _, err := client.Update(ctx, *id, props); err != nil {
					return fmt.Errorf("updating %s: %+v", *id, err)
				}
			}

			return nil
		},
	}
}

func expandDynatraceIdentity(input []identity.ModelSystemAssigned) (*monitors.IdentityProperties, error) {
	config, err := identity.ExpandSystemAssignedFromModel(input)
	if err != nil {
		return nil, err
	}

	if config.Type == identity.TypeNone {
		return pointer.To(monitors.IdentityProperties{}), nil
	}

	dynatraceIdentity := monitors.IdentityProperties{
		Type: monitors.ManagedIdentityType(config.Type),
	}

	return pointer.To(dynatraceIdentity), nil
}

func flattenDynatraceIdentity(input *monitors.IdentityProperties) ([]identity.ModelSystemAssigned, error) {
	if input == nil {
		return nil, errors.New("flattening Dynatrace identity: input is nil")
	}
	var identityProp identity.ModelSystemAssigned

	identityProp.Type = identity.Type(input.Type)

	if input.PrincipalId != nil {
		identityProp.PrincipalId = pointer.From(input.PrincipalId)
	}

	if input.TenantId != nil {
		identityProp.TenantId = pointer.From(input.TenantId)
	}

	return []identity.ModelSystemAssigned{
		identityProp,
	}, nil
}
