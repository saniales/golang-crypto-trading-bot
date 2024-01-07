FROM golang:1.21-alpine3.19 as builder

ENV DIR_SRC=/src
ENV DIR_OUT=/build

RUN mkdir -p ${DIR_SRC}
RUN mkdir -p ${DIR_OUT}

WORKDIR ${DIR_SRC}

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -o ${DIR_OUT}/trading-bot .

FROM alpine:3.19

ENV DIR_APP=/app

RUN mkdir -p ${DIR_APP}

WORKDIR ${DIR_APP}

COPY --from=builder /build/trading-bot ${DIR_APP}/trading-bot

CMD ${DIR_APP}/trading-bot start --config-file ${DIR_APP}/config.yaml
