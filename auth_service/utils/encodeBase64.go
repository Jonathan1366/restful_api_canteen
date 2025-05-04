package utils

import (
	"encoding/base64"
	"io"
	"os"
)

func EncodeFileToBase64(filePath string)(string, error) {
	file, err:= os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	var fileBytes []byte
	fileBytes, err = io.ReadAll(file)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(fileBytes), nil

}