####################################################
# GOLANG BUILDER
####################################################
FROM --platform=linux/amd64 golang:1 as go_builder

COPY . /go/src/github.com/malice-plugins/windows-defender
WORKDIR /go/src/github.com/malice-plugins/windows-defender
RUN go build -ldflags "-s -w -X main.Version=v$(cat VERSION) -X main.BuildTime=$(date -u +%Y%m%d)" -o /bin/avscan

####################################################
# PLUGIN BUILDER
####################################################
FROM --platform=linux/amd64 ubuntu:focal

LABEL maintainer "https://github.com/blacktop"

LABEL malice.plugin.repository = "https://github.com/malice-plugins/windows-defender.git"
LABEL malice.plugin.category="av"
LABEL malice.plugin.mime="*"
LABEL malice.plugin.docker.engine="*"

# Create a malice user and group first so the IDs get set the same way, even as
# the rest of this may change over time.
RUN groupadd -r malice \
  && useradd --no-log-init -r -g malice malice \
  && mkdir /malware \
  && chown -R malice:malice /malware

RUN buildDeps='libreadline-dev:i386 \
  ca-certificates \
  libc6-dev:i386 \
  build-essential \
  gcc-multilib \
  cabextract \
  mercurial \
  git-core \
  unzip \
  curl' \
  && set -x \
  && dpkg --add-architecture i386 && apt-get update -qq \
  && apt-get install -o APT::Immediate-Configure=false -y $buildDeps libc6-i386 --no-install-recommends \
  && echo "===> Install taviso/loadlibrary..." \
  && git clone https://github.com/taviso/loadlibrary.git /loadlibrary \
  && echo "===> Download 32-bit antimalware update file.." \
  && curl -L --output /loadlibrary/engine/mpam-fe.exe "http://download.microsoft.com/download/DefinitionUpdates/mpam-fe.exe" \
  # && curl -L --output /loadlibrary/engine/mpam-fe.exe "https://www.microsoft.com/security/encyclopedia/adlpackages.aspx?arch=x86" \
  && cd /loadlibrary/engine \
  && cabextract mpam-fe.exe \
  && rm mpam-fe.exe \
  && cd /loadlibrary \
  && make -j2 \
  && echo "===> Clean up unnecessary files..." \
  && apt-get purge -y --auto-remove $buildDeps $(apt-mark showauto) \
  && apt-get clean \
  && rm -rf /var/lib/apt/lists/* /var/cache/apt/archives /tmp/* /var/tmp/*

# Ensure ca-certificates is installed for elasticsearch to use https
RUN apt-get update -qq && apt-get install -yq --no-install-recommends ca-certificates \
  && rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*

# Install exiftool for engine version extraction
RUN apt-get update -qq && apt-get install -yq --no-install-recommends exiftool \
  && rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*

RUN  mkdir -p /opt/malice
COPY update.sh /opt/malice/update

COPY --from=go_builder /bin/avscan /bin/avscan

WORKDIR /malware

ENTRYPOINT ["/bin/avscan"]
CMD ["--help"]

####################################################
####################################################