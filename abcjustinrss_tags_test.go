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
		<a data-component="CombinedTag" href="/news/topic/test">
			Tag1, Tag2 and Tag3
		</a>
		<a data-component="SubjectTag" href="/news/topic/death-and-dying">
			Death and Dying
		</a>
		<a data-component="SubjectTag" href="/news/topic/test">
			Tag1
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
				text = strings.ReplaceAll(text, " and ", ", ")
				parts := strings.Split(text, ", ")
				for _, p := range parts {
					p = strings.TrimSpace(p)
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

	expected := []string{"Tag1", "Tag2", "Tag3", "Death and Dying"}
	if !reflect.DeepEqual(categories, expected) {
		t.Errorf("Expected %v, got %v", expected, categories)
	}
}
