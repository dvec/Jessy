from threading import Thread, main_thread
import time

import requests.exceptions
import vk_requests

from data import private_data
from engine import bot_engine
from api import api
from log import log

max_users_per_second = 2

vk_api = vk_requests.create_api(app_id=private_data.app_id, login=private_data.login, password=private_data.password,
                                scope=['messages', 'account', 'friends'], access_token=private_data.access_token)
appeals = ['бот', 'джесси', 'бот,', 'джесси,']


def set_online():
    while main_thread().is_alive():
        try:
            vk_api.account.setOnline()
            for user_id in vk_api.friends.getRequests()['items']:
                log.log('New friend: ' + str(user_id))
                vk_api.friends.add(user_id=user_id)
            log.log('The second thread worked')
        except requests.exceptions.ReadTimeout as e:
            log.log('TIMEOUT ERROR: ' + str(e))
            continue
        except vk_requests.exceptions.VkAPIError as e:
            log.log('INTERNAL SERVER ERROR: ' + str(e))
            continue
        except requests.exceptions.ConnectionError as e:
            log.log('CONNECTION ABORTED: ' + str(e))
            continue
        time.sleep(300)


def update_files():
    while main_thread().is_alive():
        try:
            api.update_news()
            log.log('The third thread worked')
        except Exception as e:
            log.log('UPDATING FILES ERROR: ' + str(e))
        time.sleep(10800)


def handle_message(message, user_name):
    if user_name is not None and message.strip() != '':
        message = user_name + ', ' + message[0].lower() + message[1:]
    else:
        message = message[0].upper() + message[1:]
    return message


def get_attachments(message):
    if message.find('<attach>') != message.find('<end>') != -1:
        attachments = message[message.find('<attach>') + len('<attach>'):message.find('<end>')].split('; ')
        return ','.join(attachments), message[:message.find('<attach>')]
    else:
        return -1, message


def send_message(message_api, message, user_id, chat_id):
    attachments, message = get_attachments(message)
    if attachments == -1:
        if chat_id is None:
            message_api.send(user_id=user_id, message=message)
        else:
            message_api.send(chat_id=chat_id, message=message)
    else:
        if chat_id is None:
            message_api.send(user_id=user_id, message=message, attachment=attachments)
        else:
            message_api.send(chat_id=chat_id, message=message, attachment=attachments)


def main():
    # print(get_attachments('<attach>video85635407_165186811_69dff3de4372ae9b6e<end>'))
    last_id = (-1, -1)
    delay = 1
    log.log('Loading is complete')
    second_thread = Thread(target=set_online)
    second_thread.start()
    third_thread = Thread(target=update_files)
    third_thread.start()
    log.log('Second thread started')
    while True:
        users_block = []
        requests_count = 0
        try:
            request = vk_api.messages.get(out=0, time_offset=delay * 2, count=200)
        except requests.exceptions.ReadTimeout as e:
            log.log('TIMEOUT ERROR: ' + str(e))
            continue
        except vk_requests.exceptions.VkAPIError as e:
            log.log('INTERNAL SERVER ERROR: ' + str(e))
            continue
        except requests.exceptions.ConnectionError as e:
            log.log('CONNECTION ABORTED: ' + str(e))
            continue

        for i in request['items']:
            user_id, chat_id, message = i.get('user_id'), i.get('chat_id'), i.get('body')
            if user_id not in users_block \
                    and last_id != (i.get('id'), user_id) and (not i.get('read_state') or chat_id is not None):
                log.log('{user_id}: {message}'.format(user_id=user_id, message=message))
                users_block.append(user_id)
                requests_count += 1
                try:
                    if chat_id is None:
                        message = handle_message(bot_engine.analyze(message.replace('\n', ' —').split(' '),
                                                                    vk_api, user_id), None)
                        message = message[0].upper() + message[1:]
                        send_message(vk_api.messages, message, user_id, None)
                    elif i['body'].split(' ')[0].lower() in appeals:
                        user_name = vk_api.users.get(user_ids=[user_id])[0].get('first_name')
                        message = handle_message(bot_engine.analyze(message.replace('\n', ' —').split(' ')[1:],
                                                                    vk_api, user_id, chat_id=chat_id), user_name)
                        send_message(vk_api.messages, message, None, chat_id)
                    log.log('Jessy: ' + message)
                except vk_requests.exceptions.VkAPIError as e:
                    log.log('CAN\'T SEND MESSAGE: ' + str(e))
                    continue
                except requests.exceptions.ReadTimeout as e:
                    log.log('TIMEOUT EXCEPTION: ' + str(e))
                    continue
                last_id = i.get('id'), i.get('user_id')

        delay = int(requests_count / max_users_per_second) + 1
        time.sleep(delay)
