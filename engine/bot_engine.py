from .modes import modes
from . import advanced


def analyze(message, vk_request, user_id, chat_id=None):
    user_mode = advanced.get_user_mode(message=[user_id])
    return modes.get(user_mode)(message, vk_request, chat_id, user_id)
