#!/bin/bash

DIR="./release"

mkdir -p codepass

for FILE in ${DIR}/*
do
    if [ -e "${FILE}" ]; then
        NAME=$(basename "${FILE}")
        cp config.yaml codepass
        mv release/${NAME} codepass/codepass
        tar zcf ${NAME}.tar.gz codepass
    fi
done

rm -rf codepass