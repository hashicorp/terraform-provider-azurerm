package projectskus

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SkusListByProjectOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]DevCenterSku
}

type SkusListByProjectCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []DevCenterSku
}

type SkusListByProjectCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *SkusListByProjectCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// SkusListByProject ...
func (c ProjectSKUsClient) SkusListByProject(ctx context.Context, id ProjectId) (result SkusListByProjectOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodPost,
		Pager:      &SkusListByProjectCustomPager{},
		Path:       fmt.Sprintf("%s/listSkus", id.ID()),
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
		Values *[]DevCenterSku `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// SkusListByProjectComplete retrieves all the results into a single object
func (c ProjectSKUsClient) SkusListByProjectComplete(ctx context.Context, id ProjectId) (SkusListByProjectCompleteResult, error) {
	return c.SkusListByProjectCompleteMatchingPredicate(ctx, id, DevCenterSkuOperationPredicate{})
}

// SkusListByProjectCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c ProjectSKUsClient) SkusListByProjectCompleteMatchingPredicate(ctx context.Context, id ProjectId, predicate DevCenterSkuOperationPredicate) (result SkusListByProjectCompleteResult, err error) {
	items := make([]DevCenterSku, 0)

	resp, err := c.SkusListByProject(ctx, id)
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

	result = SkusListByProjectCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
