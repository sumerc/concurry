#!/bin/bash

cmd=$@
pyenv_versions=$(pyenv versions --bare)
echo $cmd
pyenv local $pyenv_versions

for pyenv_version in ${pyenv_versions[*]}; do
    pyenv_version="$(echo $pyenv_version | cut -c-3)"
    new_cmd=${cmd//python/python$pyenv_version};
    echo $new_cmd
done
