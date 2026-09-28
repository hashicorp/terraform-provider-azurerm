// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package managedhsm

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/Azure/go-autorest/autorest"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type deleteAndPurgeNestedItem interface {
	DeleteNestedItem(ctx context.Context) (autorest.Response, error)
	NestedItemHasBeenDeleted(ctx context.Context) (autorest.Response, error)
	PurgeNestedItem(ctx context.Context) (autorest.Response, error)
	NestedItemHasBeenPurged(ctx context.Context) (autorest.Response, error)
}

func deleteAndOptionallyPurge(ctx context.Context, description string, shouldPurge bool, helper deleteAndPurgeNestedItem) error {
	timeout, ok := ctx.Deadline()
	if !ok {
		return fmt.Errorf("context is missing a timeout")
	}

	log.Printf("[DEBUG] Deleting %s..", description)
	if resp, err := helper.DeleteNestedItem(ctx); err != nil {
		if response.WasNotFound(resp.Response) {
			return nil
		}

		return fmt.Errorf("deleting %s: %+v", description, err)
	}
	log.Printf("[DEBUG] Waiting for %s to finish deleting..", description)
	stateConf := &pluginsdk.StateChangeConf{
		Pending: []string{"InProgress"},
		Target:  []string{"NotFound"},
		Refresh: func() (interface{}, string, error) {
			item, err := helper.NestedItemHasBeenDeleted(ctx)
			if err != nil {
				if response.WasNotFound(item.Response) {
					return item, "NotFound", nil
				}

				return nil, "Error", err
			}

			return item, "InProgress", nil
		},
		ContinuousTargetOccurence: 3,
		PollInterval:              5 * time.Second,
		Timeout:                   time.Until(timeout),
	}
	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("waiting for %s to be deleted: %+v", description, err)
	}
	log.Printf("[DEBUG] Deleted %s.", description)

	if !shouldPurge {
		log.Printf("[DEBUG] Skipping purging of %s as opted-out..", description)
		return nil
	}

	log.Printf("[DEBUG] Purging %s..", description)
	if err := pluginsdk.Retry(time.Until(timeout), func() *pluginsdk.RetryError {
		_, err := helper.PurgeNestedItem(ctx)
		if err == nil {
			return nil
		}
		if strings.Contains(err.Error(), "is currently being deleted") {
			return pluginsdk.RetryableError(fmt.Errorf("%s is currently being deleted, retrying", description))
		}
		return pluginsdk.NonRetryableError(fmt.Errorf("purging of %s : %+v", description, err))
	}); err != nil {
		return err
	}

	log.Printf("[DEBUG] Waiting for %s to finish purging..", description)
	stateConf = &pluginsdk.StateChangeConf{
		Pending: []string{"InProgress"},
		Target:  []string{"NotFound"},
		Refresh: func() (interface{}, string, error) {
			item, err := helper.NestedItemHasBeenPurged(ctx)
			if err != nil {
				if response.WasNotFound(item.Response) {
					return item, "NotFound", nil
				}

				return nil, "Error", err
			}

			return item, "InProgress", nil
		},
		ContinuousTargetOccurence: 3,
		PollInterval:              5 * time.Second,
		Timeout:                   time.Until(timeout),
	}
	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("waiting for %s to finish purging: %+v", description, err)
	}
	log.Printf("[DEBUG] Purged %s.", description)

	return nil
}

func managedHSMKeyRefreshFunc(ctx context.Context, childItemUri string) pluginsdk.StateRefreshFunc {
	return func() (interface{}, string, error) {
		log.Printf("[DEBUG] Checking to see if Managed HSM Key %q is available..", childItemUri)

		PTransport := &http.Transport{Proxy: http.ProxyFromEnvironment}

		client := &http.Client{
			Transport: PTransport,
		}

		req, err := http.NewRequestWithContext(ctx, http.MethodGet, childItemUri, nil)
		if err != nil {
			return nil, "pending", fmt.Errorf("building request to check Managed HSM Key at %q: %s", childItemUri, err)
		}

		conn, err := client.Do(req)
		if err != nil {
			log.Printf("[DEBUG] Didn't find Managed HSM Key at %q", childItemUri)
			return nil, "pending", fmt.Errorf("checking Managed HSM Key at %q: %s", childItemUri, err)
		}

		defer conn.Body.Close()

		log.Printf("[DEBUG] Found Managed HSM Key %q", childItemUri)
		return "available", "available", nil
	}
}
