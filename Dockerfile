# build stage
FROM golang:alpine AS build

RUN apk --no-cache add tzdata
RUN apk --no-cache add ca-certificates bash busybox-extras
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app .

# final stage
FROM alpine

WORKDIR /app

RUN apk --no-cache add bash busybox-extras

COPY --from=build /app/app /app
# COPY --from=build /app/.env /app/.env  

RUN apk add -U tzdata
RUN cp /usr/share/zoneinfo/Asia/Jakarta /etc/localtime

RUN chmod +x /app

ENTRYPOINT ["./app"]

EXPOSE 9000
