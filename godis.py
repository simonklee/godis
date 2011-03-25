#!/usr/bin/python2.7

import redis, time

r = redis.Redis()

for i in range(0, 100):
    r.rpush("list", i)

start = time.time()
for i in range(0, 10000):
    r.lrange("list", 0, 50)
print time.time() - start

r.flushdb()
