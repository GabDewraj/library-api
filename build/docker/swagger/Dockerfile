FROM swaggerapi/swagger-ui AS build_base


COPY ./build/docker/swagger/docs/swagger.json /mnt/swagger.json
ENV SWAGGER_JSON=/mnt/swagger.json

EXPOSE 8080
