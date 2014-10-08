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
	"encoding/json"
	"log"
	"os"
)

type ConsoleWriter struct {
	l *log.Logger
	tWriteRule
}

// create ConsoleWriter returning as ILogWriter.
func NewConsole() ILogWriter {
	cw := new(ConsoleWriter)
	cw.l = log.New(os.Stdout, "", log.Ldate|log.Ltime)
	cw.Severity = SEVERROR
	cw.DebugLevel = DLB
	//cw.Level = LevelTrace
	return cw
}

// init console logger.
// jsonconfig like '{Severity":SEVDEBUG, "DebugLevel":DLB}'.
func (cw *ConsoleWriter) Init(logger *SLogger, jsonconfig []byte) error {
	err := json.Unmarshal(jsonconfig, cw)
	if err != nil {
		return err
	}
	logger.SetMaxDebugLevel(cw.DebugLevel)
	return nil
}

// write message in console.
func (cw *ConsoleWriter) Write(msg TLogMsg) error {
	if msg.sev < cw.Severity || (msg.sev == SEVDEBUG && msg.debLev > cw.DebugLevel) {
		return nil
	}
	cw.l.Println("|||", msg.fnc, "|||", msg.msg)
	return nil
}

// implementing method. empty.
func (cw *ConsoleWriter) Destroy() {

}

// implementing method. empty.
func (cw *ConsoleWriter) Flush() {

}

func init() {
	CreateWriter("console", NewConsole)
}
