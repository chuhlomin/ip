FROM --platform=${TARGETPLATFORM:-linux/amd64} golang:1.20 as build-env

WORKDIR /go/src/app
ADD . /go/src/app

RUN go get -d -v ./...
RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -ldflags="-w -s" -mod=vendor -buildvcs -o /go/bin/app


FROM --platform=${TARGETPLATFORM:-linux/amd64} gcr.io/distroless/static:966f4bd97f611354c4ad829f1ed298df9386c2ec
# gcr.io/distroless/static:966f4bd97f611354c4ad829f1ed298df9386c2ec + GeoLite2
# latest-amd64 -> 966f4bd97f611354c4ad829f1ed298df9386c2ec
# https://github.com/GoogleContainerTools/distroless/tree/master/base

LABEL name="ip"
LABEL repository="https://github.com/chuhlomin/ip"
LABEL homepage="https://github.com/chuhlomin/ip"
LABEL maintainer="Konstantin Chukhlomin <mail@chuhlomin.com>"

COPY --from=build-env /go/bin/app /bin/app
COPY favicon.ico og.png /

CMD ["/bin/app"]
