# Provider Documentation Standards

In an effort to keep the [provider documentation](https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs) consistent, this page documents some standards that have been agreed on. 

This page will grow over time, and suggestions are welcome!

## Notes

Note blocks are used to provide additional information to users beyond the basic description of a resource, argument or attribute.

In the past, there have been different approaches to how notes were formatted, some examples are:

- Different words to indicate level of importance, e.g. `Info`, `Important`, `Caution`, and `Be Aware`.
- Capitalisation differences, e.g. `Note:` vs `NOTE:`.
- Whether or not a colon is included, e.g. `Note:` vs `Note`.

Going forward, all notes should follow the exact same format (`(->|~>|!>) **Note:**`) where level of importance is indicated through the different types of notes as documented below.

Breaking changes have previously been added as notes to the resource documentation.
These should no longer be included, instead follow these guidelines:

- Breaking changes in a minor version should be added to the top of the changelog.
- Breaking changes in a major version should be added to the upgrade guide.

> We may revisit the guidelines above and/or add a specific place in the documentation for all breaking changes in minor versions.

<!-- 
    - TODO: Considerations for when to add notes? We probably don't want to overdo it (More relevant to informational notes)
-->

### Informational Note

Informational note blocks should generally be used when a note provides additional useful information, recommendations and/or tips to the user.

To add an informational note, use `-> **Note:**`, within the Terraform registry documentation this will template as a block with an info icon.

For example, extra information on the supported values for an argument, possibly linking to external documentation for the resource/service:

```markdown
* `type` - (Required) The type. Possible values include `This`, `That`, and `Other`.

-> **Note:** More information on each of the supported types can be found in [type documentation](link-to-additional-info)
```

### Warning Note

Warning note blocks should generally be used when a note provides information that the user will need to avoid certain errors, however if these errors are encountered they should not break anything or cause irreversible changes.

To add a warning note, use `~> **Note:**`, within the Terraform registry documentation this will template as a block with a warning icon.

For example, an argument that is optional but required when another argument is set to `true`:

```markdown
* `optional_argument_enabled` - (Optional) Is the optional argument enabled? Defaults to `false`.

* `optional_argument` - (Optional) An optional argument.

~> **Note:** The argument `optional_argument` is required when `optional_argument_enabled` is set to `true`.
```

### Caution Note

Caution note blocks should generally be used when a note provides critical information on potential irreversible changes, data loss or other things that can negatively affect a user's environment.

To add a caution note, use `!> **Note:**`, within the Terraform registry documentation this will template as a block with a caution icon.

For example, an argument that when set to `true` cannot be reversed without recreating the resource:

```markdown
* `irreversible_argument_enabled` - (Optional) Is irreversible argument enabled? Defaults to `false`.

!> **Note:** The argument `irreversible_argument_enabled` cannot be disabled after being enabled.
```