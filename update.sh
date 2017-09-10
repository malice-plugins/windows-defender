#!/bin/bash

echo "===> Installing deps..."
apt-get update -qq && apt-get install -yq wget cabextract
echo "===> Download 32-bit antimalware update file.."
wget --progress=bar:force "https://go.microsoft.com/fwlink/?LinkID=121721&arch=x86" -O /loadlibrary/engine/mpam-fe.exe
cd /loadlibrary/engine
cabextract mpam-fe.exe
echo "===> Clean up unnecessary files..."
apt-get purge -y --auto-remove wget cabextract $(apt-mark showauto)
apt-get clean \
rm -rf /var/lib/apt/lists/* /var/cache/apt/archives /tmp/* /var/tmp/* mpam-fe.exe