# XDPFail2ban plugin for Traefik

This plugin is a small implementation of a fail2ban instance
drop packets via eBPF XDP as a middleware plugin for Traefik.
There are two components to make it work:
- XDPFail2ban plugin   
This is hosted in this repository.   
It implements the required logic to block a certain source IP.
Once the source IP falls in Ban mode,
it hits XDPDropper component via API to add it there.
- [XDPDropper](https://github.com/renanqts/xdpdropper)   
It implements a dropper source IP via eBFP XDP.   
It has an API to add/remove IP from being dropped.
Once a certain source IP is added to be dropped,
It discards every single packet before it hits the kernel,
allowing for high-performance packet processing.
[See more](https://blog.cloudflare.com/l4drop-xdp-ebpf-based-ddos-mitigations/).

## Configuration
Available configuration
```yml
testData:
  loglevel: DEBUG
  xdpdropperurl: http://localhost:8080
  rules:
    bantime: 3h
    findtime: 10m
    maxretry: 4
    urlregexps:
    - regexp: /foo
    - regexp: /bar
```

Where:
- `logLevel`: is used to show the correct level of logs (`DEBUG`, `INFO` (default),
`NONE`)
- `xdpdropperurl`: url where [XDPDropper](https://github.com/renanqts/xdpdropper)
service can be found.

under `rules`:
- `findtime`: is the time slot used to count requests (if there are too many
requests with the same source IP in this slot of time, the IP goes into ban).
You can use 'smart' strings like "4h", "2m", "1s", ...
- `bantime`: correspond to the amount of time the IP is in Ban mode.
- `maxretry`: number of requests before Ban mode.
- `urlregexps`: a regexp list to requests with regexps on the URL to be filtered.
In this example, all requests to `/foo` and `/bar` will be filtered.

### Ban logic
First request, one of `urlregexps` is matched, the IP is added to the Pool,
and the `findtime` timer is started:
```
A |------------->
  ↑
```

Second request, `urlregexps` is matched again, `findtime` is not yet finished
thus the request is fine:
```
A |--x---------->
     ↑
```

Third request, `urlregexps` is matched, `maxretry` is now almost full, this request
is fine:
```
A |--x--x------->
        ↑
```

Fourth request matched, now it's jail time, `bantime` is started,
the source IP is added into [XDPDropper](https://github.com/renanqts/xdpdropper) list:
```
A |--x--x--x---->
           ↓
B          |------------->
```

Next requests, the IP is in Ban mode, those requests will never arrive
since XDP is dropping it in the level down:
```
A |--x--x--x--x->
```

`bantime` is now expired, another `findtime` is started:
```
A |--x--x--x---->            |------------->
                             ↑
B          |--x---------->
```

## How to dev
Standards tests/lint can be achieved by running:
```bash
make 
```

for integration test, use:
```bash
docker compose up
```
It will bring `http://localhost:8000/whoami` locally to be hit for the sake of tests.
The plugin log will show up prefixed with `XDPFail2Ban`.   
Traefik dashboard is also available at `http://localhost:8080/dashboard`.   
`xdpdropperdummy` is the name of the container for a dummy API that can be found
at `http://localhost:8081/drop`.

## Credits
[fail2ban](https://github.com/tommoulard/fail2ban)