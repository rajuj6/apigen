FROM golang:1.19 as build

WORKDIR /app
COPY . .
RUN cd /app/


RUN --mount=type=cache,id=gomod_apigen,target=/go/pkg/mod \
    --mount=type=cache,id=gomod_apigen,target=/root/.cache/go-build \
    CGO_ENABLED=0 go build -o /usr/bin/app /app/*.go

FROM alpine
WORKDIR /app

COPY --from=build /usr/bin/app /usr/bin/apigen
RUN chmod +x /usr/bin/apigen