// Copyright 2025 yeeaiclub
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package log

import (
	"io"
	"log"
)

type stdImplementation struct {
	logger *log.Logger
}

func newStdImplementation(out io.Writer) *stdImplementation {
	return &stdImplementation{
		logger: log.New(out, "", log.Ldate|log.Ltime|log.Lmicroseconds|log.LUTC),
	}
}

func (s *stdImplementation) format(prefix string, msg string) string {
	return prefix + msg
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
