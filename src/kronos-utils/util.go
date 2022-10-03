package lib

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
)

func PrintJson(data []byte) {
	var prettyJSON bytes.Buffer
	err := json.Indent(&prettyJSON, data, "", "\t")
	if err != nil {
		log.Fatal(err)
	}

	log.Println(string(prettyJSON.Bytes()))
}

func SaveResource(fileName string, bin []byte) bool {
	if _, err := os.Stat(fileName); err != nil {
		file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0666)
		if err != nil {
			fmt.Println("Ошибка создания файла", fileName, ",", err.Error())
			return false
		}
		defer file.Close()

		_, err = file.Write(bin)
		if err != nil {
			fmt.Println("Ошибка записи в файл", fileName, ",", err.Error())
			return false
		}
	}

	return true
}
