# Fixed versions mean reproducibility
FROM golang:alpine3.11 AS go_api_builder
LABEL maintainer="Yevhen Orestov e.orestov@gmail.com"

WORKDIR $GOPATH/src/gobot-api
COPY . .
RUN go get ./...
RUN CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -a -o /bin/gobot-api main.go

FROM alpine:3.9.6
RUN ls -lha
COPY wait-for.sh /bin/wait-for.sh
COPY --from=go_api_builder /bin/gobot-api /bin/gobot-api
RUN mkdir -p /var/gobot/img 
EXPOSE 8000
RUN ["chmod", "+x", "/bin/wait-for.sh"]