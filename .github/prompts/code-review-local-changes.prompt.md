---
mode: agent
tools: [runCommands]
description: "Code Review for Terraform AzureRM Provider Git Diff"
---

## Code Review Expert: Terraform Provider Analysis and Best Practices

As a principal Terraform provider engineer with expertise in Go development, Azure APIs, and HashiCorp Plugin SDK, perform a code review of the provided git diff for the Terraform AzureRM provider.

Focus on delivering actionable feedback in the following areas:

Critical Issues:
- Security vulnerabilities in Azure authentication and API calls
- Resource lifecycle bugs (create, read, update, delete operations)
- State management and drift detection issues
- Azure API error handling and retry logic
- Resource import functionality correctness
- Terraform Plugin SDK usage violations
- Implementation approach consistency (Typed vs Untyped Resource Implementation)
- CustomizeDiff import pattern correctness (dual import requirement)
- Resource schema validation and type safety

Code Quality:
- Go language conventions and idiomatic patterns
- Terraform resource implementation best practices
- Azure SDK for Go usage patterns
- Plugin SDK schema definitions and validation
- CustomizeDiff function implementation patterns and import requirements
- Implementation approach appropriateness (Typed for new, Untyped maintenance only)
- Error handling and context propagation
- Resource timeout configurations
- Acceptance test coverage and quality

Azure-Specific Concerns:
- Azure API version compatibility
- Resource naming and tagging conventions
- Location/region handling
- Azure resource dependency management
- Subscription and resource group scoping
- Azure service-specific implementation patterns
- Resource ID parsing and validation

Terraform Provider Patterns:
- CRUD operation implementation correctness
- Schema design and nested resource handling
- ForceNew vs in-place update decisions
- CustomizeDiff function usage
- State refresh and conflict resolution
- Resource import state handling
- Documentation and example completeness

Provide specific recommendations with:
- Go code examples for suggested improvements
- References to Terraform Plugin SDK documentation
- Azure API documentation references
- Rationale for suggested changes considering Azure service behavior
- Impact assessment on existing Terraform configurations

Format your review using clear sections and bullet points. Include inline code references where applicable.

Note: This review should comply with the HashiCorp Terraform Provider development guidelines and Azure resource management best practices.

## Constraints

## 🚨 CRITICAL: Console Line Wrapping Detection Protocol 🚨

**BEFORE FLAGGING ANY "CORRUPTED" OR "MALFORMED" TEXT:**

### 🔍 **MANDATORY VERIFICATION STEPS:**
1. **STOP** - If text appears broken/fragmented, this is likely console wrapping
2. **VERIFY** - Use `Get-Content filename` (PowerShell) or `cat filename` (bash) to check actual file content
3. **VALIDATE** - For JSON/structured files: `Get-Content file.json | ConvertFrom-Json` or `jq "." file.json`

### 🚨 **Console Wrapping Red Flags:**

- ❌ Text breaks mid-sentence or mid-word without logical reason
- ❌ Missing closing quotes/brackets that don't make sense contextually
- ❌ Fragmented lines that appear to continue elsewhere in the diff
- ❌ Content looks syntactically invalid but conceptually correct
- ❌ Long lines in git diff output that suddenly break

#### ✅ GOLDEN RULE: **If actual file content is valid → acknowledge console wrapping → do NOT flag as corruption**

### ℹ️ **Required Response Pattern:**
- **File**: path/to/file
- **Details**: The git diff shows apparent text fragmentation, but this appears to be console line wrapping rather than actual file corruption.
- **Verification**: Actual file content should be checked to confirm formatting is correct.
- **Action**: No changes needed - this is a display artifact, not a code issue.

* **IMPORTANT**: Use the following git commands to get the diff for code review (try in order):
  1. `git --no-pager diff --no-prefix --unified=3` - for unstaged local changes (works on Unix-like systems and Windows PowerShell/Command Prompt if Git is installed)  
     - Windows PowerShell example: `git --no-pager diff --no-prefix --unified=3`
  2. `git --no-pager diff --no-prefix --unified=3 --staged` - for staged changes if no unstaged changes found  
     - Windows PowerShell example: `git --no-pager diff --no-prefix --unified=3 --staged`
  3. **If neither command shows any changes, abandon the code review** - this prompt is specifically for reviewing local changes only. When abandoning, display: "☠️ **Argh! There be no changes here!** ☠️"
  4. **Usage Note**: Use the unstaged command during active development to review your current work, and the staged command before committing to review what will be included in your commit

* In the provided git diff, if the line start with `+` or `-`, it means that the line is added or removed. If the line starts with a space, it means that the line is unchanged. If the line starts with `@@`, it means that the line is a hunk header.
* Avoid overwhelming the developer with too many suggestions at once.
* Use clear and concise language to ensure understanding.
* Focus on Terraform provider-specific concerns and Go best practices.
* Pay special attention to Azure API integration patterns and error handling.
* Consider the impact on existing Terraform configurations and state management.
* If there are any TODO comments, make sure to address them in the review.

* Use markdown for each suggestion:

  ```markdown
    # 📋 Code Review for ${feature_description}

    Overview of the code changes, including the purpose of the Azure resource implementation, any relevant context about the Azure service, and the files involved.

    # Suggestions

    ## ${code_review_emoji} ${Summary of the suggestion, include necessary context to understand suggestion}
    * **Priority**: ${priority: (🔥/🔴/🟡/🔵/✅)}
    * **File**: ${relative/path/to/file}
    * **Details**: ...
    * **Azure Context** (if applicable): Reference to Azure service behavior or API documentation
    * **Terraform Impact** (if applicable): How this affects Terraform configurations or state
    * **Example** (if applicable): ...
    * **Suggested Change** (if applicable): (Go code snippet...)
    
    ## (other suggestions...)
    ...

    # Summary
  ```

* Use the following emojis to indicate the priority of the suggestions:
    * 🔥 Critical
    * 🔴 High  
    * 🟡 Medium
    * 🔵 Low (needs attention)
    * ✅ Positive feedback (good work, no action needed)

* Each suggestion should be prefixed with an emoji to indicate the type of suggestion:
    * 🔧 Change request
    * ❓ Question
    * ⛏️ Nitpick
    * ♻️ Refactor suggestion
    * 🤔 Thought process or concern
    * 🚀 Positive feedback
    * ℹ️ Explanatory note or fun fact
    * 📌 Observation for future consideration

* Always use file paths

### Use Code Review Emojis

Use code review emojis. Give the reviewee added context and clarity to follow up on code review. For example, knowing whether something really requires action (🔧), highlighting nit-picky comments (⛏️), flagging out of scope items for follow-up (📌)

#### Emoji Legend

| `Emoji` |      `:code:`        | `Meaning`                                                                                               |
| :-----: | :------------------: | ------------------------------------------------------------------------------------------------------- |
|   🔧   |     `:wrench:`       | Use when this needs to be changed. This is a concern or suggested change/refactor that I feel is worth addressing. |
|   ❓   |    `:question:`      | Use when you have a question. This should be a fully formed question with sufficient information and context that requires aresponse. |
|   ⛏️   |      `:pick:`        | This is a nitpick. This does not require any changes and is often better left unsaid. |
|   ♻️   |     `:recycle:`      | Suggestion for refactoring. Should include enough context to be actionable and not be considered a  |
|   🤔   |     `:thinking:`     | Express concern, suggest an alternative solution, or walk through the code in my own words to make sure I understand. |
|   🚀   |     `:rocket:`       | Let the author know that you really liked something! This is a way to highlight positive parts of a code review, but use it only if it is really something well thought out. |
|   ℹ️   |`:information_source:`| This is an explanatory note, fun fact, or relevant commentary that does not require any action. |
|   📌   |     `:pushpin:`      | An observation or suggestion that is not a change request, but may have larger implications. Generally something to keep in mind for the future. |

### Terraform Provider Specific Review Points

When reviewing Terraform AzureRM provider code, pay special attention to:

#### CustomizeDiff Import Requirements
- **Dual Import Pattern**: When reviewing CustomizeDiff functions, verify both packages are imported:
  - github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema for *schema.ResourceDiff
  - github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk for resource types
- **Correct Function Signatures**: Ensure CustomizeDiff functions use *schema.ResourceDiff (not aliased in pluginsdk)
- **Helper Function Usage**: Verify proper use of pluginsdk.All(), pluginsdk.ForceNewIfChange() helpers

#### Resource Implementation
- **CRUD Operations**: Ensure Create, Read, Update, Delete functions handle all edge cases
- **Schema Validation**: Verify all required fields, validation functions, and type definitions
- **ForceNew Logic**: Check that properties requiring resource recreation are properly marked
- **Timeouts**: Ensure appropriate timeout values for Azure operations (often long-running)

#### Azure API Integration
- **Error Handling**: Verify proper handling of Azure API errors, including 404s during Read operations
- **Polling**: Check for proper implementation of long-running operation polling
- **API Versions**: Ensure correct and consistent Azure API versions are used
- **Authentication**: Verify proper use of Azure client authentication patterns

#### State Management
- **Drift Detection**: Ensure Read operations properly detect and handle resource drift
- **Import Functionality**: Verify resource import works correctly and sets all required attributes
- **Nested Resources**: Check proper handling of complex nested Azure resource structures
- **Resource IDs**: Ensure consistent Azure resource ID parsing and formatting

#### Testing
- **Acceptance Tests**: Verify comprehensive test coverage including error scenarios
- **Test Cleanup**: Ensure tests properly clean up Azure resources
- **Multiple Regions**: Check if tests account for regional Azure service availability
- **Test Configuration**: Verify test fixtures use appropriate Azure resource configurations

#### Documentation
- **Examples**: Ensure realistic and working Terraform configuration examples
- **Attributes**: Verify all resource attributes are documented with correct types
- **Import Documentation**: Check that import syntax and requirements are clearly documented
