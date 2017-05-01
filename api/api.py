import subprocess
import requests


def check_ping():
    response = subprocess.Popen(['ping -c 1 api.vk.com'],
                                shell=True, stdout=subprocess.PIPE).stdout.read().splitlines()[1].decode()
    return response[response.rfind('time='):response.find(' ms')]


def get_data(url, begin='<![CDATA[', end=']]>'):
    page = requests.get(url).text
    out = []
    while page.find('<![CDATA[') != -1:
        fr = page.find(begin)
        to = page.find(end)
        out.append(page[fr + len(begin):to] + '\\end\\')
        page = page[to + len(end):]
    out = '\n'.join(out)
    replacements = [
        ['</a>', ''],
        ['<p>', ''],
        ['</p>', ''],
        ['<em>', ''],
        ['</em>', ''],
        ['<pre>', ''],
        ['</pre>', ''],
        ['<code>', ''],
        ['</code>', ''],
        ['<br>', '\n'],
        ['&quot;', '"'],
        ['>', ' '],
        ['&lt;', '<'],
        ['&gt;', '>'],
        ['<a href=', '']
    ]
    for replacement in replacements:
        out = out.replace(replacement[0], replacement[1])
    return out


def update_news():
    with open('data/commands_data/news', 'w') as file:
        file.writelines(get_data('https://lenta.ru/rss'))


def update_bash():
    with open('data/commands_data/bash', 'w') as file:
        file.writelines(get_data('http://bash.im/rss/'))


def update_ithappens():
    with open('data/commands_data/ithappens', 'w') as file:
        file.writelines(get_data('http://ithappens.me/rss'))


def update_zadolbali():
    with open('data/commands_data/zadolbali', 'w') as file:
        file.writelines(get_data('http://zadolba.li/rss'))
