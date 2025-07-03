package main

import (
	"fmt"
	"log"
	"os"
	"time"
)

// defaultGTAPathForUser возвращает зашитый путь к GTA5.exe
func defaultGTAPathForUser() string {
	return "C:\\Games\\GTA5RP\\RageMP\\GTA5.exe"
}

func main() {
	initLogger("app.log")
	configPath := "config.json"
	var cfg *Config
	var err error

	// Если конфигурация отсутствует – первичный запуск: запрос настроек
	if _, err = os.Stat(configPath); os.IsNotExist(err) {
		log.Println("Конфигурация не найдена. Выполняется первичная настройка.")
		cfg = &Config{}

		fmt.Println("Первичная настройка директорий:")
		fmt.Println("  gunpack-new: директория с новыми файлами для gunpack, копируются в gunpack-old")
		fmt.Println("  gunpack-old: директория, куда копируются файлы gunpack")
		fmt.Println("  redux-new:   директория с новыми файлами для redux, копируются в redux-old")
		fmt.Println("  redux-old:   директория, куда копируются файлы redux")

		cfg.GunpackNew = prompt("Введите путь к директории gunpack-new: ")
		cfg.GunpackOld = prompt("Введите путь к директории gunpack-old: ")
		cfg.ReduxNew = prompt("Введите путь к директории redux-new: ")
		cfg.ReduxOld = prompt("Введите путь к директории redux-old: ")

		cfg.GtaExePath = prompt(fmt.Sprintf("Введите полный путь к GTA5.exe (по умолчанию: %s): ", defaultGTAPathForUser()))
		if cfg.GtaExePath == "" {
			cfg.GtaExePath = defaultGTAPathForUser()
		}

		auto := prompt("Использовать автозапуск при старте ПК? (y/n): ")
		if auto == "y" || auto == "Y" {
			cfg.AutoRun = true
			exePath, _ := os.Executable()
			if err := setAutoRun(true, exePath+" -autostart"); err != nil {
				log.Printf("Ошибка установки автозапуска: %v", err)
			}
		}

		if err := saveConfig(cfg, configPath); err != nil {
			log.Printf("Ошибка сохранения конфигурации: %v", err)
		}

		// После ввода настроек скрываем консоль и запускаем цикл копирования
		hideConsole()
		runCopyingLoop("GTA5.exe", cfg.GtaExePath, cfg, "")
		log.Println("Программа завершена (первичный запуск).")
		return
	} else {
		// Конфигурация существует – загружаем её
		cfg, err = loadConfig(configPath)
		if err != nil {
			log.Printf("Ошибка загрузки конфигурации: %v", err)
			return
		}
	}

	// Определяем режим запуска по аргументу "-autostart"
	isAutoStart := false
	if len(os.Args) > 1 && os.Args[1] == "-autostart" {
		isAutoStart = true
	}

	if cfg.GtaExePath == "" {
		cfg.GtaExePath = defaultGTAPathForUser()
	}

	if isAutoStart {
		hideConsole()
		runCopyingLoop("GTA5.exe", cfg.GtaExePath, cfg, "")
		log.Println("Автозапуск: программа завершилась.")
	} else {
		fmt.Println("Программа запущена, все изменения вносите в config.json")
		go runCopyingLoop("GTA5.exe", cfg.GtaExePath, cfg, "")
		time.Sleep(10 * time.Second)
		log.Println("Ручной запуск: скрываем консоль через 10 секунд.")
		hideConsole()
		select {}
	}
}
