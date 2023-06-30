## Convert Untyped Resource To Typed Resource

This application generates [typed SDK](../../sdk) resource code from untyped resource. It will generate a golang source code file named as `xxxx_typed.go` in the same folder of the given untyped SDK resource.

The generated source code may contain syntax error, please fix them manually. You need to check over the generated code and make sure it is correct even there is no compiled error.

You have to remove the old resource registration and add the new typed resource registration manually when making sure the code is correct.

## Example Usage

```
$ go run main.go <resource_type1> <resource_type2> <...>
```

E.g.

```
$ go run main.go azurerm_automation_account
```

## Arguments

* `resource_type`: The resource type to generate the schema. 

## What It Does
- [x] generate Model struct from schema and implement `meta.Resource` interface.
- [x] split schema definition into `Arguments()` and `Attriutes()` method.
- [x] implement `ModelObject()` and `ResourceType()`.
- [x] move old `create/update/update/delete` function body into model `Create/Update/Read/Delete` method.
- [x] `Timeout` value is moved to corresponding method definition.
- [x] `d.Get("key").(string)` => `val = model.Key`
- [x] `d.GetOk("key")` => `meta.ResourceData.GetOk("key")`
- [x] `d.IsNewResource()` => `meta.ResourceData.IsNewResource()`
- [x] `meta.(*clients.Client)` => `meta.Client`
- [x] `d.SetId("")` => `metadata.MarkAsGone`
- [x] `d.Set("key", value)` => `model.Key = value`  remove `if err := d.Set("key", value); err != nil {xxx}`
- [x] `tf.ImportAsExistsError` => `meta.ResourceRequiresImport(m.ResourceType(), id)`
- [x] Local funtion call => function call with receiver. `flattenXXX` => `m.flattenXXX`
- [x] Local function auto refactor: expand args type change, and `input['key']` => `input.Field`
- [x] Local function auto refactor: flatten return type change : `[]interface{}` => `[]XXXModel`
- [x] Auto build model composite from map composite
- [ ] flatten/expand function auto rewrite (partially done)

## What It Can't do now
1. Schema with custom logic or composite. like cdn_endpoint_custom_domain_resource.go
2. Resource source code split into multiple files.