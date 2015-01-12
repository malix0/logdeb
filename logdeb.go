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
	"strings"
	"sync"
	"time"
)

const cPckName = "logdeb"

// Log severity - type tSeverity
const (
	SEVFATAL = iota + 1 // fatal
	SEVERROR            // error
	SEVWARN             // warning
	SEVINFO             // information
	SEVDEBUG            // debug
)

// Debug level - type tDebLevel
const (
	DLB   = iota + 1 // base
	DLE              // extended
	DLV              // verbose
	DLVV             // very verbose
	DLVVV            // even more verbose
)

type tSeverity int8

type tDebLevel int8

type tFncName string

// Basic write rule
type sBaseRule struct {
	Severity   tSeverity
	DebugLevel tDebLevel
}

// Write rules generic and with func granularity
type sWriteRules struct {
	sBaseRule
	FncRules map[tFncName]sBaseRule
}

type SLogMsg struct {
	fnc    tFncName
	msg    string
	sev    tSeverity
	debLev tDebLevel
}

type ILogWriter interface {
	Init(logger *SLogger, config map[string]interface{}) error
	Write(msg SLogMsg) error
	Destroy()
	Flush()
}

type SLogWriter struct {
	writer     ILogWriter
	writeRules sWriteRules
}

type tLogWriter func() ILogWriter

var logWriters = make(map[string]tLogWriter)

// SLogger is the basic struct of deblog
type SLogger struct {
	lock        sync.Mutex            // ensures atomic writes; protects the following fields
	wg          sync.WaitGroup        // wait until all channels are drained
	msgChan     chan *SLogMsg         // Channels that will dispatch the log messages
	writers     map[string]SLogWriter // Log writers
	buf         bytes.Buffer          // for accumulating text to write
	sBaseRule                         // Log Severity & DebugLevel
	maxSeverity tSeverity             // Maximum severity defined for writers
	maxDebLev   tDebLevel             // Maximum debug level used by writers. Used to discard message immediately
	UseFncRules bool                  // Define if write function rules must be used
	sessionId   string                // Log session Id
}

const DEBUG = "1"

func prDeb(fnc string, par ...interface{}) {
	if DEBUG == "1" {
		fmt.Println("*D*", "[["+fnc+"]]", par)
	}
}

func GetTsStr() string {
	var t time.Time
	t = time.Now()
	year, month, day := t.Date()
	hour, min, sec := t.Clock()
	msec := t.Nanosecond() / 1e6
	return fmt.Sprintf("%d%02d%02d%02d%02d%02d%03d", year, month, day, hour, min, sec, msec)
}

func getWriteRules(config map[string]interface{}) sWriteRules {
	prDeb("getWriteRules", "config:", config)
	wr := new(sWriteRules)
	wr.extract(config)
	for k, v := range config {
		prDeb("getWriteRules", "k:", k, "v:", v, "strings.Title(k):", strings.Title(k))
		if strings.ToLower(k) == "fncrules" {
			wr.FncRules = make(map[tFncName]sBaseRule)
			rules := v.(map[string]interface{})
			for fnc, rc := range rules {
				fnc := tFncName(fnc)
				r := new(sBaseRule)
				r.extract(rc.(map[string]interface{}))
				wr.FncRules[fnc] = *r
			}
		}
	}
	prDeb("getWriteRules", "WriteRules:", *wr)
	return *wr
}

// NewLogDeb start configured writers and returns a new SLogger
// - bufferSize: is the size of channel that hold messages before send to writers
// - config:     is the configuration in JSON format. Like {"console":{"Severity":5, "DebugLevel":3}}
func NewLogDeb(bufferSize int64, config string) *SLogger {
	const cFncName = cPckName + ".NewLogDeb"
	l := new(SLogger)
	l.SetSeverity(SEVERROR)
	l.SetDebugLevel(DLB)
	l.SetSessionId("GEN" + GetTsStr())
	l.msgChan = make(chan *SLogMsg, bufferSize)
	l.writers = make(map[string]SLogWriter)
	l.wg.Add(1)
	go l.StartWriter()

	// Read writers and their configuration from config
	var writersConf map[string]json.RawMessage
	err := json.Unmarshal([]byte(config), &writersConf)
	if err != nil {
		panic(fmt.Sprintf("Error extracting writer config from json. ERR: %s", err))
	}

	l.lock.Lock()
	defer l.lock.Unlock()
	for wr, c := range writersConf {
		if wr == "main" {
			json.Unmarshal(c, &l)
			if l.Severity == 0 {
				l.SetSeverity(SEVERROR)
			} else {
				l.setMaxSeverity(l.Severity)
			}
			if l.DebugLevel == 0 {
				l.SetDebugLevel(DLB)
			} else {
				l.setMaxDebugLevel(l.DebugLevel)
			}
			prDeb(cFncName, *l)
		} else {
			if logWriter, ok := logWriters[wr]; ok {
				lw := logWriter()
				var ci interface{}
				json.Unmarshal(c, &ci)
				cm := ci.(map[string]interface{})
				lw.Init(l, cm)
				l.writers[wr] = SLogWriter{writer: lw, writeRules: getWriteRules(cm)}
				l.setMaxSeverity(l.writers[wr].writeRules.Severity)
				l.setMaxDebugLevel(l.writers[wr].writeRules.DebugLevel)
			} else {
				panic(fmt.Sprintf("logdeb: unknown writer %q (forgotten Register?)", wr))
			}
		}
	}
	if len(l.writers) == 0 {
		panic("No writer configured")
	}
	return l
}

// CreateWriter register a writer adapter
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

// SetSessionId set the log session unique identification
func (l *SLogger) SetSessionId(sessionId string) {
	l.sessionId = sessionId
}

func (l *SLogger) setMaxSeverity(sev tSeverity) {
	prDeb("setMaxSeverity", "sev:", sev, "maxSev:", l.maxSeverity)
	if sev > l.maxSeverity {
		l.maxSeverity = sev
	}
}

func (l *SLogger) SetSeverity(sev tSeverity) {
	l.Severity = sev
	l.setMaxSeverity(sev)
}

func (l *SLogger) setMaxDebugLevel(debLev tDebLevel) {
	prDeb("setMaxDebugLevel", "debLev:", debLev, "maxDebLev:", l.maxDebLev)
	if debLev > l.maxDebLev {
		l.maxDebLev = debLev
	}
	prDeb("setMaxDebugLevel", "debLev:", debLev, "maxDebLev:", l.maxDebLev)
}

func (l *SLogger) SetDebugLevel(debLev tDebLevel) {
	l.DebugLevel = debLev
	l.setMaxDebugLevel(debLev)
	prDeb("SetDebugLevel", "debLev:", debLev, "maxDebLev:", l.maxDebLev)
}

// extract: extract base rule value from config
func (r *sBaseRule) extract(config map[string]interface{}) {
	for k, v := range config {
		prDeb("extract", "k:", k, "v:", v)
		if strings.ToLower(k) == "severity" {
			prDeb("extract", "Set Severity:", v)
			r.Severity = tSeverity(v.(float64))
		} else if strings.ToLower(k) == "debuglevel" {
			prDeb("extract", "Set Debuglevel:", v)
			r.DebugLevel = tDebLevel(v.(float64))
		}
	}
}

// get: get rule or parent rule when rule value is null
func (r sBaseRule) get(parentRule sBaseRule) sBaseRule {
	if r.Severity == 0 {
		r.Severity = parentRule.Severity
	}
	if r.DebugLevel == 0 {
		r.DebugLevel = parentRule.DebugLevel
	}
	return r
}

// eval: evaluate the write rule and return true if the message
// is to be written
func (r sBaseRule) eval(msg SLogMsg, parentRule sBaseRule) bool {
	r = r.get(parentRule)
	return ((msg.sev < SEVDEBUG && msg.sev <= r.Severity) || (r.Severity == SEVDEBUG && msg.debLev <= r.DebugLevel))
}

func (l *SLogger) MustWrite(writerName string, msg SLogMsg) bool {
	prDeb("MustWrite")
	if l.UseFncRules && len(l.writers[writerName].writeRules.FncRules) > 0 {
		// Search inside FncRules for function name matching
		// or plartial matching
		for {
			if len(msg.fnc) == 0 {
				return false
			}
			if rule, ok := l.writers[writerName].writeRules.FncRules[msg.fnc]; ok {
				return rule.eval(msg, l.writers[writerName].writeRules.get(sBaseRule{l.Severity, l.DebugLevel}))
			}
			msg.fnc = msg.fnc[:len(msg.fnc)-1]
		}
	} else {
		prDeb("MustWrite", "BaseRule", l.writers[writerName].writeRules)
		return l.writers[writerName].writeRules.eval(msg, sBaseRule{l.Severity, l.DebugLevel}) // l.evalRule(msg, sBaseRule())
	}
	return false
}

func (l *SLogger) StartWriter() {
	prDeb("StartWriter", "run")
	defer l.wg.Done()
	for lm := range l.msgChan {
		for _, lw := range l.writers {
			prDeb("StartWriter", "lw:", lw, ":: lm:", *lm)
			lw.writer.Write(*lm)
		}
	}
}

func (l *SLogger) logw(fnc tFncName, msg string, sev tSeverity, debLev tDebLevel) error {
	const cFncName = cPckName + ".logw"
	prDeb(cFncName, "sev:", sev, ":: maxSeverity:", l.maxSeverity, ":: UseFncRules:", l.UseFncRules)
	if !l.UseFncRules && sev > l.maxSeverity {
		prDeb(cFncName, "EXIT")
		return nil
	}
	prDeb(cFncName, "WRITE:", msg)
	lm := &SLogMsg{fnc: fnc, msg: msg, sev: sev, debLev: debLev}
	l.msgChan <- lm
	return nil
}

func (l *SLogger) Fatal(fnc tFncName, msg string) {
	l.logw(fnc, msg, SEVFATAL, 0)
}

func (l *SLogger) Err(fnc tFncName, msg string) {
	l.logw(fnc, msg, SEVERROR, 0)
}

func (l *SLogger) Warn(fnc tFncName, msg string) {
	l.logw(fnc, msg, SEVWARN, 0)
}

func (l *SLogger) Info(fnc tFncName, msg string) {
	l.logw(fnc, msg, SEVINFO, 0)
}

// Deb: log message with severity debug
func (l *SLogger) Deb(fnc tFncName, msg string) {
	l.logw(fnc, msg, SEVDEBUG, DLB)
}

// Debl: log message with severity debug and input debug level
func (l *SLogger) Debl(fnc tFncName, msg string, debLev tDebLevel) {
	if debLev <= l.maxDebLev {
		l.logw(fnc, msg, SEVDEBUG, debLev)
	}
}

// Flush all chan data
func (l *SLogger) Flush() {
	for _, lw := range l.writers {
		lw.writer.Flush()
	}
}

// Destroy logger, flush all chan data and destroy all writers.
func (l *SLogger) Destroy() {
	close(l.msgChan)
	l.wg.Wait()
	for _, lw := range l.writers {
		lw.writer.Flush()
		lw.writer.Destroy()
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
