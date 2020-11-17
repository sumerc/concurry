<p align="center">
    <img src="https://github.com/sumerc/concurry/blob/master/screenshot.png?raw=true" alt="concurry">
</p>

<h1 align="center">concurry</h1>
<p align="center">
    Run your terminal commands in parallel (with some nifty options and colors)
</p>

![version: 1.0](https://img.shields.io/badge/version-1.0-green.svg?style=flat-square)
![language: go](https://img.shields.io/badge/language-go-blue.svg?style=flat-square) 
![license: MIT](https://img.shields.io/badge/license-MIT-blue.svg?style=flat-square) 

## Motivation

I happen to run lots of concurrent scripts in terminal with various requirements. This is especially true
when you maintain a [Python library](https://github.com/sumerc/yappi) that aims to run on all supported Python 
versions. I run unittests on different versions of Python(via pyenv) all the time.

With concurry running concurrent tests locally via pyenv becomes pretty easy:

```bash
pyenv-matrix "python -m unittest discover" | concurry
```

I know there are already lots of good software doing this kind of work already but 
I would especially want to specialize on generating a better/cleaner output with maybe
colors and maybe some UI work in the future via ncurses... Not sure.

And did I mention it colorizes the commands and their outputs?

## Examples:

### Repeat commands (both syncronously and concurrently)

```bash
» echo "sleep 1" | concurry -n=2 -rc=True
2020/04/22 23:05:34 (Task-1) Executing 'sleep 1'
2020/04/22 23:05:34 (Task-2) Executing 'sleep 1'
2020/04/22 23:05:35 (Task-2) 'sleep 1' succeeded. [1.004634221s]
2020/04/22 23:05:35 (Task-1) 'sleep 1' succeeded. [1.005176089s]
2020/04/22 23:05:35 Total elapsed: 1.005404569s
» echo "sleep 1" | concurry -n=2
2020/04/22 23:05:40 (Task-1) Executing 'sleep 1'
2020/04/22 23:05:41 (Task-2) Executing 'sleep 1'
2020/04/22 23:05:42 (Task-1) 'sleep 1' succeeded. [1.001172707s]
2020/04/22 23:05:42 (Task-2) 'sleep 1' succeeded. [1.001278251s]
2020/04/22 23:05:42 Total elapsed: 2.004980207s
```

### Run multiple commands concurrently

Feed command(s) string to concurry via pipe

```bash
» cat long_script | concurry
2020/04/22 23:12:35 (Task-4) Executing 'sleep 1'
2020/04/22 23:12:35 (Task-1) Executing 'sleep 1'
2020/04/22 23:12:35 (Task-2) Executing 'sleep 1'
2020/04/22 23:12:35 (Task-3) Executing 'sleep 1'
2020/04/22 23:12:36 (Task-4) 'sleep 1' succeeded. [1.001739171s]
2020/04/22 23:12:36 (Task-3) 'sleep 1' succeeded. [1.001723235s]
2020/04/22 23:12:36 (Task-1) 'sleep 1' succeeded. [1.001945581s]
2020/04/22 23:12:36 (Task-2) 'sleep 1' succeeded. [1.00211696s]
2020/04/22 23:12:36 Total elapsed: 1.002296034s
```

### Run Python unittests in all available Pyenv versions

`pyenv-matrix` takes a python command string and generates the commands necessary
to initialize and run the given command in all pyenv versions available in the system.
The command string should be in the form of `python <args>`. It should start with 
string `python` that is what `pyenv-matrix` will change.

```bash
» pyenv versions --bare
2.7.17-debug
3.5.9-debug
3.6.10-debug
3.7.7-debug
3.8.2-debug
3.9-dev-debug
» pyenv-matrix "python -m unittest discover"
/home/supo/.pyenv/versions/2.7.17-debug/bin/python -m unittest discover
/home/supo/.pyenv/versions/3.5.9-debug/bin/python -m unittest discover
/home/supo/.pyenv/versions/3.6.10-debug/bin/python -m unittest discover
/home/supo/.pyenv/versions/3.7.7-debug/bin/python -m unittest discover
/home/supo/.pyenv/versions/3.8.2-debug/bin/python -m unittest discover
/home/supo/.pyenv/versions/3.9-dev-debug/bin/python -m unittest discover
```

Feed above output to concurry:

```bash
» pyenv-matrix "python -m unittest discover" | concurry -o=False
2020/04/22 23:23:45 (Task-1) Executing '/home/supo/.pyenv/versions/2.7.17-debug/bin/python -m unittest discover'
2020/04/22 23:23:45 (Task-2) Executing '/home/supo/.pyenv/versions/3.5.9-debug/bin/python -m unittest discover'
2020/04/22 23:23:45 (Task-3) Executing '/home/supo/.pyenv/versions/3.6.10-debug/bin/python -m unittest discover'
2020/04/22 23:23:45 (Task-4) Executing '/home/supo/.pyenv/versions/3.7.7-debug/bin/python -m unittest discover'
2020/04/22 23:23:45 (Task-6) Executing '/home/supo/.pyenv/versions/3.9-dev-debug/bin/python -m unittest discover'
2020/04/22 23:23:45 (Task-5) Executing '/home/supo/.pyenv/versions/3.8.2-debug/bin/python -m unittest discover'
2020/04/22 23:24:08 (Task-1) '/home/supo/.pyenv/versions/2.7.17-debug/bin/python -m unittest discover' succeeded. [16.242247607s]
2020/04/22 23:24:08 (Task-2) '/home/supo/.pyenv/versions/3.5.9-debug/bin/python -m unittest discover' succeeded. [22.282859586s]
2020/04/22 23:24:08 (Task-3) '/home/supo/.pyenv/versions/3.6.10-debug/bin/python -m unittest discover' succeeded. [22.417325148s]
2020/04/22 23:24:08 (Task-4) '/home/supo/.pyenv/versions/3.7.7-debug/bin/python -m unittest discover' succeeded. [22.465179853s]
2020/04/22 23:24:08 (Task-5) '/home/supo/.pyenv/versions/3.8.2-debug/bin/python -m unittest discover' succeeded. [22.498331961s]
2020/04/22 23:24:08 (Task-6) '/home/supo/.pyenv/versions/3.9-dev-debug/bin/python -m unittest discover' succeeded. [22.854892308s]
2020/04/22 23:24:08 Total elapsed: 22.997623419s
```
