FROM alpine

COPY config.yaml /
COPY gstp-linux /gstp

ENTRYPOINT [ "gstp" ]
