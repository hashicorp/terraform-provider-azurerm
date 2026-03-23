package localrulestackresources

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type LocalRulestackslistAppIdsOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]string
}

type LocalRulestackslistAppIdsCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []string
}

type LocalRulestackslistAppIdsOperationOptions struct {
	AppIdVersion *string
	AppPrefix    *string
	Skip         *string
	Top          *int64
}

func DefaultLocalRulestackslistAppIdsOperationOptions() LocalRulestackslistAppIdsOperationOptions {
	return LocalRulestackslistAppIdsOperationOptions{}
}

func (o LocalRulestackslistAppIdsOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o LocalRulestackslistAppIdsOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}

	return &out
}

func (o LocalRulestackslistAppIdsOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.AppIdVersion != nil {
		out.Append("appIdVersion", fmt.Sprintf("%v", *o.AppIdVersion))
	}
	if o.AppPrefix != nil {
		out.Append("appPrefix", fmt.Sprintf("%v", *o.AppPrefix))
	}
	if o.Skip != nil {
		out.Append("skip", fmt.Sprintf("%v", *o.Skip))
	}
	if o.Top != nil {
		out.Append("top", fmt.Sprintf("%v", *o.Top))
	}
	return &out
}

type LocalRulestackslistAppIdsCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *LocalRulestackslistAppIdsCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// LocalRulestackslistAppIds ...
func (c LocalRulestackResourcesClient) LocalRulestackslistAppIds(ctx context.Context, id LocalRulestackId, options LocalRulestackslistAppIdsOperationOptions) (result LocalRulestackslistAppIdsOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodPost,
		OptionsObject: options,
		Pager:         &LocalRulestackslistAppIdsCustomPager{},
		Path:          fmt.Sprintf("%s/listAppIds", id.ID()),
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
		Values *[]string `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// LocalRulestackslistAppIdsComplete retrieves all the results into a single object
func (c LocalRulestackResourcesClient) LocalRulestackslistAppIdsComplete(ctx context.Context, id LocalRulestackId, options LocalRulestackslistAppIdsOperationOptions) (result LocalRulestackslistAppIdsCompleteResult, err error) {
	items := make([]string, 0)

	resp, err := c.LocalRulestackslistAppIds(ctx, id, options)
	if err != nil {
		result.LatestHttpResponse = resp.HttpResponse
		err = fmt.Errorf("loading results: %+v", err)
		return
	}
	if resp.Model != nil {
		for _, v := range *resp.Model {
			items = append(items, v)
		}
	}

	result = LocalRulestackslistAppIdsCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
