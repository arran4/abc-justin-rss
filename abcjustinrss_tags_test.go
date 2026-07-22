package abcjustinrss

import (
	"reflect"
	"strings"
	"testing"

	"github.com/PuerkitoBio/goquery"
)

func TestExtractCategories(t *testing.T) {
	html := `
	<article>
		<a data-component="CombinedTag" href="/news/topics">
			<span>Death and Dying</span>, <span>Business, Economics and Finance</span> and <span>Australia</span>
		</a>
		<a data-component="CombinedTag" href="/news/topics2">
			<span data-component="ScreenReaderOnly">Topic:</span>
			<span>Death and Dying</span>, <span>Travel and Tourism</span> and <span>Actor</span>
		</a>
		<a data-component="CombinedTag" href="/news/topics3">
			<p>Analysis by Annabel Crabb</p>
		</a>
		<a data-component="CombinedTag" href="/news/topics4">
			Tag1, Tag2 and Tag3
		</a>
		<a data-component="SubjectTag" href="/news/topic/death-and-dying">
			Death and Dying
		</a>
	</article>
	`
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		t.Fatalf("Failed to parse HTML: %v", err)
	}

	s := doc.Find("article").First()
	var categories []string
	seenCategories := make(map[string]bool)

	s.Find("a[data-component='SubjectTag'], a[data-component='Tag'], a[data-component='CombinedTag']").Each(func(j int, a *goquery.Selection) {
		clone := a.Clone()
		clone.Find("[data-component=\"ScreenReaderOnly\"], [data-component=\"fallbackText\"], [data-component=\"ScreenReaderTimestamp\"]").Remove()
		text := strings.TrimSpace(clone.Text())
		if text != "" {
			isCombined := a.AttrOr("data-component", "") == "CombinedTag" || a.HasClass("CombinedTag")
			if isCombined {
				var extractedTags []string
				clone.Contents().Each(func(k int, n *goquery.Selection) {
					if goquery.NodeName(n) == "#text" {
						tStr := strings.TrimSpace(n.Text())
						if tStr == "," || tStr == ", " || tStr == "and" {
							return
						}
						if tStr != "" {
							extractedTags = append(extractedTags, tStr)
						}
					} else {
						tStr := strings.TrimSpace(n.Text())
						if tStr != "" {
							extractedTags = append(extractedTags, tStr)
						}
					}
				})

				if len(extractedTags) == 0 {
					extractedTags = append(extractedTags, text)
				}

				// fallback for pure text strings like "Tag1, Tag2 and Tag3"
				if len(extractedTags) == 1 && strings.Contains(extractedTags[0], " and ") {
					tStr := strings.ReplaceAll(extractedTags[0], " and ", ", ")
					parts := strings.Split(tStr, ", ")
					extractedTags = []string{}
					for _, p := range parts {
						p = strings.TrimSpace(p)
						if p != "" {
							extractedTags = append(extractedTags, p)
						}
					}
				}

				for _, p := range extractedTags {
					if p != "" && !seenCategories[p] {
						seenCategories[p] = true
						categories = append(categories, p)
					}
				}
			} else {
				if !seenCategories[text] {
					seenCategories[text] = true
					categories = append(categories, text)
				}
			}
		}
	})

	expected := []string{
		"Death and Dying", "Business, Economics and Finance", "Australia",
		"Travel and Tourism", "Actor", "Analysis by Annabel Crabb",
		"Tag1", "Tag2", "Tag3",
	}
	if !reflect.DeepEqual(categories, expected) {
		t.Errorf("Expected %v, got %v", expected, categories)
	}
}
