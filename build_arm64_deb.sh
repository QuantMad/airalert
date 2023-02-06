#!bin/bash

cp ./postinst ./airalert_1.0-1_arm64/DEBIAN/postinst
chmod +x ./airalert_1.0-1_arm64/DEBIAN/postinst
cp ./postrm ./airalert_1.0-1_arm64/DEBIAN/postrm
chmod +x ./airalert_1.0-1_arm64/DEBIAN/postrm

mkdir -p ./airalert_1.0-1_arm64/etc/systemd/system
cp ./airalertd.service ./airalert_1.0-1_arm64/etc/systemd/system/airalertd.service

mkdir -p ./airalert_1.0-1_arm64/usr/local/bin/airalert/
cp -r ./res ./airalert_1.0-1_arm64/usr/local/bin/airalert/res

env GOOS=linux GOARCH=arm64 go build -o ./airalert_1.0-1_arm64/usr/local/bin/airalert/airalert
dpkg-deb --build --root-owner-group airalert_1.0-1_arm64
mv ./airalert_1.0-1_arm64.deb ./builds/airalert_1.0-1_arm64.deb