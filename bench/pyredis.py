#!/usr/bin/python2

import timeit
t = timeit.Timer('r.get("0")', setup='import redis; r = redis.Redis(); gc.enable()')
n = 20000
rep = 10
res = t.repeat(repeat=rep, number=n)
tot = 0.0

for r in res:
    print "%.2f op/sec, real %.4f sec" % (float(n/r), r)
    tot += r

print "avg %.2f op/sec, real %.4f sec" % (float(n/(tot/rep)), tot/rep)
