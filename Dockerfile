FROM scratch
COPY inacceld /bin/inacceld
ENTRYPOINT ["inacceld"]
