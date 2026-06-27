# Contributing

Thanks for considering a contribution to Agent Policy Kit.

The project goal is to stay deterministic, small, and useful for coding-agent feedback loops. Prefer simple changes that keep the engine language-agnostic.

## Principles

- No AI calls inside the tool.
- No external services at runtime.
- Keep the CLI predictable and scriptable.
- Keep rule packs separate from the engine.
- Prefer declarative rules over hardcoded framework behavior.
- Optimize output for agents: compact, actionable, and stable.

## Local Setup

```bash
go test ./...
go build ./cmd/agent-policy-kit
```

The generated local binary is ignored by Git.

## Before Opening A Pull Request

Run:

```bash
gofmt -w ./cmd ./internal
go test ./...
go build ./cmd/agent-policy-kit
```

Also check that generated files, local binaries, editor settings, and OS files are not included.

## Adding Rules

Framework-specific rules should live as YAML rule packs, not as Go code.

Good examples:

- `rules/react-design-system.yaml`
- `rules/swift-design-system.yaml`

Avoid adding framework-specific checks directly to `internal/engine`.

## Adding Match Types

The MVP supports only `regex`. Future match types should preserve the same high-level flow:

1. Read the current Git diff.
2. Extract added lines or structured changed regions.
3. Evaluate declarative rules.
4. Report compact, actionable violations.

Do not change the reporter format casually. Agents may rely on its stability.

## License

By contributing, you agree that your contributions are licensed under the MIT License.
