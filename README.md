# GitHub Action PR Changelog Add

A GitHub Action for adding the files changed for every PR when the PR is made by creating an issue.

## Usage

Add the following GitHub workflow to your repository.

```yaml
name: PR Changelog Add
on:
  pull_request:
    types:
    - opened
    - edited
    - synchronized
    - labeled
    - unlabeled
jobs:
  check_pr_size:
    runs-on: ubuntu-latest
    steps:
    - uses: Guzzler/gh-changelog@v1.0.0
      env:
        GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
```
## License

[MIT License](./LICENSE)

Copyright (c) 2020 Sharang Pai <sharangpai123@gmail.com>
