FROM debian:jessie

LABEL maintainer "https://github.com/blacktop"

ENV GO_VERSION 1.8.1

COPY . /go/src/github.com/maliceio/malice-windows-defender
RUN buildDeps='ca-certificates \
               build-essential \
               cabextract \
               mercurial \
               git-core \
               unzip \
               wget' \
  && apt-get update \
  && apt-get install -y $buildDeps libc6-i386 --no-install-recommends \
  && echo "===> Download 32-bit antimalware update file.." \
  && cd /tmp \
  && wget "http://go.microsoft.com/fwlink/?LinkID=121721&arch=x86" \
  && cabextract mpam-fe.exe \
  && echo "===> Install taviso/loadlibrary..." \
  && git clone https://github.com/taviso/loadlibrary.git /tmp/loadlibrary \
  && cd /tmp/loadlibrary \
  && make \
  && echo "===> Install Go..." \
  && ARCH="$(dpkg --print-architecture)" \
  && wget https://storage.googleapis.com/golang/go$GO_VERSION.linux-$ARCH.tar.gz -O /tmp/go.tar.gz \
  && tar -C /usr/local -xzf /tmp/go.tar.gz \
  && export PATH=$PATH:/usr/local/go/bin \
  && echo "===> Building avscan Go binary..." \
  && cd /go/src/github.com/maliceio/malice-windows-defender \
  && export GOPATH=/go \
  && go version \
  && go get \
  && go build -ldflags "-X main.Version=$(cat VERSION) -X main.BuildTime=$(date -u +%Y%m%d)" -o /bin/avscan \
  && echo "===> Clean up unnecessary files..." \
  && apt-get purge -y --auto-remove $buildDeps $(apt-mark showauto) \
  && apt-get clean \
  && rm -rf /var/lib/apt/lists/* /var/cache/apt/archives /tmp/* /var/tmp/* /go /usr/local/go

# Add EICAR Test Virus File to malware folder
ADD http://www.eicar.org/download/eicar.com.txt /malware/EICAR

WORKDIR /malware

ENTRYPOINT ["/bin/avscan"]
CMD ["--help"]
