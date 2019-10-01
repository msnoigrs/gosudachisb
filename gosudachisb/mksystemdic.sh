#!/bin/bash

if [ -d gosudachi ]; then
    cd gosudachi
    git pull origin master
else
    git clone https://github.com/msnoigrs/gosudachi.git
    cd gosudachi
fi

bash scripts/build.sh

if [ -d SudachiDict ]; then
    cd SudachiDict
    git pull origin develop
else
    git clone https://github.com/WorksApplications/SudachiDict.git
    cd SudachiDict
    git lfs pull
fi

cd ../dist
bash ../scripts/mksystemdic.sh ../SudachiDict

mv *.dic ../../
