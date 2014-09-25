package logdeb

import (
	"testing"
)

func TestConsole(t *testing.T) {
	config := `{"console":{"Severity":1, "DebugLevel":3}}`
	l := NewLogDeb(10, config)
	l.Deb("TestConsole", "test console")
	l.Debl("TestConsole", "test console with debug level 4", 4) // Don't write this because DebugLevel is set to 3
	l.Debl("TestConsole", "test console with debug level 3", 3)

	config = `{"console":{"Severity":1}}`
	l2 := NewLogDeb(10, config)
	l2.Debl("TestConsole2", "test console with debug level 1", 1)
	l2.Debl("TestConsole2", "test console with debug level 2", 2)
}
