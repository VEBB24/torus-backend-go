FROM debian
RUN apt update
RUN apt install -y ca-certificates
COPY ./hdfs /hdfs
COPY ./build/server /server
COPY ./init.sh /init.sh
RUN chmod +x ./hdfs
RUN chmod +x ./init.sh
ENTRYPOINT ["/init.sh"]
