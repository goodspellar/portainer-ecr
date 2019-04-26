ARG GO_VERSION=1.12.1

FROM golang:${GO_VERSION}-stretch as builder

RUN mkdir /user && \
    echo 'nobody:x:65534:65534:nobody:/:' > /user/passwd && \
    echo 'nobody:x:65534'> /user/group

RUN apt-get update && apt-get install ca-certificates git

WORKDIR /src

COPY ./go.mod ./go.sum ./
RUN go mod download

COPY ./ ./

RUN CGO_ENABLED=0 go build \
    -installsuffix 'static' \
    -o /app .

FROM scratch as final

COPY --from=builder /user/group /user/passwd /etc/

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

COPY --from=builder /app /app

USER nobody:nobody

ENTRYPOINT [ "/app" ]