package skillsets

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AzureOpenAIModelName string

const (
	AzureOpenAIModelNameTextNegativeembeddingNegativeThreeNegativelarge    AzureOpenAIModelName = "text-embedding-3-large"
	AzureOpenAIModelNameTextNegativeembeddingNegativeThreeNegativesmall    AzureOpenAIModelName = "text-embedding-3-small"
	AzureOpenAIModelNameTextNegativeembeddingNegativeadaNegativeHundredTwo AzureOpenAIModelName = "text-embedding-ada-002"
)

func PossibleValuesForAzureOpenAIModelName() []string {
	return []string{
		string(AzureOpenAIModelNameTextNegativeembeddingNegativeThreeNegativelarge),
		string(AzureOpenAIModelNameTextNegativeembeddingNegativeThreeNegativesmall),
		string(AzureOpenAIModelNameTextNegativeembeddingNegativeadaNegativeHundredTwo),
	}
}

func (s *AzureOpenAIModelName) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseAzureOpenAIModelName(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseAzureOpenAIModelName(input string) (*AzureOpenAIModelName, error) {
	vals := map[string]AzureOpenAIModelName{
		"text-embedding-3-large": AzureOpenAIModelNameTextNegativeembeddingNegativeThreeNegativelarge,
		"text-embedding-3-small": AzureOpenAIModelNameTextNegativeembeddingNegativeThreeNegativesmall,
		"text-embedding-ada-002": AzureOpenAIModelNameTextNegativeembeddingNegativeadaNegativeHundredTwo,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := AzureOpenAIModelName(input)
	return &out, nil
}

type CustomEntityLookupSkillLanguage string

const (
	CustomEntityLookupSkillLanguageDa CustomEntityLookupSkillLanguage = "da"
	CustomEntityLookupSkillLanguageDe CustomEntityLookupSkillLanguage = "de"
	CustomEntityLookupSkillLanguageEn CustomEntityLookupSkillLanguage = "en"
	CustomEntityLookupSkillLanguageEs CustomEntityLookupSkillLanguage = "es"
	CustomEntityLookupSkillLanguageFi CustomEntityLookupSkillLanguage = "fi"
	CustomEntityLookupSkillLanguageFr CustomEntityLookupSkillLanguage = "fr"
	CustomEntityLookupSkillLanguageIt CustomEntityLookupSkillLanguage = "it"
	CustomEntityLookupSkillLanguageKo CustomEntityLookupSkillLanguage = "ko"
	CustomEntityLookupSkillLanguagePt CustomEntityLookupSkillLanguage = "pt"
)

func PossibleValuesForCustomEntityLookupSkillLanguage() []string {
	return []string{
		string(CustomEntityLookupSkillLanguageDa),
		string(CustomEntityLookupSkillLanguageDe),
		string(CustomEntityLookupSkillLanguageEn),
		string(CustomEntityLookupSkillLanguageEs),
		string(CustomEntityLookupSkillLanguageFi),
		string(CustomEntityLookupSkillLanguageFr),
		string(CustomEntityLookupSkillLanguageIt),
		string(CustomEntityLookupSkillLanguageKo),
		string(CustomEntityLookupSkillLanguagePt),
	}
}

func (s *CustomEntityLookupSkillLanguage) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseCustomEntityLookupSkillLanguage(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseCustomEntityLookupSkillLanguage(input string) (*CustomEntityLookupSkillLanguage, error) {
	vals := map[string]CustomEntityLookupSkillLanguage{
		"da": CustomEntityLookupSkillLanguageDa,
		"de": CustomEntityLookupSkillLanguageDe,
		"en": CustomEntityLookupSkillLanguageEn,
		"es": CustomEntityLookupSkillLanguageEs,
		"fi": CustomEntityLookupSkillLanguageFi,
		"fr": CustomEntityLookupSkillLanguageFr,
		"it": CustomEntityLookupSkillLanguageIt,
		"ko": CustomEntityLookupSkillLanguageKo,
		"pt": CustomEntityLookupSkillLanguagePt,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := CustomEntityLookupSkillLanguage(input)
	return &out, nil
}

type DocumentIntelligenceLayoutSkillChunkingUnit string

const (
	DocumentIntelligenceLayoutSkillChunkingUnitCharacters DocumentIntelligenceLayoutSkillChunkingUnit = "characters"
)

func PossibleValuesForDocumentIntelligenceLayoutSkillChunkingUnit() []string {
	return []string{
		string(DocumentIntelligenceLayoutSkillChunkingUnitCharacters),
	}
}

func (s *DocumentIntelligenceLayoutSkillChunkingUnit) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseDocumentIntelligenceLayoutSkillChunkingUnit(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseDocumentIntelligenceLayoutSkillChunkingUnit(input string) (*DocumentIntelligenceLayoutSkillChunkingUnit, error) {
	vals := map[string]DocumentIntelligenceLayoutSkillChunkingUnit{
		"characters": DocumentIntelligenceLayoutSkillChunkingUnitCharacters,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := DocumentIntelligenceLayoutSkillChunkingUnit(input)
	return &out, nil
}

type DocumentIntelligenceLayoutSkillExtractionOptions string

const (
	DocumentIntelligenceLayoutSkillExtractionOptionsImages           DocumentIntelligenceLayoutSkillExtractionOptions = "images"
	DocumentIntelligenceLayoutSkillExtractionOptionsLocationMetadata DocumentIntelligenceLayoutSkillExtractionOptions = "locationMetadata"
)

func PossibleValuesForDocumentIntelligenceLayoutSkillExtractionOptions() []string {
	return []string{
		string(DocumentIntelligenceLayoutSkillExtractionOptionsImages),
		string(DocumentIntelligenceLayoutSkillExtractionOptionsLocationMetadata),
	}
}

func (s *DocumentIntelligenceLayoutSkillExtractionOptions) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseDocumentIntelligenceLayoutSkillExtractionOptions(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseDocumentIntelligenceLayoutSkillExtractionOptions(input string) (*DocumentIntelligenceLayoutSkillExtractionOptions, error) {
	vals := map[string]DocumentIntelligenceLayoutSkillExtractionOptions{
		"images":           DocumentIntelligenceLayoutSkillExtractionOptionsImages,
		"locationmetadata": DocumentIntelligenceLayoutSkillExtractionOptionsLocationMetadata,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := DocumentIntelligenceLayoutSkillExtractionOptions(input)
	return &out, nil
}

type DocumentIntelligenceLayoutSkillMarkdownHeaderDepth string

const (
	DocumentIntelligenceLayoutSkillMarkdownHeaderDepthHFive  DocumentIntelligenceLayoutSkillMarkdownHeaderDepth = "h5"
	DocumentIntelligenceLayoutSkillMarkdownHeaderDepthHFour  DocumentIntelligenceLayoutSkillMarkdownHeaderDepth = "h4"
	DocumentIntelligenceLayoutSkillMarkdownHeaderDepthHOne   DocumentIntelligenceLayoutSkillMarkdownHeaderDepth = "h1"
	DocumentIntelligenceLayoutSkillMarkdownHeaderDepthHSix   DocumentIntelligenceLayoutSkillMarkdownHeaderDepth = "h6"
	DocumentIntelligenceLayoutSkillMarkdownHeaderDepthHThree DocumentIntelligenceLayoutSkillMarkdownHeaderDepth = "h3"
	DocumentIntelligenceLayoutSkillMarkdownHeaderDepthHTwo   DocumentIntelligenceLayoutSkillMarkdownHeaderDepth = "h2"
)

func PossibleValuesForDocumentIntelligenceLayoutSkillMarkdownHeaderDepth() []string {
	return []string{
		string(DocumentIntelligenceLayoutSkillMarkdownHeaderDepthHFive),
		string(DocumentIntelligenceLayoutSkillMarkdownHeaderDepthHFour),
		string(DocumentIntelligenceLayoutSkillMarkdownHeaderDepthHOne),
		string(DocumentIntelligenceLayoutSkillMarkdownHeaderDepthHSix),
		string(DocumentIntelligenceLayoutSkillMarkdownHeaderDepthHThree),
		string(DocumentIntelligenceLayoutSkillMarkdownHeaderDepthHTwo),
	}
}

func (s *DocumentIntelligenceLayoutSkillMarkdownHeaderDepth) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseDocumentIntelligenceLayoutSkillMarkdownHeaderDepth(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseDocumentIntelligenceLayoutSkillMarkdownHeaderDepth(input string) (*DocumentIntelligenceLayoutSkillMarkdownHeaderDepth, error) {
	vals := map[string]DocumentIntelligenceLayoutSkillMarkdownHeaderDepth{
		"h5": DocumentIntelligenceLayoutSkillMarkdownHeaderDepthHFive,
		"h4": DocumentIntelligenceLayoutSkillMarkdownHeaderDepthHFour,
		"h1": DocumentIntelligenceLayoutSkillMarkdownHeaderDepthHOne,
		"h6": DocumentIntelligenceLayoutSkillMarkdownHeaderDepthHSix,
		"h3": DocumentIntelligenceLayoutSkillMarkdownHeaderDepthHThree,
		"h2": DocumentIntelligenceLayoutSkillMarkdownHeaderDepthHTwo,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := DocumentIntelligenceLayoutSkillMarkdownHeaderDepth(input)
	return &out, nil
}

type DocumentIntelligenceLayoutSkillOutputFormat string

const (
	DocumentIntelligenceLayoutSkillOutputFormatMarkdown DocumentIntelligenceLayoutSkillOutputFormat = "markdown"
	DocumentIntelligenceLayoutSkillOutputFormatText     DocumentIntelligenceLayoutSkillOutputFormat = "text"
)

func PossibleValuesForDocumentIntelligenceLayoutSkillOutputFormat() []string {
	return []string{
		string(DocumentIntelligenceLayoutSkillOutputFormatMarkdown),
		string(DocumentIntelligenceLayoutSkillOutputFormatText),
	}
}

func (s *DocumentIntelligenceLayoutSkillOutputFormat) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseDocumentIntelligenceLayoutSkillOutputFormat(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseDocumentIntelligenceLayoutSkillOutputFormat(input string) (*DocumentIntelligenceLayoutSkillOutputFormat, error) {
	vals := map[string]DocumentIntelligenceLayoutSkillOutputFormat{
		"markdown": DocumentIntelligenceLayoutSkillOutputFormatMarkdown,
		"text":     DocumentIntelligenceLayoutSkillOutputFormatText,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := DocumentIntelligenceLayoutSkillOutputFormat(input)
	return &out, nil
}

type DocumentIntelligenceLayoutSkillOutputMode string

const (
	DocumentIntelligenceLayoutSkillOutputModeOneToMany DocumentIntelligenceLayoutSkillOutputMode = "oneToMany"
)

func PossibleValuesForDocumentIntelligenceLayoutSkillOutputMode() []string {
	return []string{
		string(DocumentIntelligenceLayoutSkillOutputModeOneToMany),
	}
}

func (s *DocumentIntelligenceLayoutSkillOutputMode) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseDocumentIntelligenceLayoutSkillOutputMode(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseDocumentIntelligenceLayoutSkillOutputMode(input string) (*DocumentIntelligenceLayoutSkillOutputMode, error) {
	vals := map[string]DocumentIntelligenceLayoutSkillOutputMode{
		"onetomany": DocumentIntelligenceLayoutSkillOutputModeOneToMany,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := DocumentIntelligenceLayoutSkillOutputMode(input)
	return &out, nil
}

type EntityCategory string

const (
	EntityCategoryDatetime     EntityCategory = "datetime"
	EntityCategoryEmail        EntityCategory = "email"
	EntityCategoryLocation     EntityCategory = "location"
	EntityCategoryOrganization EntityCategory = "organization"
	EntityCategoryPerson       EntityCategory = "person"
	EntityCategoryQuantity     EntityCategory = "quantity"
	EntityCategoryUrl          EntityCategory = "url"
)

func PossibleValuesForEntityCategory() []string {
	return []string{
		string(EntityCategoryDatetime),
		string(EntityCategoryEmail),
		string(EntityCategoryLocation),
		string(EntityCategoryOrganization),
		string(EntityCategoryPerson),
		string(EntityCategoryQuantity),
		string(EntityCategoryUrl),
	}
}

func (s *EntityCategory) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseEntityCategory(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseEntityCategory(input string) (*EntityCategory, error) {
	vals := map[string]EntityCategory{
		"datetime":     EntityCategoryDatetime,
		"email":        EntityCategoryEmail,
		"location":     EntityCategoryLocation,
		"organization": EntityCategoryOrganization,
		"person":       EntityCategoryPerson,
		"quantity":     EntityCategoryQuantity,
		"url":          EntityCategoryUrl,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := EntityCategory(input)
	return &out, nil
}

type EntityRecognitionSkillLanguage string

const (
	EntityRecognitionSkillLanguageAr             EntityRecognitionSkillLanguage = "ar"
	EntityRecognitionSkillLanguageCs             EntityRecognitionSkillLanguage = "cs"
	EntityRecognitionSkillLanguageDa             EntityRecognitionSkillLanguage = "da"
	EntityRecognitionSkillLanguageDe             EntityRecognitionSkillLanguage = "de"
	EntityRecognitionSkillLanguageEl             EntityRecognitionSkillLanguage = "el"
	EntityRecognitionSkillLanguageEn             EntityRecognitionSkillLanguage = "en"
	EntityRecognitionSkillLanguageEs             EntityRecognitionSkillLanguage = "es"
	EntityRecognitionSkillLanguageFi             EntityRecognitionSkillLanguage = "fi"
	EntityRecognitionSkillLanguageFr             EntityRecognitionSkillLanguage = "fr"
	EntityRecognitionSkillLanguageHu             EntityRecognitionSkillLanguage = "hu"
	EntityRecognitionSkillLanguageIt             EntityRecognitionSkillLanguage = "it"
	EntityRecognitionSkillLanguageJa             EntityRecognitionSkillLanguage = "ja"
	EntityRecognitionSkillLanguageKo             EntityRecognitionSkillLanguage = "ko"
	EntityRecognitionSkillLanguageNl             EntityRecognitionSkillLanguage = "nl"
	EntityRecognitionSkillLanguageNo             EntityRecognitionSkillLanguage = "no"
	EntityRecognitionSkillLanguagePl             EntityRecognitionSkillLanguage = "pl"
	EntityRecognitionSkillLanguagePtNegativeBR   EntityRecognitionSkillLanguage = "pt-BR"
	EntityRecognitionSkillLanguagePtNegativePT   EntityRecognitionSkillLanguage = "pt-PT"
	EntityRecognitionSkillLanguageRu             EntityRecognitionSkillLanguage = "ru"
	EntityRecognitionSkillLanguageSv             EntityRecognitionSkillLanguage = "sv"
	EntityRecognitionSkillLanguageTr             EntityRecognitionSkillLanguage = "tr"
	EntityRecognitionSkillLanguageZhNegativeHans EntityRecognitionSkillLanguage = "zh-Hans"
	EntityRecognitionSkillLanguageZhNegativeHant EntityRecognitionSkillLanguage = "zh-Hant"
)

func PossibleValuesForEntityRecognitionSkillLanguage() []string {
	return []string{
		string(EntityRecognitionSkillLanguageAr),
		string(EntityRecognitionSkillLanguageCs),
		string(EntityRecognitionSkillLanguageDa),
		string(EntityRecognitionSkillLanguageDe),
		string(EntityRecognitionSkillLanguageEl),
		string(EntityRecognitionSkillLanguageEn),
		string(EntityRecognitionSkillLanguageEs),
		string(EntityRecognitionSkillLanguageFi),
		string(EntityRecognitionSkillLanguageFr),
		string(EntityRecognitionSkillLanguageHu),
		string(EntityRecognitionSkillLanguageIt),
		string(EntityRecognitionSkillLanguageJa),
		string(EntityRecognitionSkillLanguageKo),
		string(EntityRecognitionSkillLanguageNl),
		string(EntityRecognitionSkillLanguageNo),
		string(EntityRecognitionSkillLanguagePl),
		string(EntityRecognitionSkillLanguagePtNegativeBR),
		string(EntityRecognitionSkillLanguagePtNegativePT),
		string(EntityRecognitionSkillLanguageRu),
		string(EntityRecognitionSkillLanguageSv),
		string(EntityRecognitionSkillLanguageTr),
		string(EntityRecognitionSkillLanguageZhNegativeHans),
		string(EntityRecognitionSkillLanguageZhNegativeHant),
	}
}

func (s *EntityRecognitionSkillLanguage) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseEntityRecognitionSkillLanguage(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseEntityRecognitionSkillLanguage(input string) (*EntityRecognitionSkillLanguage, error) {
	vals := map[string]EntityRecognitionSkillLanguage{
		"ar":      EntityRecognitionSkillLanguageAr,
		"cs":      EntityRecognitionSkillLanguageCs,
		"da":      EntityRecognitionSkillLanguageDa,
		"de":      EntityRecognitionSkillLanguageDe,
		"el":      EntityRecognitionSkillLanguageEl,
		"en":      EntityRecognitionSkillLanguageEn,
		"es":      EntityRecognitionSkillLanguageEs,
		"fi":      EntityRecognitionSkillLanguageFi,
		"fr":      EntityRecognitionSkillLanguageFr,
		"hu":      EntityRecognitionSkillLanguageHu,
		"it":      EntityRecognitionSkillLanguageIt,
		"ja":      EntityRecognitionSkillLanguageJa,
		"ko":      EntityRecognitionSkillLanguageKo,
		"nl":      EntityRecognitionSkillLanguageNl,
		"no":      EntityRecognitionSkillLanguageNo,
		"pl":      EntityRecognitionSkillLanguagePl,
		"pt-br":   EntityRecognitionSkillLanguagePtNegativeBR,
		"pt-pt":   EntityRecognitionSkillLanguagePtNegativePT,
		"ru":      EntityRecognitionSkillLanguageRu,
		"sv":      EntityRecognitionSkillLanguageSv,
		"tr":      EntityRecognitionSkillLanguageTr,
		"zh-hans": EntityRecognitionSkillLanguageZhNegativeHans,
		"zh-hant": EntityRecognitionSkillLanguageZhNegativeHant,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := EntityRecognitionSkillLanguage(input)
	return &out, nil
}

type ImageAnalysisSkillLanguage string

const (
	ImageAnalysisSkillLanguageAr             ImageAnalysisSkillLanguage = "ar"
	ImageAnalysisSkillLanguageAz             ImageAnalysisSkillLanguage = "az"
	ImageAnalysisSkillLanguageBg             ImageAnalysisSkillLanguage = "bg"
	ImageAnalysisSkillLanguageBs             ImageAnalysisSkillLanguage = "bs"
	ImageAnalysisSkillLanguageCa             ImageAnalysisSkillLanguage = "ca"
	ImageAnalysisSkillLanguageCs             ImageAnalysisSkillLanguage = "cs"
	ImageAnalysisSkillLanguageCy             ImageAnalysisSkillLanguage = "cy"
	ImageAnalysisSkillLanguageDa             ImageAnalysisSkillLanguage = "da"
	ImageAnalysisSkillLanguageDe             ImageAnalysisSkillLanguage = "de"
	ImageAnalysisSkillLanguageEl             ImageAnalysisSkillLanguage = "el"
	ImageAnalysisSkillLanguageEn             ImageAnalysisSkillLanguage = "en"
	ImageAnalysisSkillLanguageEs             ImageAnalysisSkillLanguage = "es"
	ImageAnalysisSkillLanguageEt             ImageAnalysisSkillLanguage = "et"
	ImageAnalysisSkillLanguageEu             ImageAnalysisSkillLanguage = "eu"
	ImageAnalysisSkillLanguageFi             ImageAnalysisSkillLanguage = "fi"
	ImageAnalysisSkillLanguageFr             ImageAnalysisSkillLanguage = "fr"
	ImageAnalysisSkillLanguageGa             ImageAnalysisSkillLanguage = "ga"
	ImageAnalysisSkillLanguageGl             ImageAnalysisSkillLanguage = "gl"
	ImageAnalysisSkillLanguageHe             ImageAnalysisSkillLanguage = "he"
	ImageAnalysisSkillLanguageHi             ImageAnalysisSkillLanguage = "hi"
	ImageAnalysisSkillLanguageHr             ImageAnalysisSkillLanguage = "hr"
	ImageAnalysisSkillLanguageHu             ImageAnalysisSkillLanguage = "hu"
	ImageAnalysisSkillLanguageId             ImageAnalysisSkillLanguage = "id"
	ImageAnalysisSkillLanguageIt             ImageAnalysisSkillLanguage = "it"
	ImageAnalysisSkillLanguageJa             ImageAnalysisSkillLanguage = "ja"
	ImageAnalysisSkillLanguageKk             ImageAnalysisSkillLanguage = "kk"
	ImageAnalysisSkillLanguageKo             ImageAnalysisSkillLanguage = "ko"
	ImageAnalysisSkillLanguageLt             ImageAnalysisSkillLanguage = "lt"
	ImageAnalysisSkillLanguageLv             ImageAnalysisSkillLanguage = "lv"
	ImageAnalysisSkillLanguageMk             ImageAnalysisSkillLanguage = "mk"
	ImageAnalysisSkillLanguageMs             ImageAnalysisSkillLanguage = "ms"
	ImageAnalysisSkillLanguageNb             ImageAnalysisSkillLanguage = "nb"
	ImageAnalysisSkillLanguageNl             ImageAnalysisSkillLanguage = "nl"
	ImageAnalysisSkillLanguagePl             ImageAnalysisSkillLanguage = "pl"
	ImageAnalysisSkillLanguagePrs            ImageAnalysisSkillLanguage = "prs"
	ImageAnalysisSkillLanguagePt             ImageAnalysisSkillLanguage = "pt"
	ImageAnalysisSkillLanguagePtNegativeBR   ImageAnalysisSkillLanguage = "pt-BR"
	ImageAnalysisSkillLanguagePtNegativePT   ImageAnalysisSkillLanguage = "pt-PT"
	ImageAnalysisSkillLanguageRo             ImageAnalysisSkillLanguage = "ro"
	ImageAnalysisSkillLanguageRu             ImageAnalysisSkillLanguage = "ru"
	ImageAnalysisSkillLanguageSk             ImageAnalysisSkillLanguage = "sk"
	ImageAnalysisSkillLanguageSl             ImageAnalysisSkillLanguage = "sl"
	ImageAnalysisSkillLanguageSrNegativeCyrl ImageAnalysisSkillLanguage = "sr-Cyrl"
	ImageAnalysisSkillLanguageSrNegativeLatn ImageAnalysisSkillLanguage = "sr-Latn"
	ImageAnalysisSkillLanguageSv             ImageAnalysisSkillLanguage = "sv"
	ImageAnalysisSkillLanguageTh             ImageAnalysisSkillLanguage = "th"
	ImageAnalysisSkillLanguageTr             ImageAnalysisSkillLanguage = "tr"
	ImageAnalysisSkillLanguageUk             ImageAnalysisSkillLanguage = "uk"
	ImageAnalysisSkillLanguageVi             ImageAnalysisSkillLanguage = "vi"
	ImageAnalysisSkillLanguageZh             ImageAnalysisSkillLanguage = "zh"
	ImageAnalysisSkillLanguageZhNegativeHans ImageAnalysisSkillLanguage = "zh-Hans"
	ImageAnalysisSkillLanguageZhNegativeHant ImageAnalysisSkillLanguage = "zh-Hant"
)

func PossibleValuesForImageAnalysisSkillLanguage() []string {
	return []string{
		string(ImageAnalysisSkillLanguageAr),
		string(ImageAnalysisSkillLanguageAz),
		string(ImageAnalysisSkillLanguageBg),
		string(ImageAnalysisSkillLanguageBs),
		string(ImageAnalysisSkillLanguageCa),
		string(ImageAnalysisSkillLanguageCs),
		string(ImageAnalysisSkillLanguageCy),
		string(ImageAnalysisSkillLanguageDa),
		string(ImageAnalysisSkillLanguageDe),
		string(ImageAnalysisSkillLanguageEl),
		string(ImageAnalysisSkillLanguageEn),
		string(ImageAnalysisSkillLanguageEs),
		string(ImageAnalysisSkillLanguageEt),
		string(ImageAnalysisSkillLanguageEu),
		string(ImageAnalysisSkillLanguageFi),
		string(ImageAnalysisSkillLanguageFr),
		string(ImageAnalysisSkillLanguageGa),
		string(ImageAnalysisSkillLanguageGl),
		string(ImageAnalysisSkillLanguageHe),
		string(ImageAnalysisSkillLanguageHi),
		string(ImageAnalysisSkillLanguageHr),
		string(ImageAnalysisSkillLanguageHu),
		string(ImageAnalysisSkillLanguageId),
		string(ImageAnalysisSkillLanguageIt),
		string(ImageAnalysisSkillLanguageJa),
		string(ImageAnalysisSkillLanguageKk),
		string(ImageAnalysisSkillLanguageKo),
		string(ImageAnalysisSkillLanguageLt),
		string(ImageAnalysisSkillLanguageLv),
		string(ImageAnalysisSkillLanguageMk),
		string(ImageAnalysisSkillLanguageMs),
		string(ImageAnalysisSkillLanguageNb),
		string(ImageAnalysisSkillLanguageNl),
		string(ImageAnalysisSkillLanguagePl),
		string(ImageAnalysisSkillLanguagePrs),
		string(ImageAnalysisSkillLanguagePt),
		string(ImageAnalysisSkillLanguagePtNegativeBR),
		string(ImageAnalysisSkillLanguagePtNegativePT),
		string(ImageAnalysisSkillLanguageRo),
		string(ImageAnalysisSkillLanguageRu),
		string(ImageAnalysisSkillLanguageSk),
		string(ImageAnalysisSkillLanguageSl),
		string(ImageAnalysisSkillLanguageSrNegativeCyrl),
		string(ImageAnalysisSkillLanguageSrNegativeLatn),
		string(ImageAnalysisSkillLanguageSv),
		string(ImageAnalysisSkillLanguageTh),
		string(ImageAnalysisSkillLanguageTr),
		string(ImageAnalysisSkillLanguageUk),
		string(ImageAnalysisSkillLanguageVi),
		string(ImageAnalysisSkillLanguageZh),
		string(ImageAnalysisSkillLanguageZhNegativeHans),
		string(ImageAnalysisSkillLanguageZhNegativeHant),
	}
}

func (s *ImageAnalysisSkillLanguage) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseImageAnalysisSkillLanguage(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseImageAnalysisSkillLanguage(input string) (*ImageAnalysisSkillLanguage, error) {
	vals := map[string]ImageAnalysisSkillLanguage{
		"ar":      ImageAnalysisSkillLanguageAr,
		"az":      ImageAnalysisSkillLanguageAz,
		"bg":      ImageAnalysisSkillLanguageBg,
		"bs":      ImageAnalysisSkillLanguageBs,
		"ca":      ImageAnalysisSkillLanguageCa,
		"cs":      ImageAnalysisSkillLanguageCs,
		"cy":      ImageAnalysisSkillLanguageCy,
		"da":      ImageAnalysisSkillLanguageDa,
		"de":      ImageAnalysisSkillLanguageDe,
		"el":      ImageAnalysisSkillLanguageEl,
		"en":      ImageAnalysisSkillLanguageEn,
		"es":      ImageAnalysisSkillLanguageEs,
		"et":      ImageAnalysisSkillLanguageEt,
		"eu":      ImageAnalysisSkillLanguageEu,
		"fi":      ImageAnalysisSkillLanguageFi,
		"fr":      ImageAnalysisSkillLanguageFr,
		"ga":      ImageAnalysisSkillLanguageGa,
		"gl":      ImageAnalysisSkillLanguageGl,
		"he":      ImageAnalysisSkillLanguageHe,
		"hi":      ImageAnalysisSkillLanguageHi,
		"hr":      ImageAnalysisSkillLanguageHr,
		"hu":      ImageAnalysisSkillLanguageHu,
		"id":      ImageAnalysisSkillLanguageId,
		"it":      ImageAnalysisSkillLanguageIt,
		"ja":      ImageAnalysisSkillLanguageJa,
		"kk":      ImageAnalysisSkillLanguageKk,
		"ko":      ImageAnalysisSkillLanguageKo,
		"lt":      ImageAnalysisSkillLanguageLt,
		"lv":      ImageAnalysisSkillLanguageLv,
		"mk":      ImageAnalysisSkillLanguageMk,
		"ms":      ImageAnalysisSkillLanguageMs,
		"nb":      ImageAnalysisSkillLanguageNb,
		"nl":      ImageAnalysisSkillLanguageNl,
		"pl":      ImageAnalysisSkillLanguagePl,
		"prs":     ImageAnalysisSkillLanguagePrs,
		"pt":      ImageAnalysisSkillLanguagePt,
		"pt-br":   ImageAnalysisSkillLanguagePtNegativeBR,
		"pt-pt":   ImageAnalysisSkillLanguagePtNegativePT,
		"ro":      ImageAnalysisSkillLanguageRo,
		"ru":      ImageAnalysisSkillLanguageRu,
		"sk":      ImageAnalysisSkillLanguageSk,
		"sl":      ImageAnalysisSkillLanguageSl,
		"sr-cyrl": ImageAnalysisSkillLanguageSrNegativeCyrl,
		"sr-latn": ImageAnalysisSkillLanguageSrNegativeLatn,
		"sv":      ImageAnalysisSkillLanguageSv,
		"th":      ImageAnalysisSkillLanguageTh,
		"tr":      ImageAnalysisSkillLanguageTr,
		"uk":      ImageAnalysisSkillLanguageUk,
		"vi":      ImageAnalysisSkillLanguageVi,
		"zh":      ImageAnalysisSkillLanguageZh,
		"zh-hans": ImageAnalysisSkillLanguageZhNegativeHans,
		"zh-hant": ImageAnalysisSkillLanguageZhNegativeHant,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ImageAnalysisSkillLanguage(input)
	return &out, nil
}

type ImageDetail string

const (
	ImageDetailCelebrities ImageDetail = "celebrities"
	ImageDetailLandmarks   ImageDetail = "landmarks"
)

func PossibleValuesForImageDetail() []string {
	return []string{
		string(ImageDetailCelebrities),
		string(ImageDetailLandmarks),
	}
}

func (s *ImageDetail) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseImageDetail(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseImageDetail(input string) (*ImageDetail, error) {
	vals := map[string]ImageDetail{
		"celebrities": ImageDetailCelebrities,
		"landmarks":   ImageDetailLandmarks,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ImageDetail(input)
	return &out, nil
}

type IndexProjectionMode string

const (
	IndexProjectionModeIncludeIndexingParentDocuments IndexProjectionMode = "includeIndexingParentDocuments"
	IndexProjectionModeSkipIndexingParentDocuments    IndexProjectionMode = "skipIndexingParentDocuments"
)

func PossibleValuesForIndexProjectionMode() []string {
	return []string{
		string(IndexProjectionModeIncludeIndexingParentDocuments),
		string(IndexProjectionModeSkipIndexingParentDocuments),
	}
}

func (s *IndexProjectionMode) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseIndexProjectionMode(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseIndexProjectionMode(input string) (*IndexProjectionMode, error) {
	vals := map[string]IndexProjectionMode{
		"includeindexingparentdocuments": IndexProjectionModeIncludeIndexingParentDocuments,
		"skipindexingparentdocuments":    IndexProjectionModeSkipIndexingParentDocuments,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := IndexProjectionMode(input)
	return &out, nil
}

type KeyPhraseExtractionSkillLanguage string

const (
	KeyPhraseExtractionSkillLanguageDa           KeyPhraseExtractionSkillLanguage = "da"
	KeyPhraseExtractionSkillLanguageDe           KeyPhraseExtractionSkillLanguage = "de"
	KeyPhraseExtractionSkillLanguageEn           KeyPhraseExtractionSkillLanguage = "en"
	KeyPhraseExtractionSkillLanguageEs           KeyPhraseExtractionSkillLanguage = "es"
	KeyPhraseExtractionSkillLanguageFi           KeyPhraseExtractionSkillLanguage = "fi"
	KeyPhraseExtractionSkillLanguageFr           KeyPhraseExtractionSkillLanguage = "fr"
	KeyPhraseExtractionSkillLanguageIt           KeyPhraseExtractionSkillLanguage = "it"
	KeyPhraseExtractionSkillLanguageJa           KeyPhraseExtractionSkillLanguage = "ja"
	KeyPhraseExtractionSkillLanguageKo           KeyPhraseExtractionSkillLanguage = "ko"
	KeyPhraseExtractionSkillLanguageNl           KeyPhraseExtractionSkillLanguage = "nl"
	KeyPhraseExtractionSkillLanguageNo           KeyPhraseExtractionSkillLanguage = "no"
	KeyPhraseExtractionSkillLanguagePl           KeyPhraseExtractionSkillLanguage = "pl"
	KeyPhraseExtractionSkillLanguagePtNegativeBR KeyPhraseExtractionSkillLanguage = "pt-BR"
	KeyPhraseExtractionSkillLanguagePtNegativePT KeyPhraseExtractionSkillLanguage = "pt-PT"
	KeyPhraseExtractionSkillLanguageRu           KeyPhraseExtractionSkillLanguage = "ru"
	KeyPhraseExtractionSkillLanguageSv           KeyPhraseExtractionSkillLanguage = "sv"
)

func PossibleValuesForKeyPhraseExtractionSkillLanguage() []string {
	return []string{
		string(KeyPhraseExtractionSkillLanguageDa),
		string(KeyPhraseExtractionSkillLanguageDe),
		string(KeyPhraseExtractionSkillLanguageEn),
		string(KeyPhraseExtractionSkillLanguageEs),
		string(KeyPhraseExtractionSkillLanguageFi),
		string(KeyPhraseExtractionSkillLanguageFr),
		string(KeyPhraseExtractionSkillLanguageIt),
		string(KeyPhraseExtractionSkillLanguageJa),
		string(KeyPhraseExtractionSkillLanguageKo),
		string(KeyPhraseExtractionSkillLanguageNl),
		string(KeyPhraseExtractionSkillLanguageNo),
		string(KeyPhraseExtractionSkillLanguagePl),
		string(KeyPhraseExtractionSkillLanguagePtNegativeBR),
		string(KeyPhraseExtractionSkillLanguagePtNegativePT),
		string(KeyPhraseExtractionSkillLanguageRu),
		string(KeyPhraseExtractionSkillLanguageSv),
	}
}

func (s *KeyPhraseExtractionSkillLanguage) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseKeyPhraseExtractionSkillLanguage(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseKeyPhraseExtractionSkillLanguage(input string) (*KeyPhraseExtractionSkillLanguage, error) {
	vals := map[string]KeyPhraseExtractionSkillLanguage{
		"da":    KeyPhraseExtractionSkillLanguageDa,
		"de":    KeyPhraseExtractionSkillLanguageDe,
		"en":    KeyPhraseExtractionSkillLanguageEn,
		"es":    KeyPhraseExtractionSkillLanguageEs,
		"fi":    KeyPhraseExtractionSkillLanguageFi,
		"fr":    KeyPhraseExtractionSkillLanguageFr,
		"it":    KeyPhraseExtractionSkillLanguageIt,
		"ja":    KeyPhraseExtractionSkillLanguageJa,
		"ko":    KeyPhraseExtractionSkillLanguageKo,
		"nl":    KeyPhraseExtractionSkillLanguageNl,
		"no":    KeyPhraseExtractionSkillLanguageNo,
		"pl":    KeyPhraseExtractionSkillLanguagePl,
		"pt-br": KeyPhraseExtractionSkillLanguagePtNegativeBR,
		"pt-pt": KeyPhraseExtractionSkillLanguagePtNegativePT,
		"ru":    KeyPhraseExtractionSkillLanguageRu,
		"sv":    KeyPhraseExtractionSkillLanguageSv,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := KeyPhraseExtractionSkillLanguage(input)
	return &out, nil
}

type LineEnding string

const (
	LineEndingCarriageReturn         LineEnding = "carriageReturn"
	LineEndingCarriageReturnLineFeed LineEnding = "carriageReturnLineFeed"
	LineEndingLineFeed               LineEnding = "lineFeed"
	LineEndingSpace                  LineEnding = "space"
)

func PossibleValuesForLineEnding() []string {
	return []string{
		string(LineEndingCarriageReturn),
		string(LineEndingCarriageReturnLineFeed),
		string(LineEndingLineFeed),
		string(LineEndingSpace),
	}
}

func (s *LineEnding) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseLineEnding(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseLineEnding(input string) (*LineEnding, error) {
	vals := map[string]LineEnding{
		"carriagereturn":         LineEndingCarriageReturn,
		"carriagereturnlinefeed": LineEndingCarriageReturnLineFeed,
		"linefeed":               LineEndingLineFeed,
		"space":                  LineEndingSpace,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := LineEnding(input)
	return &out, nil
}

type OcrSkillLanguage string

const (
	OcrSkillLanguageAf              OcrSkillLanguage = "af"
	OcrSkillLanguageAnp             OcrSkillLanguage = "anp"
	OcrSkillLanguageAr              OcrSkillLanguage = "ar"
	OcrSkillLanguageAst             OcrSkillLanguage = "ast"
	OcrSkillLanguageAwa             OcrSkillLanguage = "awa"
	OcrSkillLanguageAz              OcrSkillLanguage = "az"
	OcrSkillLanguageBe              OcrSkillLanguage = "be"
	OcrSkillLanguageBeNegativecyrl  OcrSkillLanguage = "be-cyrl"
	OcrSkillLanguageBeNegativelatn  OcrSkillLanguage = "be-latn"
	OcrSkillLanguageBfy             OcrSkillLanguage = "bfy"
	OcrSkillLanguageBfz             OcrSkillLanguage = "bfz"
	OcrSkillLanguageBg              OcrSkillLanguage = "bg"
	OcrSkillLanguageBgc             OcrSkillLanguage = "bgc"
	OcrSkillLanguageBho             OcrSkillLanguage = "bho"
	OcrSkillLanguageBi              OcrSkillLanguage = "bi"
	OcrSkillLanguageBns             OcrSkillLanguage = "bns"
	OcrSkillLanguageBr              OcrSkillLanguage = "br"
	OcrSkillLanguageBra             OcrSkillLanguage = "bra"
	OcrSkillLanguageBrx             OcrSkillLanguage = "brx"
	OcrSkillLanguageBs              OcrSkillLanguage = "bs"
	OcrSkillLanguageBua             OcrSkillLanguage = "bua"
	OcrSkillLanguageCa              OcrSkillLanguage = "ca"
	OcrSkillLanguageCeb             OcrSkillLanguage = "ceb"
	OcrSkillLanguageCh              OcrSkillLanguage = "ch"
	OcrSkillLanguageCnrNegativecyrl OcrSkillLanguage = "cnr-cyrl"
	OcrSkillLanguageCnrNegativelatn OcrSkillLanguage = "cnr-latn"
	OcrSkillLanguageCo              OcrSkillLanguage = "co"
	OcrSkillLanguageCrh             OcrSkillLanguage = "crh"
	OcrSkillLanguageCs              OcrSkillLanguage = "cs"
	OcrSkillLanguageCsb             OcrSkillLanguage = "csb"
	OcrSkillLanguageCy              OcrSkillLanguage = "cy"
	OcrSkillLanguageDa              OcrSkillLanguage = "da"
	OcrSkillLanguageDe              OcrSkillLanguage = "de"
	OcrSkillLanguageDhi             OcrSkillLanguage = "dhi"
	OcrSkillLanguageDoi             OcrSkillLanguage = "doi"
	OcrSkillLanguageDsb             OcrSkillLanguage = "dsb"
	OcrSkillLanguageEl              OcrSkillLanguage = "el"
	OcrSkillLanguageEn              OcrSkillLanguage = "en"
	OcrSkillLanguageEs              OcrSkillLanguage = "es"
	OcrSkillLanguageEt              OcrSkillLanguage = "et"
	OcrSkillLanguageEu              OcrSkillLanguage = "eu"
	OcrSkillLanguageFa              OcrSkillLanguage = "fa"
	OcrSkillLanguageFi              OcrSkillLanguage = "fi"
	OcrSkillLanguageFil             OcrSkillLanguage = "fil"
	OcrSkillLanguageFj              OcrSkillLanguage = "fj"
	OcrSkillLanguageFo              OcrSkillLanguage = "fo"
	OcrSkillLanguageFr              OcrSkillLanguage = "fr"
	OcrSkillLanguageFur             OcrSkillLanguage = "fur"
	OcrSkillLanguageFy              OcrSkillLanguage = "fy"
	OcrSkillLanguageGa              OcrSkillLanguage = "ga"
	OcrSkillLanguageGag             OcrSkillLanguage = "gag"
	OcrSkillLanguageGd              OcrSkillLanguage = "gd"
	OcrSkillLanguageGil             OcrSkillLanguage = "gil"
	OcrSkillLanguageGl              OcrSkillLanguage = "gl"
	OcrSkillLanguageGon             OcrSkillLanguage = "gon"
	OcrSkillLanguageGv              OcrSkillLanguage = "gv"
	OcrSkillLanguageGvr             OcrSkillLanguage = "gvr"
	OcrSkillLanguageHaw             OcrSkillLanguage = "haw"
	OcrSkillLanguageHi              OcrSkillLanguage = "hi"
	OcrSkillLanguageHlb             OcrSkillLanguage = "hlb"
	OcrSkillLanguageHne             OcrSkillLanguage = "hne"
	OcrSkillLanguageHni             OcrSkillLanguage = "hni"
	OcrSkillLanguageHoc             OcrSkillLanguage = "hoc"
	OcrSkillLanguageHr              OcrSkillLanguage = "hr"
	OcrSkillLanguageHsb             OcrSkillLanguage = "hsb"
	OcrSkillLanguageHt              OcrSkillLanguage = "ht"
	OcrSkillLanguageHu              OcrSkillLanguage = "hu"
	OcrSkillLanguageIa              OcrSkillLanguage = "ia"
	OcrSkillLanguageId              OcrSkillLanguage = "id"
	OcrSkillLanguageIs              OcrSkillLanguage = "is"
	OcrSkillLanguageIt              OcrSkillLanguage = "it"
	OcrSkillLanguageIu              OcrSkillLanguage = "iu"
	OcrSkillLanguageJa              OcrSkillLanguage = "ja"
	OcrSkillLanguageJns             OcrSkillLanguage = "Jns"
	OcrSkillLanguageJv              OcrSkillLanguage = "jv"
	OcrSkillLanguageKaa             OcrSkillLanguage = "kaa"
	OcrSkillLanguageKaaNegativecyrl OcrSkillLanguage = "kaa-cyrl"
	OcrSkillLanguageKac             OcrSkillLanguage = "kac"
	OcrSkillLanguageKea             OcrSkillLanguage = "kea"
	OcrSkillLanguageKfq             OcrSkillLanguage = "kfq"
	OcrSkillLanguageKha             OcrSkillLanguage = "kha"
	OcrSkillLanguageKkNegativecyrl  OcrSkillLanguage = "kk-cyrl"
	OcrSkillLanguageKkNegativelatn  OcrSkillLanguage = "kk-latn"
	OcrSkillLanguageKl              OcrSkillLanguage = "kl"
	OcrSkillLanguageKlr             OcrSkillLanguage = "klr"
	OcrSkillLanguageKmj             OcrSkillLanguage = "kmj"
	OcrSkillLanguageKo              OcrSkillLanguage = "ko"
	OcrSkillLanguageKos             OcrSkillLanguage = "kos"
	OcrSkillLanguageKpy             OcrSkillLanguage = "kpy"
	OcrSkillLanguageKrc             OcrSkillLanguage = "krc"
	OcrSkillLanguageKru             OcrSkillLanguage = "kru"
	OcrSkillLanguageKsh             OcrSkillLanguage = "ksh"
	OcrSkillLanguageKuNegativearab  OcrSkillLanguage = "ku-arab"
	OcrSkillLanguageKuNegativelatn  OcrSkillLanguage = "ku-latn"
	OcrSkillLanguageKum             OcrSkillLanguage = "kum"
	OcrSkillLanguageKw              OcrSkillLanguage = "kw"
	OcrSkillLanguageKy              OcrSkillLanguage = "ky"
	OcrSkillLanguageLa              OcrSkillLanguage = "la"
	OcrSkillLanguageLb              OcrSkillLanguage = "lb"
	OcrSkillLanguageLkt             OcrSkillLanguage = "lkt"
	OcrSkillLanguageLt              OcrSkillLanguage = "lt"
	OcrSkillLanguageMi              OcrSkillLanguage = "mi"
	OcrSkillLanguageMn              OcrSkillLanguage = "mn"
	OcrSkillLanguageMr              OcrSkillLanguage = "mr"
	OcrSkillLanguageMs              OcrSkillLanguage = "ms"
	OcrSkillLanguageMt              OcrSkillLanguage = "mt"
	OcrSkillLanguageMww             OcrSkillLanguage = "mww"
	OcrSkillLanguageMyv             OcrSkillLanguage = "myv"
	OcrSkillLanguageNap             OcrSkillLanguage = "nap"
	OcrSkillLanguageNb              OcrSkillLanguage = "nb"
	OcrSkillLanguageNe              OcrSkillLanguage = "ne"
	OcrSkillLanguageNiu             OcrSkillLanguage = "niu"
	OcrSkillLanguageNl              OcrSkillLanguage = "nl"
	OcrSkillLanguageNo              OcrSkillLanguage = "no"
	OcrSkillLanguageNog             OcrSkillLanguage = "nog"
	OcrSkillLanguageOc              OcrSkillLanguage = "oc"
	OcrSkillLanguageOs              OcrSkillLanguage = "os"
	OcrSkillLanguagePa              OcrSkillLanguage = "pa"
	OcrSkillLanguagePl              OcrSkillLanguage = "pl"
	OcrSkillLanguagePrs             OcrSkillLanguage = "prs"
	OcrSkillLanguagePs              OcrSkillLanguage = "ps"
	OcrSkillLanguagePt              OcrSkillLanguage = "pt"
	OcrSkillLanguageQuc             OcrSkillLanguage = "quc"
	OcrSkillLanguageRab             OcrSkillLanguage = "rab"
	OcrSkillLanguageRm              OcrSkillLanguage = "rm"
	OcrSkillLanguageRo              OcrSkillLanguage = "ro"
	OcrSkillLanguageRu              OcrSkillLanguage = "ru"
	OcrSkillLanguageSa              OcrSkillLanguage = "sa"
	OcrSkillLanguageSat             OcrSkillLanguage = "sat"
	OcrSkillLanguageSck             OcrSkillLanguage = "sck"
	OcrSkillLanguageSco             OcrSkillLanguage = "sco"
	OcrSkillLanguageSk              OcrSkillLanguage = "sk"
	OcrSkillLanguageSl              OcrSkillLanguage = "sl"
	OcrSkillLanguageSm              OcrSkillLanguage = "sm"
	OcrSkillLanguageSma             OcrSkillLanguage = "sma"
	OcrSkillLanguageSme             OcrSkillLanguage = "sme"
	OcrSkillLanguageSmj             OcrSkillLanguage = "smj"
	OcrSkillLanguageSmn             OcrSkillLanguage = "smn"
	OcrSkillLanguageSms             OcrSkillLanguage = "sms"
	OcrSkillLanguageSo              OcrSkillLanguage = "so"
	OcrSkillLanguageSq              OcrSkillLanguage = "sq"
	OcrSkillLanguageSr              OcrSkillLanguage = "sr"
	OcrSkillLanguageSrNegativeCyrl  OcrSkillLanguage = "sr-Cyrl"
	OcrSkillLanguageSrNegativeLatn  OcrSkillLanguage = "sr-Latn"
	OcrSkillLanguageSrx             OcrSkillLanguage = "srx"
	OcrSkillLanguageSv              OcrSkillLanguage = "sv"
	OcrSkillLanguageSw              OcrSkillLanguage = "sw"
	OcrSkillLanguageTet             OcrSkillLanguage = "tet"
	OcrSkillLanguageTg              OcrSkillLanguage = "tg"
	OcrSkillLanguageThf             OcrSkillLanguage = "thf"
	OcrSkillLanguageTk              OcrSkillLanguage = "tk"
	OcrSkillLanguageTo              OcrSkillLanguage = "to"
	OcrSkillLanguageTr              OcrSkillLanguage = "tr"
	OcrSkillLanguageTt              OcrSkillLanguage = "tt"
	OcrSkillLanguageTyv             OcrSkillLanguage = "tyv"
	OcrSkillLanguageUg              OcrSkillLanguage = "ug"
	OcrSkillLanguageUnk             OcrSkillLanguage = "unk"
	OcrSkillLanguageUr              OcrSkillLanguage = "ur"
	OcrSkillLanguageUz              OcrSkillLanguage = "uz"
	OcrSkillLanguageUzNegativearab  OcrSkillLanguage = "uz-arab"
	OcrSkillLanguageUzNegativecyrl  OcrSkillLanguage = "uz-cyrl"
	OcrSkillLanguageVo              OcrSkillLanguage = "vo"
	OcrSkillLanguageWae             OcrSkillLanguage = "wae"
	OcrSkillLanguageXnr             OcrSkillLanguage = "xnr"
	OcrSkillLanguageXsr             OcrSkillLanguage = "xsr"
	OcrSkillLanguageYua             OcrSkillLanguage = "yua"
	OcrSkillLanguageZa              OcrSkillLanguage = "za"
	OcrSkillLanguageZhNegativeHans  OcrSkillLanguage = "zh-Hans"
	OcrSkillLanguageZhNegativeHant  OcrSkillLanguage = "zh-Hant"
	OcrSkillLanguageZu              OcrSkillLanguage = "zu"
)

func PossibleValuesForOcrSkillLanguage() []string {
	return []string{
		string(OcrSkillLanguageAf),
		string(OcrSkillLanguageAnp),
		string(OcrSkillLanguageAr),
		string(OcrSkillLanguageAst),
		string(OcrSkillLanguageAwa),
		string(OcrSkillLanguageAz),
		string(OcrSkillLanguageBe),
		string(OcrSkillLanguageBeNegativecyrl),
		string(OcrSkillLanguageBeNegativelatn),
		string(OcrSkillLanguageBfy),
		string(OcrSkillLanguageBfz),
		string(OcrSkillLanguageBg),
		string(OcrSkillLanguageBgc),
		string(OcrSkillLanguageBho),
		string(OcrSkillLanguageBi),
		string(OcrSkillLanguageBns),
		string(OcrSkillLanguageBr),
		string(OcrSkillLanguageBra),
		string(OcrSkillLanguageBrx),
		string(OcrSkillLanguageBs),
		string(OcrSkillLanguageBua),
		string(OcrSkillLanguageCa),
		string(OcrSkillLanguageCeb),
		string(OcrSkillLanguageCh),
		string(OcrSkillLanguageCnrNegativecyrl),
		string(OcrSkillLanguageCnrNegativelatn),
		string(OcrSkillLanguageCo),
		string(OcrSkillLanguageCrh),
		string(OcrSkillLanguageCs),
		string(OcrSkillLanguageCsb),
		string(OcrSkillLanguageCy),
		string(OcrSkillLanguageDa),
		string(OcrSkillLanguageDe),
		string(OcrSkillLanguageDhi),
		string(OcrSkillLanguageDoi),
		string(OcrSkillLanguageDsb),
		string(OcrSkillLanguageEl),
		string(OcrSkillLanguageEn),
		string(OcrSkillLanguageEs),
		string(OcrSkillLanguageEt),
		string(OcrSkillLanguageEu),
		string(OcrSkillLanguageFa),
		string(OcrSkillLanguageFi),
		string(OcrSkillLanguageFil),
		string(OcrSkillLanguageFj),
		string(OcrSkillLanguageFo),
		string(OcrSkillLanguageFr),
		string(OcrSkillLanguageFur),
		string(OcrSkillLanguageFy),
		string(OcrSkillLanguageGa),
		string(OcrSkillLanguageGag),
		string(OcrSkillLanguageGd),
		string(OcrSkillLanguageGil),
		string(OcrSkillLanguageGl),
		string(OcrSkillLanguageGon),
		string(OcrSkillLanguageGv),
		string(OcrSkillLanguageGvr),
		string(OcrSkillLanguageHaw),
		string(OcrSkillLanguageHi),
		string(OcrSkillLanguageHlb),
		string(OcrSkillLanguageHne),
		string(OcrSkillLanguageHni),
		string(OcrSkillLanguageHoc),
		string(OcrSkillLanguageHr),
		string(OcrSkillLanguageHsb),
		string(OcrSkillLanguageHt),
		string(OcrSkillLanguageHu),
		string(OcrSkillLanguageIa),
		string(OcrSkillLanguageId),
		string(OcrSkillLanguageIs),
		string(OcrSkillLanguageIt),
		string(OcrSkillLanguageIu),
		string(OcrSkillLanguageJa),
		string(OcrSkillLanguageJns),
		string(OcrSkillLanguageJv),
		string(OcrSkillLanguageKaa),
		string(OcrSkillLanguageKaaNegativecyrl),
		string(OcrSkillLanguageKac),
		string(OcrSkillLanguageKea),
		string(OcrSkillLanguageKfq),
		string(OcrSkillLanguageKha),
		string(OcrSkillLanguageKkNegativecyrl),
		string(OcrSkillLanguageKkNegativelatn),
		string(OcrSkillLanguageKl),
		string(OcrSkillLanguageKlr),
		string(OcrSkillLanguageKmj),
		string(OcrSkillLanguageKo),
		string(OcrSkillLanguageKos),
		string(OcrSkillLanguageKpy),
		string(OcrSkillLanguageKrc),
		string(OcrSkillLanguageKru),
		string(OcrSkillLanguageKsh),
		string(OcrSkillLanguageKuNegativearab),
		string(OcrSkillLanguageKuNegativelatn),
		string(OcrSkillLanguageKum),
		string(OcrSkillLanguageKw),
		string(OcrSkillLanguageKy),
		string(OcrSkillLanguageLa),
		string(OcrSkillLanguageLb),
		string(OcrSkillLanguageLkt),
		string(OcrSkillLanguageLt),
		string(OcrSkillLanguageMi),
		string(OcrSkillLanguageMn),
		string(OcrSkillLanguageMr),
		string(OcrSkillLanguageMs),
		string(OcrSkillLanguageMt),
		string(OcrSkillLanguageMww),
		string(OcrSkillLanguageMyv),
		string(OcrSkillLanguageNap),
		string(OcrSkillLanguageNb),
		string(OcrSkillLanguageNe),
		string(OcrSkillLanguageNiu),
		string(OcrSkillLanguageNl),
		string(OcrSkillLanguageNo),
		string(OcrSkillLanguageNog),
		string(OcrSkillLanguageOc),
		string(OcrSkillLanguageOs),
		string(OcrSkillLanguagePa),
		string(OcrSkillLanguagePl),
		string(OcrSkillLanguagePrs),
		string(OcrSkillLanguagePs),
		string(OcrSkillLanguagePt),
		string(OcrSkillLanguageQuc),
		string(OcrSkillLanguageRab),
		string(OcrSkillLanguageRm),
		string(OcrSkillLanguageRo),
		string(OcrSkillLanguageRu),
		string(OcrSkillLanguageSa),
		string(OcrSkillLanguageSat),
		string(OcrSkillLanguageSck),
		string(OcrSkillLanguageSco),
		string(OcrSkillLanguageSk),
		string(OcrSkillLanguageSl),
		string(OcrSkillLanguageSm),
		string(OcrSkillLanguageSma),
		string(OcrSkillLanguageSme),
		string(OcrSkillLanguageSmj),
		string(OcrSkillLanguageSmn),
		string(OcrSkillLanguageSms),
		string(OcrSkillLanguageSo),
		string(OcrSkillLanguageSq),
		string(OcrSkillLanguageSr),
		string(OcrSkillLanguageSrNegativeCyrl),
		string(OcrSkillLanguageSrNegativeLatn),
		string(OcrSkillLanguageSrx),
		string(OcrSkillLanguageSv),
		string(OcrSkillLanguageSw),
		string(OcrSkillLanguageTet),
		string(OcrSkillLanguageTg),
		string(OcrSkillLanguageThf),
		string(OcrSkillLanguageTk),
		string(OcrSkillLanguageTo),
		string(OcrSkillLanguageTr),
		string(OcrSkillLanguageTt),
		string(OcrSkillLanguageTyv),
		string(OcrSkillLanguageUg),
		string(OcrSkillLanguageUnk),
		string(OcrSkillLanguageUr),
		string(OcrSkillLanguageUz),
		string(OcrSkillLanguageUzNegativearab),
		string(OcrSkillLanguageUzNegativecyrl),
		string(OcrSkillLanguageVo),
		string(OcrSkillLanguageWae),
		string(OcrSkillLanguageXnr),
		string(OcrSkillLanguageXsr),
		string(OcrSkillLanguageYua),
		string(OcrSkillLanguageZa),
		string(OcrSkillLanguageZhNegativeHans),
		string(OcrSkillLanguageZhNegativeHant),
		string(OcrSkillLanguageZu),
	}
}

func (s *OcrSkillLanguage) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseOcrSkillLanguage(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseOcrSkillLanguage(input string) (*OcrSkillLanguage, error) {
	vals := map[string]OcrSkillLanguage{
		"af":       OcrSkillLanguageAf,
		"anp":      OcrSkillLanguageAnp,
		"ar":       OcrSkillLanguageAr,
		"ast":      OcrSkillLanguageAst,
		"awa":      OcrSkillLanguageAwa,
		"az":       OcrSkillLanguageAz,
		"be":       OcrSkillLanguageBe,
		"be-cyrl":  OcrSkillLanguageBeNegativecyrl,
		"be-latn":  OcrSkillLanguageBeNegativelatn,
		"bfy":      OcrSkillLanguageBfy,
		"bfz":      OcrSkillLanguageBfz,
		"bg":       OcrSkillLanguageBg,
		"bgc":      OcrSkillLanguageBgc,
		"bho":      OcrSkillLanguageBho,
		"bi":       OcrSkillLanguageBi,
		"bns":      OcrSkillLanguageBns,
		"br":       OcrSkillLanguageBr,
		"bra":      OcrSkillLanguageBra,
		"brx":      OcrSkillLanguageBrx,
		"bs":       OcrSkillLanguageBs,
		"bua":      OcrSkillLanguageBua,
		"ca":       OcrSkillLanguageCa,
		"ceb":      OcrSkillLanguageCeb,
		"ch":       OcrSkillLanguageCh,
		"cnr-cyrl": OcrSkillLanguageCnrNegativecyrl,
		"cnr-latn": OcrSkillLanguageCnrNegativelatn,
		"co":       OcrSkillLanguageCo,
		"crh":      OcrSkillLanguageCrh,
		"cs":       OcrSkillLanguageCs,
		"csb":      OcrSkillLanguageCsb,
		"cy":       OcrSkillLanguageCy,
		"da":       OcrSkillLanguageDa,
		"de":       OcrSkillLanguageDe,
		"dhi":      OcrSkillLanguageDhi,
		"doi":      OcrSkillLanguageDoi,
		"dsb":      OcrSkillLanguageDsb,
		"el":       OcrSkillLanguageEl,
		"en":       OcrSkillLanguageEn,
		"es":       OcrSkillLanguageEs,
		"et":       OcrSkillLanguageEt,
		"eu":       OcrSkillLanguageEu,
		"fa":       OcrSkillLanguageFa,
		"fi":       OcrSkillLanguageFi,
		"fil":      OcrSkillLanguageFil,
		"fj":       OcrSkillLanguageFj,
		"fo":       OcrSkillLanguageFo,
		"fr":       OcrSkillLanguageFr,
		"fur":      OcrSkillLanguageFur,
		"fy":       OcrSkillLanguageFy,
		"ga":       OcrSkillLanguageGa,
		"gag":      OcrSkillLanguageGag,
		"gd":       OcrSkillLanguageGd,
		"gil":      OcrSkillLanguageGil,
		"gl":       OcrSkillLanguageGl,
		"gon":      OcrSkillLanguageGon,
		"gv":       OcrSkillLanguageGv,
		"gvr":      OcrSkillLanguageGvr,
		"haw":      OcrSkillLanguageHaw,
		"hi":       OcrSkillLanguageHi,
		"hlb":      OcrSkillLanguageHlb,
		"hne":      OcrSkillLanguageHne,
		"hni":      OcrSkillLanguageHni,
		"hoc":      OcrSkillLanguageHoc,
		"hr":       OcrSkillLanguageHr,
		"hsb":      OcrSkillLanguageHsb,
		"ht":       OcrSkillLanguageHt,
		"hu":       OcrSkillLanguageHu,
		"ia":       OcrSkillLanguageIa,
		"id":       OcrSkillLanguageId,
		"is":       OcrSkillLanguageIs,
		"it":       OcrSkillLanguageIt,
		"iu":       OcrSkillLanguageIu,
		"ja":       OcrSkillLanguageJa,
		"jns":      OcrSkillLanguageJns,
		"jv":       OcrSkillLanguageJv,
		"kaa":      OcrSkillLanguageKaa,
		"kaa-cyrl": OcrSkillLanguageKaaNegativecyrl,
		"kac":      OcrSkillLanguageKac,
		"kea":      OcrSkillLanguageKea,
		"kfq":      OcrSkillLanguageKfq,
		"kha":      OcrSkillLanguageKha,
		"kk-cyrl":  OcrSkillLanguageKkNegativecyrl,
		"kk-latn":  OcrSkillLanguageKkNegativelatn,
		"kl":       OcrSkillLanguageKl,
		"klr":      OcrSkillLanguageKlr,
		"kmj":      OcrSkillLanguageKmj,
		"ko":       OcrSkillLanguageKo,
		"kos":      OcrSkillLanguageKos,
		"kpy":      OcrSkillLanguageKpy,
		"krc":      OcrSkillLanguageKrc,
		"kru":      OcrSkillLanguageKru,
		"ksh":      OcrSkillLanguageKsh,
		"ku-arab":  OcrSkillLanguageKuNegativearab,
		"ku-latn":  OcrSkillLanguageKuNegativelatn,
		"kum":      OcrSkillLanguageKum,
		"kw":       OcrSkillLanguageKw,
		"ky":       OcrSkillLanguageKy,
		"la":       OcrSkillLanguageLa,
		"lb":       OcrSkillLanguageLb,
		"lkt":      OcrSkillLanguageLkt,
		"lt":       OcrSkillLanguageLt,
		"mi":       OcrSkillLanguageMi,
		"mn":       OcrSkillLanguageMn,
		"mr":       OcrSkillLanguageMr,
		"ms":       OcrSkillLanguageMs,
		"mt":       OcrSkillLanguageMt,
		"mww":      OcrSkillLanguageMww,
		"myv":      OcrSkillLanguageMyv,
		"nap":      OcrSkillLanguageNap,
		"nb":       OcrSkillLanguageNb,
		"ne":       OcrSkillLanguageNe,
		"niu":      OcrSkillLanguageNiu,
		"nl":       OcrSkillLanguageNl,
		"no":       OcrSkillLanguageNo,
		"nog":      OcrSkillLanguageNog,
		"oc":       OcrSkillLanguageOc,
		"os":       OcrSkillLanguageOs,
		"pa":       OcrSkillLanguagePa,
		"pl":       OcrSkillLanguagePl,
		"prs":      OcrSkillLanguagePrs,
		"ps":       OcrSkillLanguagePs,
		"pt":       OcrSkillLanguagePt,
		"quc":      OcrSkillLanguageQuc,
		"rab":      OcrSkillLanguageRab,
		"rm":       OcrSkillLanguageRm,
		"ro":       OcrSkillLanguageRo,
		"ru":       OcrSkillLanguageRu,
		"sa":       OcrSkillLanguageSa,
		"sat":      OcrSkillLanguageSat,
		"sck":      OcrSkillLanguageSck,
		"sco":      OcrSkillLanguageSco,
		"sk":       OcrSkillLanguageSk,
		"sl":       OcrSkillLanguageSl,
		"sm":       OcrSkillLanguageSm,
		"sma":      OcrSkillLanguageSma,
		"sme":      OcrSkillLanguageSme,
		"smj":      OcrSkillLanguageSmj,
		"smn":      OcrSkillLanguageSmn,
		"sms":      OcrSkillLanguageSms,
		"so":       OcrSkillLanguageSo,
		"sq":       OcrSkillLanguageSq,
		"sr":       OcrSkillLanguageSr,
		"sr-cyrl":  OcrSkillLanguageSrNegativeCyrl,
		"sr-latn":  OcrSkillLanguageSrNegativeLatn,
		"srx":      OcrSkillLanguageSrx,
		"sv":       OcrSkillLanguageSv,
		"sw":       OcrSkillLanguageSw,
		"tet":      OcrSkillLanguageTet,
		"tg":       OcrSkillLanguageTg,
		"thf":      OcrSkillLanguageThf,
		"tk":       OcrSkillLanguageTk,
		"to":       OcrSkillLanguageTo,
		"tr":       OcrSkillLanguageTr,
		"tt":       OcrSkillLanguageTt,
		"tyv":      OcrSkillLanguageTyv,
		"ug":       OcrSkillLanguageUg,
		"unk":      OcrSkillLanguageUnk,
		"ur":       OcrSkillLanguageUr,
		"uz":       OcrSkillLanguageUz,
		"uz-arab":  OcrSkillLanguageUzNegativearab,
		"uz-cyrl":  OcrSkillLanguageUzNegativecyrl,
		"vo":       OcrSkillLanguageVo,
		"wae":      OcrSkillLanguageWae,
		"xnr":      OcrSkillLanguageXnr,
		"xsr":      OcrSkillLanguageXsr,
		"yua":      OcrSkillLanguageYua,
		"za":       OcrSkillLanguageZa,
		"zh-hans":  OcrSkillLanguageZhNegativeHans,
		"zh-hant":  OcrSkillLanguageZhNegativeHant,
		"zu":       OcrSkillLanguageZu,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := OcrSkillLanguage(input)
	return &out, nil
}

type PIIDetectionSkillMaskingMode string

const (
	PIIDetectionSkillMaskingModeNone    PIIDetectionSkillMaskingMode = "none"
	PIIDetectionSkillMaskingModeReplace PIIDetectionSkillMaskingMode = "replace"
)

func PossibleValuesForPIIDetectionSkillMaskingMode() []string {
	return []string{
		string(PIIDetectionSkillMaskingModeNone),
		string(PIIDetectionSkillMaskingModeReplace),
	}
}

func (s *PIIDetectionSkillMaskingMode) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parsePIIDetectionSkillMaskingMode(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parsePIIDetectionSkillMaskingMode(input string) (*PIIDetectionSkillMaskingMode, error) {
	vals := map[string]PIIDetectionSkillMaskingMode{
		"none":    PIIDetectionSkillMaskingModeNone,
		"replace": PIIDetectionSkillMaskingModeReplace,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := PIIDetectionSkillMaskingMode(input)
	return &out, nil
}

type Prefer string

const (
	PreferReturnRepresentation Prefer = "return=representation"
)

func PossibleValuesForPrefer() []string {
	return []string{
		string(PreferReturnRepresentation),
	}
}

func (s *Prefer) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parsePrefer(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parsePrefer(input string) (*Prefer, error) {
	vals := map[string]Prefer{
		"return=representation": PreferReturnRepresentation,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := Prefer(input)
	return &out, nil
}

type SentimentSkillLanguage string

const (
	SentimentSkillLanguageDa           SentimentSkillLanguage = "da"
	SentimentSkillLanguageDe           SentimentSkillLanguage = "de"
	SentimentSkillLanguageEl           SentimentSkillLanguage = "el"
	SentimentSkillLanguageEn           SentimentSkillLanguage = "en"
	SentimentSkillLanguageEs           SentimentSkillLanguage = "es"
	SentimentSkillLanguageFi           SentimentSkillLanguage = "fi"
	SentimentSkillLanguageFr           SentimentSkillLanguage = "fr"
	SentimentSkillLanguageIt           SentimentSkillLanguage = "it"
	SentimentSkillLanguageNl           SentimentSkillLanguage = "nl"
	SentimentSkillLanguageNo           SentimentSkillLanguage = "no"
	SentimentSkillLanguagePl           SentimentSkillLanguage = "pl"
	SentimentSkillLanguagePtNegativePT SentimentSkillLanguage = "pt-PT"
	SentimentSkillLanguageRu           SentimentSkillLanguage = "ru"
	SentimentSkillLanguageSv           SentimentSkillLanguage = "sv"
	SentimentSkillLanguageTr           SentimentSkillLanguage = "tr"
)

func PossibleValuesForSentimentSkillLanguage() []string {
	return []string{
		string(SentimentSkillLanguageDa),
		string(SentimentSkillLanguageDe),
		string(SentimentSkillLanguageEl),
		string(SentimentSkillLanguageEn),
		string(SentimentSkillLanguageEs),
		string(SentimentSkillLanguageFi),
		string(SentimentSkillLanguageFr),
		string(SentimentSkillLanguageIt),
		string(SentimentSkillLanguageNl),
		string(SentimentSkillLanguageNo),
		string(SentimentSkillLanguagePl),
		string(SentimentSkillLanguagePtNegativePT),
		string(SentimentSkillLanguageRu),
		string(SentimentSkillLanguageSv),
		string(SentimentSkillLanguageTr),
	}
}

func (s *SentimentSkillLanguage) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseSentimentSkillLanguage(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseSentimentSkillLanguage(input string) (*SentimentSkillLanguage, error) {
	vals := map[string]SentimentSkillLanguage{
		"da":    SentimentSkillLanguageDa,
		"de":    SentimentSkillLanguageDe,
		"el":    SentimentSkillLanguageEl,
		"en":    SentimentSkillLanguageEn,
		"es":    SentimentSkillLanguageEs,
		"fi":    SentimentSkillLanguageFi,
		"fr":    SentimentSkillLanguageFr,
		"it":    SentimentSkillLanguageIt,
		"nl":    SentimentSkillLanguageNl,
		"no":    SentimentSkillLanguageNo,
		"pl":    SentimentSkillLanguagePl,
		"pt-pt": SentimentSkillLanguagePtNegativePT,
		"ru":    SentimentSkillLanguageRu,
		"sv":    SentimentSkillLanguageSv,
		"tr":    SentimentSkillLanguageTr,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := SentimentSkillLanguage(input)
	return &out, nil
}

type SplitSkillLanguage string

const (
	SplitSkillLanguageAm           SplitSkillLanguage = "am"
	SplitSkillLanguageBs           SplitSkillLanguage = "bs"
	SplitSkillLanguageCs           SplitSkillLanguage = "cs"
	SplitSkillLanguageDa           SplitSkillLanguage = "da"
	SplitSkillLanguageDe           SplitSkillLanguage = "de"
	SplitSkillLanguageEn           SplitSkillLanguage = "en"
	SplitSkillLanguageEs           SplitSkillLanguage = "es"
	SplitSkillLanguageEt           SplitSkillLanguage = "et"
	SplitSkillLanguageFi           SplitSkillLanguage = "fi"
	SplitSkillLanguageFr           SplitSkillLanguage = "fr"
	SplitSkillLanguageHe           SplitSkillLanguage = "he"
	SplitSkillLanguageHi           SplitSkillLanguage = "hi"
	SplitSkillLanguageHr           SplitSkillLanguage = "hr"
	SplitSkillLanguageHu           SplitSkillLanguage = "hu"
	SplitSkillLanguageId           SplitSkillLanguage = "id"
	SplitSkillLanguageIs           SplitSkillLanguage = "is"
	SplitSkillLanguageIt           SplitSkillLanguage = "it"
	SplitSkillLanguageJa           SplitSkillLanguage = "ja"
	SplitSkillLanguageKo           SplitSkillLanguage = "ko"
	SplitSkillLanguageLv           SplitSkillLanguage = "lv"
	SplitSkillLanguageNb           SplitSkillLanguage = "nb"
	SplitSkillLanguageNl           SplitSkillLanguage = "nl"
	SplitSkillLanguagePl           SplitSkillLanguage = "pl"
	SplitSkillLanguagePt           SplitSkillLanguage = "pt"
	SplitSkillLanguagePtNegativebr SplitSkillLanguage = "pt-br"
	SplitSkillLanguageRu           SplitSkillLanguage = "ru"
	SplitSkillLanguageSk           SplitSkillLanguage = "sk"
	SplitSkillLanguageSl           SplitSkillLanguage = "sl"
	SplitSkillLanguageSr           SplitSkillLanguage = "sr"
	SplitSkillLanguageSv           SplitSkillLanguage = "sv"
	SplitSkillLanguageTr           SplitSkillLanguage = "tr"
	SplitSkillLanguageUr           SplitSkillLanguage = "ur"
	SplitSkillLanguageZh           SplitSkillLanguage = "zh"
)

func PossibleValuesForSplitSkillLanguage() []string {
	return []string{
		string(SplitSkillLanguageAm),
		string(SplitSkillLanguageBs),
		string(SplitSkillLanguageCs),
		string(SplitSkillLanguageDa),
		string(SplitSkillLanguageDe),
		string(SplitSkillLanguageEn),
		string(SplitSkillLanguageEs),
		string(SplitSkillLanguageEt),
		string(SplitSkillLanguageFi),
		string(SplitSkillLanguageFr),
		string(SplitSkillLanguageHe),
		string(SplitSkillLanguageHi),
		string(SplitSkillLanguageHr),
		string(SplitSkillLanguageHu),
		string(SplitSkillLanguageId),
		string(SplitSkillLanguageIs),
		string(SplitSkillLanguageIt),
		string(SplitSkillLanguageJa),
		string(SplitSkillLanguageKo),
		string(SplitSkillLanguageLv),
		string(SplitSkillLanguageNb),
		string(SplitSkillLanguageNl),
		string(SplitSkillLanguagePl),
		string(SplitSkillLanguagePt),
		string(SplitSkillLanguagePtNegativebr),
		string(SplitSkillLanguageRu),
		string(SplitSkillLanguageSk),
		string(SplitSkillLanguageSl),
		string(SplitSkillLanguageSr),
		string(SplitSkillLanguageSv),
		string(SplitSkillLanguageTr),
		string(SplitSkillLanguageUr),
		string(SplitSkillLanguageZh),
	}
}

func (s *SplitSkillLanguage) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseSplitSkillLanguage(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseSplitSkillLanguage(input string) (*SplitSkillLanguage, error) {
	vals := map[string]SplitSkillLanguage{
		"am":    SplitSkillLanguageAm,
		"bs":    SplitSkillLanguageBs,
		"cs":    SplitSkillLanguageCs,
		"da":    SplitSkillLanguageDa,
		"de":    SplitSkillLanguageDe,
		"en":    SplitSkillLanguageEn,
		"es":    SplitSkillLanguageEs,
		"et":    SplitSkillLanguageEt,
		"fi":    SplitSkillLanguageFi,
		"fr":    SplitSkillLanguageFr,
		"he":    SplitSkillLanguageHe,
		"hi":    SplitSkillLanguageHi,
		"hr":    SplitSkillLanguageHr,
		"hu":    SplitSkillLanguageHu,
		"id":    SplitSkillLanguageId,
		"is":    SplitSkillLanguageIs,
		"it":    SplitSkillLanguageIt,
		"ja":    SplitSkillLanguageJa,
		"ko":    SplitSkillLanguageKo,
		"lv":    SplitSkillLanguageLv,
		"nb":    SplitSkillLanguageNb,
		"nl":    SplitSkillLanguageNl,
		"pl":    SplitSkillLanguagePl,
		"pt":    SplitSkillLanguagePt,
		"pt-br": SplitSkillLanguagePtNegativebr,
		"ru":    SplitSkillLanguageRu,
		"sk":    SplitSkillLanguageSk,
		"sl":    SplitSkillLanguageSl,
		"sr":    SplitSkillLanguageSr,
		"sv":    SplitSkillLanguageSv,
		"tr":    SplitSkillLanguageTr,
		"ur":    SplitSkillLanguageUr,
		"zh":    SplitSkillLanguageZh,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := SplitSkillLanguage(input)
	return &out, nil
}

type TextSplitMode string

const (
	TextSplitModePages     TextSplitMode = "pages"
	TextSplitModeSentences TextSplitMode = "sentences"
)

func PossibleValuesForTextSplitMode() []string {
	return []string{
		string(TextSplitModePages),
		string(TextSplitModeSentences),
	}
}

func (s *TextSplitMode) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseTextSplitMode(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseTextSplitMode(input string) (*TextSplitMode, error) {
	vals := map[string]TextSplitMode{
		"pages":     TextSplitModePages,
		"sentences": TextSplitModeSentences,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := TextSplitMode(input)
	return &out, nil
}

type TextTranslationSkillLanguage string

const (
	TextTranslationSkillLanguageAf              TextTranslationSkillLanguage = "af"
	TextTranslationSkillLanguageAr              TextTranslationSkillLanguage = "ar"
	TextTranslationSkillLanguageBg              TextTranslationSkillLanguage = "bg"
	TextTranslationSkillLanguageBn              TextTranslationSkillLanguage = "bn"
	TextTranslationSkillLanguageBs              TextTranslationSkillLanguage = "bs"
	TextTranslationSkillLanguageCa              TextTranslationSkillLanguage = "ca"
	TextTranslationSkillLanguageCs              TextTranslationSkillLanguage = "cs"
	TextTranslationSkillLanguageCy              TextTranslationSkillLanguage = "cy"
	TextTranslationSkillLanguageDa              TextTranslationSkillLanguage = "da"
	TextTranslationSkillLanguageDe              TextTranslationSkillLanguage = "de"
	TextTranslationSkillLanguageEl              TextTranslationSkillLanguage = "el"
	TextTranslationSkillLanguageEn              TextTranslationSkillLanguage = "en"
	TextTranslationSkillLanguageEs              TextTranslationSkillLanguage = "es"
	TextTranslationSkillLanguageEt              TextTranslationSkillLanguage = "et"
	TextTranslationSkillLanguageFa              TextTranslationSkillLanguage = "fa"
	TextTranslationSkillLanguageFi              TextTranslationSkillLanguage = "fi"
	TextTranslationSkillLanguageFil             TextTranslationSkillLanguage = "fil"
	TextTranslationSkillLanguageFj              TextTranslationSkillLanguage = "fj"
	TextTranslationSkillLanguageFr              TextTranslationSkillLanguage = "fr"
	TextTranslationSkillLanguageGa              TextTranslationSkillLanguage = "ga"
	TextTranslationSkillLanguageHe              TextTranslationSkillLanguage = "he"
	TextTranslationSkillLanguageHi              TextTranslationSkillLanguage = "hi"
	TextTranslationSkillLanguageHr              TextTranslationSkillLanguage = "hr"
	TextTranslationSkillLanguageHt              TextTranslationSkillLanguage = "ht"
	TextTranslationSkillLanguageHu              TextTranslationSkillLanguage = "hu"
	TextTranslationSkillLanguageId              TextTranslationSkillLanguage = "id"
	TextTranslationSkillLanguageIs              TextTranslationSkillLanguage = "is"
	TextTranslationSkillLanguageIt              TextTranslationSkillLanguage = "it"
	TextTranslationSkillLanguageJa              TextTranslationSkillLanguage = "ja"
	TextTranslationSkillLanguageKn              TextTranslationSkillLanguage = "kn"
	TextTranslationSkillLanguageKo              TextTranslationSkillLanguage = "ko"
	TextTranslationSkillLanguageLt              TextTranslationSkillLanguage = "lt"
	TextTranslationSkillLanguageLv              TextTranslationSkillLanguage = "lv"
	TextTranslationSkillLanguageMg              TextTranslationSkillLanguage = "mg"
	TextTranslationSkillLanguageMi              TextTranslationSkillLanguage = "mi"
	TextTranslationSkillLanguageMl              TextTranslationSkillLanguage = "ml"
	TextTranslationSkillLanguageMs              TextTranslationSkillLanguage = "ms"
	TextTranslationSkillLanguageMt              TextTranslationSkillLanguage = "mt"
	TextTranslationSkillLanguageMww             TextTranslationSkillLanguage = "mww"
	TextTranslationSkillLanguageNb              TextTranslationSkillLanguage = "nb"
	TextTranslationSkillLanguageNl              TextTranslationSkillLanguage = "nl"
	TextTranslationSkillLanguageOtq             TextTranslationSkillLanguage = "otq"
	TextTranslationSkillLanguagePa              TextTranslationSkillLanguage = "pa"
	TextTranslationSkillLanguagePl              TextTranslationSkillLanguage = "pl"
	TextTranslationSkillLanguagePt              TextTranslationSkillLanguage = "pt"
	TextTranslationSkillLanguagePtNegativePT    TextTranslationSkillLanguage = "pt-PT"
	TextTranslationSkillLanguagePtNegativebr    TextTranslationSkillLanguage = "pt-br"
	TextTranslationSkillLanguageRo              TextTranslationSkillLanguage = "ro"
	TextTranslationSkillLanguageRu              TextTranslationSkillLanguage = "ru"
	TextTranslationSkillLanguageSk              TextTranslationSkillLanguage = "sk"
	TextTranslationSkillLanguageSl              TextTranslationSkillLanguage = "sl"
	TextTranslationSkillLanguageSm              TextTranslationSkillLanguage = "sm"
	TextTranslationSkillLanguageSrNegativeCyrl  TextTranslationSkillLanguage = "sr-Cyrl"
	TextTranslationSkillLanguageSrNegativeLatn  TextTranslationSkillLanguage = "sr-Latn"
	TextTranslationSkillLanguageSv              TextTranslationSkillLanguage = "sv"
	TextTranslationSkillLanguageSw              TextTranslationSkillLanguage = "sw"
	TextTranslationSkillLanguageTa              TextTranslationSkillLanguage = "ta"
	TextTranslationSkillLanguageTe              TextTranslationSkillLanguage = "te"
	TextTranslationSkillLanguageTh              TextTranslationSkillLanguage = "th"
	TextTranslationSkillLanguageTlh             TextTranslationSkillLanguage = "tlh"
	TextTranslationSkillLanguageTlhNegativeLatn TextTranslationSkillLanguage = "tlh-Latn"
	TextTranslationSkillLanguageTlhNegativePiqd TextTranslationSkillLanguage = "tlh-Piqd"
	TextTranslationSkillLanguageTo              TextTranslationSkillLanguage = "to"
	TextTranslationSkillLanguageTr              TextTranslationSkillLanguage = "tr"
	TextTranslationSkillLanguageTy              TextTranslationSkillLanguage = "ty"
	TextTranslationSkillLanguageUk              TextTranslationSkillLanguage = "uk"
	TextTranslationSkillLanguageUr              TextTranslationSkillLanguage = "ur"
	TextTranslationSkillLanguageVi              TextTranslationSkillLanguage = "vi"
	TextTranslationSkillLanguageYua             TextTranslationSkillLanguage = "yua"
	TextTranslationSkillLanguageYue             TextTranslationSkillLanguage = "yue"
	TextTranslationSkillLanguageZhNegativeHans  TextTranslationSkillLanguage = "zh-Hans"
	TextTranslationSkillLanguageZhNegativeHant  TextTranslationSkillLanguage = "zh-Hant"
)

func PossibleValuesForTextTranslationSkillLanguage() []string {
	return []string{
		string(TextTranslationSkillLanguageAf),
		string(TextTranslationSkillLanguageAr),
		string(TextTranslationSkillLanguageBg),
		string(TextTranslationSkillLanguageBn),
		string(TextTranslationSkillLanguageBs),
		string(TextTranslationSkillLanguageCa),
		string(TextTranslationSkillLanguageCs),
		string(TextTranslationSkillLanguageCy),
		string(TextTranslationSkillLanguageDa),
		string(TextTranslationSkillLanguageDe),
		string(TextTranslationSkillLanguageEl),
		string(TextTranslationSkillLanguageEn),
		string(TextTranslationSkillLanguageEs),
		string(TextTranslationSkillLanguageEt),
		string(TextTranslationSkillLanguageFa),
		string(TextTranslationSkillLanguageFi),
		string(TextTranslationSkillLanguageFil),
		string(TextTranslationSkillLanguageFj),
		string(TextTranslationSkillLanguageFr),
		string(TextTranslationSkillLanguageGa),
		string(TextTranslationSkillLanguageHe),
		string(TextTranslationSkillLanguageHi),
		string(TextTranslationSkillLanguageHr),
		string(TextTranslationSkillLanguageHt),
		string(TextTranslationSkillLanguageHu),
		string(TextTranslationSkillLanguageId),
		string(TextTranslationSkillLanguageIs),
		string(TextTranslationSkillLanguageIt),
		string(TextTranslationSkillLanguageJa),
		string(TextTranslationSkillLanguageKn),
		string(TextTranslationSkillLanguageKo),
		string(TextTranslationSkillLanguageLt),
		string(TextTranslationSkillLanguageLv),
		string(TextTranslationSkillLanguageMg),
		string(TextTranslationSkillLanguageMi),
		string(TextTranslationSkillLanguageMl),
		string(TextTranslationSkillLanguageMs),
		string(TextTranslationSkillLanguageMt),
		string(TextTranslationSkillLanguageMww),
		string(TextTranslationSkillLanguageNb),
		string(TextTranslationSkillLanguageNl),
		string(TextTranslationSkillLanguageOtq),
		string(TextTranslationSkillLanguagePa),
		string(TextTranslationSkillLanguagePl),
		string(TextTranslationSkillLanguagePt),
		string(TextTranslationSkillLanguagePtNegativePT),
		string(TextTranslationSkillLanguagePtNegativebr),
		string(TextTranslationSkillLanguageRo),
		string(TextTranslationSkillLanguageRu),
		string(TextTranslationSkillLanguageSk),
		string(TextTranslationSkillLanguageSl),
		string(TextTranslationSkillLanguageSm),
		string(TextTranslationSkillLanguageSrNegativeCyrl),
		string(TextTranslationSkillLanguageSrNegativeLatn),
		string(TextTranslationSkillLanguageSv),
		string(TextTranslationSkillLanguageSw),
		string(TextTranslationSkillLanguageTa),
		string(TextTranslationSkillLanguageTe),
		string(TextTranslationSkillLanguageTh),
		string(TextTranslationSkillLanguageTlh),
		string(TextTranslationSkillLanguageTlhNegativeLatn),
		string(TextTranslationSkillLanguageTlhNegativePiqd),
		string(TextTranslationSkillLanguageTo),
		string(TextTranslationSkillLanguageTr),
		string(TextTranslationSkillLanguageTy),
		string(TextTranslationSkillLanguageUk),
		string(TextTranslationSkillLanguageUr),
		string(TextTranslationSkillLanguageVi),
		string(TextTranslationSkillLanguageYua),
		string(TextTranslationSkillLanguageYue),
		string(TextTranslationSkillLanguageZhNegativeHans),
		string(TextTranslationSkillLanguageZhNegativeHant),
	}
}

func (s *TextTranslationSkillLanguage) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseTextTranslationSkillLanguage(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseTextTranslationSkillLanguage(input string) (*TextTranslationSkillLanguage, error) {
	vals := map[string]TextTranslationSkillLanguage{
		"af":       TextTranslationSkillLanguageAf,
		"ar":       TextTranslationSkillLanguageAr,
		"bg":       TextTranslationSkillLanguageBg,
		"bn":       TextTranslationSkillLanguageBn,
		"bs":       TextTranslationSkillLanguageBs,
		"ca":       TextTranslationSkillLanguageCa,
		"cs":       TextTranslationSkillLanguageCs,
		"cy":       TextTranslationSkillLanguageCy,
		"da":       TextTranslationSkillLanguageDa,
		"de":       TextTranslationSkillLanguageDe,
		"el":       TextTranslationSkillLanguageEl,
		"en":       TextTranslationSkillLanguageEn,
		"es":       TextTranslationSkillLanguageEs,
		"et":       TextTranslationSkillLanguageEt,
		"fa":       TextTranslationSkillLanguageFa,
		"fi":       TextTranslationSkillLanguageFi,
		"fil":      TextTranslationSkillLanguageFil,
		"fj":       TextTranslationSkillLanguageFj,
		"fr":       TextTranslationSkillLanguageFr,
		"ga":       TextTranslationSkillLanguageGa,
		"he":       TextTranslationSkillLanguageHe,
		"hi":       TextTranslationSkillLanguageHi,
		"hr":       TextTranslationSkillLanguageHr,
		"ht":       TextTranslationSkillLanguageHt,
		"hu":       TextTranslationSkillLanguageHu,
		"id":       TextTranslationSkillLanguageId,
		"is":       TextTranslationSkillLanguageIs,
		"it":       TextTranslationSkillLanguageIt,
		"ja":       TextTranslationSkillLanguageJa,
		"kn":       TextTranslationSkillLanguageKn,
		"ko":       TextTranslationSkillLanguageKo,
		"lt":       TextTranslationSkillLanguageLt,
		"lv":       TextTranslationSkillLanguageLv,
		"mg":       TextTranslationSkillLanguageMg,
		"mi":       TextTranslationSkillLanguageMi,
		"ml":       TextTranslationSkillLanguageMl,
		"ms":       TextTranslationSkillLanguageMs,
		"mt":       TextTranslationSkillLanguageMt,
		"mww":      TextTranslationSkillLanguageMww,
		"nb":       TextTranslationSkillLanguageNb,
		"nl":       TextTranslationSkillLanguageNl,
		"otq":      TextTranslationSkillLanguageOtq,
		"pa":       TextTranslationSkillLanguagePa,
		"pl":       TextTranslationSkillLanguagePl,
		"pt":       TextTranslationSkillLanguagePt,
		"pt-pt":    TextTranslationSkillLanguagePtNegativePT,
		"pt-br":    TextTranslationSkillLanguagePtNegativebr,
		"ro":       TextTranslationSkillLanguageRo,
		"ru":       TextTranslationSkillLanguageRu,
		"sk":       TextTranslationSkillLanguageSk,
		"sl":       TextTranslationSkillLanguageSl,
		"sm":       TextTranslationSkillLanguageSm,
		"sr-cyrl":  TextTranslationSkillLanguageSrNegativeCyrl,
		"sr-latn":  TextTranslationSkillLanguageSrNegativeLatn,
		"sv":       TextTranslationSkillLanguageSv,
		"sw":       TextTranslationSkillLanguageSw,
		"ta":       TextTranslationSkillLanguageTa,
		"te":       TextTranslationSkillLanguageTe,
		"th":       TextTranslationSkillLanguageTh,
		"tlh":      TextTranslationSkillLanguageTlh,
		"tlh-latn": TextTranslationSkillLanguageTlhNegativeLatn,
		"tlh-piqd": TextTranslationSkillLanguageTlhNegativePiqd,
		"to":       TextTranslationSkillLanguageTo,
		"tr":       TextTranslationSkillLanguageTr,
		"ty":       TextTranslationSkillLanguageTy,
		"uk":       TextTranslationSkillLanguageUk,
		"ur":       TextTranslationSkillLanguageUr,
		"vi":       TextTranslationSkillLanguageVi,
		"yua":      TextTranslationSkillLanguageYua,
		"yue":      TextTranslationSkillLanguageYue,
		"zh-hans":  TextTranslationSkillLanguageZhNegativeHans,
		"zh-hant":  TextTranslationSkillLanguageZhNegativeHant,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := TextTranslationSkillLanguage(input)
	return &out, nil
}

type VisualFeature string

const (
	VisualFeatureAdult       VisualFeature = "adult"
	VisualFeatureBrands      VisualFeature = "brands"
	VisualFeatureCategories  VisualFeature = "categories"
	VisualFeatureDescription VisualFeature = "description"
	VisualFeatureFaces       VisualFeature = "faces"
	VisualFeatureObjects     VisualFeature = "objects"
	VisualFeatureTags        VisualFeature = "tags"
)

func PossibleValuesForVisualFeature() []string {
	return []string{
		string(VisualFeatureAdult),
		string(VisualFeatureBrands),
		string(VisualFeatureCategories),
		string(VisualFeatureDescription),
		string(VisualFeatureFaces),
		string(VisualFeatureObjects),
		string(VisualFeatureTags),
	}
}

func (s *VisualFeature) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseVisualFeature(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseVisualFeature(input string) (*VisualFeature, error) {
	vals := map[string]VisualFeature{
		"adult":       VisualFeatureAdult,
		"brands":      VisualFeatureBrands,
		"categories":  VisualFeatureCategories,
		"description": VisualFeatureDescription,
		"faces":       VisualFeatureFaces,
		"objects":     VisualFeatureObjects,
		"tags":        VisualFeatureTags,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := VisualFeature(input)
	return &out, nil
}
