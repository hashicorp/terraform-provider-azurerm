// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package analysis_services_server

import (
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/analysisservices/2017-08-01/servers"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

//go:generate go run ../../../../tools/generator-tests resourceidentity

const analysisServicesServerResourceName = "azurerm_analysis_services_server"

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceAnalysisServicesServerCreate,
		Read:   resourceAnalysisServicesServerRead,
		Update: resourceAnalysisServicesServerUpdate,
		Delete: resourceAnalysisServicesServerDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingIdentity(&servers.ServerId{}),

		Identity: &schema.ResourceIdentity{
			SchemaFunc: pluginsdk.GenerateIdentitySchema(&servers.ServerId{}),
		},

		Schema: analysisServicesServerSchema(),
	}
}
