FROM ubuntu:focal

RUN apt update && apt install -y ca-certificates tzdata

ADD sofar /
RUN chmod +x /sofar

CMD ["/sofar"]
