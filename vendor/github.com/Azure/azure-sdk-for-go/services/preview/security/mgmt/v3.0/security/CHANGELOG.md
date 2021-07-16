# Change History

## Breaking Changes

### Signature Changes

#### Const Types

1. KindAAD changed type from KindEnum1 to KindEnum
1. KindATA changed type from KindEnum1 to KindEnum
1. KindAlertSimulatorRequestProperties changed type from KindEnum2 to KindEnum1
1. KindBundles changed type from KindEnum2 to KindEnum1
1. KindCEF changed type from KindEnum1 to KindEnum
1. KindDataExportSettings changed type from KindEnum to KindEnum2
1. KindExternalSecuritySolution changed type from KindEnum1 to KindEnum
1. KindSetting changed type from KindEnum to KindEnum2

#### Struct Fields

1. AadExternalSecuritySolution.Kind changed type from KindEnum1 to KindEnum
1. AlertSimulatorBundlesRequestProperties.Kind changed type from KindEnum2 to KindEnum1
1. AlertSimulatorRequestProperties.Kind changed type from KindEnum2 to KindEnum1
1. AtaExternalSecuritySolution.Kind changed type from KindEnum1 to KindEnum
1. CefExternalSecuritySolution.Kind changed type from KindEnum1 to KindEnum
1. DataExportSettings.Kind changed type from KindEnum to KindEnum2
1. ExternalSecuritySolution.Kind changed type from KindEnum1 to KindEnum
1. Setting.Kind changed type from KindEnum to KindEnum2

## Additive Changes

### New Constants

1. KindEnum2.KindAlertSyncSettings
1. OnboardingKind.Evaluation
1. OnboardingKind.Purchased

### New Funcs

1. *AlertSyncSettings.UnmarshalJSON([]byte) error
1. AlertSyncSettings.AsAlertSyncSettings() (*AlertSyncSettings, bool)
1. AlertSyncSettings.AsBasicSetting() (BasicSetting, bool)
1. AlertSyncSettings.AsDataExportSettings() (*DataExportSettings, bool)
1. AlertSyncSettings.AsSetting() (*Setting, bool)
1. AlertSyncSettings.MarshalJSON() ([]byte, error)
1. DataExportSettings.AsAlertSyncSettings() (*AlertSyncSettings, bool)
1. IotDefenderSettingsProperties.MarshalJSON() ([]byte, error)
1. Setting.AsAlertSyncSettings() (*AlertSyncSettings, bool)

### Struct Changes

#### New Structs

1. AlertSyncSettingProperties
1. AlertSyncSettings

#### New Struct Fields

1. IotDefenderSettingsProperties.EvaluationEndTime
