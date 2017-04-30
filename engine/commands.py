import time
import random

from api import api
from . import functions


def set_chat_name(**kwargs):
    title = kwargs['message']
    chat_id = kwargs['chat_id']
    if chat_id is not None:
        if title != '':
            kwargs['vk_request'].messages.editChat(chat_id=chat_id, title=title)
            return ''
        else:
            return 'Я не вижу названия беседы'
    else:
        return 'Вы не в беседе!'


def add_to_database(**kwargs):
    message = kwargs['message']
    data = [j.strip() for j in [i for i in ''.join(message).split('—')] if j.strip() != '']
    if len(data) == 2:
        (message, answer) = data
        message = functions.to_simple_text(message).lower()
    else:
        return 'Неверное кол-во параметров!'
    del data

    with open('data/bot_data/answers', 'r+') as file:
        file_lines = file.readlines()
    with open('data/bot_data/answers', 'r+') as file:

        for line_id in range(len(file_lines)):
            line = file_lines[line_id].strip().split('\\')
            if line[0] == message:
                if answer not in line[1:]:
                    if line[-1][-1] == '|':
                        return 'Вы не можете научить бота такой реплике'
                    file_lines[line_id] = file_lines[line_id].strip() + '\\' + answer + '\n'
                    file.writelines(file_lines)
                    return 'Вариант сообщения добавлен'
                else:
                    return 'Я уже знаю такую реплику'

        file_lines.append(message + '\\' + answer + '\n')

        file.writelines(file_lines)
        return 'Сообщение добавлено'


def choose_random_user(**kwargs):
    message = functions.to_simple_text(kwargs['message']).split(' ')
    chat_id = kwargs['chat_id']
    to_replace = {
        'я': 'вы',
        'ты': 'я',
        'что': ''
    }
    for i in range(len(message)):
        text_to_replace = to_replace.get(message[i])
        message[i] = text_to_replace if text_to_replace is not None else message[i]
    message = ' '.join(message)
    if chat_id is not None:
        users = kwargs['vk_request'].messages.getChatUsers(chat_id=chat_id, fields=['nickname'])
        user_id = functions.get_random_num(message) % len(users)
        handle = 'это' if not message else message
        return 'Я думаю, что {} {} {}'.format(handle, users[user_id]['first_name'], users[user_id]['last_name'])
    else:
        return 'Вы не в беседе!'


def get_state(**kwargs):
    del kwargs
    start_time = time.time()
    ping = float(api.check_ping()[5:])
    smiley = {
        ping <= 50.0: '&#128513;',
        50.0 < ping <= 70: '&#128512;',
        70.0 < ping <= 90: '&#128528;',
        90.0 < ping <= 110: '&#128522;',
        110.0 < ping <= 130: '&#128551;',
        130.0 < ping: '&#128565;'
    }
    with open('data/bot_data/answers', 'r') as file:
        database_length = len(file.readlines())
        return 'Статус соединения с api.vk.com: ' + smiley[True] + \
               '\nЗаписей в базе данных: ' + str(database_length) + \
               '\nОбработка этого сообщения заняла ' + str(time.time() - start_time)[:5] + ' сек'


def get_news(**kwargs):
    del kwargs
    with open('data/commands_data/news') as file:
        return '\n' + '\n'.join(file.read().split('\\end\\')[:3])


def get_bash(**kwargs):
    del kwargs
    with open('data/commands_data/bash') as file:
        return random.choice(file.read().split('\\end\\'))[:2000]


def get_ithappens(**kwargs):
    del kwargs
    with open('data/commands_data/ithappens') as file:
        return random.choice(file.read().split('\\end\\'))[:2000]


def get_zadolbali(**kwargs):
    del kwargs
    with open('data/commands_data/zadolbali') as file:
        return random.choice(file.read().split('\\end\\'))[:2000]


def get_inf(**kwargs):
    message = kwargs['message'].split(' ')
    to_replace = {
        'я': 'вы',
        'ты': 'я',
        'что': ''
    }
    for i in range(len(message)):
        text_to_replace = to_replace.get(message[i])
        message[i] = text_to_replace if text_to_replace is not None else message[i]

    message = list(' '.join(message).strip())
    return ''.join(message) + ' с вероятностью ' + str(functions.get_random_num(message) % 100) + '%'


def get_help(**kwargs):
    del kwargs
    return '\n' + open('data/commands_data/help').read()


def add_to_chat(**kwargs):
    kwargs['vk_request'].messages.addChatUser(chat_id=1, user_id=kwargs['user_id'])
    return 'Приятного общения!'


def start_game(**kwargs):
    del kwargs
    return 'В разработке'

commands = {
    'normal': {
        'название': set_chat_name,
        'кто': choose_random_user,
        'помощь': get_help,
        'статус': get_state,
        'учись': add_to_database,
        'инфа': get_inf,
        'беседа': add_to_chat,
        'новости': get_news,
        'баш': get_bash,
        'задолбали': get_zadolbali,
        'ithappens': get_ithappens,
        'игра': start_game
    },
    'admin': {
        'del': functions.delete_user,
        'make': functions.set_user_mode
    }
}
