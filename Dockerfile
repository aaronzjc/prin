FROM alpine:3.7 as commander
ENV APP_ENV production
RUN apk add --no-cache ca-certificates tzdata && cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime
RUN mkdir -p /app/bin
COPY ./bin/prin /app/bin
EXPOSE 8980
WORKDIR /app/bin
CMD ["./prin"]