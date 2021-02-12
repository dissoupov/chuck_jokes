FROM centos:7
LABEL version="1.0"

ENV JOKES_DIR=/jokes
ENV PATH=$PATH:/jokes/bin

ADD ./bin/jokes  /jokes/bin/
ADD ./etc/  /jokes/etc

CMD ["/jokes/bin/jokes", "--std"]

EXPOSE  5000 8080

VOLUME /var/jokes/certs