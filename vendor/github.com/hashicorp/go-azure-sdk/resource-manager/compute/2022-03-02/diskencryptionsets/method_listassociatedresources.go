package diskencryptionsets

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ListAssociatedResourcesOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]string
}

type ListAssociatedResourcesCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []string
}

type ListAssociatedResourcesCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *ListAssociatedResourcesCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// ListAssociatedResources ...
func (c DiskEncryptionSetsClient) ListAssociatedResources(ctx context.Context, id commonids.DiskEncryptionSetId) (result ListAssociatedResourcesOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &ListAssociatedResourcesCustomPager{},
		Path:       fmt.Sprintf("%s/associatedResources", id.ID()),
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
		Values *[]string `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ListAssociatedResourcesComplete retrieves all the results into a single object
func (c DiskEncryptionSetsClient) ListAssociatedResourcesComplete(ctx context.Context, id commonids.DiskEncryptionSetId) (result ListAssociatedResourcesCompleteResult, err error) {
	items := make([]string, 0)

	resp, err := c.ListAssociatedResources(ctx, id)
	if err != nil {
		result.LatestHttpResponse = resp.HttpResponse
		err = fmt.Errorf("loading results: %+v", err)
		return
	}
	if resp.Model != nil {
		for _, v := range *resp.Model {
			items = append(items, v)
		}
	}

	result = ListAssociatedResourcesCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
