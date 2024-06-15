# pip install pandas openpyxl
import json
from time import sleep

import pandas as pd


def main():
    with open('/opt/data/data.json', 'r') as file:
        data = json.load(file)

    rows = []
    for source_ip, connections in data.items():
        for connection in connections:
            row = {
                "source IP": source_ip,
                "dest IP": connection["dest"],
                "protocol": connection["protocol"],
                "port": connection["port"]
            }
            rows.append(row)

    df = pd.DataFrame(rows)

    df.to_excel("/opt/output/data.xlsx", index=False)

if __name__ == "__main__":
    sleep(7)
    while True:
        main()
        sleep(5)
