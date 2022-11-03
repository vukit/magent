package logger

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/vukit/magent/internal/config"
)

var message struct {
	Host    string `json:"host"`
	Level   string `json:"level"`
	Message string `json:"message"`
	Time    string `json:"time"`
}

func TestInfo(t *testing.T) {
	var err error

	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}

	tConfig := config.Config{}
	filename := "../../configs/global.json"
	tConfig.Read(&filename)

	tLogger := NewLogger(&tConfig, w)

	tLogger.Info("info message")

	buffer := make([]byte, 1024)
	n, err := r.Read(buffer)
	if err != nil {
		t.Fatal(err)
	}

	err = json.Unmarshal(buffer[:n], &message)
	if err != nil {
		t.Fatal(err)
	}

	if message.Host != "localhost" {
		t.Errorf("wrong host, want 'localhost', got '%s'", message.Host)
	}

	if message.Level != "info" {
		t.Errorf("wrong level, want 'info', got '%s'", message.Level)
	}

	if message.Message != "info message" {
		t.Errorf("wrong message, want 'info message', got '%s'", message.Message)
	}
}

func TestWarning(t *testing.T) {
	var err error

	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}

	tConfig := config.Config{}
	filename := "../../configs/global.json"
	tConfig.Read(&filename)

	tLogger := NewLogger(&tConfig, w)

	tLogger.logger.Output(w)

	tLogger.Warning("warning message")

	buffer := make([]byte, 1024)
	n, err := r.Read(buffer)
	if err != nil {
		t.Fatal(err)
	}

	err = json.Unmarshal(buffer[:n], &message)
	if err != nil {
		t.Fatal(err)
	}

	if message.Host != "localhost" {
		t.Errorf("wrong host, want 'localhost', got '%s'", message.Host)
	}

	if message.Level != "warn" {
		t.Errorf("wrong level, want 'warning', got '%s'", message.Level)
	}

	if message.Message != "warning message" {
		t.Errorf("wrong msg, want 'warning message', got '%s'", message.Message)
	}
}

func TestDebug(t *testing.T) {
	var err error

	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}

	tConfig := config.Config{}
	filename := "../../configs/global.json"
	tConfig.Read(&filename)
	tConfig.Common.Debug = true

	tLogger := NewLogger(&tConfig, w)

	tLogger.logger.Output(w)

	tLogger.Debug("debug message")

	buffer := make([]byte, 1024)
	n, err := r.Read(buffer)
	if err != nil {
		t.Fatal(err)
	}

	err = json.Unmarshal(buffer[:n], &message)
	if err != nil {
		t.Fatal(err)
	}

	if message.Host != "localhost" {
		t.Errorf("wrong host, want 'localhost', got '%s'", message.Host)
	}

	if message.Level != "debug" {
		t.Errorf("wrong level, want 'debug', got '%s'", message.Level)
	}

	if message.Message != "debug message" {
		t.Errorf("wrong msg, want 'debug message', got '%s'", message.Message)
	}
}
