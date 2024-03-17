#!/bin/bash

ssh $DEPLOY_USER@$DEPLOY_HOST sudo systemctl stop $SERVICE
scp ./build/trackrouter-linux-x86_64 $DEPLOY_USER@$DEPLOY_HOST:$PATH
ssh $DEPLOY_USER@$DEPLOY_HOST sudo systemctl start $SERVICE
