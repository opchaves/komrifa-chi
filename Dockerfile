FROM golang:1.22-alpine3.19 as build
# Git is required for fetching the dependencies.
RUN apk update && apk add --no-cache make git
# Create appuser.
ENV USER=appuser
ENV UID=10001

# See https://stackoverflow.com/a/55757473/12429735RUN
RUN adduser \
    --disabled-password \
    --gecos "" \
    --home "/nonexistent" \
    --shell "/sbin/nologin" \
    --no-create-home \
    --uid "${UID}" \
    "${USER}"

WORKDIR /go/src/app
COPY go.mod go.sum .
RUN go mod download
RUN go mod verify
COPY . .
ENV CGO_ENABLED=0 GOOS=linux GOPROXY=direct
RUN go build -v -o app .

# FROM alpine:3.19
FROM scratch
ENV APP_ENV=production
COPY --from=build /etc/passwd /etc/passwd
COPY --from=build /etc/group /etc/group
COPY --from=build /go/src/app/app /go/bin/app
ENTRYPOINT ["/go/bin/app"]
