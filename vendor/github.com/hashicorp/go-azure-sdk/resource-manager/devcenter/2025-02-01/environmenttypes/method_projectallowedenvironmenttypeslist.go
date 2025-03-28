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

type ProjectAllowedEnvironmentTypesListOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]AllowedEnvironmentType
}

type ProjectAllowedEnvironmentTypesListCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []AllowedEnvironmentType
}

type ProjectAllowedEnvironmentTypesListOperationOptions struct {
	Top *int64
}

func DefaultProjectAllowedEnvironmentTypesListOperationOptions() ProjectAllowedEnvironmentTypesListOperationOptions {
	return ProjectAllowedEnvironmentTypesListOperationOptions{}
}

func (o ProjectAllowedEnvironmentTypesListOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o ProjectAllowedEnvironmentTypesListOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}

	return &out
}

func (o ProjectAllowedEnvironmentTypesListOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.Top != nil {
		out.Append("$top", fmt.Sprintf("%v", *o.Top))
	}
	return &out
}

type ProjectAllowedEnvironmentTypesListCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *ProjectAllowedEnvironmentTypesListCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// ProjectAllowedEnvironmentTypesList ...
func (c EnvironmentTypesClient) ProjectAllowedEnvironmentTypesList(ctx context.Context, id ProjectId, options ProjectAllowedEnvironmentTypesListOperationOptions) (result ProjectAllowedEnvironmentTypesListOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		OptionsObject: options,
		Pager:         &ProjectAllowedEnvironmentTypesListCustomPager{},
		Path:          fmt.Sprintf("%s/allowedEnvironmentTypes", id.ID()),
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
		Values *[]AllowedEnvironmentType `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ProjectAllowedEnvironmentTypesListComplete retrieves all the results into a single object
func (c EnvironmentTypesClient) ProjectAllowedEnvironmentTypesListComplete(ctx context.Context, id ProjectId, options ProjectAllowedEnvironmentTypesListOperationOptions) (ProjectAllowedEnvironmentTypesListCompleteResult, error) {
	return c.ProjectAllowedEnvironmentTypesListCompleteMatchingPredicate(ctx, id, options, AllowedEnvironmentTypeOperationPredicate{})
}

// ProjectAllowedEnvironmentTypesListCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c EnvironmentTypesClient) ProjectAllowedEnvironmentTypesListCompleteMatchingPredicate(ctx context.Context, id ProjectId, options ProjectAllowedEnvironmentTypesListOperationOptions, predicate AllowedEnvironmentTypeOperationPredicate) (result ProjectAllowedEnvironmentTypesListCompleteResult, err error) {
	items := make([]AllowedEnvironmentType, 0)

	resp, err := c.ProjectAllowedEnvironmentTypesList(ctx, id, options)
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

	result = ProjectAllowedEnvironmentTypesListCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
