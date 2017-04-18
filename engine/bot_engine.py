def status(chat_id, vk_request, message):
    return 'Running'

commands = {
    'state': status,
}


def analyze(message, vk_request, chat_id=None):
    print(message)
    if message[0] in commands:
        return commands[message[0]](chat_id, vk_request, status)
    else:
        return 'Command not found'
