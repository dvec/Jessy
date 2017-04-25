def format_to_color(text):
    return text.format(normal='\033[0m', red='\033[31m', green='\033[96m', blue='\033[34m')

invitation = format_to_color('{green}> {normal}')


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


def main():
    init()
    cases = ((get_log, 'Get log of Jessy'), (emulate, 'Emulate Jessy'), (run, 'Run Jessy'))

    print(format_to_color('{blue}Welcome! Choose case to run it{normal}: '))
    for case in range(len(cases)):
        description = cases[case][1]
        print(format_to_color('{red}' + str(case + 1) + '{normal}. {blue}' + description + '{normal}'))

    choose = int(input(invitation)) - 1
    exit(cases[choose][0]())


if __name__ == '__main__':
    main()
