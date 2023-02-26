#!/bin/sh

DOCKERFILE_NAME="dia-scratch-Dockerfile"
POSITIONAL_ARGS=()

while [[ $# -gt 0 ]]; do
  case $1 in
    -h|--help)
      echo '
        USAGE: dia-sign.sh [ARGS] [IMAGE]
          -h, --help      print this help message
          -c, --cert      path to x509 certificate file (base64 encoded)
        EXAMPLE: dia-sign.sh -c /tmp/image.crt registry.local/test-app:v1
      '
      exit 0
      ;;
    -c|--cert)
      CERT="$2"
      shift # past argument
      shift # past value
      ;;
    -*|--*)
      echo "Unknown option $1"
      exit 1
      ;;
    *)
      POSITIONAL_ARGS+=("$1") # save positional arg
      shift # past argument
      ;;
  esac
done

set -- "${POSITIONAL_ARGS[@]}"

echo "
FROM scratch
COPY $CERT ." > $DOCKERFILE_NAME

IMAGE_DIGEST=$(docker inspect  --format '{{index .RepoDigests 0}}' $1 | rev | cut -d\: -f1 | rev)
if test -z "$IMAGE_DIGEST"
then
  echo "Image digest not found"
  exit 1
fi

IMAGE_REF=${1%:*}
DIA_SCRATCH_IMAGE=$IMAGE_REF:dia-$IMAGE_DIGEST

docker build -t $DIA_SCRATCH_IMAGE -f $DOCKERFILE_NAME . && \
    docker push $DIA_SCRATCH_IMAGE

rm $DOCKERFILE_NAME
