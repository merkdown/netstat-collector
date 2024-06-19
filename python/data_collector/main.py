import requests
import json
import time
import os

# Запрос данных с источников
def getData(ip):
    try:
        response = requests.get(f"http://{ip}:11110/devops/ntst").json()
        return response
            
    except requests.RequestException as e:
        print(f"Ошибка при выполнении запроса к {ip}: {e}")
        return 0 

# Функция для чтения данных из JSON-файла
def readJSON(file):
    if os.path.exists(file):
        with open(file, 'r') as file:
                return json.load(file)
    return {}

# Функция для записи данных в JSON-файл
def writeJSON(file, data):
    with open(file, 'w') as file:
        json.dump(data, file, indent=4)
        
# Основная логика
def collectData():
    with open("source_list", "r") as f:
        source = [line.replace("\n", "") for line in f.readlines()]


    current_data = readJSON("/opt/data/data.json")

    for ip in source:
        new_data = getData(ip)["connections"]
        if ip not in current_data:
            current_data[ip] = []
        
        if new_data:
            for connection in new_data:
                connect = {
                    "dest" : connection["foreign"],
                    "protocol" : connection["proto"],
                    "port" : connection["port"],
                }
                if connect not in current_data[ip]:
                    current_data[ip].append(connect)
  
    writeJSON("/opt/data/data.json", current_data)

if __name__ == "__main__":           
    while True:
        collectData()
        time.sleep(5)
