// Copyright (c) 2023 Yang,Zhong
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package logf

import (
	"fmt"
	"io"
	"os"
	"path"

	"github.com/fatih/color"
)

type Printer interface {
	Print(prefix string, record Record)
}

type printer struct {
	w           io.Writer
	enableColor bool
}

func NewPrinter(w io.Writer) Printer {
	return printer{w: w, enableColor: w == os.Stdout}
}

func (a printer) Print(prefix string, record Record) {
	a.print(prefix, record)
	for _, msg := range record.Stack {
		record.Format = string(msg)
		record.Args = nil
		a.print(prefix, record)
	}
}

func (a printer) print(prefix string, record Record) {
	fmt.Fprintf(a.w, a.colorMsg(record.Level, fmt.Sprintf(
		"%s %s:%d:\t[%s]\t%s%s\n",
		record.CreatedAt.Format("2006/01/02 15:04:05"),
		path.Base(record.PathFile),
		record.Line,
		LevelString(record.Level),
		prefix,
		fmt.Sprintf(record.Format, record.Args...))))
}

func (a printer) colorMsg(level Level, msg string) string {
	if a.enableColor {
		switch level {
		case Trace:
			return color.WhiteString(msg)
		case Debug:
			return color.CyanString(msg)
		case Info:
			return color.GreenString(msg)
		case Warn:
			return color.YellowString(msg)
		case Error:
			return color.RedString(msg)
		case Fatal:
			return color.BlueString(msg)
		}
		return msg
	}
	return msg
}
