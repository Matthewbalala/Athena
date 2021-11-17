import requests

url = "http://127.0.0.1:5000/cdt"

response = requests.request("GET", url + "/action/default")

print("action default: ", response.text)

headers = {
        'Content-Type': 'application/json'
        }
requests.request("POST", url + "/deploy/up", headers=headers, data=response.text)

print("done")