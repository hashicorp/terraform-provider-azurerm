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

type ListByServerCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *ListByServerCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// ListByServer ...
func (c ConfigurationsClient) ListByServer(ctx context.Context, id FlexibleServerId) (result ListByServerOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &ListByServerCustomPager{},
		Path:       fmt.Sprintf("%s/configurations", id.ID()),
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
func (c ConfigurationsClient) ListByServerComplete(ctx context.Context, id FlexibleServerId) (ListByServerCompleteResult, error) {
	return c.ListByServerCompleteMatchingPredicate(ctx, id, ConfigurationOperationPredicate{})
}

// ListByServerCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c ConfigurationsClient) ListByServerCompleteMatchingPredicate(ctx context.Context, id FlexibleServerId, predicate ConfigurationOperationPredicate) (result ListByServerCompleteResult, err error) {
	items := make([]Configuration, 0)

	resp, err := c.ListByServer(ctx, id)
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

	result = ListByServerCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
