package rewriter

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/beevik/etree"
)

// loadFixture reads the shared Spotify feed fixture. All rewriter tests
// start from the same input so regressions stay comparable.
func loadFixture(t *testing.T) []byte {
	t.Helper()
	path := filepath.Join("testdata", "spotify-feed.xml")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("load fixture: %v", err)
	}
	return data
}

// defaultBranding returns a realistic branding config for tests.
func defaultBranding() Branding {
	return Branding{
		SelfURL:     "https://feed.yourdomain.com/v1/your-show.xml",
		Title:       "Your Show",
		Link:        "https://yourdomain.com/podcast",
		Image:       "https://yourdomain.com/podcast/cover.jpg",
		ITunesAuthor:     "Your Name",
		ITunesOwnerEmail: "you@yourdomain.com",
	}
}

// TestRewrite_PreservesITunesNamespace confirms that unknown namespaces and
// iTunes/podcast-2.0 elements are preserved through a round-trip.
//
// If etree drops the podcast: namespace or the itunes:category element, the
// rewritten feed is broken for Apple Podcasts and we've regressed the core
// promise of v0.1.0.
func TestRewrite_PreservesITunesNamespace(t *testing.T) {
	input := loadFixture(t)
	r := New(defaultBranding())

	output, err := r.Rewrite(input)
	if err != nil {
		t.Fatalf("rewrite: %v", err)
	}

	s := string(output)

	// Namespace declarations MUST survive the round-trip.
	for _, ns := range []string{
		`xmlns:itunes="http://www.itunes.com/dtds/podcast-1.0.dtd"`,
		`xmlns:podcast="https://podcastindex.org/namespace/1.0"`,
		`xmlns:content="http://purl.org/rss/1.0/modules/content/"`,
	} {
		if !strings.Contains(s, ns) {
			t.Errorf("namespace dropped on round-trip: %s", ns)
		}
	}

	// Unknown namespaced element MUST survive.
	if !strings.Contains(s, "<podcast:guid>abc123-def456-7890</podcast:guid>") {
		t.Errorf("podcast:guid element was not preserved")
	}
	if !strings.Contains(s, `<podcast:locked owner="original@example.com">no</podcast:locked>`) {
		t.Errorf("podcast:locked element was not preserved")
	}

	// iTunes category MUST survive (Apple Podcasts requires it).
	if !strings.Contains(s, `<itunes:category text="Technology"/>`) {
		t.Errorf("itunes:category was not preserved verbatim")
	}
}

// TestRewrite_DoesNotRewriteEnclosureURL is the most important test in the
// entire codebase. If this fails, v0.1.0 ships broken: the enclosure URL is
// the upstream audio host, and rewriting it would force us to rehost audio
// bytes, which the scope explicitly punts to v1.0+.
func TestRewrite_DoesNotRewriteEnclosureURL(t *testing.T) {
	input := loadFixture(t)
	r := New(defaultBranding())

	output, err := r.Rewrite(input)
	if err != nil {
		t.Fatalf("rewrite: %v", err)
	}

	s := string(output)

	originalURLs := []string{
		"https://anchor.fm/s/abc123/podcast/play/episode-1.mp3",
		"https://anchor.fm/s/abc123/podcast/play/episode-2.mp3",
	}
	for _, u := range originalURLs {
		if !strings.Contains(s, u) {
			t.Errorf("enclosure URL was rewritten or lost: %s missing from output", u)
		}
	}

	// Paranoia check: the branded domain must not appear in any enclosure.
	// Parse the output and inspect every <item><enclosure url> directly.
	doc := etree.NewDocument()
	if err := doc.ReadFromString(s); err != nil {
		t.Fatalf("reparse output: %v", err)
	}

	for _, item := range doc.FindElements("//item") {
		for _, enc := range item.SelectElements("enclosure") {
			url := enc.SelectAttrValue("url", "")
			if strings.Contains(url, "yourdomain.com") {
				t.Errorf("enclosure URL was rewritten to branded domain: %s", url)
			}
			if !strings.Contains(url, "anchor.fm") {
				t.Errorf("enclosure URL lost its original host: %s", url)
			}
		}
	}
}

// TestRewrite_AtomSelfLinkRewritten covers two cases: an existing atom:self
// link is rewritten to the new self URL, and a missing atom:self link is
// injected. Apple Podcasts requires this element — if we drop it, feeds
// stop validating.
func TestRewrite_AtomSelfLinkRewritten(t *testing.T) {
	// Case 1: existing atom:self link is rewritten.
	t.Run("existing link rewritten", func(t *testing.T) {
		input := loadFixture(t)
		r := New(defaultBranding())

		output, err := r.Rewrite(input)
		if err != nil {
			t.Fatalf("rewrite: %v", err)
		}

		s := string(output)

		if strings.Contains(s, "https://anchor.fm/s/abc123/podcast/rss") {
			t.Errorf("original atom:self URL leaked into output")
		}
		if !strings.Contains(s, `href="https://feed.yourdomain.com/v1/your-show.xml"`) {
			t.Errorf("new atom:self URL not present in output")
		}
		if !strings.Contains(s, `rel="self"`) {
			t.Errorf("atom:self rel attribute missing after rewrite")
		}
	})

	// Case 2: missing atom:self link is injected.
	t.Run("missing link injected", func(t *testing.T) {
		// Strip the atom:link from the fixture.
		input := loadFixture(t)
		doc := etree.NewDocument()
		if err := doc.ReadFromBytes(input); err != nil {
			t.Fatalf("parse fixture: %v", err)
		}
		channel := doc.FindElement("//channel")
		for _, link := range channel.SelectElements("atom:link") {
			channel.RemoveChild(link)
		}
		stripped, err := doc.WriteToBytes()
		if err != nil {
			t.Fatalf("rewrite fixture: %v", err)
		}

		r := New(defaultBranding())
		output, err := r.Rewrite(stripped)
		if err != nil {
			t.Fatalf("rewrite stripped: %v", err)
		}

		s := string(output)
		if !strings.Contains(s, `href="https://feed.yourdomain.com/v1/your-show.xml"`) {
			t.Errorf("atom:self link was not injected when missing")
		}
		if !strings.Contains(s, `rel="self"`) {
			t.Errorf("injected atom:self link missing rel attribute")
		}
	})
}
