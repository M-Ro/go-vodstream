FROM alpine:3.14.0

RUN apk add --no-cache bash

# copy to image working dir
COPY build/vodstream .
COPY config.yaml .
COPY start_wrapper.sh .

# wrapper to launch services
CMD ["/bin/bash", "./start_wrapper.sh"]