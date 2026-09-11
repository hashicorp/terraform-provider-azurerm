package msiximage

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ExpandOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]ExpandMsixImage
}

type ExpandCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []ExpandMsixImage
}

type ExpandCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *ExpandCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// Expand ...
func (c MsixImageClient) Expand(ctx context.Context, id HostPoolId, input MSIXImageURI) (result ExpandOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodPost,
		Pager:      &ExpandCustomPager{},
		Path:       fmt.Sprintf("%s/expandMsixImage", id.ID()),
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
		Values *[]ExpandMsixImage `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ExpandComplete retrieves all the results into a single object
func (c MsixImageClient) ExpandComplete(ctx context.Context, id HostPoolId, input MSIXImageURI) (ExpandCompleteResult, error) {
	return c.ExpandCompleteMatchingPredicate(ctx, id, input, ExpandMsixImageOperationPredicate{})
}

// ExpandCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c MsixImageClient) ExpandCompleteMatchingPredicate(ctx context.Context, id HostPoolId, input MSIXImageURI, predicate ExpandMsixImageOperationPredicate) (result ExpandCompleteResult, err error) {
	items := make([]ExpandMsixImage, 0)

	resp, err := c.Expand(ctx, id, input)
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

	result = ExpandCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
