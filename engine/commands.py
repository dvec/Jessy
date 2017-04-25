import time

from api import api
from . import bot_engine, modes


def delete(**kwargs):
    message = ''.join(kwargs['message'])
    if message.find(' ') != -1:
        return 'To many parameters!'
    with open('data/users', 'r+') as file:
        new_file = file.readlines()
    for i in range(len(new_file)):
        line = new_file[i].split(':')
        if line[0] == message.strip():
            del new_file[i]
            with open('data/users', 'w') as file:
                file.writelines(new_file)
            return 'Yeah, my sir'
    return 'I can\'t find this person on my list!'


def set_mode(**kwargs):
    message = kwargs['message']
    user_id = kwargs['user_id']

    with open('data/users', 'a') as file:
        if message[1] in modes.modes:
            delete(message=message[0], user_id=user_id)
            file.write(message[0] + ':' + message[1] + '\n')
        else:
            return 'I can\'t find this mode on my list!'
    return 'Yeah, my sir'


def get_list(**kwargs):
    del kwargs
    with open('data/users') as file:
        return '\n' + file.read()


def get_mode(**kwargs):
    user_id = kwargs['user_id']
    with open('data/users') as file:
        users = [line.strip().split(':') for line in file.readlines()]
        for user in users:
            if user[0] == str(user_id):
                return user[1] if len(user) >= 2 else '-1'
    return 'Error. User not found'


def set_chat_name(**kwargs):
    title = kwargs['message']
    chat_id = kwargs['chat_id']
    if chat_id is not None:
        if title != '':
            kwargs['vk_request'].messages.editChat(chat_id=chat_id, title=title)
        else:
            return 'Я не вижу названия беседы'
        return ''
    else:
        return 'Вы не в беседе!'


def add_to_database(**kwargs):
    message = kwargs['message']
    data = [j.strip() for j in [i for i in ''.join(message).split('—')] if j.strip() != '']
    if len(data) == 2:
        (message, answer) = data
        message = bot_engine.to_simple_text(message)
    else:
        return 'Недостаточно параметров!'
    del data

    with open('data/answers', 'r+') as file:
        file_lines = file.readlines()
    with open('data/answers', 'r+') as file:

        for line_id in range(len(file_lines)):
            line = file_lines[line_id].strip().split('\\')
            if line[0] == message:
                if answer not in line[1:]:
                    file_lines[line_id] = file_lines[line_id].strip() + '\\' + answer + '\n'
                    file.writelines(file_lines)
                    return 'Вариант сообщения добавлен'
                else:
                    return 'Я уже знаю такую реплику'

        file_lines.append(message + '\\' + answer + '\n')

        file.writelines(file_lines)
        return 'Сообщение добавлено'


def choose_random_user(**kwargs):
    message = bot_engine.to_simple_text(kwargs['message']).split(' ')
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
        user_id = bot_engine.get_random_num(message) % len(users)
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
    with open('data/answers', 'r') as file:
        database_length = len(file.readlines())
        return 'Статус соединения с api.vk.com: ' + smiley[True] + \
               '\nЗаписей в базе данных: ' + str(database_length) + \
               '\nОбработка этого сообщения заняла ' + str(time.time() - start_time)[:5] + ' сек'


def get_random_num(**kwargs):
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
    return ''.join(message) + ' с вероятностью ' + str(bot_engine.get_random_num(message) % 100)


def get_help(**kwargs):
    del kwargs
    return '\n'.join(open('data/help').read().splitlines())


def add_to_chat(**kwargs):
    kwargs['vk_request'].messages.addChatUser(chat_id=1, user_id=kwargs['user_id'])
    return 'Приятного общения!'


def start_game(*args):
    del args
    return 'В разработке'

commands = {
    'normal': {
        'название': set_chat_name,
        'кто': choose_random_user,
        'помощь': get_help,
        'статус': get_state,
        'учись': add_to_database,
        'инфа': get_random_num,
        'беседа': add_to_chat,
        'игра': start_game
    },
    'admin': {
        'del': delete,
        'make': set_mode,
        'list': get_list,
        'gm': get_mode
    }
}
