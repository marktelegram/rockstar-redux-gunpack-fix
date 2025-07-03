package main

import (
	"io"
	"log"
	"os"
	"path/filepath"
	"sync"
)

// copyFile копирует один файл из src в dst, перезаписывая его, если он существует.
// Используется для копирования файлов модов или игровых ресурсов GTA 5 RP.
func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		log.Printf("Ошибка открытия файла %s: %v", src, err)
		return err
	}
	defer in.Close()

	if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
		log.Printf("Ошибка создания директорий для %s: %v", dst, err)
		return err
	}

	out, err := os.Create(dst)
	if err != nil {
		log.Printf("Ошибка создания файла %s: %v", dst, err)
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		log.Printf("Ошибка копирования %s -> %s: %v", src, dst, err)
		return err
	}

	if err = out.Sync(); err != nil {
		log.Printf("Ошибка синхронизации файла %s: %v", dst, err)
		return err
	}

	log.Printf("Файл скопирован: %s -> %s", src, dst)
	return nil
}

// copyDirRecursive рекурсивно копирует содержимое директории src в dst.
// Используется для копирования директорий с модами GTA 5 RP.
// Возвращает ошибку, если она возникла.
func copyDirRecursive(src, dst string, wg *sync.WaitGroup) error {
	entries, err := os.ReadDir(src)
	if err != nil {
		log.Printf("Ошибка чтения директории %s: %v", src, err)
		return err
	}
	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())
		if entry.IsDir() {
			if err := copyDirRecursive(srcPath, dstPath, wg); err != nil {
				return err
			}
		} else {
			wg.Add(1)
			go func(s, d string) {
				defer wg.Done()
				_ = copyFile(s, d)
			}(srcPath, dstPath)
		}
	}
	return nil
}
