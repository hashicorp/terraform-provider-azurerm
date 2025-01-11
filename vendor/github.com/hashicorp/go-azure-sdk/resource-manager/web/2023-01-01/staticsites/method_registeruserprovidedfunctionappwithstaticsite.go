package staticsites

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RegisterUserProvidedFunctionAppWithStaticSiteOperationResponse struct {
	Poller       pollers.Poller
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *StaticSiteUserProvidedFunctionAppARMResource
}

type RegisterUserProvidedFunctionAppWithStaticSiteOperationOptions struct {
	IsForced *bool
}

func DefaultRegisterUserProvidedFunctionAppWithStaticSiteOperationOptions() RegisterUserProvidedFunctionAppWithStaticSiteOperationOptions {
	return RegisterUserProvidedFunctionAppWithStaticSiteOperationOptions{}
}

func (o RegisterUserProvidedFunctionAppWithStaticSiteOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o RegisterUserProvidedFunctionAppWithStaticSiteOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}

	return &out
}

func (o RegisterUserProvidedFunctionAppWithStaticSiteOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.IsForced != nil {
		out.Append("isForced", fmt.Sprintf("%v", *o.IsForced))
	}
	return &out
}

// RegisterUserProvidedFunctionAppWithStaticSite ...
func (c StaticSitesClient) RegisterUserProvidedFunctionAppWithStaticSite(ctx context.Context, id UserProvidedFunctionAppId, input StaticSiteUserProvidedFunctionAppARMResource, options RegisterUserProvidedFunctionAppWithStaticSiteOperationOptions) (result RegisterUserProvidedFunctionAppWithStaticSiteOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusAccepted,
			http.StatusOK,
		},
		HttpMethod:    http.MethodPut,
		OptionsObject: options,
		Path:          id.ID(),
	}

	req, err := c.Client.NewRequest(ctx, opts)
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

	result.Poller, err = resourcemanager.PollerFromResponse(resp, c.Client)
	if err != nil {
		return
	}

	return
}

// RegisterUserProvidedFunctionAppWithStaticSiteThenPoll performs RegisterUserProvidedFunctionAppWithStaticSite then polls until it's completed
func (c StaticSitesClient) RegisterUserProvidedFunctionAppWithStaticSiteThenPoll(ctx context.Context, id UserProvidedFunctionAppId, input StaticSiteUserProvidedFunctionAppARMResource, options RegisterUserProvidedFunctionAppWithStaticSiteOperationOptions) error {
	result, err := c.RegisterUserProvidedFunctionAppWithStaticSite(ctx, id, input, options)
	if err != nil {
		return fmt.Errorf("performing RegisterUserProvidedFunctionAppWithStaticSite: %+v", err)
	}

	if err := result.Poller.PollUntilDone(ctx); err != nil {
		return fmt.Errorf("polling after RegisterUserProvidedFunctionAppWithStaticSite: %+v", err)
	}

	return nil
}
