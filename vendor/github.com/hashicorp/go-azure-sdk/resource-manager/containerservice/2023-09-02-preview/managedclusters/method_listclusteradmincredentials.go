package managedclusters

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

type ListClusterAdminCredentialsOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *CredentialResults
}

type ListClusterAdminCredentialsOperationOptions struct {
	ServerFqdn *string
}

func DefaultListClusterAdminCredentialsOperationOptions() ListClusterAdminCredentialsOperationOptions {
	return ListClusterAdminCredentialsOperationOptions{}
}

func (o ListClusterAdminCredentialsOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o ListClusterAdminCredentialsOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}
	return &out
}

func (o ListClusterAdminCredentialsOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.ServerFqdn != nil {
		out.Append("server-fqdn", fmt.Sprintf("%v", *o.ServerFqdn))
	}
	return &out
}

// ListClusterAdminCredentials ...
func (c ManagedClustersClient) ListClusterAdminCredentials(ctx context.Context, id commonids.KubernetesClusterId, options ListClusterAdminCredentialsOperationOptions) (result ListClusterAdminCredentialsOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodPost,
		Path:          fmt.Sprintf("%s/listClusterAdminCredential", id.ID()),
		OptionsObject: options,
	}

	req, err := c.Client.NewRequest(ctx, opts)
	if err != nil {
		return
	}

	var resp *client.Response
	resp, err = req.Execute(ctx)
	if resp != nil {
		result.OData = resp.OData
		result.HttpResponse = resp.Response
	}
	if err != nil {
		return
	}

	var model CredentialResults
	result.Model = &model

	if err = resp.Unmarshal(result.Model); err != nil {
		return
	}

	return
}
