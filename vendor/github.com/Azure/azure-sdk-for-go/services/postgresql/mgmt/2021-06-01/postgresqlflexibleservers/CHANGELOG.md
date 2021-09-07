# Change History

## Breaking Changes

### Signature Changes

#### Struct Fields

1. RestartParameter.FailoverMode changed type from *string to FailoverMode

## Additive Changes

### New Constants

1. FailoverMode.FailoverModeForcedFailover
1. FailoverMode.FailoverModeForcedSwitchover
1. FailoverMode.FailoverModePlannedFailover
1. FailoverMode.FailoverModePlannedSwitchover

### New Funcs

1. PossibleFailoverModeValues() []FailoverMode
