FROM python:3.12-alpine

RUN pip install dash_cytoscape dash
RUN rm -vrf /var/cache/apk/*

RUN chmod o+rwx /opt/

RUN adduser viz -D
USER viz
COPY ./main.py /opt/
WORKDIR /opt/

ENTRYPOINT [ "python3", "main.py" ]