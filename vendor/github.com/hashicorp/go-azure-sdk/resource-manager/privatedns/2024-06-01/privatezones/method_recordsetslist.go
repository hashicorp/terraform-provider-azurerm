package privatezones

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RecordSetsListOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]RecordSet
}

type RecordSetsListCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []RecordSet
}

type RecordSetsListOperationOptions struct {
	Recordsetnamesuffix *string
	Top                 *int64
}

func DefaultRecordSetsListOperationOptions() RecordSetsListOperationOptions {
	return RecordSetsListOperationOptions{}
}

func (o RecordSetsListOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o RecordSetsListOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}

	return &out
}

func (o RecordSetsListOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.Recordsetnamesuffix != nil {
		out.Append("$recordsetnamesuffix", fmt.Sprintf("%v", *o.Recordsetnamesuffix))
	}
	if o.Top != nil {
		out.Append("$top", fmt.Sprintf("%v", *o.Top))
	}
	return &out
}

type RecordSetsListCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *RecordSetsListCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// RecordSetsList ...
func (c PrivateZonesClient) RecordSetsList(ctx context.Context, id PrivateDnsZoneId, options RecordSetsListOperationOptions) (result RecordSetsListOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		OptionsObject: options,
		Pager:         &RecordSetsListCustomPager{},
		Path:          fmt.Sprintf("%s/aLL", id.ID()),
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
		Values *[]RecordSet `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// RecordSetsListComplete retrieves all the results into a single object
func (c PrivateZonesClient) RecordSetsListComplete(ctx context.Context, id PrivateDnsZoneId, options RecordSetsListOperationOptions) (RecordSetsListCompleteResult, error) {
	return c.RecordSetsListCompleteMatchingPredicate(ctx, id, options, RecordSetOperationPredicate{})
}

// RecordSetsListCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c PrivateZonesClient) RecordSetsListCompleteMatchingPredicate(ctx context.Context, id PrivateDnsZoneId, options RecordSetsListOperationOptions, predicate RecordSetOperationPredicate) (result RecordSetsListCompleteResult, err error) {
	items := make([]RecordSet, 0)

	resp, err := c.RecordSetsList(ctx, id, options)
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

	result = RecordSetsListCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
