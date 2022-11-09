FROM alpine:3.7
ENV APP_ENV production
ENV APP_PATH /app
RUN apk add --no-cache ca-certificates tzdata && cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime
RUN mkdir -p /app/bin /app/public /app/assets
COPY ./dagger/backend/prin /app/bin
COPY ./dagger/frontend /app/public
COPY ./assets /app/assets
EXPOSE 8980
WORKDIR /app/bin
CMD ["./prin"]