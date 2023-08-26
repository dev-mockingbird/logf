// Copyright (c) 2023 Yang,Zhong
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package logf

import (
	"bytes"
	"runtime"
	"runtime/debug"
	"time"
)

type Record struct {
	CreatedAt time.Time
	PathFile  string
	Line      int
	Stack     [][]byte
	Level     Level
	Format    string
	Args      []any
}

func CollectRecord(l Level, callerDepth int, format string, args ...any) Record {
	record := Record{CreatedAt: time.Now(), Level: l, Format: format, Args: args}
	if _, file, line, ok := runtime.Caller(callerDepth); ok {
		record.PathFile = file
		record.Line = line
	}
	if l >= Warn {
		stack := debug.Stack()
		record.Stack = bytes.Split(stack, []byte{'\n'})[callerDepth+6:]
	}

	return record
}
