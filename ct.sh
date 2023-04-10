#!/bin/bash

PIDS=`ps -ef |grep chat |grep -v grep | awk '{print $2}'`
echo $PIDS
if [ "$PIDS" != "" ]; then
echo "chat is runing!"
sudo kill -2 $PIDS
sudo ./chat bg

else
cd /home/yy
sudo ./chat bg

fi