FROM golang

ADD dbproxy /

RUN chmod +x /dbproxy

ENTRYPOINT ["sh", "-c", "/dbproxy"]