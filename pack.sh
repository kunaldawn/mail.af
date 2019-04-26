#!/usr/bin/env bash

cd ./web/dashboard/
ng build --prod
cd ..
packr2 --verbose
