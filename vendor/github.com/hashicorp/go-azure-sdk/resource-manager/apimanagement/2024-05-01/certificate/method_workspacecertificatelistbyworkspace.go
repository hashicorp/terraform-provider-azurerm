package certificate

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type WorkspaceCertificateListByWorkspaceOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]CertificateContract
}

type WorkspaceCertificateListByWorkspaceCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []CertificateContract
}

type WorkspaceCertificateListByWorkspaceOperationOptions struct {
	Filter                  *string
	IsKeyVaultRefreshFailed *bool
	Skip                    *int64
	Top                     *int64
}

func DefaultWorkspaceCertificateListByWorkspaceOperationOptions() WorkspaceCertificateListByWorkspaceOperationOptions {
	return WorkspaceCertificateListByWorkspaceOperationOptions{}
}

func (o WorkspaceCertificateListByWorkspaceOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o WorkspaceCertificateListByWorkspaceOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}

	return &out
}

func (o WorkspaceCertificateListByWorkspaceOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.Filter != nil {
		out.Append("$filter", fmt.Sprintf("%v", *o.Filter))
	}
	if o.IsKeyVaultRefreshFailed != nil {
		out.Append("isKeyVaultRefreshFailed", fmt.Sprintf("%v", *o.IsKeyVaultRefreshFailed))
	}
	if o.Skip != nil {
		out.Append("$skip", fmt.Sprintf("%v", *o.Skip))
	}
	if o.Top != nil {
		out.Append("$top", fmt.Sprintf("%v", *o.Top))
	}
	return &out
}

type WorkspaceCertificateListByWorkspaceCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *WorkspaceCertificateListByWorkspaceCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// WorkspaceCertificateListByWorkspace ...
func (c CertificateClient) WorkspaceCertificateListByWorkspace(ctx context.Context, id WorkspaceId, options WorkspaceCertificateListByWorkspaceOperationOptions) (result WorkspaceCertificateListByWorkspaceOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		OptionsObject: options,
		Pager:         &WorkspaceCertificateListByWorkspaceCustomPager{},
		Path:          fmt.Sprintf("%s/certificates", id.ID()),
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
		Values *[]CertificateContract `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// WorkspaceCertificateListByWorkspaceComplete retrieves all the results into a single object
func (c CertificateClient) WorkspaceCertificateListByWorkspaceComplete(ctx context.Context, id WorkspaceId, options WorkspaceCertificateListByWorkspaceOperationOptions) (WorkspaceCertificateListByWorkspaceCompleteResult, error) {
	return c.WorkspaceCertificateListByWorkspaceCompleteMatchingPredicate(ctx, id, options, CertificateContractOperationPredicate{})
}

// WorkspaceCertificateListByWorkspaceCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c CertificateClient) WorkspaceCertificateListByWorkspaceCompleteMatchingPredicate(ctx context.Context, id WorkspaceId, options WorkspaceCertificateListByWorkspaceOperationOptions, predicate CertificateContractOperationPredicate) (result WorkspaceCertificateListByWorkspaceCompleteResult, err error) {
	items := make([]CertificateContract, 0)

	resp, err := c.WorkspaceCertificateListByWorkspace(ctx, id, options)
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

	result = WorkspaceCertificateListByWorkspaceCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
