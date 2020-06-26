package mongo

import (
	"context"
	"math"
)

type Page struct {
	Index int64 `json:"index" form:"pi"`
	Size  int64 `json:"size" form:"ps"`
	Total int64 `json:"total"`
	Count int64 `json:"pages"`
}

func (c *Command) setPaging(ctx context.Context) {
	if ctx != nil && ctx.Value("pi") != nil && ctx.Value("ps") != nil {
		pageIndex := ctx.Value("pi").(int64)
		pageSize := ctx.Value("ps").(int64)
		// 若pageIndex為0則預設為1
		if pageIndex == 0 {
			pageIndex = 1
		}
		skip := (pageIndex - 1) * pageSize
		c.skip = &skip
		c.limit = &(pageSize)
		c.isPaging = true
	}
}

func (c *Command) pageTotal() (page *Page, err error) {
	if !c.isPaging || c.skip == nil || c.limit == nil {
		return
	}
	skip := *c.skip
	limit := *c.limit
	c.skip = nil
	c.limit = nil
	var count int64
	count, err = c.Count()
	if err != nil {
		return
	}
	pageSize := limit
	page = &Page{
		Total: count,
		Index: (skip / limit) + 1,
		Count: int64(math.Ceil(float64(count) / float64(pageSize))),
		Size:  pageSize,
	}
	return
}
