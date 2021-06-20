cd `dirname $0`
pid=`pidof api_server`
cp -f ./app/nohup.out ./app/nohup.out.1
cat /dev/null > ./app/nohup.out
if [ "$pid" == "" ]; then
    nohup ./api_server </dev/null > ./app/nohup.out 2>&1 &
else
    kill -SIGUSR2 ${pid}
fi