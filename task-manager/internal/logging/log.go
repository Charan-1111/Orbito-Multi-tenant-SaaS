package logging

import (
	"os"
	"sync"

	"github.com/rs/zerolog"
)

type Log struct {
	log   zerolog.Logger
	close func()
	once  sync.Once
}

func (log *Log) Initialize() {
	log.once = sync.Once{}
	log.log = zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout}).With().Timestamp().Logger()
	log.close = func() {}
}

func (log *Log) Close() {
	log.once.Do(func() {
		if log.close != nil {
			log.close()
		}
	})
}
