package projectconnectionresource

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ProjectConnectionsListOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]ConnectionPropertiesV2BasicResource
}

type ProjectConnectionsListCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []ConnectionPropertiesV2BasicResource
}

type ProjectConnectionsListOperationOptions struct {
	Category   *string
	IncludeAll *bool
	Target     *string
}

func DefaultProjectConnectionsListOperationOptions() ProjectConnectionsListOperationOptions {
	return ProjectConnectionsListOperationOptions{}
}

func (o ProjectConnectionsListOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o ProjectConnectionsListOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}

	return &out
}

func (o ProjectConnectionsListOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.Category != nil {
		out.Append("category", fmt.Sprintf("%v", *o.Category))
	}
	if o.IncludeAll != nil {
		out.Append("includeAll", fmt.Sprintf("%v", *o.IncludeAll))
	}
	if o.Target != nil {
		out.Append("target", fmt.Sprintf("%v", *o.Target))
	}
	return &out
}

type ProjectConnectionsListCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *ProjectConnectionsListCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// ProjectConnectionsList ...
func (c ProjectConnectionResourceClient) ProjectConnectionsList(ctx context.Context, id ProjectId, options ProjectConnectionsListOperationOptions) (result ProjectConnectionsListOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		OptionsObject: options,
		Pager:         &ProjectConnectionsListCustomPager{},
		Path:          fmt.Sprintf("%s/connections", id.ID()),
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
		Values *[]ConnectionPropertiesV2BasicResource `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ProjectConnectionsListComplete retrieves all the results into a single object
func (c ProjectConnectionResourceClient) ProjectConnectionsListComplete(ctx context.Context, id ProjectId, options ProjectConnectionsListOperationOptions) (ProjectConnectionsListCompleteResult, error) {
	return c.ProjectConnectionsListCompleteMatchingPredicate(ctx, id, options, ConnectionPropertiesV2BasicResourceOperationPredicate{})
}

// ProjectConnectionsListCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c ProjectConnectionResourceClient) ProjectConnectionsListCompleteMatchingPredicate(ctx context.Context, id ProjectId, options ProjectConnectionsListOperationOptions, predicate ConnectionPropertiesV2BasicResourceOperationPredicate) (result ProjectConnectionsListCompleteResult, err error) {
	items := make([]ConnectionPropertiesV2BasicResource, 0)

	resp, err := c.ProjectConnectionsList(ctx, id, options)
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

	result = ProjectConnectionsListCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
