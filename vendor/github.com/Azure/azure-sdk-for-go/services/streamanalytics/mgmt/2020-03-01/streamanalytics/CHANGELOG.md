# Change History

## Breaking Changes

### Removed Funcs

1. *ScalarFunctionConfiguration.UnmarshalJSON([]byte) error

### Struct Changes

#### Removed Structs

1. ScalarFunctionConfiguration

#### Removed Struct Fields

1. ScalarFunctionProperties.*ScalarFunctionConfiguration

## Additive Changes

### New Constants

1. TypeBasicFunctionProperties.TypeBasicFunctionPropertiesTypeAggregate

### New Funcs

1. *AggregateFunctionProperties.UnmarshalJSON([]byte) error
1. *FunctionConfiguration.UnmarshalJSON([]byte) error
1. *FunctionProperties.UnmarshalJSON([]byte) error
1. AggregateFunctionProperties.AsAggregateFunctionProperties() (*AggregateFunctionProperties, bool)
1. AggregateFunctionProperties.AsBasicFunctionProperties() (BasicFunctionProperties, bool)
1. AggregateFunctionProperties.AsFunctionProperties() (*FunctionProperties, bool)
1. AggregateFunctionProperties.AsScalarFunctionProperties() (*ScalarFunctionProperties, bool)
1. AggregateFunctionProperties.MarshalJSON() ([]byte, error)
1. FunctionProperties.AsAggregateFunctionProperties() (*AggregateFunctionProperties, bool)
1. ScalarFunctionProperties.AsAggregateFunctionProperties() (*AggregateFunctionProperties, bool)

### Struct Changes

#### New Structs

1. AggregateFunctionProperties
1. FunctionConfiguration

#### New Struct Fields

1. FunctionProperties.*FunctionConfiguration
1. ScalarFunctionProperties.*FunctionConfiguration
