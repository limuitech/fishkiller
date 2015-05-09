#!/usr/bin/env python2
# coding:utf-8

import urllib
import urllib2
import re
import time

REG = "location.href='(down.asp)'"

ADDRESS = 'http://jfh.10086yux.com/submit.asp'

POST = {
        'idType': '1',
        'cnName': u'放烟花咯！'.encode('utf-8'),
        'sec_val': u'你是傻逼吗？'.encode('utf-8'),
        'idcard': '233333333333333',
        'idcard1': '654321',
        'idNo1': '233333333333333333',
        'souji': '13811011011',
        'ssName': '233',
        'sja': '01',
        'sjb': '2020'
        }

data = urllib.urlencode(POST)
req = urllib2.Request(ADDRESS, data)
n = 0
r = re.compile(REG)
while True:
    response = urllib2.urlopen(req)
    try:
        if response.getcode() == 200:
            nk = r.findall(response.read())
            try:
                if nk[0] == 'down.asp':
                    n += 1
                    print u'已轰炸%d次' % n
            except Exception:
                print u'目测被墙'
                break
    except urllib2.HTTPError as e:
        print e.code
    except KeyboardInterrupt:
        print u'爷，玩够了'
        break
    time.sleep(1)
