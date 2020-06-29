package mongo

import (
	"github.com/georgeliu825/mongo/bsontool"
)

type Pipeline []bsontool.M

func NewPipeline() (pipe Pipeline) {
	c := Pipeline{}
	return c
}

func (p Pipeline) Match(cond bsontool.D) Pipeline {
	p = append(p, bsontool.M{"$match": cond})
	return p
}

func (p Pipeline) Group(id interface{}, m bsontool.M) Pipeline {
	b := make(bsontool.M)
	b["_id"] = id
	for k, v := range m {
		b[k] = v
	}
	p = append(p, bsontool.M{"$group": b})
	return p
}

func (p Pipeline) Project(m bsontool.M) Pipeline {
	p = append(p, bsontool.M{"$project": m})
	return p
}

func (p Pipeline) Count(field string) Pipeline {
	p = append(p, bsontool.M{"$count": field})
	return p
}

func (p Pipeline) Unwind(m interface{}) Pipeline {
	p = append(p, bsontool.M{"$unwind": m})
	return p
}

func (p Pipeline) Lookup(from, localField, foreignField, as string, let bsontool.M, pipeline Pipeline) Pipeline {
	m := make(bsontool.M)
	if from != "" {
		m["from"] = from
	}
	if localField != "" {
		m["localField"] = localField
	}
	if foreignField != "" {
		m["foreignField"] = foreignField
	}
	if as != "" {
		m["as"] = as
	}
	if len(let) > 0 {
		m["let"] = let
	}
	if len(pipeline) > 0 {
		m["pipeline"] = pipeline
	}
	p = append(p, bsontool.M{"$lookup": m})
	return p
}
