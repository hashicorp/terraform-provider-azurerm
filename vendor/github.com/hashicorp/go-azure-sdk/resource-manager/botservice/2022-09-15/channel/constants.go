package channel

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type BotServiceChannelType string

const (
	BotServiceChannelTypeAcsChatChannel          BotServiceChannelType = "AcsChatChannel"
	BotServiceChannelTypeAlexaChannel            BotServiceChannelType = "AlexaChannel"
	BotServiceChannelTypeDirectLineChannel       BotServiceChannelType = "DirectLineChannel"
	BotServiceChannelTypeDirectLineSpeechChannel BotServiceChannelType = "DirectLineSpeechChannel"
	BotServiceChannelTypeEmailChannel            BotServiceChannelType = "EmailChannel"
	BotServiceChannelTypeFacebookChannel         BotServiceChannelType = "FacebookChannel"
	BotServiceChannelTypeKikChannel              BotServiceChannelType = "KikChannel"
	BotServiceChannelTypeLineChannel             BotServiceChannelType = "LineChannel"
	BotServiceChannelTypeM365Extensions          BotServiceChannelType = "M365Extensions"
	BotServiceChannelTypeMsTeamsChannel          BotServiceChannelType = "MsTeamsChannel"
	BotServiceChannelTypeOmnichannel             BotServiceChannelType = "Omnichannel"
	BotServiceChannelTypeOutlookChannel          BotServiceChannelType = "OutlookChannel"
	BotServiceChannelTypeSearchAssistant         BotServiceChannelType = "SearchAssistant"
	BotServiceChannelTypeSkypeChannel            BotServiceChannelType = "SkypeChannel"
	BotServiceChannelTypeSlackChannel            BotServiceChannelType = "SlackChannel"
	BotServiceChannelTypeSmsChannel              BotServiceChannelType = "SmsChannel"
	BotServiceChannelTypeTelegramChannel         BotServiceChannelType = "TelegramChannel"
	BotServiceChannelTypeTelephonyChannel        BotServiceChannelType = "TelephonyChannel"
	BotServiceChannelTypeWebChatChannel          BotServiceChannelType = "WebChatChannel"
)

func PossibleValuesForBotServiceChannelType() []string {
	return []string{
		string(BotServiceChannelTypeAcsChatChannel),
		string(BotServiceChannelTypeAlexaChannel),
		string(BotServiceChannelTypeDirectLineChannel),
		string(BotServiceChannelTypeDirectLineSpeechChannel),
		string(BotServiceChannelTypeEmailChannel),
		string(BotServiceChannelTypeFacebookChannel),
		string(BotServiceChannelTypeKikChannel),
		string(BotServiceChannelTypeLineChannel),
		string(BotServiceChannelTypeM365Extensions),
		string(BotServiceChannelTypeMsTeamsChannel),
		string(BotServiceChannelTypeOmnichannel),
		string(BotServiceChannelTypeOutlookChannel),
		string(BotServiceChannelTypeSearchAssistant),
		string(BotServiceChannelTypeSkypeChannel),
		string(BotServiceChannelTypeSlackChannel),
		string(BotServiceChannelTypeSmsChannel),
		string(BotServiceChannelTypeTelegramChannel),
		string(BotServiceChannelTypeTelephonyChannel),
		string(BotServiceChannelTypeWebChatChannel),
	}
}

func (s *BotServiceChannelType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseBotServiceChannelType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseBotServiceChannelType(input string) (*BotServiceChannelType, error) {
	vals := map[string]BotServiceChannelType{
		"acschatchannel":          BotServiceChannelTypeAcsChatChannel,
		"alexachannel":            BotServiceChannelTypeAlexaChannel,
		"directlinechannel":       BotServiceChannelTypeDirectLineChannel,
		"directlinespeechchannel": BotServiceChannelTypeDirectLineSpeechChannel,
		"emailchannel":            BotServiceChannelTypeEmailChannel,
		"facebookchannel":         BotServiceChannelTypeFacebookChannel,
		"kikchannel":              BotServiceChannelTypeKikChannel,
		"linechannel":             BotServiceChannelTypeLineChannel,
		"m365extensions":          BotServiceChannelTypeM365Extensions,
		"msteamschannel":          BotServiceChannelTypeMsTeamsChannel,
		"omnichannel":             BotServiceChannelTypeOmnichannel,
		"outlookchannel":          BotServiceChannelTypeOutlookChannel,
		"searchassistant":         BotServiceChannelTypeSearchAssistant,
		"skypechannel":            BotServiceChannelTypeSkypeChannel,
		"slackchannel":            BotServiceChannelTypeSlackChannel,
		"smschannel":              BotServiceChannelTypeSmsChannel,
		"telegramchannel":         BotServiceChannelTypeTelegramChannel,
		"telephonychannel":        BotServiceChannelTypeTelephonyChannel,
		"webchatchannel":          BotServiceChannelTypeWebChatChannel,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := BotServiceChannelType(input)
	return &out, nil
}

type EmailChannelAuthMethod int64

const (
	EmailChannelAuthMethodOne  EmailChannelAuthMethod = 1
	EmailChannelAuthMethodZero EmailChannelAuthMethod = 0
)

func PossibleValuesForEmailChannelAuthMethod() []int64 {
	return []int64{
		int64(EmailChannelAuthMethodOne),
		int64(EmailChannelAuthMethodZero),
	}
}

type Key string

const (
	KeyKeyOne Key = "key1"
	KeyKeyTwo Key = "key2"
)

func PossibleValuesForKey() []string {
	return []string{
		string(KeyKeyOne),
		string(KeyKeyTwo),
	}
}

func (s *Key) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseKey(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseKey(input string) (*Key, error) {
	vals := map[string]Key{
		"key1": KeyKeyOne,
		"key2": KeyKeyTwo,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := Key(input)
	return &out, nil
}

type Kind string

const (
	KindAzurebot Kind = "azurebot"
	KindBot      Kind = "bot"
	KindDesigner Kind = "designer"
	KindFunction Kind = "function"
	KindSdk      Kind = "sdk"
)

func PossibleValuesForKind() []string {
	return []string{
		string(KindAzurebot),
		string(KindBot),
		string(KindDesigner),
		string(KindFunction),
		string(KindSdk),
	}
}

func (s *Kind) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseKind(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseKind(input string) (*Kind, error) {
	vals := map[string]Kind{
		"azurebot": KindAzurebot,
		"bot":      KindBot,
		"designer": KindDesigner,
		"function": KindFunction,
		"sdk":      KindSdk,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := Kind(input)
	return &out, nil
}

type SkuName string

const (
	SkuNameFZero SkuName = "F0"
	SkuNameSOne  SkuName = "S1"
)

func PossibleValuesForSkuName() []string {
	return []string{
		string(SkuNameFZero),
		string(SkuNameSOne),
	}
}

func (s *SkuName) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseSkuName(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseSkuName(input string) (*SkuName, error) {
	vals := map[string]SkuName{
		"f0": SkuNameFZero,
		"s1": SkuNameSOne,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := SkuName(input)
	return &out, nil
}

type SkuTier string

const (
	SkuTierFree     SkuTier = "Free"
	SkuTierStandard SkuTier = "Standard"
)

func PossibleValuesForSkuTier() []string {
	return []string{
		string(SkuTierFree),
		string(SkuTierStandard),
	}
}

func (s *SkuTier) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseSkuTier(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseSkuTier(input string) (*SkuTier, error) {
	vals := map[string]SkuTier{
		"free":     SkuTierFree,
		"standard": SkuTierStandard,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := SkuTier(input)
	return &out, nil
}
