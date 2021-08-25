FROM frrouting/frr:v7.5.1
WORKDIR /
RUN mkdir /var/log/frr
COPY frr-reloader.sh .

ENTRYPOINT ["/frr-reloader.sh"]
