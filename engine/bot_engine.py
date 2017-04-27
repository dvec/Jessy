from .modes import modes
from . import functions


def analyze(message, vk_request, user_id, chat_id=None):
    user_mode = functions.get_user_mode(message=[user_id])
    return modes.get(user_mode)(message, vk_request, chat_id, user_id)
