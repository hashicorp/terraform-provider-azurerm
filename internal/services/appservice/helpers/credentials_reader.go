// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package helpers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/web/2023-12-01/webapps"
)

func ListPublishingCredentials(ctx context.Context, client *webapps.WebAppsClient, id commonids.AppServiceId) (*webapps.User, error) {
	userModel := &webapps.User{}

	siteCredentials, err := client.ListPublishingCredentials(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("listing Site Publishing Credential information for %s: %+v", id, err)
	}

	// The credentials are regenerated at some point in the creation process, the initial response is not the final
	// value. The final result error is populated on success as a 404, so we're ignoring it here since this is a R/O
	// pair of properties
	_ = siteCredentials.Poller.PollUntilDone(ctx)
	siteCredentials, err = client.ListPublishingCredentials(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("listing Site Publishing Credential information for %s: %+v", id, err)
	}

	if siteCredentials.HttpResponse != nil {
		err = UnmarshalCredentialsResponse(siteCredentials.HttpResponse.Body, userModel)
		if err != nil {
			return nil, fmt.Errorf("could not decode Publishing Credential information for %s: %+v", id, err)
		}
	}

	return userModel, nil
}

func ListPublishingCredentialsSlot(ctx context.Context, client *webapps.WebAppsClient, id webapps.SlotId) (*webapps.User, error) {
	userModel := &webapps.User{}

	siteCredentials, err := client.ListPublishingCredentialsSlot(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("listing Site Publishing Credential information for %s: %+v", id, err)
	}
	// The credentials are regenerated at some point in the creation process, the initial response is not the final
	// value. The final result error is populated on success as a 404, so we're ignoring it here since this is a R/O
	// pair of properties
	_ = siteCredentials.Poller.PollUntilDone(ctx)
	siteCredentials, err = client.ListPublishingCredentialsSlot(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("listing Site Publishing Credential information for %s: %+v", id, err)
	}

	if siteCredentials.HttpResponse != nil {
		err = UnmarshalCredentialsResponse(siteCredentials.HttpResponse.Body, userModel)
		if err != nil {
			return nil, fmt.Errorf("could not decode Publishing Credential information for %s: %+v", id, err)
		}
	}

	return userModel, nil
}

func UnmarshalCredentialsResponse(r io.Reader, user *webapps.User) error {
	bytes, err := io.ReadAll(r)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(bytes, user); err != nil {
		return err
	}

	return nil
}
