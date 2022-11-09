### Intro

`Prin` is a batch of tools aim at improving my daily development or learning experience. 

Some tools like `json`, may be everywhere, but there are really some details I am not satisfied. so i did it myself.

Here are some I recommended

#### Qrcode

+ generate a qrcode based on ur input.
+ supprot 10 histroy records.

#### Json

+ encode or decode json strings.
+ support multiple tabs.

#### Cert

+ generate certs for https local testing.

#### Iptables

+ iptables beautify.

### Deploy 

```shell
$ docker run -d --name prin -p 8980:8980 aaronzjc/prin:latest
```

then open your browser, navigate to `http://127.0.0.1:8980`.

### License

MIT