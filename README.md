# oed-proxy

Oxford English Dictionary Proxy

[![Deploy with Vercel](https://vercel.com/button)](https://vercel.com/new/git/external?repository-url=https%3A%2F%2Fgithub.com%2Fserizawa-jp%2Foed-proxy&project-name=your-oed-proxy)

Get your own proxy on [Vercel](https://vercel.com/).

[oed-proxy.vercel.app](https://oed-proxy.vercel.app/api)

## Usage

```bash
$ curl -X POST -H "Content-Type: application/json" -d '{"app_key": "XXX","app_id":"YYY","word": "go"}' https://oed-proxy.vercel.app/api
```
