FROM alpine:3.7
ENV APP_ENV production
RUN apk add --no-cache ca-certificates tzdata && cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime
RUN mkdir -p /app/bin /app/public
COPY ./bin/prin /app/bin
COPY ./public /app/public
EXPOSE 8980
WORKDIR /app/bin
CMD ["./prin"]