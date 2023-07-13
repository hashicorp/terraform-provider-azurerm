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

	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/tombuildsstuff/kermit/sdk/web/2022-09-01/web"
)

const (
	zipDeployError    = 3
	zipDeployComplete = 4
)

func GetCredentialsAndPublish(ctx context.Context, client *web.AppsClient, resourceGroup string, siteName string, sourceFile string) error {
	site, err := client.Get(ctx, resourceGroup, siteName)
	if err != nil || site.SiteProperties == nil {
		return fmt.Errorf("reading site %s to perform zip deploy: %+v", siteName, err)
	}
	props := *site.SiteProperties
	if sslStates := props.HostNameSslStates; sslStates != nil {
		for _, v := range *sslStates {
			if v.Name != nil && *v.Name != "" && v.HostType == web.HostTypeRepository {
				user, passwd, err := GetSitePublishingCredentials(ctx, client, resourceGroup, siteName)
				if err != nil {
					return err
				}
				httpsHost := fmt.Sprintf("https://%s", *v.Name)

				if err := PublishZipDeployLocalFileKuduPush(ctx, httpsHost, *user, *passwd, client.UserAgent, sourceFile); err != nil {
					return fmt.Errorf("publishing source (%s) to site %s (Resource Group %s): %+v", sourceFile, siteName, resourceGroup, err)
				}

				continue
			}
		}
	} else {
		return fmt.Errorf("could not determine SCM Site name for Site %s (Resource Group %s) for Zip Deployment", siteName, resourceGroup)
	}

	return nil
}

func GetCredentialsAndPublishSlot(ctx context.Context, client *web.AppsClient, resourceGroup string, siteName string, sourceFile string, slotName string) error {
	site, err := client.GetSlot(ctx, resourceGroup, siteName, slotName)
	if err != nil || site.SiteProperties == nil {
		return fmt.Errorf("reading site %s to perform zip deploy: %+v", siteName, err)
	}
	props := *site.SiteProperties
	if sslStates := props.HostNameSslStates; sslStates != nil {
		for _, v := range *sslStates {
			if v.Name != nil && *v.Name != "" && v.HostType == web.HostTypeRepository {
				user, passwd, err := GetSitePublishingCredentialsSlot(ctx, client, resourceGroup, siteName, slotName)
				if err != nil {
					return err
				}
				httpsHost := fmt.Sprintf("https://%s", *v.Name)

				if err := PublishZipDeployLocalFileKuduPush(ctx, httpsHost, *user, *passwd, client.UserAgent, sourceFile); err != nil {
					return fmt.Errorf("publishing source (%s) to site %s (Resource Group %s): %+v", sourceFile, siteName, resourceGroup, err)
				}

				continue
			}
		}
	} else {
		return fmt.Errorf("could not determine SCM Site name for Site %s (Resource Group %s) for Zip Deployment", siteName, resourceGroup)
	}

	return nil
}

func GetSitePublishingCredentials(ctx context.Context, client *web.AppsClient, resourceGroup string, siteName string) (user *string, passwd *string, err error) {
	siteCredentialsFuture, err := client.ListPublishingCredentials(ctx, resourceGroup, siteName)
	if err != nil {
		return nil, nil, fmt.Errorf("listing Site Publishing Credential information for %s (Resource Group %s): %+v", siteName, resourceGroup, err)
	}
	if err := siteCredentialsFuture.WaitForCompletionRef(ctx, client.Client); err != nil {
		return nil, nil, fmt.Errorf("waiting for Site Publishing Credential information for %s (Resource Group %s): %+v", siteName, resourceGroup, err)
	}

	siteCredentials, err := siteCredentialsFuture.Result(*client)
	if err != nil {
		return nil, nil, fmt.Errorf("reading Site Publishing Credential information for %s (Resource Group %s): %+v", siteName, resourceGroup, err)
	}
	if siteCredentials.PublishingUserName == nil || siteCredentials.PublishingPassword == nil {
		return nil, nil, fmt.Errorf("site credentials for Site %s (Resource Group %s) were empty", siteName, resourceGroup)
	}
	return siteCredentials.PublishingUserName, siteCredentials.PublishingPassword, err
}

func GetSitePublishingCredentialsSlot(ctx context.Context, client *web.AppsClient, resourceGroup string, siteName string, slotName string) (user *string, passwd *string, err error) {
	siteCredentialsFuture, err := client.ListPublishingCredentialsSlot(ctx, resourceGroup, siteName, slotName)
	if err != nil {
		return nil, nil, fmt.Errorf("listing Site Publishing Credential information for %s (Resource Group %s): %+v", siteName, resourceGroup, err)
	}
	if err := siteCredentialsFuture.WaitForCompletionRef(ctx, client.Client); err != nil {
		return nil, nil, fmt.Errorf("waiting for Site Publishing Credential information for %s (Resource Group %s): %+v", siteName, resourceGroup, err)
	}

	siteCredentials, err := siteCredentialsFuture.Result(*client)
	if err != nil {
		return nil, nil, fmt.Errorf("reading Site Publishing Credential information for %s (Resource Group %s): %+v", siteName, resourceGroup, err)
	}
	if siteCredentials.PublishingUserName == nil || siteCredentials.PublishingPassword == nil {
		return nil, nil, fmt.Errorf("site credentials for Site %s (Resource Group %s) were empty", siteName, resourceGroup)
	}
	return siteCredentials.PublishingUserName, siteCredentials.PublishingPassword, err
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
