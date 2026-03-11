# runlens

`runlens` summarizes JSONL traces from agent runs, tool calls, or evaluation loops.

## Example

```bash
go run . summary -input examples/run.jsonl
```

## Use cases

- Inspect OpenClaw-style tool execution logs.
- Spot slow or flaky tools before they become production incidents.
- Export machine-readable summaries with `-json`.

## Install

From source:

```bash
go install github.com/YOUR_GITHUB_USER/runlens@latest
```

From Homebrew after you publish a tap formula:

```bash
brew tap itamaker/tap https://github.com/itamaker/homebrew-tap
brew install itamaker/tap/runlens
```

## Repo-Ready Files

- `.github/workflows/ci.yml`
- `.github/workflows/release.yml`
- `.goreleaser.yaml`
- `PUBLISHING.md`
- `scripts/render-homebrew-formula.sh`

## Release

```bash
git tag v0.1.0
git push origin v0.1.0
```

The tagged release workflow publishes multi-platform binaries and `checksums.txt`, which you can feed into the Homebrew formula renderer.
The generated formula should be committed to `https://github.com/itamaker/homebrew-tap`.
