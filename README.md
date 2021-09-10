# GH Analytics

CLI tool which analyzes Github event data for 1 hour.

[![Build Status](https://github.com/dikaeinstein/ghanalytics/actions/workflows/ci.yml/badge.svg?branch=main)](https://github.com/dikaeinstein/ghanalytics/actions)
[![Coverage Status](https://coveralls.io/repos/github/dikaeinstein/ghanalytics/badge.svg?branch=main)](https://coveralls.io/github/dikaeinstein/ghanalytics?branch=main)

## Build From Source

Prerequisites

- make
- GO 1.17+

```bash
git clone https://github.com/dikaeinstein/ghanalytics

cd ghanalytics
make fetch
make build
```

Run `./ghanalytics -help | -h` to get help and see available options

### Available Commands

- `top10Users` — Top 10 active users sorted by amount of PRs created and commits pushed.
- `top10ReposByCommitsPushed` — Top 10 repositories sorted by amount of commits pushed.
- `top10ReposByWatchEvents` — Top 10 repositories sorted by amount of watch events.
