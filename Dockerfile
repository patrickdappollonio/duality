FROM busybox:stable as usergen
RUN echo "nobody:x:65534:65534:Nobody:/:" > /etc/nobody

FROM busybox:stable
COPY --from=usergen /etc/nobody /etc/nobody
USER nobody
WORKDIR /app
ENV DUALITY_SHELL_COMMAND="/bin/sh -c"
COPY duality /usr/local/bin/duality
CMD ["duality"]
