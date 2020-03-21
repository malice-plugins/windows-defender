#!/bin/bash

echo "===> Installing deps..."
apt-get update -qq && apt-get install -yq curl cabextract

echo "===> Download 32-bit antimalware update file.."
curl -L --output /loadlibrary/engine/mpam-fe.exe "https://www.microsoft.com/security/encyclopedia/adlpackages.aspx?arch=x86"
cd /loadlibrary/engine
cabextract mpam-fe.exe

echo "===> Clean up unnecessary files..."
apt-get purge -y --auto-remove curl cabextract "$(apt-mark showauto)"
apt-get clean \
rm -rf /var/lib/apt/lists/* /var/cache/apt/archives /tmp/* /var/tmp/* mpam-fe.exe
