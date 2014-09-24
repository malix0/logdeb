package logdeb

import (
	"encoding/json"
	"log"
	"os"
)

type ConsoleWriter struct {
	l          *log.Logger
	DebugLevel tDebLevel
}

// create ConsoleWriter returning as ILogWriter.
func NewConsole() ILogWriter {
	cw := new(ConsoleWriter)
	cw.l = log.New(os.Stdout, "", log.Ldate|log.Ltime)
	cw.DebugLevel = DLB
	//cw.Level = LevelTrace
	return cw
}

// init console logger.
// jsonconfig like '{"level":LevelTrace}'.
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