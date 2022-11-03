package config

import (
	"fmt"
	"os"
	"testing"
)

func TestRead(t *testing.T) {
	var (
		err      error
		filename string
	)

	tConfig := Config{}

	filename = "../../configs/nonexistent.json"
	err = tConfig.Read(&filename)
	if err == nil {
		t.Error(err)
	}

	incorrectJson := `{"common": {}, "sensors": [{"name": "cpu", "enable": true, "metrics": ["usage"]}] "collectors": []}`
	os.WriteFile("../../configs/incorrect.json", []byte(incorrectJson), os.ModePerm)
	filename = "../../configs/incorrect.json"
	err = tConfig.Read(&filename)
	fmt.Println(err)
	if err == nil {
		t.Error(err)
	}
	os.Remove("../../configs/incorrect.json")

	filename = "../../configs/global.json"
	err = tConfig.Read(&filename)
	if err != nil {
		t.Error(err)
	}

	err = tConfig.Read(&filename)
	if err != nil {
		t.Error(err)
	}

}
