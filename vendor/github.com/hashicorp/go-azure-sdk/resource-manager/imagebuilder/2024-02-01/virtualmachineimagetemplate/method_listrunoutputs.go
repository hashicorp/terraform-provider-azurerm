package virtualmachineimagetemplate

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ListRunOutputsOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]RunOutput
}

type ListRunOutputsCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []RunOutput
}

type ListRunOutputsCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *ListRunOutputsCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// ListRunOutputs ...
func (c VirtualMachineImageTemplateClient) ListRunOutputs(ctx context.Context, id ImageTemplateId) (result ListRunOutputsOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &ListRunOutputsCustomPager{},
		Path:       fmt.Sprintf("%s/runOutputs", id.ID()),
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
		Values *[]RunOutput `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ListRunOutputsComplete retrieves all the results into a single object
func (c VirtualMachineImageTemplateClient) ListRunOutputsComplete(ctx context.Context, id ImageTemplateId) (ListRunOutputsCompleteResult, error) {
	return c.ListRunOutputsCompleteMatchingPredicate(ctx, id, RunOutputOperationPredicate{})
}

// ListRunOutputsCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c VirtualMachineImageTemplateClient) ListRunOutputsCompleteMatchingPredicate(ctx context.Context, id ImageTemplateId, predicate RunOutputOperationPredicate) (result ListRunOutputsCompleteResult, err error) {
	items := make([]RunOutput, 0)

	resp, err := c.ListRunOutputs(ctx, id)
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

	result = ListRunOutputsCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
