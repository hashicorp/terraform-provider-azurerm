// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package sdkhacks

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/resource-manager/web/2023-01-01/staticsites"
	"github.com/hashicorp/go-azure-sdk/sdk/client"
)

type StaticSitesClient struct {
	client *staticsites.StaticSitesClient
}

func NewStaticWebAppClient(client *staticsites.StaticSitesClient) StaticSitesClient {
	return StaticSitesClient{
		client: client,
	}
}

// CreateOrUpdateBasicAuth ...
func (c StaticSitesClient) CreateOrUpdateBasicAuth(ctx context.Context, id staticsites.StaticSiteId, input staticsites.StaticSiteBasicAuthPropertiesARMResource) (result staticsites.CreateOrUpdateBasicAuthOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodPut,
		Path:       fmt.Sprintf("%s/config/basicAuth", id.ID()),
	}

	req, err := c.client.Client.NewRequest(ctx, opts)
	if err != nil {
		return
	}

	if err = req.Marshal(input); err != nil {
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

	if err = resp.Unmarshal(&result.Model); err != nil {
		return
	}

	return
}

func (c StaticSitesClient) GetBasicAuth(ctx context.Context, id staticsites.StaticSiteId) (result staticsites.GetBasicAuthOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Path:       fmt.Sprintf("%s/basicAuth/default", id.ID()),
	}

	req, err := c.client.Client.NewRequest(ctx, opts)
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

	if err = resp.Unmarshal(&result.Model); err != nil {
		return
	}

	return
}
