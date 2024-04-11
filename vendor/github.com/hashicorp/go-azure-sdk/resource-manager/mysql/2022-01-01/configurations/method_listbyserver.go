package configurations

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ListByServerOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]Configuration
}

type ListByServerCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []Configuration
}

type ListByServerOperationOptions struct {
	Keyword  *string
	Page     *int64
	PageSize *int64
	Tags     *string
}

func DefaultListByServerOperationOptions() ListByServerOperationOptions {
	return ListByServerOperationOptions{}
}

func (o ListByServerOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o ListByServerOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}
	return &out
}

func (o ListByServerOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.Keyword != nil {
		out.Append("keyword", fmt.Sprintf("%v", *o.Keyword))
	}
	if o.Page != nil {
		out.Append("page", fmt.Sprintf("%v", *o.Page))
	}
	if o.PageSize != nil {
		out.Append("pageSize", fmt.Sprintf("%v", *o.PageSize))
	}
	if o.Tags != nil {
		out.Append("tags", fmt.Sprintf("%v", *o.Tags))
	}
	return &out
}

// ListByServer ...
func (c ConfigurationsClient) ListByServer(ctx context.Context, id FlexibleServerId, options ListByServerOperationOptions) (result ListByServerOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		Path:          fmt.Sprintf("%s/configurations", id.ID()),
		OptionsObject: options,
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
		Values *[]Configuration `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ListByServerComplete retrieves all the results into a single object
func (c ConfigurationsClient) ListByServerComplete(ctx context.Context, id FlexibleServerId, options ListByServerOperationOptions) (ListByServerCompleteResult, error) {
	return c.ListByServerCompleteMatchingPredicate(ctx, id, options, ConfigurationOperationPredicate{})
}

// ListByServerCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c ConfigurationsClient) ListByServerCompleteMatchingPredicate(ctx context.Context, id FlexibleServerId, options ListByServerOperationOptions, predicate ConfigurationOperationPredicate) (result ListByServerCompleteResult, err error) {
	items := make([]Configuration, 0)

	resp, err := c.ListByServer(ctx, id, options)
	if err != nil {
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

	result = ListByServerCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
