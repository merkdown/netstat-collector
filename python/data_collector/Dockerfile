FROM alpine:latest

RUN apk update && apk upgrade
RUN apk add --no-cache python3 py-pip
RUN rm -vrf /var/cache/apk/*

RUN chmod o+rwx /opt/

RUN adduser data_collector -D
USER data_collector
COPY ./main.py /opt/
WORKDIR /opt/
RUN mkdir /opt/data

ENTRYPOINT [ "python3", "main.py" ]