# hello-app

Tiny Go webserver that prints hostname and version as HTML or JSON with health and prometheus endpoints.

There are four versions available which is handy for testing different scenario's:

- `1.0.0`: OK
- `1.0.1`: FAIL, error while starting.
- `1.0.2`: FAIL, will start but health check will report `HTTP 503 unhealthy`. 
- `1.0.3`: OK
- `2.0.0`: OK

## Usage

### Paths

- `/`: prints hostname and version as HTML.
- `/api`: prints hostname and version as JSON.
- `/health`: health check which will return `HTTP 200 Healthy`.
- `/down`: mark the server as down, this will result in a HTTP 503 Unhealthy response for `/health`.

### Flags

- `bind`: the socket to bind to (default is `:8080`).

### Environment

- `BACKGROUND_COLOR`: you can optionally set a different background color for the HTML page. This can be useful when testing blue/green scenarios. 

## Examples

Run hello app:

```bash
$ docker run -d -P --name hello-app avthart/hello-app:1.0.0
```

If you would like to have a different background-color in the HTML page, run with `-e BACKGROUND_COLOR=blue`.

Get host and port:

```bash
$ HELLO_PORT=`docker port hello-app 8080`
```

Get HTML using [httpie](https://httpie.org/):

```bash
$ http $HELLO_PORT
HTTP/1.1 200 OK
Content-Length: 144
Content-Type: text/html; charset=utf-8
Date: Fri, 27 Sep 2019 12:50:14 GMT

<html><head><title>Hello World</title></head><body style="background-color: white"><h1>Hello from 2a80a0a5eac3 version v1.0.0</h1></body></html>     
```

Get API

```bash
$ http $HELLO_PORT/api
HTTP/1.1 200 OK
Content-Length: 47
Content-Type: application/json
Date: Fri, 27 Sep 2019 13:01:27 GMT

{
    "Hostname": "22e72af228ab",
    "Version": "v1.0.0"
}
```

Get health check report:

```bash
$ http $HELLO_PORT/api
HTTP/1.1 200 OK
Content-Length: 7
Content-Type: text/plain; charset=utf-8
Date: Fri, 27 Sep 2019 12:50:05 GMT

Healthy
```

Mark server as down:

```bash
$ http POST $HELLO_PORT/down

HTTP/1.1 200 OK
Content-Length: 0
Date: Fri, 27 Sep 2019 12:51:08 GMT
```

Get health check report again:

```bash
$ http $HELLO_PORT/api
HTTP/1.1 503 Service Unavailable
Content-Length: 9
Content-Type: text/plain; charset=utf-8
Date: Fri, 27 Sep 2019 12:51:35 GMT

Unhealthy
```
