# AI Docs Audit (optional)

This repo includes an optional, user-invoked Copilot prompt that audits a docs page under `website/docs/**` against:
- AzureRM documentation standards, and
- the provider schema under `internal/**`.

## How to use

### Copilot Chat (Agent mode) (recommended)

1. Open the docs page you want to audit (must be under `website/docs/**`).
2. Click the docs page so it is the **active editor tab**.
3. In **GitHub Copilot Chat** (Agent mode), run the prompt by slash-command name:

```text
/docs-schema-audit
```

4. Apply the **Minimal Fixes (PATCH-READY)** from the audit output.

### Alternative: run the prompt file directly

If your Copilot Chat surface does not offer prompt slash commands, you can invoke the prompt by telling Copilot Chat to follow the prompt file.

1. Open the docs page you want to audit (must be under `website/docs/**`) and make it the **active editor tab**.
2. In Copilot Chat (Agent mode), send a message that includes the prompt file path and an explicit instruction to follow it.

Example message:

```text
Follow the prompt in .github/prompts/docs-schema-audit.prompt.md
```

Alternative message (equivalent):

```text
Use .github/prompts/docs-schema-audit.prompt.md to audit the currently open docs page.
```

What should happen:
- Copilot will load the prompt file and produce a **Docs Review** report for the currently active docs page.
- If the currently active tab is not under `website/docs/**`, Copilot will refuse and tell you to open a docs page and re-run.

```text
.github/prompts/docs-schema-audit.prompt.md
```

## Troubleshooting

- The prompt audits the **currently active editor**. If your active tab is not a docs page under `website/docs/**` (for example, the prompt file itself or a README), it will refuse to run. Switch back to the docs page and re-run.

If you see the refusal message, the fix is usually:

```text
# click the docs page tab
/docs-schema-audit
```

## What it checks

- Required sections and standard ordering
- Argument/attribute ordering rules (required first, alphabetical, tags last)
- Hard-coded secrets in examples
- Schema parity (required/optional/computed + ForceNew notes)
- Common grammar and consistency issues
