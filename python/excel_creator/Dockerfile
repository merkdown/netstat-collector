FROM python:3.12-alpine

RUN pip install pandas openpyxl
RUN rm -vrf /var/cache/apk/*

RUN chmod o+rwx /opt/

RUN adduser excel -D
USER excel
COPY ./main.py /opt/
WORKDIR /opt/

ENTRYPOINT [ "python3", "main.py" ]