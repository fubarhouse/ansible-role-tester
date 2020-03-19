FROM golang:alpine as build

ARG TARGETARCH
ARG GOARCH=$TARGETARCH
ARG GOARM

RUN apk --no-cache add git
WORKDIR /opt
COPY . .

RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o /bin/ansible-role-tester

FROM scratch
USER 65534
ENTRYPOINT ["ansible-role-tester"]
