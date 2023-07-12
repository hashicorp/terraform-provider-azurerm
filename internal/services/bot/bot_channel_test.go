// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package bot_test

import (
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
)

func TestAccBotChannelsRegistration(t *testing.T) {
	// NOTE: this is a combined test rather than separate split out tests due to
	// Azure only being able provision against one app id at a time
	acceptance.RunTestsInSequence(t, map[string]map[string]func(t *testing.T){
		"basic": {
			"basic":                    testAccBotChannelsRegistration_basic,
			"update":                   testAccBotChannelsRegistration_update,
			"complete":                 testAccBotChannelsRegistration_complete,
			"streamingEndpointEnabled": testAccBotChannelsRegistration_streamingEndpointEnabled,
		},
		"bot": {
			"basic":                    testAccBotServiceAzureBot_basic,
			"completeUpdate":           testAccBotServiceAzureBot_completeUpdate,
			"msaAppType":               testAccBotServiceAzureBot_msaAppType,
			"requiresImport":           testAccBotServiceAzureBot_requiresImport,
			"streamingEndpointEnabled": testAccBotServiceAzureBot_streamingEndpointEnabled,
		},
		"connection": {
			"basic":    testAccBotConnection_basic,
			"complete": testAccBotConnection_complete,
		},
		"channel": {
			"alexaBasic":                     testAccBotChannelAlexa_basic,
			"alexaUpdate":                    testAccBotChannelAlexa_update,
			"alexaRequiresImport":            testAccBotChannelAlexa_requiresImport,
			"emailBasic":                     testAccBotChannelEmail_basic,
			"emailUpdate":                    testAccBotChannelEmail_update,
			"slackBasic":                     testAccBotChannelSlack_basic,
			"slackComplete":                  testAccBotChannelSlack_complete,
			"slackUpdate":                    testAccBotChannelSlack_update,
			"smsBasic":                       testAccBotChannelSMS_basic,
			"smsRequiresImport":              testAccBotChannelSMS_requiresImport,
			"msteamsBasic":                   testAccBotChannelMsTeams_basic,
			"msteamsUpdate":                  testAccBotChannelMsTeams_update,
			"directlineBasic":                testAccBotChannelDirectline_basic,
			"directlineComplete":             testAccBotChannelDirectline_complete,
			"directlineUpdate":               testAccBotChannelDirectline_update,
			"directLineSpeechBasic":          testAccBotChannelDirectLineSpeech_basic,
			"directLineSpeechComplete":       testAccBotChannelDirectLineSpeech_complete,
			"directLineSpeechUpdate":         testAccBotChannelDirectLineSpeech_update,
			"directLineSpeechRequiresImport": testAccBotChannelDirectLineSpeech_requiresImport,
			"facebookBasic":                  testAccBotChannelFacebook_basic,
			"facebookUpdate":                 testAccBotChannelFacebook_update,
			"facebookRequiresImport":         testAccBotChannelFacebook_requiresImport,
			"lineBasic":                      testAccBotChannelLine_basic,
			"lineComplete":                   testAccBotChannelLine_complete,
			"lineUpdate":                     testAccBotChannelLine_update,
			"lineRequiresImport":             testAccBotChannelLine_requiresImport,
			"webchatBasic":                   testAccBotChannelWebChat_basic,
			"webchatComplete":                testAccBotChannelWebChat_complete,
			"webchatUpdate":                  testAccBotChannelWebChat_update,
			"webchatRequiresImport":          testAccBotChannelWebChat_requiresImport,
		},
		"web_app": {
			"basic":    testAccBotWebApp_basic,
			"update":   testAccBotWebApp_update,
			"complete": testAccBotWebApp_complete,
		},
	})
}
