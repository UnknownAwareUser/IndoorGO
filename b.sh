CC=arm-linux-gnueabi-gcc GOOS=linux GOARCH=arm GOARM=7 CGO_ENABLED=1 go build --ldflags="-extldflags=-static" --tags sqlite_omit_load_extension -o m_arm7
