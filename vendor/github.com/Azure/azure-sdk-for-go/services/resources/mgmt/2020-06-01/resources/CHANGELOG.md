# Change History

## Breaking Changes

### Removed Constants

1. AliasPathAttributes.Modifiable
1. AliasPathAttributes.None
1. AliasPathTokenType.Any
1. AliasPathTokenType.Array
1. AliasPathTokenType.Boolean
1. AliasPathTokenType.Integer
1. AliasPathTokenType.NotSpecified
1. AliasPathTokenType.Number
1. AliasPathTokenType.Object
1. AliasPathTokenType.String
1. ChangeType.Create
1. ChangeType.Delete
1. ChangeType.Deploy
1. ChangeType.Ignore
1. ChangeType.Modify
1. ChangeType.NoChange
1. DeploymentMode.Complete
1. DeploymentMode.Incremental
1. OnErrorDeploymentType.LastSuccessful
1. OnErrorDeploymentType.SpecificDeployment
1. WhatIfResultFormat.FullResourcePayloads
1. WhatIfResultFormat.ResourceIDOnly

### Signature Changes

#### Funcs

1. GroupsClient.Delete
	- Params
		- From: context.Context, string
		- To: context.Context, string, string
1. GroupsClient.DeletePreparer
	- Params
		- From: context.Context, string
		- To: context.Context, string, string

## Additive Changes

### New Constants

1. AliasPathAttributes.AliasPathAttributesModifiable
1. AliasPathAttributes.AliasPathAttributesNone
1. AliasPathTokenType.AliasPathTokenTypeAny
1. AliasPathTokenType.AliasPathTokenTypeArray
1. AliasPathTokenType.AliasPathTokenTypeBoolean
1. AliasPathTokenType.AliasPathTokenTypeInteger
1. AliasPathTokenType.AliasPathTokenTypeNotSpecified
1. AliasPathTokenType.AliasPathTokenTypeNumber
1. AliasPathTokenType.AliasPathTokenTypeObject
1. AliasPathTokenType.AliasPathTokenTypeString
1. ChangeType.ChangeTypeCreate
1. ChangeType.ChangeTypeDelete
1. ChangeType.ChangeTypeDeploy
1. ChangeType.ChangeTypeIgnore
1. ChangeType.ChangeTypeModify
1. ChangeType.ChangeTypeNoChange
1. DeploymentMode.DeploymentModeComplete
1. DeploymentMode.DeploymentModeIncremental
1. OnErrorDeploymentType.OnErrorDeploymentTypeLastSuccessful
1. OnErrorDeploymentType.OnErrorDeploymentTypeSpecificDeployment
1. WhatIfResultFormat.WhatIfResultFormatFullResourcePayloads
1. WhatIfResultFormat.WhatIfResultFormatResourceIDOnly

### Struct Changes

#### New Structs

1. ZoneMapping

#### New Struct Fields

1. ProviderResourceType.ZoneMappings
