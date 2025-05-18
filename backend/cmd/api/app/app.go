package app

import "sync"

type App struct {
	mu sync.RWMutex
}

var (
	instance *App
	once     sync.Once
)

func Get() *App {
	once.Do(func() {
		instance = &App{}
	})

	return instance
}
