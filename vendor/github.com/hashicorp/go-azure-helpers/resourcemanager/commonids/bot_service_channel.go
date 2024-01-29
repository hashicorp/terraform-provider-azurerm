// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package commonids

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.ResourceId = &BotServiceChannelId{}

// BotServiceChannelId is a struct representing the Resource ID for a Bot Service Channel
type BotServiceChannelId struct {
	SubscriptionId    string
	ResourceGroupName string
	BotServiceName    string
	ChannelType       BotServiceChannelType
}

// NewBotServiceChannelID returns a new BotServiceChannelId struct
func NewBotServiceChannelID(subscriptionId string, resourceGroupName string, botServiceName string, channelType BotServiceChannelType) BotServiceChannelId {
	return BotServiceChannelId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		BotServiceName:    botServiceName,
		ChannelType:       channelType,
	}
}

// ParseBotServiceChannelID parses 'input' into a BotServiceChannelId
func ParseBotServiceChannelID(input string) (*BotServiceChannelId, error) {
	parser := resourceids.NewParserFromResourceIdType(&BotServiceChannelId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := BotServiceChannelId{}
	if err := id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseBotServiceChannelIDInsensitively parses 'input' case-insensitively into a BotServiceChannelId
// note: this method should only be used for API response data and not user input
func ParseBotServiceChannelIDInsensitively(input string) (*BotServiceChannelId, error) {
	parser := resourceids.NewParserFromResourceIdType(&BotServiceChannelId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := BotServiceChannelId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *BotServiceChannelId) FromParseResult(input resourceids.ParseResult) error {

	var ok bool
	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.BotServiceName, ok = input.Parsed["botServiceName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "botServiceName", input)
	}

	if v, ok := input.Parsed["channelType"]; true {
		if !ok {
			return resourceids.NewSegmentNotSpecifiedError(id, "channelType", input)
		}

		channelType, err := parseBotServiceChannelType(v)
		if err != nil {
			return fmt.Errorf("parsing %q: %+v", v, err)
		}
		id.ChannelType = *channelType
	}

	return nil
}

// ValidateBotServiceChannelID checks that 'input' can be parsed as a Bot Service Channel ID
func ValidateBotServiceChannelID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseBotServiceChannelID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Bot Service Channel ID
func (id BotServiceChannelId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.BotService/botServices/%s/channels/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.BotServiceName, id.ChannelType)
}

// Segments returns a slice of Resource ID Segments which comprise this Bot Service Channel ID
func (id BotServiceChannelId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftBotService", "Microsoft.BotService", "Microsoft.BotService"),
		resourceids.StaticSegment("staticBotServices", "botServices", "botServices"),
		resourceids.UserSpecifiedSegment("botServiceName", "botServiceValue"),
		resourceids.StaticSegment("staticChannels", "channels", "channels"),
		resourceids.ConstantSegment("channelType", PossibleValuesForBotServiceChannelType(), string(AcsChatBotServiceChannelType)),
	}
}

// String returns a human-readable description of this Bot Service Channel ID
func (id BotServiceChannelId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Bot Service Name: %q", id.BotServiceName),
		fmt.Sprintf("Channel Type: %q", id.ChannelType),
	}
	return fmt.Sprintf("Bot Service Channel (%s)", strings.Join(components, "\n"))
}

type BotServiceChannelType = string

const (
	AcsChatBotServiceChannelType          BotServiceChannelType = "AcsChatChannel"
	AlexaBotServiceChannelType            BotServiceChannelType = "AlexaChannel"
	DirectLineBotServiceChannelType       BotServiceChannelType = "DirectLineChannel"
	DirectLineSpeechBotServiceChannelType BotServiceChannelType = "DirectLineSpeechChannel"
	EmailBotServiceChannelType            BotServiceChannelType = "EmailChannel"
	FacebookBotServiceChannelType         BotServiceChannelType = "FacebookChannel"
	KikChannelBotServiceChannelType       BotServiceChannelType = "KikChannel"
	LineBotServiceChannelType             BotServiceChannelType = "LineChannel"
	M365ExtensionsBotServiceChannelType   BotServiceChannelType = "M365Extensions"
	MsTeamsBotServiceChannelType          BotServiceChannelType = "MsTeamsChannel"
	OmniChannelBotServiceChannelType      BotServiceChannelType = "OmniChannel"
	OutlookBotServiceChannelType          BotServiceChannelType = "OutlookChannel"
	SearchAssistantBotServiceChannelType  BotServiceChannelType = "SearchAssistant"
	SkypeBotServiceChannelType            BotServiceChannelType = "SkypeChannel"
	SlackBotServiceChannelType            BotServiceChannelType = "SlackChannel"
	SmsBotServiceChannelType              BotServiceChannelType = "SmsChannel"
	TelegramBotServiceChannelType         BotServiceChannelType = "TelegramChannel"
	TelephonyBotServiceChannelType        BotServiceChannelType = "TelephonyChannel"
	WebChatBotServiceChannelType          BotServiceChannelType = "WebChatChannel"
)

func parseBotServiceChannelType(input string) (*BotServiceChannelType, error) {
	vals := map[string]BotServiceChannelType{
		"AcsChatChannel":          AcsChatBotServiceChannelType,
		"AlexaChannel":            AlexaBotServiceChannelType,
		"DirectLineChannel":       DirectLineBotServiceChannelType,
		"DirectLineSpeechChannel": DirectLineSpeechBotServiceChannelType,
		"EmailChannel":            EmailBotServiceChannelType,
		"KikChannel":              KikChannelBotServiceChannelType,
		"FacebookChannel":         FacebookBotServiceChannelType,
		"LineChannel":             LineBotServiceChannelType,
		"M365Extensions":          M365ExtensionsBotServiceChannelType,
		"MsTeamsChannel":          MsTeamsBotServiceChannelType,
		"Omnichannel":             OmniChannelBotServiceChannelType,
		"OutlookChannel":          OutlookBotServiceChannelType,
		"SearchAssistant":         SearchAssistantBotServiceChannelType,
		"SkypeChannel":            SkypeBotServiceChannelType,
		"SlackChannel":            SlackBotServiceChannelType,
		"SmsChannel":              SmsBotServiceChannelType,
		"TelegramChannel":         TelegramBotServiceChannelType,
		"TelephonyChannel":        TelephonyBotServiceChannelType,
		"WebChatChannel":          WebChatBotServiceChannelType,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := BotServiceChannelType(input)
	return &out, nil
}

func PossibleValuesForBotServiceChannelType() []string {
	return []string{
		string(AcsChatBotServiceChannelType),
		string(AlexaBotServiceChannelType),
		string(DirectLineBotServiceChannelType),
		string(DirectLineSpeechBotServiceChannelType),
		string(EmailBotServiceChannelType),
		string(FacebookBotServiceChannelType),
		string(KikChannelBotServiceChannelType),
		string(LineBotServiceChannelType),
		string(M365ExtensionsBotServiceChannelType),
		string(MsTeamsBotServiceChannelType),
		string(OmniChannelBotServiceChannelType),
		string(OutlookBotServiceChannelType),
		string(SearchAssistantBotServiceChannelType),
		string(SkypeBotServiceChannelType),
		string(SlackBotServiceChannelType),
		string(SmsBotServiceChannelType),
		string(TelegramBotServiceChannelType),
		string(TelephonyBotServiceChannelType),
		string(WebChatBotServiceChannelType),
	}
}
