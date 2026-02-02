package giminorversions

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ListByParentOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]GiMinorVersion
}

type ListByParentCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []GiMinorVersion
}

type ListByParentOperationOptions struct {
	ShapeFamily *ShapeFamily
	Zone        *string
}

func DefaultListByParentOperationOptions() ListByParentOperationOptions {
	return ListByParentOperationOptions{}
}

func (o ListByParentOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o ListByParentOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}

	return &out
}

func (o ListByParentOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.ShapeFamily != nil {
		out.Append("shapeFamily", fmt.Sprintf("%v", *o.ShapeFamily))
	}
	if o.Zone != nil {
		out.Append("zone", fmt.Sprintf("%v", *o.Zone))
	}
	return &out
}

type ListByParentCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *ListByParentCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// ListByParent ...
func (c GiMinorVersionsClient) ListByParent(ctx context.Context, id GiVersionId, options ListByParentOperationOptions) (result ListByParentOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		OptionsObject: options,
		Pager:         &ListByParentCustomPager{},
		Path:          fmt.Sprintf("%s/giMinorVersions", id.ID()),
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
		Values *[]GiMinorVersion `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ListByParentComplete retrieves all the results into a single object
func (c GiMinorVersionsClient) ListByParentComplete(ctx context.Context, id GiVersionId, options ListByParentOperationOptions) (ListByParentCompleteResult, error) {
	return c.ListByParentCompleteMatchingPredicate(ctx, id, options, GiMinorVersionOperationPredicate{})
}

// ListByParentCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c GiMinorVersionsClient) ListByParentCompleteMatchingPredicate(ctx context.Context, id GiVersionId, options ListByParentOperationOptions, predicate GiMinorVersionOperationPredicate) (result ListByParentCompleteResult, err error) {
	items := make([]GiMinorVersion, 0)

	resp, err := c.ListByParent(ctx, id, options)
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

	result = ListByParentCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
