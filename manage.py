import sys


def format_to_color(text):
    return text.format(normal='\033[0m', red='\033[31m', green='\033[96m', blue='\033[34m')

invitation = format_to_color('{green}> {normal}')


def halt():
    sys.exit(0)


def clear_log():
    import os
    for file in os.listdir('log/logs'):
        os.remove('log/logs/' + file)


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
        print('Fatal error. Restarting')
        from time import sleep
        sleep(5)
        stable_run()


def spam():
    targets = [24261502, 59221166, 39130136, 47484197, 24261502, 59221166, 39130136,
               47484197, 47122228, 41843539, 42590964, 61413825, 77208800, 68114314,
               24261502, 59221166, 39130136, 47484197, 47122228, 41843539, 42590964,
               24261502, 59221166, 39130136, 47484197, 24261502, 59221166, 24261502]
    import vk_requests
    import time
    from data import private_data
    vk_api = vk_requests.create_api(app_id=private_data.app_id, login=private_data.login,
                                    password=private_data.password,
                                    scope=['wall'], access_token=private_data.access_token)
    message = 'Я чат-бот. Добавь в меня в друзья (я приму заявку через несколько минут) и поговори со мной. :)'
    attachments = 'photo424752907_456239023'
    for target in targets:
        print('Message to', target, end=' ')
        try:
            vk_api.wall.post(owner_id=-target, message=message, attachments=attachments)
            print('ok')
            time.sleep(1)
        except:
            print('err')
            time.sleep(1)


def main():
    init()
    cases = ((halt, 'Exit from this menu'),
             (clear_log, 'Clear log of Jessy'),
             (emulate, 'Emulate Jessy'),
             (spam, 'Make spam'),
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
