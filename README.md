# colly-linkcheck

Linkcheker that crawls a webpage and checks for dead links.

The crawler will check external links but will only parse content of pages with the same host as the provided start url.

## usage (command line)

```bash
# parse complete page
go run main.go --url "https://ems.press"

# exclude path patterns. E.g. all paths startung with journals
go run main.go --url "https://ems.press" --exclude "^\/journals*"

# exclude multiple path patterns, but also include sub patterns:
go run main.go \
    --url "https://ems.press" \
    \
    --exclude "^\/journals\/.*\/articles.*" \
	--exclude "^\/journals\/.*\/issues.*" \
	--exclude "^\/books\/.*\/.*" \
    \
	--include "^\/journals\/msl\/articles.*" \
	--include "^\/journals\/msl\/issues.*" \
	--include "^\/books\/esiam.*"
```

## use as github action

```yaml
jobs:
    linkcheck:
        runs-on: ubuntu-latest
        name: check links
        steps:
        - name: Linkcheck
            # pin to current commit
            uses: ems-press/colly-linkcheck@main
```