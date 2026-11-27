// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package app_service_certificate

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/web/2023-12-01/certificates"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	webClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/web/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceAppServiceCertificateRead(d *pluginsdk.ResourceData, meta any) error {
	client := webClient.From(meta).CertificatesClient
	ctx, cancel := timeouts.ForRead(account.StopContext(meta), d)
	defer cancel()

	id, err := certificates.ParseCertificateID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.Set("name", id.CertificateName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		d.Set("location", location.Normalize(resp.Model.Location))

		if props := model.Properties; props != nil {
			d.Set("friendly_name", props.FriendlyName)
			d.Set("subject_name", props.SubjectName)
			d.Set("host_names", props.HostNames)
			d.Set("issuer", props.Issuer)
			d.Set("tags", tags.Flatten(model.Tags))
			d.Set("issue_date", props.IssueDate)
			d.Set("expiration_date", props.ExpirationDate)
			d.Set("thumbprint", props.Thumbprint)

			if props.HostingEnvironmentProfile != nil && props.HostingEnvironmentProfile.Id != nil {
				envId, err := commonids.ParseAppServiceEnvironmentID(*props.HostingEnvironmentProfile.Id)
				if err != nil {
					return err
				}
				d.Set("hosting_environment_profile_id", envId.ID())
			}

			if props.ServerFarmId != nil {
				sfID, err := commonids.ParseAppServicePlanID(*props.ServerFarmId)
				if err != nil {
					return err
				}
				d.Set("app_service_plan_id", sfID.ID())
			}
		}
	}

	return nil
}
