import random
import time
from engine import api


def set_chat_name(title, vk_request, chat_id):
    if chat_id is not None:
        vk_request.messages.editChat(chat_id=chat_id, title=title)
        return ''
    else:
        return 'Вы не в беседе!'


def choose_random_user(message, vk_request, chat_id):
    if chat_id is not None:
        users = vk_request.messages.getChatUsers(chat_id=chat_id, fields=['nickname'])
        user_id = random.randint(0, len(users) - 1)
        handle = 'это' if not message else message
        return 'Я думаю, что {} {} {}'.format(handle, users[user_id]['first_name'], users[user_id]['last_name'])
    else:
        return 'Вы не в беседе!'


def get_state(message, vk_request, chat_id):
    start_time = time.time()
    print(start_time)
    ping = float(api.check_ping()[5:])
    smiley = {
        ping <= 50: '&#128513;',
        50 < ping <= 70: '&#128512;',
        70 < ping <= 90: '&#128528;',
        90 < ping <= 110: '&#128522;',
        110 < ping <= 130: '&#128551;',
        130 < ping: '&#128565;'
    }
    database_length = len(open('../engine/data/answers', 'r').readlines())
    return 'Статус соединения с api.vk.com: ' + smiley[True] + \
           '\nЗаписей в базе данных: ' + str(database_length) + \
           '\nОбработка этого сообщения заняла ' + str(time.time() - start_time)[:5] + ' сек'


def get_random_num(message, vk_request, chat_id):
    message = message.split(' ')
    seed = 50
    to_replace = {
        'я': 'вы',
        'ты': 'я',
        'что': ''
    }
    for i in range(len(message)):
        text_to_replace = to_replace.get(message[i])
        message[i] = text_to_replace if text_to_replace is not None else message[i]

    message = list(' '.join(message).strip())
    for i in range(len(message)):
        seed += ord(message[i]) * i
    return ''.join(message) + ' с вероятностью ' + str(seed % 100)


def get_help(message, vk_request, chat_id):
    return '\n'.join(open('../engine/data/help').read().splitlines())


commands = {
    'название': set_chat_name,
    'кто': choose_random_user,
    'помощь': get_help,
    'статус': get_state,
    'инфа': get_random_num
}


def to_simple_text(text):
    ban_symbols = ['?', '!', '(', ')', '0', '9']
    text = list(''.join(text))
    while len(text) != 0 and text[len(text) - 1] in ban_symbols:
        text = text[:len(text) - 1]
    return ''.join(text)


def analyze(message, vk_request, chat_id=None):
    command = to_simple_text(message[0]).lower()
    if command in commands:
        return commands[command](' '.join(message[1:]), vk_request, chat_id)
    else:
        with open('../engine/data/answers') as data:
            message = to_simple_text(' '.join(message)).lower()
            for answer in data:
                answer = answer.split('\\')
                if answer[0] == message:
                    return answer[random.randint(1, len(answer) - 1)]

        with open('../engine/data/ignorance') as file:
            data = file.readline().split('\\')
            return data[random.randrange(0, len(data))]
