import random
from .commands import commands
from .advanced import to_simple_text, answers_cache


def blocked(*args):
    del args
    return 'Вы не можете отправлять мне сообщения! Вы в черном списке.'


def normal(words, vk_request, chat_id, user_id):
    command = to_simple_text(words[0] if len(words) != 0 else '').lower()
    if command in commands['normal']:
        return commands['normal'][command](message=' '.join(words[1:]),
                                           vk_request=vk_request, chat_id=chat_id, user_id=user_id)
    else:
        def find_answer(message_words, word_len=100):
            message = to_simple_text(' '.join(message_words)).lower()
            for line in answers_cache:
                index = line[0].rfind('*')
                if line[0][:word_len] == message or \
                        (index == len(line[0]) - 1 and line[0][:index] == message[:index]):
                    out = to_simple_text(random.choice(line[1:]), ban_symbols=['|', '*'])
                    if out[0] == '$' and out[1:] != line[0]:
                        return normal(out[1:].split(' '), vk_request, chat_id, user_id)
                    return out
        answer = find_answer(words)
        if answer is not None:
            return answer
        elif len(words) != 1:
            answer = find_answer([words[0]], len(words[0]))
            if answer is not None:
                return answer
        return random.choice(answers_cache[0][1:])


def admin(message, vk_request, chat_id, user_id):
    if message[:1] != ['sudo']:
        return normal(message, vk_request, chat_id, user_id)
    else:
        if len(message) >= 2:
            out = commands['admin'].get(message[1])
            if out is None:
                return 'I can\'t find this command on my list'
            return out(message=message[1], user_id=user_id)
        else:
            return 'Too few parameters'


modes = {
    'admin': admin,
    'blocked': blocked,
    'normal': normal,
}
