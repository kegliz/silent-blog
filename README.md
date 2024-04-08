# Silent Blog

htmx and go based personal blog boilerplate. For motivation and more information, see [this blog post](https://zkegli.kegnet.dev/post/go-gin-htmx-templ-tailwindcss).

- [htmx](https://htmx.org/) for minimal or no js dynamic content
- [gin](https://github.com/gin-gonic/gin) for server
- [templ](https://templ.guide) for templating (like JSX in Go)
- [TailwindCSS](https://tailwindcss.com/blog/standalone-cli) for styling
- [Goldmark](https://github.com/yuin/goldmark) for markdown rendering

## Features
Write a post in markdown, and place it to the prod/posts directory. Describe the post in the prod/posts.json file. The post will be rendered in the blog. 

## Development
- Install go, make, templ and tailwindcss

Test the blog locally with 
```bash
make run
```
This will start the server on localhost:3049

## Deployment
The most basic approach is to build the binary for the target architecture and deploy it to the server. 

### Reverse Proxy
One approach is that server has a reverse proxy like nginx installed and configured to serve our service. For example, if the binary runs with default config on localhost:3049, and nginx is configured to serve the binary with your domain. Then the steps for installing and configuring the server are as follows:
- config server with firewall, install nginx: 
https://www.digitalocean.com/community/tutorials/initial-server-setup-with-ubuntu, 
https://www.digitalocean.com/community/tutorials/how-to-install-nginx-on-ubuntu-22-04
- config nginx, install certbot: 
https://www.digitalocean.com/community/tutorials/how-to-secure-nginx-with-let-s-encrypt-on-ubuntu-22-04, 
https://www.digitalocean.com/community/tutorials/how-to-configure-nginx-as-a-reverse-proxy-on-ubuntu-22-04
- in the nginx config use location http://127.0.0.1:3049

- build architecture specific binary
```bash
make build
```
- copy the prod directory to the server
```bash
scp -r /prod ubuntu@YOUR_DOMAIN_NAME:/home/ubuntu/
```
- run binary with default config
```bash
cd prod
./app
```

### Standalone
Another approach is to build the binary for the target architecture and run it on the server.
In this case you should run it with  the following config parameters in the config.yaml file if you want to use https.
```bash
localonly: False
tls: True
domain: "your-domain.dev"
# port: 443
```

### Docker

## VSCode with TailwindCSS
settings.json
``` 
"tailwindCSS.includeLanguages": {
    "templ": "html"
}
```

## TODO
- containerize 
- automate deployment
