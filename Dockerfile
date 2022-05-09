FROM golang:1.18 as build

WORKDIR /build
COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . ./
RUN go build -a -o app

FROM ubuntu:20.04
COPY --from=build /build /
EXPOSE 9717
ENTRYPOINT [ "/app" ]
