package logdeb

import (
	"encoding/json"
	"log"
	"os"
)

type ConsoleWriter struct {
	l          *log.Logger
	Severity   tSeverity
	DebugLevel tDebLevel
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
