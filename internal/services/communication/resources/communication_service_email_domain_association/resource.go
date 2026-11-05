// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package communication_service_email_domain_association

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/communication/2026-03-18/communicationservices"
	"github.com/hashicorp/go-azure-sdk/resource-manager/communication/2026-03-18/domains"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

var _ sdk.Resource = Resource{}

type Resource struct{}

func (Resource) ResourceType() string {
	return "azurerm_communication_service_email_domain_association"
}

func (Resource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return func(input any, key string) (warnings []string, errors []error) {
		v, ok := input.(string)
		if !ok {
			errors = append(errors, fmt.Errorf("expected %q to be a string", key))
			return
		}

		if _, err := commonids.ParseCompositeResourceID(v, &communicationservices.CommunicationServiceId{}, &domains.DomainId{}); err != nil {
			errors = append(errors, err)
		}

		return
	}
}
