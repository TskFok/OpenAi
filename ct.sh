#!/bin/bash

PIDS=`ps -ef |grep chat |grep -v grep | awk '{print $2}'`
echo $PIDS
if [ "$PIDS" != "" ]; then

echo "chat is runing!"

else

cd /home/yy

./chat bg

fi