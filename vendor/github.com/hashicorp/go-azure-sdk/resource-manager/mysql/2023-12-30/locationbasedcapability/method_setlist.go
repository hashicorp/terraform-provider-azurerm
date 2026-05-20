package locationbasedcapability

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SetListOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]Capability
}

type SetListCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []Capability
}

type SetListCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *SetListCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// SetList ...
func (c LocationBasedCapabilityClient) SetList(ctx context.Context, id LocationId) (result SetListOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &SetListCustomPager{},
		Path:       fmt.Sprintf("%s/capabilitySets", id.ID()),
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
		Values *[]Capability `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// SetListComplete retrieves all the results into a single object
func (c LocationBasedCapabilityClient) SetListComplete(ctx context.Context, id LocationId) (SetListCompleteResult, error) {
	return c.SetListCompleteMatchingPredicate(ctx, id, CapabilityOperationPredicate{})
}

// SetListCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c LocationBasedCapabilityClient) SetListCompleteMatchingPredicate(ctx context.Context, id LocationId, predicate CapabilityOperationPredicate) (result SetListCompleteResult, err error) {
	items := make([]Capability, 0)

	resp, err := c.SetList(ctx, id)
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

	result = SetListCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
