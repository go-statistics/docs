---
title: "go-statistics documentation"
linkTitle: "Home"
type: docs
cascade:
  type: docs
description: "Probability distributions, descriptive statistics and statistical tests in pure Go (CGO=0), on go-ndarray."
---

**go-statistics is probability and statistics in pure Go.** Its library,
[stats]({{< relref "stats" >}}), has 32 univariate distributions (24
continuous, 8 discrete) and 3 multivariate ones, each with the methods of
SciPy's frozen distributions; descriptive statistics with NumPy's quantiles;
the statistical tests; the Gaussian kernel density estimate and the
bootstrap.

Every distribution function is measured in ULPs against
[mpmath](https://mpmath.org), with SciPy measured on the same points, and the
per-region budgets are enforced by `go test` on eight architectures. Where
the library is less accurate or slower than SciPy, its pages say so.

## The pages

The pages under **stats** are the library's own documentation — its
`README.md` and `docs/` — imported from
[go-statistics/stats](https://github.com/go-statistics/stats) when this site
is built, with the text unchanged:

| Page | What it holds |
| --- | --- |
| [Overview]({{< relref "stats" >}}) | examples with the SciPy call each one is, what is in the package, accuracy, draws and speed in brief |
| [Design]({{< relref "stats/design" >}}) | one package, distributions written once in standard form, the accuracy techniques, draws, fits, arrays, the three layers of verification |
| [Accuracy]({{< relref "stats/accuracy" >}}) | how the error is measured against mpmath, the headline per family, where SciPy is more accurate, the negative controls |
| [Accuracy table]({{< relref "stats/accuracy-table" >}}) | every function and every region, beside SciPy |
| [Draws]({{< relref "stats/draws" >}}) | per distribution, the draws compared with SciPy's for the same seed and how many differ |
| [Performance]({{< relref "stats/performance" >}}) and [its table]({{< relref "stats/performance-table" >}}) | against SciPy's vectorised calls, and where it is slower and why |

## Install

```sh
go get github.com/go-statistics/stats
```

The module needs Go 1.27.2 or later and is pure Go (`CGO_ENABLED=0`).

## Built on

[go-ndarray](https://github.com/go-ndarray/ndarray) (arrays, linear algebra,
and `random`, NumPy's `Generator` bit for bit),
[go-numerics/special](https://github.com/go-numerics/special) and
[go-numerics/optimize](https://github.com/go-numerics/optimize). Results plot
with [go-plot](https://github.com/go-plot/plot), and a
[go-notebook](https://github.com/go-notebook/notebook) cell can use all of it.
