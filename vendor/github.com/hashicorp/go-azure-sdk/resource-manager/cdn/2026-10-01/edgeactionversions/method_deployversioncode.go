package edgeactionversions

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

type DeployVersionCodeOperationResponse struct {
	Poller       pollers.Poller
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *EdgeActionVersion
}

// DeployVersionCode ...
func (c EdgeActionVersionsClient) DeployVersionCode(ctx context.Context, id VersionId, input VersionCode) (result DeployVersionCodeOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusAccepted,
			http.StatusOK,
		},
		HttpMethod: http.MethodPost,
		Path:       fmt.Sprintf("%s/deployVersionCode", id.ID()),
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

// DeployVersionCodeThenPoll performs DeployVersionCode then polls until it's completed
func (c EdgeActionVersionsClient) DeployVersionCodeThenPoll(ctx context.Context, id VersionId, input VersionCode) error {
	return c.DeployVersionCodeCallbackThenPoll(ctx, id, input, nil)
}

// DeployVersionCodeCallbackThenPoll performs DeployVersionCode, runs the optional callback function, then polls until it's completed
func (c EdgeActionVersionsClient) DeployVersionCodeCallbackThenPoll(ctx context.Context, id VersionId, input VersionCode, callback func() error) error {
	result, err := c.DeployVersionCode(ctx, id, input)
	if err != nil {
		return fmt.Errorf("performing DeployVersionCode: %+v", err)
	}

	if callback != nil {
		if err := callback(); err != nil {
			return fmt.Errorf("executing callback function: %+v", err)
		}
	}

	if err := result.Poller.PollUntilDone(ctx); err != nil {
		return fmt.Errorf("polling after DeployVersionCode: %+v", err)
	}

	return nil
}
