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

type LocalRulestackslistFirewallsOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]string
}

type LocalRulestackslistFirewallsCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []string
}

type LocalRulestackslistFirewallsCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *LocalRulestackslistFirewallsCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// LocalRulestackslistFirewalls ...
func (c LocalRulestackResourcesClient) LocalRulestackslistFirewalls(ctx context.Context, id LocalRulestackId) (result LocalRulestackslistFirewallsOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodPost,
		Pager:      &LocalRulestackslistFirewallsCustomPager{},
		Path:       fmt.Sprintf("%s/listFirewalls", id.ID()),
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

// LocalRulestackslistFirewallsComplete retrieves all the results into a single object
func (c LocalRulestackResourcesClient) LocalRulestackslistFirewallsComplete(ctx context.Context, id LocalRulestackId) (result LocalRulestackslistFirewallsCompleteResult, err error) {
	items := make([]string, 0)

	resp, err := c.LocalRulestackslistFirewalls(ctx, id)
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

	result = LocalRulestackslistFirewallsCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
