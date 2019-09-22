import requests
from random import randint
from time import sleep

while True:
    requests.get("http://192.168.99.100/url/short?url=https://google.com", verify=False)
    sleep(randint(1, 5) / 5.0)
