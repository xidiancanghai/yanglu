#!/bin/bash
cd `dirname $0`
pid=`pidof api_server`
cp -f nohup.out nohup.out.1
cat /dev/null > nohup.out
if [ "$pid" == "" ]; then 
    nohup ./api_server &
else 
    kill -SIGUSR2 ${pid}
fi