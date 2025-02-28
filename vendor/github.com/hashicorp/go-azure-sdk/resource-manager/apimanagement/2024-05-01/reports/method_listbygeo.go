package reports

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ListByGeoOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]ReportRecordContract
}

type ListByGeoCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []ReportRecordContract
}

type ListByGeoOperationOptions struct {
	Filter *string
	Skip   *int64
	Top    *int64
}

func DefaultListByGeoOperationOptions() ListByGeoOperationOptions {
	return ListByGeoOperationOptions{}
}

func (o ListByGeoOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o ListByGeoOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}

	return &out
}

func (o ListByGeoOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.Filter != nil {
		out.Append("$filter", fmt.Sprintf("%v", *o.Filter))
	}
	if o.Skip != nil {
		out.Append("$skip", fmt.Sprintf("%v", *o.Skip))
	}
	if o.Top != nil {
		out.Append("$top", fmt.Sprintf("%v", *o.Top))
	}
	return &out
}

type ListByGeoCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *ListByGeoCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// ListByGeo ...
func (c ReportsClient) ListByGeo(ctx context.Context, id ServiceId, options ListByGeoOperationOptions) (result ListByGeoOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		OptionsObject: options,
		Pager:         &ListByGeoCustomPager{},
		Path:          fmt.Sprintf("%s/reports/byGeo", id.ID()),
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
		Values *[]ReportRecordContract `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ListByGeoComplete retrieves all the results into a single object
func (c ReportsClient) ListByGeoComplete(ctx context.Context, id ServiceId, options ListByGeoOperationOptions) (ListByGeoCompleteResult, error) {
	return c.ListByGeoCompleteMatchingPredicate(ctx, id, options, ReportRecordContractOperationPredicate{})
}

// ListByGeoCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c ReportsClient) ListByGeoCompleteMatchingPredicate(ctx context.Context, id ServiceId, options ListByGeoOperationOptions, predicate ReportRecordContractOperationPredicate) (result ListByGeoCompleteResult, err error) {
	items := make([]ReportRecordContract, 0)

	resp, err := c.ListByGeo(ctx, id, options)
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

	result = ListByGeoCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
