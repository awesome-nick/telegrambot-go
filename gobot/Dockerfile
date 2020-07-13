# Fixed versions mean reproducibility
FROM golang:alpine3.11 AS go_builder
LABEL maintainer="Yevhen Orestov e.orestov@gmail.com"

WORKDIR $GOPATH/src/gobot
COPY . .
RUN go get ./...
RUN CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -a -o /bin/gobot main.go

FROM alpine:3.9.6
COPY wait-for.sh /bin/wait-for.sh
COPY --from=go_builder /bin/gobot /bin/gobot
RUN mkdir -p /var/gobot/img 
EXPOSE 80
RUN ["chmod", "+x", "/bin/wait-for.sh"]