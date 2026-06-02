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

type ProjectEnvironmentTypesListOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]ProjectEnvironmentType
}

type ProjectEnvironmentTypesListCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []ProjectEnvironmentType
}

type ProjectEnvironmentTypesListOperationOptions struct {
	Top *int64
}

func DefaultProjectEnvironmentTypesListOperationOptions() ProjectEnvironmentTypesListOperationOptions {
	return ProjectEnvironmentTypesListOperationOptions{}
}

func (o ProjectEnvironmentTypesListOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o ProjectEnvironmentTypesListOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}

	return &out
}

func (o ProjectEnvironmentTypesListOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.Top != nil {
		out.Append("$top", fmt.Sprintf("%v", *o.Top))
	}
	return &out
}

type ProjectEnvironmentTypesListCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *ProjectEnvironmentTypesListCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// ProjectEnvironmentTypesList ...
func (c EnvironmentTypesClient) ProjectEnvironmentTypesList(ctx context.Context, id ProjectId, options ProjectEnvironmentTypesListOperationOptions) (result ProjectEnvironmentTypesListOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		OptionsObject: options,
		Pager:         &ProjectEnvironmentTypesListCustomPager{},
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
		Values *[]ProjectEnvironmentType `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ProjectEnvironmentTypesListComplete retrieves all the results into a single object
func (c EnvironmentTypesClient) ProjectEnvironmentTypesListComplete(ctx context.Context, id ProjectId, options ProjectEnvironmentTypesListOperationOptions) (ProjectEnvironmentTypesListCompleteResult, error) {
	return c.ProjectEnvironmentTypesListCompleteMatchingPredicate(ctx, id, options, ProjectEnvironmentTypeOperationPredicate{})
}

// ProjectEnvironmentTypesListCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c EnvironmentTypesClient) ProjectEnvironmentTypesListCompleteMatchingPredicate(ctx context.Context, id ProjectId, options ProjectEnvironmentTypesListOperationOptions, predicate ProjectEnvironmentTypeOperationPredicate) (result ProjectEnvironmentTypesListCompleteResult, err error) {
	items := make([]ProjectEnvironmentType, 0)

	resp, err := c.ProjectEnvironmentTypesList(ctx, id, options)
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

	result = ProjectEnvironmentTypesListCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
