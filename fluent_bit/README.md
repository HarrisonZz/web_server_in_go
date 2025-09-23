Docker command for test:

# Import coustomize fluent-bit.conf
# Mount the path of app.log and flb_app.db
docker run --rm -it -v $(pwd)/fluent-bit.conf:/fluent-bit/etc/fluent-bit.conf -v $(pwd)/.log/:/var/log/ cr.fluentbit.io/fluent/fluent-bit fluent-bit/bin/fluent-bit -c /fluent-bit/etc/fluent-bit.conf