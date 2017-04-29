def get_random_num(text):
    text = list(' '.join(text).strip())
    num = 50
    for i in range(len(text)):
        num += ord(text[i]) * (i + 1)
    return num


def to_simple_text(text, ban_symbols=None):
    if ban_symbols is None:
        ban_symbols = ['?', '!', '(', ')', '0', '9', ':', '.', '|']
    text = list(''.join(text).strip())
    while len(text) != 0 and text[len(text) - 1] in ban_symbols:
        text = text[:len(text) - 1]
    return ''.join(text)


def delete_user(**kwargs):
    message = ''.join(kwargs['message'])
    if ':' not in message:
        with open('data/bot_data/users', 'r+') as file:
            new_file = file.readlines()
        for i in range(len(new_file)):
            line = new_file[i].split(':')
            if line[0] == message.strip():
                del new_file[i]
                with open('data/bot_data/users', 'w') as file:
                    file.writelines(new_file)
                return 'done'
    return 'error'


def set_user_mode(**kwargs):
    message = kwargs['message']

    with open('data/bot_data/users', 'a') as file:
        delete_user(message=message[0])
        file.write(message[0] + ':' + message[1] + '\n')
    return 'done'


def get_user_mode(**kwargs):
    user_id = str(kwargs['message'][0])
    with open('data/bot_data/users') as file:
        users = [line.strip().split(':') for line in file.readlines()]
        for user in users:
            if user[0] == user_id:
                return ':'.join(user[1:])
    return 'normal'
