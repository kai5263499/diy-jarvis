FROM kai5263499/diy-jarvis-deepspeech-models as models
FROM kai5263499/diy-jarvis-builder as builder

COPY / /go/src/github.com/kai5263499/diy-jarvis

RUN cd /go/src/github.com/kai5263499/diy-jarvis/cmd/deepspeech && \
    go build && \
    ldd deepspeech | tr -s '[:blank:]' '\n' | grep '^/' | \
    xargs -I % sh -c 'mkdir -p $(dirname deps%); cp % deps%;'

RUN cd /go/src/github.com/kai5263499/diy-jarvis/cmd/deepspeech && \
    mv deps/deepspeech/libdeepspeech.so deps/lib/libdeepspeech.so && \
    rm -rf deps/deepspeech

FROM ubuntu:18.04

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
ENV TEXT_PROCESSOR_ADDRESS=""

COPY --from=builder /go/src/github.com/kai5263499/diy-jarvis/cmd/deepspeech/deps /
COPY --from=builder /go/src/github.com/kai5263499/diy-jarvis/cmd/deepspeech/deepspeech /deepspeech
COPY --from=models /deepspeech_models /deepspeech_models

EXPOSE 6000

ENTRYPOINT [ "/deepspeech" ]