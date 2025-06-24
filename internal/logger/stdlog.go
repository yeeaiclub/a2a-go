package log

import (
	"fmt"
	"io"
	"log"
	"os"
)

type stdImplementation struct {
	logger *log.Logger
	prefix string
}

func newStdImplementation(out io.Writer) *stdImplementation {
	pid := os.Getpid()
	return &stdImplementation{
		logger: log.New(out, "", log.Ldate|log.Ltime|log.Lmicroseconds|log.LUTC),
		prefix: fmt.Sprintf("asynq: pid=%d ", pid),
	}
}

func (s *stdImplementation) format(prefix string, msg string) string {
	return s.prefix + prefix + msg
}

func (s *stdImplementation) Debug(msg string) {
	s.logger.Print(s.format("DEBUG: ", msg))
}

func (s *stdImplementation) Info(msg string) {
	s.logger.Print(s.format("INFO: ", msg))
}

func (s *stdImplementation) Warn(msg string) {
	s.logger.Print(s.format("WARN: ", msg))
}

func (s *stdImplementation) Error(msg string) {
	s.logger.Print(s.format("ERROR: ", msg))
}

func (s *stdImplementation) Fatal(msg string) {
	s.logger.Print(s.format("FATAL: ", msg))
}
