package resourceproviders

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

type ListBillingMetersOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]BillingMeter
}

type ListBillingMetersCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []BillingMeter
}

type ListBillingMetersOperationOptions struct {
	BillingLocation *string
	OsType          *string
}

func DefaultListBillingMetersOperationOptions() ListBillingMetersOperationOptions {
	return ListBillingMetersOperationOptions{}
}

func (o ListBillingMetersOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o ListBillingMetersOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}

	return &out
}

func (o ListBillingMetersOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.BillingLocation != nil {
		out.Append("billingLocation", fmt.Sprintf("%v", *o.BillingLocation))
	}
	if o.OsType != nil {
		out.Append("osType", fmt.Sprintf("%v", *o.OsType))
	}
	return &out
}

type ListBillingMetersCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *ListBillingMetersCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// ListBillingMeters ...
func (c ResourceProvidersClient) ListBillingMeters(ctx context.Context, id commonids.SubscriptionId, options ListBillingMetersOperationOptions) (result ListBillingMetersOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		OptionsObject: options,
		Pager:         &ListBillingMetersCustomPager{},
		Path:          fmt.Sprintf("%s/providers/Microsoft.Web/billingMeters", id.ID()),
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
		Values *[]BillingMeter `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ListBillingMetersComplete retrieves all the results into a single object
func (c ResourceProvidersClient) ListBillingMetersComplete(ctx context.Context, id commonids.SubscriptionId, options ListBillingMetersOperationOptions) (ListBillingMetersCompleteResult, error) {
	return c.ListBillingMetersCompleteMatchingPredicate(ctx, id, options, BillingMeterOperationPredicate{})
}

// ListBillingMetersCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c ResourceProvidersClient) ListBillingMetersCompleteMatchingPredicate(ctx context.Context, id commonids.SubscriptionId, options ListBillingMetersOperationOptions, predicate BillingMeterOperationPredicate) (result ListBillingMetersCompleteResult, err error) {
	items := make([]BillingMeter, 0)

	resp, err := c.ListBillingMeters(ctx, id, options)
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

	result = ListBillingMetersCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
