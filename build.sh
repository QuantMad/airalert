#!/bin/bash

ARCH="amd64"
DEB_FOLDER=./airalert_1.0-1_"$ARCH"

# DEBIAN section
mkdir -p $DEB_FOLDER/DEBIAN

cp ./template/control $DEB_FOLDER/DEBIAN/control

echo "Architecture: $ARCH" >>$DEB_FOLDER/DEBIAN/control
echo "Depends: mpg123:$ARCH" >>$DEB_FOLDER/DEBIAN/control

cp ./template/postinst $DEB_FOLDER/DEBIAN/postinst
cp ./template/postrm $DEB_FOLDER/DEBIAN/postrm
# end of DEBIAN section

mkdir -p $DEB_FOLDER/etc/systemd/system
mkdir -p $DEB_FOLDER/usr/local/bin/airalert

env GOOS=linux GOARCH="$ARCH" go build -o $DEB_FOLDER/usr/local/bin/airalert/airalert

cp ./template/postinst $DEB_FOLDER/DEBIAN/postinst && chmod +x $DEB_FOLDER/DEBIAN/postinst
cp ./template/postrm $DEB_FOLDER/DEBIAN/postrm && chmod +x $DEB_FOLDER/DEBIAN/postrm
cp ./template/airalertd.service $DEB_FOLDER/etc/systemd/system/airalertd.service
cp -r ./res $DEB_FOLDER/usr/local/bin/airalert/res

dpkg-deb --build --root-owner-group airalert_1.0-1_"$ARCH"
rm -rf $DEB_FOLDER
mv ./airalert_1.0-1_"$ARCH".deb ./builds/airalert_1.0-1_"$ARCH".deb
