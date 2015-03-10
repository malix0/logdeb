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
	"fmt"
	"log"
	"os"
)

const CONSSEP = "|||"

type SConsoleWriter struct {
	l          *log.Logger
	flags      int
	mainLogger *SLogger
}

// NewConsoleWriter: create SConsoleWriter returning as ILogWriter.
func NewConsoleWriter() ILogWriter {
	cw := new(SConsoleWriter)
	cw.flags = log.Ldate | log.Ltime
	cw.l = log.New(os.Stdout, "", cw.flags)
	return cw
}

// getConfig: extract configuration
func (cw *SConsoleWriter) getConfig(config map[string]interface{}) error {
	if len(config) > 0 {
		confout := getConfig(config)
		if v, t := confout["flags"]; t {
			cw.flags = int(v.(float64))
			cw.l.SetFlags(cw.flags)
		}
	}
	return nil
}

// Init console logger.
func (cw *SConsoleWriter) Init(logger *SLogger, config map[string]interface{}) error {
	cw.mainLogger = logger
	if err := cw.getConfig(config); err != nil {
		return err
	}
	return nil
}

// Write message in console.
func (cw *SConsoleWriter) Write(msg SLogMsg) error {
	prDeb("Write", msg)
	if !cw.mainLogger.MustWrite("console", msg) {
		return nil
	}
	fnc := fmt.Sprintf("%v[%s]", msg.fnc, msg.sev)
	if cw.flags > 0 {
		cw.l.Println(CONSSEP, fnc, CONSSEP, msg.msg)
	} else {
		cw.l.Println(fnc, CONSSEP, msg.msg)
	}
	return nil
}

// implementing method. empty.
func (cw *SConsoleWriter) Destroy() {

}

// implementing method. empty.
func (cw *SConsoleWriter) Flush() {

}

func init() {
	CreateWriter("console", NewConsoleWriter)
}
