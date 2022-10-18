FROM registry.access.redhat.com/ubi9/ubi:9.0.0-1640

ADD remote-signing-api /opt/remote-signing-api/remote-signing-api

VOLUME /opt/remote-signing-api/data

EXPOSE 4444

ENV CONFIG_FILE /opt/remote-signing-api/config.json

WORKDIR /opt/remote-signing-api/

CMD /opt/remote-signing-api/remote-signing-api ${CONFIG_FILE}