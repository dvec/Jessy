import subprocess


def check_ping():
    response = subprocess.Popen(['ping -c 1 api.vk.com'],
                                shell=True, stdout=subprocess.PIPE).stdout.read().splitlines()[1].decode()
    return response[response.rfind('time='):response.find(' ms')]