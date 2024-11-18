package newrelic

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/newrelic/2024-03-01/monitoredsubscriptions"
	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
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
	SubscriptionId string `tfschema:"subscription_id"`
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

			// The resource is created by the NewRelic Monitor resource, so we check the monitored subscription list to check if it exists
			if !response.WasNotFound(existing.HttpResponse) &&
				existing.Model != nil &&
				existing.Model.Properties != nil &&
				existing.Model.Properties.MonitoredSubscriptionList != nil &&
				len(*existing.Model.Properties.MonitoredSubscriptionList) != 0 {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			properties := &monitoredsubscriptions.MonitoredSubscriptionProperties{
				Properties: &monitoredsubscriptions.SubscriptionList{
					MonitoredSubscriptionList: expandMonitorSubscriptionList(model.MonitoredSubscription),
				},
			}

			if err := client.CreateOrUpdateThenPoll(ctx, *monitorId, *properties); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			if err := resourceMonitoredSubscriptionsWaitForAvailable(ctx, client, id, model.MonitoredSubscription); err != nil {
				return fmt.Errorf("waiting for the %s to become available: %+v", id, err)
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
				return fmt.Errorf("retrieving %s: model was nil", *id)
			}

			if existing.Properties == nil {
				return fmt.Errorf("retrieving %s: property was nil", *id)
			}

			var config NewRelicMonitoredSubscriptionModel
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			if metadata.ResourceData.HasChange("monitored_subscription") {
				existing.Properties = &monitoredsubscriptions.SubscriptionList{
					MonitoredSubscriptionList: expandMonitorSubscriptionList(config.MonitoredSubscription),
				}
			}

			if err := client.CreateOrUpdateThenPoll(ctx, monitorId, *existing); err != nil {
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

			if _, err = client.Delete(ctx, monitorId); err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}

			// The resource cannot be deleted if the parent NewRelic Monitor exists, the DELETE only clears the monitoredSubscriptionList, so here we add a custom poller
			pollerType := &monitoredSubscriptionDeletedPoller{
				client: client,
				id:     *id,
			}
			poller := pollers.NewPoller(pollerType, 5*time.Second, pollers.DefaultNumberOfDroppedConnectionsToAllow)
			if err := poller.PollUntilDone(ctx); err != nil {
				return fmt.Errorf("polling after deleting %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func expandMonitorSubscriptionList(input []NewRelicMonitoredSubscription) *[]monitoredsubscriptions.MonitoredSubscription {
	results := make([]monitoredsubscriptions.MonitoredSubscription, 0)
	if len(input) == 0 {
		return &results
	}

	for _, v := range input {
		results = append(results, monitoredsubscriptions.MonitoredSubscription{
			SubscriptionId: pointer.To(v.SubscriptionId),
		})
	}

	return &results
}

func flattenMonitorSubscriptionList(input *[]monitoredsubscriptions.MonitoredSubscription) []NewRelicMonitoredSubscription {
	if input == nil {
		return make([]NewRelicMonitoredSubscription, 0)
	}

	results := make([]NewRelicMonitoredSubscription, 0)
	for _, v := range *input {
		results = append(results, NewRelicMonitoredSubscription{
			// The returned subscription ID is in upper case
			SubscriptionId: strings.ToLower(pointer.From(v.SubscriptionId)),
		})
	}

	return results
}

var _ pollers.PollerType = &monitoredSubscriptionDeletedPoller{}

type monitoredSubscriptionDeletedPoller struct {
	client *monitoredsubscriptions.MonitoredSubscriptionsClient
	id     parse.NewRelicMonitoredSubscriptionId
}

func (p *monitoredSubscriptionDeletedPoller) Poll(ctx context.Context) (*pollers.PollResult, error) {
	monitorId := monitoredsubscriptions.NewMonitorID(p.id.SubscriptionId, p.id.ResourceGroup, p.id.MonitorName)
	resp, err := p.client.Get(ctx, monitorId)
	if err != nil && !response.WasNotFound(resp.HttpResponse) {
		return nil, fmt.Errorf("retrieving %q: %+v", p.id, err)
	}

	if !response.WasNotFound(resp.HttpResponse) && resp.Model != nil && resp.Model.Properties != nil && resp.Model.Properties.MonitoredSubscriptionList != nil && len(*resp.Model.Properties.MonitoredSubscriptionList) != 0 {
		return &pollers.PollResult{
			HttpResponse: &client.Response{
				Response: resp.HttpResponse,
			},
			Status: pollers.PollingStatusInProgress,
		}, nil
	}

	return &pollers.PollResult{
		HttpResponse: &client.Response{
			Response: resp.HttpResponse,
		},
		Status: pollers.PollingStatusSucceeded,
	}, nil
}

func resourceMonitoredSubscriptionsWaitForAvailable(ctx context.Context, client *monitoredsubscriptions.MonitoredSubscriptionsClient, id parse.NewRelicMonitoredSubscriptionId, monitoredSubscriptions []NewRelicMonitoredSubscription) error {
	deadline, ok := ctx.Deadline()
	if !ok {
		return fmt.Errorf("internal error: context had no deadline")
	}
	state := &pluginsdk.StateChangeConf{
		MinTimeout:                5 * time.Second,
		ContinuousTargetOccurence: 2,
		Pending:                   []string{"Unavailable"},
		Target:                    []string{"Available"},
		Refresh:                   resourceMonitoredSubscriptionsRefresh(ctx, client, id, monitoredSubscriptions),
		Timeout:                   time.Until(deadline),
	}

	if _, err := state.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("waiting for the %s to become available: %+v", id, err)
	}

	return nil
}

func resourceMonitoredSubscriptionsRefresh(ctx context.Context, client *monitoredsubscriptions.MonitoredSubscriptionsClient, id parse.NewRelicMonitoredSubscriptionId, monitoredSubscriptions []NewRelicMonitoredSubscription) pluginsdk.StateRefreshFunc {
	return func() (interface{}, string, error) {
		log.Printf("[DEBUG] Checking to see if %s is available ..", id)

		resp, err := client.Get(ctx, monitoredsubscriptions.NewMonitorID(id.SubscriptionId, id.ResourceGroup, id.MonitorName))
		if err != nil {
			return resp, "Error", fmt.Errorf("retrieving %s: %+v", id, err)
		}

		if resp.Model == nil || resp.Model.Properties == nil || resp.Model.Properties.MonitoredSubscriptionList == nil {
			return resp, "Error", fmt.Errorf("unexpected nil model of %s", id)
		}

		availableCount := 0
		for _, v := range monitoredSubscriptions {
			for _, u := range *resp.Model.Properties.MonitoredSubscriptionList {
				if u.SubscriptionId != nil && strings.EqualFold(v.SubscriptionId, *u.SubscriptionId) {
					availableCount++
					break
				}
			}
		}

		if availableCount != len(monitoredSubscriptions) {
			return resp, "Unavailable", nil
		}

		return resp, "Available", nil
	}
}
