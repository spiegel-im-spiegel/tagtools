# Copilot Instructions for tagtools

## Project Overview

- This repository provides a small CLI for Hugo tag metadata workflows.
- Main focus areas are `tagslist`, `toptags`, and `verify` command maintenance.
- Keep changes practical, incremental, and easy to review.

## Coding Guidelines

- Language: Go.
- Prefer small, focused changes over broad refactors.
- Preserve existing public behavior unless explicitly asked to change it.
- Keep source-code comments in English.
- Follow existing package boundaries and naming style in this repository.

## Validation

- Use Taskfile tasks for local validation.
- Primary check command:

```sh
task test
```

- If needed, run additional checks:

```sh
task govulncheck
```

- For routine code changes, prefer this full check:

```sh
task
```

## Dependencies

- Do not add new external dependencies unless there is a clear benefit.
- If a dependency is added, explain why in the change summary.

## Documentation and Maintenance

- Update related documentation when behavior, options, or workflows change.
- Keep README and package comments consistent with the current implementation.
- Prefer explicit notes about assumptions and constraints when they are not obvious from code.

## Current Maintenance Context

- Ongoing maintenance includes migration follow-up from `spiegel-im-spiegel/github-pages-env`, workflow cleanup, and documentation refresh.
- During this phase, prioritize output compatibility (`tagslist.csv` and `toptags.json`) and risk reduction over feature expansion.
