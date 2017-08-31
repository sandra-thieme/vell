FROM rkcpi/docker-golang-createrepo

USER nobody

RUN mkdir -p /go/src/github.com/rkcpi/vell
WORKDIR /go/src/github.com/rkcpi/vell

COPY . /go/src/github.com/rkcpi/vell
RUN go-wrapper download
RUN go-wrapper install

ENV VELL_HTTP_PORT=8080
ENV VELL_REPOS_PATH=/tmp

EXPOSE 8080

CMD ["go-wrapper", "run"]
