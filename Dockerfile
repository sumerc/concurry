FROM ubuntu:14.04

RUN apt-get update
RUN apt-get -y install git build-essential \
    python-dev python-pip python3-pip python3-dev libssl-dev zlib1g-dev \
    libbz2-dev libreadline-dev libsqlite3-dev curl vim

RUN git clone https://github.com/yyuu/pyenv.git /root/.pyenv

ENV HOME /root
ENV PYENV_ROOT $HOME/.pyenv
ENV PATH $PYENV_ROOT/shims:$PYENV_ROOT/bin:$PATH

RUN echo 'export PYENV_ROOT=$HOME/.pyenv' >> /root/.profile
RUN echo '$PYENV_ROOT/shims:$PYENV_ROOT/bin:$PATH' >> /root/.profile
RUN echo 'eval "$(pyenv init -)"' >> /root/.profile

RUN pyenv install 2.7.16
# RUN pyenv install 3.0.1
# RUN pyenv install 3.1.5
# RUN pyenv install 3.2.6
# RUN pyenv install 3.3.7
RUN pyenv install 3.4.10
# RUN pyenv install 3.5.9
# RUN pyenv install 3.6.10
# RUN pyenv install 3.7.7
# RUN pyenv install 3.8.2
# RUN pyenv install 3.9-dev

# TODO: https://github.com/s1341/pyenv-alias
