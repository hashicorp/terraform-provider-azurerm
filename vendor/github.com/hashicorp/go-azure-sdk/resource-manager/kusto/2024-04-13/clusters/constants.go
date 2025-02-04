package clusters

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AzureScaleType string

const (
	AzureScaleTypeAutomatic AzureScaleType = "automatic"
	AzureScaleTypeManual    AzureScaleType = "manual"
	AzureScaleTypeNone      AzureScaleType = "none"
)

func PossibleValuesForAzureScaleType() []string {
	return []string{
		string(AzureScaleTypeAutomatic),
		string(AzureScaleTypeManual),
		string(AzureScaleTypeNone),
	}
}

func (s *AzureScaleType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseAzureScaleType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseAzureScaleType(input string) (*AzureScaleType, error) {
	vals := map[string]AzureScaleType{
		"automatic": AzureScaleTypeAutomatic,
		"manual":    AzureScaleTypeManual,
		"none":      AzureScaleTypeNone,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := AzureScaleType(input)
	return &out, nil
}

type AzureSkuName string

const (
	AzureSkuNameDevNoSLAStandardDOneOneVTwo              AzureSkuName = "Dev(No SLA)_Standard_D11_v2"
	AzureSkuNameDevNoSLAStandardETwoaVFour               AzureSkuName = "Dev(No SLA)_Standard_E2a_v4"
	AzureSkuNameStandardDOneFourVTwo                     AzureSkuName = "Standard_D14_v2"
	AzureSkuNameStandardDOneOneVTwo                      AzureSkuName = "Standard_D11_v2"
	AzureSkuNameStandardDOneSixdVFive                    AzureSkuName = "Standard_D16d_v5"
	AzureSkuNameStandardDOneThreeVTwo                    AzureSkuName = "Standard_D13_v2"
	AzureSkuNameStandardDOneTwoVTwo                      AzureSkuName = "Standard_D12_v2"
	AzureSkuNameStandardDSOneFourVTwoPositiveFourTBPS    AzureSkuName = "Standard_DS14_v2+4TB_PS"
	AzureSkuNameStandardDSOneFourVTwoPositiveThreeTBPS   AzureSkuName = "Standard_DS14_v2+3TB_PS"
	AzureSkuNameStandardDSOneThreeVTwoPositiveOneTBPS    AzureSkuName = "Standard_DS13_v2+1TB_PS"
	AzureSkuNameStandardDSOneThreeVTwoPositiveTwoTBPS    AzureSkuName = "Standard_DS13_v2+2TB_PS"
	AzureSkuNameStandardDThreeTwodVFive                  AzureSkuName = "Standard_D32d_v5"
	AzureSkuNameStandardDThreeTwodVFour                  AzureSkuName = "Standard_D32d_v4"
	AzureSkuNameStandardECEightadsVFive                  AzureSkuName = "Standard_EC8ads_v5"
	AzureSkuNameStandardECEightasVFivePositiveOneTBPS    AzureSkuName = "Standard_EC8as_v5+1TB_PS"
	AzureSkuNameStandardECEightasVFivePositiveTwoTBPS    AzureSkuName = "Standard_EC8as_v5+2TB_PS"
	AzureSkuNameStandardECOneSixadsVFive                 AzureSkuName = "Standard_EC16ads_v5"
	AzureSkuNameStandardECOneSixasVFivePositiveFourTBPS  AzureSkuName = "Standard_EC16as_v5+4TB_PS"
	AzureSkuNameStandardECOneSixasVFivePositiveThreeTBPS AzureSkuName = "Standard_EC16as_v5+3TB_PS"
	AzureSkuNameStandardEEightZeroidsVFour               AzureSkuName = "Standard_E80ids_v4"
	AzureSkuNameStandardEEightaVFour                     AzureSkuName = "Standard_E8a_v4"
	AzureSkuNameStandardEEightadsVFive                   AzureSkuName = "Standard_E8ads_v5"
	AzureSkuNameStandardEEightasVFivePositiveOneTBPS     AzureSkuName = "Standard_E8as_v5+1TB_PS"
	AzureSkuNameStandardEEightasVFivePositiveTwoTBPS     AzureSkuName = "Standard_E8as_v5+2TB_PS"
	AzureSkuNameStandardEEightasVFourPositiveOneTBPS     AzureSkuName = "Standard_E8as_v4+1TB_PS"
	AzureSkuNameStandardEEightasVFourPositiveTwoTBPS     AzureSkuName = "Standard_E8as_v4+2TB_PS"
	AzureSkuNameStandardEEightdVFive                     AzureSkuName = "Standard_E8d_v5"
	AzureSkuNameStandardEEightdVFour                     AzureSkuName = "Standard_E8d_v4"
	AzureSkuNameStandardEEightsVFivePositiveOneTBPS      AzureSkuName = "Standard_E8s_v5+1TB_PS"
	AzureSkuNameStandardEEightsVFivePositiveTwoTBPS      AzureSkuName = "Standard_E8s_v5+2TB_PS"
	AzureSkuNameStandardEEightsVFourPositiveOneTBPS      AzureSkuName = "Standard_E8s_v4+1TB_PS"
	AzureSkuNameStandardEEightsVFourPositiveTwoTBPS      AzureSkuName = "Standard_E8s_v4+2TB_PS"
	AzureSkuNameStandardEFouraVFour                      AzureSkuName = "Standard_E4a_v4"
	AzureSkuNameStandardEFouradsVFive                    AzureSkuName = "Standard_E4ads_v5"
	AzureSkuNameStandardEFourdVFive                      AzureSkuName = "Standard_E4d_v5"
	AzureSkuNameStandardEFourdVFour                      AzureSkuName = "Standard_E4d_v4"
	AzureSkuNameStandardEOneSixaVFour                    AzureSkuName = "Standard_E16a_v4"
	AzureSkuNameStandardEOneSixadsVFive                  AzureSkuName = "Standard_E16ads_v5"
	AzureSkuNameStandardEOneSixasVFivePositiveFourTBPS   AzureSkuName = "Standard_E16as_v5+4TB_PS"
	AzureSkuNameStandardEOneSixasVFivePositiveThreeTBPS  AzureSkuName = "Standard_E16as_v5+3TB_PS"
	AzureSkuNameStandardEOneSixasVFourPositiveFourTBPS   AzureSkuName = "Standard_E16as_v4+4TB_PS"
	AzureSkuNameStandardEOneSixasVFourPositiveThreeTBPS  AzureSkuName = "Standard_E16as_v4+3TB_PS"
	AzureSkuNameStandardEOneSixdVFive                    AzureSkuName = "Standard_E16d_v5"
	AzureSkuNameStandardEOneSixdVFour                    AzureSkuName = "Standard_E16d_v4"
	AzureSkuNameStandardEOneSixsVFivePositiveFourTBPS    AzureSkuName = "Standard_E16s_v5+4TB_PS"
	AzureSkuNameStandardEOneSixsVFivePositiveThreeTBPS   AzureSkuName = "Standard_E16s_v5+3TB_PS"
	AzureSkuNameStandardEOneSixsVFourPositiveFourTBPS    AzureSkuName = "Standard_E16s_v4+4TB_PS"
	AzureSkuNameStandardEOneSixsVFourPositiveThreeTBPS   AzureSkuName = "Standard_E16s_v4+3TB_PS"
	AzureSkuNameStandardESixFouriVThree                  AzureSkuName = "Standard_E64i_v3"
	AzureSkuNameStandardETwoaVFour                       AzureSkuName = "Standard_E2a_v4"
	AzureSkuNameStandardETwoadsVFive                     AzureSkuName = "Standard_E2ads_v5"
	AzureSkuNameStandardETwodVFive                       AzureSkuName = "Standard_E2d_v5"
	AzureSkuNameStandardETwodVFour                       AzureSkuName = "Standard_E2d_v4"
	AzureSkuNameStandardLEightasVThree                   AzureSkuName = "Standard_L8as_v3"
	AzureSkuNameStandardLEights                          AzureSkuName = "Standard_L8s"
	AzureSkuNameStandardLEightsVThree                    AzureSkuName = "Standard_L8s_v3"
	AzureSkuNameStandardLEightsVTwo                      AzureSkuName = "Standard_L8s_v2"
	AzureSkuNameStandardLFours                           AzureSkuName = "Standard_L4s"
	AzureSkuNameStandardLOneSixasVThree                  AzureSkuName = "Standard_L16as_v3"
	AzureSkuNameStandardLOneSixs                         AzureSkuName = "Standard_L16s"
	AzureSkuNameStandardLOneSixsVThree                   AzureSkuName = "Standard_L16s_v3"
	AzureSkuNameStandardLOneSixsVTwo                     AzureSkuName = "Standard_L16s_v2"
	AzureSkuNameStandardLThreeTwoasVThree                AzureSkuName = "Standard_L32as_v3"
	AzureSkuNameStandardLThreeTwosVThree                 AzureSkuName = "Standard_L32s_v3"
)

func PossibleValuesForAzureSkuName() []string {
	return []string{
		string(AzureSkuNameDevNoSLAStandardDOneOneVTwo),
		string(AzureSkuNameDevNoSLAStandardETwoaVFour),
		string(AzureSkuNameStandardDOneFourVTwo),
		string(AzureSkuNameStandardDOneOneVTwo),
		string(AzureSkuNameStandardDOneSixdVFive),
		string(AzureSkuNameStandardDOneThreeVTwo),
		string(AzureSkuNameStandardDOneTwoVTwo),
		string(AzureSkuNameStandardDSOneFourVTwoPositiveFourTBPS),
		string(AzureSkuNameStandardDSOneFourVTwoPositiveThreeTBPS),
		string(AzureSkuNameStandardDSOneThreeVTwoPositiveOneTBPS),
		string(AzureSkuNameStandardDSOneThreeVTwoPositiveTwoTBPS),
		string(AzureSkuNameStandardDThreeTwodVFive),
		string(AzureSkuNameStandardDThreeTwodVFour),
		string(AzureSkuNameStandardECEightadsVFive),
		string(AzureSkuNameStandardECEightasVFivePositiveOneTBPS),
		string(AzureSkuNameStandardECEightasVFivePositiveTwoTBPS),
		string(AzureSkuNameStandardECOneSixadsVFive),
		string(AzureSkuNameStandardECOneSixasVFivePositiveFourTBPS),
		string(AzureSkuNameStandardECOneSixasVFivePositiveThreeTBPS),
		string(AzureSkuNameStandardEEightZeroidsVFour),
		string(AzureSkuNameStandardEEightaVFour),
		string(AzureSkuNameStandardEEightadsVFive),
		string(AzureSkuNameStandardEEightasVFivePositiveOneTBPS),
		string(AzureSkuNameStandardEEightasVFivePositiveTwoTBPS),
		string(AzureSkuNameStandardEEightasVFourPositiveOneTBPS),
		string(AzureSkuNameStandardEEightasVFourPositiveTwoTBPS),
		string(AzureSkuNameStandardEEightdVFive),
		string(AzureSkuNameStandardEEightdVFour),
		string(AzureSkuNameStandardEEightsVFivePositiveOneTBPS),
		string(AzureSkuNameStandardEEightsVFivePositiveTwoTBPS),
		string(AzureSkuNameStandardEEightsVFourPositiveOneTBPS),
		string(AzureSkuNameStandardEEightsVFourPositiveTwoTBPS),
		string(AzureSkuNameStandardEFouraVFour),
		string(AzureSkuNameStandardEFouradsVFive),
		string(AzureSkuNameStandardEFourdVFive),
		string(AzureSkuNameStandardEFourdVFour),
		string(AzureSkuNameStandardEOneSixaVFour),
		string(AzureSkuNameStandardEOneSixadsVFive),
		string(AzureSkuNameStandardEOneSixasVFivePositiveFourTBPS),
		string(AzureSkuNameStandardEOneSixasVFivePositiveThreeTBPS),
		string(AzureSkuNameStandardEOneSixasVFourPositiveFourTBPS),
		string(AzureSkuNameStandardEOneSixasVFourPositiveThreeTBPS),
		string(AzureSkuNameStandardEOneSixdVFive),
		string(AzureSkuNameStandardEOneSixdVFour),
		string(AzureSkuNameStandardEOneSixsVFivePositiveFourTBPS),
		string(AzureSkuNameStandardEOneSixsVFivePositiveThreeTBPS),
		string(AzureSkuNameStandardEOneSixsVFourPositiveFourTBPS),
		string(AzureSkuNameStandardEOneSixsVFourPositiveThreeTBPS),
		string(AzureSkuNameStandardESixFouriVThree),
		string(AzureSkuNameStandardETwoaVFour),
		string(AzureSkuNameStandardETwoadsVFive),
		string(AzureSkuNameStandardETwodVFive),
		string(AzureSkuNameStandardETwodVFour),
		string(AzureSkuNameStandardLEightasVThree),
		string(AzureSkuNameStandardLEights),
		string(AzureSkuNameStandardLEightsVThree),
		string(AzureSkuNameStandardLEightsVTwo),
		string(AzureSkuNameStandardLFours),
		string(AzureSkuNameStandardLOneSixasVThree),
		string(AzureSkuNameStandardLOneSixs),
		string(AzureSkuNameStandardLOneSixsVThree),
		string(AzureSkuNameStandardLOneSixsVTwo),
		string(AzureSkuNameStandardLThreeTwoasVThree),
		string(AzureSkuNameStandardLThreeTwosVThree),
	}
}

func (s *AzureSkuName) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseAzureSkuName(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseAzureSkuName(input string) (*AzureSkuName, error) {
	vals := map[string]AzureSkuName{
		"dev(no sla)_standard_d11_v2": AzureSkuNameDevNoSLAStandardDOneOneVTwo,
		"dev(no sla)_standard_e2a_v4": AzureSkuNameDevNoSLAStandardETwoaVFour,
		"standard_d14_v2":             AzureSkuNameStandardDOneFourVTwo,
		"standard_d11_v2":             AzureSkuNameStandardDOneOneVTwo,
		"standard_d16d_v5":            AzureSkuNameStandardDOneSixdVFive,
		"standard_d13_v2":             AzureSkuNameStandardDOneThreeVTwo,
		"standard_d12_v2":             AzureSkuNameStandardDOneTwoVTwo,
		"standard_ds14_v2+4tb_ps":     AzureSkuNameStandardDSOneFourVTwoPositiveFourTBPS,
		"standard_ds14_v2+3tb_ps":     AzureSkuNameStandardDSOneFourVTwoPositiveThreeTBPS,
		"standard_ds13_v2+1tb_ps":     AzureSkuNameStandardDSOneThreeVTwoPositiveOneTBPS,
		"standard_ds13_v2+2tb_ps":     AzureSkuNameStandardDSOneThreeVTwoPositiveTwoTBPS,
		"standard_d32d_v5":            AzureSkuNameStandardDThreeTwodVFive,
		"standard_d32d_v4":            AzureSkuNameStandardDThreeTwodVFour,
		"standard_ec8ads_v5":          AzureSkuNameStandardECEightadsVFive,
		"standard_ec8as_v5+1tb_ps":    AzureSkuNameStandardECEightasVFivePositiveOneTBPS,
		"standard_ec8as_v5+2tb_ps":    AzureSkuNameStandardECEightasVFivePositiveTwoTBPS,
		"standard_ec16ads_v5":         AzureSkuNameStandardECOneSixadsVFive,
		"standard_ec16as_v5+4tb_ps":   AzureSkuNameStandardECOneSixasVFivePositiveFourTBPS,
		"standard_ec16as_v5+3tb_ps":   AzureSkuNameStandardECOneSixasVFivePositiveThreeTBPS,
		"standard_e80ids_v4":          AzureSkuNameStandardEEightZeroidsVFour,
		"standard_e8a_v4":             AzureSkuNameStandardEEightaVFour,
		"standard_e8ads_v5":           AzureSkuNameStandardEEightadsVFive,
		"standard_e8as_v5+1tb_ps":     AzureSkuNameStandardEEightasVFivePositiveOneTBPS,
		"standard_e8as_v5+2tb_ps":     AzureSkuNameStandardEEightasVFivePositiveTwoTBPS,
		"standard_e8as_v4+1tb_ps":     AzureSkuNameStandardEEightasVFourPositiveOneTBPS,
		"standard_e8as_v4+2tb_ps":     AzureSkuNameStandardEEightasVFourPositiveTwoTBPS,
		"standard_e8d_v5":             AzureSkuNameStandardEEightdVFive,
		"standard_e8d_v4":             AzureSkuNameStandardEEightdVFour,
		"standard_e8s_v5+1tb_ps":      AzureSkuNameStandardEEightsVFivePositiveOneTBPS,
		"standard_e8s_v5+2tb_ps":      AzureSkuNameStandardEEightsVFivePositiveTwoTBPS,
		"standard_e8s_v4+1tb_ps":      AzureSkuNameStandardEEightsVFourPositiveOneTBPS,
		"standard_e8s_v4+2tb_ps":      AzureSkuNameStandardEEightsVFourPositiveTwoTBPS,
		"standard_e4a_v4":             AzureSkuNameStandardEFouraVFour,
		"standard_e4ads_v5":           AzureSkuNameStandardEFouradsVFive,
		"standard_e4d_v5":             AzureSkuNameStandardEFourdVFive,
		"standard_e4d_v4":             AzureSkuNameStandardEFourdVFour,
		"standard_e16a_v4":            AzureSkuNameStandardEOneSixaVFour,
		"standard_e16ads_v5":          AzureSkuNameStandardEOneSixadsVFive,
		"standard_e16as_v5+4tb_ps":    AzureSkuNameStandardEOneSixasVFivePositiveFourTBPS,
		"standard_e16as_v5+3tb_ps":    AzureSkuNameStandardEOneSixasVFivePositiveThreeTBPS,
		"standard_e16as_v4+4tb_ps":    AzureSkuNameStandardEOneSixasVFourPositiveFourTBPS,
		"standard_e16as_v4+3tb_ps":    AzureSkuNameStandardEOneSixasVFourPositiveThreeTBPS,
		"standard_e16d_v5":            AzureSkuNameStandardEOneSixdVFive,
		"standard_e16d_v4":            AzureSkuNameStandardEOneSixdVFour,
		"standard_e16s_v5+4tb_ps":     AzureSkuNameStandardEOneSixsVFivePositiveFourTBPS,
		"standard_e16s_v5+3tb_ps":     AzureSkuNameStandardEOneSixsVFivePositiveThreeTBPS,
		"standard_e16s_v4+4tb_ps":     AzureSkuNameStandardEOneSixsVFourPositiveFourTBPS,
		"standard_e16s_v4+3tb_ps":     AzureSkuNameStandardEOneSixsVFourPositiveThreeTBPS,
		"standard_e64i_v3":            AzureSkuNameStandardESixFouriVThree,
		"standard_e2a_v4":             AzureSkuNameStandardETwoaVFour,
		"standard_e2ads_v5":           AzureSkuNameStandardETwoadsVFive,
		"standard_e2d_v5":             AzureSkuNameStandardETwodVFive,
		"standard_e2d_v4":             AzureSkuNameStandardETwodVFour,
		"standard_l8as_v3":            AzureSkuNameStandardLEightasVThree,
		"standard_l8s":                AzureSkuNameStandardLEights,
		"standard_l8s_v3":             AzureSkuNameStandardLEightsVThree,
		"standard_l8s_v2":             AzureSkuNameStandardLEightsVTwo,
		"standard_l4s":                AzureSkuNameStandardLFours,
		"standard_l16as_v3":           AzureSkuNameStandardLOneSixasVThree,
		"standard_l16s":               AzureSkuNameStandardLOneSixs,
		"standard_l16s_v3":            AzureSkuNameStandardLOneSixsVThree,
		"standard_l16s_v2":            AzureSkuNameStandardLOneSixsVTwo,
		"standard_l32as_v3":           AzureSkuNameStandardLThreeTwoasVThree,
		"standard_l32s_v3":            AzureSkuNameStandardLThreeTwosVThree,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := AzureSkuName(input)
	return &out, nil
}

type AzureSkuTier string

const (
	AzureSkuTierBasic    AzureSkuTier = "Basic"
	AzureSkuTierStandard AzureSkuTier = "Standard"
)

func PossibleValuesForAzureSkuTier() []string {
	return []string{
		string(AzureSkuTierBasic),
		string(AzureSkuTierStandard),
	}
}

func (s *AzureSkuTier) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseAzureSkuTier(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseAzureSkuTier(input string) (*AzureSkuTier, error) {
	vals := map[string]AzureSkuTier{
		"basic":    AzureSkuTierBasic,
		"standard": AzureSkuTierStandard,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := AzureSkuTier(input)
	return &out, nil
}

type ClusterNetworkAccessFlag string

const (
	ClusterNetworkAccessFlagDisabled ClusterNetworkAccessFlag = "Disabled"
	ClusterNetworkAccessFlagEnabled  ClusterNetworkAccessFlag = "Enabled"
)

func PossibleValuesForClusterNetworkAccessFlag() []string {
	return []string{
		string(ClusterNetworkAccessFlagDisabled),
		string(ClusterNetworkAccessFlagEnabled),
	}
}

func (s *ClusterNetworkAccessFlag) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseClusterNetworkAccessFlag(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseClusterNetworkAccessFlag(input string) (*ClusterNetworkAccessFlag, error) {
	vals := map[string]ClusterNetworkAccessFlag{
		"disabled": ClusterNetworkAccessFlagDisabled,
		"enabled":  ClusterNetworkAccessFlagEnabled,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ClusterNetworkAccessFlag(input)
	return &out, nil
}

type ClusterType string

const (
	ClusterTypeMicrosoftPointKustoClusters ClusterType = "Microsoft.Kusto/clusters"
)

func PossibleValuesForClusterType() []string {
	return []string{
		string(ClusterTypeMicrosoftPointKustoClusters),
	}
}

func (s *ClusterType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseClusterType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseClusterType(input string) (*ClusterType, error) {
	vals := map[string]ClusterType{
		"microsoft.kusto/clusters": ClusterTypeMicrosoftPointKustoClusters,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ClusterType(input)
	return &out, nil
}

type DatabaseShareOrigin string

const (
	DatabaseShareOriginDataShare DatabaseShareOrigin = "DataShare"
	DatabaseShareOriginDirect    DatabaseShareOrigin = "Direct"
	DatabaseShareOriginOther     DatabaseShareOrigin = "Other"
)

func PossibleValuesForDatabaseShareOrigin() []string {
	return []string{
		string(DatabaseShareOriginDataShare),
		string(DatabaseShareOriginDirect),
		string(DatabaseShareOriginOther),
	}
}

func (s *DatabaseShareOrigin) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseDatabaseShareOrigin(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseDatabaseShareOrigin(input string) (*DatabaseShareOrigin, error) {
	vals := map[string]DatabaseShareOrigin{
		"datashare": DatabaseShareOriginDataShare,
		"direct":    DatabaseShareOriginDirect,
		"other":     DatabaseShareOriginOther,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := DatabaseShareOrigin(input)
	return &out, nil
}

type EngineType string

const (
	EngineTypeVThree EngineType = "V3"
	EngineTypeVTwo   EngineType = "V2"
)

func PossibleValuesForEngineType() []string {
	return []string{
		string(EngineTypeVThree),
		string(EngineTypeVTwo),
	}
}

func (s *EngineType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseEngineType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseEngineType(input string) (*EngineType, error) {
	vals := map[string]EngineType{
		"v3": EngineTypeVThree,
		"v2": EngineTypeVTwo,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := EngineType(input)
	return &out, nil
}

type LanguageExtensionImageName string

const (
	LanguageExtensionImageNamePythonCustomImage         LanguageExtensionImageName = "PythonCustomImage"
	LanguageExtensionImageNamePythonThreeOneZeroEight   LanguageExtensionImageName = "Python3_10_8"
	LanguageExtensionImageNamePythonThreeOneZeroEightDL LanguageExtensionImageName = "Python3_10_8_DL"
	LanguageExtensionImageNamePythonThreeSixFive        LanguageExtensionImageName = "Python3_6_5"
	LanguageExtensionImageNameR                         LanguageExtensionImageName = "R"
)

func PossibleValuesForLanguageExtensionImageName() []string {
	return []string{
		string(LanguageExtensionImageNamePythonCustomImage),
		string(LanguageExtensionImageNamePythonThreeOneZeroEight),
		string(LanguageExtensionImageNamePythonThreeOneZeroEightDL),
		string(LanguageExtensionImageNamePythonThreeSixFive),
		string(LanguageExtensionImageNameR),
	}
}

func (s *LanguageExtensionImageName) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseLanguageExtensionImageName(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseLanguageExtensionImageName(input string) (*LanguageExtensionImageName, error) {
	vals := map[string]LanguageExtensionImageName{
		"pythoncustomimage": LanguageExtensionImageNamePythonCustomImage,
		"python3_10_8":      LanguageExtensionImageNamePythonThreeOneZeroEight,
		"python3_10_8_dl":   LanguageExtensionImageNamePythonThreeOneZeroEightDL,
		"python3_6_5":       LanguageExtensionImageNamePythonThreeSixFive,
		"r":                 LanguageExtensionImageNameR,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := LanguageExtensionImageName(input)
	return &out, nil
}

type LanguageExtensionName string

const (
	LanguageExtensionNamePYTHON LanguageExtensionName = "PYTHON"
	LanguageExtensionNameR      LanguageExtensionName = "R"
)

func PossibleValuesForLanguageExtensionName() []string {
	return []string{
		string(LanguageExtensionNamePYTHON),
		string(LanguageExtensionNameR),
	}
}

func (s *LanguageExtensionName) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseLanguageExtensionName(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseLanguageExtensionName(input string) (*LanguageExtensionName, error) {
	vals := map[string]LanguageExtensionName{
		"python": LanguageExtensionNamePYTHON,
		"r":      LanguageExtensionNameR,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := LanguageExtensionName(input)
	return &out, nil
}

type MigrationClusterRole string

const (
	MigrationClusterRoleDestination MigrationClusterRole = "Destination"
	MigrationClusterRoleSource      MigrationClusterRole = "Source"
)

func PossibleValuesForMigrationClusterRole() []string {
	return []string{
		string(MigrationClusterRoleDestination),
		string(MigrationClusterRoleSource),
	}
}

func (s *MigrationClusterRole) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseMigrationClusterRole(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseMigrationClusterRole(input string) (*MigrationClusterRole, error) {
	vals := map[string]MigrationClusterRole{
		"destination": MigrationClusterRoleDestination,
		"source":      MigrationClusterRoleSource,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := MigrationClusterRole(input)
	return &out, nil
}

type ProvisioningState string

const (
	ProvisioningStateCanceled  ProvisioningState = "Canceled"
	ProvisioningStateCreating  ProvisioningState = "Creating"
	ProvisioningStateDeleting  ProvisioningState = "Deleting"
	ProvisioningStateFailed    ProvisioningState = "Failed"
	ProvisioningStateMoving    ProvisioningState = "Moving"
	ProvisioningStateRunning   ProvisioningState = "Running"
	ProvisioningStateSucceeded ProvisioningState = "Succeeded"
)

func PossibleValuesForProvisioningState() []string {
	return []string{
		string(ProvisioningStateCanceled),
		string(ProvisioningStateCreating),
		string(ProvisioningStateDeleting),
		string(ProvisioningStateFailed),
		string(ProvisioningStateMoving),
		string(ProvisioningStateRunning),
		string(ProvisioningStateSucceeded),
	}
}

func (s *ProvisioningState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseProvisioningState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseProvisioningState(input string) (*ProvisioningState, error) {
	vals := map[string]ProvisioningState{
		"canceled":  ProvisioningStateCanceled,
		"creating":  ProvisioningStateCreating,
		"deleting":  ProvisioningStateDeleting,
		"failed":    ProvisioningStateFailed,
		"moving":    ProvisioningStateMoving,
		"running":   ProvisioningStateRunning,
		"succeeded": ProvisioningStateSucceeded,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ProvisioningState(input)
	return &out, nil
}

type PublicIPType string

const (
	PublicIPTypeDualStack PublicIPType = "DualStack"
	PublicIPTypeIPvFour   PublicIPType = "IPv4"
)

func PossibleValuesForPublicIPType() []string {
	return []string{
		string(PublicIPTypeDualStack),
		string(PublicIPTypeIPvFour),
	}
}

func (s *PublicIPType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parsePublicIPType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parsePublicIPType(input string) (*PublicIPType, error) {
	vals := map[string]PublicIPType{
		"dualstack": PublicIPTypeDualStack,
		"ipv4":      PublicIPTypeIPvFour,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := PublicIPType(input)
	return &out, nil
}

type PublicNetworkAccess string

const (
	PublicNetworkAccessDisabled PublicNetworkAccess = "Disabled"
	PublicNetworkAccessEnabled  PublicNetworkAccess = "Enabled"
)

func PossibleValuesForPublicNetworkAccess() []string {
	return []string{
		string(PublicNetworkAccessDisabled),
		string(PublicNetworkAccessEnabled),
	}
}

func (s *PublicNetworkAccess) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parsePublicNetworkAccess(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parsePublicNetworkAccess(input string) (*PublicNetworkAccess, error) {
	vals := map[string]PublicNetworkAccess{
		"disabled": PublicNetworkAccessDisabled,
		"enabled":  PublicNetworkAccessEnabled,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := PublicNetworkAccess(input)
	return &out, nil
}

type Reason string

const (
	ReasonAlreadyExists Reason = "AlreadyExists"
	ReasonInvalid       Reason = "Invalid"
)

func PossibleValuesForReason() []string {
	return []string{
		string(ReasonAlreadyExists),
		string(ReasonInvalid),
	}
}

func (s *Reason) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseReason(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseReason(input string) (*Reason, error) {
	vals := map[string]Reason{
		"alreadyexists": ReasonAlreadyExists,
		"invalid":       ReasonInvalid,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := Reason(input)
	return &out, nil
}

type State string

const (
	StateCreating    State = "Creating"
	StateDeleted     State = "Deleted"
	StateDeleting    State = "Deleting"
	StateMigrated    State = "Migrated"
	StateRunning     State = "Running"
	StateStarting    State = "Starting"
	StateStopped     State = "Stopped"
	StateStopping    State = "Stopping"
	StateUnavailable State = "Unavailable"
	StateUpdating    State = "Updating"
)

func PossibleValuesForState() []string {
	return []string{
		string(StateCreating),
		string(StateDeleted),
		string(StateDeleting),
		string(StateMigrated),
		string(StateRunning),
		string(StateStarting),
		string(StateStopped),
		string(StateStopping),
		string(StateUnavailable),
		string(StateUpdating),
	}
}

func (s *State) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseState(input string) (*State, error) {
	vals := map[string]State{
		"creating":    StateCreating,
		"deleted":     StateDeleted,
		"deleting":    StateDeleting,
		"migrated":    StateMigrated,
		"running":     StateRunning,
		"starting":    StateStarting,
		"stopped":     StateStopped,
		"stopping":    StateStopping,
		"unavailable": StateUnavailable,
		"updating":    StateUpdating,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := State(input)
	return &out, nil
}

type VnetState string

const (
	VnetStateDisabled VnetState = "Disabled"
	VnetStateEnabled  VnetState = "Enabled"
)

func PossibleValuesForVnetState() []string {
	return []string{
		string(VnetStateDisabled),
		string(VnetStateEnabled),
	}
}

func (s *VnetState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseVnetState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseVnetState(input string) (*VnetState, error) {
	vals := map[string]VnetState{
		"disabled": VnetStateDisabled,
		"enabled":  VnetStateEnabled,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := VnetState(input)
	return &out, nil
}
