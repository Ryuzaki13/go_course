package utils

import (
	"errors"
	"os"
	"strings"
)

// LoadFile загрузить файл и получить указатель на его слайс байтов
func LoadFile(filename string) ([]byte, error) {
	var file *os.File
	var e error

	file, e = os.Open(filename)
	if e != nil {
		return nil, e
	}
	defer func(file *os.File) {
		_ = file.Close()
	}(file)

	stat, e := file.Stat()
	if e != nil {
		return nil, errors.New("ошибка получения информации о файле " + filename + ": " + e.Error())
	}

	bs := make([]byte, stat.Size())
	_, e = file.Read(bs)
	if e != nil {
		return nil, errors.New("ошибка чтения файла " + filename + ": " + e.Error())
	}

	return bs, nil
}

// LoadAssets загрузка файла текстового формата. Для загрузки CSS и JS
func LoadAssets(filename string) string {
	bytes, e := LoadFile(filename)
	if e != nil {
		return ""
	} else {
		return strings.TrimSpace(string(bytes))
	}
}
