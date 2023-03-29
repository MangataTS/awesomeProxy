# 导入相关模块
import requests
import threading
from bs4 import BeautifulSoup

# 定义一个函数，用于获取网页内容并解析链接
def get_links(url):
    header = {
       'User-Agent':'Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/70.0.3538.25 Safari/537.36 Core/1.70.3775.400 QQBrowser/10.6.4208.400'
    }
    for i in range(100):
        # 发送请求，获取响应
        response = requests.get(url,headers=header)
        page_text = response.text

        print(page_text)
# 定义一个列表，存储要爬取的网址
urls = ["http://127.0.0.1:9090/hello"]

# 定义一个空列表，存储线程对象
threads = []

# 遍历每个网址，创建一个线程对象，并将其添加到线程列表中
for url in urls:
    for i in range(1):
        thread = threading.Thread(target=get_links, args=(url,))
        threads.append(thread)

# 遍历每个线程对象，启动线程
for thread in threads:
    thread.start()

# 遍历每个线程对象，等待线程结束
for thread in threads:
    thread.join()