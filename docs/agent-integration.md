# Agent Integration

Recommended instruction for coding agents:

```text
After completing an implementation iteration:

Run:

agent-policy-kit review

If violations are found:

- Read every violation.
- Fix every violation.
- Run the command again.
- Do not finish the task until the command succeeds.
```

`agent-policy-kit review` is deterministic. It reads the current Git diff, applies YAML rules to added lines only, and exits with `1` when violations are found.

For new repositories, run `agent-policy-kit init` once to create starter files under `rules/`.

For corporate multi-root workspaces, pass both paths explicitly:

```text
agent-policy-kit review --rules ../shared-policy-kit/rules --target ../client-project
```

Before starting a feedback loop in a new workspace, run:

```text
agent-policy-kit doctor --rules ../shared-policy-kit/rules --target ../client-project
```
