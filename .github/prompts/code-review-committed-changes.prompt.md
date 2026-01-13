---
description: "Code Review Prompt for Terraform AzureRM Provider Committed Changes"
---

# 🚀 **EXECUTE IMMEDIATELY** - Code Review Task

**PRIMARY TASK**: Perform code review of committed changes for Terraform AzureRM provider

## ⚡ **START HERE - MANDATORY EXECUTION STEPS**

**1. GET THE DIFF - Run these git commands immediately:**
```powershell
# Get branch and overview
git branch --show-current
git --no-pager diff --stat --no-prefix origin/main...HEAD

# Get the focused diff for review (exclude generated/vendor files)
git --no-pager diff --no-prefix origin/main...HEAD -- ":(exclude)vendor/**" ":(exclude)go.sum" ":(exclude)go.mod"

# Get commit context
git log --oneline origin/main..HEAD
git status
```

**⚠️ IMPORTANT**: If the commands do not show any changes, abandon the code review and display:
**"☠️ Argh! Shiver me source files! This branch be cleaner than a swabbed deck! Push some code, Ye Lily-livered scallywag! ☠️"**

**2. REVIEW THE CHANGES** - Apply expertise as principal Terraform provider engineer

**3. PROVIDE STRUCTURED FEEDBACK** - Use the review format below

---

## 🎯 **CORE REVIEW MISSION**

As a principal Terraform provider engineer with expertise in Go development, Azure APIs, and HashiCorp Plugin SDK, deliver actionable code review feedback.

**Critical Issues**:
- Security vulnerabilities in Azure authentication and API calls
- Resource lifecycle bugs (create, read, update, delete operations)
- State management and drift detection issues
- Azure API error handling and retry logic
- Resource import functionality correctness
- Terraform Plugin SDK usage violations
- Implementation approach consistency (Typed vs Untyped Resource Implementation)
- CustomizeDiff import pattern correctness (conditional import requirements)
- Resource schema validation and type safety

**Code Quality**:
- Go language conventions and idiomatic patterns
- Terraform resource implementation best practices
- Azure SDK for Go usage patterns
- Plugin SDK schema definitions and validation
- CustomizeDiff function validation logic and patterns
- Implementation approach appropriateness (Typed for new, Untyped maintenance only)
- Error handling and context propagation
- Resource timeout configurations
- Acceptance test coverage and quality
- **Tests use ONLY ExistsInAzure() check with ImportStep() - NO redundant field validation**
- **CRITICAL: Code comments policy enforcement - Comments only for Azure API quirks, complex business logic, or SDK workarounds that cannot be expressed through code structure**

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

---

##  **REVIEW OUTPUT FORMAT**

```markdown
# 📋 **Code Review**: ${change_description}

## 📊 **CHANGE SUMMARY**
- **Files Changed**: [number] files ([additions], [modifications], [deletions])
- **Scale**: [insertions] insertions, [deletions] deletions
- **Branch**: [current_branch_from_git_command] vs origin/main
- **Scope**: [Brief description of overall scope]

## 🎯 **PRIMARY CHANGES ANALYSIS**
[Overview of the code changes and purpose]

## 📋 **DETAILED TECHNICAL REVIEW**

### 🟢 **STRENGTHS**
[List positive aspects and well-implemented features]

### 🟡 **OBSERVATIONS**
[List areas for consideration or minor improvements]

### 🔴 **ISSUES** (if any)
[List ONLY actual problems that need to be fixed - bugs, errors, violations, missing requirements, typos, misspellings, improper pointer handling, incorrect SDK usage, using deprecated utilities, etc. Do NOT include observations about what was done correctly or opinions about changes that are already implemented properly]

## ✅ **RECOMMENDATIONS**

### 🎯 **IMMEDIATE**
[Critical actions needed before merge]

### 🔄 **FUTURE CONSIDERATIONS**
[Improvements for future iterations]

## 🏆 **OVERALL ASSESSMENT**
[Final recommendation with confidence level]

---

## Individual Suggestions Format:

## ${🔧/❓/⛏️/♻️/🤔/🚀/ℹ️/📌} ${Review Type}: ${Summary}
* **Priority**: ${🔥/🔴/🟡/🔵/⭐/✅}
* **File**: ${relative/path/to/file}
* **Details**: Clear explanation
* **Azure Context** (if applicable): Service behavior reference
* **Terraform Impact** (if applicable): Configuration/state effects
* **Suggested Change** (if applicable): Code snippet

# Summary
Concise assessment and any follow-up items.
```

**Priority System:** 🔥 Critical → 🔴 High → 🟡 Medium → 🔵 Low → ⭐ Notable → ✅ Good

**Review Type Emojis:**
* 🔧 Change request - Functional issues requiring fixes
* ❓ Question - Clarification needed about design decisions
* ⛏️ Nitpick - Minor style/consistency issues (typos, formatting, naming)
* ♻️ Refactor suggestion - Structural code improvements
* 🤔 Thought/concern - Design or approach concerns requiring discussion
* 🚀 Positive feedback - Excellent implementations worth highlighting
* ℹ️ Explanatory note - Technical context or background information
* 📌 Future consideration - Larger scope items for follow-up

---

# 📚 **APPENDIX: EDGE CASE HANDLING** *(Secondary Guidelines)*

## Console Line Wrapping Detection *(If Needed)*

**⚠️ ONLY IF git diff content appears corrupted:**
- Use `read_file filename` to verify actual content
- Note: `*(Verified: console wrapping - content clean)*`
- Continue with technical review

**Console artifacts are normal** - Focus on delivering valuable code review feedback.

## Git Requirements

**Git Commands to Execute:**
```powershell
git branch --show-current
git --no-pager diff --stat --no-prefix origin/main...HEAD
git --no-pager diff --no-prefix origin/main...HEAD
git log --oneline origin/main..HEAD
git status
```

**If no changes found:**
**"☠️ Argh! Shiver me source files! This branch be cleaner than a swabbed deck! Push some code, Ye Lily-livered scallywag! ☠️"**

## Verification Protocol *(Edge Cases Only)*

**When to verify:**
- Text breaks mid-word without logical reason
- Missing quotes/brackets that don't make contextual sense
- Emojis appear as `??`
- JSON/YAML looks syntactically broken

**Verification format:**
```markdown
*(Verified: console wrapping - actual content clean)*
```

## Review Scope Expansion

**Beyond diff changes, also check:**
- Spelling and grammar in visible text
- Command examples accuracy
- Naming consistency
- Professional language standards

## Comprehensive Quality Guidelines

- **Code Comments Policy**: Comments only for Azure API quirks, complex business logic, or SDK workarounds that cannot be expressed through code structure
- **Comment Quality**: All comments must have clear justification and add genuine value beyond code structure
- **Refactoring Preference**: Consider if code restructuring could eliminate need for comments
- **Documentation Standards**:
  - Spelling accuracy in all text content
  - Grammar and syntax correctness
  - Consistent terminology and naming
  - Professional language standards

## Provider-Specific Excellence

- **Testing Standards**: ExistsInAzure() + ImportStep() only, no redundant field validation
- **CustomizeDiff Patterns**: Correct imports based on implementation type
- **Azure Patterns**: PATCH operations, "None" value handling, SDK integration
- **Implementation Approach**: Typed for new resources, Untyped for maintenance only

---

**REMEMBER: PRIMARY MISSION is to deliver actionable technical feedback. All appendix items are secondary safeguards.**
