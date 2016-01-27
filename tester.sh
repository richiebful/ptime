$GOPATH/bin/ptime
echo "SHOULD BE FAIL 1"
$GOPATH/bin/ptime -tz -15
echo "SHOULD BE FAIL 2"
$GOPATH/bin/ptime -lat 40
echo "SHOULD FAIL 3"
$GOPATH/bin/ptime -lat -360
echo "SHOULD FAIL 4"
$GOPATH/bin/ptime -long -80
echo "SHOULD FAIL 5"
$GOPATH/bin/ptime -long 100
echo "SHOULD FAIL 6"
$GOPATH/bin/ptime -long 4959
echo "SHOULD FAIL 7"
