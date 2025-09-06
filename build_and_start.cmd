docker build . --tag modbustohttp:latest
docker rm -f modbustohttp
docker run -d -e HTTP_HOST= -e HTTP_PORT=8080 -p 8080:8080 --name modbustohttp modbustohttp:latest