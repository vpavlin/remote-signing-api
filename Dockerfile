FROM registry.access.redhat.com/ubi9/ubi:9.0.0-1640

WORKDIR /opt/remote-signing-api/

ADD remote-signing-api scripts/entrypoint.sh /opt/remote-signing-api/
ADD https://github.com/vpavlin/fly-helper/releases/latest/download/flyhelper /bin/flyhelper

RUN chmod +x entrypoint.sh /bin/flyhelper

VOLUME /opt/remote-signing-api/data

EXPOSE 4444

ENV CONFIG_FILE /opt/remote-signing-api/config.json


ENTRYPOINT ./entrypoint.sh