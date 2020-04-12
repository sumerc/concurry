# concurry

## Examples:

### Run commands concurrently (separated by semicolon)
```
»  echo 'sleep 1 | echo "1" ; sleep 2 | echo "2"' | concurry
2020/04/12 20:52:11 Executing ' zsh -c sleep 2 | echo "2" '
2020/04/12 20:52:11 Executing ' zsh -c sleep 1 | echo "1" '
2020/04/12 20:52:12 Command sleep 1 | echo "1" succeeded.
2020/04/12 20:52:13 Command sleep 2 | echo "2" succeeded.
```

### Run Python unittests in all available Pyenv versions

pyenv-matrix.sh takes a python command string and generates the commands necessary
to initialize and run the given command in all pyenv versions available in the system.

```bash
»  pyenv versions --bare
2.7.17
3.5.9
3.6.10
3.9.0a4
»  ./pyenv-matrix.sh "python run_tests.py"
python2.7 run_tests.py;
python3.5 run_tests.py;
python3.6 run_tests.py;
python3.9 run_tests.py;
```

Feed above output to concurry:

```bash
»  pyenv-matrix.sh "python run_tests.py" | concurry -v
2020/04/12 21:13:47 Executing ' zsh -c python3.9 run_tests.py '
2020/04/12 21:13:47 Executing ' zsh -c python3.5 run_tests.py '
2020/04/12 21:13:47 Executing ' zsh -c python2.7 run_tests.py '
2020/04/12 21:13:47 Executing ' zsh -c python3.6 run_tests.py '
2020/04/12 21:13:47 Command python2.7 run_tests.py succeeded.
2020/04/12 21:13:47 Command python3.9 run_tests.py succeeded.
2020/04/12 21:13:47 Command python3.5 run_tests.py succeeded.
2020/04/12 21:13:47 Command python3.6 run_tests.py succeeded.
```
