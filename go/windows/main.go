package windows_DevOps_exporter

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

// Структура вывода netstat
type Connection struct {
	Proto       string `json:"proto"`
	LocalAddr   string `json:"local"`
	ForeignAddr string `json:"foreign"`
	State       string `json:"state"`
}

// Структура вывода в веб интерфейс
type NetstatData struct {
	LastUpdated time.Time    `json:"update"`
	Connections []Connection `json:"connections"`
}

var (
	netstatData NetstatData
	mu          sync.Mutex // сахар в го для гоферов, нужен для защиты переменной от использования при изменении
)

// getNetstat выполняет команду netstat и обновляет переменную netstatData.
func getNetstat() {
	for {
		tcpConnections, err := runNetstatCommand("TCP")
		if err != nil {
			log.Println("Error executing netstat for TCP:", err)
			time.Sleep(5 * time.Second)
			continue
		}

		udpConnections, err := runNetstatCommand("UDP")
		if err != nil {
			log.Println("Error executing netstat for UDP:", err)
			time.Sleep(5 * time.Second)
			continue
		}

		newConnections := append(tcpConnections, udpConnections...)

		mu.Lock()
		if !compareConnections(netstatData.Connections, newConnections) {
			netstatData.Connections = newConnections
			netstatData.LastUpdated = time.Now()
		}
		mu.Unlock()
		time.Sleep(5 * time.Second)
	}
}

// runNetstatCommand выполняет команду netstat для указанного протокола (TCP или UDP) и возвращает список соединений.
func runNetstatCommand(proto string) ([]Connection, error) {
	cmd := exec.Command("netstat", "-anp", proto)
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(output), "\n")
	var connections []Connection
	for _, line := range lines[4:] { // Пропускаем первые четыре строки заголовков
		fields := strings.Fields(line)
		if len(fields) >= 4 {
			localAddr := fields[1]
			foreignAddr := fields[2]
			state := ""
			if proto == "TCP" {
				state = fields[3]
			}
			conn := Connection{
				Proto:       proto,
				LocalAddr:   localAddr,
				ForeignAddr: foreignAddr,
				State:       state,
			}
			connections = append(connections, conn)
		}
	}
	return connections, nil
}

// compareConnections сравнивает два среза соединений.
func compareConnections(a, b []Connection) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
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
