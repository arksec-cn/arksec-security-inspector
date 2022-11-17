ARG OS_ARCH="amd64"
ARG PROJECT_ENV="dev"

FROM registry.arksec.cn/library/${OS_ARCH}/bigbang/${PROJECT_ENV}:latest

WORKDIR /
COPY dist/inspector /

CMD ["/inspector"]