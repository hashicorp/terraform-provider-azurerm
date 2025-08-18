---
mode: agent
description: "Code Review Prompt for Terraform AzureRM Provider Local Changes"
---

# Code Review Expert: Terraform Provider Analysis and Best Practices

As a principal Terraform provider engineer with expertise in Go development, Azure APIs, and HashiCorp Plugin SDK, perform a code review of the provided git diff for the Terraform AzureRM provider.

**STRUCTURED EXECUTION FRAMEWORK**:

- Follow this structured approach to ensure comprehensive and consistent reviews
- Use the provided templates and emoji systems for clarity and prioritization
- Complete verification protocols before flagging issues to prevent false positives
- Adapt the format as needed to provide the most helpful feedback for the specific changes
- Focus on substance over strict formatting - the goal is actionable, valuable insights
- When deviating from templates, explain why to maintain transparency

Focus on delivering actionable feedback in the following areas:

## **Critical Issues**:
  - Security vulnerabilities in Azure authentication and API calls
  - Resource lifecycle bugs (create, read, update, delete operations)
  - State management and drift detection issues
  - Azure API error handling and retry logic
  - Resource import functionality correctness
  - Terraform Plugin SDK usage violations
  - Implementation approach consistency (Typed vs Untyped Resource Implementation)
  - CustomizeDiff import pattern correctness (dual import requirement)
  - Resource schema validation and type safety

## **Code Quality**:
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
  - **CRITICAL: Code comments policy enforcement - Comments only for Azure API quirks, complex business logic, or SDK workarounds that cannot be expressed through code structure**
  - **Comment quality requirement - All comments must have clear justification and add genuine value beyond what code structure provides**
  - **Refactoring preference - Before accepting comments, consider if code restructuring could eliminate the need for explanation**
  - **Documentation Quality & Language Standards**:
    - Spelling accuracy in all text content (comments, documentation, README files)
    - Grammar and syntax correctness in documentation
    - Consistent terminology and naming conventions
    - Command examples and usage instructions accuracy
    - Typo detection in visible text content (even if not part of the diff)
    - Professional language standards for user-facing content

## **Azure-Specific Concerns**:
  - Azure API version compatibility
  - Resource naming and tagging conventions
  - Location/region handling
  - Azure resource dependency management
  - Subscription and resource group scoping
  - Azure service-specific implementation patterns
  - Resource ID parsing and validation

## **Terraform Provider Patterns**:
  - CRUD operation implementation correctness
  - Schema design and nested resource handling
  - ForceNew vs in-place update decisions
  - CustomizeDiff function usage
  - State refresh and conflict resolution
  - Resource import state handling
  - Documentation and example completeness

## **Provide specific recommendations with**:
  - Go code examples for suggested improvements
  - References to Terraform Plugin SDK documentation
  - Azure API documentation references
  - Rationale for suggested changes considering Azure service behavior
  - Impact assessment on existing Terraform configurations

Format your review using clear sections and bullet points. Include inline code references where applicable.

**Note**: This review should comply with the HashiCorp Terraform Provider development guidelines and Azure resource management best practices.

## Constraints
  * Before starting the code review, complete the essential verification checklist to ensure accurate issue identification
  * Prioritize file content verification over visual diff formatting to prevent false positives
  * Use structured format guidelines while adapting content to provide maximum value for the specific changes
  * Focus on actionable feedback that improves code quality, security, and maintainability

For any suspected formatting or corruption issues, you **MUST** use this verification approach:
  - **Suspected Issue**: [describe what appears problematic in git diff]
  - **Verification**: Use read_file to check actual content
  - **Assessment**: [console wrapping **OR** genuine issue requiring action]
  - **Response**: [acknowledge display artifact **OR** provide specific remediation]

This prompt file follows its own **verification protocols**. Apply these protocols consistently when reviewing any files, including prompt files themselves.
Focus on preventing false positives through proper verification while maintaining helpful technical insights.

**Verification priority order:**
  1. read_file (most reliable for content validation)
  2. Direct file access tools  
  3. Terminal commands (least reliable for content verification)

Follow the structured review approach below. Before flagging **ANY** formatting or corruption issues:

  1. Apply essential verification checklist
  2. Use `read_file` **FIRST** for any suspected display problems
  3. Use the verification template for transparency
  4. Focus on genuine issues that affect code quality
  5. Acknowledge console wrapping when file content is actually clean

**Complete verification for suspected issues before proceeding with technical review.**

# 🔍 **ESSENTIAL VERIFICATION CHECKLIST**

**BEFORE FLAGGING DISPLAY/FORMATTING ISSUES:**
```markdown
- [ ] I will verify actual file content with read_file for any suspected corruption
- [ ] I understand: Git diff wrapping ≠ File corruption  
- [ ] I will focus on genuine code quality issues over display artifacts
- [ ] I will provide helpful context and actionable suggestions
- [ ] I will review for typos, consistency, and professional standards in visible content

CORE PRINCIPLE: Verify first, then provide valuable technical insights
PRIORITY: Actionable feedback > Format compliance
```

## **COMPREHENSIVE REVIEW SCOPE**

**BEYOND DIFF CHANGES - ALSO CHECK:**
- **Spelling and Grammar**: Review all visible text content for typos and language issues
- **Command Examples**: Verify accuracy of command syntax and examples
- **Consistency**: Check naming conventions and terminology consistency
- **Professional Standards**: Ensure user-facing content meets quality standards
- **Context Quality**: Review surrounding lines shown in diff for overall file quality

## **AUTOMATIC VERIFICATION TRIGGERS**

**IF YOU SEE ANY OF THESE IN GIT DIFF, IMMEDIATELY RUN FILE VERIFICATION:**

- ❌ `Git` diff formatting issues
- ❌ `emoji` display as `??`
- ❌ Line breaks in `diff`
- ❌ Fragmented text in `diff`

**FILE VERIFICATION COMMANDS:**
  - **Unix/Linux/macOS**: `sed -n "[line-5],[line+5]p" filename`
  - **Windows PowerShell**: `Get-Content "filename" | Select-Object -Skip [line-5] -First 10`
  - **Windows Command Prompt**: `more +[line-5] filename | head -10` (if available)

## 📋 **STRUCTURED VERIFICATION TEMPLATE**

When suspicious content is found, use this approach for transparency:

  - **Suspected Issue**: [describe what looks wrong in git diff]
  - **Verification Command**: `cat 'filename'`
  - **Windows PowerShell**: `Get-Content 'filename'`
  - **Actual File Content**: [paste verification results]
  - **Assessment**: [console wrapping **OR** actual issue]
  - **Action**: [no action needed **OR** specific fix required]

## 🔍 **VERIFICATION SCENARIO EXAMPLES**

**Scenario 1: Emoji Display Issues**
```text
Git diff shows: 🔍 COMMITTED CHANGES CODE REVIEW 🔍
Reality: Emojis display as ?? in some terminals but are actually valid
Action: Use read_file to verify actual content exists
Result: Console display issue, not file corruption
```

**Scenario 2: JSON/YAML Line Breaking**
```text
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
```text
Git diff shows:
"description": "Code Review: Terraform AzureRM Provider Git Diff

Reality: Closing quote appears missing due to line wrap
Action: Use read_file to verify the closing quote is actually present
Result: Console wrapping, not syntax error
```

**Scenario 5: Code Block Fragmentation**
```text
Git diff shows:
# Header
some text that breaks
off mid-sentence and continues

Reality: Markdown appears malformed due to terminal wrapping
Action: Use read_file to verify markdown structure is correct
Result: Console display issue, not markdown corruption
```

# Console Output Interpretation

**🚨 CRITICAL: CONSOLE LINE WRAPPING DETECTION POLICY 🚨**

**CONSOLE LINE WRAPPING WARNING**: When reviewing `git` diff output in terminal/console, be aware that long lines may wrap and appear malformed. Always verify actual file content for syntax validation, especially for `JSON`, `YAML`, or structured data files. Console wrapping can make valid syntax appear broken.

**VERIFICATION PROTOCOL FOR SUSPECTED ISSUES**:

## 🔍 **STREAMLINED VERIFICATION STEPS:**
  1. **STOP** - If text appears broken/fragmented, this is likely console wrapping
  2. **VERIFY** - Use `read_file filename` to check actual file content
  3. **NOTE** - Add inline verification result and continue review

## ✅ **VERIFICATION EXAMPLE:**
```markdown
## ℹ️ Console Display Verification
* **Priority**: ✅
* **File**: `filename.go`
* **Details**: Content appeared corrupted in git diff output
* **Action**: Verification completed using read_file
* **Result**: *(Verified: console wrapping - actual content clean)*
* **Assessment**: No issues found - normal console display behavior
```

## ✅ **GOLDEN RULE**: If actual file content is valid → acknowledge console wrapping → DO NOT FLAG as corruption

> **📖 Full Policy Details**: See the complete [Console Line Wrapping Detection Policy](../instructions/error-patterns.instructions.md) for comprehensive guidelines and enforcement procedures.

**Verification Rule**: If actual file content is valid, acknowledge console wrapping and do not flag as an issue

## Git Command Requirements:
  * `Git` must be installed and available in `PATH`
  * Windows: `Git for Windows` or `Git` integrated with `PowerShell`
  * Verify `git` availability: `git --version`

* **IMPORTANT**: Use the following git commands to get the diff for code review (try in order):
  1. `git --no-pager diff --no-prefix --unified=3` - for unstaged local changes
  2. `git --no-pager diff --no-prefix --unified=3 --staged` - for staged changes if no unstaged changes found
  3. **If neither command shows any changes, abandon the code review** - this prompt is specifically for reviewing local changes only. When abandoning, display: "☠️ **Argh! There be no changes here!** ☠️"
  4. **Usage Guidance**: Use the unstaged command during active development to review current work, and the staged command before committing to review what will be included in your commit

  * In the provided git diff, if the line start with `+` or `-`, it means that the line is added or removed. If the line starts with a space, it means that the line is unchanged. If the line starts with `@@`, it means that the line is a hunk header.
  * Avoid overwhelming the developer with too many suggestions at once.
  * Use clear and concise language to ensure understanding.
  * Focus on Terraform provider-specific concerns and Go best practices.
  * Pay special attention to Azure API integration patterns and error handling.
  * Consider the impact on existing Terraform configurations and state management.
  * If there are any TODO comments, make sure to address them in the review.

# Sample Review Output

## 📋 Sample Review Output

  ```markdown
  # 📋 **Code Review**: Azure Front Door Captcha Support

  **Overview**: Adding `captcha_cookie_expiration_in_minutes` property to Front Door Firewall Policy resource with CAPTCHA action type support.

  # Review Comments:

  # 🔧 Missing Error Handling for Captcha Configuration  
  * **Priority**: 🔴
  * **File**: `internal/services/cdn/frontdoor_firewall_policy_resource.go`
  * **Details**: The expand function for captcha settings lacks validation for Azure API constraints
  * **Azure Context**: Front Door captcha requires specific cookie expiration ranges (1-1440 minutes)
  * **Suggested Change**: Add validation in CustomizeDiff or schema ValidateFunc

  # 🚀 Proper Schema Design
  * **Priority**: ✅  
  * **File**: `internal/services/cdn/frontdoor_firewall_policy_resource.go`
  * **Details**: Good use of Optional+Computed pattern for Azure-managed defaults

  # ⛏️ Typo Detection in Documentation
  * **Priority**: 🟡
  * **File**: `README.md`
  * **Details**: Found typo: "comitted" should be "committed" (2 instances)
  * **Locations**: Lines 183, 187 - same misspelling in multiple command examples
  * **Suggested Change**: Fix spelling consistency across all instances

  # **Final Assessment:**
  Changes implement captcha support correctly with minor validation improvements needed.
  ```

## foo

* Use structured markdown format for each suggestion to ensure clarity and consistency, while adapting content to provide maximum value for the specific changes being reviewed:

  # 📋 **Code Review**: ${change_description}

  ## 📊 **CHANGE SUMMARY**
  - **Files Changed**: [number] files ([additions], [modifications], [deletions])
  - **Scale**: [insertions] insertions, [deletions] deletions
  - **Type**: [local changes/staged changes]
  - **Scope**: [Brief description of overall scope]

  ## 🎯 **PRIMARY CHANGES ANALYSIS**

  [Overview of the code changes, including the purpose of the implementation, any relevant context about the Azure service or infrastructure changes, and the files involved.]

  ## 📋 **DETAILED TECHNICAL REVIEW**

  ### 🟢 **STRENGTHS**
  [List positive aspects and well-implemented features]

  ### 🟡 **OBSERVATIONS**
  [List areas for consideration or minor improvements]

  ### 🔴 **ISSUES** (if any)
  [List any problems that need to be addressed]

  ## ✅ **RECOMMENDATIONS**

  ### 🎯 **IMMEDIATE**
  [Critical actions needed before commit/merge]

  ### 🔄 **FUTURE CONSIDERATIONS**
  [Improvements for future iterations]

  ## 🏆 **OVERALL ASSESSMENT**

  [Final recommendation with confidence level]

  ---

  ## Individual Suggestions (if needed):

  ## ⛏️ Typo Detection in Documentation
  * **Priority**: 🟡
  * **File**: `README.md`
  * **Details**: Found typo: "comitted" should be "committed" (2 instances)
  * **Locations**: Lines 183, 187 - same misspelling in multiple command examples
  * **Suggested Change**: Fix spelling consistency across all instances

  ---

  ## ${🔧/❓/⛏️/♻️/🤔/🚀/ℹ️/📌} ${Review Type}: ${Summary with necessary context}
  * **Priority**: ${🔥/🔴/🟡/🔵/⭐/✅}
  * **File**: ${relative/path/to/file}
  * **Details**: Clear explanation
  * **Azure Context** (if applicable): Service behavior reference
  * **Terraform Impact** (if applicable): Configuration/state effects
  * **Suggested Change** (if applicable): Code snippet

  # Final Assessment:

## Priority Emojis

* Use the following emojis to indicate the priority of the suggestions:
  * 🔥 Critical
  * 🔴 High
  * 🟡 Medium
  * 🔵 Low (needs attention)
  * ⭐ Notable (innovative/clever solution)
  * ✅ Correct implementation, no issues

## Type Emojis

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

## Use Code Review Emojis

Use code review emojis. Give the reviewee added context and clarity to follow up on code review. For example, knowing whether something really requires action (🔧), highlighting nit-picky comments (⛏️), flagging out of scope items for follow-up (📌)

### Type Emoji Legend

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

### Priority Emoji Legend

| `Emoji` |      `:code:`        | `Meaning`                                                                                               |
| :-----: | :------------------: | ------------------------------------------------------------------------------------------------------- |
|   🔥   |      `:fire:`        | Critical issue that must be addressed before code can be merged. Security vulnerabilities, breaking changes, or major bugs. |
|   🔴   |   `:red_circle:`     | High priority issue that should be addressed promptly. Significant bugs, incorrect patterns, or important improvements. |
|   🟡   |  `:yellow_circle:`   | Medium priority suggestion. Good improvements or standard best practices that would enhance the code quality. |
|   🔵   |   `:blue_circle:`    | Low priority item that needs attention. Minor improvements, style suggestions, or future considerations. |
|   ⭐   |      `:star:`        | Notable implementation that showcases innovative or clever solutions worth highlighting to the team. |
|   ✅   | `:white_check_mark:` | Correct implementation with no issues. Good work that follows patterns properly and requires no action. |

# Terraform Provider Specific Review Points

When reviewing Terraform AzureRM provider code, pay special attention to:

- **Code Comments Policy**: Apply thoughtful comment guidelines - prioritize code clarity through structure over explanatory comments
- **CustomizeDiff Import Requirements**: Verify correct import patterns based on implementation type
- **Resource Implementation**: Ensure proper CRUD operations, schema validation, and Azure patterns
- **Azure API Integration**: Check error handling, polling, and authentication patterns
- **State Management**: Verify drift detection, import functionality, and resource ID handling
- **Testing Standards**: Ensure comprehensive acceptance tests and proper cleanup
- **Documentation Quality**: Verify examples, attributes, and import documentation

## 🎯 **FLEXIBILITY FOR VALUABLE INSIGHTS**

While following the structured format:
- **Adapt sections** when they don't fit the specific changes being reviewed
- **Add context** that helps the developer understand Azure service implications
- **Provide alternative approaches** when current implementation could be improved
- **Explain the "why"** behind suggestions, especially for Azure-specific patterns
- **Share relevant patterns** from other parts of the codebase when applicable
- **Highlight innovative solutions** that could benefit other implementations

**Goal**: Structured consistency + contextual expertise = Maximum developer value

**📋 For detailed enforcement guidelines, see: [Code Clarity Enforcement](../instructions/code-clarity-enforcement.instructions.md)**
