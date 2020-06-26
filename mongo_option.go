package mongo

import (
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/event"
)

type MongoOption func(c *MongoDB) error

func SetMongoAuthMechanism(authMechanism string) MongoOption {
	return func(c *MongoDB) (err error) {
		if c.auth != nil {
			c.auth.AuthMechanism = authMechanism
		} else {
			err = errors.New("user info not found")
		}
		return
	}
}

func SetMongoAuthSource(authSource string) MongoOption {
	return func(c *MongoDB) (err error) {
		if c.auth != nil {
			c.auth.AuthSource = authSource
		} else {
			err = errors.New("user info not found")
		}
		return
	}
}

func SetMongoPoolMonitor(m *event.PoolMonitor) MongoOption {
	return func(c *MongoDB) (err error) {
		c.options.SetPoolMonitor(m)
		return
	}
}

func SetMongoMaxConnIdleTime(d time.Duration) MongoOption {
	return func(c *MongoDB) (err error) {
		c.options.SetMaxConnIdleTime(d)
		return
	}
}

func SetMongoMaxPoolSize(u uint64) MongoOption {
	return func(c *MongoDB) (err error) {
		c.options.SetMaxPoolSize(u)
		return
	}
}

func SetMongoMinPoolSize(u uint64) MongoOption {
	return func(c *MongoDB) (err error) {
		c.options.SetMinPoolSize(u)
		return
	}
}

func SetMongoMonitor(e *event.CommandMonitor) MongoOption {
	return func(c *MongoDB) (err error) {
		c.options.SetMonitor(e)
		return
	}
}

func SetLogMode(enable bool) MongoOption {
	return func(c *MongoDB) (err error) {
		if enable {
			c.logMode = defaultLogMode
		} else {
			c.logMode = noLogMode
		}
		return
	}
}

func SetLogger(log logger) MongoOption {
	return func(c *MongoDB) (err error) {
		c.logger = log
		return
	}
}
