FROM golang:1.19-bullseye

ADD . /src

RUN apt update && apt install -y ca-certificates tzdata && \
  cd /src && go build && /bin/mv -vf /src/sofar* /sofar

WORKDIR /
CMD ["/sofar"]
