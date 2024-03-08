package commons

import (
	"log"
	"os"
)

func NewLogger(prefix string) *log.Logger {
	return log.New(os.Stdout, prefix, log.Lmsgprefix|log.LstdFlags)
}
