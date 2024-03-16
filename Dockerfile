FROM golang:1.19-alpine as build

WORKDIR /go/src/github.com/AliyunContainerService/gpushare-scheduler-extender
COPY . .

RUN go build -o /go/bin/gpushare-sche-extender cmd/main.go

FROM alpine

COPY --from=build /go/bin/gpushare-sche-extender /usr/bin/gpushare-sche-extender

CMD ["gpushare-sche-extender"]
