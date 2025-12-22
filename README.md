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

1. Copy the entire contents of [`openlist-proxy.js`](./openlist-proxy.js).

2. Navigate to the Cloudflare Workers Dashboard.

3. Create a new Worker or edit an existing one.

4. Paste the copied code into the worker editor.

5. Configure the following environment variables:

   - `ADDRESS`: OpenList backend server address (do not include trailing slash)

   - `TOKEN`: API access token (secret key) for OpenList server

   - `DISABLE_SIGN`: Whether to disable signature verification (recommended to set as false)

6. Save and deploy your worker.

### Cloudflare Pages:

1. Rename `openlist-proxy.js` to `_index.js`.

2. Upload the renamed file to your Pages project.

3. In the Cloudflare Pages dashboard, navigate to Settings → Environment variables.

4. Add the same environment variables as described in the [Cloudflare Workers](#cloudflare-workers) section.

5. Deploy your Pages project.

### EdgeOne Pages:

1. Create a `functions` directory in your project root (if it doesn't exist).

2. Copy [openlist-proxy.js](./openlist-proxy.js) to the path `/functions/[[default]].js`.

   - **Important:** The file must be named `[[default]].js` and placed inside the `functions` folder.

3. In the EdgeOne Pages dashboard, configure the environment variables as described in the [Cloudflare Workers](#cloudflare-workers) section.

4. Deploy your project.
