FROM golang:1.16-alpine3.12 as builder
WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY static ./static
COPY src ./src
COPY README.md .
WORKDIR ./src
RUN CGO_ENABLED=0 go build -o /aubg-gaming

FROM scratch
COPY --from=builder /app /app
COPY --from=builder /aubg-gaming /aubg-gaming
WORKDIR ./app/src
EXPOSE 8080
ENTRYPOINT [ "/aubg-gaming" ]