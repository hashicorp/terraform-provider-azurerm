package profiles

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

type AFDProfilesUpgradeOperationResponse struct {
	Poller       pollers.Poller
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *Profile
}

// AFDProfilesUpgrade ...
func (c ProfilesClient) AFDProfilesUpgrade(ctx context.Context, id ProfileId, input ProfileUpgradeParameters) (result AFDProfilesUpgradeOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusAccepted,
			http.StatusOK,
		},
		HttpMethod: http.MethodPost,
		Path:       fmt.Sprintf("%s/upgrade", id.ID()),
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

// AFDProfilesUpgradeThenPoll performs AFDProfilesUpgrade then polls until it's completed
func (c ProfilesClient) AFDProfilesUpgradeThenPoll(ctx context.Context, id ProfileId, input ProfileUpgradeParameters) error {
	return c.AFDProfilesUpgradeCallbackThenPoll(ctx, id, input, nil)
}

// AFDProfilesUpgradeCallbackThenPoll performs AFDProfilesUpgrade, runs the optional callback function, then polls until it's completed
func (c ProfilesClient) AFDProfilesUpgradeCallbackThenPoll(ctx context.Context, id ProfileId, input ProfileUpgradeParameters, callback func() error) error {
	result, err := c.AFDProfilesUpgrade(ctx, id, input)
	if err != nil {
		return fmt.Errorf("performing AFDProfilesUpgrade: %+v", err)
	}

	if callback != nil {
		if err := callback(); err != nil {
			return fmt.Errorf("executing callback function: %+v", err)
		}
	}

	if err := result.Poller.PollUntilDone(ctx); err != nil {
		return fmt.Errorf("polling after AFDProfilesUpgrade: %+v", err)
	}

	return nil
}
