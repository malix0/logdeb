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

type SConsoleWriter struct {
	l          *log.Logger
	mainLogger *SLogger
}

// NewConsole: create ConsoleWriter returning as ILogWriter.
func NewConsole() ILogWriter {
	cw := new(SConsoleWriter)
	cw.l = log.New(os.Stdout, "", log.Ldate|log.Ltime)
	return cw
}

// Init console logger.
func (cw *SConsoleWriter) Init(logger *SLogger, jsonconfig []byte) error {
	err := json.Unmarshal(jsonconfig, cw)
	if err != nil {
		return err
	}
	cw.mainLogger = logger
	return nil
}

// Write message in console.
func (cw *SConsoleWriter) Write(msg SLogMsg) error {
	prDeb("Write", msg)
	if !cw.mainLogger.MustWrite("console", msg) {
		return nil
	}
	cw.l.Println("|||", msg.fnc, "|||", msg.msg)
	return nil
}

// implementing method. empty.
func (cw *SConsoleWriter) Destroy() {

}

// implementing method. empty.
func (cw *SConsoleWriter) Flush() {

}

func init() {
	CreateWriter("console", NewConsole)
}
