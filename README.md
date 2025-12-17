# OpenList-Proxy

Use another machine to proxy OpenList's traffic.

## CLI Usage (Go)

```shell
Usage of OpenList-Proxy:
  -address string
        openlist address
  -cert string
        cert file (default "server.crt")
  -disable-sign
        disable signature verification
  -help
        show help
  -https
        use https protocol.
  -key string
        key file (default "server.key")
  -port int
        the proxy port. (default 5243)
  -token string
        openlist token
  -version
        show version and exit
```

## Functions Deployment (JavaScript)

### Cloudflare Workers

1. Paste content from [openlist-proxy.js](./openlist-proxy.js) to Cloudflare Workers Dashboard.

2. Setting the following environment in the dashboard:

   - `ADDRESS`: OpenList backend server address (do not include trailing slash)

   - `TOKEN`: API access token (secret key) for OpenList server

   - `WORKER_ADDRESS`: Full address of your Cloudflare Worker (starts with 'http://' or 'https://')

   - `DISABLE_SIGN`: Whether to disable signature verification (recommended to set as false)

### Cloudflare Pages:

1. Upload [openlist-proxy.js](./openlist-proxy.js) under the name `_index.js`.

2. Setting the environment (refer to [Cloudflare Workers](#cloudflare-workers)) in the dashboard.


### EdgeOne Pages:

1. Upload `openlist-proxy.js` to the path `/functions/[[default]].js`. Please note that the `functions` directory is required.

2. Setting the environment (refer to [Cloudflare Workers](#cloudflare-workers)) in the dashboard.
