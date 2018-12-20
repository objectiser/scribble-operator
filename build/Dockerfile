FROM alpine:3.8

USER nobody

ADD build/_output/bin/scribble-operator /usr/local/bin/scribble-operator
ENTRYPOINT ["/usr/local/bin/scribble-operator"]
