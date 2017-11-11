FROM centurylink/ca-certs

COPY bin/checks2metrics /checks2metrics

MAINTAINER Chris Fordham <chris@fordham-nagy.id.au>

CMD ["/checks2metrics"]
