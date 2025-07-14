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

type ListByDevCenterOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]DevBoxDefinition
}

type ListByDevCenterCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []DevBoxDefinition
}

type ListByDevCenterCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *ListByDevCenterCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// ListByDevCenter ...
func (c DevBoxDefinitionsClient) ListByDevCenter(ctx context.Context, id DevCenterId) (result ListByDevCenterOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &ListByDevCenterCustomPager{},
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

// ListByDevCenterComplete retrieves all the results into a single object
func (c DevBoxDefinitionsClient) ListByDevCenterComplete(ctx context.Context, id DevCenterId) (ListByDevCenterCompleteResult, error) {
	return c.ListByDevCenterCompleteMatchingPredicate(ctx, id, DevBoxDefinitionOperationPredicate{})
}

// ListByDevCenterCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c DevBoxDefinitionsClient) ListByDevCenterCompleteMatchingPredicate(ctx context.Context, id DevCenterId, predicate DevBoxDefinitionOperationPredicate) (result ListByDevCenterCompleteResult, err error) {
	items := make([]DevBoxDefinition, 0)

	resp, err := c.ListByDevCenter(ctx, id)
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

	result = ListByDevCenterCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
