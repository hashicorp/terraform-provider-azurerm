package applications

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ApplicationListOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]Application
}

type ApplicationListCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []Application
}

type ApplicationListOperationOptions struct {
	Maxresults *int64
}

func DefaultApplicationListOperationOptions() ApplicationListOperationOptions {
	return ApplicationListOperationOptions{}
}

func (o ApplicationListOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o ApplicationListOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}

	return &out
}

func (o ApplicationListOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.Maxresults != nil {
		out.Append("maxresults", fmt.Sprintf("%v", *o.Maxresults))
	}
	return &out
}

type ApplicationListCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *ApplicationListCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// ApplicationList ...
func (c ApplicationsClient) ApplicationList(ctx context.Context, id BatchAccountId, options ApplicationListOperationOptions) (result ApplicationListOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		OptionsObject: options,
		Pager:         &ApplicationListCustomPager{},
		Path:          fmt.Sprintf("%s/applications", id.ID()),
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
		Values *[]Application `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ApplicationListComplete retrieves all the results into a single object
func (c ApplicationsClient) ApplicationListComplete(ctx context.Context, id BatchAccountId, options ApplicationListOperationOptions) (ApplicationListCompleteResult, error) {
	return c.ApplicationListCompleteMatchingPredicate(ctx, id, options, ApplicationOperationPredicate{})
}

// ApplicationListCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c ApplicationsClient) ApplicationListCompleteMatchingPredicate(ctx context.Context, id BatchAccountId, options ApplicationListOperationOptions, predicate ApplicationOperationPredicate) (result ApplicationListCompleteResult, err error) {
	items := make([]Application, 0)

	resp, err := c.ApplicationList(ctx, id, options)
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

	result = ApplicationListCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
