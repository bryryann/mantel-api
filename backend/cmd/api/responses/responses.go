package responses

import (
	"log/slog"
	"os"
	"sync"
)

type Responses struct {
	logger *slog.Logger
}

var (
	instance *Responses
	once     sync.Once
)

func Get() *Responses {
	once.Do(func() {
		instance = &Responses{
			logger: slog.New(slog.NewTextHandler(os.Stdout, nil)),
		}
	})

	return instance
}
