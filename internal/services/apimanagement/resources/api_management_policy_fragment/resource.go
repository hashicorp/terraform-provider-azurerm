// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package api_management_policy_fragment

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/policyfragment"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceApiManagementPolicyFragmentCreate,
		Read:   resourceApiManagementPolicyFragmentRead,
		Update: resourceApiManagementPolicyFragmentUpdate,
		Delete: resourceApiManagementPolicyFragmentDelete,
		Importer: pluginsdk.ImporterValidatingResourceIdThen(func(id string) error {
			_, err := policyfragment.ParsePolicyFragmentID(id)
			return err
		}, func(ctx context.Context, d *pluginsdk.ResourceData, meta any) ([]*pluginsdk.ResourceData, error) {
			client := meta.(*clients.Client).ApiManagement.PolicyFragmentClient

			id, err := policyfragment.ParsePolicyFragmentID(d.Id())
			if err != nil {
				return nil, err
			}

			resp, err := client.Get(ctx, *id, policyfragment.GetOperationOptions{
				Format: pointer.To(policyfragment.PolicyFragmentContentFormatXml),
			})
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return nil, fmt.Errorf("%s was not found in Api Management instance %q in Resource Group %q", id, id.ServiceName, id.ResourceGroupName)
				}

				return nil, fmt.Errorf("retrieving %s: %+v", id, err)
			}

			d.Set("format", string(policyfragment.PolicyFragmentContentFormatXml))

			return []*pluginsdk.ResourceData{d}, nil
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: apiManagementPolicyFragmentSchema(),
	}
}
