import sys

import time

for i in range(2):
    print(i)
    sys.stdout.write(str(i))
    time.sleep(1)

print(sys.version)

if sys.version.startswith('3.5.9'):
    raise Exception('')

# sys.stdout.write('mystdout')
# sys.stderr.write('mystderr')

# raise Exception('aa')
