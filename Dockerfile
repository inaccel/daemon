FROM gcr.io/distroless/base
COPY inacceld /bin/inacceld
ENTRYPOINT ["inacceld"]
