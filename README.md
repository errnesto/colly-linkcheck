# colly-linkcheck

Linkcheker that crawls ems.press and checks for dead links.

Right now the url and the ignore patterns are hardcoded.
Maybe we change that later so that can be customized.

## use as github action

```yaml
jobs:
    linkcheck:
        runs-on: ubuntu-latest
        name: check links
        steps:
        - name: Linkcheck
            # pin to current commit
            uses: ems-press/colly-linkcheck
```