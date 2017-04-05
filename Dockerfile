FROM debian

COPY ./build/server /server
COPY ./init.sh /init.sh
RUN chmod +x ./init.sh
RUN apt update
RUN apt install -y ca-certificates
ENTRYPOINT ["/init.sh"]
