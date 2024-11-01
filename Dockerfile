FROM golang:1.22.6 as builder
WORKDIR /src
COPY go.sum go.mod ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /bin/app ./cmd

FROM alpine:latest 
COPY --from=builder /bin/app /bin/app
COPY --from=builder /src/.env .env
ENTRYPOINT [ "/bin/app" ] 


