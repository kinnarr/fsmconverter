package config

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/felixangell/toml"
)

func ParseConfig() bool {
	searchDir := "fsm"
	fileList := []string{}
	_ = filepath.Walk(searchDir, func(path string, f os.FileInfo, err error) error {
		if strings.HasSuffix(path, ".toml") {
			fileList = append(fileList, path)
		}
		return nil
	})

	var tomlBuffer bytes.Buffer
	for _, file := range fileList {
		if readedBytes, err := ioutil.ReadFile(file); err == nil {
			tomlBuffer.Write(readedBytes)
		}
	}

	if _, err := toml.Decode(tomlBuffer.String(), &MainConfig); err != nil {
		fmt.Println(err)
		return false
	}

	return true
}
