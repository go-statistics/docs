# docs

[![docs](https://github.com/go-statistics/docs/actions/workflows/docs.yml/badge.svg)](https://github.com/go-statistics/docs/actions/workflows/docs.yml)
[![License](https://img.shields.io/badge/license-BSD--3--Clause-9E520F?style=flat-square)](LICENSE)

**The go-statistics documentation site**, published at
<https://go-statistics.github.io/docs/>.

A [Hugo](https://gohugo.io/) site with the [tannevaled/hextra](https://github.com/tannevaled/hextra)
fork of the [Hextra](https://github.com/imfing/hextra) theme, built the same way
as the [go-numerics documentation](https://github.com/go-numerics/docs), with
go-statistics' branding. It documents
[`go-statistics/stats`](https://github.com/go-statistics/stats).

## The pages come from the library

The library's pages are **not kept here**. Every build checks out
`go-statistics/stats` at `params.stats.version` in `hugo.yaml` and
`scripts/import` turns its `README.md` into the section's home
(`content/en/stats/_index.md`) and each `docs/*.md` into a page, with the
text unchanged:

- the leading `# Title` becomes the front matter's title;
- a link to another of those documents becomes a `relref` to its page, so the
  offline link check proves it;
- a link to any other file of the library becomes a GitHub link to that file
  **at the imported commit**;
- "Edit this page" opens the file in the library, where the text lives.

A `docs/*.md` the importer does not know still gets a page, after the known
ones. `content/en/stats/` is ignored by git; `content/en/_index.md` (the
site's home) is the only page written here.

While `params.stats.version` is `main`, a change to the library reaches the
site at the next daily run of the workflow, or at once with
`gh workflow run docs.yml -R go-statistics/docs`. A build fails if the
library's documents are missing: the home page links to each of them with a
`relref`, and Hugo refuses one that does not resolve.

## Versions

A version of the documentation is **the minor version of
[`go-statistics/stats`](https://github.com/go-statistics/stats) it describes**,
or `main` while no release is tagged:

| Source | Published under |
| --- | --- |
| branch `main`, with `params.stats.version` in `hugo.yaml` | `https://go-statistics.github.io/docs/<MAJOR.MINOR>/`, or `/docs/main/` while it says `main` |
| the newest release (else `main`) | also `https://go-statistics.github.io/docs/latest/` |

`https://go-statistics.github.io/docs/` redirects to `latest/`. Each push to
`main` replaces the directory of the version it describes, and nothing else.
The pages describe **`main`**, published under `main/`.

To document a release, set `params.stats.version` to its tag.

## Layout

| Path | What |
| --- | --- |
| `content/en/_index.md` | the site's home |
| `content/en/stats/` | the library's pages, written by `scripts/import` (not committed) |
| `scripts/import/` | the Go program that imports them |
| `i18n/en.yaml` | the site's strings, beside the theme's |
| `hugo.yaml` | `baseURL`, the stats version described, the navbar (version › theme › search › GitHub), the brand mounts |
| `layouts/_partials/custom/content-begin.html` | the outdated-version banner (reads `versions.json`) |
| `layouts/_partials/favicons.html` | the brand's favicons |
| `assets/css/custom.css` | the brand amber as Hextra's primary colour, the table headers |
| `themes/hextra` | submodule: [tannevaled/hextra](https://github.com/tannevaled/hextra), pinned to a release tag |
| `branding` | submodule: [go-statistics/brand](https://github.com/go-statistics/brand) (the mark, its PNGs and ICO) |
| `scripts/publish-version.sh` | puts one build into a `gh-pages` checkout and rewrites `versions.json`, `latest` and the root redirect |

## Working locally

```sh
git clone --recurse-submodules https://github.com/go-statistics/docs.git
git clone https://github.com/go-statistics/stats.git
cd docs
(cd scripts/import && go run . "$PWD/../../../stats" ../../content/en/stats)
mkdir -p data && go run github.com/tannevaled/hextra/cmd/pagehistory@v0.21.0 > data/pagehistory.json
hugo server --baseURL http://localhost:1313/docs/main/
```

then open <http://localhost:1313/docs/main/>. Hugo 0.146 or later and Go
1.27.2 or later; no Python, no Node.

## Publication

`.github/workflows/docs.yml`:

| Job | When | What |
| --- | --- | --- |
| `build` | pull requests, `main`, daily | imports the library, builds with `--baseURL …/docs/<version>/`, checks internal links and anchors (lychee, offline), uploads the site |
| `links` | pull requests, `main`, daily | imports the library and fetches every external URL the pages name |
| `deploy` | `main` (push, daily, manual) | `scripts/publish-version.sh` into `gh-pages` (started empty on the first deploy), then a commit and a push |

GitHub Pages serves the `gh-pages` branch.

## Licence

BSD-3-Clause.
