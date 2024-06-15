import json
import time

import dash
from dash import html, dcc
import dash_cytoscape as cyto
from dash.dependencies import Input, Output

def load_data():
    with open('/opt/data/data.json', 'r') as file:
        return json.load(file)

def create_elements(data):
    elements = []
    nodes = set()
    edges = set()

    for source_ip, connections in data.items():
        for connection in connections:
            dest_ip = connection['dest']
            protocol = connection['protocol']
            port = connection['port']

            if source_ip not in nodes:
                elements.append({'data': {'id': source_ip, 'label': source_ip}})
                nodes.add(source_ip)
            if dest_ip not in nodes:
                elements.append({'data': {'id': dest_ip, 'label': dest_ip}})
                nodes.add(dest_ip)

            # Добавление связи (только если ее еще нет)
            edge = (source_ip, dest_ip, port)
            if edge not in edges:
                elements.append({
                    'data': {
                        'source': source_ip,
                        'target': dest_ip,
                        'label': f"{protocol} {port}"
                    }
                })
                edges.add(edge)
    return elements

# приложение
app = dash.Dash(__name__, suppress_callback_exceptions=True)

app.layout = html.Div([
    cyto.Cytoscape(
        id='cytoscape',
        style={'width': '100%', 'height': '600px'},
        layout={'name': 'circle'},
        stylesheet=[
            {
                'selector': 'node',
                'style': {
                    'content': 'data(label)',
                    'text-valign': 'center',
                    'color': 'black',
                    'background-color': 'lightblue',
                    'width': '50px',
                    'height': '50px',
                    'font-size': '12px'
                }
            },
            {
                'selector': 'edge',
                'style': {
                    'label': 'data(label)',
                    'line-color': 'gray',
                    'width': 2,
                    'font-size': '10px',
                    'text-background-opacity': 1,
                    'text-background-color': '#ffffff',
                    'text-background-shape': 'round',
                    'target-arrow-color': 'gray',
                    'target-arrow-shape': 'triangle',
                    'curve-style': 'bezier'
                }
            }
        ]
    ),
    dcc.Interval(
        id='interval-component',
        interval=5*1000,  # Обновление каждые 5 секунд
        n_intervals=0
    )
])

@app.callback(
    Output('cytoscape', 'elements'),
    Input('interval-component', 'n_intervals')
)
def update_elements(n):
    data = load_data()
    elements = create_elements(data)
    return elements

if __name__ == '__main__':
    time.sleep(7)
    app.run_server(debug=True, host='0.0.0.0', port=8080)
