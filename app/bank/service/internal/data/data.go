package data

import (
	"github.com/go-kratos/kratos/v2/log"
)

// Data .
type Data struct {
	log *log.Helper
}

// NewData .
func NewData(logger log.Logger) (*Data, func(), error) {
	l := log.NewHelper(log.With(logger, "module", "data/bank-service"))

	d := &Data{
		log: l,
	}

	return d, func() {
		l.Info("message", "closing the data resources")
	}, nil
}
