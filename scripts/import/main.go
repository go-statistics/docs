// Command import turns go-statistics/stats' own documentation into the pages
// of this site:
//
//	go run ./scripts/import <stats checkout> content/en/stats
//
// The library's README.md becomes the section's home and each docs/*.md a
// page, with the text unchanged. Only three things are done to it:
//
//   - the leading "# Title" line becomes the front matter's title (Hextra
//     prints it), with a weight that orders the sidebar and a description;
//   - a link to another of these documents becomes a relref to its page,
//     with its fragment kept, so the offline link check proves it;
//   - a link to any other file of the repository becomes a link to that
//     file on GitHub AT THE IMPORTED COMMIT, so it cannot later point at a
//     file that has moved.
//
// A comment at the top of each page names the file and the commit it came
// from, and its "Edit this page" link opens that file in the library. A docs/*.md this command does not know still gets a page, titled by
// its own heading, after the known ones: a new document in the library is
// published without anyone editing this file.
//
// The output directory is replaced whole. It is not committed (.gitignore):
// the site is built from the library, never from a copy of it.
package main

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

const repo = "https://github.com/go-statistics/stats"

// page is one imported document.
type page struct {
	src    string // path in the stats repository
	slug   string // page name under the section; "" for the section's home
	title  string // "" to take the document's own first heading
	weight int
	desc   string
}

// known are the documents of stats at the time of writing, in the order the
// sidebar shows them. The titles are the documents' own headings, except
// where a generated table has none.
var known = []page{
	{"README.md", "", "stats", 10, "Probability distributions, descriptive statistics and statistical tests in pure Go, measured in ULPs against mpmath."},
	{"docs/design.md", "design", "", 10, "One package; distributions written once in standard form; the accuracy techniques; draws, fits, arrays and the three layers of verification."},
	{"docs/accuracy.md", "accuracy", "", 20, "How the error is measured against mpmath, the headline per family, where SciPy is more accurate, the techniques and the negative controls."},
	{"docs/accuracy-table.md", "accuracy-table", "Accuracy table", 30, "Every function and every region: the maximum error in ULPs against mpmath, beside SciPy's on the same points."},
	{"docs/rvs-table.md", "draws", "Draws", 40, "For each distribution, how many draws were compared with SciPy's for the same seed, how many differ and by how much."},
	{"docs/perf.md", "performance", "", 50, "Against SciPy's vectorised calls on the same data: the method, the results, and where it is slower and why."},
	{"docs/perf-table.md", "performance-table", "Performance table", 60, "Every case measured: SciPy, Go on one goroutine and Go on all cores."},
}

var (
	heading = regexp.MustCompile(`^#\s+(.+?)\s*$`)
	// A Markdown link or image target: ](target) or ](target "title").
	link = regexp.MustCompile(`\]\(([^)\s]+)((?:\s+"[^"]*")?)\)`)
)

func main() {
	if len(os.Args) != 3 {
		fmt.Fprintln(os.Stderr, "usage: import <stats checkout> <output directory>")
		os.Exit(2)
	}
	src, out := os.Args[1], os.Args[2]
	if err := run(src, out); err != nil {
		fmt.Fprintln(os.Stderr, "import:", err)
		os.Exit(1)
	}
}

func run(src, out string) error {
	rev, err := exec.Command("git", "-C", src, "rev-parse", "HEAD").Output()
	if err != nil {
		return fmt.Errorf("%s is not a git checkout of stats: %w", src, err)
	}
	commit := strings.TrimSpace(string(rev))

	pages, err := discover(src)
	if err != nil {
		return err
	}
	// Where each document is published, for the links between them.
	where := map[string]string{}
	for _, p := range pages {
		where[p.src] = path.Join("stats", p.slug)
	}

	if err := os.RemoveAll(out); err != nil {
		return err
	}
	if err := os.MkdirAll(out, 0o755); err != nil {
		return err
	}
	for _, p := range pages {
		b, err := os.ReadFile(filepath.Join(src, filepath.FromSlash(p.src)))
		if err != nil {
			return err
		}
		body, title := strip(string(b))
		if p.title == "" {
			p.title = title
		}
		if p.title == "" {
			return fmt.Errorf("%s has no heading and no title is known for it", p.src)
		}
		body = relink(body, p.src, commit, where)

		var fm strings.Builder
		fmt.Fprintf(&fm, "---\ntitle: %q\nweight: %d\n", p.title, p.weight)
		if p.desc != "" {
			fmt.Fprintf(&fm, "description: %q\n", p.desc)
		}
		// "Edit this page" opens the library's file, where the text lives.
		fmt.Fprintf(&fm, "editURL: %q\n", repo+"/edit/main/"+p.src)
		fm.WriteString("---\n\n")
		fmt.Fprintf(&fm, "<!-- Imported from go-statistics/stats %s, %s, by scripts/import: the text is the library's, with its links pointed into this site. Change it there. -->\n\n", commit[:12], p.src)

		name := p.slug + ".md"
		if p.slug == "" {
			name = "_index.md"
		}
		if err := os.WriteFile(filepath.Join(out, name), []byte(fm.String()+body), 0o644); err != nil {
			return err
		}
		fmt.Printf("%-24s -> %s\n", p.src, filepath.Join(out, name))
	}
	fmt.Printf("imported %d documents from go-statistics/stats %s\n", len(pages), commit[:12])
	return nil
}

// discover lists the known documents that exist, then every other
// docs/*.md, alphabetically, after them. A missing README is an error: a
// checkout without one is not the library.
func discover(src string) ([]page, error) {
	if _, err := os.Stat(filepath.Join(src, "README.md")); err != nil {
		return nil, fmt.Errorf("no README.md in %s", src)
	}
	seen := map[string]bool{}
	var pages []page
	for _, p := range known {
		if _, err := os.Stat(filepath.Join(src, filepath.FromSlash(p.src))); err == nil {
			pages = append(pages, p)
			seen[p.src] = true
		}
	}
	others, err := filepath.Glob(filepath.Join(src, "docs", "*.md"))
	if err != nil {
		return nil, err
	}
	sort.Strings(others)
	w := 100
	for _, o := range others {
		rel := "docs/" + filepath.Base(o)
		if seen[rel] {
			continue
		}
		pages = append(pages, page{src: rel, slug: strings.TrimSuffix(filepath.Base(o), ".md"), weight: w})
		w += 10
	}
	return pages, nil
}

// strip takes the first "# Title" line off the document, if it comes before
// any other text, and returns it as the title.
func strip(s string) (body, title string) {
	lines := strings.SplitAfter(s, "\n")
	for i, l := range lines {
		t := strings.TrimSpace(l)
		if t == "" || strings.HasPrefix(t, "<!--") {
			continue
		}
		if m := heading.FindStringSubmatch(t); m != nil {
			return strings.Join(append(lines[:i:i], lines[i+1:]...), ""), m[1]
		}
		break
	}
	return s, ""
}

// relink rewrites the relative links of a document at from (a path in the
// stats repository). Absolute URLs and in-page fragments are left alone.
// Links inside fenced code blocks are not links, and are not touched.
func relink(s, from, commit string, where map[string]string) string {
	lines := strings.SplitAfter(s, "\n")
	fenced := false
	for i, l := range lines {
		if strings.HasPrefix(strings.TrimSpace(l), "```") {
			fenced = !fenced
			continue
		}
		if fenced {
			continue
		}
		lines[i] = link.ReplaceAllStringFunc(l, func(m string) string {
			sub := link.FindStringSubmatch(m)
			target, title := sub[1], sub[2]
			if strings.Contains(target, "://") || strings.HasPrefix(target, "#") ||
				strings.HasPrefix(target, "mailto:") {
				return m
			}
			file, frag, _ := strings.Cut(target, "#")
			resolved := path.Clean(path.Join(path.Dir(from), file))
			if p, ok := where[resolved]; ok {
				ref := p
				if frag != "" {
					ref += "#" + frag
				}
				return fmt.Sprintf(`]({{< relref %q >}}%s)`, ref, title)
			}
			url := repo + "/blob/" + commit + "/" + resolved
			if frag != "" {
				url += "#" + frag
			}
			return "](" + url + title + ")"
		})
	}
	return strings.Join(lines, "")
}
