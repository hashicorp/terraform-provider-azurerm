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

type DeleteStaticSiteCustomDomainOperationResponse struct {
	Poller       pollers.Poller
	HttpResponse *http.Response
	OData        *odata.OData
}

// DeleteStaticSiteCustomDomain ...
func (c StaticSitesClient) DeleteStaticSiteCustomDomain(ctx context.Context, id CustomDomainId) (result DeleteStaticSiteCustomDomainOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusAccepted,
			http.StatusOK,
		},
		HttpMethod: http.MethodDelete,
		Path:       id.ID(),
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

	result.Poller, err = resourcemanager.PollerFromResponse(resp, c.Client)
	if err != nil {
		return
	}

	return
}

// DeleteStaticSiteCustomDomainThenPoll performs DeleteStaticSiteCustomDomain then polls until it's completed
func (c StaticSitesClient) DeleteStaticSiteCustomDomainThenPoll(ctx context.Context, id CustomDomainId) error {
	result, err := c.DeleteStaticSiteCustomDomain(ctx, id)
	if err != nil {
		return fmt.Errorf("performing DeleteStaticSiteCustomDomain: %+v", err)
	}

	if err := result.Poller.PollUntilDone(ctx); err != nil {
		return fmt.Errorf("polling after DeleteStaticSiteCustomDomain: %+v", err)
	}

	return nil
}
