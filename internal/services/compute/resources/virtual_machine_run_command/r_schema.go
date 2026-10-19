// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package virtual_machine_run_command

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func (r Resource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.VirtualMachineRunCommandName,
		},

		"virtual_machine_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: commonids.ValidateVirtualMachineID,
		},

		"source": {
			Type:     pluginsdk.TypeList,
			Required: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"command_id": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
						ExactlyOneOf: []string{
							"source.0.command_id",
							"source.0.script",
							"source.0.script_uri",
						},
					},
					"script": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
						ExactlyOneOf: []string{
							"source.0.command_id",
							"source.0.script",
							"source.0.script_uri",
						},
					},
					"script_uri": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.IsURLWithHTTPS,
						ExactlyOneOf: []string{
							"source.0.command_id",
							"source.0.script",
							"source.0.script_uri",
						},
					},
					"script_uri_managed_identity": {
						Type:      pluginsdk.TypeList,
						Optional:  true,
						Sensitive: true,
						MaxItems:  1,
						RequiredWith: []string{
							"source.0.script_uri",
						},
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"client_id": {
									Type:         pluginsdk.TypeString,
									Optional:     true,
									Sensitive:    true,
									ValidateFunc: validation.StringIsNotEmpty,
									ConflictsWith: []string{
										"source.0.script_uri_managed_identity.0.object_id",
									},
								},
								"object_id": {
									Type:         pluginsdk.TypeString,
									Optional:     true,
									Sensitive:    true,
									ValidateFunc: validation.StringIsNotEmpty,
									ConflictsWith: []string{
										"source.0.script_uri_managed_identity.0.client_id",
									},
								},
							},
						},
					},
				},
			},
		},

		"error_blob_managed_identity": {
			Type:      pluginsdk.TypeList,
			Optional:  true,
			MaxItems:  1,
			Sensitive: true,
			RequiredWith: []string{
				"error_blob_uri",
			},
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"client_id": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Sensitive:    true,
						ValidateFunc: validation.IsUUID,
						ConflictsWith: []string{
							"error_blob_managed_identity.0.object_id",
						},
					},
					"object_id": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Sensitive:    true,
						ValidateFunc: validation.IsUUID,
						ConflictsWith: []string{
							"error_blob_managed_identity.0.client_id",
						},
					},
				},
			},
		},

		"error_blob_uri": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.IsURLWithHTTPS,
		},

		"output_blob_managed_identity": {
			Type:      pluginsdk.TypeList,
			Optional:  true,
			MaxItems:  1,
			Sensitive: true,
			RequiredWith: []string{
				"output_blob_uri",
			},
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"client_id": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Sensitive:    true,
						ValidateFunc: validation.StringIsNotEmpty,
						ConflictsWith: []string{
							"output_blob_managed_identity.0.object_id",
						},
					},
					"object_id": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Sensitive:    true,
						ValidateFunc: validation.StringIsNotEmpty,
						ConflictsWith: []string{
							"output_blob_managed_identity.0.client_id",
						},
					},
				},
			},
		},

		"output_blob_uri": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.IsURLWithHTTPS,
		},

		"parameter": {
			Type:     pluginsdk.TypeList,
			Optional: true,
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

		"protected_parameter": {
			Type:      pluginsdk.TypeList,
			Optional:  true,
			Sensitive: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						Sensitive:    true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"value": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						Sensitive:    true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},
			},
		},

		"run_as_password": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Sensitive:    true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"run_as_user": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"location": commonschema.Location(),

		"tags": commonschema.Tags(),
	}
}

func (r Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"instance_view": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"exit_code": {
						Type:     pluginsdk.TypeInt,
						Computed: true,
					},
					"execution_state": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
					"execution_message": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
					"output": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
					"error_message": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
					"start_time": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
					"end_time": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
				},
			},
		},
	}
}
