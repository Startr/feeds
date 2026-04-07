// Package rewriter implements the DOM-style XML rewriter for RSS podcast
// feeds using beevik/etree.
//
// It implements the rewrite scope table from the v0.1.0 engineering review:
//
//	REWRITE:
//	  <channel><title>
//	  <channel><link>
//	  <channel><image> and <channel><image><url>
//	  <atom:link rel="self">            (injected if missing)
//	  <itunes:author>
//	  <itunes:owner><itunes:name>       (if branding.ITunesAuthor set)
//	  <itunes:owner><itunes:email>
//	  <itunes:image>
//
//	LEAVE ALONE:
//	  <item><enclosure url>             — upstream audio host
//	  <item><guid>, <title>, <description>, <pubDate>
//	  Any unknown namespaced element    — round-tripped via etree
package rewriter

import (
	"fmt"

	"github.com/beevik/etree"
)

// Branding holds the values that will be injected into the rewritten feed.
// Empty strings mean "leave that element alone".
type Branding struct {
	SelfURL          string // injected into <atom:link rel="self">
	Title            string // <channel><title>
	Link             string // <channel><link>
	Image            string // <channel><image><url> and <itunes:image href>
	ITunesAuthor     string // <itunes:author> and <itunes:owner><itunes:name>
	ITunesOwnerEmail string // <itunes:owner><itunes:email>
}

const (
	nsItunes = "http://www.itunes.com/dtds/podcast-1.0.dtd"
	nsAtom   = "http://www.w3.org/2005/Atom"
)

// Rewriter holds branding config and rewrites XML bytes in place.
type Rewriter struct {
	Branding Branding
}

// New returns a rewriter configured with the given branding.
func New(b Branding) *Rewriter {
	return &Rewriter{Branding: b}
}

// Rewrite parses the input XML, applies the rewrite scope table, and returns
// the rewritten XML as bytes. Unknown namespaces and enclosure URLs are
// preserved untouched.
func (r *Rewriter) Rewrite(input []byte) ([]byte, error) {
	doc := etree.NewDocument()
	if err := doc.ReadFromBytes(input); err != nil {
		return nil, fmt.Errorf("parse feed: %w", err)
	}

	rss := doc.SelectElement("rss")
	if rss == nil {
		return nil, fmt.Errorf("no <rss> root element")
	}

	channel := rss.SelectElement("channel")
	if channel == nil {
		return nil, fmt.Errorf("no <channel> under <rss>")
	}

	// Ensure the xmlns:atom declaration exists on <rss> so we can inject an
	// atom:link without breaking the namespace contract.
	if rss.SelectAttr("xmlns:atom") == nil {
		rss.CreateAttr("xmlns:atom", nsAtom)
	}

	// --- channel-level rewrites ---------------------------------------------

	if r.Branding.Title != "" {
		setChildText(channel, "title", r.Branding.Title)
	}
	if r.Branding.Link != "" {
		setChildText(channel, "link", r.Branding.Link)
	}

	// <channel><image><url>...</url></channel> — RSS 2.0 image block
	if r.Branding.Image != "" {
		img := channel.SelectElement("image")
		if img == nil {
			img = channel.CreateElement("image")
		}
		setChildText(img, "url", r.Branding.Image)
		// RSS spec: image also has <title> and <link> referencing channel.
		if r.Branding.Title != "" {
			setChildText(img, "title", r.Branding.Title)
		}
		if r.Branding.Link != "" {
			setChildText(img, "link", r.Branding.Link)
		}
	}

	// --- atom:link rel="self" ----------------------------------------------
	// Apple Podcasts requires atom:self to match the feed's public URL.
	// Inject if missing, rewrite if present.
	if r.Branding.SelfURL != "" {
		rewriteAtomSelfLink(channel, r.Branding.SelfURL)
	}

	// --- iTunes rewrites ----------------------------------------------------

	if r.Branding.ITunesAuthor != "" {
		setChildText(channel, "itunes:author", r.Branding.ITunesAuthor)

		// itunes:owner is a wrapper with itunes:name + itunes:email children.
		owner := channel.SelectElement("itunes:owner")
		if owner == nil {
			owner = channel.CreateElement("itunes:owner")
		}
		setChildText(owner, "itunes:name", r.Branding.ITunesAuthor)
		if r.Branding.ITunesOwnerEmail != "" {
			setChildText(owner, "itunes:email", r.Branding.ITunesOwnerEmail)
		}
	}

	if r.Branding.Image != "" {
		// itunes:image carries the URL as an href attribute, not text.
		itunesImg := channel.SelectElement("itunes:image")
		if itunesImg == nil {
			itunesImg = channel.CreateElement("itunes:image")
		}
		itunesImg.CreateAttr("href", r.Branding.Image)
	}

	// --- items are LEFT ALONE ----------------------------------------------
	// Intentional no-op for each <item>. enclosure url, guid, title,
	// description, pubDate, and any unknown namespaced element are
	// preserved verbatim by etree's round-trip.

	doc.Indent(2)
	return doc.WriteToBytes()
}

// setChildText finds a direct child element by tag (including namespaced
// tags like "itunes:author") and sets its text. If the child doesn't exist,
// it's created.
func setChildText(parent *etree.Element, tag, text string) {
	child := parent.SelectElement(tag)
	if child == nil {
		child = parent.CreateElement(tag)
	}
	child.SetText(text)
}

// rewriteAtomSelfLink finds the <atom:link rel="self"> element in the
// channel and rewrites its href to selfURL. If no such element exists,
// one is injected as a direct child of <channel>.
func rewriteAtomSelfLink(channel *etree.Element, selfURL string) {
	// Look for an existing atom:link with rel="self".
	for _, link := range channel.SelectElements("atom:link") {
		if rel := link.SelectAttrValue("rel", ""); rel == "self" {
			link.CreateAttr("href", selfURL)
			link.CreateAttr("rel", "self")
			link.CreateAttr("type", "application/rss+xml")
			return
		}
	}

	// Not found — inject a new one. Position within <channel> doesn't
	// affect Apple Podcasts validation; it just needs to exist.
	link := channel.CreateElement("atom:link")
	link.CreateAttr("href", selfURL)
	link.CreateAttr("rel", "self")
	link.CreateAttr("type", "application/rss+xml")
}
