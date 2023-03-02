package media

import (
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/media/2021-11-01/encodings"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/media/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
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
									"mode": {
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
																			Default:      0,
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
																			Default:      0,
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
															"splite_column": {
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
															"output_files": {
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
															"output_files": {
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
									"filter": {
										Type:     pluginsdk.TypeList,
										Optional: true,
										MaxItems: 1,
										Elem: &pluginsdk.Resource{
											Schema: map[string]*schema.Schema{
												"crop": {
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
															"fate_color": {
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
															"fate_color": {
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
												"overlays": {
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
																		"inputLabel": {
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
																		"inputLabel": {
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
																		"opacity": {
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
	client := meta.(*clients.Client).Media.V20211101Client.Encodings
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
	client := meta.(*clients.Client).Media.V20211101Client.Encodings
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
	client := meta.(*clients.Client).Media.V20211101Client.Encodings
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
		return nil, fmt.Errorf("output must contain at least one type of preset: builtin_preset,face_detector_preset,video_analyzer_preset or audio_analyzer_preset")
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

		filters, err := expandCustomPresetFilters(preset["filter"].([]interface{}))
		if err != nil {
			return nil, err
		}

		formats, err := expandCustomPresetFormats(preset["format"].([]interface{}))
		if err != nil {
			return nil, err
		}
		builtInPreset := &encodings.StandardEncoderPreset{
			Codecs:  codecs,
			Filters: filters,
			Formats: formats,
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

		if v := preset["mode"].(string); v != "" {
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

	return nil, fmt.Errorf("output must contain at least one type of preset: builtin_preset,face_detector_preset,video_analyzer_preset or audio_analyzer_preset")
}

type flattenedPresets struct {
	audioAnalyzerPresets []interface{}
	builtInPresets       []interface{}
	faceDetectorPresets  []interface{}
	videoAnalyzerPresets []interface{}
}

func flattenPreset(input encodings.Preset) flattenedPresets {
	out := flattenedPresets{
		audioAnalyzerPresets: []interface{}{},
		builtInPresets:       []interface{}{},
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
			"mode":                 mode,
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

	for i, v := range input {
		value, err := tags.TagValueToString(v)
		if err != nil {
			return nil, err
		}
		output[i] = value
	}

	return &output, nil
}

func flattenExperimentalOptions(input *map[string]string) map[string]interface{} {
	result := make(map[string]interface{}, 0)
	if input == nil {
		return result
	}
	for i, v := range *input {
		result[i] = v
	}

	return result
}

func expandCustomPresetCodecs(input []interface{}) ([]encodings.Codec, error) {
	if len(input) == 0 || input[0] == nil {
		return make([]encodings.Codec, 0), nil
	}

	result := make([]encodings.Codec, 0)

	return result, nil
}

func expandCustomPresetFilters(input []interface{}) (*encodings.Filters, error) {
	if len(input) == 0 || input[0] == nil {
		return nil, nil
	}

	result := encodings.Filters{}

	return &result, nil
}

func expandCustomPresetFormats(input []interface{}) ([]encodings.Format, error) {
	if len(input) == 0 || input[0] == nil {
		return make([]encodings.Format, 0), nil
	}

	result := make([]encodings.Format, 0)

	return result, nil
}
