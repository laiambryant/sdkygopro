package query

import (
	"fmt"
	"net/url"
	"strings"
	"testing"

	"github.com/laiambryant/sdkyGOprodeck/enums"
)

func TestBuildEmpty(t *testing.T) {
	q := New()
	if got := q.Build(); got != "" {
		t.Fatalf("expected empty string for empty query, got %q", got)
	}
}

func TestChainingReturnsSamePointer(t *testing.T) {
	q := New()
	q2 := q.Name("test")
	if q != q2 {
		t.Fatalf("expected methods to return same pointer receiver")
	}
}

func TestAllMethodsBuildsExpectedQuery(t *testing.T) {
	q := New()
	q.Name("Dark Magician")
	q.FuzzyName("Dark")
	q.ID(46986414)
	q.KonamiID(1234)
	q.CardType("Normal Monster")
	q.Race("Spellcaster")
	q.Attribute(enums.AttributeDark)
	q.ATK("gte2500")
	q.DEF("lte2100")
	q.Level("7")
	q.Link(3)
	q.LinkMarker("Top", "Bottom")
	q.Scale(8)
	q.CardSet("LOB")
	q.Archetype("Blue-Eyes")
	q.BanlistFilter(enums.BanlistTCG)
	q.Format(enums.FormatTCG)
	q.Sort(enums.SortATK)
	q.Misc(true)
	q.Staple(true)
	q.HasEffect(true)
	q.StartDate("2020-01-01")
	q.EndDate("2024-12-31")
	q.DateRegion("tcg_date")
	q.Language(enums.LanguageFr)
	q.Num(10)
	q.Offset(20)
	pairs := [][2]string{
		{"name", "Dark Magician"},
		{"fname", "Dark"},
		{"id", "46986414"},
		{"konami_id", "1234"},
		{"type", "Normal Monster"},
		{"race", "Spellcaster"},
		{"attribute", "dark"},
		{"atk", "gte2500"},
		{"def", "lte2100"},
		{"level", "7"},
		{"link", "3"},
		{"linkmarker", "Top,Bottom"},
		{"scale", "8"},
		{"cardset", "LOB"},
		{"archetype", "Blue-Eyes"},
		{"banlist", "tcg"},
		{"format", "tcg"},
		{"sort", "atk"},
		{"misc", "yes"},
		{"staple", "yes"},
		{"has_effect", "true"},
		{"startdate", "2020-01-01"},
		{"enddate", "2024-12-31"},
		{"dateregion", "tcg_date"},
		{"language", "fr"},
		{"num", "10"},
		{"offset", "20"},
	}
	var parts []string
	for _, p := range pairs {
		parts = append(parts, fmt.Sprintf("%s=%s", url.QueryEscape(p[0]), url.QueryEscape(p[1])))
	}
	expected := "?" + strings.Join(parts, "&")
	if got := q.Build(); got != expected {
		t.Fatalf("unexpected query string:\n got: %s\nwant: %s", got, expected)
	}
}

func TestQueryEscaping(t *testing.T) {
	q := New()
	q.Name("a b&c=d")
	expected := "?" + fmt.Sprintf("%s=%s", url.QueryEscape("name"), url.QueryEscape("a b&c=d"))
	if got := q.Build(); got != expected {
		t.Fatalf("escaping mismatch: got %q want %q", got, expected)
	}
}

func TestHasEffectFalse(t *testing.T) {
	q := New().HasEffect(false)
	expected := "?has_effect=false"
	if got := q.Build(); got != expected {
		t.Fatalf("expected %q, got %q", expected, got)
	}
}

func TestMiscFalseNoOp(t *testing.T) {
	q := New().Misc(false)
	if got := q.Build(); got != "" {
		t.Fatalf("expected empty for Misc(false), got %q", got)
	}
}

func TestStapleFalseNoOp(t *testing.T) {
	q := New().Staple(false)
	if got := q.Build(); got != "" {
		t.Fatalf("expected empty for Staple(false), got %q", got)
	}
}
