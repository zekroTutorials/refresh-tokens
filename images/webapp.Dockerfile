FROM node:14-alpine as build
WORKDIR /build
COPY ./webapp .
RUN apk add git
RUN npm ci &&\
    npm run build

FROM nginx:mainline-alpine AS final
WORKDIR /app
COPY --from=build /build/build .
ADD ./config/nginx.conf /etc/nginx/conf.d/webapp.conf
EXPOSE 80