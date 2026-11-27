// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package sentinel_threat_intelligence_indicator

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func (r Resource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"workspace_id": {
			Description:  "The ID of the Log Analytics Workspace.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"confidence": {
			Description:  "Confidence levels of the Threat Intelligence Indicator.",
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			Default:      -1, // set the default value to -1 to split `nil` and `0`.
			ValidateFunc: validation.IntBetween(0, 100),
		},

		"created_by": {
			Description:  "The creator of the Threat Intelligence Indicator.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"description": {
			Description:  "The description of the Threat Intelligence Indicator.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"display_name": {
			Description:  "The display name of the Threat Intelligence Indicator.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"extension": {
			Description: "The extension config of the Threat Intelligence Indicator in JSON format.",
			Type:        pluginsdk.TypeString,
			Optional:    true,
			// NOTE: O+C API sets this if omitted without issues for overwriting/reverting to default so this can remain
			Computed:         true,
			ValidateFunc:     validation.StringIsJSON,
			DiffSuppressFunc: pluginsdk.SuppressJsonDiff,
		},

		"external_reference": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"description": {
						Description:  "The description of the external reference of the Threat Intelligence Indicator.",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"id": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"hashes": {
						Description: "The list of hashes of the external reference of the Threat Intelligence Indicator.",
						Type:        pluginsdk.TypeMap,
						Optional:    true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
					},

					"source_name": {
						Description:  "The source name of the external reference of the Threat Intelligence Indicator.",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"url": {
						Description:  "The URL of the external reference of the Threat Intelligence Indicator.",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},
			},
		},

		"granular_marking": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"language": {
						Description:  "The language of granular marking of the Threat Intelligence Indicator.",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"marking_ref": {
						Description:  "The reference of the granular marking of the Threat Intelligence Indicator.",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"selectors": {
						Description: "A list of selectors of the granular marking of the Threat Intelligence Indicator.",
						Type:        pluginsdk.TypeList,
						Optional:    true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
					},
				},
			},
		},

		"kill_chain_phase": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Description:  "The name which should be used for the Lockheed Martin cyber kill chain phase.",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},
			},
		},

		"tags": {
			Description: "A list of tags of the Threat Intelligence Indicator.",
			Type:        pluginsdk.TypeList,
			Optional:    true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"language": {
			Description:  "The language of the Threat Intelligence Indicator.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"object_marking_refs": {
			Description: "A list of Threat Intelligence marking references.",
			Type:        pluginsdk.TypeList,
			Optional:    true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"pattern": {
			Description:  "The pattern used by the Threat Intelligence Indicator. When `pattern_type` set to `file`, `pattern` must be specified with `<HashName>:<Value>` format, such as `MD5:78ecc5c05cd8b79af480df2f8fba0b9d`.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"pattern_type": {
			Description: "The type of pattern used by the Threat Intelligence Indicator.",
			Type:        pluginsdk.TypeString,
			Required:    true,
			ValidateFunc: validation.StringInSlice([]string{
				string(PatternTypeDomainName),
				string(PatternTypeFile),
				string(PatternTypeIpV4Addr),
				string(PatternTypeIpV6Addr),
				string(PatternTypeUrl),
			}, false),
		},

		"pattern_version": {
			Description:  "The version of a Threat Intelligence entity.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"revoked": {
			Description: "Whether the Threat Intelligence entity revoked.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     false,
		},

		"source": {
			Description:  "Source of the Threat Intelligence Indicator.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"threat_types": {
			Description: "A list of threat types of this Threat Intelligence Indicator.",
			Type:        pluginsdk.TypeList,
			Optional:    true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"validate_from_utc": {
			Description:  "The start of validate date in RFC3339.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.IsRFC3339Time,
		},

		"validate_until_utc": {
			Description: "The end of validate date of the Threat Intelligence Indicator in RFC3339 format.",
			Type:        pluginsdk.TypeString,
			Optional:    true,
			// Note: O+C because the API returns a date around 1 year in the future if omitted from config
			Computed:     true,
			ValidateFunc: validation.IsRFC3339Time,
		},
	}
}

func (r Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"guid": {
			Description: "The GUID of this Sentinel Threat Intelligence Indicator.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"created_on": {
			Description: "The date of this Threat Intelligence Indicator created.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"defanged": {
			Description: "Whether the Threat Intelligence entity is defanged.",
			Type:        pluginsdk.TypeBool,
			Computed:    true,
		},

		"external_id": {
			Description: "The external ID of the Threat Intelligence Indicator.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"external_last_updated_time_utc": {
			Description: "The External last updated time in UTC.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"last_updated_time_utc": {
			Description: "The last updated time of the Threat Intelligence Indicator in UTC.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"parsed_pattern": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"pattern_type_key": {
						Description: "The type key of parsed pattern.",
						Type:        pluginsdk.TypeString,
						Computed:    true,
					},

					"pattern_type_values": {
						Type:     pluginsdk.TypeList,
						Computed: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"value": {
									Description: "The value of the parsed pattern type.",
									Type:        pluginsdk.TypeString,
									Computed:    true,
								},

								"value_type": {
									Description: "The type of the value of the parsed pattern type value.",
									Type:        pluginsdk.TypeString,
									Computed:    true,
								},
							},
						},
					},
				},
			},
		},

		"indicator_type": {
			Description: "A list of indicator types of this Threat Intelligence Indicator.",
			Type:        pluginsdk.TypeList,
			Computed:    true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},
	}
}
