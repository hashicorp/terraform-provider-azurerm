package servers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CapturedLogsListByServerOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]CapturedLog
}

type CapturedLogsListByServerCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []CapturedLog
}

type CapturedLogsListByServerCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *CapturedLogsListByServerCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// CapturedLogsListByServer ...
func (c ServersClient) CapturedLogsListByServer(ctx context.Context, id FlexibleServerId) (result CapturedLogsListByServerOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &CapturedLogsListByServerCustomPager{},
		Path:       fmt.Sprintf("%s/logFiles", id.ID()),
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
		Values *[]CapturedLog `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// CapturedLogsListByServerComplete retrieves all the results into a single object
func (c ServersClient) CapturedLogsListByServerComplete(ctx context.Context, id FlexibleServerId) (CapturedLogsListByServerCompleteResult, error) {
	return c.CapturedLogsListByServerCompleteMatchingPredicate(ctx, id, CapturedLogOperationPredicate{})
}

// CapturedLogsListByServerCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c ServersClient) CapturedLogsListByServerCompleteMatchingPredicate(ctx context.Context, id FlexibleServerId, predicate CapturedLogOperationPredicate) (result CapturedLogsListByServerCompleteResult, err error) {
	items := make([]CapturedLog, 0)

	resp, err := c.CapturedLogsListByServer(ctx, id)
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

	result = CapturedLogsListByServerCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
