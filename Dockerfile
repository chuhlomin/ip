FROM golang:1.17 as build-env

WORKDIR /go/src/app
ADD . /go/src/app

RUN go get -d -v ./...
RUN CGO_ENABLED=0 go build -o /go/bin/app


FROM gcr.io/distroless/static:966f4bd97f611354c4ad829f1ed298df9386c2ec
# latest-amd64 -> 966f4bd97f611354c4ad829f1ed298df9386c2ec
# https://github.com/GoogleContainerTools/distroless/tree/master/base

LABEL name="ip"
LABEL repository="https://github.com/chuhlomin/ip"
LABEL homepage="https://github.com/chuhlomin/ip"
LABEL maintainer="Konstantin Chukhlomin <mail@chuhlomin.com>"

COPY --from=build-env /go/bin/app /bin/app

CMD ["/bin/app"]
