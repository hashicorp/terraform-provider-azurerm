package indexes

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

type CharFilterName string

const (
	CharFilterNameHtmlStrip CharFilterName = "html_strip"
)

func PossibleValuesForCharFilterName() []string {
	return []string{
		string(CharFilterNameHtmlStrip),
	}
}

func (s *CharFilterName) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseCharFilterName(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseCharFilterName(input string) (*CharFilterName, error) {
	vals := map[string]CharFilterName{
		"html_strip": CharFilterNameHtmlStrip,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := CharFilterName(input)
	return &out, nil
}

type CjkBigramTokenFilterScripts string

const (
	CjkBigramTokenFilterScriptsHan      CjkBigramTokenFilterScripts = "han"
	CjkBigramTokenFilterScriptsHangul   CjkBigramTokenFilterScripts = "hangul"
	CjkBigramTokenFilterScriptsHiragana CjkBigramTokenFilterScripts = "hiragana"
	CjkBigramTokenFilterScriptsKatakana CjkBigramTokenFilterScripts = "katakana"
)

func PossibleValuesForCjkBigramTokenFilterScripts() []string {
	return []string{
		string(CjkBigramTokenFilterScriptsHan),
		string(CjkBigramTokenFilterScriptsHangul),
		string(CjkBigramTokenFilterScriptsHiragana),
		string(CjkBigramTokenFilterScriptsKatakana),
	}
}

func (s *CjkBigramTokenFilterScripts) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseCjkBigramTokenFilterScripts(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseCjkBigramTokenFilterScripts(input string) (*CjkBigramTokenFilterScripts, error) {
	vals := map[string]CjkBigramTokenFilterScripts{
		"han":      CjkBigramTokenFilterScriptsHan,
		"hangul":   CjkBigramTokenFilterScriptsHangul,
		"hiragana": CjkBigramTokenFilterScriptsHiragana,
		"katakana": CjkBigramTokenFilterScriptsKatakana,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := CjkBigramTokenFilterScripts(input)
	return &out, nil
}

type EdgeNGramTokenFilterSide string

const (
	EdgeNGramTokenFilterSideBack  EdgeNGramTokenFilterSide = "back"
	EdgeNGramTokenFilterSideFront EdgeNGramTokenFilterSide = "front"
)

func PossibleValuesForEdgeNGramTokenFilterSide() []string {
	return []string{
		string(EdgeNGramTokenFilterSideBack),
		string(EdgeNGramTokenFilterSideFront),
	}
}

func (s *EdgeNGramTokenFilterSide) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseEdgeNGramTokenFilterSide(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseEdgeNGramTokenFilterSide(input string) (*EdgeNGramTokenFilterSide, error) {
	vals := map[string]EdgeNGramTokenFilterSide{
		"back":  EdgeNGramTokenFilterSideBack,
		"front": EdgeNGramTokenFilterSideFront,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := EdgeNGramTokenFilterSide(input)
	return &out, nil
}

type LexicalAnalyzerName string

const (
	LexicalAnalyzerNameArPointlucene                    LexicalAnalyzerName = "ar.lucene"
	LexicalAnalyzerNameArPointmicrosoft                 LexicalAnalyzerName = "ar.microsoft"
	LexicalAnalyzerNameBgPointlucene                    LexicalAnalyzerName = "bg.lucene"
	LexicalAnalyzerNameBgPointmicrosoft                 LexicalAnalyzerName = "bg.microsoft"
	LexicalAnalyzerNameBnPointmicrosoft                 LexicalAnalyzerName = "bn.microsoft"
	LexicalAnalyzerNameCaPointlucene                    LexicalAnalyzerName = "ca.lucene"
	LexicalAnalyzerNameCaPointmicrosoft                 LexicalAnalyzerName = "ca.microsoft"
	LexicalAnalyzerNameCsPointlucene                    LexicalAnalyzerName = "cs.lucene"
	LexicalAnalyzerNameCsPointmicrosoft                 LexicalAnalyzerName = "cs.microsoft"
	LexicalAnalyzerNameDaPointlucene                    LexicalAnalyzerName = "da.lucene"
	LexicalAnalyzerNameDaPointmicrosoft                 LexicalAnalyzerName = "da.microsoft"
	LexicalAnalyzerNameDePointlucene                    LexicalAnalyzerName = "de.lucene"
	LexicalAnalyzerNameDePointmicrosoft                 LexicalAnalyzerName = "de.microsoft"
	LexicalAnalyzerNameElPointlucene                    LexicalAnalyzerName = "el.lucene"
	LexicalAnalyzerNameElPointmicrosoft                 LexicalAnalyzerName = "el.microsoft"
	LexicalAnalyzerNameEnPointlucene                    LexicalAnalyzerName = "en.lucene"
	LexicalAnalyzerNameEnPointmicrosoft                 LexicalAnalyzerName = "en.microsoft"
	LexicalAnalyzerNameEsPointlucene                    LexicalAnalyzerName = "es.lucene"
	LexicalAnalyzerNameEsPointmicrosoft                 LexicalAnalyzerName = "es.microsoft"
	LexicalAnalyzerNameEtPointmicrosoft                 LexicalAnalyzerName = "et.microsoft"
	LexicalAnalyzerNameEuPointlucene                    LexicalAnalyzerName = "eu.lucene"
	LexicalAnalyzerNameFaPointlucene                    LexicalAnalyzerName = "fa.lucene"
	LexicalAnalyzerNameFiPointlucene                    LexicalAnalyzerName = "fi.lucene"
	LexicalAnalyzerNameFiPointmicrosoft                 LexicalAnalyzerName = "fi.microsoft"
	LexicalAnalyzerNameFrPointlucene                    LexicalAnalyzerName = "fr.lucene"
	LexicalAnalyzerNameFrPointmicrosoft                 LexicalAnalyzerName = "fr.microsoft"
	LexicalAnalyzerNameGaPointlucene                    LexicalAnalyzerName = "ga.lucene"
	LexicalAnalyzerNameGlPointlucene                    LexicalAnalyzerName = "gl.lucene"
	LexicalAnalyzerNameGuPointmicrosoft                 LexicalAnalyzerName = "gu.microsoft"
	LexicalAnalyzerNameHePointmicrosoft                 LexicalAnalyzerName = "he.microsoft"
	LexicalAnalyzerNameHiPointlucene                    LexicalAnalyzerName = "hi.lucene"
	LexicalAnalyzerNameHiPointmicrosoft                 LexicalAnalyzerName = "hi.microsoft"
	LexicalAnalyzerNameHrPointmicrosoft                 LexicalAnalyzerName = "hr.microsoft"
	LexicalAnalyzerNameHuPointlucene                    LexicalAnalyzerName = "hu.lucene"
	LexicalAnalyzerNameHuPointmicrosoft                 LexicalAnalyzerName = "hu.microsoft"
	LexicalAnalyzerNameHyPointlucene                    LexicalAnalyzerName = "hy.lucene"
	LexicalAnalyzerNameIdPointlucene                    LexicalAnalyzerName = "id.lucene"
	LexicalAnalyzerNameIdPointmicrosoft                 LexicalAnalyzerName = "id.microsoft"
	LexicalAnalyzerNameIsPointmicrosoft                 LexicalAnalyzerName = "is.microsoft"
	LexicalAnalyzerNameItPointlucene                    LexicalAnalyzerName = "it.lucene"
	LexicalAnalyzerNameItPointmicrosoft                 LexicalAnalyzerName = "it.microsoft"
	LexicalAnalyzerNameJaPointlucene                    LexicalAnalyzerName = "ja.lucene"
	LexicalAnalyzerNameJaPointmicrosoft                 LexicalAnalyzerName = "ja.microsoft"
	LexicalAnalyzerNameKeyword                          LexicalAnalyzerName = "keyword"
	LexicalAnalyzerNameKnPointmicrosoft                 LexicalAnalyzerName = "kn.microsoft"
	LexicalAnalyzerNameKoPointlucene                    LexicalAnalyzerName = "ko.lucene"
	LexicalAnalyzerNameKoPointmicrosoft                 LexicalAnalyzerName = "ko.microsoft"
	LexicalAnalyzerNameLtPointmicrosoft                 LexicalAnalyzerName = "lt.microsoft"
	LexicalAnalyzerNameLvPointlucene                    LexicalAnalyzerName = "lv.lucene"
	LexicalAnalyzerNameLvPointmicrosoft                 LexicalAnalyzerName = "lv.microsoft"
	LexicalAnalyzerNameMlPointmicrosoft                 LexicalAnalyzerName = "ml.microsoft"
	LexicalAnalyzerNameMrPointmicrosoft                 LexicalAnalyzerName = "mr.microsoft"
	LexicalAnalyzerNameMsPointmicrosoft                 LexicalAnalyzerName = "ms.microsoft"
	LexicalAnalyzerNameNbPointmicrosoft                 LexicalAnalyzerName = "nb.microsoft"
	LexicalAnalyzerNameNlPointlucene                    LexicalAnalyzerName = "nl.lucene"
	LexicalAnalyzerNameNlPointmicrosoft                 LexicalAnalyzerName = "nl.microsoft"
	LexicalAnalyzerNameNoPointlucene                    LexicalAnalyzerName = "no.lucene"
	LexicalAnalyzerNamePaPointmicrosoft                 LexicalAnalyzerName = "pa.microsoft"
	LexicalAnalyzerNamePattern                          LexicalAnalyzerName = "pattern"
	LexicalAnalyzerNamePlPointlucene                    LexicalAnalyzerName = "pl.lucene"
	LexicalAnalyzerNamePlPointmicrosoft                 LexicalAnalyzerName = "pl.microsoft"
	LexicalAnalyzerNamePtNegativeBRPointlucene          LexicalAnalyzerName = "pt-BR.lucene"
	LexicalAnalyzerNamePtNegativeBRPointmicrosoft       LexicalAnalyzerName = "pt-BR.microsoft"
	LexicalAnalyzerNamePtNegativePTPointlucene          LexicalAnalyzerName = "pt-PT.lucene"
	LexicalAnalyzerNamePtNegativePTPointmicrosoft       LexicalAnalyzerName = "pt-PT.microsoft"
	LexicalAnalyzerNameRoPointlucene                    LexicalAnalyzerName = "ro.lucene"
	LexicalAnalyzerNameRoPointmicrosoft                 LexicalAnalyzerName = "ro.microsoft"
	LexicalAnalyzerNameRuPointlucene                    LexicalAnalyzerName = "ru.lucene"
	LexicalAnalyzerNameRuPointmicrosoft                 LexicalAnalyzerName = "ru.microsoft"
	LexicalAnalyzerNameSimple                           LexicalAnalyzerName = "simple"
	LexicalAnalyzerNameSkPointmicrosoft                 LexicalAnalyzerName = "sk.microsoft"
	LexicalAnalyzerNameSlPointmicrosoft                 LexicalAnalyzerName = "sl.microsoft"
	LexicalAnalyzerNameSrNegativecyrillicPointmicrosoft LexicalAnalyzerName = "sr-cyrillic.microsoft"
	LexicalAnalyzerNameSrNegativelatinPointmicrosoft    LexicalAnalyzerName = "sr-latin.microsoft"
	LexicalAnalyzerNameStandardPointlucene              LexicalAnalyzerName = "standard.lucene"
	LexicalAnalyzerNameStandardasciifoldingPointlucene  LexicalAnalyzerName = "standardasciifolding.lucene"
	LexicalAnalyzerNameStop                             LexicalAnalyzerName = "stop"
	LexicalAnalyzerNameSvPointlucene                    LexicalAnalyzerName = "sv.lucene"
	LexicalAnalyzerNameSvPointmicrosoft                 LexicalAnalyzerName = "sv.microsoft"
	LexicalAnalyzerNameTaPointmicrosoft                 LexicalAnalyzerName = "ta.microsoft"
	LexicalAnalyzerNameTePointmicrosoft                 LexicalAnalyzerName = "te.microsoft"
	LexicalAnalyzerNameThPointlucene                    LexicalAnalyzerName = "th.lucene"
	LexicalAnalyzerNameThPointmicrosoft                 LexicalAnalyzerName = "th.microsoft"
	LexicalAnalyzerNameTrPointlucene                    LexicalAnalyzerName = "tr.lucene"
	LexicalAnalyzerNameTrPointmicrosoft                 LexicalAnalyzerName = "tr.microsoft"
	LexicalAnalyzerNameUkPointmicrosoft                 LexicalAnalyzerName = "uk.microsoft"
	LexicalAnalyzerNameUrPointmicrosoft                 LexicalAnalyzerName = "ur.microsoft"
	LexicalAnalyzerNameViPointmicrosoft                 LexicalAnalyzerName = "vi.microsoft"
	LexicalAnalyzerNameWhitespace                       LexicalAnalyzerName = "whitespace"
	LexicalAnalyzerNameZhNegativeHansPointlucene        LexicalAnalyzerName = "zh-Hans.lucene"
	LexicalAnalyzerNameZhNegativeHansPointmicrosoft     LexicalAnalyzerName = "zh-Hans.microsoft"
	LexicalAnalyzerNameZhNegativeHantPointlucene        LexicalAnalyzerName = "zh-Hant.lucene"
	LexicalAnalyzerNameZhNegativeHantPointmicrosoft     LexicalAnalyzerName = "zh-Hant.microsoft"
)

func PossibleValuesForLexicalAnalyzerName() []string {
	return []string{
		string(LexicalAnalyzerNameArPointlucene),
		string(LexicalAnalyzerNameArPointmicrosoft),
		string(LexicalAnalyzerNameBgPointlucene),
		string(LexicalAnalyzerNameBgPointmicrosoft),
		string(LexicalAnalyzerNameBnPointmicrosoft),
		string(LexicalAnalyzerNameCaPointlucene),
		string(LexicalAnalyzerNameCaPointmicrosoft),
		string(LexicalAnalyzerNameCsPointlucene),
		string(LexicalAnalyzerNameCsPointmicrosoft),
		string(LexicalAnalyzerNameDaPointlucene),
		string(LexicalAnalyzerNameDaPointmicrosoft),
		string(LexicalAnalyzerNameDePointlucene),
		string(LexicalAnalyzerNameDePointmicrosoft),
		string(LexicalAnalyzerNameElPointlucene),
		string(LexicalAnalyzerNameElPointmicrosoft),
		string(LexicalAnalyzerNameEnPointlucene),
		string(LexicalAnalyzerNameEnPointmicrosoft),
		string(LexicalAnalyzerNameEsPointlucene),
		string(LexicalAnalyzerNameEsPointmicrosoft),
		string(LexicalAnalyzerNameEtPointmicrosoft),
		string(LexicalAnalyzerNameEuPointlucene),
		string(LexicalAnalyzerNameFaPointlucene),
		string(LexicalAnalyzerNameFiPointlucene),
		string(LexicalAnalyzerNameFiPointmicrosoft),
		string(LexicalAnalyzerNameFrPointlucene),
		string(LexicalAnalyzerNameFrPointmicrosoft),
		string(LexicalAnalyzerNameGaPointlucene),
		string(LexicalAnalyzerNameGlPointlucene),
		string(LexicalAnalyzerNameGuPointmicrosoft),
		string(LexicalAnalyzerNameHePointmicrosoft),
		string(LexicalAnalyzerNameHiPointlucene),
		string(LexicalAnalyzerNameHiPointmicrosoft),
		string(LexicalAnalyzerNameHrPointmicrosoft),
		string(LexicalAnalyzerNameHuPointlucene),
		string(LexicalAnalyzerNameHuPointmicrosoft),
		string(LexicalAnalyzerNameHyPointlucene),
		string(LexicalAnalyzerNameIdPointlucene),
		string(LexicalAnalyzerNameIdPointmicrosoft),
		string(LexicalAnalyzerNameIsPointmicrosoft),
		string(LexicalAnalyzerNameItPointlucene),
		string(LexicalAnalyzerNameItPointmicrosoft),
		string(LexicalAnalyzerNameJaPointlucene),
		string(LexicalAnalyzerNameJaPointmicrosoft),
		string(LexicalAnalyzerNameKeyword),
		string(LexicalAnalyzerNameKnPointmicrosoft),
		string(LexicalAnalyzerNameKoPointlucene),
		string(LexicalAnalyzerNameKoPointmicrosoft),
		string(LexicalAnalyzerNameLtPointmicrosoft),
		string(LexicalAnalyzerNameLvPointlucene),
		string(LexicalAnalyzerNameLvPointmicrosoft),
		string(LexicalAnalyzerNameMlPointmicrosoft),
		string(LexicalAnalyzerNameMrPointmicrosoft),
		string(LexicalAnalyzerNameMsPointmicrosoft),
		string(LexicalAnalyzerNameNbPointmicrosoft),
		string(LexicalAnalyzerNameNlPointlucene),
		string(LexicalAnalyzerNameNlPointmicrosoft),
		string(LexicalAnalyzerNameNoPointlucene),
		string(LexicalAnalyzerNamePaPointmicrosoft),
		string(LexicalAnalyzerNamePattern),
		string(LexicalAnalyzerNamePlPointlucene),
		string(LexicalAnalyzerNamePlPointmicrosoft),
		string(LexicalAnalyzerNamePtNegativeBRPointlucene),
		string(LexicalAnalyzerNamePtNegativeBRPointmicrosoft),
		string(LexicalAnalyzerNamePtNegativePTPointlucene),
		string(LexicalAnalyzerNamePtNegativePTPointmicrosoft),
		string(LexicalAnalyzerNameRoPointlucene),
		string(LexicalAnalyzerNameRoPointmicrosoft),
		string(LexicalAnalyzerNameRuPointlucene),
		string(LexicalAnalyzerNameRuPointmicrosoft),
		string(LexicalAnalyzerNameSimple),
		string(LexicalAnalyzerNameSkPointmicrosoft),
		string(LexicalAnalyzerNameSlPointmicrosoft),
		string(LexicalAnalyzerNameSrNegativecyrillicPointmicrosoft),
		string(LexicalAnalyzerNameSrNegativelatinPointmicrosoft),
		string(LexicalAnalyzerNameStandardPointlucene),
		string(LexicalAnalyzerNameStandardasciifoldingPointlucene),
		string(LexicalAnalyzerNameStop),
		string(LexicalAnalyzerNameSvPointlucene),
		string(LexicalAnalyzerNameSvPointmicrosoft),
		string(LexicalAnalyzerNameTaPointmicrosoft),
		string(LexicalAnalyzerNameTePointmicrosoft),
		string(LexicalAnalyzerNameThPointlucene),
		string(LexicalAnalyzerNameThPointmicrosoft),
		string(LexicalAnalyzerNameTrPointlucene),
		string(LexicalAnalyzerNameTrPointmicrosoft),
		string(LexicalAnalyzerNameUkPointmicrosoft),
		string(LexicalAnalyzerNameUrPointmicrosoft),
		string(LexicalAnalyzerNameViPointmicrosoft),
		string(LexicalAnalyzerNameWhitespace),
		string(LexicalAnalyzerNameZhNegativeHansPointlucene),
		string(LexicalAnalyzerNameZhNegativeHansPointmicrosoft),
		string(LexicalAnalyzerNameZhNegativeHantPointlucene),
		string(LexicalAnalyzerNameZhNegativeHantPointmicrosoft),
	}
}

func (s *LexicalAnalyzerName) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseLexicalAnalyzerName(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseLexicalAnalyzerName(input string) (*LexicalAnalyzerName, error) {
	vals := map[string]LexicalAnalyzerName{
		"ar.lucene":                   LexicalAnalyzerNameArPointlucene,
		"ar.microsoft":                LexicalAnalyzerNameArPointmicrosoft,
		"bg.lucene":                   LexicalAnalyzerNameBgPointlucene,
		"bg.microsoft":                LexicalAnalyzerNameBgPointmicrosoft,
		"bn.microsoft":                LexicalAnalyzerNameBnPointmicrosoft,
		"ca.lucene":                   LexicalAnalyzerNameCaPointlucene,
		"ca.microsoft":                LexicalAnalyzerNameCaPointmicrosoft,
		"cs.lucene":                   LexicalAnalyzerNameCsPointlucene,
		"cs.microsoft":                LexicalAnalyzerNameCsPointmicrosoft,
		"da.lucene":                   LexicalAnalyzerNameDaPointlucene,
		"da.microsoft":                LexicalAnalyzerNameDaPointmicrosoft,
		"de.lucene":                   LexicalAnalyzerNameDePointlucene,
		"de.microsoft":                LexicalAnalyzerNameDePointmicrosoft,
		"el.lucene":                   LexicalAnalyzerNameElPointlucene,
		"el.microsoft":                LexicalAnalyzerNameElPointmicrosoft,
		"en.lucene":                   LexicalAnalyzerNameEnPointlucene,
		"en.microsoft":                LexicalAnalyzerNameEnPointmicrosoft,
		"es.lucene":                   LexicalAnalyzerNameEsPointlucene,
		"es.microsoft":                LexicalAnalyzerNameEsPointmicrosoft,
		"et.microsoft":                LexicalAnalyzerNameEtPointmicrosoft,
		"eu.lucene":                   LexicalAnalyzerNameEuPointlucene,
		"fa.lucene":                   LexicalAnalyzerNameFaPointlucene,
		"fi.lucene":                   LexicalAnalyzerNameFiPointlucene,
		"fi.microsoft":                LexicalAnalyzerNameFiPointmicrosoft,
		"fr.lucene":                   LexicalAnalyzerNameFrPointlucene,
		"fr.microsoft":                LexicalAnalyzerNameFrPointmicrosoft,
		"ga.lucene":                   LexicalAnalyzerNameGaPointlucene,
		"gl.lucene":                   LexicalAnalyzerNameGlPointlucene,
		"gu.microsoft":                LexicalAnalyzerNameGuPointmicrosoft,
		"he.microsoft":                LexicalAnalyzerNameHePointmicrosoft,
		"hi.lucene":                   LexicalAnalyzerNameHiPointlucene,
		"hi.microsoft":                LexicalAnalyzerNameHiPointmicrosoft,
		"hr.microsoft":                LexicalAnalyzerNameHrPointmicrosoft,
		"hu.lucene":                   LexicalAnalyzerNameHuPointlucene,
		"hu.microsoft":                LexicalAnalyzerNameHuPointmicrosoft,
		"hy.lucene":                   LexicalAnalyzerNameHyPointlucene,
		"id.lucene":                   LexicalAnalyzerNameIdPointlucene,
		"id.microsoft":                LexicalAnalyzerNameIdPointmicrosoft,
		"is.microsoft":                LexicalAnalyzerNameIsPointmicrosoft,
		"it.lucene":                   LexicalAnalyzerNameItPointlucene,
		"it.microsoft":                LexicalAnalyzerNameItPointmicrosoft,
		"ja.lucene":                   LexicalAnalyzerNameJaPointlucene,
		"ja.microsoft":                LexicalAnalyzerNameJaPointmicrosoft,
		"keyword":                     LexicalAnalyzerNameKeyword,
		"kn.microsoft":                LexicalAnalyzerNameKnPointmicrosoft,
		"ko.lucene":                   LexicalAnalyzerNameKoPointlucene,
		"ko.microsoft":                LexicalAnalyzerNameKoPointmicrosoft,
		"lt.microsoft":                LexicalAnalyzerNameLtPointmicrosoft,
		"lv.lucene":                   LexicalAnalyzerNameLvPointlucene,
		"lv.microsoft":                LexicalAnalyzerNameLvPointmicrosoft,
		"ml.microsoft":                LexicalAnalyzerNameMlPointmicrosoft,
		"mr.microsoft":                LexicalAnalyzerNameMrPointmicrosoft,
		"ms.microsoft":                LexicalAnalyzerNameMsPointmicrosoft,
		"nb.microsoft":                LexicalAnalyzerNameNbPointmicrosoft,
		"nl.lucene":                   LexicalAnalyzerNameNlPointlucene,
		"nl.microsoft":                LexicalAnalyzerNameNlPointmicrosoft,
		"no.lucene":                   LexicalAnalyzerNameNoPointlucene,
		"pa.microsoft":                LexicalAnalyzerNamePaPointmicrosoft,
		"pattern":                     LexicalAnalyzerNamePattern,
		"pl.lucene":                   LexicalAnalyzerNamePlPointlucene,
		"pl.microsoft":                LexicalAnalyzerNamePlPointmicrosoft,
		"pt-br.lucene":                LexicalAnalyzerNamePtNegativeBRPointlucene,
		"pt-br.microsoft":             LexicalAnalyzerNamePtNegativeBRPointmicrosoft,
		"pt-pt.lucene":                LexicalAnalyzerNamePtNegativePTPointlucene,
		"pt-pt.microsoft":             LexicalAnalyzerNamePtNegativePTPointmicrosoft,
		"ro.lucene":                   LexicalAnalyzerNameRoPointlucene,
		"ro.microsoft":                LexicalAnalyzerNameRoPointmicrosoft,
		"ru.lucene":                   LexicalAnalyzerNameRuPointlucene,
		"ru.microsoft":                LexicalAnalyzerNameRuPointmicrosoft,
		"simple":                      LexicalAnalyzerNameSimple,
		"sk.microsoft":                LexicalAnalyzerNameSkPointmicrosoft,
		"sl.microsoft":                LexicalAnalyzerNameSlPointmicrosoft,
		"sr-cyrillic.microsoft":       LexicalAnalyzerNameSrNegativecyrillicPointmicrosoft,
		"sr-latin.microsoft":          LexicalAnalyzerNameSrNegativelatinPointmicrosoft,
		"standard.lucene":             LexicalAnalyzerNameStandardPointlucene,
		"standardasciifolding.lucene": LexicalAnalyzerNameStandardasciifoldingPointlucene,
		"stop":                        LexicalAnalyzerNameStop,
		"sv.lucene":                   LexicalAnalyzerNameSvPointlucene,
		"sv.microsoft":                LexicalAnalyzerNameSvPointmicrosoft,
		"ta.microsoft":                LexicalAnalyzerNameTaPointmicrosoft,
		"te.microsoft":                LexicalAnalyzerNameTePointmicrosoft,
		"th.lucene":                   LexicalAnalyzerNameThPointlucene,
		"th.microsoft":                LexicalAnalyzerNameThPointmicrosoft,
		"tr.lucene":                   LexicalAnalyzerNameTrPointlucene,
		"tr.microsoft":                LexicalAnalyzerNameTrPointmicrosoft,
		"uk.microsoft":                LexicalAnalyzerNameUkPointmicrosoft,
		"ur.microsoft":                LexicalAnalyzerNameUrPointmicrosoft,
		"vi.microsoft":                LexicalAnalyzerNameViPointmicrosoft,
		"whitespace":                  LexicalAnalyzerNameWhitespace,
		"zh-hans.lucene":              LexicalAnalyzerNameZhNegativeHansPointlucene,
		"zh-hans.microsoft":           LexicalAnalyzerNameZhNegativeHansPointmicrosoft,
		"zh-hant.lucene":              LexicalAnalyzerNameZhNegativeHantPointlucene,
		"zh-hant.microsoft":           LexicalAnalyzerNameZhNegativeHantPointmicrosoft,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := LexicalAnalyzerName(input)
	return &out, nil
}

type LexicalNormalizerName string

const (
	LexicalNormalizerNameAsciifolding LexicalNormalizerName = "asciifolding"
	LexicalNormalizerNameElision      LexicalNormalizerName = "elision"
	LexicalNormalizerNameLowercase    LexicalNormalizerName = "lowercase"
	LexicalNormalizerNameStandard     LexicalNormalizerName = "standard"
	LexicalNormalizerNameUppercase    LexicalNormalizerName = "uppercase"
)

func PossibleValuesForLexicalNormalizerName() []string {
	return []string{
		string(LexicalNormalizerNameAsciifolding),
		string(LexicalNormalizerNameElision),
		string(LexicalNormalizerNameLowercase),
		string(LexicalNormalizerNameStandard),
		string(LexicalNormalizerNameUppercase),
	}
}

func (s *LexicalNormalizerName) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseLexicalNormalizerName(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseLexicalNormalizerName(input string) (*LexicalNormalizerName, error) {
	vals := map[string]LexicalNormalizerName{
		"asciifolding": LexicalNormalizerNameAsciifolding,
		"elision":      LexicalNormalizerNameElision,
		"lowercase":    LexicalNormalizerNameLowercase,
		"standard":     LexicalNormalizerNameStandard,
		"uppercase":    LexicalNormalizerNameUppercase,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := LexicalNormalizerName(input)
	return &out, nil
}

type LexicalTokenizerName string

const (
	LexicalTokenizerNameClassic                            LexicalTokenizerName = "classic"
	LexicalTokenizerNameEdgeNGram                          LexicalTokenizerName = "edgeNGram"
	LexicalTokenizerNameKeywordVTwo                        LexicalTokenizerName = "keyword_v2"
	LexicalTokenizerNameLetter                             LexicalTokenizerName = "letter"
	LexicalTokenizerNameLowercase                          LexicalTokenizerName = "lowercase"
	LexicalTokenizerNameMicrosoftLanguageStemmingTokenizer LexicalTokenizerName = "microsoft_language_stemming_tokenizer"
	LexicalTokenizerNameMicrosoftLanguageTokenizer         LexicalTokenizerName = "microsoft_language_tokenizer"
	LexicalTokenizerNameNGram                              LexicalTokenizerName = "nGram"
	LexicalTokenizerNamePathHierarchyVTwo                  LexicalTokenizerName = "path_hierarchy_v2"
	LexicalTokenizerNamePattern                            LexicalTokenizerName = "pattern"
	LexicalTokenizerNameStandardVTwo                       LexicalTokenizerName = "standard_v2"
	LexicalTokenizerNameUaxUrlEmail                        LexicalTokenizerName = "uax_url_email"
	LexicalTokenizerNameWhitespace                         LexicalTokenizerName = "whitespace"
)

func PossibleValuesForLexicalTokenizerName() []string {
	return []string{
		string(LexicalTokenizerNameClassic),
		string(LexicalTokenizerNameEdgeNGram),
		string(LexicalTokenizerNameKeywordVTwo),
		string(LexicalTokenizerNameLetter),
		string(LexicalTokenizerNameLowercase),
		string(LexicalTokenizerNameMicrosoftLanguageStemmingTokenizer),
		string(LexicalTokenizerNameMicrosoftLanguageTokenizer),
		string(LexicalTokenizerNameNGram),
		string(LexicalTokenizerNamePathHierarchyVTwo),
		string(LexicalTokenizerNamePattern),
		string(LexicalTokenizerNameStandardVTwo),
		string(LexicalTokenizerNameUaxUrlEmail),
		string(LexicalTokenizerNameWhitespace),
	}
}

func (s *LexicalTokenizerName) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseLexicalTokenizerName(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseLexicalTokenizerName(input string) (*LexicalTokenizerName, error) {
	vals := map[string]LexicalTokenizerName{
		"classic":                               LexicalTokenizerNameClassic,
		"edgengram":                             LexicalTokenizerNameEdgeNGram,
		"keyword_v2":                            LexicalTokenizerNameKeywordVTwo,
		"letter":                                LexicalTokenizerNameLetter,
		"lowercase":                             LexicalTokenizerNameLowercase,
		"microsoft_language_stemming_tokenizer": LexicalTokenizerNameMicrosoftLanguageStemmingTokenizer,
		"microsoft_language_tokenizer":          LexicalTokenizerNameMicrosoftLanguageTokenizer,
		"ngram":                                 LexicalTokenizerNameNGram,
		"path_hierarchy_v2":                     LexicalTokenizerNamePathHierarchyVTwo,
		"pattern":                               LexicalTokenizerNamePattern,
		"standard_v2":                           LexicalTokenizerNameStandardVTwo,
		"uax_url_email":                         LexicalTokenizerNameUaxUrlEmail,
		"whitespace":                            LexicalTokenizerNameWhitespace,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := LexicalTokenizerName(input)
	return &out, nil
}

type MicrosoftStemmingTokenizerLanguage string

const (
	MicrosoftStemmingTokenizerLanguageArabic              MicrosoftStemmingTokenizerLanguage = "arabic"
	MicrosoftStemmingTokenizerLanguageBangla              MicrosoftStemmingTokenizerLanguage = "bangla"
	MicrosoftStemmingTokenizerLanguageBulgarian           MicrosoftStemmingTokenizerLanguage = "bulgarian"
	MicrosoftStemmingTokenizerLanguageCatalan             MicrosoftStemmingTokenizerLanguage = "catalan"
	MicrosoftStemmingTokenizerLanguageCroatian            MicrosoftStemmingTokenizerLanguage = "croatian"
	MicrosoftStemmingTokenizerLanguageCzech               MicrosoftStemmingTokenizerLanguage = "czech"
	MicrosoftStemmingTokenizerLanguageDanish              MicrosoftStemmingTokenizerLanguage = "danish"
	MicrosoftStemmingTokenizerLanguageDutch               MicrosoftStemmingTokenizerLanguage = "dutch"
	MicrosoftStemmingTokenizerLanguageEnglish             MicrosoftStemmingTokenizerLanguage = "english"
	MicrosoftStemmingTokenizerLanguageEstonian            MicrosoftStemmingTokenizerLanguage = "estonian"
	MicrosoftStemmingTokenizerLanguageFinnish             MicrosoftStemmingTokenizerLanguage = "finnish"
	MicrosoftStemmingTokenizerLanguageFrench              MicrosoftStemmingTokenizerLanguage = "french"
	MicrosoftStemmingTokenizerLanguageGerman              MicrosoftStemmingTokenizerLanguage = "german"
	MicrosoftStemmingTokenizerLanguageGreek               MicrosoftStemmingTokenizerLanguage = "greek"
	MicrosoftStemmingTokenizerLanguageGujarati            MicrosoftStemmingTokenizerLanguage = "gujarati"
	MicrosoftStemmingTokenizerLanguageHebrew              MicrosoftStemmingTokenizerLanguage = "hebrew"
	MicrosoftStemmingTokenizerLanguageHindi               MicrosoftStemmingTokenizerLanguage = "hindi"
	MicrosoftStemmingTokenizerLanguageHungarian           MicrosoftStemmingTokenizerLanguage = "hungarian"
	MicrosoftStemmingTokenizerLanguageIcelandic           MicrosoftStemmingTokenizerLanguage = "icelandic"
	MicrosoftStemmingTokenizerLanguageIndonesian          MicrosoftStemmingTokenizerLanguage = "indonesian"
	MicrosoftStemmingTokenizerLanguageItalian             MicrosoftStemmingTokenizerLanguage = "italian"
	MicrosoftStemmingTokenizerLanguageKannada             MicrosoftStemmingTokenizerLanguage = "kannada"
	MicrosoftStemmingTokenizerLanguageLatvian             MicrosoftStemmingTokenizerLanguage = "latvian"
	MicrosoftStemmingTokenizerLanguageLithuanian          MicrosoftStemmingTokenizerLanguage = "lithuanian"
	MicrosoftStemmingTokenizerLanguageMalay               MicrosoftStemmingTokenizerLanguage = "malay"
	MicrosoftStemmingTokenizerLanguageMalayalam           MicrosoftStemmingTokenizerLanguage = "malayalam"
	MicrosoftStemmingTokenizerLanguageMarathi             MicrosoftStemmingTokenizerLanguage = "marathi"
	MicrosoftStemmingTokenizerLanguageNorwegianBokmaal    MicrosoftStemmingTokenizerLanguage = "norwegianBokmaal"
	MicrosoftStemmingTokenizerLanguagePolish              MicrosoftStemmingTokenizerLanguage = "polish"
	MicrosoftStemmingTokenizerLanguagePortuguese          MicrosoftStemmingTokenizerLanguage = "portuguese"
	MicrosoftStemmingTokenizerLanguagePortugueseBrazilian MicrosoftStemmingTokenizerLanguage = "portugueseBrazilian"
	MicrosoftStemmingTokenizerLanguagePunjabi             MicrosoftStemmingTokenizerLanguage = "punjabi"
	MicrosoftStemmingTokenizerLanguageRomanian            MicrosoftStemmingTokenizerLanguage = "romanian"
	MicrosoftStemmingTokenizerLanguageRussian             MicrosoftStemmingTokenizerLanguage = "russian"
	MicrosoftStemmingTokenizerLanguageSerbianCyrillic     MicrosoftStemmingTokenizerLanguage = "serbianCyrillic"
	MicrosoftStemmingTokenizerLanguageSerbianLatin        MicrosoftStemmingTokenizerLanguage = "serbianLatin"
	MicrosoftStemmingTokenizerLanguageSlovak              MicrosoftStemmingTokenizerLanguage = "slovak"
	MicrosoftStemmingTokenizerLanguageSlovenian           MicrosoftStemmingTokenizerLanguage = "slovenian"
	MicrosoftStemmingTokenizerLanguageSpanish             MicrosoftStemmingTokenizerLanguage = "spanish"
	MicrosoftStemmingTokenizerLanguageSwedish             MicrosoftStemmingTokenizerLanguage = "swedish"
	MicrosoftStemmingTokenizerLanguageTamil               MicrosoftStemmingTokenizerLanguage = "tamil"
	MicrosoftStemmingTokenizerLanguageTelugu              MicrosoftStemmingTokenizerLanguage = "telugu"
	MicrosoftStemmingTokenizerLanguageTurkish             MicrosoftStemmingTokenizerLanguage = "turkish"
	MicrosoftStemmingTokenizerLanguageUkrainian           MicrosoftStemmingTokenizerLanguage = "ukrainian"
	MicrosoftStemmingTokenizerLanguageUrdu                MicrosoftStemmingTokenizerLanguage = "urdu"
)

func PossibleValuesForMicrosoftStemmingTokenizerLanguage() []string {
	return []string{
		string(MicrosoftStemmingTokenizerLanguageArabic),
		string(MicrosoftStemmingTokenizerLanguageBangla),
		string(MicrosoftStemmingTokenizerLanguageBulgarian),
		string(MicrosoftStemmingTokenizerLanguageCatalan),
		string(MicrosoftStemmingTokenizerLanguageCroatian),
		string(MicrosoftStemmingTokenizerLanguageCzech),
		string(MicrosoftStemmingTokenizerLanguageDanish),
		string(MicrosoftStemmingTokenizerLanguageDutch),
		string(MicrosoftStemmingTokenizerLanguageEnglish),
		string(MicrosoftStemmingTokenizerLanguageEstonian),
		string(MicrosoftStemmingTokenizerLanguageFinnish),
		string(MicrosoftStemmingTokenizerLanguageFrench),
		string(MicrosoftStemmingTokenizerLanguageGerman),
		string(MicrosoftStemmingTokenizerLanguageGreek),
		string(MicrosoftStemmingTokenizerLanguageGujarati),
		string(MicrosoftStemmingTokenizerLanguageHebrew),
		string(MicrosoftStemmingTokenizerLanguageHindi),
		string(MicrosoftStemmingTokenizerLanguageHungarian),
		string(MicrosoftStemmingTokenizerLanguageIcelandic),
		string(MicrosoftStemmingTokenizerLanguageIndonesian),
		string(MicrosoftStemmingTokenizerLanguageItalian),
		string(MicrosoftStemmingTokenizerLanguageKannada),
		string(MicrosoftStemmingTokenizerLanguageLatvian),
		string(MicrosoftStemmingTokenizerLanguageLithuanian),
		string(MicrosoftStemmingTokenizerLanguageMalay),
		string(MicrosoftStemmingTokenizerLanguageMalayalam),
		string(MicrosoftStemmingTokenizerLanguageMarathi),
		string(MicrosoftStemmingTokenizerLanguageNorwegianBokmaal),
		string(MicrosoftStemmingTokenizerLanguagePolish),
		string(MicrosoftStemmingTokenizerLanguagePortuguese),
		string(MicrosoftStemmingTokenizerLanguagePortugueseBrazilian),
		string(MicrosoftStemmingTokenizerLanguagePunjabi),
		string(MicrosoftStemmingTokenizerLanguageRomanian),
		string(MicrosoftStemmingTokenizerLanguageRussian),
		string(MicrosoftStemmingTokenizerLanguageSerbianCyrillic),
		string(MicrosoftStemmingTokenizerLanguageSerbianLatin),
		string(MicrosoftStemmingTokenizerLanguageSlovak),
		string(MicrosoftStemmingTokenizerLanguageSlovenian),
		string(MicrosoftStemmingTokenizerLanguageSpanish),
		string(MicrosoftStemmingTokenizerLanguageSwedish),
		string(MicrosoftStemmingTokenizerLanguageTamil),
		string(MicrosoftStemmingTokenizerLanguageTelugu),
		string(MicrosoftStemmingTokenizerLanguageTurkish),
		string(MicrosoftStemmingTokenizerLanguageUkrainian),
		string(MicrosoftStemmingTokenizerLanguageUrdu),
	}
}

func (s *MicrosoftStemmingTokenizerLanguage) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseMicrosoftStemmingTokenizerLanguage(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseMicrosoftStemmingTokenizerLanguage(input string) (*MicrosoftStemmingTokenizerLanguage, error) {
	vals := map[string]MicrosoftStemmingTokenizerLanguage{
		"arabic":              MicrosoftStemmingTokenizerLanguageArabic,
		"bangla":              MicrosoftStemmingTokenizerLanguageBangla,
		"bulgarian":           MicrosoftStemmingTokenizerLanguageBulgarian,
		"catalan":             MicrosoftStemmingTokenizerLanguageCatalan,
		"croatian":            MicrosoftStemmingTokenizerLanguageCroatian,
		"czech":               MicrosoftStemmingTokenizerLanguageCzech,
		"danish":              MicrosoftStemmingTokenizerLanguageDanish,
		"dutch":               MicrosoftStemmingTokenizerLanguageDutch,
		"english":             MicrosoftStemmingTokenizerLanguageEnglish,
		"estonian":            MicrosoftStemmingTokenizerLanguageEstonian,
		"finnish":             MicrosoftStemmingTokenizerLanguageFinnish,
		"french":              MicrosoftStemmingTokenizerLanguageFrench,
		"german":              MicrosoftStemmingTokenizerLanguageGerman,
		"greek":               MicrosoftStemmingTokenizerLanguageGreek,
		"gujarati":            MicrosoftStemmingTokenizerLanguageGujarati,
		"hebrew":              MicrosoftStemmingTokenizerLanguageHebrew,
		"hindi":               MicrosoftStemmingTokenizerLanguageHindi,
		"hungarian":           MicrosoftStemmingTokenizerLanguageHungarian,
		"icelandic":           MicrosoftStemmingTokenizerLanguageIcelandic,
		"indonesian":          MicrosoftStemmingTokenizerLanguageIndonesian,
		"italian":             MicrosoftStemmingTokenizerLanguageItalian,
		"kannada":             MicrosoftStemmingTokenizerLanguageKannada,
		"latvian":             MicrosoftStemmingTokenizerLanguageLatvian,
		"lithuanian":          MicrosoftStemmingTokenizerLanguageLithuanian,
		"malay":               MicrosoftStemmingTokenizerLanguageMalay,
		"malayalam":           MicrosoftStemmingTokenizerLanguageMalayalam,
		"marathi":             MicrosoftStemmingTokenizerLanguageMarathi,
		"norwegianbokmaal":    MicrosoftStemmingTokenizerLanguageNorwegianBokmaal,
		"polish":              MicrosoftStemmingTokenizerLanguagePolish,
		"portuguese":          MicrosoftStemmingTokenizerLanguagePortuguese,
		"portuguesebrazilian": MicrosoftStemmingTokenizerLanguagePortugueseBrazilian,
		"punjabi":             MicrosoftStemmingTokenizerLanguagePunjabi,
		"romanian":            MicrosoftStemmingTokenizerLanguageRomanian,
		"russian":             MicrosoftStemmingTokenizerLanguageRussian,
		"serbiancyrillic":     MicrosoftStemmingTokenizerLanguageSerbianCyrillic,
		"serbianlatin":        MicrosoftStemmingTokenizerLanguageSerbianLatin,
		"slovak":              MicrosoftStemmingTokenizerLanguageSlovak,
		"slovenian":           MicrosoftStemmingTokenizerLanguageSlovenian,
		"spanish":             MicrosoftStemmingTokenizerLanguageSpanish,
		"swedish":             MicrosoftStemmingTokenizerLanguageSwedish,
		"tamil":               MicrosoftStemmingTokenizerLanguageTamil,
		"telugu":              MicrosoftStemmingTokenizerLanguageTelugu,
		"turkish":             MicrosoftStemmingTokenizerLanguageTurkish,
		"ukrainian":           MicrosoftStemmingTokenizerLanguageUkrainian,
		"urdu":                MicrosoftStemmingTokenizerLanguageUrdu,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := MicrosoftStemmingTokenizerLanguage(input)
	return &out, nil
}

type MicrosoftTokenizerLanguage string

const (
	MicrosoftTokenizerLanguageBangla              MicrosoftTokenizerLanguage = "bangla"
	MicrosoftTokenizerLanguageBulgarian           MicrosoftTokenizerLanguage = "bulgarian"
	MicrosoftTokenizerLanguageCatalan             MicrosoftTokenizerLanguage = "catalan"
	MicrosoftTokenizerLanguageChineseSimplified   MicrosoftTokenizerLanguage = "chineseSimplified"
	MicrosoftTokenizerLanguageChineseTraditional  MicrosoftTokenizerLanguage = "chineseTraditional"
	MicrosoftTokenizerLanguageCroatian            MicrosoftTokenizerLanguage = "croatian"
	MicrosoftTokenizerLanguageCzech               MicrosoftTokenizerLanguage = "czech"
	MicrosoftTokenizerLanguageDanish              MicrosoftTokenizerLanguage = "danish"
	MicrosoftTokenizerLanguageDutch               MicrosoftTokenizerLanguage = "dutch"
	MicrosoftTokenizerLanguageEnglish             MicrosoftTokenizerLanguage = "english"
	MicrosoftTokenizerLanguageFrench              MicrosoftTokenizerLanguage = "french"
	MicrosoftTokenizerLanguageGerman              MicrosoftTokenizerLanguage = "german"
	MicrosoftTokenizerLanguageGreek               MicrosoftTokenizerLanguage = "greek"
	MicrosoftTokenizerLanguageGujarati            MicrosoftTokenizerLanguage = "gujarati"
	MicrosoftTokenizerLanguageHindi               MicrosoftTokenizerLanguage = "hindi"
	MicrosoftTokenizerLanguageIcelandic           MicrosoftTokenizerLanguage = "icelandic"
	MicrosoftTokenizerLanguageIndonesian          MicrosoftTokenizerLanguage = "indonesian"
	MicrosoftTokenizerLanguageItalian             MicrosoftTokenizerLanguage = "italian"
	MicrosoftTokenizerLanguageJapanese            MicrosoftTokenizerLanguage = "japanese"
	MicrosoftTokenizerLanguageKannada             MicrosoftTokenizerLanguage = "kannada"
	MicrosoftTokenizerLanguageKorean              MicrosoftTokenizerLanguage = "korean"
	MicrosoftTokenizerLanguageMalay               MicrosoftTokenizerLanguage = "malay"
	MicrosoftTokenizerLanguageMalayalam           MicrosoftTokenizerLanguage = "malayalam"
	MicrosoftTokenizerLanguageMarathi             MicrosoftTokenizerLanguage = "marathi"
	MicrosoftTokenizerLanguageNorwegianBokmaal    MicrosoftTokenizerLanguage = "norwegianBokmaal"
	MicrosoftTokenizerLanguagePolish              MicrosoftTokenizerLanguage = "polish"
	MicrosoftTokenizerLanguagePortuguese          MicrosoftTokenizerLanguage = "portuguese"
	MicrosoftTokenizerLanguagePortugueseBrazilian MicrosoftTokenizerLanguage = "portugueseBrazilian"
	MicrosoftTokenizerLanguagePunjabi             MicrosoftTokenizerLanguage = "punjabi"
	MicrosoftTokenizerLanguageRomanian            MicrosoftTokenizerLanguage = "romanian"
	MicrosoftTokenizerLanguageRussian             MicrosoftTokenizerLanguage = "russian"
	MicrosoftTokenizerLanguageSerbianCyrillic     MicrosoftTokenizerLanguage = "serbianCyrillic"
	MicrosoftTokenizerLanguageSerbianLatin        MicrosoftTokenizerLanguage = "serbianLatin"
	MicrosoftTokenizerLanguageSlovenian           MicrosoftTokenizerLanguage = "slovenian"
	MicrosoftTokenizerLanguageSpanish             MicrosoftTokenizerLanguage = "spanish"
	MicrosoftTokenizerLanguageSwedish             MicrosoftTokenizerLanguage = "swedish"
	MicrosoftTokenizerLanguageTamil               MicrosoftTokenizerLanguage = "tamil"
	MicrosoftTokenizerLanguageTelugu              MicrosoftTokenizerLanguage = "telugu"
	MicrosoftTokenizerLanguageThai                MicrosoftTokenizerLanguage = "thai"
	MicrosoftTokenizerLanguageUkrainian           MicrosoftTokenizerLanguage = "ukrainian"
	MicrosoftTokenizerLanguageUrdu                MicrosoftTokenizerLanguage = "urdu"
	MicrosoftTokenizerLanguageVietnamese          MicrosoftTokenizerLanguage = "vietnamese"
)

func PossibleValuesForMicrosoftTokenizerLanguage() []string {
	return []string{
		string(MicrosoftTokenizerLanguageBangla),
		string(MicrosoftTokenizerLanguageBulgarian),
		string(MicrosoftTokenizerLanguageCatalan),
		string(MicrosoftTokenizerLanguageChineseSimplified),
		string(MicrosoftTokenizerLanguageChineseTraditional),
		string(MicrosoftTokenizerLanguageCroatian),
		string(MicrosoftTokenizerLanguageCzech),
		string(MicrosoftTokenizerLanguageDanish),
		string(MicrosoftTokenizerLanguageDutch),
		string(MicrosoftTokenizerLanguageEnglish),
		string(MicrosoftTokenizerLanguageFrench),
		string(MicrosoftTokenizerLanguageGerman),
		string(MicrosoftTokenizerLanguageGreek),
		string(MicrosoftTokenizerLanguageGujarati),
		string(MicrosoftTokenizerLanguageHindi),
		string(MicrosoftTokenizerLanguageIcelandic),
		string(MicrosoftTokenizerLanguageIndonesian),
		string(MicrosoftTokenizerLanguageItalian),
		string(MicrosoftTokenizerLanguageJapanese),
		string(MicrosoftTokenizerLanguageKannada),
		string(MicrosoftTokenizerLanguageKorean),
		string(MicrosoftTokenizerLanguageMalay),
		string(MicrosoftTokenizerLanguageMalayalam),
		string(MicrosoftTokenizerLanguageMarathi),
		string(MicrosoftTokenizerLanguageNorwegianBokmaal),
		string(MicrosoftTokenizerLanguagePolish),
		string(MicrosoftTokenizerLanguagePortuguese),
		string(MicrosoftTokenizerLanguagePortugueseBrazilian),
		string(MicrosoftTokenizerLanguagePunjabi),
		string(MicrosoftTokenizerLanguageRomanian),
		string(MicrosoftTokenizerLanguageRussian),
		string(MicrosoftTokenizerLanguageSerbianCyrillic),
		string(MicrosoftTokenizerLanguageSerbianLatin),
		string(MicrosoftTokenizerLanguageSlovenian),
		string(MicrosoftTokenizerLanguageSpanish),
		string(MicrosoftTokenizerLanguageSwedish),
		string(MicrosoftTokenizerLanguageTamil),
		string(MicrosoftTokenizerLanguageTelugu),
		string(MicrosoftTokenizerLanguageThai),
		string(MicrosoftTokenizerLanguageUkrainian),
		string(MicrosoftTokenizerLanguageUrdu),
		string(MicrosoftTokenizerLanguageVietnamese),
	}
}

func (s *MicrosoftTokenizerLanguage) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseMicrosoftTokenizerLanguage(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseMicrosoftTokenizerLanguage(input string) (*MicrosoftTokenizerLanguage, error) {
	vals := map[string]MicrosoftTokenizerLanguage{
		"bangla":              MicrosoftTokenizerLanguageBangla,
		"bulgarian":           MicrosoftTokenizerLanguageBulgarian,
		"catalan":             MicrosoftTokenizerLanguageCatalan,
		"chinesesimplified":   MicrosoftTokenizerLanguageChineseSimplified,
		"chinesetraditional":  MicrosoftTokenizerLanguageChineseTraditional,
		"croatian":            MicrosoftTokenizerLanguageCroatian,
		"czech":               MicrosoftTokenizerLanguageCzech,
		"danish":              MicrosoftTokenizerLanguageDanish,
		"dutch":               MicrosoftTokenizerLanguageDutch,
		"english":             MicrosoftTokenizerLanguageEnglish,
		"french":              MicrosoftTokenizerLanguageFrench,
		"german":              MicrosoftTokenizerLanguageGerman,
		"greek":               MicrosoftTokenizerLanguageGreek,
		"gujarati":            MicrosoftTokenizerLanguageGujarati,
		"hindi":               MicrosoftTokenizerLanguageHindi,
		"icelandic":           MicrosoftTokenizerLanguageIcelandic,
		"indonesian":          MicrosoftTokenizerLanguageIndonesian,
		"italian":             MicrosoftTokenizerLanguageItalian,
		"japanese":            MicrosoftTokenizerLanguageJapanese,
		"kannada":             MicrosoftTokenizerLanguageKannada,
		"korean":              MicrosoftTokenizerLanguageKorean,
		"malay":               MicrosoftTokenizerLanguageMalay,
		"malayalam":           MicrosoftTokenizerLanguageMalayalam,
		"marathi":             MicrosoftTokenizerLanguageMarathi,
		"norwegianbokmaal":    MicrosoftTokenizerLanguageNorwegianBokmaal,
		"polish":              MicrosoftTokenizerLanguagePolish,
		"portuguese":          MicrosoftTokenizerLanguagePortuguese,
		"portuguesebrazilian": MicrosoftTokenizerLanguagePortugueseBrazilian,
		"punjabi":             MicrosoftTokenizerLanguagePunjabi,
		"romanian":            MicrosoftTokenizerLanguageRomanian,
		"russian":             MicrosoftTokenizerLanguageRussian,
		"serbiancyrillic":     MicrosoftTokenizerLanguageSerbianCyrillic,
		"serbianlatin":        MicrosoftTokenizerLanguageSerbianLatin,
		"slovenian":           MicrosoftTokenizerLanguageSlovenian,
		"spanish":             MicrosoftTokenizerLanguageSpanish,
		"swedish":             MicrosoftTokenizerLanguageSwedish,
		"tamil":               MicrosoftTokenizerLanguageTamil,
		"telugu":              MicrosoftTokenizerLanguageTelugu,
		"thai":                MicrosoftTokenizerLanguageThai,
		"ukrainian":           MicrosoftTokenizerLanguageUkrainian,
		"urdu":                MicrosoftTokenizerLanguageUrdu,
		"vietnamese":          MicrosoftTokenizerLanguageVietnamese,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := MicrosoftTokenizerLanguage(input)
	return &out, nil
}

type PhoneticEncoder string

const (
	PhoneticEncoderBeiderMorse     PhoneticEncoder = "beiderMorse"
	PhoneticEncoderCaverphoneOne   PhoneticEncoder = "caverphone1"
	PhoneticEncoderCaverphoneTwo   PhoneticEncoder = "caverphone2"
	PhoneticEncoderCologne         PhoneticEncoder = "cologne"
	PhoneticEncoderDoubleMetaphone PhoneticEncoder = "doubleMetaphone"
	PhoneticEncoderHaasePhonetik   PhoneticEncoder = "haasePhonetik"
	PhoneticEncoderKoelnerPhonetik PhoneticEncoder = "koelnerPhonetik"
	PhoneticEncoderMetaphone       PhoneticEncoder = "metaphone"
	PhoneticEncoderNysiis          PhoneticEncoder = "nysiis"
	PhoneticEncoderRefinedSoundex  PhoneticEncoder = "refinedSoundex"
	PhoneticEncoderSoundex         PhoneticEncoder = "soundex"
)

func PossibleValuesForPhoneticEncoder() []string {
	return []string{
		string(PhoneticEncoderBeiderMorse),
		string(PhoneticEncoderCaverphoneOne),
		string(PhoneticEncoderCaverphoneTwo),
		string(PhoneticEncoderCologne),
		string(PhoneticEncoderDoubleMetaphone),
		string(PhoneticEncoderHaasePhonetik),
		string(PhoneticEncoderKoelnerPhonetik),
		string(PhoneticEncoderMetaphone),
		string(PhoneticEncoderNysiis),
		string(PhoneticEncoderRefinedSoundex),
		string(PhoneticEncoderSoundex),
	}
}

func (s *PhoneticEncoder) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parsePhoneticEncoder(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parsePhoneticEncoder(input string) (*PhoneticEncoder, error) {
	vals := map[string]PhoneticEncoder{
		"beidermorse":     PhoneticEncoderBeiderMorse,
		"caverphone1":     PhoneticEncoderCaverphoneOne,
		"caverphone2":     PhoneticEncoderCaverphoneTwo,
		"cologne":         PhoneticEncoderCologne,
		"doublemetaphone": PhoneticEncoderDoubleMetaphone,
		"haasephonetik":   PhoneticEncoderHaasePhonetik,
		"koelnerphonetik": PhoneticEncoderKoelnerPhonetik,
		"metaphone":       PhoneticEncoderMetaphone,
		"nysiis":          PhoneticEncoderNysiis,
		"refinedsoundex":  PhoneticEncoderRefinedSoundex,
		"soundex":         PhoneticEncoderSoundex,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := PhoneticEncoder(input)
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

type RankingOrder string

const (
	RankingOrderBoostedRerankerScore RankingOrder = "BoostedRerankerScore"
	RankingOrderRerankerScore        RankingOrder = "RerankerScore"
)

func PossibleValuesForRankingOrder() []string {
	return []string{
		string(RankingOrderBoostedRerankerScore),
		string(RankingOrderRerankerScore),
	}
}

func (s *RankingOrder) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseRankingOrder(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseRankingOrder(input string) (*RankingOrder, error) {
	vals := map[string]RankingOrder{
		"boostedrerankerscore": RankingOrderBoostedRerankerScore,
		"rerankerscore":        RankingOrderRerankerScore,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := RankingOrder(input)
	return &out, nil
}

type RegexFlags string

const (
	RegexFlagsCANONEQ         RegexFlags = "CANON_EQ"
	RegexFlagsCASEINSENSITIVE RegexFlags = "CASE_INSENSITIVE"
	RegexFlagsCOMMENTS        RegexFlags = "COMMENTS"
	RegexFlagsDOTALL          RegexFlags = "DOTALL"
	RegexFlagsLITERAL         RegexFlags = "LITERAL"
	RegexFlagsMULTILINE       RegexFlags = "MULTILINE"
	RegexFlagsUNICODECASE     RegexFlags = "UNICODE_CASE"
	RegexFlagsUNIXLINES       RegexFlags = "UNIX_LINES"
)

func PossibleValuesForRegexFlags() []string {
	return []string{
		string(RegexFlagsCANONEQ),
		string(RegexFlagsCASEINSENSITIVE),
		string(RegexFlagsCOMMENTS),
		string(RegexFlagsDOTALL),
		string(RegexFlagsLITERAL),
		string(RegexFlagsMULTILINE),
		string(RegexFlagsUNICODECASE),
		string(RegexFlagsUNIXLINES),
	}
}

func (s *RegexFlags) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseRegexFlags(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseRegexFlags(input string) (*RegexFlags, error) {
	vals := map[string]RegexFlags{
		"canon_eq":         RegexFlagsCANONEQ,
		"case_insensitive": RegexFlagsCASEINSENSITIVE,
		"comments":         RegexFlagsCOMMENTS,
		"dotall":           RegexFlagsDOTALL,
		"literal":          RegexFlagsLITERAL,
		"multiline":        RegexFlagsMULTILINE,
		"unicode_case":     RegexFlagsUNICODECASE,
		"unix_lines":       RegexFlagsUNIXLINES,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := RegexFlags(input)
	return &out, nil
}

type ScoringFunctionAggregation string

const (
	ScoringFunctionAggregationAverage       ScoringFunctionAggregation = "average"
	ScoringFunctionAggregationFirstMatching ScoringFunctionAggregation = "firstMatching"
	ScoringFunctionAggregationMaximum       ScoringFunctionAggregation = "maximum"
	ScoringFunctionAggregationMinimum       ScoringFunctionAggregation = "minimum"
	ScoringFunctionAggregationSum           ScoringFunctionAggregation = "sum"
)

func PossibleValuesForScoringFunctionAggregation() []string {
	return []string{
		string(ScoringFunctionAggregationAverage),
		string(ScoringFunctionAggregationFirstMatching),
		string(ScoringFunctionAggregationMaximum),
		string(ScoringFunctionAggregationMinimum),
		string(ScoringFunctionAggregationSum),
	}
}

func (s *ScoringFunctionAggregation) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseScoringFunctionAggregation(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseScoringFunctionAggregation(input string) (*ScoringFunctionAggregation, error) {
	vals := map[string]ScoringFunctionAggregation{
		"average":       ScoringFunctionAggregationAverage,
		"firstmatching": ScoringFunctionAggregationFirstMatching,
		"maximum":       ScoringFunctionAggregationMaximum,
		"minimum":       ScoringFunctionAggregationMinimum,
		"sum":           ScoringFunctionAggregationSum,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ScoringFunctionAggregation(input)
	return &out, nil
}

type ScoringFunctionInterpolation string

const (
	ScoringFunctionInterpolationConstant    ScoringFunctionInterpolation = "constant"
	ScoringFunctionInterpolationLinear      ScoringFunctionInterpolation = "linear"
	ScoringFunctionInterpolationLogarithmic ScoringFunctionInterpolation = "logarithmic"
	ScoringFunctionInterpolationQuadratic   ScoringFunctionInterpolation = "quadratic"
)

func PossibleValuesForScoringFunctionInterpolation() []string {
	return []string{
		string(ScoringFunctionInterpolationConstant),
		string(ScoringFunctionInterpolationLinear),
		string(ScoringFunctionInterpolationLogarithmic),
		string(ScoringFunctionInterpolationQuadratic),
	}
}

func (s *ScoringFunctionInterpolation) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseScoringFunctionInterpolation(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseScoringFunctionInterpolation(input string) (*ScoringFunctionInterpolation, error) {
	vals := map[string]ScoringFunctionInterpolation{
		"constant":    ScoringFunctionInterpolationConstant,
		"linear":      ScoringFunctionInterpolationLinear,
		"logarithmic": ScoringFunctionInterpolationLogarithmic,
		"quadratic":   ScoringFunctionInterpolationQuadratic,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ScoringFunctionInterpolation(input)
	return &out, nil
}

type SearchFieldDataType string

const (
	SearchFieldDataTypeEdmPointBoolean        SearchFieldDataType = "Edm.Boolean"
	SearchFieldDataTypeEdmPointByte           SearchFieldDataType = "Edm.Byte"
	SearchFieldDataTypeEdmPointComplexType    SearchFieldDataType = "Edm.ComplexType"
	SearchFieldDataTypeEdmPointDateTimeOffset SearchFieldDataType = "Edm.DateTimeOffset"
	SearchFieldDataTypeEdmPointDouble         SearchFieldDataType = "Edm.Double"
	SearchFieldDataTypeEdmPointGeographyPoint SearchFieldDataType = "Edm.GeographyPoint"
	SearchFieldDataTypeEdmPointHalf           SearchFieldDataType = "Edm.Half"
	SearchFieldDataTypeEdmPointIntOneSix      SearchFieldDataType = "Edm.Int16"
	SearchFieldDataTypeEdmPointIntSixFour     SearchFieldDataType = "Edm.Int64"
	SearchFieldDataTypeEdmPointIntThreeTwo    SearchFieldDataType = "Edm.Int32"
	SearchFieldDataTypeEdmPointSByte          SearchFieldDataType = "Edm.SByte"
	SearchFieldDataTypeEdmPointSingle         SearchFieldDataType = "Edm.Single"
	SearchFieldDataTypeEdmPointString         SearchFieldDataType = "Edm.String"
)

func PossibleValuesForSearchFieldDataType() []string {
	return []string{
		string(SearchFieldDataTypeEdmPointBoolean),
		string(SearchFieldDataTypeEdmPointByte),
		string(SearchFieldDataTypeEdmPointComplexType),
		string(SearchFieldDataTypeEdmPointDateTimeOffset),
		string(SearchFieldDataTypeEdmPointDouble),
		string(SearchFieldDataTypeEdmPointGeographyPoint),
		string(SearchFieldDataTypeEdmPointHalf),
		string(SearchFieldDataTypeEdmPointIntOneSix),
		string(SearchFieldDataTypeEdmPointIntSixFour),
		string(SearchFieldDataTypeEdmPointIntThreeTwo),
		string(SearchFieldDataTypeEdmPointSByte),
		string(SearchFieldDataTypeEdmPointSingle),
		string(SearchFieldDataTypeEdmPointString),
	}
}

func (s *SearchFieldDataType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseSearchFieldDataType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseSearchFieldDataType(input string) (*SearchFieldDataType, error) {
	vals := map[string]SearchFieldDataType{
		"edm.boolean":        SearchFieldDataTypeEdmPointBoolean,
		"edm.byte":           SearchFieldDataTypeEdmPointByte,
		"edm.complextype":    SearchFieldDataTypeEdmPointComplexType,
		"edm.datetimeoffset": SearchFieldDataTypeEdmPointDateTimeOffset,
		"edm.double":         SearchFieldDataTypeEdmPointDouble,
		"edm.geographypoint": SearchFieldDataTypeEdmPointGeographyPoint,
		"edm.half":           SearchFieldDataTypeEdmPointHalf,
		"edm.int16":          SearchFieldDataTypeEdmPointIntOneSix,
		"edm.int64":          SearchFieldDataTypeEdmPointIntSixFour,
		"edm.int32":          SearchFieldDataTypeEdmPointIntThreeTwo,
		"edm.sbyte":          SearchFieldDataTypeEdmPointSByte,
		"edm.single":         SearchFieldDataTypeEdmPointSingle,
		"edm.string":         SearchFieldDataTypeEdmPointString,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := SearchFieldDataType(input)
	return &out, nil
}

type SnowballTokenFilterLanguage string

const (
	SnowballTokenFilterLanguageArmenian   SnowballTokenFilterLanguage = "armenian"
	SnowballTokenFilterLanguageBasque     SnowballTokenFilterLanguage = "basque"
	SnowballTokenFilterLanguageCatalan    SnowballTokenFilterLanguage = "catalan"
	SnowballTokenFilterLanguageDanish     SnowballTokenFilterLanguage = "danish"
	SnowballTokenFilterLanguageDutch      SnowballTokenFilterLanguage = "dutch"
	SnowballTokenFilterLanguageEnglish    SnowballTokenFilterLanguage = "english"
	SnowballTokenFilterLanguageFinnish    SnowballTokenFilterLanguage = "finnish"
	SnowballTokenFilterLanguageFrench     SnowballTokenFilterLanguage = "french"
	SnowballTokenFilterLanguageGerman     SnowballTokenFilterLanguage = "german"
	SnowballTokenFilterLanguageGermanTwo  SnowballTokenFilterLanguage = "german2"
	SnowballTokenFilterLanguageHungarian  SnowballTokenFilterLanguage = "hungarian"
	SnowballTokenFilterLanguageItalian    SnowballTokenFilterLanguage = "italian"
	SnowballTokenFilterLanguageKp         SnowballTokenFilterLanguage = "kp"
	SnowballTokenFilterLanguageLovins     SnowballTokenFilterLanguage = "lovins"
	SnowballTokenFilterLanguageNorwegian  SnowballTokenFilterLanguage = "norwegian"
	SnowballTokenFilterLanguagePorter     SnowballTokenFilterLanguage = "porter"
	SnowballTokenFilterLanguagePortuguese SnowballTokenFilterLanguage = "portuguese"
	SnowballTokenFilterLanguageRomanian   SnowballTokenFilterLanguage = "romanian"
	SnowballTokenFilterLanguageRussian    SnowballTokenFilterLanguage = "russian"
	SnowballTokenFilterLanguageSpanish    SnowballTokenFilterLanguage = "spanish"
	SnowballTokenFilterLanguageSwedish    SnowballTokenFilterLanguage = "swedish"
	SnowballTokenFilterLanguageTurkish    SnowballTokenFilterLanguage = "turkish"
)

func PossibleValuesForSnowballTokenFilterLanguage() []string {
	return []string{
		string(SnowballTokenFilterLanguageArmenian),
		string(SnowballTokenFilterLanguageBasque),
		string(SnowballTokenFilterLanguageCatalan),
		string(SnowballTokenFilterLanguageDanish),
		string(SnowballTokenFilterLanguageDutch),
		string(SnowballTokenFilterLanguageEnglish),
		string(SnowballTokenFilterLanguageFinnish),
		string(SnowballTokenFilterLanguageFrench),
		string(SnowballTokenFilterLanguageGerman),
		string(SnowballTokenFilterLanguageGermanTwo),
		string(SnowballTokenFilterLanguageHungarian),
		string(SnowballTokenFilterLanguageItalian),
		string(SnowballTokenFilterLanguageKp),
		string(SnowballTokenFilterLanguageLovins),
		string(SnowballTokenFilterLanguageNorwegian),
		string(SnowballTokenFilterLanguagePorter),
		string(SnowballTokenFilterLanguagePortuguese),
		string(SnowballTokenFilterLanguageRomanian),
		string(SnowballTokenFilterLanguageRussian),
		string(SnowballTokenFilterLanguageSpanish),
		string(SnowballTokenFilterLanguageSwedish),
		string(SnowballTokenFilterLanguageTurkish),
	}
}

func (s *SnowballTokenFilterLanguage) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseSnowballTokenFilterLanguage(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseSnowballTokenFilterLanguage(input string) (*SnowballTokenFilterLanguage, error) {
	vals := map[string]SnowballTokenFilterLanguage{
		"armenian":   SnowballTokenFilterLanguageArmenian,
		"basque":     SnowballTokenFilterLanguageBasque,
		"catalan":    SnowballTokenFilterLanguageCatalan,
		"danish":     SnowballTokenFilterLanguageDanish,
		"dutch":      SnowballTokenFilterLanguageDutch,
		"english":    SnowballTokenFilterLanguageEnglish,
		"finnish":    SnowballTokenFilterLanguageFinnish,
		"french":     SnowballTokenFilterLanguageFrench,
		"german":     SnowballTokenFilterLanguageGerman,
		"german2":    SnowballTokenFilterLanguageGermanTwo,
		"hungarian":  SnowballTokenFilterLanguageHungarian,
		"italian":    SnowballTokenFilterLanguageItalian,
		"kp":         SnowballTokenFilterLanguageKp,
		"lovins":     SnowballTokenFilterLanguageLovins,
		"norwegian":  SnowballTokenFilterLanguageNorwegian,
		"porter":     SnowballTokenFilterLanguagePorter,
		"portuguese": SnowballTokenFilterLanguagePortuguese,
		"romanian":   SnowballTokenFilterLanguageRomanian,
		"russian":    SnowballTokenFilterLanguageRussian,
		"spanish":    SnowballTokenFilterLanguageSpanish,
		"swedish":    SnowballTokenFilterLanguageSwedish,
		"turkish":    SnowballTokenFilterLanguageTurkish,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := SnowballTokenFilterLanguage(input)
	return &out, nil
}

type StemmerTokenFilterLanguage string

const (
	StemmerTokenFilterLanguageArabic            StemmerTokenFilterLanguage = "arabic"
	StemmerTokenFilterLanguageArmenian          StemmerTokenFilterLanguage = "armenian"
	StemmerTokenFilterLanguageBasque            StemmerTokenFilterLanguage = "basque"
	StemmerTokenFilterLanguageBrazilian         StemmerTokenFilterLanguage = "brazilian"
	StemmerTokenFilterLanguageBulgarian         StemmerTokenFilterLanguage = "bulgarian"
	StemmerTokenFilterLanguageCatalan           StemmerTokenFilterLanguage = "catalan"
	StemmerTokenFilterLanguageCzech             StemmerTokenFilterLanguage = "czech"
	StemmerTokenFilterLanguageDanish            StemmerTokenFilterLanguage = "danish"
	StemmerTokenFilterLanguageDutch             StemmerTokenFilterLanguage = "dutch"
	StemmerTokenFilterLanguageDutchKp           StemmerTokenFilterLanguage = "dutchKp"
	StemmerTokenFilterLanguageEnglish           StemmerTokenFilterLanguage = "english"
	StemmerTokenFilterLanguageFinnish           StemmerTokenFilterLanguage = "finnish"
	StemmerTokenFilterLanguageFrench            StemmerTokenFilterLanguage = "french"
	StemmerTokenFilterLanguageGalician          StemmerTokenFilterLanguage = "galician"
	StemmerTokenFilterLanguageGerman            StemmerTokenFilterLanguage = "german"
	StemmerTokenFilterLanguageGermanTwo         StemmerTokenFilterLanguage = "german2"
	StemmerTokenFilterLanguageGreek             StemmerTokenFilterLanguage = "greek"
	StemmerTokenFilterLanguageHindi             StemmerTokenFilterLanguage = "hindi"
	StemmerTokenFilterLanguageHungarian         StemmerTokenFilterLanguage = "hungarian"
	StemmerTokenFilterLanguageIndonesian        StemmerTokenFilterLanguage = "indonesian"
	StemmerTokenFilterLanguageIrish             StemmerTokenFilterLanguage = "irish"
	StemmerTokenFilterLanguageItalian           StemmerTokenFilterLanguage = "italian"
	StemmerTokenFilterLanguageLatvian           StemmerTokenFilterLanguage = "latvian"
	StemmerTokenFilterLanguageLightEnglish      StemmerTokenFilterLanguage = "lightEnglish"
	StemmerTokenFilterLanguageLightFinnish      StemmerTokenFilterLanguage = "lightFinnish"
	StemmerTokenFilterLanguageLightFrench       StemmerTokenFilterLanguage = "lightFrench"
	StemmerTokenFilterLanguageLightGerman       StemmerTokenFilterLanguage = "lightGerman"
	StemmerTokenFilterLanguageLightHungarian    StemmerTokenFilterLanguage = "lightHungarian"
	StemmerTokenFilterLanguageLightItalian      StemmerTokenFilterLanguage = "lightItalian"
	StemmerTokenFilterLanguageLightNorwegian    StemmerTokenFilterLanguage = "lightNorwegian"
	StemmerTokenFilterLanguageLightNynorsk      StemmerTokenFilterLanguage = "lightNynorsk"
	StemmerTokenFilterLanguageLightPortuguese   StemmerTokenFilterLanguage = "lightPortuguese"
	StemmerTokenFilterLanguageLightRussian      StemmerTokenFilterLanguage = "lightRussian"
	StemmerTokenFilterLanguageLightSpanish      StemmerTokenFilterLanguage = "lightSpanish"
	StemmerTokenFilterLanguageLightSwedish      StemmerTokenFilterLanguage = "lightSwedish"
	StemmerTokenFilterLanguageLovins            StemmerTokenFilterLanguage = "lovins"
	StemmerTokenFilterLanguageMinimalEnglish    StemmerTokenFilterLanguage = "minimalEnglish"
	StemmerTokenFilterLanguageMinimalFrench     StemmerTokenFilterLanguage = "minimalFrench"
	StemmerTokenFilterLanguageMinimalGalician   StemmerTokenFilterLanguage = "minimalGalician"
	StemmerTokenFilterLanguageMinimalGerman     StemmerTokenFilterLanguage = "minimalGerman"
	StemmerTokenFilterLanguageMinimalNorwegian  StemmerTokenFilterLanguage = "minimalNorwegian"
	StemmerTokenFilterLanguageMinimalNynorsk    StemmerTokenFilterLanguage = "minimalNynorsk"
	StemmerTokenFilterLanguageMinimalPortuguese StemmerTokenFilterLanguage = "minimalPortuguese"
	StemmerTokenFilterLanguageNorwegian         StemmerTokenFilterLanguage = "norwegian"
	StemmerTokenFilterLanguagePorterTwo         StemmerTokenFilterLanguage = "porter2"
	StemmerTokenFilterLanguagePortuguese        StemmerTokenFilterLanguage = "portuguese"
	StemmerTokenFilterLanguagePortugueseRslp    StemmerTokenFilterLanguage = "portugueseRslp"
	StemmerTokenFilterLanguagePossessiveEnglish StemmerTokenFilterLanguage = "possessiveEnglish"
	StemmerTokenFilterLanguageRomanian          StemmerTokenFilterLanguage = "romanian"
	StemmerTokenFilterLanguageRussian           StemmerTokenFilterLanguage = "russian"
	StemmerTokenFilterLanguageSorani            StemmerTokenFilterLanguage = "sorani"
	StemmerTokenFilterLanguageSpanish           StemmerTokenFilterLanguage = "spanish"
	StemmerTokenFilterLanguageSwedish           StemmerTokenFilterLanguage = "swedish"
	StemmerTokenFilterLanguageTurkish           StemmerTokenFilterLanguage = "turkish"
)

func PossibleValuesForStemmerTokenFilterLanguage() []string {
	return []string{
		string(StemmerTokenFilterLanguageArabic),
		string(StemmerTokenFilterLanguageArmenian),
		string(StemmerTokenFilterLanguageBasque),
		string(StemmerTokenFilterLanguageBrazilian),
		string(StemmerTokenFilterLanguageBulgarian),
		string(StemmerTokenFilterLanguageCatalan),
		string(StemmerTokenFilterLanguageCzech),
		string(StemmerTokenFilterLanguageDanish),
		string(StemmerTokenFilterLanguageDutch),
		string(StemmerTokenFilterLanguageDutchKp),
		string(StemmerTokenFilterLanguageEnglish),
		string(StemmerTokenFilterLanguageFinnish),
		string(StemmerTokenFilterLanguageFrench),
		string(StemmerTokenFilterLanguageGalician),
		string(StemmerTokenFilterLanguageGerman),
		string(StemmerTokenFilterLanguageGermanTwo),
		string(StemmerTokenFilterLanguageGreek),
		string(StemmerTokenFilterLanguageHindi),
		string(StemmerTokenFilterLanguageHungarian),
		string(StemmerTokenFilterLanguageIndonesian),
		string(StemmerTokenFilterLanguageIrish),
		string(StemmerTokenFilterLanguageItalian),
		string(StemmerTokenFilterLanguageLatvian),
		string(StemmerTokenFilterLanguageLightEnglish),
		string(StemmerTokenFilterLanguageLightFinnish),
		string(StemmerTokenFilterLanguageLightFrench),
		string(StemmerTokenFilterLanguageLightGerman),
		string(StemmerTokenFilterLanguageLightHungarian),
		string(StemmerTokenFilterLanguageLightItalian),
		string(StemmerTokenFilterLanguageLightNorwegian),
		string(StemmerTokenFilterLanguageLightNynorsk),
		string(StemmerTokenFilterLanguageLightPortuguese),
		string(StemmerTokenFilterLanguageLightRussian),
		string(StemmerTokenFilterLanguageLightSpanish),
		string(StemmerTokenFilterLanguageLightSwedish),
		string(StemmerTokenFilterLanguageLovins),
		string(StemmerTokenFilterLanguageMinimalEnglish),
		string(StemmerTokenFilterLanguageMinimalFrench),
		string(StemmerTokenFilterLanguageMinimalGalician),
		string(StemmerTokenFilterLanguageMinimalGerman),
		string(StemmerTokenFilterLanguageMinimalNorwegian),
		string(StemmerTokenFilterLanguageMinimalNynorsk),
		string(StemmerTokenFilterLanguageMinimalPortuguese),
		string(StemmerTokenFilterLanguageNorwegian),
		string(StemmerTokenFilterLanguagePorterTwo),
		string(StemmerTokenFilterLanguagePortuguese),
		string(StemmerTokenFilterLanguagePortugueseRslp),
		string(StemmerTokenFilterLanguagePossessiveEnglish),
		string(StemmerTokenFilterLanguageRomanian),
		string(StemmerTokenFilterLanguageRussian),
		string(StemmerTokenFilterLanguageSorani),
		string(StemmerTokenFilterLanguageSpanish),
		string(StemmerTokenFilterLanguageSwedish),
		string(StemmerTokenFilterLanguageTurkish),
	}
}

func (s *StemmerTokenFilterLanguage) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseStemmerTokenFilterLanguage(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseStemmerTokenFilterLanguage(input string) (*StemmerTokenFilterLanguage, error) {
	vals := map[string]StemmerTokenFilterLanguage{
		"arabic":            StemmerTokenFilterLanguageArabic,
		"armenian":          StemmerTokenFilterLanguageArmenian,
		"basque":            StemmerTokenFilterLanguageBasque,
		"brazilian":         StemmerTokenFilterLanguageBrazilian,
		"bulgarian":         StemmerTokenFilterLanguageBulgarian,
		"catalan":           StemmerTokenFilterLanguageCatalan,
		"czech":             StemmerTokenFilterLanguageCzech,
		"danish":            StemmerTokenFilterLanguageDanish,
		"dutch":             StemmerTokenFilterLanguageDutch,
		"dutchkp":           StemmerTokenFilterLanguageDutchKp,
		"english":           StemmerTokenFilterLanguageEnglish,
		"finnish":           StemmerTokenFilterLanguageFinnish,
		"french":            StemmerTokenFilterLanguageFrench,
		"galician":          StemmerTokenFilterLanguageGalician,
		"german":            StemmerTokenFilterLanguageGerman,
		"german2":           StemmerTokenFilterLanguageGermanTwo,
		"greek":             StemmerTokenFilterLanguageGreek,
		"hindi":             StemmerTokenFilterLanguageHindi,
		"hungarian":         StemmerTokenFilterLanguageHungarian,
		"indonesian":        StemmerTokenFilterLanguageIndonesian,
		"irish":             StemmerTokenFilterLanguageIrish,
		"italian":           StemmerTokenFilterLanguageItalian,
		"latvian":           StemmerTokenFilterLanguageLatvian,
		"lightenglish":      StemmerTokenFilterLanguageLightEnglish,
		"lightfinnish":      StemmerTokenFilterLanguageLightFinnish,
		"lightfrench":       StemmerTokenFilterLanguageLightFrench,
		"lightgerman":       StemmerTokenFilterLanguageLightGerman,
		"lighthungarian":    StemmerTokenFilterLanguageLightHungarian,
		"lightitalian":      StemmerTokenFilterLanguageLightItalian,
		"lightnorwegian":    StemmerTokenFilterLanguageLightNorwegian,
		"lightnynorsk":      StemmerTokenFilterLanguageLightNynorsk,
		"lightportuguese":   StemmerTokenFilterLanguageLightPortuguese,
		"lightrussian":      StemmerTokenFilterLanguageLightRussian,
		"lightspanish":      StemmerTokenFilterLanguageLightSpanish,
		"lightswedish":      StemmerTokenFilterLanguageLightSwedish,
		"lovins":            StemmerTokenFilterLanguageLovins,
		"minimalenglish":    StemmerTokenFilterLanguageMinimalEnglish,
		"minimalfrench":     StemmerTokenFilterLanguageMinimalFrench,
		"minimalgalician":   StemmerTokenFilterLanguageMinimalGalician,
		"minimalgerman":     StemmerTokenFilterLanguageMinimalGerman,
		"minimalnorwegian":  StemmerTokenFilterLanguageMinimalNorwegian,
		"minimalnynorsk":    StemmerTokenFilterLanguageMinimalNynorsk,
		"minimalportuguese": StemmerTokenFilterLanguageMinimalPortuguese,
		"norwegian":         StemmerTokenFilterLanguageNorwegian,
		"porter2":           StemmerTokenFilterLanguagePorterTwo,
		"portuguese":        StemmerTokenFilterLanguagePortuguese,
		"portugueserslp":    StemmerTokenFilterLanguagePortugueseRslp,
		"possessiveenglish": StemmerTokenFilterLanguagePossessiveEnglish,
		"romanian":          StemmerTokenFilterLanguageRomanian,
		"russian":           StemmerTokenFilterLanguageRussian,
		"sorani":            StemmerTokenFilterLanguageSorani,
		"spanish":           StemmerTokenFilterLanguageSpanish,
		"swedish":           StemmerTokenFilterLanguageSwedish,
		"turkish":           StemmerTokenFilterLanguageTurkish,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := StemmerTokenFilterLanguage(input)
	return &out, nil
}

type StopwordsList string

const (
	StopwordsListArabic     StopwordsList = "arabic"
	StopwordsListArmenian   StopwordsList = "armenian"
	StopwordsListBasque     StopwordsList = "basque"
	StopwordsListBrazilian  StopwordsList = "brazilian"
	StopwordsListBulgarian  StopwordsList = "bulgarian"
	StopwordsListCatalan    StopwordsList = "catalan"
	StopwordsListCzech      StopwordsList = "czech"
	StopwordsListDanish     StopwordsList = "danish"
	StopwordsListDutch      StopwordsList = "dutch"
	StopwordsListEnglish    StopwordsList = "english"
	StopwordsListFinnish    StopwordsList = "finnish"
	StopwordsListFrench     StopwordsList = "french"
	StopwordsListGalician   StopwordsList = "galician"
	StopwordsListGerman     StopwordsList = "german"
	StopwordsListGreek      StopwordsList = "greek"
	StopwordsListHindi      StopwordsList = "hindi"
	StopwordsListHungarian  StopwordsList = "hungarian"
	StopwordsListIndonesian StopwordsList = "indonesian"
	StopwordsListIrish      StopwordsList = "irish"
	StopwordsListItalian    StopwordsList = "italian"
	StopwordsListLatvian    StopwordsList = "latvian"
	StopwordsListNorwegian  StopwordsList = "norwegian"
	StopwordsListPersian    StopwordsList = "persian"
	StopwordsListPortuguese StopwordsList = "portuguese"
	StopwordsListRomanian   StopwordsList = "romanian"
	StopwordsListRussian    StopwordsList = "russian"
	StopwordsListSorani     StopwordsList = "sorani"
	StopwordsListSpanish    StopwordsList = "spanish"
	StopwordsListSwedish    StopwordsList = "swedish"
	StopwordsListThai       StopwordsList = "thai"
	StopwordsListTurkish    StopwordsList = "turkish"
)

func PossibleValuesForStopwordsList() []string {
	return []string{
		string(StopwordsListArabic),
		string(StopwordsListArmenian),
		string(StopwordsListBasque),
		string(StopwordsListBrazilian),
		string(StopwordsListBulgarian),
		string(StopwordsListCatalan),
		string(StopwordsListCzech),
		string(StopwordsListDanish),
		string(StopwordsListDutch),
		string(StopwordsListEnglish),
		string(StopwordsListFinnish),
		string(StopwordsListFrench),
		string(StopwordsListGalician),
		string(StopwordsListGerman),
		string(StopwordsListGreek),
		string(StopwordsListHindi),
		string(StopwordsListHungarian),
		string(StopwordsListIndonesian),
		string(StopwordsListIrish),
		string(StopwordsListItalian),
		string(StopwordsListLatvian),
		string(StopwordsListNorwegian),
		string(StopwordsListPersian),
		string(StopwordsListPortuguese),
		string(StopwordsListRomanian),
		string(StopwordsListRussian),
		string(StopwordsListSorani),
		string(StopwordsListSpanish),
		string(StopwordsListSwedish),
		string(StopwordsListThai),
		string(StopwordsListTurkish),
	}
}

func (s *StopwordsList) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseStopwordsList(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseStopwordsList(input string) (*StopwordsList, error) {
	vals := map[string]StopwordsList{
		"arabic":     StopwordsListArabic,
		"armenian":   StopwordsListArmenian,
		"basque":     StopwordsListBasque,
		"brazilian":  StopwordsListBrazilian,
		"bulgarian":  StopwordsListBulgarian,
		"catalan":    StopwordsListCatalan,
		"czech":      StopwordsListCzech,
		"danish":     StopwordsListDanish,
		"dutch":      StopwordsListDutch,
		"english":    StopwordsListEnglish,
		"finnish":    StopwordsListFinnish,
		"french":     StopwordsListFrench,
		"galician":   StopwordsListGalician,
		"german":     StopwordsListGerman,
		"greek":      StopwordsListGreek,
		"hindi":      StopwordsListHindi,
		"hungarian":  StopwordsListHungarian,
		"indonesian": StopwordsListIndonesian,
		"irish":      StopwordsListIrish,
		"italian":    StopwordsListItalian,
		"latvian":    StopwordsListLatvian,
		"norwegian":  StopwordsListNorwegian,
		"persian":    StopwordsListPersian,
		"portuguese": StopwordsListPortuguese,
		"romanian":   StopwordsListRomanian,
		"russian":    StopwordsListRussian,
		"sorani":     StopwordsListSorani,
		"spanish":    StopwordsListSpanish,
		"swedish":    StopwordsListSwedish,
		"thai":       StopwordsListThai,
		"turkish":    StopwordsListTurkish,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := StopwordsList(input)
	return &out, nil
}

type SuggesterSearchMode string

const (
	SuggesterSearchModeAnalyzingInfixMatching SuggesterSearchMode = "analyzingInfixMatching"
)

func PossibleValuesForSuggesterSearchMode() []string {
	return []string{
		string(SuggesterSearchModeAnalyzingInfixMatching),
	}
}

func (s *SuggesterSearchMode) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseSuggesterSearchMode(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseSuggesterSearchMode(input string) (*SuggesterSearchMode, error) {
	vals := map[string]SuggesterSearchMode{
		"analyzinginfixmatching": SuggesterSearchModeAnalyzingInfixMatching,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := SuggesterSearchMode(input)
	return &out, nil
}

type TokenCharacterKind string

const (
	TokenCharacterKindDigit       TokenCharacterKind = "digit"
	TokenCharacterKindLetter      TokenCharacterKind = "letter"
	TokenCharacterKindPunctuation TokenCharacterKind = "punctuation"
	TokenCharacterKindSymbol      TokenCharacterKind = "symbol"
	TokenCharacterKindWhitespace  TokenCharacterKind = "whitespace"
)

func PossibleValuesForTokenCharacterKind() []string {
	return []string{
		string(TokenCharacterKindDigit),
		string(TokenCharacterKindLetter),
		string(TokenCharacterKindPunctuation),
		string(TokenCharacterKindSymbol),
		string(TokenCharacterKindWhitespace),
	}
}

func (s *TokenCharacterKind) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseTokenCharacterKind(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseTokenCharacterKind(input string) (*TokenCharacterKind, error) {
	vals := map[string]TokenCharacterKind{
		"digit":       TokenCharacterKindDigit,
		"letter":      TokenCharacterKindLetter,
		"punctuation": TokenCharacterKindPunctuation,
		"symbol":      TokenCharacterKindSymbol,
		"whitespace":  TokenCharacterKindWhitespace,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := TokenCharacterKind(input)
	return &out, nil
}

type TokenFilterName string

const (
	TokenFilterNameApostrophe                TokenFilterName = "apostrophe"
	TokenFilterNameArabicNormalization       TokenFilterName = "arabic_normalization"
	TokenFilterNameAsciifolding              TokenFilterName = "asciifolding"
	TokenFilterNameCjkBigram                 TokenFilterName = "cjk_bigram"
	TokenFilterNameCjkWidth                  TokenFilterName = "cjk_width"
	TokenFilterNameClassic                   TokenFilterName = "classic"
	TokenFilterNameCommonGrams               TokenFilterName = "common_grams"
	TokenFilterNameEdgeNGramVTwo             TokenFilterName = "edgeNGram_v2"
	TokenFilterNameElision                   TokenFilterName = "elision"
	TokenFilterNameGermanNormalization       TokenFilterName = "german_normalization"
	TokenFilterNameHindiNormalization        TokenFilterName = "hindi_normalization"
	TokenFilterNameIndicNormalization        TokenFilterName = "indic_normalization"
	TokenFilterNameKeywordRepeat             TokenFilterName = "keyword_repeat"
	TokenFilterNameKstem                     TokenFilterName = "kstem"
	TokenFilterNameLength                    TokenFilterName = "length"
	TokenFilterNameLimit                     TokenFilterName = "limit"
	TokenFilterNameLowercase                 TokenFilterName = "lowercase"
	TokenFilterNameNGramVTwo                 TokenFilterName = "nGram_v2"
	TokenFilterNamePersianNormalization      TokenFilterName = "persian_normalization"
	TokenFilterNamePhonetic                  TokenFilterName = "phonetic"
	TokenFilterNamePorterStem                TokenFilterName = "porter_stem"
	TokenFilterNameReverse                   TokenFilterName = "reverse"
	TokenFilterNameScandinavianFolding       TokenFilterName = "scandinavian_folding"
	TokenFilterNameScandinavianNormalization TokenFilterName = "scandinavian_normalization"
	TokenFilterNameShingle                   TokenFilterName = "shingle"
	TokenFilterNameSnowball                  TokenFilterName = "snowball"
	TokenFilterNameSoraniNormalization       TokenFilterName = "sorani_normalization"
	TokenFilterNameStemmer                   TokenFilterName = "stemmer"
	TokenFilterNameStopwords                 TokenFilterName = "stopwords"
	TokenFilterNameTrim                      TokenFilterName = "trim"
	TokenFilterNameTruncate                  TokenFilterName = "truncate"
	TokenFilterNameUnique                    TokenFilterName = "unique"
	TokenFilterNameUppercase                 TokenFilterName = "uppercase"
	TokenFilterNameWordDelimiter             TokenFilterName = "word_delimiter"
)

func PossibleValuesForTokenFilterName() []string {
	return []string{
		string(TokenFilterNameApostrophe),
		string(TokenFilterNameArabicNormalization),
		string(TokenFilterNameAsciifolding),
		string(TokenFilterNameCjkBigram),
		string(TokenFilterNameCjkWidth),
		string(TokenFilterNameClassic),
		string(TokenFilterNameCommonGrams),
		string(TokenFilterNameEdgeNGramVTwo),
		string(TokenFilterNameElision),
		string(TokenFilterNameGermanNormalization),
		string(TokenFilterNameHindiNormalization),
		string(TokenFilterNameIndicNormalization),
		string(TokenFilterNameKeywordRepeat),
		string(TokenFilterNameKstem),
		string(TokenFilterNameLength),
		string(TokenFilterNameLimit),
		string(TokenFilterNameLowercase),
		string(TokenFilterNameNGramVTwo),
		string(TokenFilterNamePersianNormalization),
		string(TokenFilterNamePhonetic),
		string(TokenFilterNamePorterStem),
		string(TokenFilterNameReverse),
		string(TokenFilterNameScandinavianFolding),
		string(TokenFilterNameScandinavianNormalization),
		string(TokenFilterNameShingle),
		string(TokenFilterNameSnowball),
		string(TokenFilterNameSoraniNormalization),
		string(TokenFilterNameStemmer),
		string(TokenFilterNameStopwords),
		string(TokenFilterNameTrim),
		string(TokenFilterNameTruncate),
		string(TokenFilterNameUnique),
		string(TokenFilterNameUppercase),
		string(TokenFilterNameWordDelimiter),
	}
}

func (s *TokenFilterName) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseTokenFilterName(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseTokenFilterName(input string) (*TokenFilterName, error) {
	vals := map[string]TokenFilterName{
		"apostrophe":                 TokenFilterNameApostrophe,
		"arabic_normalization":       TokenFilterNameArabicNormalization,
		"asciifolding":               TokenFilterNameAsciifolding,
		"cjk_bigram":                 TokenFilterNameCjkBigram,
		"cjk_width":                  TokenFilterNameCjkWidth,
		"classic":                    TokenFilterNameClassic,
		"common_grams":               TokenFilterNameCommonGrams,
		"edgengram_v2":               TokenFilterNameEdgeNGramVTwo,
		"elision":                    TokenFilterNameElision,
		"german_normalization":       TokenFilterNameGermanNormalization,
		"hindi_normalization":        TokenFilterNameHindiNormalization,
		"indic_normalization":        TokenFilterNameIndicNormalization,
		"keyword_repeat":             TokenFilterNameKeywordRepeat,
		"kstem":                      TokenFilterNameKstem,
		"length":                     TokenFilterNameLength,
		"limit":                      TokenFilterNameLimit,
		"lowercase":                  TokenFilterNameLowercase,
		"ngram_v2":                   TokenFilterNameNGramVTwo,
		"persian_normalization":      TokenFilterNamePersianNormalization,
		"phonetic":                   TokenFilterNamePhonetic,
		"porter_stem":                TokenFilterNamePorterStem,
		"reverse":                    TokenFilterNameReverse,
		"scandinavian_folding":       TokenFilterNameScandinavianFolding,
		"scandinavian_normalization": TokenFilterNameScandinavianNormalization,
		"shingle":                    TokenFilterNameShingle,
		"snowball":                   TokenFilterNameSnowball,
		"sorani_normalization":       TokenFilterNameSoraniNormalization,
		"stemmer":                    TokenFilterNameStemmer,
		"stopwords":                  TokenFilterNameStopwords,
		"trim":                       TokenFilterNameTrim,
		"truncate":                   TokenFilterNameTruncate,
		"unique":                     TokenFilterNameUnique,
		"uppercase":                  TokenFilterNameUppercase,
		"word_delimiter":             TokenFilterNameWordDelimiter,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := TokenFilterName(input)
	return &out, nil
}

type VectorEncodingFormat string

const (
	VectorEncodingFormatPackedBit VectorEncodingFormat = "packedBit"
)

func PossibleValuesForVectorEncodingFormat() []string {
	return []string{
		string(VectorEncodingFormatPackedBit),
	}
}

func (s *VectorEncodingFormat) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseVectorEncodingFormat(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseVectorEncodingFormat(input string) (*VectorEncodingFormat, error) {
	vals := map[string]VectorEncodingFormat{
		"packedbit": VectorEncodingFormatPackedBit,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := VectorEncodingFormat(input)
	return &out, nil
}

type VectorSearchAlgorithmKind string

const (
	VectorSearchAlgorithmKindExhaustiveKnn VectorSearchAlgorithmKind = "exhaustiveKnn"
	VectorSearchAlgorithmKindHnsw          VectorSearchAlgorithmKind = "hnsw"
)

func PossibleValuesForVectorSearchAlgorithmKind() []string {
	return []string{
		string(VectorSearchAlgorithmKindExhaustiveKnn),
		string(VectorSearchAlgorithmKindHnsw),
	}
}

func (s *VectorSearchAlgorithmKind) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseVectorSearchAlgorithmKind(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseVectorSearchAlgorithmKind(input string) (*VectorSearchAlgorithmKind, error) {
	vals := map[string]VectorSearchAlgorithmKind{
		"exhaustiveknn": VectorSearchAlgorithmKindExhaustiveKnn,
		"hnsw":          VectorSearchAlgorithmKindHnsw,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := VectorSearchAlgorithmKind(input)
	return &out, nil
}

type VectorSearchAlgorithmMetric string

const (
	VectorSearchAlgorithmMetricCosine     VectorSearchAlgorithmMetric = "cosine"
	VectorSearchAlgorithmMetricDotProduct VectorSearchAlgorithmMetric = "dotProduct"
	VectorSearchAlgorithmMetricEuclidean  VectorSearchAlgorithmMetric = "euclidean"
	VectorSearchAlgorithmMetricHamming    VectorSearchAlgorithmMetric = "hamming"
)

func PossibleValuesForVectorSearchAlgorithmMetric() []string {
	return []string{
		string(VectorSearchAlgorithmMetricCosine),
		string(VectorSearchAlgorithmMetricDotProduct),
		string(VectorSearchAlgorithmMetricEuclidean),
		string(VectorSearchAlgorithmMetricHamming),
	}
}

func (s *VectorSearchAlgorithmMetric) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseVectorSearchAlgorithmMetric(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseVectorSearchAlgorithmMetric(input string) (*VectorSearchAlgorithmMetric, error) {
	vals := map[string]VectorSearchAlgorithmMetric{
		"cosine":     VectorSearchAlgorithmMetricCosine,
		"dotproduct": VectorSearchAlgorithmMetricDotProduct,
		"euclidean":  VectorSearchAlgorithmMetricEuclidean,
		"hamming":    VectorSearchAlgorithmMetricHamming,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := VectorSearchAlgorithmMetric(input)
	return &out, nil
}

type VectorSearchCompressionKind string

const (
	VectorSearchCompressionKindBinaryQuantization VectorSearchCompressionKind = "binaryQuantization"
	VectorSearchCompressionKindScalarQuantization VectorSearchCompressionKind = "scalarQuantization"
)

func PossibleValuesForVectorSearchCompressionKind() []string {
	return []string{
		string(VectorSearchCompressionKindBinaryQuantization),
		string(VectorSearchCompressionKindScalarQuantization),
	}
}

func (s *VectorSearchCompressionKind) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseVectorSearchCompressionKind(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseVectorSearchCompressionKind(input string) (*VectorSearchCompressionKind, error) {
	vals := map[string]VectorSearchCompressionKind{
		"binaryquantization": VectorSearchCompressionKindBinaryQuantization,
		"scalarquantization": VectorSearchCompressionKindScalarQuantization,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := VectorSearchCompressionKind(input)
	return &out, nil
}

type VectorSearchCompressionRescoreStorageMethod string

const (
	VectorSearchCompressionRescoreStorageMethodDiscardOriginals  VectorSearchCompressionRescoreStorageMethod = "discardOriginals"
	VectorSearchCompressionRescoreStorageMethodPreserveOriginals VectorSearchCompressionRescoreStorageMethod = "preserveOriginals"
)

func PossibleValuesForVectorSearchCompressionRescoreStorageMethod() []string {
	return []string{
		string(VectorSearchCompressionRescoreStorageMethodDiscardOriginals),
		string(VectorSearchCompressionRescoreStorageMethodPreserveOriginals),
	}
}

func (s *VectorSearchCompressionRescoreStorageMethod) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseVectorSearchCompressionRescoreStorageMethod(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseVectorSearchCompressionRescoreStorageMethod(input string) (*VectorSearchCompressionRescoreStorageMethod, error) {
	vals := map[string]VectorSearchCompressionRescoreStorageMethod{
		"discardoriginals":  VectorSearchCompressionRescoreStorageMethodDiscardOriginals,
		"preserveoriginals": VectorSearchCompressionRescoreStorageMethodPreserveOriginals,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := VectorSearchCompressionRescoreStorageMethod(input)
	return &out, nil
}

type VectorSearchCompressionTargetDataType string

const (
	VectorSearchCompressionTargetDataTypeIntEight VectorSearchCompressionTargetDataType = "int8"
)

func PossibleValuesForVectorSearchCompressionTargetDataType() []string {
	return []string{
		string(VectorSearchCompressionTargetDataTypeIntEight),
	}
}

func (s *VectorSearchCompressionTargetDataType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseVectorSearchCompressionTargetDataType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseVectorSearchCompressionTargetDataType(input string) (*VectorSearchCompressionTargetDataType, error) {
	vals := map[string]VectorSearchCompressionTargetDataType{
		"int8": VectorSearchCompressionTargetDataTypeIntEight,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := VectorSearchCompressionTargetDataType(input)
	return &out, nil
}

type VectorSearchVectorizerKind string

const (
	VectorSearchVectorizerKindAzureOpenAI  VectorSearchVectorizerKind = "azureOpenAI"
	VectorSearchVectorizerKindCustomWebApi VectorSearchVectorizerKind = "customWebApi"
)

func PossibleValuesForVectorSearchVectorizerKind() []string {
	return []string{
		string(VectorSearchVectorizerKindAzureOpenAI),
		string(VectorSearchVectorizerKindCustomWebApi),
	}
}

func (s *VectorSearchVectorizerKind) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseVectorSearchVectorizerKind(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseVectorSearchVectorizerKind(input string) (*VectorSearchVectorizerKind, error) {
	vals := map[string]VectorSearchVectorizerKind{
		"azureopenai":  VectorSearchVectorizerKindAzureOpenAI,
		"customwebapi": VectorSearchVectorizerKindCustomWebApi,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := VectorSearchVectorizerKind(input)
	return &out, nil
}
