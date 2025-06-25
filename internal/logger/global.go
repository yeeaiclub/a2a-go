// Copyright 2025 yumosx
//
// Licensed under the Apache License, Version 2.0 (the \"License\");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an \"AS IS\" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package log

import (
	"os"
	"sync"
)

var (
	globalLogger Logger
	globalMu     sync.Mutex
)

func init() {
	globalLogger = New(os.Stderr, DebugLevel)
}

func Global() Logger {
	return globalLogger
}

func SetGlobalLogger(logger Logger) {
	globalMu.Lock()
	defer globalMu.Unlock()
	globalLogger = logger
}

func Debug(args ...interface{})                 { globalLogger.Debug(args...) }
func Info(args ...interface{})                  { globalLogger.Info(args...) }
func Warn(args ...interface{})                  { globalLogger.Warn(args...) }
func Error(args ...interface{})                 { globalLogger.Error(args...) }
func Fatal(args ...interface{})                 { globalLogger.Fatal(args...) }
func Debugf(format string, args ...interface{}) { globalLogger.Debugf(format, args...) }
func Infof(format string, args ...interface{})  { globalLogger.Infof(format, args...) }
func Warnf(format string, args ...interface{})  { globalLogger.Warnf(format, args...) }
func Errorf(format string, args ...interface{}) { globalLogger.Errorf(format, args...) }
func Fatalf(format string, args ...interface{}) { globalLogger.Fatalf(format, args...) }

func SetLevel(level Level) {
	globalLogger.SetLevel(level)
}
