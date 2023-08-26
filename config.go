// Copyright (c) 2023 Yang,Zhong
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package logf

import (
	"io"
	"os"
)

type Config struct {
	Printer     Printer
	LogLevel    Level
	CallerDepth int
	Prefix      string
	BufferSize  int
}

type Option func(c *Config)

func LogLevel(l Level) Option {
	return func(c *Config) {
		c.LogLevel = l
	}
}

func Prefix(prefix string) Option {
	return func(c *Config) {
		c.Prefix = prefix
	}
}

func Writer(w io.Writer) Option {
	return func(c *Config) {
		c.Printer = NewPrinter(w)
	}
}

func CustomPrinter(p Printer) Option {
	return func(c *Config) {
		c.Printer = p
	}
}

func Caller(depth int) Option {
	return func(c *Config) {
		c.CallerDepth = depth
	}
}

func BufferSize(bufsize int) Option {
	return func(c *Config) {
		c.BufferSize = bufsize
	}
}

func getConfig(opts ...Option) *Config {
	c := &Config{LogLevel: Info, CallerDepth: CallerDepth, BufferSize: 10}
	for _, o := range opts {
		o(c)
	}
	if c.Printer == nil {
		c.Printer = NewPrinter(os.Stdout)
	}
	return c
}

func (config Config) WithPrefix(p string) *Config {
	config.Prefix += p
	return &config
}
