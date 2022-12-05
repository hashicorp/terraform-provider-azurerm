package media

import (
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/media/2020-05-01/encodings"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
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

			//lintignore:XS003
			"output": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				MinItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"on_error_action": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringInSlice(encodings.PossibleValuesForOnErrorType(), false),
						},
						//lintignore:XS003
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
										Type:         pluginsdk.TypeString,
										Optional:     true,
										ValidateFunc: validation.StringInSlice(encodings.PossibleValuesForAudioAnalysisMode(), false),
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
										Type:         pluginsdk.TypeString,
										Optional:     true,
										ValidateFunc: validation.StringInSlice(encodings.PossibleValuesForAudioAnalysisMode(), false),
									},
									"insights_type": {
										Type:         pluginsdk.TypeString,
										Optional:     true,
										ValidateFunc: validation.StringInSlice(encodings.PossibleValuesForInsightsType(), false),
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
										Type:         pluginsdk.TypeString,
										Optional:     true,
										ValidateFunc: validation.StringInSlice(encodings.PossibleValuesForAnalysisResolution(), false),
									},
								},
							},
						},
						"relative_priority": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringInSlice(encodings.PossibleValuesForPriority(), false),
						},
					},
				},
			},
		},
	}
}

func resourceMediaTransformCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Media.V20200501Client.Encodings
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
	client := meta.(*clients.Client).Media.V20200501Client.Encodings
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
	d.Set("media_services_account_name", id.AccountName)
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
	client := meta.(*clients.Client).Media.V20200501Client.Encodings
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

		if transform["on_error_action"] != nil {
			transformOutput.OnError = pointer.To(encodings.OnErrorType(transform["on_error_action"].(string)))
		}

		if transform["relative_priority"] != nil {
			transformOutput.RelativePriority = pointer.To(encodings.Priority(transform["relative_priority"].(string)))
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
	faceDetectorPresets := transform["face_detector_preset"].([]interface{})
	videoAnalyzerPresets := transform["video_analyzer_preset"].([]interface{})

	presetsCount := 0
	if len(audioAnalyzerPresets) > 0 {
		presetsCount++
	}
	if len(builtInPresets) > 0 {
		presetsCount++
	}
	if len(faceDetectorPresets) > 0 {
		presetsCount++
	}
	if len(videoAnalyzerPresets) > 0 {
		presetsCount++
	}
	if presetsCount == 0 {
		return nil, fmt.Errorf("output must contain at least one type of preset: builtin_preset,face_detector_preset,video_analyzer_preset or audio_analyzer_preset.")
	}
	if presetsCount > 1 {
		return nil, fmt.Errorf("more than one type of preset in the same output is not allowed.")
	}

	if len(audioAnalyzerPresets) > 0 {
		preset := audioAnalyzerPresets[0].(map[string]interface{})
		audioAnalyzerPreset := &encodings.AudioAnalyzerPreset{}
		if preset["audio_language"] != nil && preset["audio_language"].(string) != "" {
			audioAnalyzerPreset.AudioLanguage = utils.String(preset["audio_language"].(string))
		}
		if preset["audio_analysis_mode"] != nil {
			audioAnalyzerPreset.Mode = pointer.To(encodings.AudioAnalysisMode(preset["audio_analysis_mode"].(string)))
		}
		return audioAnalyzerPreset, nil
	}

	if len(builtInPresets) > 0 {
		preset := builtInPresets[0].(map[string]interface{})
		presetName := preset["preset_name"].(string)
		builtInPreset := &encodings.BuiltInStandardEncoderPreset{
			PresetName: encodings.EncoderNamedPreset(presetName),
		}
		return builtInPreset, nil
	}

	if len(faceDetectorPresets) > 0 {
		preset := faceDetectorPresets[0].(map[string]interface{})
		faceDetectorPreset := &encodings.FaceDetectorPreset{}
		if preset["analysis_resolution"] != nil {
			faceDetectorPreset.Resolution = pointer.To(encodings.AnalysisResolution(preset["analysis_resolution"].(string)))
		}
		return faceDetectorPreset, nil
	}

	if len(videoAnalyzerPresets) > 0 {
		presets := transform["video_analyzer_preset"].([]interface{})
		preset := presets[0].(map[string]interface{})
		videoAnalyzerPreset := &encodings.VideoAnalyzerPreset{}
		if preset["audio_language"] != nil {
			videoAnalyzerPreset.AudioLanguage = utils.String(preset["audio_language"].(string))
		}
		if preset["audio_analysis_mode"] != nil {
			videoAnalyzerPreset.Mode = pointer.To(encodings.AudioAnalysisMode(preset["audio_analysis_mode"].(string)))
		}
		if preset["insights_type"] != nil {
			videoAnalyzerPreset.InsightsToExtract = pointer.To(encodings.InsightsType(preset["insights_type"].(string)))
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
			"audio_analysis_mode": mode,
			"audio_language":      language,
		})
	}

	if v, ok := input.(encodings.BuiltInStandardEncoderPreset); ok {
		out.builtInPresets = append(out.builtInPresets, map[string]interface{}{
			"preset_name": string(v.PresetName),
		})
	}

	if v, ok := input.(encodings.FaceDetectorPreset); ok {
		resolution := ""
		if v.Resolution != nil {
			resolution = string(*v.Resolution)
		}
		out.faceDetectorPresets = append(out.faceDetectorPresets, map[string]interface{}{
			"analysis_resolution": resolution,
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
			"audio_analysis_mode": mode,
			"audio_language":      audioLanguage,
			"insights_type":       insightsType,
		})
	}

	return out
}
