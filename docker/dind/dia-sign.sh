#!/bin/bash

DOCKERFILE_NAME="dia-scratch-Dockerfile"
POSITIONAL_ARGS=()

while [[ $# -gt 0 ]]; do
  case $1 in
    -h|--help)
      echo '
        USAGE: dia-sign.sh [ARGS] [IMAGE]
          -h, --help      print this help message
          -d, --digest    sha256 digest of image, will used instead [IMAGE] or file .push_stdout.log
        EXAMPLE: dia-sign.sh registry.local/test-app:v1
      '
      exit 0
      ;;
    -d|--digest)
      if ! [[ $2 =~ ^[A-Fa-f0-9]{64}$ ]] ; then echo "Digest is not sha256"; exit 2; fi
      DIGEST_ARG="$2"
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

# Check image build log exists
if [ -f ".push_stdout.log" ];
then
  DIGEST_LOG=$(grep -oP '(?<=digest: sha256:).*(?=.* size)' .push_stdout.log);
fi

# Switch digest source
if ! test -z "$DIGEST_ARG"
then
  IMAGE_DIGEST=$DIGEST_ARG
elif test -z "$DIGEST_LOG"
then
  echo "Try to get digest from docker image inspect"
  IMAGE_DIGEST=$(docker inspect  --format '{{index .RepoDigests 0}}' $1 | rev | cut -d\: -f1 | rev);
else
  IMAGE_DIGEST=$DIGEST_LOG
fi

if test -z "$IMAGE_DIGEST"
then
  echo "Image digest not found"
  exit 2
fi

# Issue image cert
export IMG_DGST_40_SLICE=${IMAGE_DIGEST:0:40}
export VAULT_TOKEN=$(vault write -field=token ${VAULT_AUTH_PATH}/login role=$VAULT_AUTH_ROLE jwt=$CI_JOB_JWT)
export VAULT_PKI_TTL=600 VAULT_PKI_CN=dia-$CI_PROJECT_ID-$IMG_DGST_40_SLICE VAULT_PKI_ALT_NAMES=dia-$CI_PROJECT_ID-$CI_JOB_ID
export CERT=${IMG_DGST_40_SLICE}.pem
consul-template -template=/usr/local/share/cert.ctmpl:$CERT -once

# Place cert as image layer
echo "
FROM scratch
COPY $CERT ." > $DOCKERFILE_NAME

IMAGE_REF=${1%:*}
DIA_SCRATCH_IMAGE=$IMAGE_REF:dia-$IMAGE_DIGEST

# Push the dia image
docker build -t $DIA_SCRATCH_IMAGE -f $DOCKERFILE_NAME . && \
    docker push $DIA_SCRATCH_IMAGE

rm $DOCKERFILE_NAME
