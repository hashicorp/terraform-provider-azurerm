// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package bot_channel_direct_line_speech

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cognitive/2025-06-01/cognitiveservicesaccounts"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/bot/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func botChannelDirectLineSpeechSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"bot_name": {
			Description:  "The name of the Bot Resource this channel will be associated with.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.BotName,
		},

		"cognitive_service_access_key": {
			Description:  "The access key to access the Cognitive Service.",
			Type:         schema.TypeString,
			Optional:     true,
			Sensitive:    true,
			ValidateFunc: validation.StringIsNotEmpty,
			RequiredWith: []string{"cognitive_service_location"},
			AtLeastOneOf: []string{"cognitive_service_access_key", "cognitive_account_id"},
		},

		"cognitive_service_location": {
			Description:      "The supported Azure location where the Cognitive Service resource exists. Either `cognitive_account_id` or both `cognitive_service_access_key` and `cognitive_service_location` must be specified.",
			Type:             schema.TypeString,
			Optional:         true,
			ValidateFunc:     location.EnhancedValidate,
			StateFunc:        location.StateFunc,
			DiffSuppressFunc: location.DiffSuppressFunc,
			RequiredWith:     []string{"cognitive_service_access_key"},
		},

		"cognitive_account_id": {
			Description:  "The ID of the Cognitive Account this Bot Channel should be associated with.",
			Type:         schema.TypeString,
			Optional:     true,
			ValidateFunc: cognitiveservicesaccounts.ValidateAccountID,
			AtLeastOneOf: []string{"cognitive_service_access_key", "cognitive_account_id"},
		},

		"custom_speech_model_id": {
			Description:  "The custom speech model ID for the Direct Line Speech Channel.",
			Type:         schema.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"custom_voice_deployment_id": {
			Description:  "The custom voice deployment ID for the Direct Line Speech Channel.",
			Type:         schema.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
	}
}
