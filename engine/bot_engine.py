import random


def set_chat_name(title, vk_request, chat_id):
    if chat_id is not None:
        vk_request.messages.editChat(chat_id=chat_id, title=title)
        return ''
    else:
        return 'Вы не в беседе!'

commands = {
    'name': set_chat_name,
    'название': set_chat_name
}


def to_simple_text(text):
    ban_symbols = ['?', '!', '(', ')', '0', '9']
    text = list(text)
    if len(text) != 0:
        while text[len(text) - 1] in ban_symbols:
            text = text[:len(text) - 1]
        return ''.join(text)


def analyze(message, vk_request, chat_id=None):
    if message[0] in commands:
        return commands[message[0]](' '.join(message[1:]), vk_request, chat_id)
    else:
        with open('../engine/data/answers', 'r') as data:
            message_text = to_simple_text(' '.join(message))
            for answer in data:
                answer = answer.split('\\')
                if answer[0].lower() == message_text:
                    return answer[random.randint(1, len(answer) - 1)]

        with open('../engine/data/ignorance', 'r') as file:
            data = file.readline().split('\\')
            return data[random.randrange(0, len(data))]
