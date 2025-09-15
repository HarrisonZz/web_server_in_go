Execute this cmd:

docker build -t custom-nginx .

docker run --name custom-nginx \
  --add-host=host.docker.internal:host-gateway \ 
  -p 80:80 -p 443:443 \
  -d custom-nginx:app