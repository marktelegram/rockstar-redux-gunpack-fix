package main

import (
	"log"
	"path/filepath"
	"sync"
	"time"
)

// runCopyingLoop запускает бесконечный цикл:
// если процесс GTA запущен - ждет, если нет - копирует файлы.
func runCopyingLoop(gtaProcessName, gtaProcessPath string, cfg *Config, username string) {
	log.Printf("Запускается цикл копирования: проверяем процесс %s (по имени или по пути)", gtaProcessName)
	for {
		if checkIfProcessRunning(gtaProcessName) || isProcessRunningByPath(gtaProcessPath) {
			log.Printf("Процесс %s обнаружен. Ждем 5 секунд перед повторной проверкой...", gtaProcessName)
			time.Sleep(5 * time.Second)
			continue
		}

		log.Println("GTA не запущен. Выполняем копирование...")

		var wg sync.WaitGroup
		log.Printf("Копирование gunpack: %s -> %s", cfg.GunpackNew, cfg.GunpackOld)
		copyDirRecursive(cfg.GunpackNew, cfg.GunpackOld, &wg)
		wg.Wait()

		var wg2 sync.WaitGroup
		log.Printf("Копирование redux: %s -> %s", cfg.ReduxNew, cfg.ReduxOld)
		copyDirRecursive(cfg.ReduxNew, cfg.ReduxOld, &wg2)
		wg2.Wait()

		defaultReduxBackupDir := filepath.Join("C:\\Users", username, "AppData", "Local", "altv-majestic", "backup")
		var wg3 sync.WaitGroup
		log.Printf("Копирование redux: %s -> %s", cfg.ReduxNew, defaultReduxBackupDir)
		copyDirRecursive(cfg.ReduxNew, defaultReduxBackupDir, &wg3)
		wg3.Wait()

		log.Println("Копирование завершено. Ждем 5 секунд перед следующей проверкой...")
		time.Sleep(5 * time.Second)
	}
}
