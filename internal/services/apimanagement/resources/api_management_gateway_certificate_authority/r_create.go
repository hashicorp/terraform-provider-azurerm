// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package api_management_gateway_certificate_authority

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/gatewaycertificateauthority"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/apimanagementservice"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	apimanagementClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceApiManagementGatewayCertificateAuthorityCreateUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := apimanagementClient.From(meta).GatewayCertificateAuthorityClient
	ctx, cancel := timeouts.ForCreateUpdate(account.StopContext(meta), d)
	defer cancel()

	apimId, err := apimanagementservice.ParseServiceID(d.Get("api_management_id").(string))
	if err != nil {
		return fmt.Errorf("parsing `api_management_id`: %v", err)
	}

	id := gatewaycertificateauthority.NewCertificateAuthorityID(apimId.SubscriptionId, apimId.ResourceGroupName, apimId.ServiceName, d.Get("gateway_name").(string), d.Get("certificate_name").(string))

	if d.IsNewResource() {
		if !features.From(meta).SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
			existing, err := client.Get(ctx, id)
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %s", id, err)
				}
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return tf.ImportAsExistsError("azurerm_api_management_gateway_certificate_authority", id.ID())
			}
		}
	}

	parameters := gatewaycertificateauthority.GatewayCertificateAuthorityContract{
		Properties: &gatewaycertificateauthority.GatewayCertificateAuthorityContractProperties{
			IsTrusted: pointer.To(d.Get("is_trusted").(bool)),
		},
	}

	if _, err = client.CreateOrUpdate(ctx, id, parameters, gatewaycertificateauthority.CreateOrUpdateOperationOptions{}); err != nil {
		return fmt.Errorf("creating or updating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceApiManagementGatewayCertificateAuthorityRead(d, meta)
}
