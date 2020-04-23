FROM registry.access.redhat.com/ubi8/ubi-minimal:latest

ENV OPERATOR=/usr/local/bin/nanum.co.kr/openmcp/migration \
    USER_UID=1001 \
    USER_NAME=nanum.co.kr/openmcp/migration

# install operator binary
COPY build/_output/bin/nanum.co.kr/openmcp/migration ${OPERATOR}

COPY build/bin /usr/local/bin
RUN  /usr/local/bin/user_setup

ENTRYPOINT ["/usr/local/bin/entrypoint"]

USER ${USER_UID}
