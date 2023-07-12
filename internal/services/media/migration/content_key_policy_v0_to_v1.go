// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package migration

import (
	"context"
	"log"

	"github.com/hashicorp/go-azure-sdk/resource-manager/media/2022-08-01/contentkeypolicies"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

var _ pluginsdk.StateUpgrade = ContentKeyPolicyV0ToV1{}

type ContentKeyPolicyV0ToV1 struct{}

func (ContentKeyPolicyV0ToV1) Schema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"resource_group_name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"media_services_account_name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"description": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},

		"policy_option": {
			Type:     pluginsdk.TypeSet,
			Required: true,
			MinItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},

					"clear_key_configuration_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
					},

					"widevine_configuration_template": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},

					//lintignore:XS003
					"playready_configuration_license": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"allow_test_devices": {
									Type:     pluginsdk.TypeBool,
									Optional: true,
								},

								"begin_date": {
									Type:     pluginsdk.TypeString,
									Optional: true,
								},

								"content_key_location_from_header_enabled": {
									Type:     pluginsdk.TypeBool,
									Optional: true,
								},

								"content_key_location_from_key_id": {
									Type:     pluginsdk.TypeString,
									Optional: true,
								},

								"content_type": {
									Type:     pluginsdk.TypeString,
									Optional: true,
								},

								"expiration_date": {
									Type:     pluginsdk.TypeString,
									Optional: true,
								},

								"grace_period": {
									Type:      pluginsdk.TypeString,
									Optional:  true,
									Sensitive: true,
								},

								"license_type": {
									Type:     pluginsdk.TypeString,
									Optional: true,
								},

								//lintignore:XS003
								"play_right": {
									Type:     pluginsdk.TypeList,
									Optional: true,
									MaxItems: 1,
									Elem: &pluginsdk.Resource{
										Schema: map[string]*pluginsdk.Schema{
											"agc_and_color_stripe_restriction": {
												Type:     pluginsdk.TypeInt,
												Optional: true,
											},

											"allow_passing_video_content_to_unknown_output": {
												Type:     pluginsdk.TypeString,
												Optional: true,
											},

											"analog_video_opl": {
												Type:     pluginsdk.TypeInt,
												Optional: true,
											},

											"compressed_digital_audio_opl": {
												Type:     pluginsdk.TypeInt,
												Optional: true,
											},

											"digital_video_only_content_restriction": {
												Type:     pluginsdk.TypeBool,
												Optional: true,
											},

											"first_play_expiration": {
												Type:     pluginsdk.TypeString,
												Optional: true,
											},

											"image_constraint_for_analog_component_video_restriction": {
												Type:     pluginsdk.TypeBool,
												Optional: true,
											},

											"image_constraint_for_analog_computer_monitor_restriction": {
												Type:     pluginsdk.TypeBool,
												Optional: true,
											},

											"scms_restriction": {
												Type:     pluginsdk.TypeInt,
												Optional: true,
											},

											"uncompressed_digital_audio_opl": {
												Type:     pluginsdk.TypeInt,
												Optional: true,
											},

											"uncompressed_digital_video_opl": {
												Type:     pluginsdk.TypeInt,
												Optional: true,
											},
										},
									},
								},
								"relative_begin_date": {
									Type:     pluginsdk.TypeString,
									Optional: true,
								},

								"relative_expiration_date": {
									Type:     pluginsdk.TypeString,
									Optional: true,
								},
							},
						},
					},
					//lintignore:XS003
					"fairplay_configuration": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"ask": {
									Type:      pluginsdk.TypeString,
									Optional:  true,
									Sensitive: true,
								},
								"pfx": {
									Type:      pluginsdk.TypeString,
									Optional:  true,
									Sensitive: true,
								},
								"pfx_password": {
									Type:      pluginsdk.TypeString,
									Optional:  true,
									Sensitive: true,
								},
								//lintignore:XS003
								"offline_rental_configuration": {
									Type:     pluginsdk.TypeList,
									Optional: true,
									MaxItems: 1,
									Elem: &pluginsdk.Resource{
										Schema: map[string]*pluginsdk.Schema{
											"playback_duration_seconds": {
												Type:     pluginsdk.TypeInt,
												Optional: true,
											},
											"storage_duration_seconds": {
												Type:     pluginsdk.TypeInt,
												Optional: true,
											},
										},
									},
								},
								"rental_and_lease_key_type": {
									Type:     pluginsdk.TypeString,
									Optional: true,
								},
								"rental_duration_seconds": {
									Type:     pluginsdk.TypeInt,
									Optional: true,
								},
							},
						},
					},
					//lintignore:XS003
					"token_restriction": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"audience": {
									Type:     pluginsdk.TypeString,
									Optional: true,
								},
								"issuer": {
									Type:     pluginsdk.TypeString,
									Optional: true,
								},
								"token_type": {
									Type:     pluginsdk.TypeString,
									Optional: true,
								},
								"primary_symmetric_token_key": {
									Type:      pluginsdk.TypeString,
									Optional:  true,
									Sensitive: true,
								},
								"primary_rsa_token_key_exponent": {
									Type:      pluginsdk.TypeString,
									Optional:  true,
									Sensitive: true,
								},
								"primary_rsa_token_key_modulus": {
									Type:      pluginsdk.TypeString,
									Optional:  true,
									Sensitive: true,
								},
								"primary_x509_token_key_raw": {
									Type:      pluginsdk.TypeString,
									Optional:  true,
									Sensitive: true,
								},
								"open_id_connect_discovery_document": {
									Type:     pluginsdk.TypeString,
									Optional: true,
								},
								//lintignore:XS003
								"required_claim": {
									Type:     pluginsdk.TypeList,
									Optional: true,
									Elem: &pluginsdk.Resource{
										Schema: map[string]*pluginsdk.Schema{
											"type": {
												Type:     pluginsdk.TypeString,
												Optional: true,
											},
											"value": {
												Type:     pluginsdk.TypeString,
												Optional: true,
											},
										},
									},
								},
							},
						},
					},
					"open_restriction_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
					},
				},
			},
		},
	}
}

func (ContentKeyPolicyV0ToV1) UpgradeFunc() pluginsdk.StateUpgraderFunc {
	return func(ctx context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
		oldIdRaw := rawState["id"].(string)
		oldId, err := contentkeypolicies.ParseContentKeyPolicyIDInsensitively(oldIdRaw)
		if err != nil {
			return nil, err
		}

		newId := oldId.ID()
		log.Printf("[DEBUG] Updating ID from %q to %q..", oldIdRaw, newId)
		rawState["id"] = newId

		return rawState, nil
	}
}
