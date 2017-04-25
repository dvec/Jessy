from threading import Thread, main_thread
import time

import requests.exceptions
import vk_requests

from data import private_data
from engine import bot_engine
from log import log

max_users_per_second = 2

vk_api = vk_requests.create_api(app_id=private_data.app_id, login=private_data.login, password=private_data.password,
                                scope=['messages', 'account', 'friends'], access_token=private_data.access_token)
appeals = ['бот', 'джесси', 'бот,', 'джесси,']
log = log


def set_online():
    while main_thread().is_alive():
        try:
            vk_api.account.setOnline()
            for user_id in vk_api.friends.getRequests()['items']:
                vk_api.friends.add(user_id=user_id)
            log.log('The second thread worked')
            time.sleep(300)
        except requests.exceptions.ReadTimeout as e:
            log.log('TIMEOUT ERROR: ' + str(e))
            continue
        except vk_requests.exceptions.VkAPIError as e:
            log.log('INTERNAL SERVER ERROR: ' + str(e))
            continue
        except requests.exceptions.ConnectionError as e:
            log.log('CONNECTION ABORTED: ' + str(e))
            continue


def handle_message(message, user_name):
    if len(message) >= 1:  # It was 2
        if user_name is not None:
            message = message[0].lower() + message[1:]
            return user_name + ', ' + message
        else:
            return message
    else:
        return message


def main():
    last_id = (-1, -1)
    log.log('Loading is complete')
    second_thread = Thread(target=set_online)
    second_thread.start()
    log.log('Second thread started')
    while True:
        delay = 1
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
                        vk_api.messages.send(user_id=user_id, message=message)
                    elif i['body'].split(' ')[0].lower() in appeals:
                        user_name = vk_api.users.get(user_ids=[user_id])[0].get('first_name')
                        message = handle_message(bot_engine.analyze(message.replace('\n', ' —').split(' ')[1:],
                                                                    vk_api, user_id, chat_id=chat_id), user_name)
                        vk_api.messages.send(chat_id=chat_id, message=message)
                    log.log('Jessy: ' + message)
                except vk_requests.exceptions.VkAPIError as e:
                    print('CAN\'T SEND MESSAGE: ' + str(e))
                except requests.exceptions.ReadTimeout as e:
                    print('TIMEOUT EXCEPTION: ' + str(e))
                    continue
                last_id = i.get('id'), i.get('user_id')

        delay = int(requests_count / max_users_per_second) + 1
        time.sleep(delay)
