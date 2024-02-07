package singlesignon

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

type ConfigurationsCreateOrUpdateOperationResponse struct {
	Poller       pollers.Poller
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *DatadogSingleSignOnResource
}

// ConfigurationsCreateOrUpdate ...
func (c SingleSignOnClient) ConfigurationsCreateOrUpdate(ctx context.Context, id SingleSignOnConfigurationId, input DatadogSingleSignOnResource) (result ConfigurationsCreateOrUpdateOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusCreated,
			http.StatusOK,
		},
		HttpMethod: http.MethodPut,
		Path:       id.ID(),
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

// ConfigurationsCreateOrUpdateThenPoll performs ConfigurationsCreateOrUpdate then polls until it's completed
func (c SingleSignOnClient) ConfigurationsCreateOrUpdateThenPoll(ctx context.Context, id SingleSignOnConfigurationId, input DatadogSingleSignOnResource) error {
	result, err := c.ConfigurationsCreateOrUpdate(ctx, id, input)
	if err != nil {
		return fmt.Errorf("performing ConfigurationsCreateOrUpdate: %+v", err)
	}

	if err := result.Poller.PollUntilDone(ctx); err != nil {
		return fmt.Errorf("polling after ConfigurationsCreateOrUpdate: %+v", err)
	}

	return nil
}
