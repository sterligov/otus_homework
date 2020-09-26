package main

import (
	"fmt"
	"io"
	"strings"
)

const cursorPreviousLine = "\033[F"

type Loadbar interface {
	Update(int)
}

type BarWriter struct {
	writer io.Writer
	bar    Loadbar
}

func NewBarWriter(w io.Writer, bar Loadbar) *BarWriter {
	return &BarWriter{
		writer: w,
		bar:    bar,
	}
}

func (w *BarWriter) Write(p []byte) (n int, err error) {
	n, err = w.writer.Write(p)
	w.bar.Update(n)

	return n, err
}

type LoadbarLine struct {
	out       io.Writer
	dataSize  int64
	progress  int64
	lineSize  int
	lineChunk int
	percent   int
}

func NewLoadbarLine(dataSize int64, lineSize int, writer io.Writer) *LoadbarLine {
	return &LoadbarLine{
		out:      writer,
		dataSize: dataSize,
		lineSize: lineSize,
	}
}

func (l *LoadbarLine) Update(progress int) {
	if l.progress == 0 {
		l.start()
	}
	if l.progress >= l.dataSize {
		return
	}

	l.progress += int64(progress)

	l.updatePercent()
	l.updateLine()
}

func (l *LoadbarLine) start() {
	fmt.Fprintf(l.out, "\n")
}

func (l *LoadbarLine) updatePercent() {
	percent := int(l.progress * 100 / l.dataSize)
	if l.percent == percent {
		return
	}

	l.percent = percent
	l.printBar()
}

func (l *LoadbarLine) updateLine() {
	lineChunk := int(l.progress * int64(l.lineSize) / l.dataSize)
	if l.lineChunk == lineChunk {
		return
	}

	l.lineChunk = lineChunk
	l.printBar()
}

func (l *LoadbarLine) printBar() {
	line := strings.Repeat("=", l.lineChunk) + strings.Repeat(" ", l.lineSize-l.lineChunk)

	fmt.Fprint(l.out, cursorPreviousLine)
	fmt.Fprintf(l.out, "%3d%%  [%s]\n", l.percent, line)
}
