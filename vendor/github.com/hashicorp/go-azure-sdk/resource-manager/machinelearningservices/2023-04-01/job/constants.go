package job

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type BlockedTransformers string

const (
	BlockedTransformersCatTargetEncoder  BlockedTransformers = "CatTargetEncoder"
	BlockedTransformersCountVectorizer   BlockedTransformers = "CountVectorizer"
	BlockedTransformersHashOneHotEncoder BlockedTransformers = "HashOneHotEncoder"
	BlockedTransformersLabelEncoder      BlockedTransformers = "LabelEncoder"
	BlockedTransformersNaiveBayes        BlockedTransformers = "NaiveBayes"
	BlockedTransformersOneHotEncoder     BlockedTransformers = "OneHotEncoder"
	BlockedTransformersTextTargetEncoder BlockedTransformers = "TextTargetEncoder"
	BlockedTransformersTfIdf             BlockedTransformers = "TfIdf"
	BlockedTransformersWoETargetEncoder  BlockedTransformers = "WoETargetEncoder"
	BlockedTransformersWordEmbedding     BlockedTransformers = "WordEmbedding"
)

func PossibleValuesForBlockedTransformers() []string {
	return []string{
		string(BlockedTransformersCatTargetEncoder),
		string(BlockedTransformersCountVectorizer),
		string(BlockedTransformersHashOneHotEncoder),
		string(BlockedTransformersLabelEncoder),
		string(BlockedTransformersNaiveBayes),
		string(BlockedTransformersOneHotEncoder),
		string(BlockedTransformersTextTargetEncoder),
		string(BlockedTransformersTfIdf),
		string(BlockedTransformersWoETargetEncoder),
		string(BlockedTransformersWordEmbedding),
	}
}

func (s *BlockedTransformers) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseBlockedTransformers(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseBlockedTransformers(input string) (*BlockedTransformers, error) {
	vals := map[string]BlockedTransformers{
		"cattargetencoder":  BlockedTransformersCatTargetEncoder,
		"countvectorizer":   BlockedTransformersCountVectorizer,
		"hashonehotencoder": BlockedTransformersHashOneHotEncoder,
		"labelencoder":      BlockedTransformersLabelEncoder,
		"naivebayes":        BlockedTransformersNaiveBayes,
		"onehotencoder":     BlockedTransformersOneHotEncoder,
		"texttargetencoder": BlockedTransformersTextTargetEncoder,
		"tfidf":             BlockedTransformersTfIdf,
		"woetargetencoder":  BlockedTransformersWoETargetEncoder,
		"wordembedding":     BlockedTransformersWordEmbedding,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := BlockedTransformers(input)
	return &out, nil
}

type ClassificationModels string

const (
	ClassificationModelsBernoulliNaiveBayes   ClassificationModels = "BernoulliNaiveBayes"
	ClassificationModelsDecisionTree          ClassificationModels = "DecisionTree"
	ClassificationModelsExtremeRandomTrees    ClassificationModels = "ExtremeRandomTrees"
	ClassificationModelsGradientBoosting      ClassificationModels = "GradientBoosting"
	ClassificationModelsKNN                   ClassificationModels = "KNN"
	ClassificationModelsLightGBM              ClassificationModels = "LightGBM"
	ClassificationModelsLinearSVM             ClassificationModels = "LinearSVM"
	ClassificationModelsLogisticRegression    ClassificationModels = "LogisticRegression"
	ClassificationModelsMultinomialNaiveBayes ClassificationModels = "MultinomialNaiveBayes"
	ClassificationModelsRandomForest          ClassificationModels = "RandomForest"
	ClassificationModelsSGD                   ClassificationModels = "SGD"
	ClassificationModelsSVM                   ClassificationModels = "SVM"
	ClassificationModelsXGBoostClassifier     ClassificationModels = "XGBoostClassifier"
)

func PossibleValuesForClassificationModels() []string {
	return []string{
		string(ClassificationModelsBernoulliNaiveBayes),
		string(ClassificationModelsDecisionTree),
		string(ClassificationModelsExtremeRandomTrees),
		string(ClassificationModelsGradientBoosting),
		string(ClassificationModelsKNN),
		string(ClassificationModelsLightGBM),
		string(ClassificationModelsLinearSVM),
		string(ClassificationModelsLogisticRegression),
		string(ClassificationModelsMultinomialNaiveBayes),
		string(ClassificationModelsRandomForest),
		string(ClassificationModelsSGD),
		string(ClassificationModelsSVM),
		string(ClassificationModelsXGBoostClassifier),
	}
}

func (s *ClassificationModels) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseClassificationModels(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseClassificationModels(input string) (*ClassificationModels, error) {
	vals := map[string]ClassificationModels{
		"bernoullinaivebayes":   ClassificationModelsBernoulliNaiveBayes,
		"decisiontree":          ClassificationModelsDecisionTree,
		"extremerandomtrees":    ClassificationModelsExtremeRandomTrees,
		"gradientboosting":      ClassificationModelsGradientBoosting,
		"knn":                   ClassificationModelsKNN,
		"lightgbm":              ClassificationModelsLightGBM,
		"linearsvm":             ClassificationModelsLinearSVM,
		"logisticregression":    ClassificationModelsLogisticRegression,
		"multinomialnaivebayes": ClassificationModelsMultinomialNaiveBayes,
		"randomforest":          ClassificationModelsRandomForest,
		"sgd":                   ClassificationModelsSGD,
		"svm":                   ClassificationModelsSVM,
		"xgboostclassifier":     ClassificationModelsXGBoostClassifier,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ClassificationModels(input)
	return &out, nil
}

type ClassificationMultilabelPrimaryMetrics string

const (
	ClassificationMultilabelPrimaryMetricsAUCWeighted                   ClassificationMultilabelPrimaryMetrics = "AUCWeighted"
	ClassificationMultilabelPrimaryMetricsAccuracy                      ClassificationMultilabelPrimaryMetrics = "Accuracy"
	ClassificationMultilabelPrimaryMetricsAveragePrecisionScoreWeighted ClassificationMultilabelPrimaryMetrics = "AveragePrecisionScoreWeighted"
	ClassificationMultilabelPrimaryMetricsIOU                           ClassificationMultilabelPrimaryMetrics = "IOU"
	ClassificationMultilabelPrimaryMetricsNormMacroRecall               ClassificationMultilabelPrimaryMetrics = "NormMacroRecall"
	ClassificationMultilabelPrimaryMetricsPrecisionScoreWeighted        ClassificationMultilabelPrimaryMetrics = "PrecisionScoreWeighted"
)

func PossibleValuesForClassificationMultilabelPrimaryMetrics() []string {
	return []string{
		string(ClassificationMultilabelPrimaryMetricsAUCWeighted),
		string(ClassificationMultilabelPrimaryMetricsAccuracy),
		string(ClassificationMultilabelPrimaryMetricsAveragePrecisionScoreWeighted),
		string(ClassificationMultilabelPrimaryMetricsIOU),
		string(ClassificationMultilabelPrimaryMetricsNormMacroRecall),
		string(ClassificationMultilabelPrimaryMetricsPrecisionScoreWeighted),
	}
}

func (s *ClassificationMultilabelPrimaryMetrics) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseClassificationMultilabelPrimaryMetrics(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseClassificationMultilabelPrimaryMetrics(input string) (*ClassificationMultilabelPrimaryMetrics, error) {
	vals := map[string]ClassificationMultilabelPrimaryMetrics{
		"aucweighted":                   ClassificationMultilabelPrimaryMetricsAUCWeighted,
		"accuracy":                      ClassificationMultilabelPrimaryMetricsAccuracy,
		"averageprecisionscoreweighted": ClassificationMultilabelPrimaryMetricsAveragePrecisionScoreWeighted,
		"iou":                           ClassificationMultilabelPrimaryMetricsIOU,
		"normmacrorecall":               ClassificationMultilabelPrimaryMetricsNormMacroRecall,
		"precisionscoreweighted":        ClassificationMultilabelPrimaryMetricsPrecisionScoreWeighted,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ClassificationMultilabelPrimaryMetrics(input)
	return &out, nil
}

type ClassificationPrimaryMetrics string

const (
	ClassificationPrimaryMetricsAUCWeighted                   ClassificationPrimaryMetrics = "AUCWeighted"
	ClassificationPrimaryMetricsAccuracy                      ClassificationPrimaryMetrics = "Accuracy"
	ClassificationPrimaryMetricsAveragePrecisionScoreWeighted ClassificationPrimaryMetrics = "AveragePrecisionScoreWeighted"
	ClassificationPrimaryMetricsNormMacroRecall               ClassificationPrimaryMetrics = "NormMacroRecall"
	ClassificationPrimaryMetricsPrecisionScoreWeighted        ClassificationPrimaryMetrics = "PrecisionScoreWeighted"
)

func PossibleValuesForClassificationPrimaryMetrics() []string {
	return []string{
		string(ClassificationPrimaryMetricsAUCWeighted),
		string(ClassificationPrimaryMetricsAccuracy),
		string(ClassificationPrimaryMetricsAveragePrecisionScoreWeighted),
		string(ClassificationPrimaryMetricsNormMacroRecall),
		string(ClassificationPrimaryMetricsPrecisionScoreWeighted),
	}
}

func (s *ClassificationPrimaryMetrics) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseClassificationPrimaryMetrics(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseClassificationPrimaryMetrics(input string) (*ClassificationPrimaryMetrics, error) {
	vals := map[string]ClassificationPrimaryMetrics{
		"aucweighted":                   ClassificationPrimaryMetricsAUCWeighted,
		"accuracy":                      ClassificationPrimaryMetricsAccuracy,
		"averageprecisionscoreweighted": ClassificationPrimaryMetricsAveragePrecisionScoreWeighted,
		"normmacrorecall":               ClassificationPrimaryMetricsNormMacroRecall,
		"precisionscoreweighted":        ClassificationPrimaryMetricsPrecisionScoreWeighted,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ClassificationPrimaryMetrics(input)
	return &out, nil
}

type DistributionType string

const (
	DistributionTypeMpi        DistributionType = "Mpi"
	DistributionTypePyTorch    DistributionType = "PyTorch"
	DistributionTypeTensorFlow DistributionType = "TensorFlow"
)

func PossibleValuesForDistributionType() []string {
	return []string{
		string(DistributionTypeMpi),
		string(DistributionTypePyTorch),
		string(DistributionTypeTensorFlow),
	}
}

func (s *DistributionType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseDistributionType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseDistributionType(input string) (*DistributionType, error) {
	vals := map[string]DistributionType{
		"mpi":        DistributionTypeMpi,
		"pytorch":    DistributionTypePyTorch,
		"tensorflow": DistributionTypeTensorFlow,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := DistributionType(input)
	return &out, nil
}

type EarlyTerminationPolicyType string

const (
	EarlyTerminationPolicyTypeBandit              EarlyTerminationPolicyType = "Bandit"
	EarlyTerminationPolicyTypeMedianStopping      EarlyTerminationPolicyType = "MedianStopping"
	EarlyTerminationPolicyTypeTruncationSelection EarlyTerminationPolicyType = "TruncationSelection"
)

func PossibleValuesForEarlyTerminationPolicyType() []string {
	return []string{
		string(EarlyTerminationPolicyTypeBandit),
		string(EarlyTerminationPolicyTypeMedianStopping),
		string(EarlyTerminationPolicyTypeTruncationSelection),
	}
}

func (s *EarlyTerminationPolicyType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseEarlyTerminationPolicyType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseEarlyTerminationPolicyType(input string) (*EarlyTerminationPolicyType, error) {
	vals := map[string]EarlyTerminationPolicyType{
		"bandit":              EarlyTerminationPolicyTypeBandit,
		"medianstopping":      EarlyTerminationPolicyTypeMedianStopping,
		"truncationselection": EarlyTerminationPolicyTypeTruncationSelection,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := EarlyTerminationPolicyType(input)
	return &out, nil
}

type FeatureLags string

const (
	FeatureLagsAuto FeatureLags = "Auto"
	FeatureLagsNone FeatureLags = "None"
)

func PossibleValuesForFeatureLags() []string {
	return []string{
		string(FeatureLagsAuto),
		string(FeatureLagsNone),
	}
}

func (s *FeatureLags) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseFeatureLags(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseFeatureLags(input string) (*FeatureLags, error) {
	vals := map[string]FeatureLags{
		"auto": FeatureLagsAuto,
		"none": FeatureLagsNone,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := FeatureLags(input)
	return &out, nil
}

type FeaturizationMode string

const (
	FeaturizationModeAuto   FeaturizationMode = "Auto"
	FeaturizationModeCustom FeaturizationMode = "Custom"
	FeaturizationModeOff    FeaturizationMode = "Off"
)

func PossibleValuesForFeaturizationMode() []string {
	return []string{
		string(FeaturizationModeAuto),
		string(FeaturizationModeCustom),
		string(FeaturizationModeOff),
	}
}

func (s *FeaturizationMode) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseFeaturizationMode(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseFeaturizationMode(input string) (*FeaturizationMode, error) {
	vals := map[string]FeaturizationMode{
		"auto":   FeaturizationModeAuto,
		"custom": FeaturizationModeCustom,
		"off":    FeaturizationModeOff,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := FeaturizationMode(input)
	return &out, nil
}

type ForecastHorizonMode string

const (
	ForecastHorizonModeAuto   ForecastHorizonMode = "Auto"
	ForecastHorizonModeCustom ForecastHorizonMode = "Custom"
)

func PossibleValuesForForecastHorizonMode() []string {
	return []string{
		string(ForecastHorizonModeAuto),
		string(ForecastHorizonModeCustom),
	}
}

func (s *ForecastHorizonMode) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseForecastHorizonMode(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseForecastHorizonMode(input string) (*ForecastHorizonMode, error) {
	vals := map[string]ForecastHorizonMode{
		"auto":   ForecastHorizonModeAuto,
		"custom": ForecastHorizonModeCustom,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ForecastHorizonMode(input)
	return &out, nil
}

type ForecastingModels string

const (
	ForecastingModelsArimax               ForecastingModels = "Arimax"
	ForecastingModelsAutoArima            ForecastingModels = "AutoArima"
	ForecastingModelsAverage              ForecastingModels = "Average"
	ForecastingModelsDecisionTree         ForecastingModels = "DecisionTree"
	ForecastingModelsElasticNet           ForecastingModels = "ElasticNet"
	ForecastingModelsExponentialSmoothing ForecastingModels = "ExponentialSmoothing"
	ForecastingModelsExtremeRandomTrees   ForecastingModels = "ExtremeRandomTrees"
	ForecastingModelsGradientBoosting     ForecastingModels = "GradientBoosting"
	ForecastingModelsKNN                  ForecastingModels = "KNN"
	ForecastingModelsLassoLars            ForecastingModels = "LassoLars"
	ForecastingModelsLightGBM             ForecastingModels = "LightGBM"
	ForecastingModelsNaive                ForecastingModels = "Naive"
	ForecastingModelsProphet              ForecastingModels = "Prophet"
	ForecastingModelsRandomForest         ForecastingModels = "RandomForest"
	ForecastingModelsSGD                  ForecastingModels = "SGD"
	ForecastingModelsSeasonalAverage      ForecastingModels = "SeasonalAverage"
	ForecastingModelsSeasonalNaive        ForecastingModels = "SeasonalNaive"
	ForecastingModelsTCNForecaster        ForecastingModels = "TCNForecaster"
	ForecastingModelsXGBoostRegressor     ForecastingModels = "XGBoostRegressor"
)

func PossibleValuesForForecastingModels() []string {
	return []string{
		string(ForecastingModelsArimax),
		string(ForecastingModelsAutoArima),
		string(ForecastingModelsAverage),
		string(ForecastingModelsDecisionTree),
		string(ForecastingModelsElasticNet),
		string(ForecastingModelsExponentialSmoothing),
		string(ForecastingModelsExtremeRandomTrees),
		string(ForecastingModelsGradientBoosting),
		string(ForecastingModelsKNN),
		string(ForecastingModelsLassoLars),
		string(ForecastingModelsLightGBM),
		string(ForecastingModelsNaive),
		string(ForecastingModelsProphet),
		string(ForecastingModelsRandomForest),
		string(ForecastingModelsSGD),
		string(ForecastingModelsSeasonalAverage),
		string(ForecastingModelsSeasonalNaive),
		string(ForecastingModelsTCNForecaster),
		string(ForecastingModelsXGBoostRegressor),
	}
}

func (s *ForecastingModels) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseForecastingModels(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseForecastingModels(input string) (*ForecastingModels, error) {
	vals := map[string]ForecastingModels{
		"arimax":               ForecastingModelsArimax,
		"autoarima":            ForecastingModelsAutoArima,
		"average":              ForecastingModelsAverage,
		"decisiontree":         ForecastingModelsDecisionTree,
		"elasticnet":           ForecastingModelsElasticNet,
		"exponentialsmoothing": ForecastingModelsExponentialSmoothing,
		"extremerandomtrees":   ForecastingModelsExtremeRandomTrees,
		"gradientboosting":     ForecastingModelsGradientBoosting,
		"knn":                  ForecastingModelsKNN,
		"lassolars":            ForecastingModelsLassoLars,
		"lightgbm":             ForecastingModelsLightGBM,
		"naive":                ForecastingModelsNaive,
		"prophet":              ForecastingModelsProphet,
		"randomforest":         ForecastingModelsRandomForest,
		"sgd":                  ForecastingModelsSGD,
		"seasonalaverage":      ForecastingModelsSeasonalAverage,
		"seasonalnaive":        ForecastingModelsSeasonalNaive,
		"tcnforecaster":        ForecastingModelsTCNForecaster,
		"xgboostregressor":     ForecastingModelsXGBoostRegressor,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ForecastingModels(input)
	return &out, nil
}

type ForecastingPrimaryMetrics string

const (
	ForecastingPrimaryMetricsNormalizedMeanAbsoluteError    ForecastingPrimaryMetrics = "NormalizedMeanAbsoluteError"
	ForecastingPrimaryMetricsNormalizedRootMeanSquaredError ForecastingPrimaryMetrics = "NormalizedRootMeanSquaredError"
	ForecastingPrimaryMetricsRTwoScore                      ForecastingPrimaryMetrics = "R2Score"
	ForecastingPrimaryMetricsSpearmanCorrelation            ForecastingPrimaryMetrics = "SpearmanCorrelation"
)

func PossibleValuesForForecastingPrimaryMetrics() []string {
	return []string{
		string(ForecastingPrimaryMetricsNormalizedMeanAbsoluteError),
		string(ForecastingPrimaryMetricsNormalizedRootMeanSquaredError),
		string(ForecastingPrimaryMetricsRTwoScore),
		string(ForecastingPrimaryMetricsSpearmanCorrelation),
	}
}

func (s *ForecastingPrimaryMetrics) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseForecastingPrimaryMetrics(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseForecastingPrimaryMetrics(input string) (*ForecastingPrimaryMetrics, error) {
	vals := map[string]ForecastingPrimaryMetrics{
		"normalizedmeanabsoluteerror":    ForecastingPrimaryMetricsNormalizedMeanAbsoluteError,
		"normalizedrootmeansquarederror": ForecastingPrimaryMetricsNormalizedRootMeanSquaredError,
		"r2score":                        ForecastingPrimaryMetricsRTwoScore,
		"spearmancorrelation":            ForecastingPrimaryMetricsSpearmanCorrelation,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ForecastingPrimaryMetrics(input)
	return &out, nil
}

type Goal string

const (
	GoalMaximize Goal = "Maximize"
	GoalMinimize Goal = "Minimize"
)

func PossibleValuesForGoal() []string {
	return []string{
		string(GoalMaximize),
		string(GoalMinimize),
	}
}

func (s *Goal) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseGoal(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseGoal(input string) (*Goal, error) {
	vals := map[string]Goal{
		"maximize": GoalMaximize,
		"minimize": GoalMinimize,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := Goal(input)
	return &out, nil
}

type IdentityConfigurationType string

const (
	IdentityConfigurationTypeAMLToken     IdentityConfigurationType = "AMLToken"
	IdentityConfigurationTypeManaged      IdentityConfigurationType = "Managed"
	IdentityConfigurationTypeUserIdentity IdentityConfigurationType = "UserIdentity"
)

func PossibleValuesForIdentityConfigurationType() []string {
	return []string{
		string(IdentityConfigurationTypeAMLToken),
		string(IdentityConfigurationTypeManaged),
		string(IdentityConfigurationTypeUserIdentity),
	}
}

func (s *IdentityConfigurationType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseIdentityConfigurationType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseIdentityConfigurationType(input string) (*IdentityConfigurationType, error) {
	vals := map[string]IdentityConfigurationType{
		"amltoken":     IdentityConfigurationTypeAMLToken,
		"managed":      IdentityConfigurationTypeManaged,
		"useridentity": IdentityConfigurationTypeUserIdentity,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := IdentityConfigurationType(input)
	return &out, nil
}

type InputDeliveryMode string

const (
	InputDeliveryModeDirect         InputDeliveryMode = "Direct"
	InputDeliveryModeDownload       InputDeliveryMode = "Download"
	InputDeliveryModeEvalDownload   InputDeliveryMode = "EvalDownload"
	InputDeliveryModeEvalMount      InputDeliveryMode = "EvalMount"
	InputDeliveryModeReadOnlyMount  InputDeliveryMode = "ReadOnlyMount"
	InputDeliveryModeReadWriteMount InputDeliveryMode = "ReadWriteMount"
)

func PossibleValuesForInputDeliveryMode() []string {
	return []string{
		string(InputDeliveryModeDirect),
		string(InputDeliveryModeDownload),
		string(InputDeliveryModeEvalDownload),
		string(InputDeliveryModeEvalMount),
		string(InputDeliveryModeReadOnlyMount),
		string(InputDeliveryModeReadWriteMount),
	}
}

func (s *InputDeliveryMode) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseInputDeliveryMode(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseInputDeliveryMode(input string) (*InputDeliveryMode, error) {
	vals := map[string]InputDeliveryMode{
		"direct":         InputDeliveryModeDirect,
		"download":       InputDeliveryModeDownload,
		"evaldownload":   InputDeliveryModeEvalDownload,
		"evalmount":      InputDeliveryModeEvalMount,
		"readonlymount":  InputDeliveryModeReadOnlyMount,
		"readwritemount": InputDeliveryModeReadWriteMount,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := InputDeliveryMode(input)
	return &out, nil
}

type InstanceSegmentationPrimaryMetrics string

const (
	InstanceSegmentationPrimaryMetricsMeanAveragePrecision InstanceSegmentationPrimaryMetrics = "MeanAveragePrecision"
)

func PossibleValuesForInstanceSegmentationPrimaryMetrics() []string {
	return []string{
		string(InstanceSegmentationPrimaryMetricsMeanAveragePrecision),
	}
}

func (s *InstanceSegmentationPrimaryMetrics) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseInstanceSegmentationPrimaryMetrics(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseInstanceSegmentationPrimaryMetrics(input string) (*InstanceSegmentationPrimaryMetrics, error) {
	vals := map[string]InstanceSegmentationPrimaryMetrics{
		"meanaverageprecision": InstanceSegmentationPrimaryMetricsMeanAveragePrecision,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := InstanceSegmentationPrimaryMetrics(input)
	return &out, nil
}

type JobInputType string

const (
	JobInputTypeCustomModel JobInputType = "custom_model"
	JobInputTypeLiteral     JobInputType = "literal"
	JobInputTypeMlflowModel JobInputType = "mlflow_model"
	JobInputTypeMltable     JobInputType = "mltable"
	JobInputTypeTritonModel JobInputType = "triton_model"
	JobInputTypeUriFile     JobInputType = "uri_file"
	JobInputTypeUriFolder   JobInputType = "uri_folder"
)

func PossibleValuesForJobInputType() []string {
	return []string{
		string(JobInputTypeCustomModel),
		string(JobInputTypeLiteral),
		string(JobInputTypeMlflowModel),
		string(JobInputTypeMltable),
		string(JobInputTypeTritonModel),
		string(JobInputTypeUriFile),
		string(JobInputTypeUriFolder),
	}
}

func (s *JobInputType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseJobInputType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseJobInputType(input string) (*JobInputType, error) {
	vals := map[string]JobInputType{
		"custom_model": JobInputTypeCustomModel,
		"literal":      JobInputTypeLiteral,
		"mlflow_model": JobInputTypeMlflowModel,
		"mltable":      JobInputTypeMltable,
		"triton_model": JobInputTypeTritonModel,
		"uri_file":     JobInputTypeUriFile,
		"uri_folder":   JobInputTypeUriFolder,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := JobInputType(input)
	return &out, nil
}

type JobLimitsType string

const (
	JobLimitsTypeCommand JobLimitsType = "Command"
	JobLimitsTypeSweep   JobLimitsType = "Sweep"
)

func PossibleValuesForJobLimitsType() []string {
	return []string{
		string(JobLimitsTypeCommand),
		string(JobLimitsTypeSweep),
	}
}

func (s *JobLimitsType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseJobLimitsType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseJobLimitsType(input string) (*JobLimitsType, error) {
	vals := map[string]JobLimitsType{
		"command": JobLimitsTypeCommand,
		"sweep":   JobLimitsTypeSweep,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := JobLimitsType(input)
	return &out, nil
}

type JobOutputType string

const (
	JobOutputTypeCustomModel JobOutputType = "custom_model"
	JobOutputTypeMlflowModel JobOutputType = "mlflow_model"
	JobOutputTypeMltable     JobOutputType = "mltable"
	JobOutputTypeTritonModel JobOutputType = "triton_model"
	JobOutputTypeUriFile     JobOutputType = "uri_file"
	JobOutputTypeUriFolder   JobOutputType = "uri_folder"
)

func PossibleValuesForJobOutputType() []string {
	return []string{
		string(JobOutputTypeCustomModel),
		string(JobOutputTypeMlflowModel),
		string(JobOutputTypeMltable),
		string(JobOutputTypeTritonModel),
		string(JobOutputTypeUriFile),
		string(JobOutputTypeUriFolder),
	}
}

func (s *JobOutputType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseJobOutputType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseJobOutputType(input string) (*JobOutputType, error) {
	vals := map[string]JobOutputType{
		"custom_model": JobOutputTypeCustomModel,
		"mlflow_model": JobOutputTypeMlflowModel,
		"mltable":      JobOutputTypeMltable,
		"triton_model": JobOutputTypeTritonModel,
		"uri_file":     JobOutputTypeUriFile,
		"uri_folder":   JobOutputTypeUriFolder,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := JobOutputType(input)
	return &out, nil
}

type JobStatus string

const (
	JobStatusCancelRequested JobStatus = "CancelRequested"
	JobStatusCanceled        JobStatus = "Canceled"
	JobStatusCompleted       JobStatus = "Completed"
	JobStatusFailed          JobStatus = "Failed"
	JobStatusFinalizing      JobStatus = "Finalizing"
	JobStatusNotResponding   JobStatus = "NotResponding"
	JobStatusNotStarted      JobStatus = "NotStarted"
	JobStatusPaused          JobStatus = "Paused"
	JobStatusPreparing       JobStatus = "Preparing"
	JobStatusProvisioning    JobStatus = "Provisioning"
	JobStatusQueued          JobStatus = "Queued"
	JobStatusRunning         JobStatus = "Running"
	JobStatusStarting        JobStatus = "Starting"
	JobStatusUnknown         JobStatus = "Unknown"
)

func PossibleValuesForJobStatus() []string {
	return []string{
		string(JobStatusCancelRequested),
		string(JobStatusCanceled),
		string(JobStatusCompleted),
		string(JobStatusFailed),
		string(JobStatusFinalizing),
		string(JobStatusNotResponding),
		string(JobStatusNotStarted),
		string(JobStatusPaused),
		string(JobStatusPreparing),
		string(JobStatusProvisioning),
		string(JobStatusQueued),
		string(JobStatusRunning),
		string(JobStatusStarting),
		string(JobStatusUnknown),
	}
}

func (s *JobStatus) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseJobStatus(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseJobStatus(input string) (*JobStatus, error) {
	vals := map[string]JobStatus{
		"cancelrequested": JobStatusCancelRequested,
		"canceled":        JobStatusCanceled,
		"completed":       JobStatusCompleted,
		"failed":          JobStatusFailed,
		"finalizing":      JobStatusFinalizing,
		"notresponding":   JobStatusNotResponding,
		"notstarted":      JobStatusNotStarted,
		"paused":          JobStatusPaused,
		"preparing":       JobStatusPreparing,
		"provisioning":    JobStatusProvisioning,
		"queued":          JobStatusQueued,
		"running":         JobStatusRunning,
		"starting":        JobStatusStarting,
		"unknown":         JobStatusUnknown,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := JobStatus(input)
	return &out, nil
}

type JobType string

const (
	JobTypeAutoML   JobType = "AutoML"
	JobTypeCommand  JobType = "Command"
	JobTypePipeline JobType = "Pipeline"
	JobTypeSweep    JobType = "Sweep"
)

func PossibleValuesForJobType() []string {
	return []string{
		string(JobTypeAutoML),
		string(JobTypeCommand),
		string(JobTypePipeline),
		string(JobTypeSweep),
	}
}

func (s *JobType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseJobType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseJobType(input string) (*JobType, error) {
	vals := map[string]JobType{
		"automl":   JobTypeAutoML,
		"command":  JobTypeCommand,
		"pipeline": JobTypePipeline,
		"sweep":    JobTypeSweep,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := JobType(input)
	return &out, nil
}

type LearningRateScheduler string

const (
	LearningRateSchedulerNone         LearningRateScheduler = "None"
	LearningRateSchedulerStep         LearningRateScheduler = "Step"
	LearningRateSchedulerWarmupCosine LearningRateScheduler = "WarmupCosine"
)

func PossibleValuesForLearningRateScheduler() []string {
	return []string{
		string(LearningRateSchedulerNone),
		string(LearningRateSchedulerStep),
		string(LearningRateSchedulerWarmupCosine),
	}
}

func (s *LearningRateScheduler) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseLearningRateScheduler(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseLearningRateScheduler(input string) (*LearningRateScheduler, error) {
	vals := map[string]LearningRateScheduler{
		"none":         LearningRateSchedulerNone,
		"step":         LearningRateSchedulerStep,
		"warmupcosine": LearningRateSchedulerWarmupCosine,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := LearningRateScheduler(input)
	return &out, nil
}

type ListViewType string

const (
	ListViewTypeActiveOnly   ListViewType = "ActiveOnly"
	ListViewTypeAll          ListViewType = "All"
	ListViewTypeArchivedOnly ListViewType = "ArchivedOnly"
)

func PossibleValuesForListViewType() []string {
	return []string{
		string(ListViewTypeActiveOnly),
		string(ListViewTypeAll),
		string(ListViewTypeArchivedOnly),
	}
}

func (s *ListViewType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseListViewType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseListViewType(input string) (*ListViewType, error) {
	vals := map[string]ListViewType{
		"activeonly":   ListViewTypeActiveOnly,
		"all":          ListViewTypeAll,
		"archivedonly": ListViewTypeArchivedOnly,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ListViewType(input)
	return &out, nil
}

type LogVerbosity string

const (
	LogVerbosityCritical LogVerbosity = "Critical"
	LogVerbosityDebug    LogVerbosity = "Debug"
	LogVerbosityError    LogVerbosity = "Error"
	LogVerbosityInfo     LogVerbosity = "Info"
	LogVerbosityNotSet   LogVerbosity = "NotSet"
	LogVerbosityWarning  LogVerbosity = "Warning"
)

func PossibleValuesForLogVerbosity() []string {
	return []string{
		string(LogVerbosityCritical),
		string(LogVerbosityDebug),
		string(LogVerbosityError),
		string(LogVerbosityInfo),
		string(LogVerbosityNotSet),
		string(LogVerbosityWarning),
	}
}

func (s *LogVerbosity) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseLogVerbosity(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseLogVerbosity(input string) (*LogVerbosity, error) {
	vals := map[string]LogVerbosity{
		"critical": LogVerbosityCritical,
		"debug":    LogVerbosityDebug,
		"error":    LogVerbosityError,
		"info":     LogVerbosityInfo,
		"notset":   LogVerbosityNotSet,
		"warning":  LogVerbosityWarning,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := LogVerbosity(input)
	return &out, nil
}

type ModelSize string

const (
	ModelSizeExtraLarge ModelSize = "ExtraLarge"
	ModelSizeLarge      ModelSize = "Large"
	ModelSizeMedium     ModelSize = "Medium"
	ModelSizeNone       ModelSize = "None"
	ModelSizeSmall      ModelSize = "Small"
)

func PossibleValuesForModelSize() []string {
	return []string{
		string(ModelSizeExtraLarge),
		string(ModelSizeLarge),
		string(ModelSizeMedium),
		string(ModelSizeNone),
		string(ModelSizeSmall),
	}
}

func (s *ModelSize) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseModelSize(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseModelSize(input string) (*ModelSize, error) {
	vals := map[string]ModelSize{
		"extralarge": ModelSizeExtraLarge,
		"large":      ModelSizeLarge,
		"medium":     ModelSizeMedium,
		"none":       ModelSizeNone,
		"small":      ModelSizeSmall,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ModelSize(input)
	return &out, nil
}

type NCrossValidationsMode string

const (
	NCrossValidationsModeAuto   NCrossValidationsMode = "Auto"
	NCrossValidationsModeCustom NCrossValidationsMode = "Custom"
)

func PossibleValuesForNCrossValidationsMode() []string {
	return []string{
		string(NCrossValidationsModeAuto),
		string(NCrossValidationsModeCustom),
	}
}

func (s *NCrossValidationsMode) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseNCrossValidationsMode(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseNCrossValidationsMode(input string) (*NCrossValidationsMode, error) {
	vals := map[string]NCrossValidationsMode{
		"auto":   NCrossValidationsModeAuto,
		"custom": NCrossValidationsModeCustom,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := NCrossValidationsMode(input)
	return &out, nil
}

type NodesValueType string

const (
	NodesValueTypeAll NodesValueType = "All"
)

func PossibleValuesForNodesValueType() []string {
	return []string{
		string(NodesValueTypeAll),
	}
}

func (s *NodesValueType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseNodesValueType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseNodesValueType(input string) (*NodesValueType, error) {
	vals := map[string]NodesValueType{
		"all": NodesValueTypeAll,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := NodesValueType(input)
	return &out, nil
}

type ObjectDetectionPrimaryMetrics string

const (
	ObjectDetectionPrimaryMetricsMeanAveragePrecision ObjectDetectionPrimaryMetrics = "MeanAveragePrecision"
)

func PossibleValuesForObjectDetectionPrimaryMetrics() []string {
	return []string{
		string(ObjectDetectionPrimaryMetricsMeanAveragePrecision),
	}
}

func (s *ObjectDetectionPrimaryMetrics) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseObjectDetectionPrimaryMetrics(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseObjectDetectionPrimaryMetrics(input string) (*ObjectDetectionPrimaryMetrics, error) {
	vals := map[string]ObjectDetectionPrimaryMetrics{
		"meanaverageprecision": ObjectDetectionPrimaryMetricsMeanAveragePrecision,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ObjectDetectionPrimaryMetrics(input)
	return &out, nil
}

type OutputDeliveryMode string

const (
	OutputDeliveryModeReadWriteMount OutputDeliveryMode = "ReadWriteMount"
	OutputDeliveryModeUpload         OutputDeliveryMode = "Upload"
)

func PossibleValuesForOutputDeliveryMode() []string {
	return []string{
		string(OutputDeliveryModeReadWriteMount),
		string(OutputDeliveryModeUpload),
	}
}

func (s *OutputDeliveryMode) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseOutputDeliveryMode(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseOutputDeliveryMode(input string) (*OutputDeliveryMode, error) {
	vals := map[string]OutputDeliveryMode{
		"readwritemount": OutputDeliveryModeReadWriteMount,
		"upload":         OutputDeliveryModeUpload,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := OutputDeliveryMode(input)
	return &out, nil
}

type RandomSamplingAlgorithmRule string

const (
	RandomSamplingAlgorithmRuleRandom RandomSamplingAlgorithmRule = "Random"
	RandomSamplingAlgorithmRuleSobol  RandomSamplingAlgorithmRule = "Sobol"
)

func PossibleValuesForRandomSamplingAlgorithmRule() []string {
	return []string{
		string(RandomSamplingAlgorithmRuleRandom),
		string(RandomSamplingAlgorithmRuleSobol),
	}
}

func (s *RandomSamplingAlgorithmRule) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseRandomSamplingAlgorithmRule(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseRandomSamplingAlgorithmRule(input string) (*RandomSamplingAlgorithmRule, error) {
	vals := map[string]RandomSamplingAlgorithmRule{
		"random": RandomSamplingAlgorithmRuleRandom,
		"sobol":  RandomSamplingAlgorithmRuleSobol,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := RandomSamplingAlgorithmRule(input)
	return &out, nil
}

type RegressionModels string

const (
	RegressionModelsDecisionTree       RegressionModels = "DecisionTree"
	RegressionModelsElasticNet         RegressionModels = "ElasticNet"
	RegressionModelsExtremeRandomTrees RegressionModels = "ExtremeRandomTrees"
	RegressionModelsGradientBoosting   RegressionModels = "GradientBoosting"
	RegressionModelsKNN                RegressionModels = "KNN"
	RegressionModelsLassoLars          RegressionModels = "LassoLars"
	RegressionModelsLightGBM           RegressionModels = "LightGBM"
	RegressionModelsRandomForest       RegressionModels = "RandomForest"
	RegressionModelsSGD                RegressionModels = "SGD"
	RegressionModelsXGBoostRegressor   RegressionModels = "XGBoostRegressor"
)

func PossibleValuesForRegressionModels() []string {
	return []string{
		string(RegressionModelsDecisionTree),
		string(RegressionModelsElasticNet),
		string(RegressionModelsExtremeRandomTrees),
		string(RegressionModelsGradientBoosting),
		string(RegressionModelsKNN),
		string(RegressionModelsLassoLars),
		string(RegressionModelsLightGBM),
		string(RegressionModelsRandomForest),
		string(RegressionModelsSGD),
		string(RegressionModelsXGBoostRegressor),
	}
}

func (s *RegressionModels) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseRegressionModels(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseRegressionModels(input string) (*RegressionModels, error) {
	vals := map[string]RegressionModels{
		"decisiontree":       RegressionModelsDecisionTree,
		"elasticnet":         RegressionModelsElasticNet,
		"extremerandomtrees": RegressionModelsExtremeRandomTrees,
		"gradientboosting":   RegressionModelsGradientBoosting,
		"knn":                RegressionModelsKNN,
		"lassolars":          RegressionModelsLassoLars,
		"lightgbm":           RegressionModelsLightGBM,
		"randomforest":       RegressionModelsRandomForest,
		"sgd":                RegressionModelsSGD,
		"xgboostregressor":   RegressionModelsXGBoostRegressor,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := RegressionModels(input)
	return &out, nil
}

type RegressionPrimaryMetrics string

const (
	RegressionPrimaryMetricsNormalizedMeanAbsoluteError    RegressionPrimaryMetrics = "NormalizedMeanAbsoluteError"
	RegressionPrimaryMetricsNormalizedRootMeanSquaredError RegressionPrimaryMetrics = "NormalizedRootMeanSquaredError"
	RegressionPrimaryMetricsRTwoScore                      RegressionPrimaryMetrics = "R2Score"
	RegressionPrimaryMetricsSpearmanCorrelation            RegressionPrimaryMetrics = "SpearmanCorrelation"
)

func PossibleValuesForRegressionPrimaryMetrics() []string {
	return []string{
		string(RegressionPrimaryMetricsNormalizedMeanAbsoluteError),
		string(RegressionPrimaryMetricsNormalizedRootMeanSquaredError),
		string(RegressionPrimaryMetricsRTwoScore),
		string(RegressionPrimaryMetricsSpearmanCorrelation),
	}
}

func (s *RegressionPrimaryMetrics) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseRegressionPrimaryMetrics(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseRegressionPrimaryMetrics(input string) (*RegressionPrimaryMetrics, error) {
	vals := map[string]RegressionPrimaryMetrics{
		"normalizedmeanabsoluteerror":    RegressionPrimaryMetricsNormalizedMeanAbsoluteError,
		"normalizedrootmeansquarederror": RegressionPrimaryMetricsNormalizedRootMeanSquaredError,
		"r2score":                        RegressionPrimaryMetricsRTwoScore,
		"spearmancorrelation":            RegressionPrimaryMetricsSpearmanCorrelation,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := RegressionPrimaryMetrics(input)
	return &out, nil
}

type SamplingAlgorithmType string

const (
	SamplingAlgorithmTypeBayesian SamplingAlgorithmType = "Bayesian"
	SamplingAlgorithmTypeGrid     SamplingAlgorithmType = "Grid"
	SamplingAlgorithmTypeRandom   SamplingAlgorithmType = "Random"
)

func PossibleValuesForSamplingAlgorithmType() []string {
	return []string{
		string(SamplingAlgorithmTypeBayesian),
		string(SamplingAlgorithmTypeGrid),
		string(SamplingAlgorithmTypeRandom),
	}
}

func (s *SamplingAlgorithmType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseSamplingAlgorithmType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseSamplingAlgorithmType(input string) (*SamplingAlgorithmType, error) {
	vals := map[string]SamplingAlgorithmType{
		"bayesian": SamplingAlgorithmTypeBayesian,
		"grid":     SamplingAlgorithmTypeGrid,
		"random":   SamplingAlgorithmTypeRandom,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := SamplingAlgorithmType(input)
	return &out, nil
}

type SeasonalityMode string

const (
	SeasonalityModeAuto   SeasonalityMode = "Auto"
	SeasonalityModeCustom SeasonalityMode = "Custom"
)

func PossibleValuesForSeasonalityMode() []string {
	return []string{
		string(SeasonalityModeAuto),
		string(SeasonalityModeCustom),
	}
}

func (s *SeasonalityMode) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseSeasonalityMode(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseSeasonalityMode(input string) (*SeasonalityMode, error) {
	vals := map[string]SeasonalityMode{
		"auto":   SeasonalityModeAuto,
		"custom": SeasonalityModeCustom,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := SeasonalityMode(input)
	return &out, nil
}

type ShortSeriesHandlingConfiguration string

const (
	ShortSeriesHandlingConfigurationAuto ShortSeriesHandlingConfiguration = "Auto"
	ShortSeriesHandlingConfigurationDrop ShortSeriesHandlingConfiguration = "Drop"
	ShortSeriesHandlingConfigurationNone ShortSeriesHandlingConfiguration = "None"
	ShortSeriesHandlingConfigurationPad  ShortSeriesHandlingConfiguration = "Pad"
)

func PossibleValuesForShortSeriesHandlingConfiguration() []string {
	return []string{
		string(ShortSeriesHandlingConfigurationAuto),
		string(ShortSeriesHandlingConfigurationDrop),
		string(ShortSeriesHandlingConfigurationNone),
		string(ShortSeriesHandlingConfigurationPad),
	}
}

func (s *ShortSeriesHandlingConfiguration) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseShortSeriesHandlingConfiguration(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseShortSeriesHandlingConfiguration(input string) (*ShortSeriesHandlingConfiguration, error) {
	vals := map[string]ShortSeriesHandlingConfiguration{
		"auto": ShortSeriesHandlingConfigurationAuto,
		"drop": ShortSeriesHandlingConfigurationDrop,
		"none": ShortSeriesHandlingConfigurationNone,
		"pad":  ShortSeriesHandlingConfigurationPad,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ShortSeriesHandlingConfiguration(input)
	return &out, nil
}

type StackMetaLearnerType string

const (
	StackMetaLearnerTypeElasticNet           StackMetaLearnerType = "ElasticNet"
	StackMetaLearnerTypeElasticNetCV         StackMetaLearnerType = "ElasticNetCV"
	StackMetaLearnerTypeLightGBMClassifier   StackMetaLearnerType = "LightGBMClassifier"
	StackMetaLearnerTypeLightGBMRegressor    StackMetaLearnerType = "LightGBMRegressor"
	StackMetaLearnerTypeLinearRegression     StackMetaLearnerType = "LinearRegression"
	StackMetaLearnerTypeLogisticRegression   StackMetaLearnerType = "LogisticRegression"
	StackMetaLearnerTypeLogisticRegressionCV StackMetaLearnerType = "LogisticRegressionCV"
	StackMetaLearnerTypeNone                 StackMetaLearnerType = "None"
)

func PossibleValuesForStackMetaLearnerType() []string {
	return []string{
		string(StackMetaLearnerTypeElasticNet),
		string(StackMetaLearnerTypeElasticNetCV),
		string(StackMetaLearnerTypeLightGBMClassifier),
		string(StackMetaLearnerTypeLightGBMRegressor),
		string(StackMetaLearnerTypeLinearRegression),
		string(StackMetaLearnerTypeLogisticRegression),
		string(StackMetaLearnerTypeLogisticRegressionCV),
		string(StackMetaLearnerTypeNone),
	}
}

func (s *StackMetaLearnerType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseStackMetaLearnerType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseStackMetaLearnerType(input string) (*StackMetaLearnerType, error) {
	vals := map[string]StackMetaLearnerType{
		"elasticnet":           StackMetaLearnerTypeElasticNet,
		"elasticnetcv":         StackMetaLearnerTypeElasticNetCV,
		"lightgbmclassifier":   StackMetaLearnerTypeLightGBMClassifier,
		"lightgbmregressor":    StackMetaLearnerTypeLightGBMRegressor,
		"linearregression":     StackMetaLearnerTypeLinearRegression,
		"logisticregression":   StackMetaLearnerTypeLogisticRegression,
		"logisticregressioncv": StackMetaLearnerTypeLogisticRegressionCV,
		"none":                 StackMetaLearnerTypeNone,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := StackMetaLearnerType(input)
	return &out, nil
}

type StochasticOptimizer string

const (
	StochasticOptimizerAdam  StochasticOptimizer = "Adam"
	StochasticOptimizerAdamw StochasticOptimizer = "Adamw"
	StochasticOptimizerNone  StochasticOptimizer = "None"
	StochasticOptimizerSgd   StochasticOptimizer = "Sgd"
)

func PossibleValuesForStochasticOptimizer() []string {
	return []string{
		string(StochasticOptimizerAdam),
		string(StochasticOptimizerAdamw),
		string(StochasticOptimizerNone),
		string(StochasticOptimizerSgd),
	}
}

func (s *StochasticOptimizer) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseStochasticOptimizer(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseStochasticOptimizer(input string) (*StochasticOptimizer, error) {
	vals := map[string]StochasticOptimizer{
		"adam":  StochasticOptimizerAdam,
		"adamw": StochasticOptimizerAdamw,
		"none":  StochasticOptimizerNone,
		"sgd":   StochasticOptimizerSgd,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := StochasticOptimizer(input)
	return &out, nil
}

type TargetAggregationFunction string

const (
	TargetAggregationFunctionMax  TargetAggregationFunction = "Max"
	TargetAggregationFunctionMean TargetAggregationFunction = "Mean"
	TargetAggregationFunctionMin  TargetAggregationFunction = "Min"
	TargetAggregationFunctionNone TargetAggregationFunction = "None"
	TargetAggregationFunctionSum  TargetAggregationFunction = "Sum"
)

func PossibleValuesForTargetAggregationFunction() []string {
	return []string{
		string(TargetAggregationFunctionMax),
		string(TargetAggregationFunctionMean),
		string(TargetAggregationFunctionMin),
		string(TargetAggregationFunctionNone),
		string(TargetAggregationFunctionSum),
	}
}

func (s *TargetAggregationFunction) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseTargetAggregationFunction(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseTargetAggregationFunction(input string) (*TargetAggregationFunction, error) {
	vals := map[string]TargetAggregationFunction{
		"max":  TargetAggregationFunctionMax,
		"mean": TargetAggregationFunctionMean,
		"min":  TargetAggregationFunctionMin,
		"none": TargetAggregationFunctionNone,
		"sum":  TargetAggregationFunctionSum,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := TargetAggregationFunction(input)
	return &out, nil
}

type TargetLagsMode string

const (
	TargetLagsModeAuto   TargetLagsMode = "Auto"
	TargetLagsModeCustom TargetLagsMode = "Custom"
)

func PossibleValuesForTargetLagsMode() []string {
	return []string{
		string(TargetLagsModeAuto),
		string(TargetLagsModeCustom),
	}
}

func (s *TargetLagsMode) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseTargetLagsMode(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseTargetLagsMode(input string) (*TargetLagsMode, error) {
	vals := map[string]TargetLagsMode{
		"auto":   TargetLagsModeAuto,
		"custom": TargetLagsModeCustom,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := TargetLagsMode(input)
	return &out, nil
}

type TargetRollingWindowSizeMode string

const (
	TargetRollingWindowSizeModeAuto   TargetRollingWindowSizeMode = "Auto"
	TargetRollingWindowSizeModeCustom TargetRollingWindowSizeMode = "Custom"
)

func PossibleValuesForTargetRollingWindowSizeMode() []string {
	return []string{
		string(TargetRollingWindowSizeModeAuto),
		string(TargetRollingWindowSizeModeCustom),
	}
}

func (s *TargetRollingWindowSizeMode) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseTargetRollingWindowSizeMode(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseTargetRollingWindowSizeMode(input string) (*TargetRollingWindowSizeMode, error) {
	vals := map[string]TargetRollingWindowSizeMode{
		"auto":   TargetRollingWindowSizeModeAuto,
		"custom": TargetRollingWindowSizeModeCustom,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := TargetRollingWindowSizeMode(input)
	return &out, nil
}

type TaskType string

const (
	TaskTypeClassification                TaskType = "Classification"
	TaskTypeForecasting                   TaskType = "Forecasting"
	TaskTypeImageClassification           TaskType = "ImageClassification"
	TaskTypeImageClassificationMultilabel TaskType = "ImageClassificationMultilabel"
	TaskTypeImageInstanceSegmentation     TaskType = "ImageInstanceSegmentation"
	TaskTypeImageObjectDetection          TaskType = "ImageObjectDetection"
	TaskTypeRegression                    TaskType = "Regression"
	TaskTypeTextClassification            TaskType = "TextClassification"
	TaskTypeTextClassificationMultilabel  TaskType = "TextClassificationMultilabel"
	TaskTypeTextNER                       TaskType = "TextNER"
)

func PossibleValuesForTaskType() []string {
	return []string{
		string(TaskTypeClassification),
		string(TaskTypeForecasting),
		string(TaskTypeImageClassification),
		string(TaskTypeImageClassificationMultilabel),
		string(TaskTypeImageInstanceSegmentation),
		string(TaskTypeImageObjectDetection),
		string(TaskTypeRegression),
		string(TaskTypeTextClassification),
		string(TaskTypeTextClassificationMultilabel),
		string(TaskTypeTextNER),
	}
}

func (s *TaskType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseTaskType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseTaskType(input string) (*TaskType, error) {
	vals := map[string]TaskType{
		"classification":                TaskTypeClassification,
		"forecasting":                   TaskTypeForecasting,
		"imageclassification":           TaskTypeImageClassification,
		"imageclassificationmultilabel": TaskTypeImageClassificationMultilabel,
		"imageinstancesegmentation":     TaskTypeImageInstanceSegmentation,
		"imageobjectdetection":          TaskTypeImageObjectDetection,
		"regression":                    TaskTypeRegression,
		"textclassification":            TaskTypeTextClassification,
		"textclassificationmultilabel":  TaskTypeTextClassificationMultilabel,
		"textner":                       TaskTypeTextNER,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := TaskType(input)
	return &out, nil
}

type UseStl string

const (
	UseStlNone        UseStl = "None"
	UseStlSeason      UseStl = "Season"
	UseStlSeasonTrend UseStl = "SeasonTrend"
)

func PossibleValuesForUseStl() []string {
	return []string{
		string(UseStlNone),
		string(UseStlSeason),
		string(UseStlSeasonTrend),
	}
}

func (s *UseStl) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseUseStl(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseUseStl(input string) (*UseStl, error) {
	vals := map[string]UseStl{
		"none":        UseStlNone,
		"season":      UseStlSeason,
		"seasontrend": UseStlSeasonTrend,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := UseStl(input)
	return &out, nil
}

type ValidationMetricType string

const (
	ValidationMetricTypeCoco    ValidationMetricType = "Coco"
	ValidationMetricTypeCocoVoc ValidationMetricType = "CocoVoc"
	ValidationMetricTypeNone    ValidationMetricType = "None"
	ValidationMetricTypeVoc     ValidationMetricType = "Voc"
)

func PossibleValuesForValidationMetricType() []string {
	return []string{
		string(ValidationMetricTypeCoco),
		string(ValidationMetricTypeCocoVoc),
		string(ValidationMetricTypeNone),
		string(ValidationMetricTypeVoc),
	}
}

func (s *ValidationMetricType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseValidationMetricType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseValidationMetricType(input string) (*ValidationMetricType, error) {
	vals := map[string]ValidationMetricType{
		"coco":    ValidationMetricTypeCoco,
		"cocovoc": ValidationMetricTypeCocoVoc,
		"none":    ValidationMetricTypeNone,
		"voc":     ValidationMetricTypeVoc,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ValidationMetricType(input)
	return &out, nil
}
