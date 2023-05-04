#!/bin/bash

VERSION="$1"
DIR="./release"

for FILE in ${DIR}/*
do
    if [ -e "${FILE}" ]; then
        NAME=$(basename "${FILE}")
        mkdir -p ${NAME}
        cp config.yaml ${NAME}
        cp release/${NAME} ${NAME}/codepass
        tar zcf ${NAME}_${VERSION}.tar.gz ${NAME}
        rm -rf ${NAME}
    fi
done