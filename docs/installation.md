# Installation and Distribution

Agent Policy Kit ships the `agent-policy-kit` Go CLI. It can be installed from source or distributed as a prebuilt executable.

## Install With Go

```bash
go install github.com/lucasdeprit/agent-policy-kit/cmd/agent-policy-kit@latest
```

The binary is installed into:

- macOS/Linux: `$(go env GOPATH)/bin/agent-policy-kit`
- Windows: `%USERPROFILE%\go\bin\agent-policy-kit.exe`

That directory must be available in `PATH`.

Then initialize rules inside the repository you want to validate:

```bash
agent-policy-kit init
```

After editing the generated rules, run:

```bash
agent-policy-kit review
```

For a multi-root workspace where policies live outside the target project:

```bash
agent-policy-kit review --rules ../shared-policy-kit/rules --target ../client-project
```

Validate the setup with:

```bash
agent-policy-kit doctor --rules ../shared-policy-kit/rules --target ../client-project
```

## Prebuilt Executables

GitHub Releases publish one archive per platform:

- `agent-policy-kit-darwin-arm64.tar.gz`
- `agent-policy-kit-darwin-amd64.tar.gz`
- `agent-policy-kit-linux-arm64.tar.gz`
- `agent-policy-kit-linux-amd64.tar.gz`
- `agent-policy-kit-windows-amd64.zip`
- `checksums.txt`

Users can place the executable in a directory from `PATH`, for example:

- macOS/Linux: `/usr/local/bin/agent-policy-kit`
- Windows: `C:\Tools\agent-policy-kit\agent-policy-kit.exe`

## Where Rules Live

For the MVP, rules should live inside each target repository:

```text
my-app/
  rules/
    default.yaml
```

This keeps reviews reproducible because the rules are versioned with the code they validate.

## Local Tooling Option

In restricted corporate environments, teams can also commit or provision the binary under a project-local tool directory:

```text
my-app/
  tools/
    agent-policy-kit
  rules/
    default.yaml
```

Then run:

```bash
./tools/agent-policy-kit review
```

The recommended default remains a globally installed binary plus repository-local rules.

## Versioned Releases

Releases follow semantic version tags such as `v0.1.0`.

Install a specific version with:

```bash
go install github.com/lucasdeprit/agent-policy-kit/cmd/agent-policy-kit@v0.1.0
```

Release assets are built automatically when a `v*` tag is pushed.
