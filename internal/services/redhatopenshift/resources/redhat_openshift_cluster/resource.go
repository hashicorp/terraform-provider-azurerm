// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package redhat_openshift_cluster

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/redhatopenshift/2025-07-25/openshiftclusters"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

//go:generate go run ../../../../tools/generator-tests resourceidentity -test-expect-non-empty

var (
	_ sdk.ResourceWithCustomizeDiff = Resource{}
	_ sdk.ResourceWithIdentity      = Resource{}
	_ sdk.ResourceWithUpdate        = Resource{}
)

type Resource struct{}

func (r Resource) ResourceType() string {
	return "azurerm_redhat_openshift_cluster"
}

func (r Resource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return openshiftclusters.ValidateOpenShiftClusterID
}

func (r Resource) Identity() resourceids.ResourceId {
	return &openshiftclusters.OpenShiftClusterId{}
}
