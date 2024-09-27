// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package resourcemanager

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/hashicorp/go-azure-sdk/sdk/auth"
)

func AuthorizeResourceManagerRequest(ctx context.Context, req *http.Request, authorizer auth.Authorizer) error {
	if req == nil {
		return fmt.Errorf("request was nil")
	}
	if authorizer == nil {
		return fmt.Errorf("authorizer was nil")
	}

	token, err := authorizer.Token(ctx, req)
	if err != nil {
		return err
	}

	if req.Header == nil {
		req.Header = make(http.Header)
	}

	req.Header.Set("Authorization", fmt.Sprintf("%s %s", token.Type(), token.AccessToken))

	auxTokens, err := authorizer.AuxiliaryTokens(ctx, req)
	if err != nil {
		return err
	}

	auxTokenValues := make([]string, 0)
	for _, auxToken := range auxTokens {
		auxTokenValues = append(auxTokenValues, fmt.Sprintf("%s %s", auxToken.Type(), auxToken.AccessToken))
	}
	if len(auxTokenValues) > 0 {
		req.Header.Set("X-Ms-Authorization-Auxiliary", strings.Join(auxTokenValues, ", "))
	}

	return nil
}
