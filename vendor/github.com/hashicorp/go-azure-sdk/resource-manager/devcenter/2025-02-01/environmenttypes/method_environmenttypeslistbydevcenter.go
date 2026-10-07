package environmenttypes

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type EnvironmentTypesListByDevCenterOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]EnvironmentType
}

type EnvironmentTypesListByDevCenterCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []EnvironmentType
}

type EnvironmentTypesListByDevCenterOperationOptions struct {
	Top *int64
}

func DefaultEnvironmentTypesListByDevCenterOperationOptions() EnvironmentTypesListByDevCenterOperationOptions {
	return EnvironmentTypesListByDevCenterOperationOptions{}
}

func (o EnvironmentTypesListByDevCenterOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o EnvironmentTypesListByDevCenterOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}

	return &out
}

func (o EnvironmentTypesListByDevCenterOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.Top != nil {
		out.Append("$top", fmt.Sprintf("%v", *o.Top))
	}
	return &out
}

type EnvironmentTypesListByDevCenterCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *EnvironmentTypesListByDevCenterCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// EnvironmentTypesListByDevCenter ...
func (c EnvironmentTypesClient) EnvironmentTypesListByDevCenter(ctx context.Context, id DevCenterId, options EnvironmentTypesListByDevCenterOperationOptions) (result EnvironmentTypesListByDevCenterOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		OptionsObject: options,
		Pager:         &EnvironmentTypesListByDevCenterCustomPager{},
		Path:          fmt.Sprintf("%s/environmentTypes", id.ID()),
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
		Values *[]EnvironmentType `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// EnvironmentTypesListByDevCenterComplete retrieves all the results into a single object
func (c EnvironmentTypesClient) EnvironmentTypesListByDevCenterComplete(ctx context.Context, id DevCenterId, options EnvironmentTypesListByDevCenterOperationOptions) (EnvironmentTypesListByDevCenterCompleteResult, error) {
	return c.EnvironmentTypesListByDevCenterCompleteMatchingPredicate(ctx, id, options, EnvironmentTypeOperationPredicate{})
}

// EnvironmentTypesListByDevCenterCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c EnvironmentTypesClient) EnvironmentTypesListByDevCenterCompleteMatchingPredicate(ctx context.Context, id DevCenterId, options EnvironmentTypesListByDevCenterOperationOptions, predicate EnvironmentTypeOperationPredicate) (result EnvironmentTypesListByDevCenterCompleteResult, err error) {
	items := make([]EnvironmentType, 0)

	resp, err := c.EnvironmentTypesListByDevCenter(ctx, id, options)
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

	result = EnvironmentTypesListByDevCenterCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
