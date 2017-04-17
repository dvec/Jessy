# Jessy
Bot for vk. To run it you must install **vk_requests** and **pybrain**:
```bash
sudo pip3 install vk_requests
sudo pip3 install https://github.com/pybrain/pybrain/archive/0.3.3.zip
```
After that is necessary to write **private_data.py** file. This file should look like this:
```python
app_id = 0000000 #id of your app (you can register it at dev.vk.com)
login = '89210000007' #your mail or phone that registered in the vk.com
password = 'xxxxxxxx' #your password from vk.com
access_token = 'xxxxxxxxxxxxxxxxxxxx' #your acess token that you can get in the page of your app
```

After that you can run it by enter a command:
```bash
python /path/to/project/Jessy/api/main.py
```