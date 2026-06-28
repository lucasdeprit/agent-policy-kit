# Rules

By default, rules are YAML files loaded from the `rules/` directory of the current working directory. Use `--rules` to load policies from another directory, file, or glob.

Only regex rules are supported in the MVP.

Create starter rule files with:

```bash
agent-policy-kit init
```

This creates a documented template:

```text
rules/
  default.yaml
  examples/
    rule-template.yaml
```

`agent-policy-kit review` loads only YAML files directly under `rules/`. Files under `rules/examples/` are documentation/examples and are not evaluated.

Rule sources are automatically excluded from validation. If `--rules ./rules` is used, changes under `./rules` are ignored by `review`; if `--rules ./rules/*.yaml` is used, the matched YAML files are ignored.

Rules and target project can be independent:

```bash
agent-policy-kit review --rules ../shared-policy-kit/rules --target ../client-project
```

Pass `--rules` multiple times to combine rule sources:

```bash
agent-policy-kit review \
  --target ../client-project \
  --rules ../shared-policy-kit/rules \
  --rules ./rules
```

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

Fields:

- `id`: stable rule identifier shown in reports.
- `severity`: only `error` is supported.
- `include`: glob patterns for files to validate.
- `exclude`: optional glob patterns to skip.
- `match.type`: only `regex` is supported.
- `match.pattern`: Go regular expression evaluated against each added line.
- `expected`: short desired outcome.
- `suggestion`: short actionable fix.
