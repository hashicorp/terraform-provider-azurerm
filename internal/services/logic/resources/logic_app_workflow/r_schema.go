// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package logic_app_workflow

import (
	"regexp"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/logic/2019-05-01/integrationaccounts"
	"github.com/hashicorp/go-azure-sdk/resource-manager/logic/2019-05-01/integrationserviceenvironments"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func logicAppWorkflowSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.All(
				validation.StringIsNotEmpty,
				validation.StringMatch(
					regexp.MustCompile("^[-()_.A-Za-z0-9]{1,80}$"),
					"The Logic app name can contain only letters, numbers, periods (.), hyphens (-), brackets (()) and underscores (_), up to 80 characters",
				),
			),
		},

		"location": commonschema.Location(),

		"resource_group_name": commonschema.ResourceGroupName(),

		"integration_service_environment_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: integrationserviceenvironments.ValidateIntegrationServiceEnvironmentID,
		},

		"access_control": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"action": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"allowed_caller_ip_address_range": {
									Type:     pluginsdk.TypeSet,
									Required: true,
									Elem: &pluginsdk.Schema{
										Type: pluginsdk.TypeString,
										ValidateFunc: validation.Any(
											validation.IsCIDR,
											validation.IsIPv4Range,
										),
									},
								},
							},
						},
					},

					"content": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"allowed_caller_ip_address_range": {
									Type:     pluginsdk.TypeSet,
									Required: true,
									Elem: &pluginsdk.Schema{
										Type: pluginsdk.TypeString,
										ValidateFunc: validation.Any(
											validation.IsCIDR,
											validation.IsIPv4Range,
										),
									},
								},
							},
						},
					},

					"trigger": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"allowed_caller_ip_address_range": {
									Type:     pluginsdk.TypeSet,
									Optional: true,
									Elem: &pluginsdk.Schema{
										Type: pluginsdk.TypeString,
										ValidateFunc: validation.Any(
											validation.IsCIDR,
											validation.IsIPv4Range,
										),
									},
								},

								"open_authentication_policy": {
									Type:     pluginsdk.TypeSet,
									Optional: true,
									Elem: &pluginsdk.Resource{
										Schema: map[string]*pluginsdk.Schema{
											"name": {
												Type:         pluginsdk.TypeString,
												Required:     true,
												ValidateFunc: validation.StringIsNotEmpty,
											},

											"claim": {
												Type:     pluginsdk.TypeSet,
												Required: true,
												Elem: &pluginsdk.Resource{
													Schema: map[string]*pluginsdk.Schema{
														"name": {
															Type:         pluginsdk.TypeString,
															Required:     true,
															ValidateFunc: validation.StringIsNotEmpty,
														},

														"value": {
															Type:         pluginsdk.TypeString,
															Required:     true,
															ValidateFunc: validation.StringIsNotEmpty,
														},
													},
												},
											},
										},
									},
								},
							},
						},
					},

					"workflow_management": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"allowed_caller_ip_address_range": {
									Type:     pluginsdk.TypeSet,
									Required: true,
									Elem: &pluginsdk.Schema{
										Type: pluginsdk.TypeString,
										ValidateFunc: validation.Any(
											validation.IsCIDR,
											validation.IsIPv4Range,
										),
									},
								},
							},
						},
					},
				},
			},
		},

		"identity": commonschema.SystemOrUserAssignedIdentityOptional(),

		"logic_app_integration_account_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: integrationaccounts.ValidateIntegrationAccountID,
		},

		// TODO: should Parameters be split out into their own object to allow validation on the different sub-types?
		"parameters": {
			Type:     pluginsdk.TypeMap,
			Optional: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  true,
		},

		"workflow_schema": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			ForceNew: true,
			Default:  "https://schema.management.azure.com/providers/Microsoft.Logic/schemas/2016-06-01/workflowdefinition.json#",
		},

		"workflow_version": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			ForceNew: true,
			Default:  "1.0.0.0",
		},

		"workflow_parameters": {
			Type:     pluginsdk.TypeMap,
			Optional: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"access_endpoint": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"connector_endpoint_ip_addresses": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem:     &pluginsdk.Schema{Type: pluginsdk.TypeString},
		},
		"connector_outbound_ip_addresses": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem:     &pluginsdk.Schema{Type: pluginsdk.TypeString},
		},
		"workflow_endpoint_ip_addresses": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem:     &pluginsdk.Schema{Type: pluginsdk.TypeString},
		},
		"workflow_outbound_ip_addresses": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem:     &pluginsdk.Schema{Type: pluginsdk.TypeString},
		},

		"tags": commonschema.Tags(),
	}
}
