#!/bin/bash

rm -r delivery
rm hw2.zip

go build -o rrt .
for number in {0..4}; do
	./rrt -p $number | python plot.py -png
done
rm rrt

mkdir -p delivery/code
mkdir -p delivery/plots
mkdir -p delivery/bin

mv *.png delivery/plots

cp -r vendor delivery/code/
cp *.go delivery/code/
cp *.py delivery/code/
cp *.json delivery/code/
cp *.txt delivery/code/
# cp *.sh delivery/code/
cp *.md delivery/

GOARCH=amd64 GOOS=linux go build -o delivery/bin/rrt_linux_amd64

zip -r hw2.zip delivery
mv delivery /Users/hdhauk/gdrive/8-semester-V2018/enae788v-motion-planning-for-autonomous-systems/homework/

exit 0
