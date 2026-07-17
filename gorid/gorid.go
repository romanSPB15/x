// Пакет gorid предоставляет функции для удобной идентификации горутин.
package gorid

import (
	"runtime"
	"strconv"
	"strings"
	"sync"
)

var m sync.Map

// GetID возвращает идентификатор текущей горутины.
func GetID() int {
	var buf [128]byte
	n := runtime.Stack(buf[:], false)
	line := string(buf[:n])

	parts := strings.SplitN(line, " ", 3)
	if len(parts) < 2 {
		return -1
	}
	id, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		return -1
	}
	return int(id)
}

// Register регистрирует текущую горутину с указанным именем.
// Если имя уже существует, оно будет перезаписано.
func Register(name string) {
	m.Store(name, GetID())
}

// GetGoroutineID возвращает идентификатор горутины с указанным именем.
func GetGoroutineID(name string) (int, bool) {
	v, ok := m.Load(name)
	if ok {
		return v.(int), true

	}
	return 0, false
}

// Current возвращает название текущей горутины.
func Current() string {
	id := GetID()
	name := ""
	m.Range(func(key, value any) bool {
		if value == id {
			name = key.(string)
			return false
		}
		return true
	})
	return name
}

// Is возвращает, является ли эта горутина name.
func Is(name string) bool {
	v, ok := m.Load(name)
	if !ok {
		return false
	}
	return GetID() == v
}

// Must паникует, если эта горутина не является name.
// Если имя не найдено, возвращает пустую строку.
func Must(name string) {
	if !Is(name) {
		got := Current()
		if got == "" {
			got = "gor" + strconv.Itoa(GetID())
		}
		panic("gorid: expected " + name + ", got " + got)
	}
}
