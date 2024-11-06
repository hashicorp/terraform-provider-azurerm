package containerappsrevisions

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ListRevisionsOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]Revision
}

type ListRevisionsCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []Revision
}

type ListRevisionsOperationOptions struct {
	Filter *string
}

func DefaultListRevisionsOperationOptions() ListRevisionsOperationOptions {
	return ListRevisionsOperationOptions{}
}

func (o ListRevisionsOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o ListRevisionsOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}

	return &out
}

func (o ListRevisionsOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.Filter != nil {
		out.Append("$filter", fmt.Sprintf("%v", *o.Filter))
	}
	return &out
}

type ListRevisionsCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *ListRevisionsCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// ListRevisions ...
func (c ContainerAppsRevisionsClient) ListRevisions(ctx context.Context, id ContainerAppId, options ListRevisionsOperationOptions) (result ListRevisionsOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		OptionsObject: options,
		Pager:         &ListRevisionsCustomPager{},
		Path:          fmt.Sprintf("%s/revisions", id.ID()),
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
		Values *[]Revision `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ListRevisionsComplete retrieves all the results into a single object
func (c ContainerAppsRevisionsClient) ListRevisionsComplete(ctx context.Context, id ContainerAppId, options ListRevisionsOperationOptions) (ListRevisionsCompleteResult, error) {
	return c.ListRevisionsCompleteMatchingPredicate(ctx, id, options, RevisionOperationPredicate{})
}

// ListRevisionsCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c ContainerAppsRevisionsClient) ListRevisionsCompleteMatchingPredicate(ctx context.Context, id ContainerAppId, options ListRevisionsOperationOptions, predicate RevisionOperationPredicate) (result ListRevisionsCompleteResult, err error) {
	items := make([]Revision, 0)

	resp, err := c.ListRevisions(ctx, id, options)
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

	result = ListRevisionsCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
