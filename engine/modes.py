import random
from .commands import commands
from .functions import to_simple_text


def blocked(*args):
    del args
    return 'Вы не можете отправлять мне сообщения! Вы в черном списке.'


def normal(message, vk_request, chat_id, user_id):
    command = to_simple_text(''.join(message[:1])).lower()
    if command in commands['normal']:
        return commands['normal'][command](message=' '.join(message[1:]),
                                           vk_request=vk_request, chat_id=chat_id, user_id=user_id)
    else:
        with open('data/bot_data/answers') as data:
            message = to_simple_text(' '.join(message)).lower()
            for answer in data:
                answer = answer.split('\\')
                if answer[0] == message:
                    return answer[random.randint(1, len(answer) - 1)]

        with open('data/bot_data/ignorance') as file:
            data = file.readline().split('\\')
            return random.choice(data)


def admin(message, vk_request, chat_id, user_id):
    if message[:1] != ['sudo']:
        return normal(message, vk_request, chat_id, user_id)
    else:
        if len(message) >= 2:
            out = commands['admin'].get(message[1])
            if out is None:
                return 'I can\'t find this command on my list'
            return out(message=message[2:], user_id=user_id)
        else:
            return 'To few parameters'

modes = {
    'admin': admin,
    'blocked': blocked,
    'normal': normal,
}