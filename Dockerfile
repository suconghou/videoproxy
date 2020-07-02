FROM alpine
ADD videoproxy /
ENTRYPOINT ["/videoproxy"]
EXPOSE 6060