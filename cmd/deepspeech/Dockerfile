FROM kai5263499/diy-jarvis-builder as builder

COPY / /go/src/github.com/kai5263499/diy-jarvis

RUN cd /go/src/github.com/kai5263499/diy-jarvis/cmd/deepspeech && \
    go build && \
    ldd deepspeech | tr -s '[:blank:]' '\n' | grep '^/' | \
    xargs -I % sh -c 'mkdir -p $(dirname deps%); cp % deps%;'

RUN mkdir -p /deepspeech_models && \
    wget https://github.com/mozilla/DeepSpeech/releases/download/v0.5.1/deepspeech-0.5.1-models.tar.gz && \
    tar zxf deepspeech-0.5.1-models.tar.gz -C / && \
    mv /deepspeech-0.5.1-models /deepspeech_models && \
    rm -rf deepspeech-0.5.1-models.tar.gz

FROM scratch

LABEL MAINTAINER="Wes Widner <kai5263499@gmail.com>"

ENV ALPHABET="/deepspeech_models/alphabet.txt"
ENV LANGUAGE_MODEL="/deepspeech_models/lm.binary"
ENV MODEL="/deepspeech_models/output_graph.pbmm"
ENV TRIE="/deepspeech_models/trie"
ENV LISTEN_PORT="6000"
ENV BEAM_WIDTH="500"
ENV NCEP="26"
ENV NCONTEXT="9"
ENV LM_WEIGHT="0.75"
ENV VALID_WORDCOUNT_WEIGHT="1.85"

COPY --from=builder /go/src/github.com/kai5263499/diy-jarvis/cmd/deepspeech/deps /
COPY --from=builder /go/src/github.com/kai5263499/diy-jarvis/cmd/deepspeech/deepspeech /deepspeech-service
COPY --from=builder /deepspeech_models /deepspeech_models

EXPOSE 6000

ENTRYPOINT [ "/deepspeech-service" ]