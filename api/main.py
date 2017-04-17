import threading
import time
import vk_requests

from engine import bot_engine
from api import private_data

max_users_per_second = 1

vk_api = vk_requests.create_api(app_id=private_data.app_id, login=private_data.login, password=private_data.password,
                                    scope=['messages', 'account'], timeout=10, access_token=private_data.access_token)
names = ['Jessy', 'Джесси']


def handle_message(words, id, vk_api, chat_id=None):
    return bot_engine.analyze(' '.join(words))


def set_online():
    while threading._main_thread.is_alive():
        vk_api.account.setOnline()
        time.sleep(900)

threading.Thread(target=set_online).start()


def main():
    delay, last = 3, -1
    while True:
        try:
            request = vk_api.messages.get(out=0, count=10, time_offset=(delay * 2))
        except vk_requests.exceptions.Timeout:
            print('ERROR')

        users_block, requests = [], 0
        for i in request['items']:
            if not i['user_id'] in users_block and last != i['id'] and i.get('read_state') == 0:
                print('{user_id}({id}): {message}'.format(user_id=i['user_id'], id=i['id'], message=i['body']))
                try:
                    peer_id = 0
                    if i.get('chat_id') is None:
                        peer_id = i.get('user_id')
                        vk_api.messages.send(user_id=i['user_id'],
                                             message=handle_message(i['body'].split(' '), i['user_id'], vk_api))
                    elif i['body'].split(' ')[0] in names:
                        peer_id = i.get('chat_id')
                        vk_api.messages.send(chat_id=i['chat_id'],
                                             message=handle_message(i['body'].split(' ')[1:],
                                                                    i['user_id'], vk_api, chat_id=i['chat_id']))

                    last = i['id']
                    users_block.append(i['user_id'])
                    requests += 1
                except vk_requests.exceptions.VkAPIError:
                    print('BLOCKED')
                    continue

        time.sleep(delay)
        delay = int(requests / max_users_per_second) + 1


if __name__ == '__main__':
    main()