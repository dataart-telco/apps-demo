package opencell_api

import (
	"io/ioutil"
	"path/filepath"
)

type IOUtils struct {
	
}

func (this IOUtils) GetFileData(path string) []byte {
	absFilePath, _ := filepath.Abs(path)
	rawData, err := ioutil.ReadFile(absFilePath)
	check(err)
	return rawData
}

func (this IOUtils) GetAbsolutePath(path string) string {
	absFilePath, _ := filepath.Abs(path)
	return absFilePath
}