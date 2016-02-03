import math

radius = 0.3
theta = 0
print "0.0, 0.0, 0.0,"
while theta < 2*math.pi:
    print "%f, %f, 0.0," % (math.cos(theta)*0.3, math.sin(theta)*0.3)
    theta += 0.1
print "%f, %f, 0.0," % (0.3, 0.0)
print theta/0.1
