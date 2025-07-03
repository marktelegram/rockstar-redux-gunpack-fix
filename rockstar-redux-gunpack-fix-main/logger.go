package main

import (
	"io"
	"log"
	"os"
)

// initLogger открывает (или создаёт) лог-файл и настраивает логгер.
// Если logFilePath пустой — используется "app.log".
// Логи пишутся и в файл, и в консоль.
func initLogger(logFilePath string) (*os.File, error) {
	if logFilePath == "" {
		logFilePath = "app.log"
	}

	logFile, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}

	multiWriter := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(multiWriter)
	log.SetFlags(log.LstdFlags)
	log.Println("Запуск программы auto-redux-gunpack")

	return logFile, nil
}
