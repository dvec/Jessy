from . import modes


def get_random_num(text):
    text = list(' '.join(text).strip())
    num = 50
    for i in range(len(text)):
        num += ord(text[i]) * (i + 1)
    return num


def to_simple_text(text):
    ban_symbols = ['?', '!', '(', ')', '0', '9', ':']
    text = list(''.join(text).strip())
    while len(text) != 0 and text[len(text) - 1] in ban_symbols:
        text = text[:len(text) - 1]
    return ''.join(text).lower()


def analyze(message, vk_request, user_id, chat_id=None):
    user_modes = [mode.split(':') for mode in open('data/users').readlines()]
    for mode in user_modes:
        if int(mode[0]) == user_id:
            return modes.modes.get(mode[1].strip())(message, vk_request, chat_id, user_id)
    return modes.modes.get('normal')(message, vk_request, chat_id)
