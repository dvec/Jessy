import sys


def format_to_color(text):
    return text.format(normal='\033[0m', red='\033[31m', green='\033[96m', blue='\033[34m')

invitation = format_to_color('{green}> {normal}')


def halt():
    sys.exit(0)


def get_log():
    pass


def init():
    from sys import path
    from os import getcwd
    path.append(getcwd())


def emulate():
    from engine import bot_engine
    while True:
        print(bot_engine.analyze(input(invitation).lower().split(' '), None, None))


def run():
    from api import vk_api
    vk_api.main()


def stable_run():
    try:
        run()
    except Exception:
        from time import sleep
        sleep(5)
        stable_run()


def main():
    init()
    cases = ((halt, 'Exit from this menu'),
             (get_log, 'Get log of Jessy'),
             (emulate, 'Emulate Jessy'),
             (stable_run, 'Run Jessy with automatically restarting'),
             (run, 'Normal run without restarting'))

    if len(sys.argv) == 1:
        print(format_to_color('{blue}Welcome! Choose case to run it{normal}: '))
        for case in range(len(cases)):
            description = cases[case][1]
            print(format_to_color('{red}' + str(case + 1) + '{normal}. {blue}' + description + '{normal}'))

        choose = int(input(invitation)) - 1
    else:
        choose = int(sys.argv[1]) - 1

    exit(cases[choose][0]())


if __name__ == '__main__':
    main()
