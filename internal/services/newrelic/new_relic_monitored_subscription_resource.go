package newrelic

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/newrelic/2024-03-01/monitoredsubscriptions"
	"github.com/hashicorp/go-azure-sdk/resource-manager/newrelic/2024-03-01/monitors"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/newrelic/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/newrelic/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type NewRelicMonitoredSubscriptionResource struct{}

var (
	_ sdk.Resource           = NewRelicMonitoredSubscriptionResource{}
	_ sdk.ResourceWithUpdate = NewRelicMonitoredSubscriptionResource{}
)

func (r NewRelicMonitoredSubscriptionResource) ResourceType() string {
	return "azurerm_new_relic_monitored_subscription"
}

func (r NewRelicMonitoredSubscriptionResource) ModelObject() interface{} {
	return &NewRelicMonitoredSubscriptionModel{}
}

func (r NewRelicMonitoredSubscriptionResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return validate.NewRelicMonitoredSubscriptionID
}

type NewRelicMonitoredSubscriptionModel struct {
	MonitorId             string                          `tfschema:"monitor_id"`
	MonitoredSubscription []NewRelicMonitoredSubscription `tfschema:"monitored_subscription"`
}

type NewRelicMonitoredSubscription struct {
	SubscriptionId         string                               `tfschema:"subscription_id"`
	AadLogEnabled          bool                                 `tfschema:"azure_active_directory_log_enabled"`
	ActivityLogEnabled     bool                                 `tfschema:"activity_log_enabled"`
	LogTagFilter           []NewRelicMonitoredFilteringTagModel `tfschema:"log_tag_filter"`
	MetricEnabled          bool                                 `tfschema:"metric_enabled"`
	MetricTagFilter        []NewRelicMonitoredFilteringTagModel `tfschema:"metric_tag_filter"`
	SubscriptionLogEnabled bool                                 `tfschema:"subscription_log_enabled"`
}

type NewRelicMonitoredFilteringTagModel struct {
	Action monitoredsubscriptions.TagAction `tfschema:"action"`
	Name   string                           `tfschema:"name"`
	Value  string                           `tfschema:"value"`
}

func (r NewRelicMonitoredSubscriptionResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"monitor_id": commonschema.ResourceIDReferenceRequiredForceNew(&monitoredsubscriptions.MonitorId{}),

		"monitored_subscription": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MinItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"subscription_id": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.IsUUID,
					},

					"azure_active_directory_log_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  false,
					},

					"activity_log_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  false,
					},

					"log_tag_filter": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"name": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},

								"action": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringInSlice(monitoredsubscriptions.PossibleValuesForTagAction(), false),
								},

								"value": {
									Type:     pluginsdk.TypeString,
									Required: true,
									// value can be empty string
								},
							},
						},
					},

					"metric_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  false,
					},

					"metric_tag_filter": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"name": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},

								"action": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringInSlice(monitoredsubscriptions.PossibleValuesForTagAction(), false),
								},

								"value": {
									Type:     pluginsdk.TypeString,
									Required: true,
									// value can be empty string
								},
							},
						},
					},

					"subscription_log_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  false,
					},
				},
			},
		},
	}
}

func (r NewRelicMonitoredSubscriptionResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r NewRelicMonitoredSubscriptionResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model NewRelicMonitoredSubscriptionModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			client := metadata.Client.NewRelic.MonitoredSubscriptionsClient
			monitorId, err := monitoredsubscriptions.ParseMonitorID(model.MonitorId)
			if err != nil {
				return err
			}

			id := parse.NewNewRelicMonitoredSubscriptionID(monitorId.SubscriptionId, monitorId.ResourceGroupName, monitorId.MonitorName, "default")

			existing, err := client.Get(ctx, *monitorId)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for existing %s: %+v", id, err)
			}
			if !response.WasNotFound(existing.HttpResponse) &&
				existing.Model != nil &&
				existing.Model.Properties != nil &&
				existing.Model.Properties.MonitoredSubscriptionList != nil &&
				len(*existing.Model.Properties.MonitoredSubscriptionList) != 0 {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			email, err := r.getEmail(ctx, metadata.Client.NewRelic.MonitorsClient, monitorId)
			if err != nil {
				return err
			}

			properties := &monitoredsubscriptions.MonitoredSubscriptionProperties{
				Properties: &monitoredsubscriptions.SubscriptionList{
					MonitoredSubscriptionList: expandMonitorSubscriptionList(model.MonitoredSubscription, email),
				},
			}

			if err := client.CreateOrUpdateThenPoll(ctx, *monitorId, *properties); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r NewRelicMonitoredSubscriptionResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.NewRelic.MonitoredSubscriptionsClient

			id, err := parse.NewRelicMonitoredSubscriptionID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			monitorId := monitoredsubscriptions.NewMonitorID(id.SubscriptionId, id.ResourceGroup, id.MonitorName)

			resp, err := client.Get(ctx, monitorId)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}

				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			if resp.Model == nil ||
				resp.Model.Properties == nil ||
				resp.Model.Properties.MonitoredSubscriptionList == nil ||
				len(*resp.Model.Properties.MonitoredSubscriptionList) == 0 {
				return metadata.MarkAsGone(id)
			}

			state := NewRelicMonitoredSubscriptionModel{
				MonitorId: monitorId.ID(),
			}

			state.MonitoredSubscription = flattenMonitorSubscriptionList(resp.Model.Properties.MonitoredSubscriptionList)

			return metadata.Encode(&state)
		},
	}
}

func (r NewRelicMonitoredSubscriptionResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.NewRelic.MonitoredSubscriptionsClient

			id, err := parse.NewRelicMonitoredSubscriptionID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			monitorId := monitoredsubscriptions.NewMonitorID(id.SubscriptionId, id.ResourceGroup, id.MonitorName)
			resp, err := client.Get(ctx, monitorId)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			existing := resp.Model
			if existing == nil {
				return fmt.Errorf("retrieving %s: model was nil", id)
			}

			if existing.Properties == nil {
				return fmt.Errorf("retrieving %s: property was nil", id)
			}

			var config NewRelicMonitoredSubscriptionModel
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			email, err := r.getEmail(ctx, metadata.Client.NewRelic.MonitorsClient, &monitorId)
			if err != nil {
				return err
			}

			if metadata.ResourceData.HasChange("monitored_subscription") {
				existing.Properties = &monitoredsubscriptions.SubscriptionList{
					MonitoredSubscriptionList: expandMonitorSubscriptionList(config.MonitoredSubscription, email),
				}
			}

			if err := client.UpdateThenPoll(ctx, monitorId, *existing); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (r NewRelicMonitoredSubscriptionResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.NewRelic.MonitoredSubscriptionsClient
			id, err := parse.NewRelicMonitoredSubscriptionID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			monitorId := monitoredsubscriptions.NewMonitorID(id.SubscriptionId, id.ResourceGroup, id.MonitorName)

			if err = client.DeleteThenPoll(ctx, monitorId); err != nil {
				return fmt.Errorf("deleting %s: %+v", id, err)
			}

			return nil
		},
	}
}

func expandMonitorSubscriptionList(input []NewRelicMonitoredSubscription, email string) *[]monitoredsubscriptions.MonitoredSubscription {
	results := make([]monitoredsubscriptions.MonitoredSubscription, 0)
	if len(input) == 0 {
		return &results
	}

	for _, v := range input {
		result := monitoredsubscriptions.MonitoredSubscription{
			SubscriptionId: pointer.To(v.SubscriptionId),
		}

		if len(v.LogTagFilter) > 0 {
			logRules := monitoredsubscriptions.LogRules{
				FilteringTags:        expandMonitoredSubscriptionFilteringTagModelArray(v.LogTagFilter),
				SendAadLogs:          pointer.To(monitoredsubscriptions.SendAadLogsStatusDisabled),
				SendActivityLogs:     pointer.To(monitoredsubscriptions.SendActivityLogsStatusDisabled),
				SendSubscriptionLogs: pointer.To(monitoredsubscriptions.SendSubscriptionLogsStatusDisabled),
			}

			if v.AadLogEnabled {
				logRules.SendAadLogs = pointer.To(monitoredsubscriptions.SendAadLogsStatusEnabled)
			}

			if v.ActivityLogEnabled {
				logRules.SendActivityLogs = pointer.To(monitoredsubscriptions.SendActivityLogsStatusEnabled)
			}

			if v.SubscriptionLogEnabled {
				logRules.SendSubscriptionLogs = pointer.To(monitoredsubscriptions.SendSubscriptionLogsStatusEnabled)
			}

			if result.TagRules == nil {
				result.TagRules = &monitoredsubscriptions.MonitoringTagRulesProperties{}
			}
			result.TagRules.LogRules = pointer.To(logRules)
		}

		if len(v.MetricTagFilter) > 0 {
			metricRules := monitoredsubscriptions.MetricRules{
				FilteringTags: expandMonitoredSubscriptionFilteringTagModelArray(v.MetricTagFilter),
				SendMetrics:   pointer.To(monitoredsubscriptions.SendMetricsStatusDisabled),
				UserEmail:     pointer.To(email),
			}

			if v.MetricEnabled {
				metricRules.SendMetrics = pointer.To(monitoredsubscriptions.SendMetricsStatusEnabled)
			}

			if result.TagRules == nil {
				result.TagRules = &monitoredsubscriptions.MonitoringTagRulesProperties{}
			}
			result.TagRules.MetricRules = pointer.To(metricRules)
		}

		results = append(results, result)
	}

	return &results
}

func flattenMonitorSubscriptionList(input *[]monitoredsubscriptions.MonitoredSubscription) []NewRelicMonitoredSubscription {
	if input == nil {
		return make([]NewRelicMonitoredSubscription, 0)
	}

	results := make([]NewRelicMonitoredSubscription, 0)
	for _, v := range *input {
		result := NewRelicMonitoredSubscription{
			SubscriptionId: pointer.From(v.SubscriptionId),
		}

		if tagRule := v.TagRules; tagRule != nil {
			if logRule := tagRule.LogRules; logRule != nil {
				result.AadLogEnabled = logRule.SendAadLogs != nil && *logRule.SendAadLogs == monitoredsubscriptions.SendAadLogsStatusEnabled
				result.ActivityLogEnabled = logRule.SendActivityLogs != nil && *logRule.SendActivityLogs == monitoredsubscriptions.SendActivityLogsStatusEnabled
				result.LogTagFilter = flattenMonitoredSubscriptionFilteringTagModelArray(logRule.FilteringTags)
				result.SubscriptionLogEnabled = logRule.SendSubscriptionLogs != nil && *logRule.SendSubscriptionLogs == monitoredsubscriptions.SendSubscriptionLogsStatusEnabled
			}

			if metricRule := tagRule.MetricRules; metricRule != nil {
				result.MetricEnabled = metricRule.SendMetrics != nil && *metricRule.SendMetrics == monitoredsubscriptions.SendMetricsStatusEnabled
				result.MetricTagFilter = flattenMonitoredSubscriptionFilteringTagModelArray(metricRule.FilteringTags)
			}
		}
	}

	return results
}

func (r NewRelicMonitoredSubscriptionResource) getEmail(ctx context.Context, monitorClient *monitors.MonitorsClient, monitorId *monitoredsubscriptions.MonitorId) (string, error) {
	id := monitors.NewMonitorID(monitorId.SubscriptionId, monitorId.ResourceGroupName, monitorId.MonitorName)
	monitor, err := monitorClient.Get(ctx, id)
	if err != nil {
		return "", fmt.Errorf("getting monitor: %+v", err)
	}

	if monitor.Model == nil || monitor.Model.Properties.UserInfo == nil || monitor.Model.Properties.UserInfo.EmailAddress == nil || *monitor.Model.Properties.UserInfo.EmailAddress == "" {
		return "", fmt.Errorf("failed to get user email address from monitor")
	}

	return *monitor.Model.Properties.UserInfo.EmailAddress, nil
}

func expandMonitoredSubscriptionFilteringTagModelArray(inputList []NewRelicMonitoredFilteringTagModel) *[]monitoredsubscriptions.FilteringTag {
	var outputList []monitoredsubscriptions.FilteringTag
	for _, v := range inputList {
		input := v
		output := monitoredsubscriptions.FilteringTag{
			Action: pointer.To(input.Action),
		}

		if input.Name != "" {
			output.Name = pointer.To(input.Name)
		}

		if input.Value != "" {
			output.Value = pointer.To(input.Value)
		}
		outputList = append(outputList, output)
	}
	return &outputList
}

func flattenMonitoredSubscriptionFilteringTagModelArray(inputList *[]monitoredsubscriptions.FilteringTag) []NewRelicMonitoredFilteringTagModel {
	outputList := make([]NewRelicMonitoredFilteringTagModel, 0)
	if inputList == nil {
		return outputList
	}

	for _, input := range *inputList {
		outputList = append(outputList, NewRelicMonitoredFilteringTagModel{
			Action: pointer.From(input.Action),
			Name:   pointer.From(input.Name),
			Value:  pointer.From(input.Value),
		})
	}

	return outputList
}
