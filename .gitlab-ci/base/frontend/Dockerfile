FROM nginx:alpine

RUN apk --no-cache add tzdata
ENV TZ=Asia/Jakarta

RUN rm /etc/nginx/conf.d/default.conf

COPY ./files/etc/nginx/nginx.conf /etc/nginx/conf.d
