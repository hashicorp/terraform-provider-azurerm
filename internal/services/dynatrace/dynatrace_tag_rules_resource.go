package dynatrace

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dynatrace/2021-09-01/monitors"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dynatrace/2021-09-01/tagrules"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type TagRulesResource struct{}

type TagRulesResourceModel struct {
	Name        string       `tfschema:"name"`
	Monitor     string       `tfschema:"monitor_id"`
	LogRules    []LogRule    `tfschema:"log_rules"`
	MetricRules []MetricRule `tfschema:"metric_rules"`
}

type MetricRule struct {
	FilteringTags []FilteringTag `tfschema:"filtering_tag"`
}

type LogRule struct {
	FilteringTags        []FilteringTag `tfschema:"filtering_tag"`
	SendAadLogs          string         `tfschema:"send_aad_logs_enabled"`
	SendActivityLogs     string         `tfschema:"send_activity_logs_enabled"`
	SendSubscriptionLogs string         `tfschema:"send_subscription_logs_enabled"`
}

type FilteringTag struct {
	Name   string `tfschema:"name"`
	Value  string `tfschema:"value"`
	Action string `tfschema:"action"`
}

func (r TagRulesResource) Arguments() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"monitor_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: monitors.ValidateMonitorID,
		},

		"log_rules": SchemaLogRule(),

		"metric_rules": SchemaMetricRules(),
	}
}

func (r TagRulesResource) Attributes() map[string]*schema.Schema {
	return map[string]*schema.Schema{}
}

func (r TagRulesResource) ModelObject() interface{} {
	return &TagRulesResourceModel{}
}

func (r TagRulesResource) ResourceType() string {
	return "azurerm_dynatrace_tag_rules"
}

func (r TagRulesResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model TagRulesResourceModel
			if err := metadata.Decode(&model); err != nil {
				return err
			}

			client := metadata.Client.Dynatrace.TagRulesClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			monitorsId, err := monitors.ParseMonitorID(model.Monitor)
			id := tagrules.NewTagRuleID(subscriptionId, monitorsId.ResourceGroupName, monitorsId.MonitorName, model.Name)
			if err != nil {
				return err
			}

			existing, err := client.Get(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			tagRulesProps := tagrules.MonitoringTagRulesProperties{
				LogRules:    ExpandLogRule(model.LogRules),
				MetricRules: ExpandMetricRules(model.MetricRules),
			}
			tagRules := tagrules.TagRule{
				Name:       &model.Name,
				Properties: tagRulesProps,
			}

			if _, err := client.CreateOrUpdate(ctx, id, tagRules); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)

			return nil
		},
	}
}

func (r TagRulesResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Dynatrace.TagRulesClient
			id, err := tagrules.ParseTagRuleID(metadata.ResourceData.Id())
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
				monitorId := monitors.NewMonitorID(id.SubscriptionId, id.ResourceGroupName, id.MonitorName)

				state := TagRulesResourceModel{
					Name:        id.RuleSetName,
					Monitor:     monitorId.ID(),
					LogRules:    FlattenLogRules(props.LogRules),
					MetricRules: FlattenMetricRules(props.MetricRules),
				}

				return metadata.Encode(&state)
			}

			return nil
		},
	}
}

func (r TagRulesResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Dynatrace.TagRulesClient
			id, err := tagrules.ParseTagRuleID(metadata.ResourceData.Id())
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

func (r TagRulesResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return tagrules.ValidateTagRuleID
}
