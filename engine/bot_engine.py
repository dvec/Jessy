import random
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
    ping = float(api.check_ping()[5:])
    print(ping)
    smiley = {
        ping <= 50: '&#128513;',
        50 < ping <= 70: '&#128512;',
        70 < ping <= 90: '&#128528;',
        90 < ping <= 110: '&#128522;',
        110 < ping <= 130: '&#128551;',
        130 < ping: '&#128565;'
    }
    database_length = len(open('../engine/data/answers', 'r').readlines())
    return 'Статус соединения с api.vk.com: ' + smiley[True] + '\nЗаписей в базе данных: ' + str(database_length)


def get_help(message, vk_request, chat_id):
    return '\n'.join(open('../engine/data/help').read().splitlines())


commands = {
    'название': set_chat_name,
    'кто': choose_random_user,
    'помощь': get_help,
    'статус': get_state
}


def to_simple_text(text):
    ban_symbols = ['?', '!', '(', ')', '0', '9']
    text = list(' '.join(text))
    while len(text) != 0 and text[len(text) - 1] in ban_symbols:
        text = text[:len(text) - 1]
    return ''.join(text).split(' ')


def analyze(message, vk_request, chat_id=None):
    message = to_simple_text(message)
    if message[0].lower() in commands:
        return commands[message[0].lower()](' '.join(message[1:]), vk_request, chat_id)
    else:
        with open('../engine/data/answers') as data:
            message = (' '.join(message)).lower()
            for answer in data:
                answer = answer.split('\\')
                if answer[0] == message:
                    return answer[random.randint(1, len(answer) - 1)]

        with open('../engine/data/ignorance') as file:
            data = file.readline().split('\\')
            return data[random.randrange(0, len(data))]
