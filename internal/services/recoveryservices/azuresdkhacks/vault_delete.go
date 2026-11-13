// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package azuresdkhacks

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"slices"
	"strings"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservices/2025-08-01/vaults"
)

// When the vault is deleted, the polling URL may also return 404. Tracked on https://github.com/Azure/azure-rest-api-specs/issues/47184.
func DeleteVaultThenPoll(ctx context.Context, vaultClient *vaults.VaultsClient, id vaults.VaultId) error {
	deleteClient := *vaultClient
	resourceManagerClient := *vaultClient.Client
	baseClient := *vaultClient.Client.Client
	deleteClient.Client = &resourceManagerClient
	resourceManagerClient.Client = &baseClient
	if baseClient.ResponseMiddlewares != nil {
		middlewares := slices.Clone(*baseClient.ResponseMiddlewares)
		baseClient.ResponseMiddlewares = &middlewares
	}

	result, err := deleteClient.Delete(ctx, id)
	if err != nil {
		return fmt.Errorf("performing Delete: %w", err)
	}

	pollingURL := result.HttpResponse.Header.Get("Azure-AsyncOperation")
	if pollingURL == "" {
		pollingURL = result.HttpResponse.Header.Get("Location")
	}
	if result.HttpResponse.StatusCode == http.StatusNoContent {
		pollingURL = ""
	}
	if pollingURL != "" {
		pollingURI, err := url.Parse(pollingURL)
		if err != nil {
			return fmt.Errorf("parsing the vault deletion polling URL: %w", err)
		}
		if strings.EqualFold(pollingURI.Path, result.HttpResponse.Request.URL.Path) {
			pollingURL = ""
		}
	}

	errVaultDeleted := errors.New("vault deletion confirmed")
	if pollingURL != "" {
		baseClient.AppendResponseMiddleware(func(_ *http.Request, resp *http.Response) (*http.Response, error) {
			req := resp.Request
			if resp.StatusCode != http.StatusNotFound || req == nil || req.URL == nil || req.Method != http.MethodGet || req.URL.String() != pollingURL {
				return resp, nil
			}

			vault, err := vaultClient.Get(req.Context(), id)
			if response.WasNotFound(vault.HttpResponse) {
				resp.Body.Close()
				return resp, errVaultDeleted
			}
			if err != nil {
				resp.Body.Close()
				return resp, fmt.Errorf("checking whether %s was deleted: %w", id, err)
			}
			return resp, nil
		})
	}

	if err := result.Poller.PollUntilDone(ctx); err != nil && !errors.Is(err, errVaultDeleted) {
		return fmt.Errorf("polling after Delete: %w", err)
	}
	return nil
}
