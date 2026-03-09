package ygoprodeck

import (
	"context"
	"net/url"

	"github.com/laiambryant/sdkyGOprodeck/client"
	"github.com/laiambryant/sdkyGOprodeck/endpoint"
	"github.com/laiambryant/sdkyGOprodeck/models"
	"github.com/laiambryant/sdkyGOprodeck/query"
)

type YGOProDeck struct {
	Client  *client.Client
	cards   *endpoint.Endpoint[models.CardResponse]
	sets    *endpoint.Endpoint[[]models.CardSet]
	setInfo *endpoint.Endpoint[[]models.CardSetInfo]
	archs   *endpoint.Endpoint[[]models.Archetype]
	dbVer   *endpoint.Endpoint[[]models.DBVersion]
}

func New(opts ...client.Option) *YGOProDeck {
	c := client.NewHTTPClient(nil, opts...)
	return &YGOProDeck{
		Client:  c,
		cards:   endpoint.New[models.CardResponse](c),
		sets:    endpoint.New[[]models.CardSet](c),
		setInfo: endpoint.New[[]models.CardSetInfo](c),
		archs:   endpoint.New[[]models.Archetype](c),
		dbVer:   endpoint.New[[]models.DBVersion](c),
	}
}

func (y *YGOProDeck) GetCards(ctx context.Context, q *query.Query) (models.CardInfoResponse, error) {
	qs := ""
	if q != nil {
		qs = q.Build()
	}
	resp, err := y.cards.Fetch(ctx, "/cardinfo.php"+qs)
	if err != nil {
		return models.CardInfoResponse{}, err
	}
	return models.CardInfoResponse{Cards: resp.Data, Meta: resp.Meta}, nil
}

func (y *YGOProDeck) GetRandomCard(ctx context.Context) (models.Card, error) {
	resp, err := y.cards.Fetch(ctx, "/randomcard.php")
	if err != nil {
		return models.Card{}, err
	}
	if len(resp.Data) == 0 {
		return models.Card{}, client.ErrNotFound
	}
	return resp.Data[0], nil
}

func (y *YGOProDeck) GetCardSets(ctx context.Context) ([]models.CardSet, error) {
	return y.sets.Fetch(ctx, "/cardsets.php")
}

func (y *YGOProDeck) GetCardSetInfo(ctx context.Context, setcode string) ([]models.CardSetInfo, error) {
	return y.setInfo.Fetch(ctx, "/cardsetsinfo.php?setcode="+url.QueryEscape(setcode))
}

func (y *YGOProDeck) GetArchetypes(ctx context.Context) ([]models.Archetype, error) {
	return y.archs.Fetch(ctx, "/archetypes.php")
}

func (y *YGOProDeck) GetDBVersion(ctx context.Context) (models.DBVersion, error) {
	versions, err := y.dbVer.Fetch(ctx, "/checkDBVer.php")
	if err != nil {
		return models.DBVersion{}, err
	}
	if len(versions) == 0 {
		return models.DBVersion{}, client.ErrNotFound
	}
	return versions[0], nil
}
