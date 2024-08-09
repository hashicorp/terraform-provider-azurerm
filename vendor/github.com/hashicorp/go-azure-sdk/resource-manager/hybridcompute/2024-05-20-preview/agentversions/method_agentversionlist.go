package agentversions

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AgentVersionListOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]AgentVersion
}

type AgentVersionListCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []AgentVersion
}

type AgentVersionListCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *AgentVersionListCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// AgentVersionList ...
func (c AgentVersionsClient) AgentVersionList(ctx context.Context, id OsTypeId) (result AgentVersionListOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &AgentVersionListCustomPager{},
		Path:       fmt.Sprintf("%s/agentVersions", id.ID()),
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
		Values *[]AgentVersion `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// AgentVersionListComplete retrieves all the results into a single object
func (c AgentVersionsClient) AgentVersionListComplete(ctx context.Context, id OsTypeId) (AgentVersionListCompleteResult, error) {
	return c.AgentVersionListCompleteMatchingPredicate(ctx, id, AgentVersionOperationPredicate{})
}

// AgentVersionListCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c AgentVersionsClient) AgentVersionListCompleteMatchingPredicate(ctx context.Context, id OsTypeId, predicate AgentVersionOperationPredicate) (result AgentVersionListCompleteResult, err error) {
	items := make([]AgentVersion, 0)

	resp, err := c.AgentVersionList(ctx, id)
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

	result = AgentVersionListCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
