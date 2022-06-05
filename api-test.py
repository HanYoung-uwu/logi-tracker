import requests

url = "http://127.0.0.1:6880/api"

user = "admin"
password = "12345678"

# workaround cookie domain restrain when testing
cookie = None

def print_res(res):
    print(res.status_code)
    print(f"\n {res.text}")

def create_admin():
    res = requests.post(url + "/admin/create_admin", json={"Name": user, "Password": password})
    print_res(res)

def login():
    global cookie
    res = requests.post(url + "/login", json={"Name": user, "Password": password})
    if res.status_code % 100 == 2:
        cookie = {'token': res.cookies.get('token')}
    print_res(res)


def get_all_items():
    res = requests.get(url + "/user/all_items", cookies=cookie)
    print_res(res)

def create_stockpile(location: str, code: str):
    res = requests.post(url+ "/user/create_stockpile", cookies=cookie, json={"code": code, "name": location})
    print_res(res)

# negative size means retrieval
def update_item(item: str, size: int, location: str):
    res = requests.post(url + "/user/update_item", json={"item": item, "size": size, "location": location}, cookies=cookie)
    print_res(res)

def get_all_stockpiles():
    res = requests.get(url + "/user/all_stockpiles", cookies=cookie)
    print_res(res)

def get_clan_invitation():
    import json
    res = requests.get(url + "/clan/invitation", cookies=cookie)
    body = json.loads(res.text)
    print(body)
    res = requests.post(url + "/register", cookies={"token": body["token"]}, json={"Name": "clanman", "Password": "123456789123456789"})
    print_res(res)

def invite_clan():
    import json
    res = requests.post(url + "/admin/invite_clan", cookies=cookie, json={"Clan": "invited_clan"})
    body = json.loads(res.text)
    print(body)
    res = requests.post(url + "/register", cookies={"token": body["token"]}, json={"Name": "clanadmin", "Password": "123456789123456789"})
    print_res(res)
