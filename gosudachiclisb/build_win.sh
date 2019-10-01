#!/bin/bash

B="gosudachiclisb"
TYPES="core full"
DIC="system.dic"
T="assets/dict"

function mksb() {
    mkdir -p ${T}

    cp "system_${1}.dic" ${T}/${DIC}

    cd assets
    zip -0r ../zip/"assets_${1}.zip" .

    cd ..
    rm -r assets

    cat "${B}.exe" zip/"assets_${1}.zip" > "${B}${1}.exe"
    zip -A "${B}${1}.exe"
}

GOOS=windows GOARCH=amd64 go build

if [ ! -f "system_core.dic" -o ! -f "system_full.dic" ]; then
    bash ./mksystemdic.sh
fi

[ ! -d "zip" ] && mkdir zip

for t in ${TYPES}
do
    mksb ${t}
done
