#!/bin/bash

PROJECT_NAMESPACE=my-project
SECRET_NAME=forwarder-elasticsearch
SVC_NAME=$(echo "my-project-service-$(cat /dev/urandom | tr -dc 'a-z' | head -c 10)")
FORWARDER_NAME=${SVC_NAME}
UNTIL="21:00"

export KUBECONFIG=./my-kube-config.yaml
WORKDIR=./MyProject/Production/Elasticsearch
cd ${WORKDIR}

if [ -z $UNTIL ]; then
  echo -e "error: please specify until hours"
  exit 1
fi

check_existing_forwarder() {
  if kubectl get svc -n ${PROJECT_NAMESPACE} | grep -q 9200; then
    SVC_NAME=$(kubectl get svc -n ${PROJECT_NAMESPACE} | grep 9200 | awk '{print $1}')
    FORWARDER_NAME=${SVC_NAME}
    # echo -e "error: forwarder service \"${SVC_NAME}\" already exist"
    print_output
    exit 0
  fi
}

generate_config() {
  cp -ap Chart.yaml.tmpl Chart.yaml
  cp -ap values.yaml.tmpl values.yaml

  sed -i "s/FORWARDER_NAME/${FORWARDER_NAME}/g" Chart.yaml
  sed -i "s/FORWARDER_NAME/${FORWARDER_NAME}/g;s/PROJECT_NAMESPACE/${PROJECT_NAMESPACE}/g;s/SECRET_NAME_SVC/${SECRET_NAME}/g" values.yaml
}

execute_deploy() {
  helm upgrade --install -n ${PROJECT_NAMESPACE} ${SVC_NAME} . > /dev/null 2>&1
}

print_output() {
  cat output.tmpl | sed "s/FORWARDER_NAME/${FORWARDER_NAME}/g;s/PROJECT_NAMESPACE/${PROJECT_NAMESPACE}/g"
}

log_service_name() {
  echo -e ${SVC_NAME} $(date +%Y-%m-%d_%H:%M) until ${UNTIL} >> output-svc.txt
}

check_existing_forwarder \
  && generate_config \
  && log_service_name \
  && print_output \
  && execute_deploy
