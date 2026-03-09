package query

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/laiambryant/sdkyGOprodeck/enums"
)

type Query struct {
	params []param
}

type param struct {
	key   string
	value string
}

func New() *Query {
	return &Query{}
}

func (q *Query) add(key, value string) {
	q.params = append(q.params, param{key, value})
}

func (q *Query) Name(name string) *Query {
	q.add("name", name)
	return q
}

func (q *Query) FuzzyName(fname string) *Query {
	q.add("fname", fname)
	return q
}

func (q *Query) ID(id int) *Query {
	q.add("id", fmt.Sprintf("%d", id))
	return q
}

func (q *Query) KonamiID(id int) *Query {
	q.add("konami_id", fmt.Sprintf("%d", id))
	return q
}

func (q *Query) CardType(ct string) *Query {
	q.add("type", ct)
	return q
}

func (q *Query) Race(race string) *Query {
	q.add("race", race)
	return q
}

func (q *Query) Attribute(attr enums.Attribute) *Query {
	q.add("attribute", string(attr))
	return q
}

func (q *Query) ATK(val string) *Query {
	q.add("atk", val)
	return q
}

func (q *Query) DEF(val string) *Query {
	q.add("def", val)
	return q
}

func (q *Query) Level(val string) *Query {
	q.add("level", val)
	return q
}

func (q *Query) Link(val int) *Query {
	q.add("link", fmt.Sprintf("%d", val))
	return q
}

func (q *Query) LinkMarker(markers ...string) *Query {
	q.add("linkmarker", strings.Join(markers, ","))
	return q
}

func (q *Query) Scale(val int) *Query {
	q.add("scale", fmt.Sprintf("%d", val))
	return q
}

func (q *Query) CardSet(set string) *Query {
	q.add("cardset", set)
	return q
}

func (q *Query) Archetype(arch string) *Query {
	q.add("archetype", arch)
	return q
}

func (q *Query) BanlistFilter(list enums.Banlist) *Query {
	q.add("banlist", string(list))
	return q
}

func (q *Query) Format(f enums.Format) *Query {
	q.add("format", string(f))
	return q
}

func (q *Query) Sort(field enums.SortOrder) *Query {
	q.add("sort", string(field))
	return q
}

func (q *Query) Misc(yes bool) *Query {
	if yes {
		q.add("misc", "yes")
	}
	return q
}

func (q *Query) Staple(yes bool) *Query {
	if yes {
		q.add("staple", "yes")
	}
	return q
}

func (q *Query) HasEffect(yes bool) *Query {
	if yes {
		q.add("has_effect", "true")
	} else {
		q.add("has_effect", "false")
	}
	return q
}

func (q *Query) StartDate(date string) *Query {
	q.add("startdate", date)
	return q
}

func (q *Query) EndDate(date string) *Query {
	q.add("enddate", date)
	return q
}

func (q *Query) DateRegion(region string) *Query {
	q.add("dateregion", region)
	return q
}

func (q *Query) Language(lang enums.Language) *Query {
	q.add("language", string(lang))
	return q
}

func (q *Query) Num(n int) *Query {
	q.add("num", fmt.Sprintf("%d", n))
	return q
}

func (q *Query) Offset(n int) *Query {
	q.add("offset", fmt.Sprintf("%d", n))
	return q
}

func (q *Query) Build() string {
	if len(q.params) == 0 {
		return ""
	}
	var parts []string
	for _, p := range q.params {
		k := url.QueryEscape(p.key)
		v := url.QueryEscape(p.value)
		parts = append(parts, fmt.Sprintf("%s=%s", k, v))
	}
	return "?" + strings.Join(parts, "&")
}
