import threading
import time
import vk_requests
import requests

from engine import bot_engine
from api import private_data

max_users_per_second = 1

vk_api = vk_requests.create_api(app_id=private_data.app_id, login=private_data.login, password=private_data.password,
                                scope=['messages', 'account', 'users'], timeout=10, access_token=private_data.access_token)
appeals = ['jessy', 'джесси', 'jessy,', 'джесси,']


def set_online():
    while threading._main_thread.is_alive():
        vk_api.account.setOnline()
        time.sleep(300)

threading.Thread(target=set_online).start()


def handle_message(message, user_name):
    if len(message) >= 2:
        if user_name is not None:
            message = message[0].lower() + message[1:]
            return user_name + ', ' + message
        else:
            return message
    else:
        return message


def main():
    delay, last = 3, -1
    while True:
        try:
            request = vk_api.messages.get(out=1, time_offset=(delay * 2))
        except requests.exceptions.ReadTimeout:
            print('TIMEOUT ERROR')
            continue
        except vk_requests.exceptions.VkAPIError:
            print('INTERNAL SERVER ERROR')
            continue

        users_block, requests_count = [], 0
        for i in request['items']:
            if not i['user_id'] in users_block and last != i['id'] and (not i.get('read_state') or i.get('chat_id') is not None):
                print('{user_id}({id}): {message}'.format(user_id=i['user_id'], id=i['id'], message=i['body']))
                try:
                    message = ''
                    if i.get('chat_id') is None:
                        message = handle_message(bot_engine.analyze(i['body'].split(' '), vk_api), None)
                        vk_api.messages.send(user_id=i['user_id'], message=message)
                    elif i['body'].split(' ')[0].lower() in appeals:
                        message = handle_message(bot_engine.analyze(i['body'].split(' ')[1:], vk_api, chat_id=i['chat_id']),
                                                                    vk_api.users.get(user_ids=[i['user_id']])[0]['first_name'])
                        vk_api.messages.send(chat_id=i['chat_id'], message=message)
                    print('Jessy: ' + message)
                    last = i['id']
                    users_block.append(i['user_id'])
                    requests_count += 1
                except vk_requests.exceptions.VkAPIError as e:
                    print('CAN\'T SEND MESSAGE: ' + e.message)
                    continue
                except requests.exceptions.ReadTimeout as e:
                    print('TIMEOUT EXCEPTION: ' + e.message)
                    continue

        time.sleep(delay)
        delay = int(requests_count / max_users_per_second) + 1


if __name__ == '__main__':
    main()
