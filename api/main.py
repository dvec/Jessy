import threading
import time
import vk_requests
import requests

from engine import bot_engine
from api import private_data

max_users_per_second = 5

vk_api = vk_requests.create_api(app_id=private_data.app_id, login=private_data.login, password=private_data.password,
                                scope=['messages', 'account', 'friends'], access_token=private_data.access_token)
appeals = ['jessy', 'джесси', 'jessy,', 'джесси,']


def set_online():
    while threading._main_thread.is_alive():
        try:
            vk_api.account.setOnline()
            for user_id in vk_api.friends.getRequests()['items']:
                vk_api.friends.add(user_id=user_id)
            time.sleep(300)
        except requests.exceptions.ReadTimeout:
            print('TIMEOUT ERROR')
            continue
        except vk_requests.exceptions.VkAPIError:
            print('INTERNAL SERVER ERROR')
            continue

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
    while True:
        delay = 1
        try:
            request = vk_api.messages.get(out=0, time_offset=delay, count=200)
        except requests.exceptions.ReadTimeout:
            print('TIMEOUT ERROR')
            continue
        except vk_requests.exceptions.VkAPIError:
            print('INTERNAL SERVER ERROR')
            continue

        users_block, requests_count, last_id = [], 0, -1
        for i in request['items']:
            user_id, chat_id, message = i.get('user_id'), i.get('chat_id'), i.get('body')
            if user_id not in users_block and last_id != i.get('random_id') and (not i.get('read_state') or chat_id is not None):
                print('({time}){user_id}: {message}'.format(time=time.strftime('%X'), user_id=user_id, message=message))
                users_block.append(user_id)
                requests_count += 1
                try:
                    if chat_id is None:
                        message = handle_message(bot_engine.analyze(message.split(' '), vk_api), None)
                        vk_api.messages.send(user_id=i['user_id'], message=message)
                    elif i['body'].split(' ')[0].lower() in appeals:
                        user_name = vk_api.users.get(user_ids=[user_id])[0].get('first_name')
                        message = bot_engine.analyze(message.split(' ')[1:], vk_api, chat_id=chat_id)
                        vk_api.messages.send(chat_id=chat_id, message=handle_message(message, user_name))
                    print('Jessy: ' + message)
                except vk_requests.exceptions.VkAPIError as e:
                    print('CAN\'T SEND MESSAGE: ' + e.message)
                    continue
                except requests.exceptions.ReadTimeout as e:
                    print('TIMEOUT EXCEPTION: ' + e.message)
                    continue
                last_id = i.get('random_id')

        delay = int(requests_count / max_users_per_second) + 1
        time.sleep(delay)


if __name__ == '__main__':
    main()
