// SPDX-License-Identifier: MIT

package gobuild

import (
	"io"
	"os"

	"github.com/issue9/term/v3/colors"
)

// Log 日志类型
type Log struct {
	Type    int8
	Message string
}

// 日志类型
const (
	LogTypeSuccess int8 = iota
	LogTypeInfo
	LogTypeWarn
	LogTypeError
	LogTypeIgnore
)

// ConsoleLogs 将日志输出到控制台
type ConsoleLogs struct {
	Logs       chan *Log
	showIgnore bool
	writers    map[int8]*logWriter
	stop       chan struct{}
}

// NewConsoleLogs 声明 ConsoleLogs 实例
func NewConsoleLogs(showIgnore bool) *ConsoleLogs {
	return newConsoleLogs(showIgnore, os.Stderr, os.Stdout)
}

func newConsoleLogs(showIgnore bool, err, out io.Writer) *ConsoleLogs {
	logs := &ConsoleLogs{
		Logs:       make(chan *Log, 100),
		showIgnore: showIgnore,
		writers: map[int8]*logWriter{
			LogTypeSuccess: newWriter(out, colors.Green, "[SUCC] "),
			LogTypeInfo:    newWriter(out, colors.Blue, "[INFO] "),
			LogTypeWarn:    newWriter(err, colors.Magenta, "[WARN] "),
			LogTypeError:   newWriter(err, colors.Red, "[ERRO] "),
			LogTypeIgnore:  newWriter(out, colors.Default, "[IGNO] "),
		},
	}

	go logs.output()

	return logs
}

// Stop 停止输出
func (logs *ConsoleLogs) Stop() {
	logs.stop <- struct{}{}
}

func (logs *ConsoleLogs) output() {
	for {
		select {
		case log := <-logs.Logs:
			if !logs.showIgnore && log.Type == LogTypeIgnore {
				continue
			}

			w := logs.writers[log.Type]
			colors.Fprint(w.out, colors.Normal, w.color, colors.Default, w.prefix)
			colors.Fprintln(w.out, colors.Normal, colors.Default, colors.Default, log.Message)
		case <-logs.stop:
			return
		}
	}
}

// 带色彩输出的控制台。
type logWriter struct {
	out    io.Writer
	color  colors.Color
	prefix string
}

func newWriter(out io.Writer, color colors.Color, prefix string) *logWriter {
	return &logWriter{
		out:    out,
		color:  color,
		prefix: prefix,
	}
}
