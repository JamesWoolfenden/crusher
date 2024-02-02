FROM alpine

RUN apk --no-cache add build-base git curl jq bash
RUN curl -s -k https://api.github.com/repos/JamesWoolfenden/crusher/releases/latest | jq '.assets[] | select(.name | contains("linux_386")) | select(.content_type | contains("gzip")) | .browser_download_url' -r | awk '{print "curl -L -k " $0 " -o ./crusher.tar.gz"}' | sh
RUN tar -xf ./crusher.tar.gz -C /usr/bin/ && rm ./crusher.tar.gz && chmod +x /usr/bin/crusher && echo 'alias crusher="/usr/bin/crusher"' >> ~/.bashrc
COPY entrypoint.sh /entrypoint.sh

# Code file to execute when the docker container starts up (`entrypoint.sh`)
ENTRYPOINT ["/entrypoint.sh"]
