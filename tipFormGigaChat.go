package main

import (
    "os"
    "path/filepath"
    "strings"
    "time"
)

// findLatestUfwLogFile находит последний по времени создания файл, содержащий 'ufw.log*',
// который не является пустым и не имеет расширение '.gz'.
func findLatestUfwLogFile(directory string) (*os.File, error) {
    // Итерируем по всем файлам в директории
    files, err := filepath.Glob(filepath.Join(directory, "ufw.log*"))
    if err != nil {
        return nil, err
    }

    // Создаем слайс для хранения времен создания файлов
    creationTimes := make([]time.Time, 0, len(files))
    // Создаем слайс для хранения имен файлов
    filenames := make([]string, 0, len(files))

    // Заполняем слайсы временем создания и именами файлов
    for _, filename := range files {
        fi, err := os.Stat(filename)
        if err != nil {
            return nil, err
        }
        creationTimes = append(creationTimes, fi.ModTime())
        filenames = append(filenames, filename)
    }

    // Находим последний по времени создания файл
    latestTime := creationTimes[0]
    latestFile := filenames[0]
    for i, time := range creationTimes {
        if latestTime.Before(time) {
            latestTime = time
            latestFile = filenames[i]
        }
    }

    // Проверяем, что последний файл не является пустым и не имеет расширение '.gz'
    if !strings.HasSuffix(latestFile, ".gz") && isNotEmpty(latestFile) {
        // Открываем найденный файл
        file, err := os.Open(latestFile)
        if err != nil {
            return nil, err
        }
        return file, nil
    }

    // Если последний файл не удовлетворяет условиям, возвращаем ошибку
    return nil, errors.New("no suitable file found")
}

// isNotEmpty проверяет, является ли файл пустым.
func isNotEmpty(filename string) bool {
    file, err := os.Open(filename)
    if err != nil {
        return false
    }
    defer file.Close()

    _, err = file.Seek(0, os.SEEK_END)
    if err != nil {
        return false
    }

    size := file.Stat().Size()
    if size == 0 {
        return false
    }

    file.Seek(0, os.SEEK_SET)
    return true
}

func main() {
    // Пример использования функции
    directory := "/path/to/directory"
    file, err := findLatestUfwLogFile(directory)
    if err != nil {
        // Обработка ошибки
        panic(err)
    }
    defer file.Close() // Не забываем закрыть файл

    // Теперь файл открыт и доступен для чтения
}
