// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package helpers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/web/2023-12-01/webapps"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

const (
	zipDeployError    = 3
	zipDeployComplete = 4
)

func GetCredentialsAndPublish(ctx context.Context, client *webapps.WebAppsClient, appID commonids.AppServiceId, sourceFile string) error {
	site, err := client.Get(ctx, appID)
	if err != nil || site.Model == nil {
		return fmt.Errorf("reading site %s to perform zip deploy: %+v", appID.SiteName, err)
	}
	props := *site.Model.Properties
	if sslStates := props.HostNameSslStates; sslStates != nil {
		for _, v := range *sslStates {
			if v.Name != nil && *v.Name != "" && pointer.From(v.HostType) == webapps.HostTypeRepository {
				user, passwd, err := GetSitePublishingCredentials(ctx, client, appID)
				if err != nil {
					return err
				}
				httpsHost := fmt.Sprintf("https://%s", *v.Name)

				if err := PublishZipDeployLocalFileKuduPush(ctx, httpsHost, *user, *passwd, client.Client.UserAgent, sourceFile); err != nil {
					return fmt.Errorf("publishing source (%s) to site %s: %+v", sourceFile, appID, err)
				}

				continue
			}
		}
	} else {
		return fmt.Errorf("could not determine SCM Site name for Site %s for Zip Deployment", appID)
	}

	return nil
}

func GetCredentialsAndPublishSlot(ctx context.Context, client *webapps.WebAppsClient, id webapps.SlotId, sourceFile string) error {
	site, err := client.GetSlot(ctx, id)
	if err != nil || site.Model == nil || site.Model.Properties == nil {
		return fmt.Errorf("reading site %s to perform zip deploy: %+v", id.SiteName, err)
	}
	props := *site.Model.Properties
	if sslStates := props.HostNameSslStates; sslStates != nil {
		for _, v := range *sslStates {
			if v.Name != nil && *v.Name != "" && pointer.From(v.HostType) == webapps.HostTypeRepository {
				user, passwd, err := GetSitePublishingCredentialsSlot(ctx, client, id)
				if err != nil {
					return err
				}
				httpsHost := fmt.Sprintf("https://%s", *v.Name)

				if err := PublishZipDeployLocalFileKuduPush(ctx, httpsHost, *user, *passwd, client.Client.UserAgent, sourceFile); err != nil {
					return fmt.Errorf("publishing source (%s) to site %s: %+v", sourceFile, id, err)
				}

				continue
			}
		}
	} else {
		return fmt.Errorf("could not determine SCM Site name for Slot %s for Zip Deployment", id)
	}

	return nil
}

func GetSitePublishingCredentials(ctx context.Context, client *webapps.WebAppsClient, appID commonids.AppServiceId) (user *string, passwd *string, err error) {
	siteCredentials, err := ListPublishingCredentials(ctx, client, appID)
	if err != nil {
		return nil, nil, err
	}

	if siteCredentials.Properties != nil {
		return pointer.To(siteCredentials.Properties.PublishingUserName), siteCredentials.Properties.PublishingPassword, nil
	}
	return nil, nil, fmt.Errorf("could not decode Publishing Credential information for %s", appID)
}

func GetSitePublishingCredentialsSlot(ctx context.Context, client *webapps.WebAppsClient, id webapps.SlotId) (user *string, passwd *string, err error) {
	siteCredentials, err := ListPublishingCredentialsSlot(ctx, client, id)
	if err != nil {
		return nil, nil, err
	}

	if siteCredentials.Properties != nil {
		return pointer.To(siteCredentials.Properties.PublishingUserName), siteCredentials.Properties.PublishingPassword, nil
	}
	return nil, nil, fmt.Errorf("could not decode Publishing Credential information for %s", id)
}

func PublishZipDeployLocalFileKuduPush(ctx context.Context, host string, user string, passwd string, userAgent string, zipSource string) error {
	f, err := os.Open(zipSource)
	if err != nil {
		return err
	}

	publishEndpoint := fmt.Sprintf("%s/api/zipdeploy?isAsync=true", host)
	statusEndpoint := fmt.Sprintf("%s/api/deployments/latest", host)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, publishEndpoint, f)
	if err != nil {
		return fmt.Errorf("preparing publish request: %+v", err)
	}

	req.SetBasicAuth(user, passwd)
	req.Header["Cache-Control"] = []string{"no-cache"}
	req.Header["User-Agent"] = []string{userAgent}
	req.Header["Content-Type"] = []string{"application/octet-stream"}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("sending publish request: %+v", err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusAccepted {
		if resp.StatusCode == http.StatusConflict {
			return fmt.Errorf("publising Zip Deployment failed with %s - Another operation is in progress or your application is not configured for Zip deployments", resp.Status)
		}
		return fmt.Errorf("publishing failed with status code %s", resp.Status)
	}

	statusReq, err := http.NewRequestWithContext(ctx, http.MethodGet, statusEndpoint, http.NoBody)
	if err != nil {
		return err
	}

	statusReq.SetBasicAuth(user, passwd)
	deadline, ok := ctx.Deadline()
	if !ok {
		return fmt.Errorf("publish request context had no deadline")
	}

	deployWait := &pluginsdk.StateChangeConf{
		Pending:      []string{"pending"},
		Target:       []string{"complete"},
		PollInterval: 10 * time.Second,
		Delay:        10 * time.Second,
		Timeout:      time.Until(deadline),
		Refresh:      checkZipDeploymentStatusRefresh(statusReq),
	}

	if _, err := deployWait.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("waiting for Zip Deployment to complete")
	}

	return nil
}

func checkZipDeploymentStatusRefresh(r *http.Request) pluginsdk.StateRefreshFunc {
	return func() (interface{}, string, error) {
		resp, err := http.DefaultClient.Do(r)
		if err != nil {
			return nil, "", err
		}

		if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusAccepted {
			return nil, "", fmt.Errorf("failed to read Zip Deployment status: %s", resp.Status)
		}
		respBody, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, "", fmt.Errorf("reading status response body for Zip Deploy")
		}

		body := make(map[string]interface{})
		err = json.Unmarshal(respBody, &body)
		if err != nil {
			return nil, "", fmt.Errorf("could not parse status response for Zip Deploy")
		}

		if statusRaw, ok := body["status"]; ok && statusRaw != nil {
			if status, ok := statusRaw.(float64); ok {
				switch status {
				case zipDeployError:
					return nil, "", fmt.Errorf("zip deployment failed")
				case zipDeployComplete:
					return status, "complete", nil
				default:
					return status, "pending", nil
				}
			}
		}

		return nil, "", fmt.Errorf("could not determine status from deployment response")
	}
}
