# runlens

[![All Contributors](https://img.shields.io/badge/all_contributors-1-orange.svg?style=flat-square)](#contributors-)

`runlens` is a Go CLI that summarizes JSONL traces from agent runs, tool calls, and evaluation loops.

It helps you spot slow tools, flaky execution paths, and token-heavy runs before they turn into production issues.

## Support

[![Buy Me A Coffee](https://img.shields.io/badge/Buy%20Me%20A%20Coffee-FFDD00?style=for-the-badge&logo=buy-me-a-coffee&logoColor=black)](https://buymeacoffee.com/amaker)

## Quickstart

### Install

Install with your preferred method:

```bash
# From the custom tap
brew tap itamaker/tap https://github.com/itamaker/homebrew-tap
brew install itamaker/tap/runlens
```

```bash
# Or install from source
go install github.com/itamaker/runlens@latest
```

<details>
<summary>You can also download binaries from <a href="https://github.com/itamaker/runlens/releases">GitHub Releases</a>.</summary>

Current release archives:

- macOS (Apple Silicon/arm64): `runlens_0.1.0_darwin_arm64.tar.gz`
- macOS (Intel/x86_64): `runlens_0.1.0_darwin_amd64.tar.gz`
- Linux (arm64): `runlens_0.1.0_linux_arm64.tar.gz`
- Linux (x86_64): `runlens_0.1.0_linux_amd64.tar.gz`

Each archive contains a single executable: `runlens`.

</details>

If the repository is still private, release-based installs require GitHub access to the repository assets.

### First Run

Run:

```bash
runlens summary -input examples/run.jsonl
```

## Requirements

- Go `1.22+`

## Run

```bash
go run . summary -input examples/run.jsonl
```

Machine-readable output:

```bash
go run . summary -input examples/run.jsonl -json
```

## Build From Source

```bash
make build
```

```bash
go build -o dist/runlens .
```

## What It Does

1. Parses JSONL event streams from agent or tool executions.
2. Computes aggregate latency, success rate, and token totals.
3. Produces per-tool summaries for failure analysis.
4. Exports either human-readable output or JSON for automation.

## Notes

- `examples/run.jsonl` is a good shape reference for your own logs.
- Maintainer release steps live in `PUBLISHING.md`.

## Contributors ✨

| [![Zhaoyang Jia][avatar-zhaoyang]][author-zhaoyang] |
| --- |
| [Zhaoyang Jia][author-zhaoyang] |



[author-zhaoyang]: https://github.com/itamaker
[avatar-zhaoyang]: https://images.weserv.nl/?url=https://github.com/itamaker.png&h=120&w=120&fit=cover&mask=circle&maxage=7d
