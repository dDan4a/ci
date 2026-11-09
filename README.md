# ci

Reusable GitHub Actions workflows and reference linter configs for dDan4a repositories. One workflow per linter or check, so a
repo calls exactly the ones its stack needs.

```yaml
jobs:
  lint:
    uses: dDan4a/ci/.github/workflows/go-lint.yml@v1
```

## Workflows

Every workflow checks out the caller with `persist-credentials: false` and runs with `permissions: contents: read`.
Workflows with a `working-directory` input run there (default `.`) and look for the config there.

| Workflow | Runs | Inputs | Config: repo's own, else reference |
| --- | --- | --- | --- |
| `go-lint.yml` | golangci-lint v2.15 | `working-directory` | `.golangci.*` / [`configs/.golangci.yml`](configs/.golangci.yml) |
| `go-test.yml` | `go test -race -cover` over a Go version × OS matrix | `working-directory`, `go-versions`, `os`, `tags` | — |
| `go-test-postgres.yml` | `go test -race -cover` with a `postgres:17` service | `working-directory`, `tags`, `dsn-env` | — |
| `go-vulncheck.yml` | govulncheck on the latest patch of the go.mod minor | `working-directory` | — |
| `go-mod-tidy.yml` | `go mod tidy -diff` | `working-directory` | — |
| `qml-lint.yml` | qmllint + qmlformat check (Qt 6 from Arch Linux) | `working-directory` | `.qmllint.ini`, `.qmlformat.ini` / [`configs/`](configs) |
| `js-biome.yml` | `biome ci` (plain JS, no package.json needed) | `working-directory` | `biome.json*` / [`configs/biome.json`](configs/biome.json) |
| `ts-eslint.yml` | `npm ci`, ESLint + typescript-eslint, `tsc --noEmit` | `working-directory` | `eslint.config.*` / [`configs/eslint.config.mjs`](configs/eslint.config.mjs) |
| `python-ruff.yml` | `ruff check` + `ruff format --check` | `working-directory` | `ruff.toml`, `.ruff.toml`, `[tool.ruff]` / [`configs/ruff.toml`](configs/ruff.toml) |
| `shell-lint.yml` | shellcheck + `shfmt -d`, scripts found by extension and shebang | `working-directory` | `.shellcheckrc`, `.editorconfig` / [`configs/`](configs) |
| `markdownlint.yml` | markdownlint-cli2 | `working-directory` | `.markdownlint*` / [`configs/.markdownlint-cli2.yaml`](configs/.markdownlint-cli2.yaml) |
| `yamllint.yml` | yamllint `--strict` | `working-directory` | `.yamllint*` / [`configs/.yamllint`](configs/.yamllint) |
| `hadolint.yml` | hadolint on tracked `Dockerfile*`, `Containerfile*`, `*.dockerfile` | `working-directory` | `.hadolint.yaml` |
| `buf-lint.yml` | `buf lint` + `buf format --diff` | `working-directory` | `buf.yaml` |
| `ansible-lint.yml` | ansible-lint | `working-directory` | `.ansible-lint` |
| `gdscript-lint.yml` | GDQuest gdscript-formatter `--check` + `lint` | `working-directory` | `.editorconfig` / [`configs/.editorconfig`](configs/.editorconfig) |
| `actionlint.yml` | actionlint (with shellcheck on `run:` scripts) | — | `.github/actionlint.yaml` |
| `zizmor.yml` | zizmor, GitHub Actions security audit | — | `.github/zizmor.yml` / [`configs/zizmor.yml`](configs/zizmor.yml) |
| `gitleaks.yml` | gitleaks over the full git history | — | `.gitleaks.toml` |
| `typos.yml` | typos spell checker | — | `_typos.toml`, `typos.toml` |

The reference config is fetched from the same commit as the workflow itself (`job.workflow_sha`), so `@v1` always pairs the
v1 workflow with the v1 config.

### Inputs

| Input | Workflows | Default | Meaning |
| --- | --- | --- | --- |
| `working-directory` | all with a directory | `.` | Directory to run in, relative to the repository root. Go workflows read `go.mod` and `go.sum` there. |
| `go-versions` | `go-test` | `'[""]'` | JSON list of versions for `actions/setup-go`; `""` takes the version from `go.mod`, e.g. `'["", "stable"]'`. |
| `os` | `go-test` | `'["ubuntu-latest"]'` | JSON list of runners, e.g. `'["ubuntu-latest", "macos-latest"]'`. |
| `tags` | `go-test`, `go-test-postgres` | `''` | Build tags for `go test -tags`, e.g. `integration`. |
| `dsn-env` | `go-test-postgres` | `DATABASE_URL` | Variable that receives `postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable`. |

## Callers

Each example is a complete `.github/workflows/ci.yml`. Pick the jobs you need and combine stacks freely.

### Go

```yaml
name: ci

on:
  push:
    branches: [main]
  pull_request:

permissions:
  contents: read

jobs:
  lint:
    uses: dDan4a/ci/.github/workflows/go-lint.yml@v1
  test:
    uses: dDan4a/ci/.github/workflows/go-test.yml@v1
  vulncheck:
    uses: dDan4a/ci/.github/workflows/go-vulncheck.yml@v1
  tidy:
    uses: dDan4a/ci/.github/workflows/go-mod-tidy.yml@v1
```

Tests against Postgres, tagged `integration`, reading the DSN from `TSUKUYOMI_PG_DSN`:

```yaml
  test-postgres:
    uses: dDan4a/ci/.github/workflows/go-test-postgres.yml@v1
    with:
      tags: integration
      dsn-env: TSUKUYOMI_PG_DSN
```

A Go version and OS matrix:

```yaml
  test:
    uses: dDan4a/ci/.github/workflows/go-test.yml@v1
    with:
      go-versions: '["", "stable"]'
      os: '["ubuntu-latest", "macos-latest", "windows-latest"]'
```

### QML (Quickshell plugins)

```yaml
jobs:
  qml:
    uses: dDan4a/ci/.github/workflows/qml-lint.yml@v1
```

Quickshell and the shell's own `qs.*` modules do not resolve outside the shell, so every type from them is unknown to qmllint.
The reference [`.qmllint.ini`](configs/.qmllint.ini) turns the categories that follow from that (`ImportFailure`,
`UnresolvedType`, `MissingType`, `MissingProperty`, `UnqualifiedAccess`, `InheritanceCycle`, `TypeInstantiatedRecursively`)
into `info`: still printed, not failing. Everything else stays a warning, and any warning fails the job. The reference
[`.qmlformat.ini`](configs/.qmlformat.ini) uses 2-space indent, semicolons only where needed and a 150 column limit;
`qmlformat -i *.qml` applies it locally.

### JavaScript and TypeScript

```yaml
jobs:
  biome:
    uses: dDan4a/ci/.github/workflows/js-biome.yml@v1
```

```yaml
jobs:
  eslint:
    uses: dDan4a/ci/.github/workflows/ts-eslint.yml@v1
```

`ts-eslint` needs a `package-lock.json`. Without `eslint` in the project it installs ESLint, `@eslint/js`, typescript-eslint and
TypeScript for the run. QML-flavoured JS (`.pragma library`) is not valid JavaScript for Biome: exclude such files in `biome.json`.

### Python, shell, Markdown

```yaml
jobs:
  ruff:
    uses: dDan4a/ci/.github/workflows/python-ruff.yml@v1
  shell:
    uses: dDan4a/ci/.github/workflows/shell-lint.yml@v1
  markdown:
    uses: dDan4a/ci/.github/workflows/markdownlint.yml@v1
```

### GDScript

```yaml
jobs:
  gdscript:
    uses: dDan4a/ci/.github/workflows/gdscript-lint.yml@v1
```

gdscript-formatter reads `.editorconfig`: `max_line_length` for both wrapping and lint, and
`gdscript_formatter_exclude = true` to skip a directory (the reference skips `addons/`).

### Docker, Protobuf, Ansible

```yaml
jobs:
  hadolint:
    uses: dDan4a/ci/.github/workflows/hadolint.yml@v1
  buf:
    uses: dDan4a/ci/.github/workflows/buf-lint.yml@v1
    with:
      working-directory: proto
  ansible:
    uses: dDan4a/ci/.github/workflows/ansible-lint.yml@v1
```

### Every repo

```yaml
jobs:
  actionlint:
    uses: dDan4a/ci/.github/workflows/actionlint.yml@v1
  zizmor:
    uses: dDan4a/ci/.github/workflows/zizmor.yml@v1
  gitleaks:
    uses: dDan4a/ci/.github/workflows/gitleaks.yml@v1
  typos:
    uses: dDan4a/ci/.github/workflows/typos.yml@v1
  yamllint:
    uses: dDan4a/ci/.github/workflows/yamllint.yml@v1
```

zizmor wants every `uses:` pinned by commit SHA. The reference [`zizmor.yml`](configs/zizmor.yml) allows `dDan4a/*@v1`;
a repo with its own `.github/zizmor.yml` needs the same `unpinned-uses` policy.

## Reference configs

[`configs/`](configs) holds the reference config of every linter, line length 150 wherever it applies. Keep a copy in each repo,
so local runs match CI without the network:

```bash
curl -fsSLO https://raw.githubusercontent.com/dDan4a/ci/v1/configs/.golangci.yml
```

The Go config enables `wsl_v5`; `golangci-lint run --fix` fixes most of its findings. Its `forbidigo` rule against `errors.As`
needs Go 1.26+ for `errors.AsType`.

## Versions

- Tool versions are `*_VERSION` variables in the `env:` of each workflow; the Arch Linux image tag is in `qml-lint.yml`.
- Third-party actions are pinned by full commit SHA with the version in a comment. Dependabot proposes updates weekly.

To bump, change the version, open a PR, let `self-check` pass, merge, then move the major tag:

```bash
git tag -f v1 && git push -f origin v1
```

A change that can fail callers who did nothing (a new lint rule in a reference config, a removed or renamed input) goes to
`v2` instead.

## Development

[`self-check.yml`](.github/workflows/self-check.yml) runs on every PR: actionlint, zizmor, gitleaks, typos, yamllint,
markdownlint and shell-lint on this repo, and every other workflow against a small fixture in [`testdata/`](testdata). A new
workflow gets a fixture without its own config, so the run also exercises the fallback to the reference config.
