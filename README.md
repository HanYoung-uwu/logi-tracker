# Logi-Tracker
Logistics tracker for [Foxhole Game](https://foxholegame.com).

This is backend and frontend mono repo. Backend is written in go and located in `go`, the front end is written in typescript and located in `ts`.

A simple API test python script named `api-test.py` is in this directory(not up to date).

This is a work in progress.

I aim to support:
- [x] Invitation link to invite clans for admin
- [x] Invitation link to invite members for clan leader
- [x] Create, read, update, delete of stockpiles as well as its contents
- [ ] Reminder of soon-to-be expired stockpiles
- [ ] Stockpile sharing between friendly clans
- [x] Stockpile history inside clan

## Todo list
- [ ] catalog tab or each section, like guns, shells, vehicles
- [ ] account page for setting password, delete account etc.
- [ ] clan admin manage page for kicking off memebers and add promote memebers to admin
- [x] don't use cookie for account creating
- [ ] cache item icons in local storage
- [ ] stockpile item search
- [ ] ability for clan admin to publish item requests
- [ ] use discord api to notify expire stockpile
- [ ] polish UI (Long term)

## How to run it locally
### Backend
Install go on your system
```bash
cd go/cmd/logi-tracker
go mod tidy
go build
./logi-tracker
```

### Frontend
Install yarn

```bash
cd ts
yarn start
```

### Proxy
Install nginx and put this inside the outmost `{}` block in `/etc/nginx/nginx.conf`
```
server {
        listen          6880;
        server_name     proxy_server;
        location / {
            proxy_pass http://localhost:3000;
        }
        location /api {
            proxy_pass http://localhost:8080;
        }
}
```
This assume the frontend listen on 3000 and backend listen on 8080.
then start nginx with `sudo systemctl restart nginx`

Access the website at **http://127.0.0.1:6880**.
You can use `api-test.py` to create admin account on the first ever server start.
```
python -i api-test.py
>>> create_admin()
```
The default name is "admin" and password "12345678"