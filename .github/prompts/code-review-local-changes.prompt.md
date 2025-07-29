---
mode: agent
description: "Code Review for Terraform AzureRM Provider Local Changes"
---

## Code Review Expert: Terraform Provider Analysis and Best Practices

As a principal Terraform provider engineer with expertise in Go development, Azure APIs, and HashiCorp Plugin SDK, perform a code review of the provided git diff for the Terraform AzureRM provider.

Focus on delivering actionable feedback in the following areas:

**Critical Issues**:
- Security vulnerabilities in Azure authentication and API calls
- Resource lifecycle bugs (create, read, update, delete operations)
- State management and drift detection issues
- Azure API error handling and retry logic
- Resource import functionality correctness
- Terraform Plugin SDK usage violations
- Implementation approach consistency (Typed vs Untyped Resource Implementation)
- CustomizeDiff import pattern correctness (dual import requirement)
- Resource schema validation and type safety

**Code Quality**:
- Go language conventions and idiomatic patterns
- Terraform resource implementation best practices
- Azure SDK for Go usage patterns
- Plugin SDK schema definitions and validation
- CustomizeDiff function implementation patterns and import requirements
- Implementation approach appropriateness (Typed for new, Untyped maintenance only)
- Error handling and context propagation
- Resource timeout configurations
- Acceptance test coverage and quality
- **Tests use ONLY ExistsInAzure() check with ImportStep() - NO redundant field validation**

**Azure-Specific Concerns**:
- Azure API version compatibility
- Resource naming and tagging conventions
- Location/region handling
- Azure resource dependency management
- Subscription and resource group scoping
- Azure service-specific implementation patterns
- Resource ID parsing and validation

**Terraform Provider Patterns**:
- CRUD operation implementation correctness
- Schema design and nested resource handling
- ForceNew vs in-place update decisions
- CustomizeDiff function usage
- State refresh and conflict resolution
- Resource import state handling
- Documentation and example completeness

**Provide specific recommendations with**:
- Go code examples for suggested improvements
- References to Terraform Plugin SDK documentation
- Azure API documentation references
- Rationale for suggested changes considering Azure service behavior
- Impact assessment on existing Terraform configurations

Format your review using clear sections and bullet points. Include inline code references where applicable.

**Note**: This review should comply with the HashiCorp Terraform Provider development guidelines and Azure resource management best practices.

## Constraints
* Before you start the code review, please explicitly check off each item in the MANDATORY PRE-REVIEW CHECKLIST and show me your verification.
* Only flag corruption issues IF `read_file` shows the same problems as the git diff. If `read_file` shows clean content, acknowledge console wrapping.

For any suspected issues, you **MUST** use this exact format:
- **Suspected Issue**: [describe]
- **Verification Command**: read_file
- **Actual File Content**: [paste results]
- **Assessment**: [`console wrapping` **OR** `actual issue`]
- **Action**: [required]"

This prompt file contains its own **verification protocols**. You **MUST** follow those protocols when reviewing this very file. Do not create exceptions for reviewing prompt files themselves.
If you flag **false positives** without proper verification, **STOP** and **RESTART** following the checklist correctly.

**Priority order for file verification:**
1. read_file (most reliable)
2. Direct file access tools
3. Terminal commands (least reliable for content verification)

Follow the `code-review-local-changes.prompt.md` instructions. Before flagging **ANY** issues:

1. Check off each mandatory checklist item
2. Use `read_file` **FIRST** for any suspected corruption
3. Use the mandatory verification template format
4. Only flag issues that exist in `read_file` output
5. If `read_file` shows clean content but terminal shows issues, acknowledge console wrapping

**Show me your checklist verification before proceeding with the review.**

## 🔍 **MANDATORY PRE-REVIEW CHECKLIST**

**BEFORE FLAGGING ISSUES:**
```markdown
- [ ] I will verify actual file content first with `cat` or `Get-Content`
- [ ] I understand: Git diff wrapping ≠ File corruption
- [ ] I will NOT assume formatting in diff = actual problems

RULE: Always verify file content before flagging corruption
```
**ONLY PROCEED AFTER CHECKING ALL BOXES ABOVE**

### **AUTOMATIC VERIFICATION TRIGGERS**

**IF YOU SEE ANY OF THESE IN GIT DIFF, IMMEDIATELY RUN FILE VERIFICATION:**

- ❌ `Git` diff formatting issues
- ❌ `emoji` display as `??`
- ❌ Line breaks in `diff`
- ❌ Fragmented text in `diff`

**FILE VERIFICATION COMMANDS:**
* **Unix/Linux/macOS**: `sed -n "[line-5],[line+5]p" filename`
* **Windows PowerShell**: `Get-Content "filename" | Select-Object -Skip [line-5] -First 10`
* **Windows Command Prompt**: `more +[line-5] filename | head -10` (if available)

### 📋 **MANDATORY VERIFICATION TEMPLATE**

When suspicious content is found, use this template:

- **Suspected Issue**: [describe what looks wrong in git diff]
- **Verification Command**: `cat 'filename'`
  - Windows PowerShell: `Get-Content 'filename'`
- **Actual File Content**: [paste verification results]
- **Assessment**: [console wrapping **OR** actual issue]
- **Action**: [no action needed **OR** specific fix required]

### 🔍 **VERIFICATION SCENARIO EXAMPLES**

**Scenario 1: Emoji Display Issues**
```text
Git diff shows: 🔍 COMMITTED CHANGES CODE REVIEW 🔍
Reality: Emojis display as ?? in some terminals but are actually valid
Action: Use read_file to verify actual content exists
Result: Console display issue, not file corruption
```

**Scenario 2: JSON/YAML Line Breaking**
```json
Git diff shows:
{
    "github.copilot.chat.commitMessage.generation.instructions": "Follow terraform coding standards and best
practices when reviewing code changes"
}

Reality: JSON appears broken due to terminal line wrapping
Action: Use read_file to verify JSON is valid on single line
Result: Console wrapping, not malformed JSON
```

**Scenario 3: Text Fragmentation**
```text
Git diff shows:
make valid syntax appear broken.*: When reviewing git diff output in terminal
console, be aware that long lines may wrap

Reality: Text appears fragmented mid-sentence
Action: Use read_file to see the actual continuous text without wrapping
Result: Console display artifact, not broken text
```

**Scenario 4: Missing Quotes/Brackets**
```yaml
Git diff shows:
"description": "Code Review for Terraform AzureRM Provider Git Diff
---

Reality: Closing quote appears missing due to line wrap
Action: Use read_file to verify the closing quote is actually present
Result: Console wrapping, not syntax error
```

**Scenario 5: Code Block Fragmentation**
```markdown
Git diff shows:
# Header
some text that breaks
off mid-sentence and continues

Reality: Markdown appears malformed due to terminal wrapping
Action: Use read_file to verify markdown structure is correct
Result: Console display issue, not markdown corruption
```

## Console Output Interpretation

**🚨 CRITICAL: CONSOLE LINE WRAPPING DETECTION PROTOCOL 🚨**

**CONSOLE LINE WRAPPING WARNING**: When reviewing `git` diff output in terminal/console, be aware that long lines may wrap and appear malformed. Always verify actual file content for syntax validation, especially for `JSON`, `YAML`, or structured data files. Console wrapping can make valid syntax appear broken.

**VERIFICATION PROTOCOL FOR SUSPECTED ISSUES**:
### 🔍 **MANDATORY VERIFICATION STEPS:**
1. **STOP** - If text appears broken/fragmented, this is likely console wrapping
2. **VERIFY** - Use `cat filename` to check actual file content
  - Windows PowerShell: `Get-Content 'filename'`
3. **VALIDATE** - For `JSON`/structured files: `jq "." file.json`
  - Windows PowerShell: `Get-Content file.json | ConvertFrom-Json`

### 🚨 **CONSOLE WRAPPING RED FLAGS:** 🚨
- ❌ Text breaks mid-sentence or mid-word without logical reason
- ❌ Missing closing quotes/brackets that don't make sense contextually
- ❌ Fragmented lines that appear to continue elsewhere in the diff
- ❌ Content looks syntactically invalid but conceptually correct
- ❌ Long lines in git diff output that suddenly break

#### ✅ **GOLDEN RULE**: If actual file content is valid → acknowledge console wrapping → DO NOT FLAG as corruption

**Verification Rule**: If actual file content is valid, acknowledge console wrapping and do not flag as an issue

**Git Command Requirements:**
* `Git` must be installed and available in `PATH`
* Windows: `Git for Windows` or `Git` integrated with `PowerShell`
* Verify `git` availability: `git --version`

* **IMPORTANT**: Use the following git commands to get the diff for code review (try in order):
  1. `git --no-pager diff --no-prefix --unified=3` - for unstaged local changes
  2. `git --no-pager diff --no-prefix --unified=3 --staged` - for staged changes if no unstaged changes found
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

    # 📋 Code Review for ${change_description}

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

* Use the following emojis to indicate the priority of the suggestions:
    * 🔥 Critical
    * 🔴 High
    * 🟡 Medium
    * 🔵 Low (needs attention)
    * ✅ good work, no action needed

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
- **Conditional Import Pattern**: Import requirements depend on the CustomizeDiff implementation:
  - **Dual Imports Required**: When using *schema.ResourceDiff directly in CustomizeDiff functions:
    - github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema (for *schema.ResourceDiff)
    - github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk (for helpers)
  - **Single Import Sufficient**: When using *pluginsdk.ResourceDiff (legacy resources):
    - github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema (only this import needed)
- **Function Signature Analysis**: Check the function signature to determine import requirements:
  - *schema.ResourceDiff → Usually typed resources, may need dual imports
  - *pluginsdk.ResourceDiff → Usually legacy resources, single import sufficient
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
