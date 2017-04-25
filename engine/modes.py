import random
from . import commands
from . import bot_engine


def blocked(message, vk_request, chat_id, *args):
    return 'Вы не можете отправлять мне сообщения! Вы в черном списке.'


def normal(message, vk_request, chat_id, *args):
    command = bot_engine.to_simple_text(''.join(message[:1])).lower()
    if command in commands.commands['normal']:
        return commands.commands['normal'][command](' '.join(message[1:]), vk_request, chat_id)
    else:
        with open('data/answers') as data:
            message = bot_engine.to_simple_text(' '.join(message)).lower()
            for answer in data:
                answer = answer.split('\\')
                if answer[0] == message:
                    return answer[random.randint(1, len(answer) - 1)]

        with open('data/ignorance') as file:
            data = file.readline().split('\\')
            return random.choice(data)


def admin(message, vk_request, chat_id, user_id):
    if message[:1] != ['sudo']:
        return normal(message, vk_request, chat_id)
    else:
        if len(message) >= 2:
            out = commands.commands['admin'].get(message[1])
            return out(message[2:], user_id) if out is not None else 'I can\'t find this command on my list'
        else:
            return 'To few parameters'

modes = {
    'admin': admin,
    'blocked': blocked,
    'normal': normal,
}
