package commitmentplans

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ListAssociationsOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]CommitmentPlanAccountAssociation
}

type ListAssociationsCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []CommitmentPlanAccountAssociation
}

type ListAssociationsCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *ListAssociationsCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// ListAssociations ...
func (c CommitmentPlansClient) ListAssociations(ctx context.Context, id CommitmentPlanId) (result ListAssociationsOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &ListAssociationsCustomPager{},
		Path:       fmt.Sprintf("%s/accountAssociations", id.ID()),
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
		Values *[]CommitmentPlanAccountAssociation `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ListAssociationsComplete retrieves all the results into a single object
func (c CommitmentPlansClient) ListAssociationsComplete(ctx context.Context, id CommitmentPlanId) (ListAssociationsCompleteResult, error) {
	return c.ListAssociationsCompleteMatchingPredicate(ctx, id, CommitmentPlanAccountAssociationOperationPredicate{})
}

// ListAssociationsCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c CommitmentPlansClient) ListAssociationsCompleteMatchingPredicate(ctx context.Context, id CommitmentPlanId, predicate CommitmentPlanAccountAssociationOperationPredicate) (result ListAssociationsCompleteResult, err error) {
	items := make([]CommitmentPlanAccountAssociation, 0)

	resp, err := c.ListAssociations(ctx, id)
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

	result = ListAssociationsCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
