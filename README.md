# Logi-Tracker
Logistics tracker for [Foxhole Game](https://foxholegame.com).

This is backend and front end mono repo. Backend is written in go and located in `go`, the front end is written in typescript and located in `ts`.

A simple API test python script named `api-test.py` is in this directory .

This is a work in progress.

I aim to support:
- [x] Invitation link to invite clans for admin
- [x] Invitation link to invite members for clan leader
- [] Create, read, update, delete of stockpiles as well as its contents
- [] Reminder of soon-to-be expired stockpiles
- [] Stockpile sharing between friendly clans
- [] Stockpile history inside clan

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