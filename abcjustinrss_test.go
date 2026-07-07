package abcjustinrss

import (
	"strings"
	"testing"
	"time"

	"github.com/PuerkitoBio/goquery"
)

func TestExtractPubDate(t *testing.T) {
	tests := []struct {
		name     string
		html     string
		expected string // We'll just check if it parses correctly and produces expected RFC1123
	}{
		{
			name: "With datetime attribute",
			html: `<article><h3><a href="/news/2026/07/07/test">Test</a></h3>
	<div data-component="Typography">Desc</div>
	<a data-component="Link" href="/news/2026/07/07/test">Link</a>
	<time class="ScreenReaderOnly_srOnly__bnJwm" data-component="ScreenReaderTimestamp" datetime="2026-07-07T04:34:03.000Z">Tue 7 Jul 2026 at 2:34pm</time>
	</article>`,
			expected: "Tue, 07 Jul 2026 04:34:03 UTC",
		},
		{
			name: "Without datetime attribute fallback to relative time",
			html: `<article><h3><a href="/news/2026/07/07/test">Test</a></h3>
	<div data-component="Typography">Desc</div>
	<a data-component="Link" href="/news/2026/07/07/test">Link</a>
	<time>15 minutes ago</time>
	</article>`,
			// Expected will vary because time.Now() is used in DeducePublicationDate, so we just check it's not empty and parsable
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc, err := goquery.NewDocumentFromReader(strings.NewReader(tt.html))
			if err != nil {
				t.Fatalf("Failed to parse HTML: %v", err)
			}

			s := doc.Find("article").First()

			var pubDate string
			dateTime, exists := s.Find("time[datetime]").Attr("datetime")
			if exists {
				parsedTime, err := time.Parse(time.RFC3339, dateTime)
				if err == nil {
					pubDate = parsedTime.Format(time.RFC1123)
				}
			}

			if pubDate == "" {
				relativeTime := strings.TrimSpace(s.Find("time").Text())
				pubDate = DeducePublicationDate(relativeTime)
			}

			if pubDate == "" {
				t.Errorf("pubDate is empty")
			}

			if tt.expected != "" && pubDate != tt.expected {
				t.Errorf("expected pubDate %q, got %q", tt.expected, pubDate)
			}

			_, err = time.Parse(time.RFC1123, pubDate)
			if err != nil {
				t.Errorf("pubDate %q does not match RFC1123 format: %v", pubDate, err)
			}
		})
	}
}
