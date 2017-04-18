import threading
import time
import vk_requests

from engine import bot_engine
from api import private_data

max_users_per_second = 1

vk_api = vk_requests.create_api(app_id=private_data.app_id, login=private_data.login, password=private_data.password,
                                    scope=['messages', 'account'], timeout=10, access_token=private_data.access_token)
appeals = ['Jessy', 'Джесси', 'Jessy, ', 'Джесси, ']


def set_online():
    while threading._main_thread.is_alive():
        vk_api.account.setOnline()
        time.sleep(900)

threading.Thread(target=set_online).start()


def main():
    delay, last = 3, -1
    while True:
        try:
            request = vk_api.messages.get(out=1, count=10, time_offset=(delay * 2))
        except vk_requests.exceptions.Timeout as e:
            print('ERROR: ' + e.message)
            continue

        users_block, requests = [], 0
        for i in request['items']:
            if not i['user_id'] in users_block and last != i['id'] and (i.get('read_state') == 0 or i.get('chat_id') is not None):
                print('{user_id}({id}): {message}'.format(user_id=i['user_id'], id=i['id'], message=i['body']))
                try:
                    message = ''
                    if i.get('chat_id') is None:
                        message = bot_engine.analyze(i['body'].lower().split(' '), vk_api)
                        vk_api.messages.send(user_id=i['user_id'], message=message)
                    elif i['body'].split(' ')[0] in appeals:
                        message = bot_engine.analyze(i['body'].lower().split(' ')[1:], vk_api, chat_id=i['chat_id'])
                        vk_api.messages.send(chat_id=i['chat_id'], message=message)
                    print('Jessy: ' + message)
                    last = i['id']
                    users_block.append(i['user_id'])
                    requests += 1
                except vk_requests.exceptions.VkAPIError as e:
                    print('CAN\'T SEND MESSAGE: ' + e.message)
                    continue

        time.sleep(delay)
        delay = int(requests / max_users_per_second) + 1


if __name__ == '__main__':
    main()
