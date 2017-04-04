FROM debian

COPY ./build/server /server
COPY ./init.sh /init.sh
RUN chmod +x ./init.sh

ENTRYPOINT ["/init.sh"]
