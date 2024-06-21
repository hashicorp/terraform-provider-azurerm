package actionrules

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

type ListByResourceGroupOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]ActionRule
}

type ListByResourceGroupCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []ActionRule
}

type ListByResourceGroupOperationOptions struct {
	ActionGroup         *string
	AlertRuleId         *string
	Description         *string
	ImpactedScope       *string
	MonitorService      *MonitorService
	Name                *string
	Severity            *Severity
	TargetResource      *string
	TargetResourceGroup *string
	TargetResourceType  *string
}

func DefaultListByResourceGroupOperationOptions() ListByResourceGroupOperationOptions {
	return ListByResourceGroupOperationOptions{}
}

func (o ListByResourceGroupOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o ListByResourceGroupOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}
	return &out
}

func (o ListByResourceGroupOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.ActionGroup != nil {
		out.Append("actionGroup", fmt.Sprintf("%v", *o.ActionGroup))
	}
	if o.AlertRuleId != nil {
		out.Append("alertRuleId", fmt.Sprintf("%v", *o.AlertRuleId))
	}
	if o.Description != nil {
		out.Append("description", fmt.Sprintf("%v", *o.Description))
	}
	if o.ImpactedScope != nil {
		out.Append("impactedScope", fmt.Sprintf("%v", *o.ImpactedScope))
	}
	if o.MonitorService != nil {
		out.Append("monitorService", fmt.Sprintf("%v", *o.MonitorService))
	}
	if o.Name != nil {
		out.Append("name", fmt.Sprintf("%v", *o.Name))
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
	return &out
}

// ListByResourceGroup ...
func (c ActionRulesClient) ListByResourceGroup(ctx context.Context, id commonids.ResourceGroupId, options ListByResourceGroupOperationOptions) (result ListByResourceGroupOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		Path:          fmt.Sprintf("%s/providers/Microsoft.AlertsManagement/actionRules", id.ID()),
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
		Values *[]ActionRule `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ListByResourceGroupComplete retrieves all the results into a single object
func (c ActionRulesClient) ListByResourceGroupComplete(ctx context.Context, id commonids.ResourceGroupId, options ListByResourceGroupOperationOptions) (ListByResourceGroupCompleteResult, error) {
	return c.ListByResourceGroupCompleteMatchingPredicate(ctx, id, options, ActionRuleOperationPredicate{})
}

// ListByResourceGroupCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c ActionRulesClient) ListByResourceGroupCompleteMatchingPredicate(ctx context.Context, id commonids.ResourceGroupId, options ListByResourceGroupOperationOptions, predicate ActionRuleOperationPredicate) (result ListByResourceGroupCompleteResult, err error) {
	items := make([]ActionRule, 0)

	resp, err := c.ListByResourceGroup(ctx, id, options)
	if err != nil {
		result.LatestHttpResponse = resp.HttpResponse
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

	result = ListByResourceGroupCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
