#!/usr/bin/env bash

package_name=tbmk
output_dir=built

platforms=("darwin/amd64" "darwin/arm64" "freebsd/386" "freebsd/amd64" "freebsd/arm" "freebsd/arm64" "linux/386" "linux/amd64" "linux/arm" "linux/arm64" "netbsd/386" "netbsd/amd64" "netbsd/arm" "netbsd/arm64" "openbsd/386" "openbsd/amd64" "openbsd/arm" "openbsd/arm64" )

for platform in "${platforms[@]}"
do
    platform_split=(${platform//\// })
    GOOS=${platform_split[0]}
    GOARCH=${platform_split[1]}
    output_name=$package_name'-'$GOOS'-'$GOARCH
    output_path="$output_dir/$output_name"

    env CGO_ENABLED=0 GOOS=$GOOS GOARCH=$GOARCH go build -o $output_path .
    if [ $? -ne 0 ]; then
        echo 'An error has occurred! Aborting the script execution...'
        exit 1
    fi
    chmod +x $output_path
    tar -cvzf "$output_path.tar.xz" --transform="s/$output_dir\/$output_name/$package_name/" "$output_path" install uninstall shell config.yaml
    rm $output_path
done
