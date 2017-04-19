# Jessy
Bot for vk. To run it you must install **vk_requests**:
```bash
pip install vk_requests
```
After that is necessary to write **private_data.py** file. It must be located in the api dir This file should look like this:
```python
app_id = 0000000 #id of your app (you can register it at dev.vk.com)
login = '89210000007' #your mail or phone that registered in the vk.com
password = 'xxxxxxxx' #your password from vk.com
access_token = 'xxxxxxxxxxxxxxxxxxxx' #your acess token that you can get in the page of your app
```
Also you can write dictionary for Jessy. Create 2 files in ```engine/data```: **answers** and **ignorance**. Answers is file with answers on question. Don't write punctuation marks in this file! This file should look like this:
```
Hello\Hello!\Hi
How are you?\I'm fine\Normal\Very good
```
**ignorance** file is file that programm will read when can't find answer in the **answers** file. Gust enter in this file reports of ignorance. After that it should look like this:
```
Ok\Nicely\I don't understand you!
```
After that you can run it by enter a command:
```bash
python /path/to/project/Jessy/api/main.py
```