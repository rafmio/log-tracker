package parser

import (
	"testing"
	"time"
)

// Сначала создаем функцию, которая принимает строку файла - запись в журнале
// логов, парсим ее на структуру или карту

var logEntry struct {
	Date   time.Time
	SrcIP  string
	Len    string
	Ttl    string
	Id     string
	Spt    string
	Dpt    string
	Window string
}

func TestParseLog(t *testing.T) {

}
