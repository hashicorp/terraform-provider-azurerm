package dynatrace

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dynatrace/2021-09-01/monitors"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type MonitorsResource struct{}

var _ sdk.ResourceWithUpdate = MonitorsResource{}

type MonitorsResourceModel struct {
	Name                          string            `tfschema:"name"`
	ResourceGroup                 string            `tfschema:"resource_group_name"`
	Location                      string            `tfschema:"location"`
	MonitoringStatus              bool              `tfschema:"monitoring_enabled"`
	MarketplaceSubscriptionStatus string            `tfschema:"marketplace_subscription"`
	IdentityType                  string            `tfschema:"identity_type"`
	PlanData                      []PlanData        `tfschema:"plan_data"`
	UserInfo                      []UserInfo        `tfschema:"user_info"`
	Tags                          map[string]string `tfschema:"tags"`
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

func (r MonitorsResource) Arguments() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"resource_group_name": azure.SchemaResourceGroupName(),

		"location": azure.SchemaLocation(),

		"identity_type": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			ForceNew: true,
			Default:  "SystemAssigned",
			ValidateFunc: validation.StringInSlice([]string{
				"SystemAssigned",
				"UserAssigned",
			}, false),
		},

		"monitoring_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			ForceNew: true,
			Default:  true,
		},

		"marketplace_subscription": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			ForceNew: true,
			ValidateFunc: validation.StringInSlice([]string{
				"Active",
				"Suspended",
			}, false),
		},

		"plan_data": SchemaPlanData(),

		"user_info": SchemaUserInfo(),

		"tags": tags.Schema(),
	}
}

func (r MonitorsResource) Attributes() map[string]*schema.Schema {
	return map[string]*schema.Schema{}
}

func (r MonitorsResource) ModelObject() interface{} {
	return &MonitorsResourceModel{}
}

func (r MonitorsResource) ResourceType() string {
	return "azurerm_dynatrace_monitors"
}

func (r MonitorsResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model MonitorsResourceModel
			if err := metadata.Decode(&model); err != nil {
				return err
			}

			client := metadata.Client.Dynatrace.MonitorClient
			subscriptionId := metadata.Client.Account.SubscriptionId
			id := monitors.NewMonitorID(subscriptionId, model.ResourceGroup, model.Name)

			existing, err := client.Get(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			marketplaceSubscriptionServiceStatus := monitors.MarketplaceSubscriptionStatus(model.MarketplaceSubscriptionStatus)
			monitoringStatus := monitors.MonitoringStatusEnabled
			if !model.MonitoringStatus {
				monitoringStatus = monitors.MonitoringStatusDisabled
			}
			monitorsProps := monitors.MonitorProperties{
				MarketplaceSubscriptionStatus: &marketplaceSubscriptionServiceStatus,
				MonitoringStatus:              &monitoringStatus,
				PlanData:                      ExpandDynatracePlanData(model.PlanData),
				UserInfo:                      ExpandDynatraceUserInfo(model.UserInfo),
			}

			identity := monitors.IdentityProperties{
				Type: monitors.ManagedIdentityType(model.IdentityType),
			}

			monitor := monitors.MonitorResource{
				Identity:   &identity,
				Location:   model.Location,
				Name:       &model.Name,
				Properties: monitorsProps,
				Tags:       &model.Tags,
			}

			if _, err := client.CreateOrUpdate(ctx, id, monitor); err != nil {
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
			client := metadata.Client.Dynatrace.MonitorClient
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
				identityProps := model.Identity
				userInfo := metadata.ResourceData.Get("user_info").([]interface{})
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
					IdentityType:                  string(identityProps.Type),
					PlanData:                      FlattenDynatracePlanData(props.PlanData),
					UserInfo:                      FlattenDynatraceUserInfo(userInfo),
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
			client := metadata.Client.Dynatrace.MonitorClient
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
			client := metadata.Client.Dynatrace.MonitorClient
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
					Tags: &state.Tags,
				}

				if _, err := client.Update(ctx, *id, props); err != nil {
					return fmt.Errorf("updating %s: %+v", *id, err)
				}
			}

			return nil
		},
	}
}
