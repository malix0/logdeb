// Copyright 2014 Massimo Fidanza.
//
// Licensed under the Apache License, Version 2.0 (the "License"): you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.

package logdeb

import (
	"bytes"
	"encoding/json"
	"fmt"
	"sync"
)

const cPckName = "logdeb"

// Log severity - type tSeverity
const (
	SEVDEBUG = iota + 1 // debug
	SEVINFO             // information
	SEVWARN             // warning
	SEVERROR            // error
	SEVFATAL            // fatal
)

// Debug level - type tDebLevel
const (
	DLB   = iota + 1 // base
	DLE              // extended
	DLV              // verbose
	DLVV             // very verbose
	DLVVV            // even more verbose
)

type SLogger struct {
	lock      sync.Mutex            // ensures atomic writes; protects the following fields
	funcName  string                // prefix to write at beginning of each line
	severity  tSeverity             // Log severity
	msgChan   chan *TLogMsg         // Channels that will dispatch the log messages
	writers   map[string]ILogWriter // Log writers
	buf       bytes.Buffer          // for accumulating text to write
	maxDebLev tDebLevel
}

type tSeverity int8

type tDebLevel int8

type TLogMsg struct {
	fnc    string
	msg    string
	sev    tSeverity
	debLev tDebLevel
}

type ILogWriter interface {
	Init(logger *SLogger, jsonconfig []byte) error
	Write(msg TLogMsg) error
	Destroy()
	Flush()
}

type tLogWriter func() ILogWriter

var logWriters = make(map[string]tLogWriter)

func NewLogDeb(bufferSize int64, config string) *SLogger {
	l := new(SLogger)
	//l.severity = SEVERROR
	l.msgChan = make(chan *TLogMsg, bufferSize)
	l.writers = make(map[string]ILogWriter)
	go l.StartWriter()

	// Read writers and their configuration from config
	var writers map[string]json.RawMessage
	err := json.Unmarshal([]byte(config), &writers)
	if err != nil {
		panic(fmt.Sprintf("Error extracting writer config from json. ERR: %s", err))
	}

	l.lock.Lock()
	defer l.lock.Unlock()
	for wr, c := range writers {
		//if log, ok := adapters[adaptername]; ok {
		if logWriter, ok := logWriters[wr]; ok {
			lw := logWriter()
			lw.Init(l, c)
			l.writers[wr] = lw
		} else {
			panic(fmt.Sprintf("logdeb: unknown writer %q (forgotten Register?)", wr))
		}
	}
	return l
}

func CreateWriter(name string, writer tLogWriter) {
	const cFncName = cPckName + ".CreateWriter"
	if writer == nil {
		panic(fmt.Sprintf("%s: writer %s can not be nil", cFncName, name))
	}
	if _, dup := logWriters[name]; dup {
		panic(fmt.Sprintf("%s: writer %s already exists", cFncName, name))
	}
	logWriters[name] = writer
}

func (l *SLogger) SetMaxDebugLevel(debugLevel tDebLevel) {
	if debugLevel > l.maxDebLev {
		l.maxDebLev = debugLevel
	}
}

func (l *SLogger) SetSeverity(severity tSeverity) {
	l.severity = severity
}

func (l *SLogger) StartWriter() {
	for {
		select {
		case lm := <-l.msgChan:
			for _, w := range l.writers {
				w.Write(*lm)
			}
		}
	}
}

func (l *SLogger) logw(fn string, msg string, sev tSeverity, debLev tDebLevel) error {
	if sev < l.severity {
		return nil
	}
	mt := &TLogMsg{fnc: fn, msg: msg, sev: sev, debLev: debLev}
	//mt := fmt.Sprintf("[%s - %s] %s", fn, sev, msg)
	// TODO: Write messages to buffer. Filter by severity
	l.msgChan <- mt
	return nil
}

func (l *SLogger) Deb(fn string, msg string) {
	l.logw(fn, msg, SEVDEBUG, DLB)
}

// Debug with debug level
func (l *SLogger) Debl(fn string, msg string, debLev tDebLevel) {
	if debLev <= l.maxDebLev {
		l.logw(fn, msg, SEVDEBUG, debLev)
	}
}

func (sev tSeverity) String() string {
	switch sev {
	case SEVDEBUG:
		return "D"
	case SEVINFO:
		return "I"
	case SEVWARN:
		return "W"
	case SEVERROR:
		return "E"
	case SEVFATAL:
		return "F"
	default:
		return "Unknown severity: " + string(sev)
	}
}
