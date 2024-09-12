package alertsmanagements

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AlertsGetSummaryOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *AlertsSummary
}

type AlertsGetSummaryOperationOptions struct {
	AlertRule               *string
	AlertState              *AlertState
	CustomTimeRange         *string
	Groupby                 *AlertsSummaryGroupByFields
	IncludeSmartGroupsCount *bool
	MonitorCondition        *MonitorCondition
	MonitorService          *MonitorService
	Severity                *Severity
	TargetResource          *string
	TargetResourceGroup     *string
	TargetResourceType      *string
	TimeRange               *TimeRange
}

func DefaultAlertsGetSummaryOperationOptions() AlertsGetSummaryOperationOptions {
	return AlertsGetSummaryOperationOptions{}
}

func (o AlertsGetSummaryOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o AlertsGetSummaryOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}
	return &out
}

func (o AlertsGetSummaryOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.AlertRule != nil {
		out.Append("alertRule", fmt.Sprintf("%v", *o.AlertRule))
	}
	if o.AlertState != nil {
		out.Append("alertState", fmt.Sprintf("%v", *o.AlertState))
	}
	if o.CustomTimeRange != nil {
		out.Append("customTimeRange", fmt.Sprintf("%v", *o.CustomTimeRange))
	}
	if o.Groupby != nil {
		out.Append("groupby", fmt.Sprintf("%v", *o.Groupby))
	}
	if o.IncludeSmartGroupsCount != nil {
		out.Append("includeSmartGroupsCount", fmt.Sprintf("%v", *o.IncludeSmartGroupsCount))
	}
	if o.MonitorCondition != nil {
		out.Append("monitorCondition", fmt.Sprintf("%v", *o.MonitorCondition))
	}
	if o.MonitorService != nil {
		out.Append("monitorService", fmt.Sprintf("%v", *o.MonitorService))
	}
	if o.Severity != nil {
		out.Append("severity", fmt.Sprintf("%v", *o.Severity))
	}
	if o.TargetResource != nil {
		out.Append("targetResource", fmt.Sprintf("%v", *o.TargetResource))
	}
	if o.TargetResourceGroup != nil {
		out.Append("targetResourceGroup", fmt.Sprintf("%v", *o.TargetResourceGroup))
	}
	if o.TargetResourceType != nil {
		out.Append("targetResourceType", fmt.Sprintf("%v", *o.TargetResourceType))
	}
	if o.TimeRange != nil {
		out.Append("timeRange", fmt.Sprintf("%v", *o.TimeRange))
	}
	return &out
}

// AlertsGetSummary ...
func (c AlertsManagementsClient) AlertsGetSummary(ctx context.Context, id commonids.SubscriptionId, options AlertsGetSummaryOperationOptions) (result AlertsGetSummaryOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		OptionsObject: options,
		Path:          fmt.Sprintf("%s/providers/Microsoft.AlertsManagement/alertsSummary", id.ID()),
	}

	req, err := c.Client.NewRequest(ctx, opts)
	if err != nil {
		return
	}

	var resp *client.Response
	resp, err = req.Execute(ctx)
	if resp != nil {
		result.OData = resp.OData
		result.HttpResponse = resp.Response
	}
	if err != nil {
		return
	}

	var model AlertsSummary
	result.Model = &model

	if err = resp.Unmarshal(result.Model); err != nil {
		return
	}

	return
}
