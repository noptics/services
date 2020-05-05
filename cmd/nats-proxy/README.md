nats-proxy
---
Work around for nats monitoring endpoint's lack of CORS headers.

This is a very basic proxy that will listen for incoming `GET` http requests to the `/proxy` path. It will intern do a `GET` request to whever URL is configure in the `url` query parameter and return the results.

In each request it sets the CORS `Access-Control-Allow-*` headers to facilitate accessing the nats monitoring endpoints for clusters that do not live on the noptics-ui domain.