// Copyright IBM Corp. 2014, 2026
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
			Description:  "The name of this Virtual Machine Run Command.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.VirtualMachineRunCommandName,
		},

		"virtual_machine_id": {
			Description:  "The Virtual Machine ID within which this Virtual Machine Run Command should exist.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: commonids.ValidateVirtualMachineID,
		},

		"source": {
			Description: "The source of the run command script.",
			Type:        pluginsdk.TypeList,
			Required:    true,
			MaxItems:    1,
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
									Description:  "The client ID of the Managed Identity.",
									Type:         pluginsdk.TypeString,
									Optional:     true,
									Sensitive:    true,
									ValidateFunc: validation.StringIsNotEmpty,
									ConflictsWith: []string{
										"source.0.script_uri_managed_identity.0.object_id",
									},
								},
								"object_id": {
									Description:  "The object ID of the Managed Identity.",
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
			Description: "User-assigned Managed Identity that has access to errorBlobUri storage blob.",
			Type:        pluginsdk.TypeList,
			Optional:    true,
			MaxItems:    1,
			Sensitive:   true,
			RequiredWith: []string{
				"error_blob_uri",
			},
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"client_id": {
						Description:  "The client ID of the Managed Identity.",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Sensitive:    true,
						ValidateFunc: validation.IsUUID,
						ConflictsWith: []string{
							"error_blob_managed_identity.0.object_id",
						},
					},
					"object_id": {
						Description:  "The object ID of the Managed Identity.",
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
			Description:  "The Azure storage blob where script error stream will be uploaded.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.IsURLWithHTTPS,
		},

		"output_blob_managed_identity": {
			Description: "User-assigned Managed Identity that has access to outputBlobUri storage blob.",
			Type:        pluginsdk.TypeList,
			Optional:    true,
			MaxItems:    1,
			Sensitive:   true,
			RequiredWith: []string{
				"output_blob_uri",
			},
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"client_id": {
						Description:  "The client ID of the Managed Identity.",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Sensitive:    true,
						ValidateFunc: validation.StringIsNotEmpty,
						ConflictsWith: []string{
							"output_blob_managed_identity.0.object_id",
						},
					},
					"object_id": {
						Description:  "The object ID of the Managed Identity.",
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
			Description:  "The Azure storage blob where script output stream will be uploaded.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.IsURLWithHTTPS,
		},

		"parameter": {
			Description: "The parameters used by the script.",
			Type:        pluginsdk.TypeList,
			Optional:    true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Description:  "The run parameter name.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"value": {
						Description:  "The run parameter value.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},
			},
		},

		"protected_parameter": {
			Description: "The protected parameters used by the script.",
			Type:        pluginsdk.TypeList,
			Optional:    true,
			Sensitive:   true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Description:  "The run parameter name.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						Sensitive:    true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"value": {
						Description:  "The run parameter value.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						Sensitive:    true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},
			},
		},

		"run_as_password": {
			Description:  "The user account password on the VM when executing the Virtual Machine Run Command.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Sensitive:    true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"run_as_user": {
			Description:  "The user account on the VM when executing the Virtual Machine Run Command.",
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
