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

type AlertsGetAllOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]Alert
}

type AlertsGetAllCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []Alert
}

type AlertsGetAllOperationOptions struct {
	AlertRule           *string
	AlertState          *AlertState
	CustomTimeRange     *string
	IncludeContext      *bool
	IncludeEgressConfig *bool
	MonitorCondition    *MonitorCondition
	MonitorService      *MonitorService
	PageCount           *int64
	Select              *string
	Severity            *Severity
	SmartGroupId        *string
	SortBy              *AlertsSortByFields
	SortOrder           *SortOrder
	TargetResource      *string
	TargetResourceGroup *string
	TargetResourceType  *string
	TimeRange           *TimeRange
}

func DefaultAlertsGetAllOperationOptions() AlertsGetAllOperationOptions {
	return AlertsGetAllOperationOptions{}
}

func (o AlertsGetAllOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o AlertsGetAllOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}
	return &out
}

func (o AlertsGetAllOperationOptions) ToQuery() *client.QueryParams {
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
	if o.IncludeContext != nil {
		out.Append("includeContext", fmt.Sprintf("%v", *o.IncludeContext))
	}
	if o.IncludeEgressConfig != nil {
		out.Append("includeEgressConfig", fmt.Sprintf("%v", *o.IncludeEgressConfig))
	}
	if o.MonitorCondition != nil {
		out.Append("monitorCondition", fmt.Sprintf("%v", *o.MonitorCondition))
	}
	if o.MonitorService != nil {
		out.Append("monitorService", fmt.Sprintf("%v", *o.MonitorService))
	}
	if o.PageCount != nil {
		out.Append("pageCount", fmt.Sprintf("%v", *o.PageCount))
	}
	if o.Select != nil {
		out.Append("select", fmt.Sprintf("%v", *o.Select))
	}
	if o.Severity != nil {
		out.Append("severity", fmt.Sprintf("%v", *o.Severity))
	}
	if o.SmartGroupId != nil {
		out.Append("smartGroupId", fmt.Sprintf("%v", *o.SmartGroupId))
	}
	if o.SortBy != nil {
		out.Append("sortBy", fmt.Sprintf("%v", *o.SortBy))
	}
	if o.SortOrder != nil {
		out.Append("sortOrder", fmt.Sprintf("%v", *o.SortOrder))
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

// AlertsGetAll ...
func (c AlertsManagementsClient) AlertsGetAll(ctx context.Context, id commonids.SubscriptionId, options AlertsGetAllOperationOptions) (result AlertsGetAllOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		Path:          fmt.Sprintf("%s/providers/Microsoft.AlertsManagement/alerts", id.ID()),
		OptionsObject: options,
	}

	req, err := c.Client.NewRequest(ctx, opts)
	if err != nil {
		return
	}

	var resp *client.Response
	resp, err = req.ExecutePaged(ctx)
	if resp != nil {
		result.OData = resp.OData
		result.HttpResponse = resp.Response
	}
	if err != nil {
		return
	}

	var values struct {
		Values *[]Alert `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// AlertsGetAllComplete retrieves all the results into a single object
func (c AlertsManagementsClient) AlertsGetAllComplete(ctx context.Context, id commonids.SubscriptionId, options AlertsGetAllOperationOptions) (AlertsGetAllCompleteResult, error) {
	return c.AlertsGetAllCompleteMatchingPredicate(ctx, id, options, AlertOperationPredicate{})
}

// AlertsGetAllCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c AlertsManagementsClient) AlertsGetAllCompleteMatchingPredicate(ctx context.Context, id commonids.SubscriptionId, options AlertsGetAllOperationOptions, predicate AlertOperationPredicate) (result AlertsGetAllCompleteResult, err error) {
	items := make([]Alert, 0)

	resp, err := c.AlertsGetAll(ctx, id, options)
	if err != nil {
		err = fmt.Errorf("loading results: %+v", err)
		return
	}
	if resp.Model != nil {
		for _, v := range *resp.Model {
			if predicate.Matches(v) {
				items = append(items, v)
			}
		}
	}

	result = AlertsGetAllCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
