// https://stackoverflow.com/questions/8757389/reading-a-file-line-by-line-in-go
package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
)

var ipAressesMap map[string]int

var (
	EXIT_SUCCESS = 0
	EXIT_FAILURE = 1
)

func main() {
	ipAressesMap = make(map[string]int)
	// Проверяем аргументы командной строки
	if len(os.Args) != 2 {
		fmt.Printf("Usage: %s filename\n", os.Args[0])
		os.Exit(EXIT_FAILURE)
	}

	// Открываем файл, проверяем на ошибки
	file, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal("opening file:", err.Error())
		os.Exit(EXIT_FAILURE)
	}
	defer file.Close()

	// Создаем структуру scanner для дальнейшей обработки текста
	scanner := bufio.NewScanner(file)
	for scanner.Scan() { // перебираем каждую строку файла
		selectIP(scanner.Text()) // вызываем фукцию, которая обрабатывает лог
	}

	if err := scanner.Err(); err != nil {
		log.Fatal("scanner.Err():", err.Error())
	}

	convertIpAddressesMapToJson()

	os.Exit(EXIT_SUCCESS)
}

func selectIP(fullLine string) {

	tokens := strings.Fields(fullLine) // разбиваем строку на токены от пробела до пробела
	for _, srcIP := range tokens {     // перебираем полученный слайс
		if strings.Index(srcIP, "SRC=") != -1 { // если строка содержит SRC=
			IPaddr := strings.TrimPrefix(srcIP, "SRC=")
			ipAressesMap[IPaddr]++
			break
		} else {
			continue
		}
	}
}

func convertIpAddressesMapToJson() {
	jsonIpAddresses, err := json.Marshal(ipAressesMap)
	if err != nil {
		log.Fatal("json map:", err.Error())
	}
	file, err := os.OpenFile("jsonFileLog.json", os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		log.Fatal("opening file:", err.Error())
	}
	defer file.Close()

	_, err = file.Write(jsonIpAddresses)
	if err != nil {
		log.Fatal("os.Write()", err.Error())
	}
}
