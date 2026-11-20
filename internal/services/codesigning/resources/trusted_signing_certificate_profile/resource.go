// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package trusted_signing_certificate_profile

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/codesigning/2025-10-13/certificateprofiles"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

//go:generate go run ../../../../tools/generator-tests resourceidentity -properties name -compare-values subscription_id:trusted_signing_account_id,resource_group_name:trusted_signing_account_id,code_signing_account_name:trusted_signing_account_id -test-env-vars ARM_TEST_TRUSTED_SIGNING_IDENTITY_ID

type Resource struct{}

var _ sdk.ResourceWithIdentity = Resource{}

func (Resource) ResourceType() string {
	return "azurerm_trusted_signing_certificate_profile"
}

func (Resource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return certificateprofiles.ValidateCertificateProfileID
}

func (Resource) Identity() resourceids.ResourceId {
	return &certificateprofiles.CertificateProfileId{}
}
