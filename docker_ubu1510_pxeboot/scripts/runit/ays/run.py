import subprocess
import signal
import sys


def main():
    stats = subprocess.call(['ays', 'start'])
    if stats != 0:
        sys.exit(stats)

    def sig_handler(s, f):
        print 'Terminating @YS'
        sys.exit(subprocess.call(['ays', 'stop']))

    for s in (signal.SIGTERM, signal.SIGHUP, signal.SIGQUIT):
        signal.signal(s, sig_handler)

    signal.pause()


if __name__ == '__main__':
    main()
