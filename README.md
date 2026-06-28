# Agent Policy Kit

Deterministic policy tooling for coding agents.

Agent Policy Kit provides the `agent-policy-kit` CLI. It validates the current Git working-tree diff, applies declarative YAML rules to added lines, and returns compact feedback designed for agents such as Copilot, Codex, Claude Code, and OpenCode.

The tool does not use AI, external services, AST plugins, or IDE integrations. The MVP is intentionally small: regex rules, YAML configuration, and two commands.

## Why

Prompts tell an agent how to work. Guardrails verify that the generated code follows concrete policies.

The intended loop is:

1. The agent implements a change.
2. The agent runs `agent-policy-kit review`.
3. If violations exist, the agent reads the feedback, fixes the code, and runs the command again.
4. The task is not finished until the command succeeds.

## Quick Start

```bash
go install github.com/lucasdeprit/agent-policy-kit/cmd/agent-policy-kit@latest
```

From a repository that contains a `rules/` directory:

```bash
agent-policy-kit review
```

For corporate multi-root workspaces, keep rules and target independent:

```bash
agent-policy-kit review --rules ../shared-policy-kit/rules --target ../client-project
```

To create starter rules in a repository:

```bash
agent-policy-kit init
```

This creates:

```text
rules/
  default.yaml
  examples/
    rule-template.yaml
```

Existing files are skipped, not overwritten.

Validate setup before reviewing:

```bash
agent-policy-kit doctor --rules ../shared-policy-kit/rules --target ../client-project
```

Exit codes:

- `0`: no violations.
- `1`: violations found or invalid configuration.

## Output

```text
FAILED 1

GEN001 src/app.ts:42
+ console.log("debug")
expected: Use the approved abstraction.
fix: Replace the raw API with the project-approved alternative.
```

When there are no violations:

```text
OK
```

## Rules

By default, rules are loaded from `rules/*.yaml` and `rules/*.yml` in the current working directory. Use `--rules` to load policies from another directory, file, or glob.

The engine is language-agnostic. Files under `rules/examples/` are examples for humans and are not loaded by `review`.

Rule sources are automatically excluded from validation. This lets rule files contain examples of forbidden patterns without causing `review` to fail on the rules themselves.

```yaml
rules:
  - id: GEN001
    name: no_debug_console_logs
    severity: error
    include:
      - "src/**/*.ts"
    exclude:
      - "src/generated/**"
    match:
      type: regex
      pattern: "\\bconsole\\.log\\s*\\("
    message: "Do not commit debug console logging."
    expected: "Use the project logger or remove the log."
    suggestion: "Replace console.log(...) or delete it before finishing."
```

## Documentation

- [Rule format](docs/rules.md)
- [Agent integration](docs/agent-integration.md)
- [Installation and distribution](docs/installation.md)
- [Contributing](CONTRIBUTING.md)
- [License](LICENSE)
- [Changelog](CHANGELOG.md)
- [Example diff](examples/debug-log.diff)

## Project Layout

```text
cmd/agent-policy-kit/   CLI entrypoint
internal/diff/          Git working-tree diff reader
internal/engine/        Rule evaluation engine
internal/parser/        YAML rule loader and validator
internal/reporter/      Compact agent-readable output
internal/rules/         Rule data model
rules/                  Example rule packs
docs/                   Project documentation
examples/               Example inputs and outputs
```

## MVP Scope

Implemented:

- `agent-policy-kit init` to create starter rules.
- `agent-policy-kit review` to validate the current diff.
- `agent-policy-kit doctor` to validate target/rules setup.
- Independent `--target` and repeated `--rules` flags.
- Review only current Git working-tree changes.
- Detect added lines from staged, unstaged, and untracked text files.
- YAML rules.
- Regex matchers.
- Compact text output for agents.

Not implemented yet:

- AST rules.
- Quick fixes.
- IDE plugins.
- Remote services.
- AI-based validation.

## Author

Created by Lucas Deprit.

## License

MIT License. You can use, copy, modify, distribute, sublicense, and sell copies of this software, including for commercial purposes, as long as the license notice is preserved.
