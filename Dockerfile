# build stage
FROM golang:alpine AS build-env
WORKDIR /go/src/github.com/kunaldawn/mail.af/
ADD . .
RUN go build -o mail.af cmd/af.go

# final stage
FROM alpine
WORKDIR /app
COPY --from=build-env /go/src/github.com/kunaldawn/mail.af/mail.af /app/
COPY --from=build-env /go/src/github.com/kunaldawn/mail.af/af.config.json /app/
CMD /app/mail.af