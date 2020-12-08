package media

import (
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/mediaservices/mgmt/2020-05-01/media"
	"github.com/hashicorp/go-azure-helpers/response"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azuread/azuread/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/media/parse"
	azSchema "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceMediaTransform() *schema.Resource {
	return &schema.Resource{
		Create: resourceMediaTransformCreateUpdate,
		Read:   resourceMediaTransformRead,
		Update: resourceMediaTransformCreateUpdate,
		Delete: resourceMediaTransformDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Importer: azSchema.ValidateResourceIDPriorToImport(func(id string) error {
			_, err := parse.TransformID(id)
			return err
		}),

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringMatch(
					regexp.MustCompile("^[-a-zA-Z0-9(_)]{1,128}$"),
					"Transform name must be 1 - 128 characters long, can contain letters, numbers, underscores, and hyphens (but the first and last character must be a letter or number).",
				),
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"media_services_account_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringMatch(
					regexp.MustCompile("^[-a-z0-9]{3,24}$"),
					"Media Services Account name must be 3 - 24 characters long, contain only lowercase letters and numbers.",
				),
			},

			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validate.NoEmptyStrings,
			},

			"output": {
				Type:     schema.TypeSet,
				Optional: true,
				MinItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"on_error_action": {
							Type:     schema.TypeString,
							Optional: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(media.ContinueJob), string(media.StopProcessingJob),
							}, true),
						},
						"preset": {
							Type:     schema.TypeSet,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": {
										Type:     schema.TypeString,
										Optional: true,
										ValidateFunc: validation.StringInSlice([]string{
											"BuiltInStandardEncoderPreset", "AudioAnalyzerPreset",
											"VideoAnalyzerPreset", "FaceDetectorPreset",
										}, true),
									},

									"preset_name": {
										Type:     schema.TypeString,
										Optional: true,
										ValidateFunc: validation.StringInSlice([]string{
											string(media.AACGoodQualityAudio), string(media.AdaptiveStreaming),
											string(media.ContentAwareEncoding), string(media.ContentAwareEncodingExperimental),
											string(media.CopyAllBitrateNonInterleaved), string(media.H264MultipleBitrate1080p),
											string(media.H264MultipleBitrate720p), string(media.H264MultipleBitrateSD),
											string(media.H264SingleBitrate1080p), string(media.H264SingleBitrate720p),
											string(media.H264MultipleBitrateSD),
										}, true),
									},

									"audio_language": {
										Type:     schema.TypeString,
										Optional: true,
										ValidateFunc: validation.StringInSlice([]string{
											"ar-EG", "ar-SY", "de-DE", "en-AU", "en-GB", "en-US", "es-ES", "es-MX",
											"fr-FR", "hi-IN", "it-IT", "ja-JP", "ko-KR", "pt-BR", "ru-RU", "zh-CN",
										}, true),
									},

									"audio_analysis_mode": {
										Type:     schema.TypeString,
										Optional: true,
										ValidateFunc: validation.StringInSlice([]string{
											string(media.Basic), string(media.Standard),
										}, true),
									},
									"insights_type": {
										Type:     schema.TypeString,
										Optional: true,
										ValidateFunc: validation.StringInSlice([]string{
											string(media.AllInsights), string(media.AudioInsightsOnly), string(media.VideoInsightsOnly),
										}, true),
									},
									"analysis_resolution": {
										Type:     schema.TypeString,
										Optional: true,
										ValidateFunc: validation.StringInSlice([]string{
											string(media.SourceResolution), string(media.StandardDefinition),
										}, true),
									},
								},
							},
						},
						"relative_priority": {
							Type:     schema.TypeString,
							Optional: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(media.High), string(media.Normal), string(media.Low),
							}, true),
						},
					},
				},
			},
		},
	}
}

func resourceMediaTransformCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Media.TransformsClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	transformName := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	accountName := d.Get("media_services_account_name").(string)
	description := d.Get("description").(string)

	parameters := media.Transform{
		TransformProperties: &media.TransformProperties{
			Description: utils.String(description),
		},
	}

	if v, ok := d.GetOk("output"); ok {
		transformOutput, err := expandTransformOuputs(v.(*schema.Set).List())
		if err != nil {
			return err
		}
		parameters.Outputs = transformOutput
	}

	if _, err := client.CreateOrUpdate(ctx, resourceGroup, accountName, transformName, parameters); err != nil {
		return fmt.Errorf("Error creating Transform %q in Media Services Account %q (Resource Group %q): %+v", transformName, accountName, resourceGroup, err)
	}

	transform, err := client.Get(ctx, resourceGroup, accountName, transformName)
	if err != nil {
		return fmt.Errorf("Error retrieving Transform %q from Media Services Account %q (Resource Group %q): %+v", transformName, accountName, resourceGroup, err)
	}

	d.SetId(*transform.ID)

	return resourceMediaTransformRead(d, meta)
}

func resourceMediaTransformRead(d *schema.ResourceData, meta interface{}) error {
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

func resourceMediaTransformDelete(d *schema.ResourceData, meta interface{}) error {
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
		transform := transformOuputRaw.(map[string]interface{})

		if transform["preset"] == nil {
			return nil, fmt.Errorf("output must contain a preset property")
		}

		if len(transform["preset"].(*schema.Set).List()) == 0 {
			return nil, fmt.Errorf("output must contain a preset property")
		}
		preset, err := expandPreset(transform["preset"].(*schema.Set).List())
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
		output["preset"] = flattenPreset(transformOuput.Preset)
		results = append(results, output)
	}

	return results
}

func expandPreset(presets []interface{}) (media.BasicPreset, error) {
	preset := presets[0].(map[string]interface{})
	presetType := preset["type"].(string)
	switch presetType {
	case "BuiltInStandardEncoderPreset":
		if preset["preset_name"] == nil {
			return nil, fmt.Errorf("preset_name is required for BuiltInStandardEncoderPreset")
		}
		presetName := preset["preset_name"].(string)
		builtInPreset := &media.BuiltInStandardEncoderPreset{
			PresetName: media.EncoderNamedPreset(presetName),
			OdataType:  media.OdataTypeMicrosoftMediaBuiltInStandardEncoderPreset,
		}
		return builtInPreset, nil
	case "AudioAnalyzerPreset":
		audioAnalyzerPreset := &media.AudioAnalyzerPreset{
			OdataType: media.OdataTypeMicrosoftMediaAudioAnalyzerPreset,
		}
		if preset["audio_language"] != nil && preset["audio_language"].(string) != "" {
			audioAnalyzerPreset.AudioLanguage = utils.String(preset["audio_language"].(string))
		}
		if preset["audio_analysis_mode"] != nil {
			audioAnalyzerPreset.Mode = media.AudioAnalysisMode(preset["audio_analysis_mode"].(string))
		}
		return audioAnalyzerPreset, nil
	case "FaceDetectorPreset":
		faceDetectorPreset := &media.FaceDetectorPreset{
			OdataType: media.OdataTypeMicrosoftMediaFaceDetectorPreset,
		}
		if preset["analysis_resolution"] != nil {
			faceDetectorPreset.Resolution = media.AnalysisResolution(preset["analysis_resolution"].(string))
		}
		return faceDetectorPreset, nil
	case "VideoAnalyzerPreset":
		videoAnalyzerPreset := &media.VideoAnalyzerPreset{
			OdataType: media.OdataTypeMicrosoftMediaVideoAnalyzerPreset,
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
	}
	return nil, fmt.Errorf("type property of preset is invalid")
}

func flattenPreset(preset media.BasicPreset) []interface{} {
	if preset == nil {
		return []interface{}{}
	}

	results := make([]interface{}, 0)
	result := make(map[string]interface{})
	switch preset.(type) {
	case media.AudioAnalyzerPreset:
		mediaAudioAnalyzerPreset, _ := preset.AsAudioAnalyzerPreset()
		result["audio_analysis_mode"] = string(mediaAudioAnalyzerPreset.Mode)
		result["type"] = "AudioAnalyzerPreset"
		if mediaAudioAnalyzerPreset.AudioLanguage != nil {
			result["audio_language"] = mediaAudioAnalyzerPreset.AudioLanguage
		}
		results = append(results, result)
		return results
	case media.BuiltInStandardEncoderPreset:
		builtInStandardEncoderPreset, _ := preset.AsBuiltInStandardEncoderPreset()
		result["preset_name"] = string(builtInStandardEncoderPreset.PresetName)
		result["type"] = "BuiltInStandardEncoderPreset"
		results = append(results, result)
		return results
	case media.FaceDetectorPreset:
		faceDetectorPreset, _ := preset.AsFaceDetectorPreset()
		result["analysis_resolution"] = string(faceDetectorPreset.Resolution)
		result["type"] = "FaceDetectorPreset"
		results = append(results, result)
		return results
	case media.VideoAnalyzerPreset:
		videoAnalyzerPreset, _ := preset.AsVideoAnalyzerPreset()
		result["audio_analysis_mode"] = string(videoAnalyzerPreset.Mode)
		result["insights_type"] = string(videoAnalyzerPreset.InsightsToExtract)
		result["type"] = "VideoAnalyzerPreset"
		if videoAnalyzerPreset.AudioLanguage != nil {
			result["audio_language"] = videoAnalyzerPreset.AudioLanguage
		}
		results = append(results, result)
		return results
	}

	return results
}
