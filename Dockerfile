FROM golang:1.17.1-alpine

MAINTAINER Eric Shi <shibingli@realclouds.org>
RUN apk add make --no-cache
ADD . /go/
RUN make build
RUN rm -Rf /go/cmd
RUN rm -Rf /go/server
RUN rm -Rf /go/utils
RUN rm -Rf /go/webconsole
EXPOSE 8080

CMD ["/go/bin/apibox","start"]
