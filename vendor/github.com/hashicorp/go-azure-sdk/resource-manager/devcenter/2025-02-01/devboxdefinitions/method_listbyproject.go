package devboxdefinitions

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ListByProjectOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]DevBoxDefinition
}

type ListByProjectCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []DevBoxDefinition
}

type ListByProjectCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *ListByProjectCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// ListByProject ...
func (c DevBoxDefinitionsClient) ListByProject(ctx context.Context, id ProjectId) (result ListByProjectOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &ListByProjectCustomPager{},
		Path:       fmt.Sprintf("%s/devBoxDefinitions", id.ID()),
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
		Values *[]DevBoxDefinition `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ListByProjectComplete retrieves all the results into a single object
func (c DevBoxDefinitionsClient) ListByProjectComplete(ctx context.Context, id ProjectId) (ListByProjectCompleteResult, error) {
	return c.ListByProjectCompleteMatchingPredicate(ctx, id, DevBoxDefinitionOperationPredicate{})
}

// ListByProjectCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c DevBoxDefinitionsClient) ListByProjectCompleteMatchingPredicate(ctx context.Context, id ProjectId, predicate DevBoxDefinitionOperationPredicate) (result ListByProjectCompleteResult, err error) {
	items := make([]DevBoxDefinition, 0)

	resp, err := c.ListByProject(ctx, id)
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

	result = ListByProjectCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
