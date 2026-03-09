package ygoprodeck

import (
	"context"
	"net/http"
	"testing"

	"github.com/laiambryant/sdkyGOprodeck/client"
	"github.com/laiambryant/sdkyGOprodeck/query"
)

type fakeHTTPClient struct {
	fn func(req *http.Request) (*http.Response, error)
}

func (f *fakeHTTPClient) Do(req *http.Request) (*http.Response, error) { return f.fn(req) }

func TestNewDefaultsAndFields(t *testing.T) {
	sdk := New()
	if sdk.Client == nil {
		t.Fatalf("expected Client to be non-nil")
	}
	if sdk.Client.BaseURL != "https://db.ygoprodeck.com/api/v7" {
		t.Fatalf("unexpected default BaseURL: %s", sdk.Client.BaseURL)
	}
	if sdk.cards == nil || sdk.sets == nil || sdk.setInfo == nil || sdk.archs == nil || sdk.dbVer == nil {
		t.Fatalf("expected all endpoints to be initialized")
	}
}

func TestNewWithOptionsOverrides(t *testing.T) {
	fake := &fakeHTTPClient{fn: func(req *http.Request) (*http.Response, error) {
		return client.NewMockResponse(200, `{}`), nil
	}}
	sdk := New(client.WithBaseURL("http://example"), client.WithUserAgent("custom-agent"), client.WithHTTPClient(fake))
	if sdk.Client.BaseURL != "http://example" {
		t.Fatalf("expected base url override, got %s", sdk.Client.BaseURL)
	}
	if sdk.Client.UserAgent != "custom-agent" {
		t.Fatalf("expected user agent override, got %s", sdk.Client.UserAgent)
	}
	if sdk.Client.HTTP != fake {
		t.Fatalf("expected provided HTTP client to be used")
	}
}

func TestGetCards(t *testing.T) {
	fake := &fakeHTTPClient{fn: func(req *http.Request) (*http.Response, error) {
		if req.URL.Path != "/cardinfo.php" {
			t.Fatalf("unexpected path: %s", req.URL.Path)
		}
		if req.URL.Query().Get("name") != "Dark Magician" {
			t.Fatalf("unexpected query: %s", req.URL.RawQuery)
		}
		body := `{"data": [{"id": 46986414, "name": "Dark Magician", "type": "Normal Monster", "frameType": "normal", "desc": "The ultimate wizard.", "atk": 2500, "def": 2100, "level": 7, "race": "Spellcaster", "attribute": "DARK", "card_images": [], "card_prices": [], "ygoprodeck_url": "url"}]}`
		return client.NewMockResponse(200, body), nil
	}}
	sdk := New(client.WithBaseURL("http://example"), client.WithHTTPClient(fake))
	q := query.New().Name("Dark Magician")
	resp, err := sdk.GetCards(context.Background(), q)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resp.Cards) != 1 || resp.Cards[0].Name != "Dark Magician" {
		t.Fatalf("unexpected cards: %#v", resp.Cards)
	}
}

func TestGetCardsNilQuery(t *testing.T) {
	fake := &fakeHTTPClient{fn: func(req *http.Request) (*http.Response, error) {
		if req.URL.Path != "/cardinfo.php" {
			t.Fatalf("unexpected path: %s", req.URL.Path)
		}
		return client.NewMockResponse(200, `{"data": []}`), nil
	}}
	sdk := New(client.WithBaseURL("http://example"), client.WithHTTPClient(fake))
	resp, err := sdk.GetCards(context.Background(), nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resp.Cards) != 0 {
		t.Fatalf("expected 0 cards, got %d", len(resp.Cards))
	}
}

func TestGetCardsWithMeta(t *testing.T) {
	fake := &fakeHTTPClient{fn: func(req *http.Request) (*http.Response, error) {
		body := `{"data": [{"id": 1, "name": "Card", "type": "T", "frameType": "f", "desc": "d", "race": "r", "card_images": [], "card_prices": [], "ygoprodeck_url": "u"}], "meta": {"current_rows": 1, "total_rows": 100, "rows_remaining": 99, "total_pages": 100, "pages_remaining": 99}}`
		return client.NewMockResponse(200, body), nil
	}}
	sdk := New(client.WithBaseURL("http://example"), client.WithHTTPClient(fake))
	resp, err := sdk.GetCards(context.Background(), nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Meta == nil || resp.Meta.TotalRows != 100 {
		t.Fatalf("expected meta with total_rows 100")
	}
}

func TestGetCardsError(t *testing.T) {
	fake := &fakeHTTPClient{fn: func(req *http.Request) (*http.Response, error) {
		return client.NewMockResponse(500, "error"), nil
	}}
	sdk := New(client.WithBaseURL("http://example"), client.WithHTTPClient(fake))
	_, err := sdk.GetCards(context.Background(), nil)
	if err == nil {
		t.Fatalf("expected error")
	}
}

func TestGetRandomCard(t *testing.T) {
	fake := &fakeHTTPClient{fn: func(req *http.Request) (*http.Response, error) {
		if req.URL.Path != "/randomcard.php" {
			t.Fatalf("unexpected path: %s", req.URL.Path)
		}
		body := `{"data": [{"id": 123, "name": "Random Card", "type": "Spell Card", "frameType": "spell", "desc": "A spell.", "race": "Normal", "card_images": [], "card_prices": [], "ygoprodeck_url": "url"}]}`
		return client.NewMockResponse(200, body), nil
	}}
	sdk := New(client.WithBaseURL("http://example"), client.WithHTTPClient(fake))
	card, err := sdk.GetRandomCard(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if card.Name != "Random Card" {
		t.Fatalf("expected Random Card, got %s", card.Name)
	}
}

func TestGetRandomCardEmptyData(t *testing.T) {
	fake := &fakeHTTPClient{fn: func(req *http.Request) (*http.Response, error) {
		return client.NewMockResponse(200, `{"data": []}`), nil
	}}
	sdk := New(client.WithBaseURL("http://example"), client.WithHTTPClient(fake))
	_, err := sdk.GetRandomCard(context.Background())
	if err != client.ErrNotFound {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestGetRandomCardError(t *testing.T) {
	fake := &fakeHTTPClient{fn: func(req *http.Request) (*http.Response, error) {
		return client.NewMockResponse(500, "error"), nil
	}}
	sdk := New(client.WithBaseURL("http://example"), client.WithHTTPClient(fake))
	_, err := sdk.GetRandomCard(context.Background())
	if err == nil {
		t.Fatalf("expected error")
	}
}

func TestGetCardSets(t *testing.T) {
	fake := &fakeHTTPClient{fn: func(req *http.Request) (*http.Response, error) {
		if req.URL.Path != "/cardsets.php" {
			t.Fatalf("unexpected path: %s", req.URL.Path)
		}
		return client.NewMockResponse(200, `[{"set_name": "LOB", "set_code": "LOB", "num_of_cards": 126}]`), nil
	}}
	sdk := New(client.WithBaseURL("http://example"), client.WithHTTPClient(fake))
	sets, err := sdk.GetCardSets(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(sets) != 1 || sets[0].SetName != "LOB" {
		t.Fatalf("unexpected sets: %#v", sets)
	}
}

func TestGetCardSetsError(t *testing.T) {
	fake := &fakeHTTPClient{fn: func(req *http.Request) (*http.Response, error) {
		return client.NewMockResponse(500, "error"), nil
	}}
	sdk := New(client.WithBaseURL("http://example"), client.WithHTTPClient(fake))
	_, err := sdk.GetCardSets(context.Background())
	if err == nil {
		t.Fatalf("expected error")
	}
}

func TestGetCardSetInfo(t *testing.T) {
	fake := &fakeHTTPClient{fn: func(req *http.Request) (*http.Response, error) {
		if req.URL.Path != "/cardsetsinfo.php" {
			t.Fatalf("unexpected path: %s", req.URL.Path)
		}
		if req.URL.Query().Get("setcode") != "LOB" {
			t.Fatalf("unexpected setcode: %s", req.URL.Query().Get("setcode"))
		}
		return client.NewMockResponse(200, `[{"id": 1, "name": "Dark Magician", "set_name": "LOB", "set_code": "LOB-005", "set_rarity": "Ultra Rare", "set_price": "5.00"}]`), nil
	}}
	sdk := New(client.WithBaseURL("http://example"), client.WithHTTPClient(fake))
	info, err := sdk.GetCardSetInfo(context.Background(), "LOB")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(info) != 1 || info[0].Name != "Dark Magician" {
		t.Fatalf("unexpected set info: %#v", info)
	}
}

func TestGetCardSetInfoError(t *testing.T) {
	fake := &fakeHTTPClient{fn: func(req *http.Request) (*http.Response, error) {
		return client.NewMockResponse(500, "error"), nil
	}}
	sdk := New(client.WithBaseURL("http://example"), client.WithHTTPClient(fake))
	_, err := sdk.GetCardSetInfo(context.Background(), "LOB")
	if err == nil {
		t.Fatalf("expected error")
	}
}

func TestGetCardSetInfoEscaping(t *testing.T) {
	fake := &fakeHTTPClient{fn: func(req *http.Request) (*http.Response, error) {
		if req.URL.Query().Get("setcode") != "LOB EN" {
			t.Fatalf("unexpected setcode: %q", req.URL.Query().Get("setcode"))
		}
		return client.NewMockResponse(200, `[]`), nil
	}}
	sdk := New(client.WithBaseURL("http://example"), client.WithHTTPClient(fake))
	_, err := sdk.GetCardSetInfo(context.Background(), "LOB EN")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestGetArchetypes(t *testing.T) {
	fake := &fakeHTTPClient{fn: func(req *http.Request) (*http.Response, error) {
		if req.URL.Path != "/archetypes.php" {
			t.Fatalf("unexpected path: %s", req.URL.Path)
		}
		return client.NewMockResponse(200, `[{"archetype_name": "Blue-Eyes"}, {"archetype_name": "Dark Magician"}]`), nil
	}}
	sdk := New(client.WithBaseURL("http://example"), client.WithHTTPClient(fake))
	archs, err := sdk.GetArchetypes(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(archs) != 2 || archs[0].ArchetypeName != "Blue-Eyes" {
		t.Fatalf("unexpected archetypes: %#v", archs)
	}
}

func TestGetArchetypesError(t *testing.T) {
	fake := &fakeHTTPClient{fn: func(req *http.Request) (*http.Response, error) {
		return client.NewMockResponse(500, "error"), nil
	}}
	sdk := New(client.WithBaseURL("http://example"), client.WithHTTPClient(fake))
	_, err := sdk.GetArchetypes(context.Background())
	if err == nil {
		t.Fatalf("expected error")
	}
}

func TestGetDBVersion(t *testing.T) {
	fake := &fakeHTTPClient{fn: func(req *http.Request) (*http.Response, error) {
		if req.URL.Path != "/checkDBVer.php" {
			t.Fatalf("unexpected path: %s", req.URL.Path)
		}
		return client.NewMockResponse(200, `[{"database_version": "v1.0", "last_update": "2024-01-15"}]`), nil
	}}
	sdk := New(client.WithBaseURL("http://example"), client.WithHTTPClient(fake))
	ver, err := sdk.GetDBVersion(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ver.DatabaseVersion != "v1.0" || ver.LastUpdate != "2024-01-15" {
		t.Fatalf("unexpected version: %#v", ver)
	}
}

func TestGetDBVersionEmptyArray(t *testing.T) {
	fake := &fakeHTTPClient{fn: func(req *http.Request) (*http.Response, error) {
		return client.NewMockResponse(200, `[]`), nil
	}}
	sdk := New(client.WithBaseURL("http://example"), client.WithHTTPClient(fake))
	_, err := sdk.GetDBVersion(context.Background())
	if err != client.ErrNotFound {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestGetDBVersionError(t *testing.T) {
	fake := &fakeHTTPClient{fn: func(req *http.Request) (*http.Response, error) {
		return client.NewMockResponse(500, "error"), nil
	}}
	sdk := New(client.WithBaseURL("http://example"), client.WithHTTPClient(fake))
	_, err := sdk.GetDBVersion(context.Background())
	if err == nil {
		t.Fatalf("expected error")
	}
}
