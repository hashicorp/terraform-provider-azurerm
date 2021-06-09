package media

import (
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/mediaservices/mgmt/2021-05-01/media"
	"github.com/hashicorp/go-azure-helpers/response"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/media/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
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
			_, err := parse.TransformID(id)
			return err
		}),

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

			"resource_group_name": azure.SchemaResourceGroupName(),

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

			//lintignore:XS003
			"output": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				MinItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"on_error_action": {
							Type:     pluginsdk.TypeString,
							Optional: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(media.OnErrorTypeContinueJob),
								string(media.OnErrorTypeStopProcessingJob),
							}, false),
						},
						//lintignore:XS003
						"builtin_preset": {
							Type:     pluginsdk.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"preset_name": {
										Type:     pluginsdk.TypeString,
										Optional: true,
										ValidateFunc: validation.StringInSlice([]string{
											string(media.EncoderNamedPresetAACGoodQualityAudio),
											string(media.EncoderNamedPresetAdaptiveStreaming),
											string(media.EncoderNamedPresetContentAwareEncoding),
											string(media.EncoderNamedPresetContentAwareEncodingExperimental),
											string(media.EncoderNamedPresetCopyAllBitrateNonInterleaved),
											string(media.EncoderNamedPresetH264MultipleBitrate1080p),
											string(media.EncoderNamedPresetH264MultipleBitrate720p),
											string(media.EncoderNamedPresetH264MultipleBitrateSD),
											string(media.EncoderNamedPresetH264SingleBitrate1080p),
											string(media.EncoderNamedPresetH264SingleBitrate720p),
											string(media.EncoderNamedPresetH264SingleBitrateSD),
										}, false),
									},
								},
							},
						},
						//lintignore:XS003
						"audio_analyzer_preset": {
							Type:     pluginsdk.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"audio_language": {
										Type:     pluginsdk.TypeString,
										Optional: true,
										ValidateFunc: validation.StringInSlice([]string{
											"ar-EG",
											"ar-SY",
											"de-DE",
											"en-AU",
											"en-GB",
											"en-US",
											"es-ES",
											"es-MX",
											"fr-FR",
											"hi-IN",
											"it-IT",
											"ja-JP",
											"ko-KR",
											"pt-BR",
											"ru-RU",
											"zh-CN",
										}, false),
									},
									"audio_analysis_mode": {
										Type:     pluginsdk.TypeString,
										Optional: true,
										ValidateFunc: validation.StringInSlice([]string{
											string(media.AudioAnalysisModeBasic),
											string(media.AudioAnalysisModeStandard),
										}, false),
									},
								},
							},
						},
						//lintignore:XS003
						"video_analyzer_preset": {
							Type:     pluginsdk.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"audio_language": {
										Type:     pluginsdk.TypeString,
										Optional: true,
										ValidateFunc: validation.StringInSlice([]string{
											"ar-EG",
											"ar-SY",
											"de-DE",
											"en-AU",
											"en-GB",
											"en-US",
											"es-ES",
											"es-MX",
											"fr-FR",
											"hi-IN",
											"it-IT",
											"ja-JP",
											"ko-KR",
											"pt-BR",
											"ru-RU",
											"zh-CN",
										}, false),
									},
									"audio_analysis_mode": {
										Type:     pluginsdk.TypeString,
										Optional: true,
										ValidateFunc: validation.StringInSlice([]string{
											string(media.AudioAnalysisModeBasic),
											string(media.AudioAnalysisModeStandard),
										}, false),
									},
									"insights_type": {
										Type:     pluginsdk.TypeString,
										Optional: true,
										ValidateFunc: validation.StringInSlice([]string{
											string(media.InsightsTypeAllInsights),
											string(media.InsightsTypeAudioInsightsOnly),
											string(media.InsightsTypeVideoInsightsOnly),
										}, false),
									},
								},
							},
						},
						//lintignore:XS003
						"face_detector_preset": {
							Type:     pluginsdk.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"analysis_resolution": {
										Type:     pluginsdk.TypeString,
										Optional: true,
										ValidateFunc: validation.StringInSlice([]string{
											string(media.AnalysisResolutionSourceResolution),
											string(media.AnalysisResolutionStandardDefinition),
										}, false),
									},
								},
							},
						},
						"relative_priority": {
							Type:     pluginsdk.TypeString,
							Optional: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(media.PriorityHigh),
								string(media.PriorityNormal),
								string(media.PriorityLow),
							}, false),
						},
					},
				},
			},
		},
	}
}

func resourceMediaTransformCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Media.TransformsClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	resourceId := parse.NewTransformID(subscriptionId, d.Get("resource_group_name").(string), d.Get("media_services_account_name").(string), d.Get("name").(string))
	if d.IsNewResource() {
		existing, err := client.Get(ctx, resourceId.ResourceGroup, resourceId.MediaserviceName, resourceId.Name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for existing %s: %+v", resourceId, err)
			}
		}

		if !utils.ResponseWasNotFound(existing.Response) {
			return tf.ImportAsExistsError("azurerm_media_transform", resourceId.ID())
		}
	}

	parameters := media.Transform{
		TransformProperties: &media.TransformProperties{
			Description: utils.String(d.Get("description").(string)),
		},
	}

	if v, ok := d.GetOk("output"); ok {
		transformOutput, err := expandTransformOuputs(v.([]interface{}))
		if err != nil {
			return err
		}
		parameters.Outputs = transformOutput
	}

	if _, err := client.CreateOrUpdate(ctx, resourceId.ResourceGroup, resourceId.MediaserviceName, resourceId.Name, parameters); err != nil {
		return fmt.Errorf("creating/updating %s: %+v", resourceId, err)
	}

	d.SetId(resourceId.ID())
	return resourceMediaTransformRead(d, meta)
}

func resourceMediaTransformRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Media.TransformsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.TransformID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.MediaserviceName, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] Transform %q was not found in Media Services Account %q and Resource Group %q - removing from state", id.Name, id.MediaserviceName, id.ResourceGroup)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error retrieving Transform %q in Media Services Account %q (Resource Group %q): %+v", id.Name, id.MediaserviceName, id.ResourceGroup, err)
	}

	d.Set("name", id.Name)
	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("media_services_account_name", id.MediaserviceName)

	if props := resp.TransformProperties; props != nil {
		if description := props.Description; description != nil {
			d.Set("description", description)
		}

		outputs := flattenTransformOutputs(props.Outputs)
		if err := d.Set("output", outputs); err != nil {
			return fmt.Errorf("Error flattening `output`: %s", err)
		}
	}

	return nil
}

func resourceMediaTransformDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Media.TransformsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.TransformID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Delete(ctx, id.ResourceGroup, id.MediaserviceName, id.Name)
	if err != nil {
		if response.WasNotFound(resp.Response) {
			return nil
		}
		return fmt.Errorf("Error deleting Transform %q in Media Services Account %q (Resource Group %q): %+v", id.Name, id.MediaserviceName, id.ResourceGroup, err)
	}

	return nil
}

func expandTransformOuputs(input []interface{}) (*[]media.TransformOutput, error) {
	results := make([]media.TransformOutput, 0)

	for _, transformOuputRaw := range input {
		if transformOuputRaw == nil {
			continue
		}
		transform := transformOuputRaw.(map[string]interface{})

		preset, err := expandPreset(transform)
		if err != nil {
			return nil, err
		}

		transformOuput := media.TransformOutput{
			Preset: preset,
		}

		if transform["on_error_action"] != nil {
			transformOuput.OnError = media.OnErrorType(transform["on_error_action"].(string))
		}

		if transform["relative_priority"] != nil {
			transformOuput.RelativePriority = media.Priority(transform["relative_priority"].(string))
		}

		results = append(results, transformOuput)
	}

	return &results, nil
}

func flattenTransformOutputs(input *[]media.TransformOutput) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	results := make([]interface{}, 0)
	for _, transformOuput := range *input {
		output := make(map[string]interface{})
		output["on_error_action"] = string(transformOuput.OnError)
		output["relative_priority"] = string(transformOuput.RelativePriority)
		attribute, preset := flattenPreset(transformOuput.Preset)
		if attribute != "" {
			output[attribute] = preset
		}
		results = append(results, output)
	}

	return results
}

func expandPreset(transform map[string]interface{}) (media.BasicPreset, error) {
	presetsCount := 0
	presetType := ""
	if transform["builtin_preset"] != nil && len(transform["builtin_preset"].([]interface{})) > 0 && transform["builtin_preset"].([]interface{})[0] != nil {
		presetsCount++
		presetType = string(media.OdataTypeBasicPresetOdataTypeMicrosoftMediaBuiltInStandardEncoderPreset)
	}
	if transform["audio_analyzer_preset"] != nil && len(transform["audio_analyzer_preset"].([]interface{})) > 0 && transform["audio_analyzer_preset"].([]interface{})[0] != nil {
		presetsCount++
		presetType = string(media.OdataTypeBasicPresetOdataTypeMicrosoftMediaAudioAnalyzerPreset)
	}
	if transform["video_analyzer_preset"] != nil && len(transform["video_analyzer_preset"].([]interface{})) > 0 && transform["video_analyzer_preset"].([]interface{})[0] != nil {
		presetsCount++
		presetType = string(media.OdataTypeBasicPresetOdataTypeMicrosoftMediaVideoAnalyzerPreset)
	}
	if transform["face_detector_preset"] != nil && len(transform["face_detector_preset"].([]interface{})) > 0 && transform["face_detector_preset"].([]interface{})[0] != nil {
		presetsCount++
		presetType = string(media.OdataTypeBasicPresetOdataTypeMicrosoftMediaFaceDetectorPreset)
	}

	if presetsCount == 0 {
		return nil, fmt.Errorf("output must contain at least one type of preset: builtin_preset,face_detector_preset,video_analyzer_preset or audio_analyzer_preset.")
	}

	if presetsCount > 1 {
		return nil, fmt.Errorf("more than one type of preset in the same output is not allowed.")
	}

	switch presetType {
	case string(media.OdataTypeBasicPresetOdataTypeMicrosoftMediaBuiltInStandardEncoderPreset):
		presets := transform["builtin_preset"].([]interface{})
		preset := presets[0].(map[string]interface{})
		if preset["preset_name"] == nil {
			return nil, fmt.Errorf("preset_name is required for BuiltInStandardEncoderPreset")
		}
		presetName := preset["preset_name"].(string)
		builtInPreset := &media.BuiltInStandardEncoderPreset{
			PresetName: media.EncoderNamedPreset(presetName),
			OdataType:  media.OdataTypeBasicPresetOdataTypeMicrosoftMediaBuiltInStandardEncoderPreset,
		}
		return builtInPreset, nil
	case string(media.OdataTypeBasicPresetOdataTypeMicrosoftMediaAudioAnalyzerPreset):
		presets := transform["audio_analyzer_preset"].([]interface{})
		preset := presets[0].(map[string]interface{})
		audioAnalyzerPreset := &media.AudioAnalyzerPreset{
			OdataType: media.OdataTypeBasicPresetOdataTypeMicrosoftMediaAudioAnalyzerPreset,
		}
		if preset["audio_language"] != nil && preset["audio_language"].(string) != "" {
			audioAnalyzerPreset.AudioLanguage = utils.String(preset["audio_language"].(string))
		}
		if preset["audio_analysis_mode"] != nil {
			audioAnalyzerPreset.Mode = media.AudioAnalysisMode(preset["audio_analysis_mode"].(string))
		}
		return audioAnalyzerPreset, nil
	case string(media.OdataTypeBasicPresetOdataTypeMicrosoftMediaFaceDetectorPreset):
		presets := transform["face_detector_preset"].([]interface{})
		preset := presets[0].(map[string]interface{})
		faceDetectorPreset := &media.FaceDetectorPreset{
			OdataType: media.OdataTypeBasicPresetOdataTypeMicrosoftMediaFaceDetectorPreset,
		}
		if preset["analysis_resolution"] != nil {
			faceDetectorPreset.Resolution = media.AnalysisResolution(preset["analysis_resolution"].(string))
		}
		return faceDetectorPreset, nil
	case string(media.OdataTypeBasicPresetOdataTypeMicrosoftMediaVideoAnalyzerPreset):
		presets := transform["video_analyzer_preset"].([]interface{})
		preset := presets[0].(map[string]interface{})
		videoAnalyzerPreset := &media.VideoAnalyzerPreset{
			OdataType: media.OdataTypeBasicPresetOdataTypeMicrosoftMediaVideoAnalyzerPreset,
		}
		if preset["audio_language"] != nil {
			videoAnalyzerPreset.AudioLanguage = utils.String(preset["audio_language"].(string))
		}
		if preset["audio_analysis_mode"] != nil {
			videoAnalyzerPreset.Mode = media.AudioAnalysisMode(preset["audio_analysis_mode"].(string))
		}
		if preset["insights_type"] != nil {
			videoAnalyzerPreset.InsightsToExtract = media.InsightsType(preset["insights_type"].(string))
		}
		return videoAnalyzerPreset, nil
	default:
		return nil, fmt.Errorf("output must contain at least one type of preset: builtin_preset,face_detector_preset,video_analyzer_preset or audio_analyzer_preset")
	}
}

func flattenPreset(preset media.BasicPreset) (string, []interface{}) {
	if preset == nil {
		return "", []interface{}{}
	}

	results := make([]interface{}, 0)
	result := make(map[string]interface{})
	switch preset.(type) {
	case media.AudioAnalyzerPreset:
		mediaAudioAnalyzerPreset, _ := preset.AsAudioAnalyzerPreset()
		result["audio_analysis_mode"] = string(mediaAudioAnalyzerPreset.Mode)
		if mediaAudioAnalyzerPreset.AudioLanguage != nil {
			result["audio_language"] = mediaAudioAnalyzerPreset.AudioLanguage
		}
		results = append(results, result)
		return "audio_analyzer_preset", results
	case media.BuiltInStandardEncoderPreset:
		builtInStandardEncoderPreset, _ := preset.AsBuiltInStandardEncoderPreset()
		result["preset_name"] = string(builtInStandardEncoderPreset.PresetName)
		results = append(results, result)
		return "builtin_preset", results
	case media.FaceDetectorPreset:
		faceDetectorPreset, _ := preset.AsFaceDetectorPreset()
		result["analysis_resolution"] = string(faceDetectorPreset.Resolution)
		results = append(results, result)
		return "face_detector_preset", results
	case media.VideoAnalyzerPreset:
		videoAnalyzerPreset, _ := preset.AsVideoAnalyzerPreset()
		result["audio_analysis_mode"] = string(videoAnalyzerPreset.Mode)
		result["insights_type"] = string(videoAnalyzerPreset.InsightsToExtract)
		if videoAnalyzerPreset.AudioLanguage != nil {
			result["audio_language"] = videoAnalyzerPreset.AudioLanguage
		}
		results = append(results, result)
		return "video_analyzer_preset", results
	}

	return "", results
}
