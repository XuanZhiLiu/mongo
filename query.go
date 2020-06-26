package mongo

import (
	"github.com/XuanZhiLiu/mongo/bsontool"
	"github.com/XuanZhiLiu/mongo/internal/glob"
)

// enable：由client自行判斷是否需要帶此條件
// val若為該型別的預設值則不會帶入。若需要使用該型別之預設型別請使用指標，例如：bool的預設值false將不會帶入條件
func (c *Command) CustomConditions(cond bsontool.D) *Command {
	c.condition = c.condition.CustomConditions(cond)
	return c
}

func (c *Command) Equals(enable bool, vals ...bsontool.KeyValue) *Command {
	c.condition = c.condition.Equals(enable, vals...)
	return c
}

func (c *Command) Equal(col string, enable bool, val interface{}) *Command {
	c.condition = c.condition.Equal(col, enable, val)
	return c
}

func (c *Command) NotEqual(col string, enable bool, val interface{}) *Command {
	c.condition = c.condition.NotEqual(col, enable, val)
	return c
}

func (c *Command) Exists(col string, enable, exists bool) *Command {
	c.condition = c.condition.Exists(col, enable, exists)
	return c
}

func (c *Command) In(col string, enable bool, vals ...interface{}) *Command {
	c.condition = c.condition.In(col, enable, vals...)
	return c
}

func (c *Command) InSlice(col string, enable bool, slice interface{}) *Command {
	if !enable {
		return c
	}
	iSlice, ok := glob.ToSlice(slice)
	if ok {
		c.condition = c.condition.In(col, enable, iSlice...)
	}
	return c
}

func (c *Command) NotIn(col string, enable bool, vals ...interface{}) *Command {
	c.condition = c.condition.NotIn(col, enable, vals...)
	return c
}

func (c *Command) All(col string, enable bool, vals ...interface{}) *Command {
	c.condition = c.condition.All(col, enable, vals...)
	return c
}

func (c *Command) SizeOfArray(col string, enable bool, size int) *Command {
	c.condition = c.condition.SizeOfArray(col, enable, size)
	return c
}

func (c *Command) Regex(col string, enable bool, pattern, options string) *Command {
	c.condition = c.condition.Regex(col, enable, pattern, options)
	return c
}

func (c *Command) Between(col string, enable bool, gte, gt, lte, lt interface{}) *Command {
	c.condition = c.condition.Between(col, enable, gte, gt, lte, lt)
	return c
}

func (c *Command) GreaterThanOrEqual(col string, enable bool, val interface{}) *Command {
	return c.Between(col, enable, val, nil, nil, nil)
}

func (c *Command) GreaterThan(col string, enable bool, val interface{}) *Command {
	return c.Between(col, enable, nil, val, nil, nil)
}

func (c *Command) LessThanOrEqual(col string, enable bool, val interface{}) *Command {
	return c.Between(col, enable, nil, nil, val, nil)
}

func (c *Command) LessThan(col string, enable bool, val interface{}) *Command {
	return c.Between(col, enable, nil, nil, nil, val)
}

func (c *Command) Or(conds ...interface{}) *Command {
	c.condition = c.condition.Or(conds...)
	return c
}

func (c *Command) And(conds ...interface{}) *Command {
	c.condition = c.condition.And(conds...)
	return c
}
