package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"strings"
	"sync"
	"time"
)

// Структура вывода netstat -anut
type Connection struct {
	Proto       string `json:"proto"`
	Port		string `json:"port"`  // в начальном выводе нет, добавляю для упрощения текстуры вывода
	// RecvQ       string `json:"recv"` // бесполезная информация
	// SendQ       string `json:"send"` // бесполезная информация
	LocalAddr   string `json:"local"`
	ForeignAddr string `json:"foreign"`
	// State       string `json:"state"` // обрабатывается в логике ниже
}

// Структура вывода в веб интерфейс
type NetstatData struct {
	LastUpdated time.Time    `json:"update"` // нужно было при разработке, можно в теории убрать и упростить структуру
	Connections []Connection `json:"connections"`
}

var (
	netstatData NetstatData
	mu          sync.Mutex // сахар в го для гоферов, нужен для защиты переменной от использования при изменении
)

// getNetstat выполняет команду netstat и обновляет переменную netstatData.
func getNetstat() {
	for {
		cmd := exec.Command("netstat", "-anut")
		output, err := cmd.Output()
		if err != nil {
			log.Println("Error executing netstat:", err)
			time.Sleep(5 * time.Second)
			continue
		}

		lines := strings.Split(string(output), "\n")
		var newConnections []Connection
		for _, line := range lines[2:] { // Скипаем заголовок
			fields := strings.Fields(line)
			if len(fields) >= 6 {

				proto := fields[0]
				state := fields[5]
				local := strings.Split(fields[3], ":")[0]

				if proto == "tcp6" || proto == "udp6" { // Скипаем ipv6
					continue 
				}
				if state != "ESTABLISHED" { // Оставляем активные соединения
					continue
				}
				if local == "127.0.0.1" || local == "0.0.0.0"{ // На всякий проверяем на активные соединения
					continue
				}

				conn := Connection{
					Proto:       proto,
					Port: 		 strings.Split(fields[4], ":")[1],
					// RecvQ:       fields[1],
					// SendQ:       fields[2],
					LocalAddr:   local,
					ForeignAddr: strings.Split(fields[4], ":")[0],
					// State:       fields[5],
				}
				newConnections = append(newConnections, conn)
			}
		}

		mu.Lock()
		netstatData.Connections = newConnections
		netstatData.LastUpdated = time.Now()
		mu.Unlock()
		time.Sleep(5 * time.Second)
	}
}

// handler возвращает данные netstat в формате JSON.
func handler(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(netstatData)
}

func main() {
	go getNetstat()

	http.HandleFunc("/devops/ntst", handler)
	fmt.Println("Сервер запущен на http://0.0.0.0:11110/devops/ntst")
	log.Fatal(http.ListenAndServe("0.0.0.0:11110", nil))
}
