
ARG GO_VERSION=1.24.3

FROM golang:${GO_VERSION}-alpine AS build
WORKDIR /app

ARG MAIN_DIR=cmd/cbrstats
ENV MAIN_DIR=${MAIN_DIR}


COPY go.mod go.sum ./
RUN go mod download


COPY . .


ENV GOTOOLCHAIN=auto


WORKDIR /app/${MAIN_DIR}

RUN CGO_ENABLED=0 go build -o /bankstats .


FROM gcr.io/distroless/static-debian12
COPY --from=build /bankstats /bankstats
EXPOSE 8080
ENTRYPOINT ["/bankstats", "-addr", ":8080"]
