package mongo

import "github.com/XuanZhiLiu/mongo/bsontool"

func (c *Command) Set(fieldName string, document interface{}) *Command {
	return c.updateOperator("$set", fieldName, document)
}

func (c *Command) Inc(fieldName string, document interface{}) *Command {
	return c.updateOperator("$inc", fieldName, document)
}

func (c *Command) Min(fieldName string, document interface{}) *Command {
	return c.updateOperator("$min", fieldName, document)
}

func (c *Command) Max(fieldName string, document interface{}) *Command {
	return c.updateOperator("$max", fieldName, document)
}

func (c *Command) Mul(fieldName string, document interface{}) *Command {
	return c.updateOperator("$mul", fieldName, document)
}

func (c *Command) Unset(fieldName string) *Command {
	return c.updateOperator("$unset", fieldName, "")
}

func (c *Command) AddToSet(fieldName string, document interface{}) *Command {
	return c.updateOperator("$addToSet", fieldName, document)
}

// position: -1 | 1
func (c *Command) Pop(fieldName string, position int) *Command {
	return c.updateOperator("$pop", fieldName, position)
}

func (c *Command) Pull(fieldName string, document interface{}) *Command {
	return c.updateOperator("$pull", fieldName, document)
}

func (c *Command) PullAll(fieldName string, document interface{}) *Command {
	return c.updateOperator("$pullAll", fieldName, document)
}

func (c *Command) Push(fieldName string, document interface{}) *Command {
	return c.updateOperator("$push", fieldName, document)
}

func (c *Command) updateOperator(operate string, fieldName string, document interface{}) *Command {
	if c.updateModel.updateOperator == nil {
		c.updateModel.updateOperator = make(bsontool.M)
	}
	operateData := make(bsontool.M)
	if tmp, ok := c.updateModel.updateOperator[operate]; ok {
		operateData = tmp.(bsontool.M)
	}
	operateData[fieldName] = document

	c.updateModel.updateOperator[operate] = operateData
	return c
}
