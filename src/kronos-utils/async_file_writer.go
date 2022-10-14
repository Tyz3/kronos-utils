package utils

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"time"
)

type AsyncFileWriter struct {
	fileName   string
	flag       int
	datePrefix bool

	lines       []string
	currentLine int64
	stopped     bool
	sync.Mutex
}

func NewAsyncFileWriter(fileName string, clear bool, datePrefix bool) *AsyncFileWriter {
	var flag int
	if clear {
		flag = os.O_CREATE | os.O_TRUNC
	} else {
		flag = os.O_CREATE | os.O_APPEND
	}

	return &AsyncFileWriter{
		fileName:   fileName,
		flag:       flag,
		datePrefix: datePrefix,
	}
}

func (afw *AsyncFileWriter) Start() {
	file, err := os.OpenFile(afw.fileName, afw.flag, 0666)
	if err != nil {
		panic(err)
		return
	}
	file.Close()

	go afw.run()
}

func (afw *AsyncFileWriter) Stop() {
	afw.stopped = true
	afw.saveToFile()
}

func (afw *AsyncFileWriter) Write(args ...string) {
	var sb strings.Builder
	if afw.datePrefix {
		sb.WriteString("[")
		sb.WriteString(time.Now().Format(time.ANSIC))
		sb.WriteString("] ")
	}
	sb.WriteString(strings.Join(args, " "))
	afw.Lock()
	afw.lines = append(afw.lines, sb.String())
	afw.Unlock()
}

func (afw *AsyncFileWriter) saveToFile() {
	file, err := os.OpenFile(afw.fileName, os.O_APPEND, 0666)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	if len(afw.lines) == 0 {
		return
	}

	for _, line := range afw.lines[afw.currentLine:] {
		var sb strings.Builder
		sb.WriteString(line)
		sb.WriteString("\n")
		_, err = file.WriteString(sb.String())
		if err != nil {
			fmt.Println(err)
			continue
		}
		afw.currentLine++
	}
}

func (afw *AsyncFileWriter) run() {
	for !afw.stopped {
		afw.saveToFile()
		time.Sleep(5 * time.Second)
	}
}
