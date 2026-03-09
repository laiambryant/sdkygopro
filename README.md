# YGOProDeck Go SDK

A Go SDK for the YGOProDeck API.

[![Go Reference](https://pkg.go.dev/badge/github.com/laiambryant/sdkyGOprodeck.svg)](https://pkg.go.dev/github.com/laiambryant/sdkyGOprodeck)
[![Go Report Card](https://goreportcard.com/badge/github.com/laiambryant/sdkyGOprodeck)](https://goreportcard.com/report/github.com/laiambryant/sdkyGOprodeck)
[![GitHub license](https://img.shields.io/github/license/laiambryant/ygoprodeck.svg)](https://github.com/laiambryant/sdkyGOprodeck/blob/main/LICENSE)
[![GitHub issues](https://img.shields.io/github/issues/laiambryant/ygoprodeck.svg)](https://github.com/laiambryant/sdkyGOprodeck/issues)
[![GitHub stars](https://img.shields.io/github/stars/laiambryant/ygoprodeck.svg)](https://github.com/laiambryant/sdkyGOprodeck/stargazers)
[![Coverage Status](https://coveralls.io/repos/github/laiambryant/ygoprodeck/badge.svg)](https://coveralls.io/github/laiambryant/ygoprodeck)

## Installation

```bash
go get github.com/laiambryant/sdkyGOprodeck
```

## Quick Start

Create an SDK instance with [`ygoprodeck.New`](ygoprodeck.go) and call the API methods:

```go
sdk := ygoprodeck.New()

// Search for cards
q := query.New().Name("Dark Magician")
resp, err := sdk.GetCards(context.Background(), q)

// Get a random card
card, err := sdk.GetRandomCard(context.Background())

// List all card sets
sets, err := sdk.GetCardSets(context.Background())

// Get info for a specific set
info, err := sdk.GetCardSetInfo(context.Background(), "LOB")

// List all archetypes
archetypes, err := sdk.GetArchetypes(context.Background())

// Check database version
version, err := sdk.GetDBVersion(context.Background())
```

## Configuration

Pass options to [`ygoprodeck.New`](ygoprodeck.go):

- `WithBaseURL(url)` - Override API base URL
- `WithUserAgent(ua)` - Set custom User-Agent
- `WithHTTPClient(client)` - Provide custom HTTP client
- `WithCache(ttl)` - Enable response caching

## API

### SDK

- [`YGOProDeck`](ygoprodeck.go) - Main SDK type with API methods

### Methods

- `GetCards(ctx, query)` - Search cards via `/cardinfo.php`
- `GetRandomCard(ctx)` - Get a random card via `/randomcard.php`
- `GetCardSets(ctx)` - List all card sets via `/cardsets.php`
- `GetCardSetInfo(ctx, setcode)` - Get set details via `/cardsetsinfo.php`
- `GetArchetypes(ctx)` - List all archetypes via `/archetypes.php`
- `GetDBVersion(ctx)` - Get database version via `/checkDBVer.php`

### Client

- [`client.Client`](client/client.go) - HTTP client for API requests
- [`client.Option`](client/client_options.go) - Configuration options

### Endpoint

- [`endpoint.Endpoint`](endpoint/endpoint.go) - Generic endpoint with Fetch method
- [`endpoint.DecodeError`](endpoint/errors.go) - JSON decoding error

### Query Builder

Use the `query` builder to construct search parameters for `GetCards`:

```go
q := query.New().
  FuzzyName("Dragon").
  Attribute(enums.AttributeLight).
  Sort(enums.SortATK).
  Num(10).
  Offset(0)

resp, err := sdk.GetCards(context.Background(), q)
```

Passing `nil` as the query returns all cards.

#### Available Query Methods

- `Name(name)`, `FuzzyName(fname)` - Search by name
- `ID(id)`, `KonamiID(id)` - Search by ID
- `CardType(type)`, `Race(race)`, `Attribute(attr)` - Filter by card properties
- `ATK(val)`, `DEF(val)`, `Level(val)` - Filter by stats (supports operator prefixes like `gte2500`)
- `Link(val)`, `LinkMarker(markers...)`, `Scale(val)` - Filter Link/Pendulum cards
- `CardSet(set)`, `Archetype(arch)` - Filter by set or archetype
- `BanlistFilter(list)`, `Format(fmt)` - Filter by banlist or format
- `Sort(field)` - Sort results
- `Misc(bool)`, `Staple(bool)`, `HasEffect(bool)` - Misc filters
- `StartDate(date)`, `EndDate(date)`, `DateRegion(region)` - Date filters
- `Language(lang)` - Set response language
- `Num(n)`, `Offset(n)` - Pagination

### Models

- [`models.Card`](models/card.go) - Card details
- [`models.CardSet`](models/card_set.go) - Card set summary
- [`models.CardSetInfo`](models/card_set_info.go) - Card set details
- [`models.Archetype`](models/archetype.go) - Archetype
- [`models.DBVersion`](models/db_version.go) - Database version
- [`models.CardInfoResponse`](models/subs.go) - Cards response with pagination meta

### Enums

- [`enums.Language`](enums/enums.go) - Language codes
- [`enums.Attribute`](enums/enums.go) - Monster attributes
- [`enums.Format`](enums/enums.go) - Game formats
- [`enums.SortOrder`](enums/enums.go) - Sort fields
- [`enums.Banlist`](enums/enums.go) - Banlist types
