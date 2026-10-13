// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package api_management_backend

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/backend"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceApiManagementBackendCreateUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).ApiManagement.BackendClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := backend.NewBackendID(subscriptionId, d.Get("resource_group_name").(string), d.Get("api_management_name").(string), d.Get("name").(string))

	if d.IsNewResource() {
		if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
			existing, err := client.Get(ctx, id)
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %s", id, err)
				}
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return tf.ImportAsExistsError("azurerm_api_management_backend", id.ID())
			}
		}
	}

	credentialsRaw := d.Get("credentials").([]any)
	credentials := expandApiManagementBackendCredentials(credentialsRaw)
	protocol := d.Get("protocol").(string)
	proxyRaw := d.Get("proxy").([]any)
	proxy := expandApiManagementBackendProxy(proxyRaw)
	tlsRaw := d.Get("tls").([]any)
	tls := expandApiManagementBackendTls(tlsRaw)
	url := d.Get("url").(string)

	backendContract := backend.BackendContract{
		Properties: &backend.BackendContractProperties{
			Credentials: credentials,
			Protocol:    pointer.ToEnum[backend.BackendProtocol](protocol),
			Proxy:       proxy,
			Tls:         tls,
			Url:         pointer.To(url),
		},
	}
	if v, ok := d.GetOk("circuit_breaker_rule"); ok {
		backendContract.Properties.CircuitBreaker = expandApiManagementBackendCircuitBreaker(v.([]any))
	}
	if description, ok := d.GetOk("description"); ok {
		backendContract.Properties.Description = pointer.To(description.(string))
	}
	if resourceID, ok := d.GetOk("resource_id"); ok {
		backendContract.Properties.ResourceId = pointer.To(resourceID.(string))
	}
	if title, ok := d.GetOk("title"); ok {
		backendContract.Properties.Title = pointer.To(title.(string))
	}

	if serviceFabricClusterRaw, ok := d.GetOk("service_fabric_cluster"); ok {
		serviceFabricCluster, err := expandApiManagementBackendServiceFabricCluster(serviceFabricClusterRaw.([]any))
		if err != nil {
			return err
		}
		backendContract.Properties.Properties = &backend.BackendProperties{
			ServiceFabricCluster: serviceFabricCluster,
		}
	}

	if _, err := client.CreateOrUpdate(ctx, id, backendContract, backend.CreateOrUpdateOperationOptions{}); err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceApiManagementBackendRead(d, meta)
}
