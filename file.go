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
	"errors"
	"fmt"
	"log"
	"os"
	"sync"
)

const FILESEP = "|||"

// an *os.File writer with locker.
type MuxWriter struct {
	sync.Mutex
	fd *os.File
}

type SFileWriter struct {
	l          *log.Logger
	mainLogger *SLogger
	fileName   string
	flags      int
	mw         *MuxWriter
}

// write to os.File.
func (mw *MuxWriter) Write(b []byte) (int, error) {
	mw.Lock()
	defer mw.Unlock()
	return mw.fd.Write(b)
}

// SetFd: set file descriptor
func (mw *MuxWriter) SetFd(fd *os.File) {
	if mw.fd != nil {
		mw.fd.Close()
	}
	mw.fd = fd
}

// NewFileWriter: create SFileWriter returning as ILogWriter.
func NewFileWriter() ILogWriter {
	fw := new(SFileWriter)
	fw.mw = new(MuxWriter)
	fw.flags = log.Ldate | log.Ltime
	fw.l = log.New(fw.mw, "", fw.flags)
	return fw
}

// getConfig: extract configuration
func (fw *SFileWriter) getConfig(config map[string]interface{}) error {
	if len(config) > 0 {
		confout := getConfig(config)
		if v, t := confout["filename"]; t {
			fw.fileName = v.(string)
		}
		if v, t := confout["flags"]; t {
			fw.flags = int(v.(float64))
			fw.l.SetFlags(fw.flags)
		}

	}
	return nil
}

// openFile, open the file for writing
func (fw *SFileWriter) openFile() error {
	prDeb("file.go - openFile", "Begin")
	// open the file
	var err error
	if fw.mw.fd != nil {
		fw.mw.fd.Close()
	}
	prDeb("file.go - openFile", "Open file "+fw.fileName)
	fw.mw.fd, err = os.OpenFile(fw.fileName, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0660)
	if err != nil {
		prDeb("file.go - openFile", "File error. Err: ", err)
		return err
	}
	prDeb("file.go - openFile", "File opened. Fd: ", fw.mw.fd)
	return nil
}

// Init file logger.
func (fw *SFileWriter) Init(logger *SLogger, config map[string]interface{}) error {
	fw.mainLogger = logger
	if err := fw.getConfig(config); err != nil {
		return err
	}
	if len(fw.fileName) == 0 {
		return errors.New("filename not configured")
	}
	return nil
}

// Write message on the file.
func (fw *SFileWriter) Write(msg SLogMsg) error {
	prDeb("file.go - Write", "MSG: ", msg)
	if !fw.mainLogger.MustWrite("file", msg) {
		return nil
	}
	if fw.mw.fd == nil {
		if err := fw.openFile(); err != nil {
			return err
		}
	}
	prDeb("file.go - Write", "Write message to file")
	fnc := fmt.Sprintf("%v[%s]", msg.fnc, msg.sev)
	if fw.flags > 0 {
		fw.l.Println(FILESEP, fnc, FILESEP, msg.msg)
	} else {
		fw.l.Println(fnc, FILESEP, msg.msg)
	}
	return nil
}

// implementing method. empty.
func (fw *SFileWriter) Destroy() {
	prDeb("Destroy", "close the file")
	fw.mw.fd.Close()
}

// implementing method. empty.
func (fw *SFileWriter) Flush() {
	fw.mw.fd.Sync()
}

func init() {
	CreateWriter("file", NewFileWriter)
}
