import time

file_id = int(time.time())
file = open(str('log/logs/' + str(file_id) + '.log'), 'a+')


def log(message):
    message = '[{}]: {}'.format(time.strftime('%X'), message).replace('\n', ' ')
    print(message if len(message) <= 100 else message[:100] + '...')
    file.write(message + '\n')
    file.flush()
