package workbooksapis

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type WorkbooksListBySubscriptionOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]Workbook
}

type WorkbooksListBySubscriptionCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []Workbook
}

type WorkbooksListBySubscriptionOperationOptions struct {
	CanFetchContent *bool
	Category        *CategoryType
	Tags            *string
}

func DefaultWorkbooksListBySubscriptionOperationOptions() WorkbooksListBySubscriptionOperationOptions {
	return WorkbooksListBySubscriptionOperationOptions{}
}

func (o WorkbooksListBySubscriptionOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o WorkbooksListBySubscriptionOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}

	return &out
}

func (o WorkbooksListBySubscriptionOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.CanFetchContent != nil {
		out.Append("canFetchContent", fmt.Sprintf("%v", *o.CanFetchContent))
	}
	if o.Category != nil {
		out.Append("category", fmt.Sprintf("%v", *o.Category))
	}
	if o.Tags != nil {
		out.Append("tags", fmt.Sprintf("%v", *o.Tags))
	}
	return &out
}

type WorkbooksListBySubscriptionCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *WorkbooksListBySubscriptionCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// WorkbooksListBySubscription ...
func (c WorkbooksAPIsClient) WorkbooksListBySubscription(ctx context.Context, id commonids.SubscriptionId, options WorkbooksListBySubscriptionOperationOptions) (result WorkbooksListBySubscriptionOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		OptionsObject: options,
		Pager:         &WorkbooksListBySubscriptionCustomPager{},
		Path:          fmt.Sprintf("%s/providers/Microsoft.Insights/workbooks", id.ID()),
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
		Values *[]Workbook `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// WorkbooksListBySubscriptionComplete retrieves all the results into a single object
func (c WorkbooksAPIsClient) WorkbooksListBySubscriptionComplete(ctx context.Context, id commonids.SubscriptionId, options WorkbooksListBySubscriptionOperationOptions) (WorkbooksListBySubscriptionCompleteResult, error) {
	return c.WorkbooksListBySubscriptionCompleteMatchingPredicate(ctx, id, options, WorkbookOperationPredicate{})
}

// WorkbooksListBySubscriptionCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c WorkbooksAPIsClient) WorkbooksListBySubscriptionCompleteMatchingPredicate(ctx context.Context, id commonids.SubscriptionId, options WorkbooksListBySubscriptionOperationOptions, predicate WorkbookOperationPredicate) (result WorkbooksListBySubscriptionCompleteResult, err error) {
	items := make([]Workbook, 0)

	resp, err := c.WorkbooksListBySubscription(ctx, id, options)
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

	result = WorkbooksListBySubscriptionCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
