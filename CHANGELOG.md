# Changelog

## v0.1.0

Initial MVP release of Agent Policy Kit.

### Features

- `agent-policy-kit init` creates starter rule files.
- `agent-policy-kit review` validates the current Git working-tree diff.
- YAML rule files loaded from `rules/*.yaml` and `rules/*.yml`.
- Regex-based rule matching.
- Diff-only validation of added lines.
- Compact agent-readable output.
- Embedded documented rule template.
- Embedded generic default rule and documented starter template.
- MIT License with commercial use allowed.

### Installation

```bash
go install github.com/lucasdeprit/agent-policy-kit/cmd/agent-policy-kit@v0.1.0
```

Prebuilt binaries are available from the GitHub Release assets.
