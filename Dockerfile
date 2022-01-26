FROM moby/buildkit:v0.9.3
WORKDIR /cms
COPY cms README.md /cms/
ENV PATH=/cms:$PATH
ENTRYPOINT [ "/bhojpur/cms" ]