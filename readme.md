# htadmin

* htadmin provides a simple API interface to manage the `.htpasswd` file.
* The `.htpasswd` file can be configured to restrict access with HTTP basic authentication for Nginx web servers. Refer: <https://docs.nginx.com/nginx/admin-guide/security-controls/configuring-http-basic-authentication/>
* The credentials are automatically loaded, no need to reload the Nginx web server.
* The `.htpasswd` file will be generated in the same directory as the htadmin executable is running.

---

## Compile `htadmin` binary for Linux

```shell
env CGO_ENABLED=0 go build -a -gcflags=all="-l -B" -ldflags "-s -w" .
```

---

## Create the API credentials file

* htadmin reads the credentials file `creds.yaml` from the same directory where the binary is executed from.
* modify the `creds.yaml` file as per the requirement

```shell
users:
  pulse: changeme
  admin: changeme
```

*NOTE: Replace the `changeme` with your actual password.*

---

## Start the API server

```shell
./htadmin
```

---

## Generate the base64 encoded credentials

For example, if your username is `pulse` and the password is `changeme`, you have to encode the string `pulse:changeme` into base64, which is `cHVsc2U6Y2hhbmdlbWU=` .

---

## Create a new user using curl

```shell
curl -i -X  POST http://localhost:19978/uac/create \
  -H 'authorization: Basic cHVsc2U6Y2hhbmdlbWU=' \
  -H "Accept: application/json" -H "Content-type: application/json" \
  -d '{ "name": "acme" }'
```

Output:

```shell
HTTP/1.1 200 OK
Content-Type: text/plain; charset=utf-8
Date: Tue, 22 Nov 2022 11:55:32 GMT
Content-Length: 40

Success
5L0XC7ARAo4qStt18v91rkUsydrBwW6s
```

---

## Delete an existing user using curl

```shell
curl -i -X  POST http://localhost:19978/uac/delete \
  -H 'authorization: Basic cHVsc2U6Y2hhbmdlbWU=' \
  -H "Accept: application/json" -H "Content-type: application/json" \
  -d '{ "name": "acme" }'
```

Output:

```shell
HTTP/1.1 200 OK
Content-Type: text/plain; charset=utf-8
Date: Tue, 22 Nov 2022 11:57:12 GMT
Content-Length: 34

User "acme" deleted successfully
```

---
