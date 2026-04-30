package jobcredentials

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ListByAgentOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]JobCredential
}

type ListByAgentCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []JobCredential
}

type ListByAgentCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *ListByAgentCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// ListByAgent ...
func (c JobCredentialsClient) ListByAgent(ctx context.Context, id JobAgentId) (result ListByAgentOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &ListByAgentCustomPager{},
		Path:       fmt.Sprintf("%s/credentials", id.ID()),
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
		Values *[]JobCredential `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ListByAgentComplete retrieves all the results into a single object
func (c JobCredentialsClient) ListByAgentComplete(ctx context.Context, id JobAgentId) (ListByAgentCompleteResult, error) {
	return c.ListByAgentCompleteMatchingPredicate(ctx, id, JobCredentialOperationPredicate{})
}

// ListByAgentCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c JobCredentialsClient) ListByAgentCompleteMatchingPredicate(ctx context.Context, id JobAgentId, predicate JobCredentialOperationPredicate) (result ListByAgentCompleteResult, err error) {
	items := make([]JobCredential, 0)

	resp, err := c.ListByAgent(ctx, id)
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

	result = ListByAgentCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
