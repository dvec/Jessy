import subprocess
import requests


def check_ping():
    response = subprocess.Popen(['ping -c 1 api.vk.com'],
                                shell=True, stdout=subprocess.PIPE).stdout.read().splitlines()[1].decode()
    return response[response.rfind('time='):response.find(' ms')]


def update_news():
    page = requests.get('https://lenta.ru/rss').text
    out = []
    for i in page.split('\n'):
        if i.strip()[:7] == '<title>':
            out.append('&#128213;|' + i.replace('<title>', '').replace('</title>', '').strip() + '\n')
    with open('data/commands_data/news', 'w') as file:
        file.writelines(out[2:7])
