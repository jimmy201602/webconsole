FROM golang:1.8.3-alpine

MAINTAINER Eric Shi <shibingli@realclouds.org>
ADD . /go/
RUN make build
RUN rm -Rf /go/cmd
RUN rm -Rf /go/server
RUN rm -Rf /go/utils
RUN rm -Rf /go/webconsole
EXPOSE 8080

CMD ["/go/bin/apibox","start"]