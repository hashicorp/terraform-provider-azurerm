// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package media

import (
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/media/2022-07-01/encodings"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/media/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceMediaTransform() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceMediaTransformCreateUpdate,
		Read:   resourceMediaTransformRead,
		Update: resourceMediaTransformCreateUpdate,
		Delete: resourceMediaTransformDelete,

		DeprecationMessage: azureMediaRetirementMessage,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := encodings.ParseTransformID(id)
			return err
		}),

		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.TransformV0ToV1{},
		}),
		SchemaVersion: 1,

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringMatch(
					regexp.MustCompile("^[-a-zA-Z0-9(_)]{1,128}$"),
					"Transform name must be 1 - 128 characters long, can contain letters, numbers, underscores, and hyphens (but the first and last character must be a letter or number).",
				),
			},

			"resource_group_name": commonschema.ResourceGroupName(),

			"media_services_account_name": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringMatch(
					regexp.MustCompile("^[-a-z0-9]{3,24}$"),
					"Media Services Account name must be 3 - 24 characters long, contain only lowercase letters and numbers.",
				),
			},

			"description": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			// lintignore:XS003
			"output": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				MinItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"on_error_action": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							Default:      string(encodings.OnErrorTypeStopProcessingJob),
							ValidateFunc: validation.StringInSlice(encodings.PossibleValuesForOnErrorType(), false),
						},
						// lintignore:XS003
						"builtin_preset": {
							Type:     pluginsdk.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"preset_name": {
										Type:         pluginsdk.TypeString,
										Required:     true,
										ValidateFunc: validation.StringInSlice(encodings.PossibleValuesForEncoderNamedPreset(), false),
									},
									"preset_configuration": {
										Type:     pluginsdk.TypeList,
										Optional: true,
										MaxItems: 1,
										Elem: &pluginsdk.Resource{
											Schema: map[string]*pluginsdk.Schema{
												"complexity": {
													Type:         pluginsdk.TypeString,
													Optional:     true,
													ValidateFunc: validation.StringInSlice(encodings.PossibleValuesForComplexity(), false),
												},
												"interleave_output": {
													Type:         pluginsdk.TypeString,
													Optional:     true,
													ValidateFunc: validation.StringInSlice(encodings.PossibleValuesForInterleaveOutput(), false),
												},
												"key_frame_interval_in_seconds": {
													Type:         pluginsdk.TypeFloat,
													Optional:     true,
													ValidateFunc: validation.FloatAtLeast(0),
												},
												"max_bitrate_bps": {
													Type:         pluginsdk.TypeInt,
													Optional:     true,
													ValidateFunc: validation.IntAtLeast(0),
												},
												"max_height": {
													Type:         pluginsdk.TypeInt,
													Optional:     true,
													ValidateFunc: validation.IntAtLeast(0),
												},
												"max_layers": {
													Type:         pluginsdk.TypeInt,
													Optional:     true,
													ValidateFunc: validation.IntAtLeast(0),
												},
												"min_bitrate_bps": {
													Type:         pluginsdk.TypeInt,
													Optional:     true,
													ValidateFunc: validation.IntAtLeast(0),
												},
												"min_height": {
													Type:         pluginsdk.TypeInt,
													Optional:     true,
													ValidateFunc: validation.IntAtLeast(0),
												},
											},
										},
									},
								},
							},
						},
						// lintignore:XS003
						"audio_analyzer_preset": {
							Type:     pluginsdk.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									// https://go.microsoft.com/fwlink/?linkid=2109463
									"audio_language": {
										Type:         pluginsdk.TypeString,
										Optional:     true,
										ValidateFunc: validation.StringIsNotEmpty,
									},
									"audio_analysis_mode": {
										Type:         pluginsdk.TypeString,
										Optional:     true,
										Default:      string(encodings.AudioAnalysisModeStandard),
										ValidateFunc: validation.StringInSlice(encodings.PossibleValuesForAudioAnalysisMode(), false),
									},
									"experimental_options": {
										Type:     pluginsdk.TypeMap,
										Optional: true,
										Elem: &pluginsdk.Schema{
											Type: pluginsdk.TypeString,
										},
									},
								},
							},
						},
						// lintignore:XS003
						"video_analyzer_preset": {
							Type:     pluginsdk.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									// https://go.microsoft.com/fwlink/?linkid=2109463
									"audio_language": {
										Type:         pluginsdk.TypeString,
										Optional:     true,
										ValidateFunc: validation.StringIsNotEmpty,
									},
									"audio_analysis_mode": {
										Type:         pluginsdk.TypeString,
										Optional:     true,
										Default:      string(encodings.AudioAnalysisModeStandard),
										ValidateFunc: validation.StringInSlice(encodings.PossibleValuesForAudioAnalysisMode(), false),
									},
									"insights_type": {
										Type:         pluginsdk.TypeString,
										Optional:     true,
										Default:      string(encodings.InsightsTypeAllInsights),
										ValidateFunc: validation.StringInSlice(encodings.PossibleValuesForInsightsType(), false),
									},
									"experimental_options": {
										Type:     pluginsdk.TypeMap,
										Optional: true,
										Elem: &pluginsdk.Schema{
											Type: pluginsdk.TypeString,
										},
									},
								},
							},
						},
						// lintignore:XS003
						"face_detector_preset": {
							Type:     pluginsdk.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"analysis_resolution": {
										Type:         pluginsdk.TypeString,
										Optional:     true,
										Default:      string(encodings.AnalysisResolutionSourceResolution),
										ValidateFunc: validation.StringInSlice(encodings.PossibleValuesForAnalysisResolution(), false),
									},
									"blur_type": {
										Type:         pluginsdk.TypeString,
										Optional:     true,
										ValidateFunc: validation.StringInSlice(encodings.PossibleValuesForBlurType(), false),
									},
									"experimental_options": {
										Type:     pluginsdk.TypeMap,
										Optional: true,
										Elem: &pluginsdk.Schema{
											Type: pluginsdk.TypeString,
										},
									},
									"face_redactor_mode": {
										Type:         pluginsdk.TypeString,
										Optional:     true,
										Default:      string(encodings.FaceRedactorModeAnalyze),
										ValidateFunc: validation.StringInSlice(encodings.PossibleValuesForFaceRedactorMode(), false),
									},
								},
							},
						},
						"custom_preset": {
							Type:     pluginsdk.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"codec": {
										Type:     pluginsdk.TypeList,
										Required: true,
										MinItems: 1,
										Elem: &pluginsdk.Resource{
											Schema: map[string]*schema.Schema{
												"aac_audio": {
													Type:     pluginsdk.TypeList,
													Optional: true,
													MaxItems: 1,
													Elem: &pluginsdk.Resource{
														Schema: map[string]*schema.Schema{
															"bitrate": {
																Type:         pluginsdk.TypeInt,
																Optional:     true,
																Default:      128000,
																ValidateFunc: validation.IntAtLeast(1),
															},
															"channels": {
																Type:         pluginsdk.TypeInt,
																Optional:     true,
																Default:      2,
																ValidateFunc: validation.IntBetween(1, 6),
															},
															"label": {
																Type:         pluginsdk.TypeString,
																Optional:     true,
																ValidateFunc: validation.StringIsNotEmpty,
															},
															"profile": {
																Type:         pluginsdk.TypeString,
																Optional:     true,
																Default:      string(encodings.AacAudioProfileAacLc),
																ValidateFunc: validation.StringInSlice(encodings.PossibleValuesForAacAudioProfile(), false),
															},
															"sampling_rate": {
																Type:         pluginsdk.TypeInt,
																Optional:     true,
																Default:      48000,
																ValidateFunc: validation.IntBetween(11025, 96000),
															},
														},
													},
												},
												"copy_audio": {
													Type:     pluginsdk.TypeList,
													Optional: true,
													MaxItems: 1,
													Elem: &pluginsdk.Resource{
														Schema: map[string]*schema.Schema{
															"label": {
																Type:         pluginsdk.TypeString,
																Optional:     true,
																ValidateFunc: validation.StringIsNotEmpty,
															},
														},
													},
												},
												"copy_video": {
													Type:     pluginsdk.TypeList,
													Optional: true,
													MaxItems: 1,
													Elem: &pluginsdk.Resource{
														Schema: map[string]*schema.Schema{
															"label": {
																Type:         pluginsdk.TypeString,
																Optional:     true,
																ValidateFunc: validation.StringIsNotEmpty,
															},
														},
													},
												},
												"dd_audio": {
													Type:     pluginsdk.TypeList,
													Optional: true,
													MaxItems: 1,
													Elem: &pluginsdk.Resource{
														Schema: map[string]*schema.Schema{
															"bitrate": {
																Type:         pluginsdk.TypeInt,
																Optional:     true,
																Default:      192000,
																ValidateFunc: validation.IntAtLeast(1),
															},
															"channels": {
																Type:         pluginsdk.TypeInt,
																Default:      2,
																Optional:     true,
																ValidateFunc: validation.IntBetween(1, 6),
															},
															"label": {
																Type:         pluginsdk.TypeString,
																Optional:     true,
																ValidateFunc: validation.StringIsNotEmpty,
															},
															"sampling_rate": {
																Type:         pluginsdk.TypeInt,
																Default:      48000,
																Optional:     true,
																ValidateFunc: validation.IntBetween(32000, 48000),
															},
														},
													},
												},
												"h264_video": {
													Type:     pluginsdk.TypeList,
													Optional: true,
													MaxItems: 1,
													Elem: &pluginsdk.Resource{
														Schema: map[string]*schema.Schema{
															"complexity": {
																Type:         pluginsdk.TypeString,
																Optional:     true,
																Default:      string(encodings.H264ComplexityBalanced),
																ValidateFunc: validation.StringInSlice(encodings.PossibleValuesForH264Complexity(), false),
															},
															"key_frame_interval": {
																Type:         pluginsdk.TypeString,
																Optional:     true,
																Default:      "PT2S",
																ValidateFunc: validate.ISO8601DurationBetween("PT0.5S", "PT20S"),
															},
															"label": {
																Type:         pluginsdk.TypeString,
																Optional:     true,
																ValidateFunc: validation.StringIsNotEmpty,
															},
															"layer": {
																Type:     pluginsdk.TypeList,
																Optional: true,
																MinItems: 1,
																Elem: &pluginsdk.Resource{
																	Schema: map[string]*schema.Schema{
																		"bitrate": {
																			Type:         pluginsdk.TypeInt,
																			Required:     true,
																			ValidateFunc: validation.IntAtLeast(1),
																		},
																		"adaptive_b_frame_enabled": {
																			Type:     pluginsdk.TypeBool,
																			Optional: true,
																			Default:  true,
																		},
																		"b_frames": {
																			Type:         pluginsdk.TypeInt,
																			Optional:     true,
																			Computed:     true,
																			ValidateFunc: validation.IntAtLeast(0),
																		},
																		"buffer_window": {
																			Type:         pluginsdk.TypeString,
																			Optional:     true,
																			Default:      "PT5S",
																			ValidateFunc: validate.ISO8601DurationBetween("PT0.1S", "PT100S"),
																		},
																		"crf": {
																			Type:         pluginsdk.TypeFloat,
																			Optional:     true,
																			Default:      23,
																			ValidateFunc: validation.FloatBetween(0, 51),
																		},
																		"entropy_mode": {
																			Type:         pluginsdk.TypeString,
																			Optional:     true,
																			Computed:     true,
																			ValidateFunc: validation.StringInSlice(encodings.PossibleValuesForEntropyMode(), false),
																		},
																		"frame_rate": {
																			Type:         pluginsdk.TypeString,
																			Optional:     true,
																			ValidateFunc: validation.StringIsNotEmpty,
																		},
																		"height": {
																			Type:         pluginsdk.TypeString,
																			Optional:     true,
																			Computed:     true,
																			ValidateFunc: validation.StringIsNotEmpty,
																		},
																		"label": {
																			Type:         pluginsdk.TypeString,
																			Optional:     true,
																			ValidateFunc: validation.StringIsNotEmpty,
																		},
																		"level": {
																			Type:         pluginsdk.TypeString,
																			Optional:     true,
																			Default:      "auto",
																			ValidateFunc: validation.StringIsNotEmpty,
																		},
																		"max_bitrate": {
																			Type:     pluginsdk.TypeInt,
																			Optional: true,
																			Computed: true,
																		},
																		"profile": {
																			Type:         pluginsdk.TypeString,
																			Optional:     true,
																			Default:      string(encodings.H264VideoProfileAuto),
																			ValidateFunc: validation.StringInSlice(encodings.PossibleValuesForH264VideoProfile(), false),
																		},
																		"reference_frames": {
																			Type:         pluginsdk.TypeInt,
																			Optional:     true,
																			Computed:     true,
																			ValidateFunc: validation.IntAtLeast(0),
																		},
																		"slices": {
																			Type:         pluginsdk.TypeInt,
																			Optional:     true,
																			Computed:     true,
																			ValidateFunc: validation.IntAtLeast(0),
																		},
																		"width": {
																			Type:         pluginsdk.TypeString,
																			Optional:     true,
																			Computed:     true,
																			ValidateFunc: validation.StringIsNotEmpty,
																		},
																	},
																},
															},
															"rate_control_mode": {
																Type:         pluginsdk.TypeString,
																Optional:     true,
																Default:      string(encodings.H264RateControlModeABR),
																ValidateFunc: validation.StringInSlice(encodings.PossibleValuesForH264RateControlMode(), false),
															},
															"scene_change_detection_enabled": {
																Type:     pluginsdk.TypeBool,
																Optional: true,
																Default:  false,
															},
															"stretch_mode": {
																Type:         pluginsdk.TypeString,
																Optional:     true,
																Default:      string(encodings.StretchModeAutoSize),
																ValidateFunc: validation.StringInSlice(encodings.PossibleValuesForStretchMode(), false),
															},
															"sync_mode": {
																Type:         pluginsdk.TypeString,
																Optional:     true,
																Default:      string(encodings.VideoSyncModeAuto),
																ValidateFunc: validation.StringInSlice(encodings.PossibleValuesForVideoSyncMode(), false),
															},
														},
													},
												},
												"h265_video": {
													Type:     pluginsdk.TypeList,
													Optional: true,
													MaxItems: 1,
													Elem: &pluginsdk.Resource{
														Schema: map[string]*schema.Schema{
															"complexity": {
																Type:         pluginsdk.TypeString,
																Optional:     true,
																Default:      string(encodings.H265ComplexityBalanced),
																ValidateFunc: validation.StringInSlice(encodings.PossibleValuesForH265Complexity(), false),
															},
															"key_frame_interval": {
																Type:         pluginsdk.TypeString,
																Optional:     true,
																Default:      "PT2S",
																ValidateFunc: validate.ISO8601DurationBetween("PT0.5S", "PT20S"),
															},
															"label": {
																Type:         pluginsdk.TypeString,
																Optional:     true,
																ValidateFunc: validation.StringIsNotEmpty,
															},
															"layer": {
																Type:     pluginsdk.TypeList,
																Optional: true,
																MinItems: 1,
																Elem: &pluginsdk.Resource{
																	Schema: map[string]*schema.Schema{
																		"bitrate": {
																			Type:         pluginsdk.TypeInt,
																			Required:     true,
																			ValidateFunc: validation.IntAtLeast(1),
																		},
																		"adaptive_b_frame_enabled": {
																			Type:     pluginsdk.TypeBool,
																			Optional: true,
																			Default:  true,
																		},
																		"b_frames": {
																			Type:         pluginsdk.TypeInt,
																			Optional:     true,
																			Computed:     true,
																			ValidateFunc: validation.IntAtLeast(0),
																		},
																		"buffer_window": {
																			Type:         pluginsdk.TypeString,
																			Optional:     true,
																			Default:      "PT5S",
																			ValidateFunc: validate.ISO8601DurationBetween("PT0.1S", "PT100S"),
																		},
																		"crf": {
																			Type:         pluginsdk.TypeFloat,
																			Optional:     true,
																			Default:      28,
																			ValidateFunc: validation.FloatBetween(0, 51),
																		},
																		"frame_rate": {
																			Type:         pluginsdk.TypeString,
																			Optional:     true,
																			ValidateFunc: validation.StringIsNotEmpty,
																		},
																		"height": {
																			Type:         pluginsdk.TypeString,
																			Optional:     true,
																			Computed:     true,
																			ValidateFunc: validation.StringIsNotEmpty,
																		},
																		"label": {
																			Type:         pluginsdk.TypeString,
																			Optional:     true,
																			ValidateFunc: validation.StringIsNotEmpty,
																		},
																		"level": {
																			Type:         pluginsdk.TypeString,
																			Optional:     true,
																			Default:      "auto",
																			ValidateFunc: validation.StringIsNotEmpty,
																		},
																		"max_bitrate": {
																			Type:     pluginsdk.TypeInt,
																			Optional: true,
																			Computed: true,
																		},
																		"profile": {
																			Type:         pluginsdk.TypeString,
																			Optional:     true,
																			Default:      string(encodings.H265VideoProfileAuto),
																			ValidateFunc: validation.StringInSlice(encodings.PossibleValuesForH265VideoProfile(), false),
																		},
																		"reference_frames": {
																			Type:         pluginsdk.TypeInt,
																			Optional:     true,
																			Computed:     true,
																			ValidateFunc: validation.IntAtLeast(0),
																		},
																		"slices": {
																			Type:         pluginsdk.TypeInt,
																			Optional:     true,
																			Computed:     true,
																			ValidateFunc: validation.IntAtLeast(0),
																		},
																		"width": {
																			Type:         pluginsdk.TypeString,
																			Optional:     true,
																			Computed:     true,
																			ValidateFunc: validation.StringIsNotEmpty,
																		},
																	},
																},
															},
															"scene_change_detection_enabled": {
																Type:     pluginsdk.TypeBool,
																Optional: true,
																Default:  false,
															},
															"stretch_mode": {
																Type:         pluginsdk.TypeString,
																Optional:     true,
																Default:      string(encodings.StretchModeAutoSize),
																ValidateFunc: validation.StringInSlice(encodings.PossibleValuesForStretchMode(), false),
															},
															"sync_mode": {
																Type:         pluginsdk.TypeString,
																Optional:     true,
																Default:      string(encodings.VideoSyncModeAuto),
																ValidateFunc: validation.StringInSlice(encodings.PossibleValuesForVideoSyncMode(), false),
															},
														},
													},
												},
												"jpg_image": {
													Type:     pluginsdk.TypeList,
													Optional: true,
													MaxItems: 1,
													Elem: &pluginsdk.Resource{
														Schema: map[string]*schema.Schema{
															"start": {
																Type:         pluginsdk.TypeString,
																Required:     true,
																ValidateFunc: validation.StringIsNotEmpty,
															},
															"key_frame_interval": {
																Type:         pluginsdk.TypeString,
																Optional:     true,
																Default:      "PT2S",
																ValidateFunc: validate.ISO8601DurationBetween("PT0.5S", "PT20S"),
															},
															"label": {
																Type:         pluginsdk.TypeString,
																Optional:     true,
																ValidateFunc: validation.StringIsNotEmpty,
															},
															"layer": {
																Type:     pluginsdk.TypeList,
																Optional: true,
																MinItems: 1,
																Elem: &pluginsdk.Resource{
																	Schema: map[string]*schema.Schema{
																		"height": {
																			Type:         pluginsdk.TypeString,
																			Optional:     true,
																			ValidateFunc: validation.StringIsNotEmpty,
																		},
																		"label": {
																			Type:         pluginsdk.TypeString,
																			Optional:     true,
																			ValidateFunc: validation.StringIsNotEmpty,
																		},
																		"quality": {
																			Type:         pluginsdk.TypeInt,
																			Optional:     true,
																			Default:      70,
																			ValidateFunc: validation.IntBetween(0, 100),
																		},
																		"width": {
																			Type:         pluginsdk.TypeString,
																			Optional:     true,
																			ValidateFunc: validation.StringIsNotEmpty,
																		},
																	},
																},
															},
															"range": {
																Type:         pluginsdk.TypeString,
																Optional:     true,
																Default:      "100%",
																ValidateFunc: validation.StringIsNotEmpty,
															},
															"sprite_column": {
																Type:         pluginsdk.TypeInt,
																Optional:     true,
																Default:      0,
																ValidateFunc: validation.IntAtLeast(0),
															},
															"step": {
																Type:         pluginsdk.TypeString,
																Optional:     true,
																ValidateFunc: validation.StringIsNotEmpty,
															},
															"stretch_mode": {
																Type:         pluginsdk.TypeString,
																Optional:     true,
																Default:      string(encodings.StretchModeAutoSize),
																ValidateFunc: validation.StringInSlice(encodings.PossibleValuesForStretchMode(), false),
															},
															"sync_mode": {
																Type:         pluginsdk.TypeString,
																Optional:     true,
																Default:      string(encodings.VideoSyncModeAuto),
																ValidateFunc: validation.StringInSlice(encodings.PossibleValuesForVideoSyncMode(), false),
															},
														},
													},
												},
												"png_image": {
													Type:     pluginsdk.TypeList,
													Optional: true,
													MaxItems: 1,
													Elem: &pluginsdk.Resource{
														Schema: map[string]*schema.Schema{
															"start": {
																Type:         pluginsdk.TypeString,
																Required:     true,
																ValidateFunc: validation.StringIsNotEmpty,
															},
															"key_frame_interval": {
																Type:         pluginsdk.TypeString,
																Optional:     true,
																Default:      "PT2S",
																ValidateFunc: validate.ISO8601Duration,
															},
															"label": {
																Type:         pluginsdk.TypeString,
																Optional:     true,
																ValidateFunc: validation.StringIsNotEmpty,
															},
															"layer": {
																Type:     pluginsdk.TypeList,
																Optional: true,
																MinItems: 1,
																Elem: &pluginsdk.Resource{
																	Schema: map[string]*schema.Schema{
																		"height": {
																			Type:         pluginsdk.TypeString,
																			Optional:     true,
																			ValidateFunc: validation.StringIsNotEmpty,
																		},
																		"label": {
																			Type:         pluginsdk.TypeString,
																			Optional:     true,
																			ValidateFunc: validation.StringIsNotEmpty,
																		},
																		"width": {
																			Type:         pluginsdk.TypeString,
																			Optional:     true,
																			ValidateFunc: validation.StringIsNotEmpty,
																		},
																	},
																},
															},
															"range": {
																Type:         pluginsdk.TypeString,
																Optional:     true,
																Default:      "100%",
																ValidateFunc: validation.StringIsNotEmpty,
															},
															"step": {
																Type:         pluginsdk.TypeString,
																Optional:     true,
																ValidateFunc: validation.StringIsNotEmpty,
															},
															"stretch_mode": {
																Type:         pluginsdk.TypeString,
																Optional:     true,
																Default:      string(encodings.StretchModeAutoSize),
																ValidateFunc: validation.StringInSlice(encodings.PossibleValuesForStretchMode(), false),
															},
															"sync_mode": {
																Type:         pluginsdk.TypeString,
																Optional:     true,
																Default:      string(encodings.VideoSyncModeAuto),
																ValidateFunc: validation.StringInSlice(encodings.PossibleValuesForVideoSyncMode(), false),
															},
														},
													},
												},
											},
										},
									},
									"format": {
										Type:     pluginsdk.TypeList,
										Required: true,
										MinItems: 1,
										Elem: &pluginsdk.Resource{
											Schema: map[string]*schema.Schema{
												"jpg": {
													Type:     pluginsdk.TypeList,
													Optional: true,
													MaxItems: 1,
													Elem: &pluginsdk.Resource{
														Schema: map[string]*schema.Schema{
															"filename_pattern": {
																Type:         pluginsdk.TypeString,
																Required:     true,
																ValidateFunc: validation.StringIsNotEmpty,
															},
														},
													},
												},
												"mp4": {
													Type:     pluginsdk.TypeList,
													Optional: true,
													MaxItems: 1,
													Elem: &pluginsdk.Resource{
														Schema: map[string]*schema.Schema{
															"filename_pattern": {
																Type:         pluginsdk.TypeString,
																Required:     true,
																ValidateFunc: validation.StringIsNotEmpty,
															},
															"output_file": {
																Type:     pluginsdk.TypeList,
																Optional: true,
																MinItems: 1,
																Elem: &pluginsdk.Resource{
																	Schema: map[string]*schema.Schema{
																		"labels": {
																			Type:     pluginsdk.TypeList,
																			Required: true,
																			MinItems: 1,
																			Elem: &pluginsdk.Schema{
																				Type:         pluginsdk.TypeString,
																				ValidateFunc: validation.StringIsNotEmpty,
																			},
																		},
																	},
																},
															},
														},
													},
												},
												"png": {
													Type:     pluginsdk.TypeList,
													Optional: true,
													MaxItems: 1,
													Elem: &pluginsdk.Resource{
														Schema: map[string]*schema.Schema{
															"filename_pattern": {
																Type:         pluginsdk.TypeString,
																Required:     true,
																ValidateFunc: validation.StringIsNotEmpty,
															},
														},
													},
												},
												"transport_stream": {
													Type:     pluginsdk.TypeList,
													Optional: true,
													MaxItems: 1,
													Elem: &pluginsdk.Resource{
														Schema: map[string]*schema.Schema{
															"filename_pattern": {
																Type:         pluginsdk.TypeString,
																Required:     true,
																ValidateFunc: validation.StringIsNotEmpty,
															},
															"output_file": {
																Type:     pluginsdk.TypeList,
																Optional: true,
																MinItems: 1,
																Elem: &pluginsdk.Resource{
																	Schema: map[string]*schema.Schema{
																		"labels": {
																			Type:     pluginsdk.TypeList,
																			Required: true,
																			MinItems: 1,
																			Elem: &pluginsdk.Schema{
																				Type:         pluginsdk.TypeString,
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
									},
									"experimental_options": {
										Type:     pluginsdk.TypeMap,
										Optional: true,
										Elem: &pluginsdk.Schema{
											Type: pluginsdk.TypeString,
										},
									},
									"filter": {
										Type:     pluginsdk.TypeList,
										Optional: true,
										MaxItems: 1,
										Elem: &pluginsdk.Resource{
											Schema: map[string]*schema.Schema{
												"crop_rectangle": {
													Type:     pluginsdk.TypeList,
													Optional: true,
													MaxItems: 1,
													Elem: &pluginsdk.Resource{
														Schema: map[string]*schema.Schema{
															"height": {
																Type:         pluginsdk.TypeString,
																Optional:     true,
																ValidateFunc: validation.StringIsNotEmpty,
															},
															"left": {
																Type:         pluginsdk.TypeString,
																Optional:     true,
																ValidateFunc: validation.StringIsNotEmpty,
															},
															"top": {
																Type:         pluginsdk.TypeString,
																Optional:     true,
																ValidateFunc: validation.StringIsNotEmpty,
															},
															"width": {
																Type:         pluginsdk.TypeString,
																Optional:     true,
																ValidateFunc: validation.StringIsNotEmpty,
															},
														},
													},
												},
												"deinterlace": {
													Type:     pluginsdk.TypeList,
													Optional: true,
													MaxItems: 1,
													Elem: &pluginsdk.Resource{
														Schema: map[string]*schema.Schema{
															"parity": {
																Type:         pluginsdk.TypeString,
																Optional:     true,
																Default:      string(encodings.DeinterlaceParityAuto),
																ValidateFunc: validation.StringInSlice(encodings.PossibleValuesForDeinterlaceParity(), false),
															},
															"mode": {
																Type:         pluginsdk.TypeString,
																Optional:     true,
																Default:      string(encodings.DeinterlaceModeAutoPixelAdaptive),
																ValidateFunc: validation.StringInSlice(encodings.PossibleValuesForDeinterlaceMode(), false),
															},
														},
													},
												},
												"fade_in": {
													Type:     pluginsdk.TypeList,
													Optional: true,
													MaxItems: 1,
													Elem: &pluginsdk.Resource{
														Schema: map[string]*schema.Schema{
															"duration": {
																Type:         pluginsdk.TypeString,
																Required:     true,
																ValidateFunc: validation.StringIsNotEmpty,
															},
															"fade_color": {
																Type:         pluginsdk.TypeString,
																Required:     true,
																ValidateFunc: validation.StringIsNotEmpty,
															},
															"start": {
																Type:         pluginsdk.TypeString,
																Optional:     true,
																Default:      "0",
																ValidateFunc: validation.StringIsNotEmpty,
															},
														},
													},
												},
												"fade_out": {
													Type:     pluginsdk.TypeList,
													Optional: true,
													MaxItems: 1,
													Elem: &pluginsdk.Resource{
														Schema: map[string]*schema.Schema{
															"duration": {
																Type:         pluginsdk.TypeString,
																Required:     true,
																ValidateFunc: validation.StringIsNotEmpty,
															},
															"fade_color": {
																Type:         pluginsdk.TypeString,
																Required:     true,
																ValidateFunc: validation.StringIsNotEmpty,
															},
															"start": {
																Type:         pluginsdk.TypeString,
																Optional:     true,
																Default:      "0",
																ValidateFunc: validation.StringIsNotEmpty,
															},
														},
													},
												},
												"overlay": {
													Type:     pluginsdk.TypeList,
													Optional: true,
													MinItems: 1,
													Elem: &pluginsdk.Resource{
														Schema: map[string]*schema.Schema{
															"audio": {
																Type:     pluginsdk.TypeList,
																Optional: true,
																MaxItems: 1,
																Elem: &pluginsdk.Resource{
																	Schema: map[string]*schema.Schema{
																		"input_label": {
																			Type:         pluginsdk.TypeString,
																			Required:     true,
																			ValidateFunc: validation.StringIsNotEmpty,
																		},
																		"audio_gain_level": {
																			Type:         pluginsdk.TypeFloat,
																			Optional:     true,
																			Default:      1.0,
																			ValidateFunc: validation.FloatBetween(0, 1.0),
																		},
																		"end": {
																			Type:         pluginsdk.TypeString,
																			Optional:     true,
																			ValidateFunc: validate.ISO8601Duration,
																		},
																		"fade_in_duration": {
																			Type:         pluginsdk.TypeString,
																			Optional:     true,
																			ValidateFunc: validate.ISO8601Duration,
																		},
																		"fade_out_duration": {
																			Type:         pluginsdk.TypeString,
																			Optional:     true,
																			ValidateFunc: validate.ISO8601Duration,
																		},
																		"start": {
																			Type:         pluginsdk.TypeString,
																			Optional:     true,
																			ValidateFunc: validate.ISO8601Duration,
																		},
																	},
																},
															},
															"video": {
																Type:     pluginsdk.TypeList,
																Optional: true,
																MaxItems: 1,
																Elem: &pluginsdk.Resource{
																	Schema: map[string]*schema.Schema{
																		"input_label": {
																			Type:         pluginsdk.TypeString,
																			Required:     true,
																			ValidateFunc: validation.StringIsNotEmpty,
																		},
																		"audio_gain_level": {
																			Type:         pluginsdk.TypeFloat,
																			Optional:     true,
																			Default:      1.0,
																			ValidateFunc: validation.FloatBetween(0, 1.0),
																		},
																		"crop_rectangle": {
																			Type:     pluginsdk.TypeList,
																			Optional: true,
																			MaxItems: 1,
																			Elem: &pluginsdk.Resource{
																				Schema: map[string]*schema.Schema{
																					"height": {
																						Type:         pluginsdk.TypeString,
																						Optional:     true,
																						ValidateFunc: validation.StringIsNotEmpty,
																					},
																					"left": {
																						Type:         pluginsdk.TypeString,
																						Optional:     true,
																						ValidateFunc: validation.StringIsNotEmpty,
																					},
																					"top": {
																						Type:         pluginsdk.TypeString,
																						Optional:     true,
																						ValidateFunc: validation.StringIsNotEmpty,
																					},
																					"width": {
																						Type:         pluginsdk.TypeString,
																						Optional:     true,
																						ValidateFunc: validation.StringIsNotEmpty,
																					},
																				},
																			},
																		},
																		"end": {
																			Type:         pluginsdk.TypeString,
																			Optional:     true,
																			ValidateFunc: validate.ISO8601Duration,
																		},
																		"fade_in_duration": {
																			Type:         pluginsdk.TypeString,
																			Optional:     true,
																			ValidateFunc: validate.ISO8601Duration,
																		},
																		"fade_out_duration": {
																			Type:         pluginsdk.TypeString,
																			Optional:     true,
																			ValidateFunc: validate.ISO8601Duration,
																		},
																		"opacity": {
																			Type:         pluginsdk.TypeFloat,
																			Optional:     true,
																			Default:      1.0,
																			ValidateFunc: validation.FloatBetween(0, 1.0),
																		},
																		"position": {
																			Type:     pluginsdk.TypeList,
																			Optional: true,
																			MaxItems: 1,
																			Elem: &pluginsdk.Resource{
																				Schema: map[string]*schema.Schema{
																					"height": {
																						Type:         pluginsdk.TypeString,
																						Optional:     true,
																						ValidateFunc: validation.StringIsNotEmpty,
																					},
																					"left": {
																						Type:         pluginsdk.TypeString,
																						Optional:     true,
																						ValidateFunc: validation.StringIsNotEmpty,
																					},
																					"top": {
																						Type:         pluginsdk.TypeString,
																						Optional:     true,
																						ValidateFunc: validation.StringIsNotEmpty,
																					},
																					"width": {
																						Type:         pluginsdk.TypeString,
																						Optional:     true,
																						ValidateFunc: validation.StringIsNotEmpty,
																					},
																				},
																			},
																		},
																		"start": {
																			Type:         pluginsdk.TypeString,
																			Optional:     true,
																			ValidateFunc: validate.ISO8601Duration,
																		},
																	},
																},
															},
														},
													},
												},
												"rotation": {
													Type:         pluginsdk.TypeString,
													Optional:     true,
													Default:      string(encodings.RotationAuto),
													ValidateFunc: validation.StringInSlice(encodings.PossibleValuesForRotation(), false),
												},
											},
										},
									},
								},
							},
						},
						"relative_priority": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							Default:      string(encodings.PriorityNormal),
							ValidateFunc: validation.StringInSlice(encodings.PossibleValuesForPriority(), false),
						},
					},
				},
			},
		},
	}
}

func resourceMediaTransformCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Media.V20220701Client.Encodings
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := encodings.NewTransformID(subscriptionId, d.Get("resource_group_name").(string), d.Get("media_services_account_name").(string), d.Get("name").(string))
	if d.IsNewResource() {
		existing, err := client.TransformsGet(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for existing %s: %+v", id, err)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_media_transform", id.ID())
		}
	}

	payload := encodings.Transform{
		Properties: &encodings.TransformProperties{
			Description: utils.String(d.Get("description").(string)),
		},
	}

	if v, ok := d.GetOk("output"); ok {
		transformOutput, err := expandTransformOutputs(v.([]interface{}))
		if err != nil {
			return err
		}
		payload.Properties.Outputs = *transformOutput
	}

	if _, err := client.TransformsCreateOrUpdate(ctx, id, payload); err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceMediaTransformRead(d, meta)
}

func resourceMediaTransformRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Media.V20220701Client.Encodings
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := encodings.ParseTransformID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.TransformsGet(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[INFO] %s was not found - removing from state", *id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.TransformName)
	d.Set("media_services_account_name", id.MediaServiceName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
			d.Set("description", props.Description)

			outputs := flattenTransformOutputs(props.Outputs)
			if err := d.Set("output", outputs); err != nil {
				return fmt.Errorf("flattening `output`: %s", err)
			}
		}
	}

	return nil
}

func resourceMediaTransformDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Media.V20220701Client.Encodings
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := encodings.ParseTransformID(d.Id())
	if err != nil {
		return err
	}

	if _, err := client.TransformsDelete(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return nil
}

func expandTransformOutputs(input []interface{}) (*[]encodings.TransformOutput, error) {
	results := make([]encodings.TransformOutput, 0)

	for _, transformOutputRaw := range input {
		if transformOutputRaw == nil {
			continue
		}
		transform := transformOutputRaw.(map[string]interface{})

		preset, err := expandPreset(transform)
		if err != nil {
			return nil, err
		}

		transformOutput := encodings.TransformOutput{
			Preset: preset,
		}

		if v := transform["on_error_action"].(string); v != "" {
			transformOutput.OnError = pointer.To(encodings.OnErrorType(v))
		}

		if v := transform["relative_priority"].(string); v != "" {
			transformOutput.RelativePriority = pointer.To(encodings.Priority(v))
		}

		results = append(results, transformOutput)
	}

	return &results, nil
}

func flattenTransformOutputs(input []encodings.TransformOutput) []interface{} {
	results := make([]interface{}, 0)
	for _, transformOutput := range input {
		onErrorAction := ""
		if transformOutput.OnError != nil {
			onErrorAction = string(*transformOutput.OnError)
		}

		relativePriority := ""
		if transformOutput.RelativePriority != nil {
			relativePriority = string(*transformOutput.RelativePriority)
		}

		preset := flattenPreset(transformOutput.Preset)
		results = append(results, map[string]interface{}{
			"audio_analyzer_preset": preset.audioAnalyzerPresets,
			"builtin_preset":        preset.builtInPresets,
			"custom_preset":         preset.customPresets,
			"face_detector_preset":  preset.faceDetectorPresets,
			"on_error_action":       onErrorAction,
			"relative_priority":     relativePriority,
			"video_analyzer_preset": preset.videoAnalyzerPresets,
		})
	}

	return results
}

func expandPreset(transform map[string]interface{}) (encodings.Preset, error) {
	audioAnalyzerPresets := transform["audio_analyzer_preset"].([]interface{})
	builtInPresets := transform["builtin_preset"].([]interface{})
	customPresets := transform["custom_preset"].([]interface{})
	faceDetectorPresets := transform["face_detector_preset"].([]interface{})
	videoAnalyzerPresets := transform["video_analyzer_preset"].([]interface{})

	presetsCount := 0
	if len(audioAnalyzerPresets) > 0 {
		presetsCount++
	}
	if len(builtInPresets) > 0 {
		presetsCount++
	}
	if len(customPresets) > 0 {
		presetsCount++
	}
	if len(faceDetectorPresets) > 0 {
		presetsCount++
	}
	if len(videoAnalyzerPresets) > 0 {
		presetsCount++
	}
	if presetsCount == 0 {
		return nil, fmt.Errorf("output must contain at least one type of preset: builtin_preset, custom_preset, face_detector_preset, video_analyzer_preset or audio_analyzer_preset")
	}
	if presetsCount > 1 {
		return nil, fmt.Errorf("more than one type of preset in the same output is not allowed")
	}

	if len(audioAnalyzerPresets) > 0 {
		preset := audioAnalyzerPresets[0].(map[string]interface{})

		options, err := expandExperimentalOptions(preset["experimental_options"].(map[string]interface{}))
		if err != nil {
			return nil, err
		}

		audioAnalyzerPreset := &encodings.AudioAnalyzerPreset{
			ExperimentalOptions: options,
		}

		if v := preset["audio_language"].(string); v != "" {
			audioAnalyzerPreset.AudioLanguage = utils.String(v)
		}
		if v := preset["audio_analysis_mode"].(string); v != "" {
			audioAnalyzerPreset.Mode = pointer.To(encodings.AudioAnalysisMode(v))
		}
		return audioAnalyzerPreset, nil
	}

	if len(builtInPresets) > 0 {
		preset := builtInPresets[0].(map[string]interface{})
		presetName := preset["preset_name"].(string)
		builtInPreset := &encodings.BuiltInStandardEncoderPreset{
			PresetName:     encodings.EncoderNamedPreset(presetName),
			Configurations: expandBuiltInPresetConfiguration(preset["preset_configuration"].([]interface{})),
		}
		return builtInPreset, nil
	}

	if len(customPresets) > 0 {
		preset := customPresets[0].(map[string]interface{})

		codecs, err := expandCustomPresetCodecs(preset["codec"].([]interface{}))
		if err != nil {
			return nil, err
		}

		experimentalOptions, err := expandExperimentalOptions(preset["experimental_options"].(map[string]interface{}))
		if err != nil {
			return nil, err
		}

		filters, err := expandCustomPresetFilters(preset["filter"].([]interface{}))
		if err != nil {
			return nil, err
		}

		formats, err := expandCustomPresetFormats(preset["format"].([]interface{}))
		if err != nil {
			return nil, err
		}
		builtInPreset := &encodings.StandardEncoderPreset{
			Codecs:              codecs,
			ExperimentalOptions: experimentalOptions,
			Filters:             filters,
			Formats:             formats,
		}
		return builtInPreset, nil
	}

	if len(faceDetectorPresets) > 0 {
		preset := faceDetectorPresets[0].(map[string]interface{})

		options, err := expandExperimentalOptions(preset["experimental_options"].(map[string]interface{}))
		if err != nil {
			return nil, err
		}

		faceDetectorPreset := &encodings.FaceDetectorPreset{
			ExperimentalOptions: options,
		}

		if v := preset["analysis_resolution"].(string); v != "" {
			faceDetectorPreset.Resolution = pointer.To(encodings.AnalysisResolution(v))
		}

		if v := preset["blur_type"].(string); v != "" {
			faceDetectorPreset.BlurType = pointer.To(encodings.BlurType(v))
		}

		if v := preset["face_redactor_mode"].(string); v != "" {
			faceDetectorPreset.Mode = pointer.To(encodings.FaceRedactorMode(v))
		}

		return faceDetectorPreset, nil
	}

	if len(videoAnalyzerPresets) > 0 {
		presets := transform["video_analyzer_preset"].([]interface{})
		preset := presets[0].(map[string]interface{})

		options, err := expandExperimentalOptions(preset["experimental_options"].(map[string]interface{}))
		if err != nil {
			return nil, err
		}

		videoAnalyzerPreset := &encodings.VideoAnalyzerPreset{
			ExperimentalOptions: options,
		}

		if v := preset["audio_language"].(string); v != "" {
			videoAnalyzerPreset.AudioLanguage = utils.String(v)
		}
		if v := preset["audio_analysis_mode"].(string); v != "" {
			videoAnalyzerPreset.Mode = pointer.To(encodings.AudioAnalysisMode(v))
		}
		if v := preset["insights_type"].(string); v != "" {
			videoAnalyzerPreset.InsightsToExtract = pointer.To(encodings.InsightsType(v))
		}
		return videoAnalyzerPreset, nil
	}

	return nil, fmt.Errorf("output must contain at least one type of preset: builtin_preset, custom_preset, face_detector_preset, video_analyzer_preset or audio_analyzer_preset")
}

type flattenedPresets struct {
	audioAnalyzerPresets []interface{}
	builtInPresets       []interface{}
	customPresets        []interface{}
	faceDetectorPresets  []interface{}
	videoAnalyzerPresets []interface{}
}

func flattenPreset(input encodings.Preset) flattenedPresets {
	out := flattenedPresets{
		audioAnalyzerPresets: []interface{}{},
		builtInPresets:       []interface{}{},
		customPresets:        []interface{}{},
		faceDetectorPresets:  []interface{}{},
		videoAnalyzerPresets: []interface{}{},
	}
	if input == nil {
		return out
	}

	if v, ok := input.(encodings.AudioAnalyzerPreset); ok {
		language := ""
		if v.AudioLanguage != nil {
			language = *v.AudioLanguage
		}
		mode := ""
		if v.Mode != nil {
			mode = string(*v.Mode)
		}
		out.audioAnalyzerPresets = append(out.audioAnalyzerPresets, map[string]interface{}{
			"audio_analysis_mode":  mode,
			"audio_language":       language,
			"experimental_options": flattenExperimentalOptions(v.ExperimentalOptions),
		})
	}

	if v, ok := input.(encodings.StandardEncoderPreset); ok {
		out.customPresets = append(out.customPresets, map[string]interface{}{
			"codec":                flattenCustomPresetCodecs(v.Codecs),
			"experimental_options": flattenExperimentalOptions(v.ExperimentalOptions),
			"filter":               flattenCustomPresetFilters(v.Filters),
			"format":               flattenCustomPresetFormats(v.Formats),
		})
	}

	if v, ok := input.(encodings.BuiltInStandardEncoderPreset); ok {
		out.builtInPresets = append(out.builtInPresets, map[string]interface{}{
			"preset_name":          string(v.PresetName),
			"preset_configuration": flattenBuiltInPresetConfiguration(v.Configurations),
		})
	}

	if v, ok := input.(encodings.FaceDetectorPreset); ok {
		resolution := ""
		if v.Resolution != nil {
			resolution = string(*v.Resolution)
		}

		blurType := ""
		if v.BlurType != nil {
			blurType = string(*v.BlurType)
		}

		mode := ""
		if v.Mode != nil {
			mode = string(*v.Mode)
		}

		out.faceDetectorPresets = append(out.faceDetectorPresets, map[string]interface{}{
			"analysis_resolution":  resolution,
			"blur_type":            blurType,
			"experimental_options": flattenExperimentalOptions(v.ExperimentalOptions),
			"face_redactor_mode":   mode,
		})
	}

	if v, ok := input.(encodings.VideoAnalyzerPreset); ok {
		audioLanguage := ""
		if v.AudioLanguage != nil {
			audioLanguage = *v.AudioLanguage
		}
		insightsType := ""
		if v.InsightsToExtract != nil {
			insightsType = string(*v.InsightsToExtract)
		}
		mode := ""
		if v.Mode != nil {
			mode = string(*v.Mode)
		}
		out.videoAnalyzerPresets = append(out.videoAnalyzerPresets, map[string]interface{}{
			"audio_analysis_mode":  mode,
			"audio_language":       audioLanguage,
			"insights_type":        insightsType,
			"experimental_options": flattenExperimentalOptions(v.ExperimentalOptions),
		})
	}

	return out
}

func expandBuiltInPresetConfiguration(input []interface{}) *encodings.PresetConfigurations {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	configuration := input[0].(map[string]interface{})
	result := encodings.PresetConfigurations{}

	if v := configuration["complexity"].(string); v != "" {
		result.Complexity = pointer.To(encodings.Complexity(v))
	}

	if v := configuration["interleave_output"].(string); v != "" {
		result.InterleaveOutput = pointer.To(encodings.InterleaveOutput(v))
	}

	if v := configuration["key_frame_interval_in_seconds"].(float64); v != 0 {
		result.KeyFrameIntervalInSeconds = utils.Float(v)
	}

	if v := configuration["max_bitrate_bps"].(int); v != 0 {
		result.MaxBitrateBps = utils.Int64(int64(v))
	}

	if v := configuration["max_height"].(int); v != 0 {
		result.MaxHeight = utils.Int64(int64(v))
	}

	if v := configuration["max_layers"].(int); v != 0 {
		result.MaxLayers = utils.Int64(int64(v))
	}

	if v := configuration["min_bitrate_bps"].(int); v != 0 {
		result.MinBitrateBps = utils.Int64(int64(v))
	}

	if v := configuration["min_height"].(int); v != 0 {
		result.MinHeight = utils.Int64(int64(v))
	}

	return &result
}

func flattenBuiltInPresetConfiguration(input *encodings.PresetConfigurations) []interface{} {
	if input == nil {
		return make([]interface{}, 0)
	}

	complexity := ""
	if input.Complexity != nil {
		complexity = string(*input.Complexity)
	}

	interleaveOutput := ""
	if input.InterleaveOutput != nil {
		interleaveOutput = string(*input.InterleaveOutput)
	}

	keyFrameIntervalInSeconds := 0.0
	if input.KeyFrameIntervalInSeconds != nil {
		keyFrameIntervalInSeconds = *input.KeyFrameIntervalInSeconds
	}

	maxBitrateBps := 0
	if input.MaxBitrateBps != nil {
		maxBitrateBps = int(*input.MaxBitrateBps)
	}

	maxHeight := 0
	if input.MaxHeight != nil {
		maxHeight = int(*input.MaxHeight)
	}

	maxLayers := 0
	if input.MaxLayers != nil {
		maxLayers = int(*input.MaxLayers)
	}

	minBitrateBps := 0
	if input.MinBitrateBps != nil {
		minBitrateBps = int(*input.MinBitrateBps)
	}

	minHeight := 0
	if input.MinHeight != nil {
		minHeight = int(*input.MinHeight)
	}

	return []interface{}{
		map[string]interface{}{
			"complexity":                    complexity,
			"interleave_output":             interleaveOutput,
			"key_frame_interval_in_seconds": keyFrameIntervalInSeconds,
			"max_bitrate_bps":               maxBitrateBps,
			"max_height":                    maxHeight,
			"max_layers":                    maxLayers,
			"min_bitrate_bps":               minBitrateBps,
			"min_height":                    minHeight,
		},
	}
}

func expandExperimentalOptions(input map[string]interface{}) (*map[string]string, error) {
	output := make(map[string]string, len(input))

	for k, v := range input {
		key := k
		value, ok := v.(string)
		if !ok {
			return nil, fmt.Errorf("expect experimental options value %q to be a string", value)
		}
		output[key] = value
	}

	return &output, nil
}

func flattenExperimentalOptions(input *map[string]string) map[string]interface{} {
	result := make(map[string]interface{}, 0)
	if input == nil {
		return result
	}
	for k, v := range *input {
		key := k
		value := v
		result[key] = value
	}

	return result
}

func expandCustomPresetCodecs(input []interface{}) ([]encodings.Codec, error) {
	if len(input) == 0 || input[0] == nil {
		return make([]encodings.Codec, 0), nil
	}

	results := make([]encodings.Codec, 0)

	for _, v := range input {
		if v == nil {
			continue
		}

		codec := v.(map[string]interface{})
		aacAudio := codec["aac_audio"].([]interface{})
		copyAudio := codec["copy_audio"].([]interface{})
		copyVideo := codec["copy_video"].([]interface{})
		ddAudio := codec["dd_audio"].([]interface{})
		h264Video := codec["h264_video"].([]interface{})
		h265Video := codec["h265_video"].([]interface{})
		jpgImage := codec["jpg_image"].([]interface{})
		pngImage := codec["png_image"].([]interface{})

		codecsCount := 0
		if len(aacAudio) > 0 {
			codecsCount++
		}
		if len(copyAudio) > 0 {
			codecsCount++
		}
		if len(copyVideo) > 0 {
			codecsCount++
		}
		if len(ddAudio) > 0 {
			codecsCount++
		}
		if len(h264Video) > 0 {
			codecsCount++
		}
		if len(h265Video) > 0 {
			codecsCount++
		}
		if len(jpgImage) > 0 {
			codecsCount++
		}
		if len(pngImage) > 0 {
			codecsCount++
		}
		if codecsCount == 0 {
			return nil, fmt.Errorf("custom preset codec must contain at least one type of: aac_audio, copy_audio, copy_video, dd_audio, h264_video, h265_video, jpg_image or png_image")
		}
		if codecsCount > 1 {
			return nil, fmt.Errorf("more than one type of codec in the same custom preset codec is not allowed")
		}

		if len(aacAudio) > 0 {
			results = append(results, expandCustomPresetCodecsAacAudio(aacAudio))
		}
		if len(copyAudio) > 0 {
			results = append(results, expandCustomPresetCodecsCopyAudio(copyAudio))
		}
		if len(copyVideo) > 0 {
			results = append(results, expandCustomPresetCodecsCopyVideo(copyVideo))
		}
		if len(ddAudio) > 0 {
			results = append(results, expandCustomPresetCodecsDdAudio(ddAudio))
		}
		if len(h264Video) > 0 {
			results = append(results, expandCustomPresetCodecsH264Video(h264Video))
		}
		if len(h265Video) > 0 {
			results = append(results, expandCustomPresetCodecsH265Video(h265Video))
		}
		if len(jpgImage) > 0 {
			results = append(results, expandCustomPresetCodecsJpgImage(jpgImage))
		}
		if len(pngImage) > 0 {
			results = append(results, expandCustomPresetCodecsPngImage(pngImage))
		}
	}

	return results, nil
}

type flattenedCustomPresetsCodec struct {
	aacAudio  []interface{}
	copyAudio []interface{}
	copyVideo []interface{}
	ddAudio   []interface{}
	h264Video []interface{}
	h265Video []interface{}
	jpgImage  []interface{}
	pngImage  []interface{}
}

func flattenCustomPresetCodecs(input []encodings.Codec) []interface{} {
	if len(input) == 0 {
		return make([]interface{}, 0)
	}

	results := make([]interface{}, 0)

	for _, v := range input {
		result := flattenedCustomPresetsCodec{
			aacAudio:  []interface{}{},
			copyAudio: []interface{}{},
			copyVideo: []interface{}{},
			ddAudio:   []interface{}{},
			h264Video: []interface{}{},
			h265Video: []interface{}{},
			jpgImage:  []interface{}{},
			pngImage:  []interface{}{},
		}

		if codec, ok := v.(encodings.AacAudio); ok {
			result.aacAudio = flattenCustomPresetCodecsAacAudio(codec)
		}

		if codec, ok := v.(encodings.CopyAudio); ok {
			result.copyAudio = flattenCustomPresetCodecsCopyAudio(codec)
		}

		if codec, ok := v.(encodings.CopyVideo); ok {
			result.copyVideo = flattenCustomPresetCodecsCopyVideo(codec)
		}

		if codec, ok := v.(encodings.DDAudio); ok {
			result.ddAudio = flattenCustomPresetCodecsDdAudio(codec)
		}

		if codec, ok := v.(encodings.H264Video); ok {
			result.h264Video = flattenCustomPresetCodecsH264Video(codec)
		}

		if codec, ok := v.(encodings.H265Video); ok {
			result.h265Video = flattenCustomPresetCodecsH265Video(codec)
		}

		if codec, ok := v.(encodings.JpgImage); ok {
			result.jpgImage = flattenCustomPresetCodecsJpgImage(codec)
		}

		if codec, ok := v.(encodings.PngImage); ok {
			result.pngImage = flattenCustomPresetCodecsPngImage(codec)
		}

		results = append(results, map[string]interface{}{
			"aac_audio":  result.aacAudio,
			"copy_audio": result.copyAudio,
			"copy_video": result.copyVideo,
			"dd_audio":   result.ddAudio,
			"h264_video": result.h264Video,
			"h265_video": result.h265Video,
			"jpg_image":  result.jpgImage,
			"png_image":  result.pngImage,
		})
	}

	return results
}

func expandCustomPresetCodecsAacAudio(input []interface{}) encodings.AacAudio {
	if len(input) == 0 || input[0] == nil {
		return encodings.AacAudio{}
	}

	aacAudio := input[0].(map[string]interface{})
	result := encodings.AacAudio{}

	if v := aacAudio["bitrate"].(int); v != 0 {
		result.Bitrate = utils.Int64(int64(v))
	}

	if v := aacAudio["channels"].(int); v != 0 {
		result.Channels = utils.Int64(int64(v))
	}

	if v := aacAudio["label"].(string); v != "" {
		result.Label = utils.String(v)
	}

	if v := aacAudio["profile"].(string); v != "" {
		result.Profile = pointer.To(encodings.AacAudioProfile(v))
	}

	if v := aacAudio["sampling_rate"].(int); v != 0 {
		result.SamplingRate = utils.Int64(int64(v))
	}

	return result
}

func flattenCustomPresetCodecsAacAudio(input encodings.AacAudio) []interface{} {
	bitrate := 0
	if input.Bitrate != nil {
		bitrate = int(*input.Bitrate)
	}

	channels := 0
	if input.Channels != nil {
		channels = int(*input.Channels)
	}

	label := ""
	if input.Label != nil {
		label = *input.Label
	}

	profile := ""
	if input.Profile != nil {
		profile = string(*input.Profile)
	}

	samplingRate := 0
	if input.SamplingRate != nil {
		samplingRate = int(*input.SamplingRate)
	}

	return []interface{}{
		map[string]interface{}{
			"bitrate":       bitrate,
			"channels":      channels,
			"label":         label,
			"profile":       profile,
			"sampling_rate": samplingRate,
		},
	}
}

func expandCustomPresetCodecsCopyAudio(input []interface{}) encodings.CopyAudio {
	if len(input) == 0 || input[0] == nil {
		return encodings.CopyAudio{}
	}

	copyAudio := input[0].(map[string]interface{})
	result := encodings.CopyAudio{}

	if v := copyAudio["label"].(string); v != "" {
		result.Label = utils.String(v)
	}

	return result
}

func flattenCustomPresetCodecsCopyAudio(input encodings.CopyAudio) []interface{} {
	label := ""
	if input.Label != nil {
		label = *input.Label
	}

	return []interface{}{
		map[string]interface{}{
			"label": label,
		},
	}
}

func expandCustomPresetCodecsCopyVideo(input []interface{}) encodings.CopyVideo {
	if len(input) == 0 || input[0] == nil {
		return encodings.CopyVideo{}
	}

	copyVideo := input[0].(map[string]interface{})
	result := encodings.CopyVideo{}

	if v := copyVideo["label"].(string); v != "" {
		result.Label = utils.String(v)
	}

	return result
}

func flattenCustomPresetCodecsCopyVideo(input encodings.CopyVideo) []interface{} {
	label := ""
	if input.Label != nil {
		label = *input.Label
	}

	return []interface{}{
		map[string]interface{}{
			"label": label,
		},
	}
}

func expandCustomPresetCodecsDdAudio(input []interface{}) encodings.DDAudio {
	if len(input) == 0 || input[0] == nil {
		return encodings.DDAudio{}
	}

	ddAudio := input[0].(map[string]interface{})
	result := encodings.DDAudio{}

	if v := ddAudio["bitrate"].(int); v != 0 {
		result.Bitrate = utils.Int64(int64(v))
	}

	if v := ddAudio["channels"].(int); v != 0 {
		result.Channels = utils.Int64(int64(v))
	}

	if v := ddAudio["label"].(string); v != "" {
		result.Label = utils.String(v)
	}

	if v := ddAudio["sampling_rate"].(int); v != 0 {
		result.SamplingRate = utils.Int64(int64(v))
	}

	return result
}

func flattenCustomPresetCodecsDdAudio(input encodings.DDAudio) []interface{} {
	bitrate := 0
	if input.Bitrate != nil {
		bitrate = int(*input.Bitrate)
	}

	channels := 0
	if input.Channels != nil {
		channels = int(*input.Channels)
	}

	label := ""
	if input.Label != nil {
		label = *input.Label
	}

	samplingRate := 0
	if input.SamplingRate != nil {
		samplingRate = int(*input.SamplingRate)
	}

	return []interface{}{
		map[string]interface{}{
			"bitrate":       bitrate,
			"channels":      channels,
			"label":         label,
			"sampling_rate": samplingRate,
		},
	}
}

func expandCustomPresetCodecsH264Video(input []interface{}) encodings.H264Video {
	if len(input) == 0 || input[0] == nil {
		return encodings.H264Video{}
	}

	h264Video := input[0].(map[string]interface{})
	result := encodings.H264Video{
		Layers:               expandCustomPresetCodecsH264VideoLayers(h264Video["layer"].([]interface{})),
		SceneChangeDetection: utils.Bool(h264Video["scene_change_detection_enabled"].(bool)),
	}

	if v := h264Video["complexity"].(string); v != "" {
		result.Complexity = pointer.To(encodings.H264Complexity(v))
	}

	if v := h264Video["key_frame_interval"].(string); v != "" {
		result.KeyFrameInterval = utils.String(v)
	}

	if v := h264Video["label"].(string); v != "" {
		result.Label = utils.String(v)
	}

	if v := h264Video["rate_control_mode"].(string); v != "" {
		result.RateControlMode = pointer.To(encodings.H264RateControlMode(v))
	}

	if v := h264Video["stretch_mode"].(string); v != "" {
		result.StretchMode = pointer.To(encodings.StretchMode(v))
	}

	if v := h264Video["sync_mode"].(string); v != "" {
		result.SyncMode = pointer.To(encodings.VideoSyncMode(v))
	}

	return result
}

func flattenCustomPresetCodecsH264Video(input encodings.H264Video) []interface{} {
	complexity := ""
	if input.Complexity != nil {
		complexity = string(*input.Complexity)
	}

	keyFrameInterval := ""
	if input.KeyFrameInterval != nil {
		keyFrameInterval = *input.KeyFrameInterval
	}

	label := ""
	if input.Label != nil {
		label = *input.Label
	}

	rateControlMode := ""
	if input.RateControlMode != nil {
		rateControlMode = string(*input.RateControlMode)
	}

	sceneChangeDetectionEnabled := false
	if input.SceneChangeDetection != nil {
		sceneChangeDetectionEnabled = *input.SceneChangeDetection
	}

	stretchMode := ""
	if input.StretchMode != nil {
		stretchMode = string(*input.StretchMode)
	}

	syncMode := ""
	if input.SyncMode != nil {
		syncMode = string(*input.SyncMode)
	}
	return []interface{}{
		map[string]interface{}{
			"complexity":                     complexity,
			"key_frame_interval":             keyFrameInterval,
			"label":                          label,
			"layer":                          flattenCustomPresetCodecsH264VideoLayers(input.Layers),
			"rate_control_mode":              rateControlMode,
			"scene_change_detection_enabled": sceneChangeDetectionEnabled,
			"stretch_mode":                   stretchMode,
			"sync_mode":                      syncMode,
		},
	}
}

func expandCustomPresetCodecsH264VideoLayers(input []interface{}) *[]encodings.H264Layer {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	results := make([]encodings.H264Layer, 0)
	for _, layerRaw := range input {
		layer := layerRaw.(map[string]interface{})
		result := encodings.H264Layer{
			Bitrate:        int64(layer["bitrate"].(int)),
			AdaptiveBFrame: utils.Bool(layer["adaptive_b_frame_enabled"].(bool)),
		}

		if v := layer["b_frames"].(int); v != 0 {
			result.BFrames = utils.Int64(int64(v))
		}

		if v := layer["buffer_window"].(string); v != "" {
			result.BufferWindow = utils.String(v)
		}

		if v := layer["crf"].(float64); v != 0 {
			result.Crf = utils.Float(v)
		}

		if v := layer["entropy_mode"].(string); v != "" {
			result.EntropyMode = pointer.To(encodings.EntropyMode(v))
		}

		if v := layer["frame_rate"].(string); v != "" {
			result.FrameRate = utils.String(v)
		}

		if v := layer["height"].(string); v != "" {
			result.Height = utils.String(v)
		}

		if v := layer["label"].(string); v != "" {
			result.Label = utils.String(v)
		}

		if v := layer["level"].(string); v != "" {
			result.Level = utils.String(v)
		}

		if v := layer["max_bitrate"].(int); v != 0 {
			result.MaxBitrate = utils.Int64(int64(v))
		}

		if v := layer["profile"].(string); v != "" {
			result.Profile = pointer.To(encodings.H264VideoProfile(v))
		}

		if v := layer["reference_frames"].(int); v != 0 {
			result.ReferenceFrames = utils.Int64(int64(v))
		}

		if v := layer["slices"].(int); v != 0 {
			result.Slices = utils.Int64(int64(v))
		}

		if v := layer["width"].(string); v != "" {
			result.Width = utils.String(v)
		}

		results = append(results, result)
	}

	return &results
}

func flattenCustomPresetCodecsH264VideoLayers(input *[]encodings.H264Layer) []interface{} {
	if input == nil {
		return make([]interface{}, 0)
	}

	results := make([]interface{}, 0)

	for _, v := range *input {
		adaptiveBFrameEnabled := false
		if v.AdaptiveBFrame != nil {
			adaptiveBFrameEnabled = *v.AdaptiveBFrame
		}

		bFrames := 0
		if v.BFrames != nil {
			bFrames = int(*v.BFrames)
		}

		bufferWindow := ""
		if v.BufferWindow != nil {
			bufferWindow = *v.BufferWindow
		}

		crf := 0.0
		if v.Crf != nil {
			crf = *v.Crf
		}

		entropyMode := ""
		if v.EntropyMode != nil {
			entropyMode = string(*v.EntropyMode)
		}

		frameRate := ""
		if v.FrameRate != nil {
			frameRate = *v.FrameRate
		}

		height := ""
		if v.Height != nil {
			height = *v.Height
		}

		label := ""
		if v.Label != nil {
			label = *v.Label
		}

		level := ""
		if v.Level != nil {
			level = *v.Level
		}

		maxBitrate := 0
		if v.MaxBitrate != nil {
			maxBitrate = int(*v.MaxBitrate)
		}

		profile := ""
		if v.Profile != nil {
			profile = string(*v.Profile)
		}

		referenceFrames := 0
		if v.ReferenceFrames != nil {
			referenceFrames = int(*v.ReferenceFrames)
		}

		slices := 0
		if v.Slices != nil {
			slices = int(*v.Slices)
		}

		width := ""
		if v.Width != nil {
			width = *v.Width
		}

		results = append(results, map[string]interface{}{
			"bitrate":                  v.Bitrate,
			"adaptive_b_frame_enabled": adaptiveBFrameEnabled,
			"b_frames":                 bFrames,
			"buffer_window":            bufferWindow,
			"crf":                      crf,
			"entropy_mode":             entropyMode,
			"frame_rate":               frameRate,
			"height":                   height,
			"label":                    label,
			"level":                    level,
			"max_bitrate":              maxBitrate,
			"profile":                  profile,
			"reference_frames":         referenceFrames,
			"slices":                   slices,
			"width":                    width,
		})
	}

	return results
}

func expandCustomPresetCodecsH265Video(input []interface{}) encodings.H265Video {
	if len(input) == 0 || input[0] == nil {
		return encodings.H265Video{}
	}

	h265Video := input[0].(map[string]interface{})
	result := encodings.H265Video{
		Layers:               expandCustomPresetCodecsH265VideoLayers(h265Video["layer"].([]interface{})),
		SceneChangeDetection: utils.Bool(h265Video["scene_change_detection_enabled"].(bool)),
	}

	if v := h265Video["complexity"].(string); v != "" {
		result.Complexity = pointer.To(encodings.H265Complexity(v))
	}

	if v := h265Video["key_frame_interval"].(string); v != "" {
		result.KeyFrameInterval = utils.String(v)
	}

	if v := h265Video["label"].(string); v != "" {
		result.Label = utils.String(v)
	}

	if v := h265Video["stretch_mode"].(string); v != "" {
		result.StretchMode = pointer.To(encodings.StretchMode(v))
	}

	if v := h265Video["sync_mode"].(string); v != "" {
		result.SyncMode = pointer.To(encodings.VideoSyncMode(v))
	}

	return result
}

func flattenCustomPresetCodecsH265Video(input encodings.H265Video) []interface{} {
	complexity := ""
	if input.Complexity != nil {
		complexity = string(*input.Complexity)
	}

	keyFrameInterval := ""
	if input.KeyFrameInterval != nil {
		keyFrameInterval = *input.KeyFrameInterval
	}

	label := ""
	if input.Label != nil {
		label = *input.Label
	}

	sceneChangeDetectionEnabled := false
	if input.SceneChangeDetection != nil {
		sceneChangeDetectionEnabled = *input.SceneChangeDetection
	}

	stretchMode := ""
	if input.StretchMode != nil {
		stretchMode = string(*input.StretchMode)
	}

	syncMode := ""
	if input.SyncMode != nil {
		syncMode = string(*input.SyncMode)
	}
	return []interface{}{
		map[string]interface{}{
			"complexity":                     complexity,
			"key_frame_interval":             keyFrameInterval,
			"label":                          label,
			"layer":                          flattenCustomPresetCodecsH265VideoLayers(input.Layers),
			"scene_change_detection_enabled": sceneChangeDetectionEnabled,
			"stretch_mode":                   stretchMode,
			"sync_mode":                      syncMode,
		},
	}
}

func expandCustomPresetCodecsH265VideoLayers(input []interface{}) *[]encodings.H265Layer {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	results := make([]encodings.H265Layer, 0)
	for _, layerRaw := range input {
		layer := layerRaw.(map[string]interface{})
		result := encodings.H265Layer{
			Bitrate:        int64(layer["bitrate"].(int)),
			AdaptiveBFrame: utils.Bool(layer["adaptive_b_frame_enabled"].(bool)),
		}

		if v := layer["b_frames"].(int); v != 0 {
			result.BFrames = utils.Int64(int64(v))
		}

		if v := layer["buffer_window"].(string); v != "" {
			result.BufferWindow = utils.String(v)
		}

		if v := layer["crf"].(float64); v != 0 {
			result.Crf = utils.Float(v)
		}

		if v := layer["frame_rate"].(string); v != "" {
			result.FrameRate = utils.String(v)
		}

		if v := layer["height"].(string); v != "" {
			result.Height = utils.String(v)
		}

		if v := layer["label"].(string); v != "" {
			result.Label = utils.String(v)
		}

		if v := layer["level"].(string); v != "" {
			result.Level = utils.String(v)
		}

		if v := layer["max_bitrate"].(int); v != 0 {
			result.MaxBitrate = utils.Int64(int64(v))
		}

		if v := layer["profile"].(string); v != "" {
			result.Profile = pointer.To(encodings.H265VideoProfile(v))
		}

		if v := layer["reference_frames"].(int); v != 0 {
			result.ReferenceFrames = utils.Int64(int64(v))
		}

		if v := layer["slices"].(int); v != 0 {
			result.Slices = utils.Int64(int64(v))
		}

		if v := layer["width"].(string); v != "" {
			result.Width = utils.String(v)
		}

		results = append(results, result)
	}

	return &results
}

func flattenCustomPresetCodecsH265VideoLayers(input *[]encodings.H265Layer) []interface{} {
	if input == nil {
		return make([]interface{}, 0)
	}

	results := make([]interface{}, 0)

	for _, v := range *input {
		adaptiveBFrameEnabled := false
		if v.AdaptiveBFrame != nil {
			adaptiveBFrameEnabled = *v.AdaptiveBFrame
		}

		bFrames := 0
		if v.BFrames != nil {
			bFrames = int(*v.BFrames)
		}

		bufferWindow := ""
		if v.BufferWindow != nil {
			bufferWindow = *v.BufferWindow
		}

		crf := 0.0
		if v.Crf != nil {
			crf = *v.Crf
		}

		frameRate := ""
		if v.FrameRate != nil {
			frameRate = *v.FrameRate
		}

		height := ""
		if v.Height != nil {
			height = *v.Height
		}

		label := ""
		if v.Label != nil {
			label = *v.Label
		}

		level := ""
		if v.Level != nil {
			level = *v.Level
		}

		maxBitrate := 0
		if v.MaxBitrate != nil {
			maxBitrate = int(*v.MaxBitrate)
		}

		profile := ""
		if v.Profile != nil {
			profile = string(*v.Profile)
		}

		referenceFrames := 0
		if v.ReferenceFrames != nil {
			referenceFrames = int(*v.ReferenceFrames)
		}

		slices := 0
		if v.Slices != nil {
			slices = int(*v.Slices)
		}

		width := ""
		if v.Width != nil {
			width = *v.Width
		}

		results = append(results, map[string]interface{}{
			"bitrate":                  v.Bitrate,
			"adaptive_b_frame_enabled": adaptiveBFrameEnabled,
			"b_frames":                 bFrames,
			"buffer_window":            bufferWindow,
			"crf":                      crf,
			"frame_rate":               frameRate,
			"height":                   height,
			"label":                    label,
			"level":                    level,
			"max_bitrate":              maxBitrate,
			"profile":                  profile,
			"reference_frames":         referenceFrames,
			"slices":                   slices,
			"width":                    width,
		})
	}

	return results
}
func expandCustomPresetCodecsJpgImage(input []interface{}) encodings.JpgImage {
	if len(input) == 0 || input[0] == nil {
		return encodings.JpgImage{}
	}

	jpgImage := input[0].(map[string]interface{})
	result := encodings.JpgImage{
		Start:  jpgImage["start"].(string),
		Layers: expandCustomPresetCodecsJpgImageLayer(jpgImage["layer"].([]interface{})),
	}

	if v := jpgImage["key_frame_interval"].(string); v != "" {
		result.KeyFrameInterval = utils.String(v)
	}

	if v := jpgImage["label"].(string); v != "" {
		result.Label = utils.String(v)
	}

	if v := jpgImage["range"].(string); v != "" {
		result.Range = utils.String(v)
	}

	if v := jpgImage["sprite_column"].(int); v != 0 {
		result.SpriteColumn = utils.Int64(int64(v))
	}

	if v := jpgImage["step"].(string); v != "" {
		result.Step = utils.String(v)
	}

	if v := jpgImage["stretch_mode"].(string); v != "" {
		result.StretchMode = pointer.To(encodings.StretchMode(v))
	}

	if v := jpgImage["sync_mode"].(string); v != "" {
		result.SyncMode = pointer.To(encodings.VideoSyncMode(v))
	}

	return result
}

func flattenCustomPresetCodecsJpgImage(input encodings.JpgImage) []interface{} {
	keyFrameInterval := ""
	if input.KeyFrameInterval != nil {
		keyFrameInterval = *input.KeyFrameInterval
	}

	label := ""
	if input.Label != nil {
		label = *input.Label
	}

	rang := ""
	if input.Range != nil {
		rang = *input.Range
	}

	spriteColumn := 0
	if input.SpriteColumn != nil {
		spriteColumn = int(*input.SpriteColumn)
	}

	step := ""
	if input.Step != nil {
		step = *input.Step
	}

	stretchMode := ""
	if input.StretchMode != nil {
		stretchMode = string(*input.StretchMode)
	}

	syncMode := ""
	if input.SyncMode != nil {
		syncMode = string(*input.SyncMode)
	}

	return []interface{}{
		map[string]interface{}{
			"key_frame_interval": keyFrameInterval,
			"label":              label,
			"layer":              flattenCustomPresetCodecsJpgImageLayer(input.Layers),
			"range":              rang,
			"start":              input.Start,
			"step":               step,
			"sprite_column":      spriteColumn,
			"stretch_mode":       stretchMode,
			"sync_mode":          syncMode,
		},
	}
}

func expandCustomPresetCodecsJpgImageLayer(input []interface{}) *[]encodings.JpgLayer {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	results := make([]encodings.JpgLayer, 0)
	for _, layerRaw := range input {
		if layerRaw == nil {
			continue
		}

		layer := layerRaw.(map[string]interface{})
		result := encodings.JpgLayer{}

		if v := layer["height"].(string); v != "" {
			result.Height = utils.String(v)
		}

		if v := layer["label"].(string); v != "" {
			result.Label = utils.String(v)
		}

		if v := layer["quality"].(int); v != 0 {
			result.Quality = utils.Int64(int64(v))
		}

		if v := layer["width"].(string); v != "" {
			result.Width = utils.String(v)
		}

		results = append(results, result)
	}

	return &results
}

func flattenCustomPresetCodecsJpgImageLayer(input *[]encodings.JpgLayer) []interface{} {
	if input == nil {
		return make([]interface{}, 0)
	}

	results := make([]interface{}, 0)

	for _, v := range *input {
		height := ""
		if v.Height != nil {
			height = *v.Height
		}

		label := ""
		if v.Label != nil {
			label = *v.Label
		}

		quality := 0
		if v.Quality != nil {
			quality = int(*v.Quality)
		}

		width := ""
		if v.Width != nil {
			width = *v.Width
		}

		results = append(results, map[string]interface{}{
			"height":  height,
			"label":   label,
			"quality": quality,
			"width":   width,
		})
	}

	return results
}

func expandCustomPresetCodecsPngImage(input []interface{}) encodings.PngImage {
	if len(input) == 0 || input[0] == nil {
		return encodings.PngImage{}
	}

	jpgImage := input[0].(map[string]interface{})
	result := encodings.PngImage{
		Start:  jpgImage["start"].(string),
		Layers: expandCustomPresetCodecsPngImageLayer(jpgImage["layer"].([]interface{})),
	}
	if v := jpgImage["key_frame_interval"].(string); v != "" {
		result.KeyFrameInterval = utils.String(v)
	}

	if v := jpgImage["label"].(string); v != "" {
		result.Label = utils.String(v)
	}

	if v := jpgImage["range"].(string); v != "" {
		result.Range = utils.String(v)
	}

	if v := jpgImage["step"].(string); v != "" {
		result.Step = utils.String(v)
	}

	if v := jpgImage["stretch_mode"].(string); v != "" {
		result.StretchMode = pointer.To(encodings.StretchMode(v))
	}

	if v := jpgImage["sync_mode"].(string); v != "" {
		result.SyncMode = pointer.To(encodings.VideoSyncMode(v))
	}

	return result
}

func flattenCustomPresetCodecsPngImage(input encodings.PngImage) []interface{} {
	keyFrameInterval := ""
	if input.KeyFrameInterval != nil {
		keyFrameInterval = *input.KeyFrameInterval
	}

	label := ""
	if input.Label != nil {
		label = *input.Label
	}

	rang := ""
	if input.Range != nil {
		rang = *input.Range
	}

	step := ""
	if input.Step != nil {
		step = *input.Step
	}

	stretchMode := ""
	if input.StretchMode != nil {
		stretchMode = string(*input.StretchMode)
	}

	syncMode := ""
	if input.SyncMode != nil {
		syncMode = string(*input.SyncMode)
	}

	return []interface{}{
		map[string]interface{}{
			"key_frame_interval": keyFrameInterval,
			"label":              label,
			"layer":              flattenCustomPresetCodecsPngImageLayer(input.Layers),
			"range":              rang,
			"start":              input.Start,
			"step":               step,
			"stretch_mode":       stretchMode,
			"sync_mode":          syncMode,
		},
	}
}

func expandCustomPresetCodecsPngImageLayer(input []interface{}) *[]encodings.Layer {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	results := make([]encodings.Layer, 0)
	for _, layerRaw := range input {
		if layerRaw == nil {
			continue
		}

		layer := layerRaw.(map[string]interface{})
		result := encodings.Layer{}

		if v := layer["height"].(string); v != "" {
			result.Height = utils.String(v)
		}

		if v := layer["label"].(string); v != "" {
			result.Label = utils.String(v)
		}

		if v := layer["width"].(string); v != "" {
			result.Width = utils.String(v)
		}

		results = append(results, result)
	}

	return &results
}

func flattenCustomPresetCodecsPngImageLayer(input *[]encodings.Layer) []interface{} {
	if input == nil {
		return make([]interface{}, 0)
	}
	results := make([]interface{}, 0)

	for _, v := range *input {
		height := ""
		if v.Height != nil {
			height = *v.Height
		}

		label := ""
		if v.Label != nil {
			label = *v.Label
		}

		width := ""
		if v.Width != nil {
			width = *v.Width
		}

		results = append(results, map[string]interface{}{
			"height": height,
			"label":  label,
			"width":  width,
		})
	}

	return results
}

func expandCustomPresetFilters(input []interface{}) (*encodings.Filters, error) {
	if len(input) == 0 || input[0] == nil {
		return nil, nil
	}

	filters := input[0].(map[string]interface{})

	result := encodings.Filters{}

	if overlayRaw, ok := filters["overlay"].([]interface{}); ok {
		overlay, err := expandCustomPresetFiltersOverlays(overlayRaw)
		if err != nil {
			return nil, err
		}
		result.Overlays = overlay
	}

	if cropRectangle, ok := filters["crop_rectangle"].([]interface{}); ok {
		result.Crop = expandCustomPresetFiltersCropRectangle(cropRectangle)
	}

	if deinterlace, ok := filters["deinterlace"].([]interface{}); ok {
		result.Deinterlace = expandCustomPresetFiltersDeinterlace(deinterlace)
	}

	if fadeIn, ok := filters["fade_in"].([]interface{}); ok {
		result.FadeIn = expandCustomPresetFiltersFade(fadeIn)
	}

	if fadeOut, ok := filters["fade_out"].([]interface{}); ok {
		result.FadeOut = expandCustomPresetFiltersFade(fadeOut)
	}

	if v := filters["rotation"].(string); v != "" {
		result.Rotation = pointer.To(encodings.Rotation(v))
	}

	return &result, nil
}

func flattenCustomPresetFilters(input *encodings.Filters) []interface{} {
	if input == nil {
		return make([]interface{}, 0)
	}

	rotation := ""
	if input.Rotation != nil {
		rotation = string(*input.Rotation)
	}

	return []interface{}{
		map[string]interface{}{
			"crop_rectangle": flattenCustomPresetFilterCropRectangle(input.Crop),
			"deinterlace":    flattenCustomPresetFilterDeinterlace(input.Deinterlace),
			"fade_in":        flattenCustomPresetFilterFade(input.FadeIn),
			"fade_out":       flattenCustomPresetFilterFade(input.FadeOut),
			"overlay":        flattenCustomPresetFilterOverlays(input.Overlays),
			"rotation":       rotation,
		},
	}
}

func expandCustomPresetFiltersCropRectangle(input []interface{}) *encodings.Rectangle {
	if input == nil || input[0] == nil {
		return nil
	}

	cropRectangle := input[0].(map[string]interface{})
	result := encodings.Rectangle{}

	if v := cropRectangle["height"].(string); v != "" {
		result.Height = utils.String(v)
	}

	if v := cropRectangle["left"].(string); v != "" {
		result.Left = utils.String(v)
	}

	if v := cropRectangle["top"].(string); v != "" {
		result.Top = utils.String(v)
	}

	if v := cropRectangle["width"].(string); v != "" {
		result.Width = utils.String(v)
	}

	return &result
}

func flattenCustomPresetFilterCropRectangle(input *encodings.Rectangle) []interface{} {
	if input == nil {
		return make([]interface{}, 0)
	}

	height := ""
	if input.Height != nil {
		height = *input.Height
	}

	left := ""
	if input.Left != nil {
		left = *input.Left
	}

	top := ""
	if input.Top != nil {
		top = *input.Top
	}

	width := ""
	if input.Width != nil {
		width = *input.Width
	}

	return []interface{}{
		map[string]interface{}{
			"height": height,
			"left":   left,
			"top":    top,
			"width":  width,
		},
	}
}

func expandCustomPresetFiltersDeinterlace(input []interface{}) *encodings.Deinterlace {
	if input == nil || input[0] == nil {
		return nil
	}

	crop := input[0].(map[string]interface{})
	result := encodings.Deinterlace{}

	if v := crop["parity"].(string); v != "" {
		result.Parity = pointer.To(encodings.DeinterlaceParity(v))
	}

	if v := crop["mode"].(string); v != "" {
		result.Mode = pointer.To(encodings.DeinterlaceMode(v))
	}

	return &result
}

func flattenCustomPresetFilterDeinterlace(input *encodings.Deinterlace) []interface{} {
	if input == nil {
		return make([]interface{}, 0)
	}

	parity := ""
	if input.Parity != nil {
		parity = string(*input.Parity)
	}

	mode := ""
	if input.Mode != nil {
		mode = string(*input.Mode)
	}

	return []interface{}{
		map[string]interface{}{
			"parity": parity,
			"mode":   mode,
		},
	}
}

func expandCustomPresetFiltersFade(input []interface{}) *encodings.Fade {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	v := input[0].(map[string]interface{})
	result := encodings.Fade{
		Duration:  v["duration"].(string),
		FadeColor: v["fade_color"].(string),
	}

	if start := v["start"].(string); start != "" {
		result.Start = utils.String(start)
	}

	return &result
}

func flattenCustomPresetFilterFade(input *encodings.Fade) []interface{} {
	if input == nil {
		return make([]interface{}, 0)
	}

	start := ""
	if input.Start != nil {
		start = *input.Start
	}

	return []interface{}{
		map[string]interface{}{
			"duration":   input.Duration,
			"fade_color": input.FadeColor,
			"start":      start,
		},
	}
}

func expandCustomPresetFiltersOverlays(input []interface{}) (*[]encodings.Overlay, error) {
	if len(input) == 0 || input[0] == nil {
		return nil, nil
	}

	results := make([]encodings.Overlay, 0)

	for _, v := range input {
		if v == nil {
			continue
		}
		overlay := v.(map[string]interface{})
		audio := overlay["audio"].([]interface{})
		video := overlay["video"].([]interface{})

		overlaysCount := 0
		if len(audio) > 0 {
			overlaysCount++
		}
		if len(video) > 0 {
			overlaysCount++
		}
		if overlaysCount == 0 {
			return nil, fmt.Errorf("custom preset filter overlay must contain at least one type of: audio or video")
		}
		if overlaysCount > 1 {
			return nil, fmt.Errorf("more than one type of overlay in the same custom preset filter overlay is not allowed")
		}

		if len(audio) > 0 {
			results = append(results, expandCustomPresetFiltersOverlaysAudio(audio))
		} else if len(video) > 0 {
			results = append(results, expandCustomPresetFiltersOverlaysVideo(video))
		}
	}

	return &results, nil
}

type flattenedCustomPresetFilterOverlay struct {
	audio []interface{}
	video []interface{}
}

func flattenCustomPresetFilterOverlays(input *[]encodings.Overlay) []interface{} {
	if input == nil {
		return make([]interface{}, 0)
	}

	results := make([]interface{}, 0)

	for _, v := range *input {
		result := flattenedCustomPresetFilterOverlay{
			audio: []interface{}{},
			video: []interface{}{},
		}

		if overlay, ok := v.(encodings.AudioOverlay); ok {
			result.audio = flattenCustomPresetFilterOverlayAudio(overlay)
		}

		if overlay, ok := v.(encodings.VideoOverlay); ok {
			result.video = flattenCustomPresetFilterOverlayVideo(overlay)
		}

		results = append(results, map[string]interface{}{
			"audio": result.audio,
			"video": result.video,
		})
	}

	return results
}

func expandCustomPresetFiltersOverlaysAudio(input []interface{}) encodings.AudioOverlay {
	if len(input) == 0 || input[0] == nil {
		return encodings.AudioOverlay{}
	}

	audio := input[0].(map[string]interface{})
	result := encodings.AudioOverlay{
		InputLabel: audio["input_label"].(string),
	}

	if v := audio["audio_gain_level"].(float64); v != 0 {
		result.AudioGainLevel = utils.Float(v)
	}

	if v := audio["end"].(string); v != "" {
		result.End = utils.String(v)
	}

	if v := audio["fade_in_duration"].(string); v != "" {
		result.FadeInDuration = utils.String(v)
	}

	if v := audio["fade_out_duration"].(string); v != "" {
		result.FadeOutDuration = utils.String(v)
	}

	if v := audio["start"].(string); v != "" {
		result.Start = utils.String(v)
	}

	return result
}

func flattenCustomPresetFilterOverlayAudio(input encodings.AudioOverlay) []interface{} {
	audioGainLevel := 0.0
	if input.AudioGainLevel != nil {
		audioGainLevel = *input.AudioGainLevel
	}

	end := ""
	if input.End != nil {
		end = *input.End
	}

	fadeInDuration := ""
	if input.FadeInDuration != nil {
		fadeInDuration = *input.FadeInDuration
	}

	fadeOutDuration := ""
	if input.FadeOutDuration != nil {
		fadeOutDuration = *input.FadeOutDuration
	}

	start := ""
	if input.Start != nil {
		start = *input.Start
	}

	return []interface{}{
		map[string]interface{}{
			"audio_gain_level":  audioGainLevel,
			"end":               end,
			"fade_in_duration":  fadeInDuration,
			"fade_out_duration": fadeOutDuration,
			"input_label":       input.InputLabel,
			"start":             start,
		},
	}
}

func expandCustomPresetFiltersOverlaysVideo(input []interface{}) encodings.VideoOverlay {
	if len(input) == 0 || input[0] == nil {
		return encodings.VideoOverlay{}
	}

	video := input[0].(map[string]interface{})
	result := encodings.VideoOverlay{
		InputLabel:    video["input_label"].(string),
		Position:      expandCustomPresetFiltersCropRectangle(video["position"].([]interface{})),
		CropRectangle: expandCustomPresetFiltersCropRectangle(video["crop_rectangle"].([]interface{})),
	}

	if v := video["audio_gain_level"].(float64); v != 0 {
		result.AudioGainLevel = utils.Float(v)
	}

	if v := video["end"].(string); v != "" {
		result.End = utils.String(v)
	}

	if v := video["fade_in_duration"].(string); v != "" {
		result.FadeInDuration = utils.String(v)
	}

	if v := video["fade_out_duration"].(string); v != "" {
		result.FadeOutDuration = utils.String(v)
	}

	if v := video["opacity"].(float64); v != 0 {
		result.Opacity = utils.Float(v)
	}

	if v := video["start"].(string); v != "" {
		result.Start = utils.String(v)
	}

	return result
}

func flattenCustomPresetFilterOverlayVideo(input encodings.VideoOverlay) []interface{} {
	audioGainLevel := 0.0
	if input.AudioGainLevel != nil {
		audioGainLevel = *input.AudioGainLevel
	}

	end := ""
	if input.End != nil {
		end = *input.End
	}

	fadeInDuration := ""
	if input.FadeInDuration != nil {
		fadeInDuration = *input.FadeInDuration
	}

	fadeOutDuration := ""
	if input.FadeOutDuration != nil {
		fadeOutDuration = *input.FadeOutDuration
	}

	opacity := 0.0
	if input.Opacity != nil {
		opacity = *input.Opacity
	}

	start := ""
	if input.Start != nil {
		start = *input.Start
	}

	return []interface{}{
		map[string]interface{}{
			"audio_gain_level":  audioGainLevel,
			"crop_rectangle":    flattenCustomPresetFilterCropRectangle(input.CropRectangle),
			"end":               end,
			"fade_in_duration":  fadeInDuration,
			"fade_out_duration": fadeOutDuration,
			"input_label":       input.InputLabel,
			"opacity":           opacity,
			"position":          flattenCustomPresetFilterCropRectangle(input.Position),
			"start":             start,
		},
	}
}

func expandCustomPresetFormats(input []interface{}) ([]encodings.Format, error) {
	if len(input) == 0 || input[0] == nil {
		return make([]encodings.Format, 0), nil
	}

	results := make([]encodings.Format, 0)
	for _, v := range input {
		if v == nil {
			continue
		}

		format := v.(map[string]interface{})

		jpg := format["jpg"].([]interface{})
		mp4 := format["mp4"].([]interface{})
		png := format["png"].([]interface{})
		transportStream := format["transport_stream"].([]interface{})

		formatCount := 0

		if len(jpg) > 0 {
			formatCount++
		}
		if len(mp4) > 0 {
			formatCount++
		}
		if len(png) > 0 {
			formatCount++
		}
		if len(transportStream) > 0 {
			formatCount++
		}

		if formatCount == 0 {
			return nil, fmt.Errorf("custom preset format must contain at least one type of: jpg, mp4, png or transport_stream")
		}
		if formatCount > 1 {
			return nil, fmt.Errorf("more than one type of format in the same custom preset format is not allowed")
		}

		if len(jpg) > 0 {
			results = append(results, expandCustomPresetFormatsJpg(jpg))
		}
		if len(mp4) > 0 {
			results = append(results, expandCustomPresetFormatsMp4(mp4))
		}
		if len(png) > 0 {
			results = append(results, expandCustomPresetFormatsPng(png))
		}
		if len(transportStream) > 0 {
			results = append(results, expandCustomPresetFormatsTransportStream(transportStream))
		}
	}

	return results, nil
}

type flattenedCustomPresetFormat struct {
	jpg             []interface{}
	mp4             []interface{}
	png             []interface{}
	transportStream []interface{}
}

func flattenCustomPresetFormats(input []encodings.Format) []interface{} {
	results := make([]interface{}, 0)

	for _, v := range input {
		result := flattenedCustomPresetFormat{
			jpg:             []interface{}{},
			mp4:             []interface{}{},
			png:             []interface{}{},
			transportStream: []interface{}{},
		}

		if format, ok := v.(encodings.JpgFormat); ok {
			result.jpg = flattenCustomPresetFormatsJpg(format)
		}
		if format, ok := v.(encodings.Mp4Format); ok {
			result.mp4 = flattenCustomPresetFormatsMp4(format)
		}
		if format, ok := v.(encodings.PngFormat); ok {
			result.png = flattenCustomPresetFormatsPng(format)
		}
		if format, ok := v.(encodings.TransportStreamFormat); ok {
			result.transportStream = flattenCustomPresetFormatsTransportStream(format)
		}

		results = append(results, map[string]interface{}{
			"jpg":              result.jpg,
			"mp4":              result.mp4,
			"png":              result.png,
			"transport_stream": result.transportStream,
		})
	}

	return results
}

func expandCustomPresetFormatsJpg(input []interface{}) encodings.JpgFormat {
	if len(input) == 0 || input[0] == nil {
		return encodings.JpgFormat{}
	}

	jpg := input[0].(map[string]interface{})

	result := encodings.JpgFormat{
		FilenamePattern: jpg["filename_pattern"].(string),
	}

	return result
}

func flattenCustomPresetFormatsJpg(input encodings.JpgFormat) []interface{} {
	return []interface{}{
		map[string]interface{}{
			"filename_pattern": input.FilenamePattern,
		},
	}
}

func expandCustomPresetFormatsMp4(input []interface{}) encodings.Mp4Format {
	if len(input) == 0 || input[0] == nil {
		return encodings.Mp4Format{}
	}

	mp4 := input[0].(map[string]interface{})
	result := encodings.Mp4Format{
		FilenamePattern: mp4["filename_pattern"].(string),
		OutputFiles:     expandCustomPresetFormatsOutputFiles(mp4["output_file"].([]interface{})),
	}

	return result
}

func flattenCustomPresetFormatsMp4(input encodings.Mp4Format) []interface{} {
	return []interface{}{
		map[string]interface{}{
			"filename_pattern": input.FilenamePattern,
			"output_file":      flattenCustomPresetFormatOutputFiles(input.OutputFiles),
		},
	}
}

func expandCustomPresetFormatsPng(input []interface{}) encodings.PngFormat {
	if len(input) == 0 || input[0] == nil {
		return encodings.PngFormat{}
	}

	jpg := input[0].(map[string]interface{})
	result := encodings.PngFormat{
		FilenamePattern: jpg["filename_pattern"].(string),
	}

	return result
}

func flattenCustomPresetFormatsPng(input encodings.PngFormat) []interface{} {
	return []interface{}{
		map[string]interface{}{
			"filename_pattern": input.FilenamePattern,
		},
	}
}

func expandCustomPresetFormatsTransportStream(input []interface{}) encodings.TransportStreamFormat {
	if len(input) == 0 || input[0] == nil {
		return encodings.TransportStreamFormat{}
	}

	transportStream := input[0].(map[string]interface{})
	result := encodings.TransportStreamFormat{
		FilenamePattern: transportStream["filename_pattern"].(string),
		OutputFiles:     expandCustomPresetFormatsOutputFiles(transportStream["output_file"].([]interface{})),
	}

	return result
}

func flattenCustomPresetFormatsTransportStream(input encodings.TransportStreamFormat) []interface{} {
	return []interface{}{
		map[string]interface{}{
			"filename_pattern": input.FilenamePattern,
			"output_file":      flattenCustomPresetFormatOutputFiles(input.OutputFiles),
		},
	}
}

func expandCustomPresetFormatsOutputFiles(input []interface{}) *[]encodings.OutputFile {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	results := make([]encodings.OutputFile, 0)
	for _, v := range input {
		if v == nil {
			continue
		}

		outputFile := v.(map[string]interface{})
		labels := make([]string, 0)
		for _, label := range outputFile["labels"].([]interface{}) {
			labels = append(labels, label.(string))
		}

		results = append(results, encodings.OutputFile{
			Labels: labels,
		})
	}

	return &results
}

func flattenCustomPresetFormatOutputFiles(input *[]encodings.OutputFile) []interface{} {
	if input == nil {
		return make([]interface{}, 0)
	}

	results := make([]interface{}, 0)
	for _, v := range *input {
		labels := make([]interface{}, 0)
		for _, label := range v.Labels {
			labels = append(labels, label)
		}

		results = append(results, map[string]interface{}{
			"labels": labels,
		})
	}

	return results
}
