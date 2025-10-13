package dbsystemshapes

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ListByLocationOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]DbSystemShape
}

type ListByLocationCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []DbSystemShape
}

type ListByLocationOperationOptions struct {
	Zone *string
}

func DefaultListByLocationOperationOptions() ListByLocationOperationOptions {
	return ListByLocationOperationOptions{}
}

func (o ListByLocationOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o ListByLocationOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}

	return &out
}

func (o ListByLocationOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.Zone != nil {
		out.Append("zone", fmt.Sprintf("%v", *o.Zone))
	}
	return &out
}

type ListByLocationCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *ListByLocationCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// ListByLocation ...
func (c DbSystemShapesClient) ListByLocation(ctx context.Context, id LocationId, options ListByLocationOperationOptions) (result ListByLocationOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		OptionsObject: options,
		Pager:         &ListByLocationCustomPager{},
		Path:          fmt.Sprintf("%s/dbSystemShapes", id.ID()),
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
		Values *[]DbSystemShape `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ListByLocationComplete retrieves all the results into a single object
func (c DbSystemShapesClient) ListByLocationComplete(ctx context.Context, id LocationId, options ListByLocationOperationOptions) (ListByLocationCompleteResult, error) {
	return c.ListByLocationCompleteMatchingPredicate(ctx, id, options, DbSystemShapeOperationPredicate{})
}

// ListByLocationCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c DbSystemShapesClient) ListByLocationCompleteMatchingPredicate(ctx context.Context, id LocationId, options ListByLocationOperationOptions, predicate DbSystemShapeOperationPredicate) (result ListByLocationCompleteResult, err error) {
	items := make([]DbSystemShape, 0)

	resp, err := c.ListByLocation(ctx, id, options)
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

	result = ListByLocationCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
