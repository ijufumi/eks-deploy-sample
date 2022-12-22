#!/bin/bash

CURRENT=$(cd $(dirname $0);pwd)
sed -i -e "s|.*.amazonaws.com/[0-9A-Za-z.-]*:[0-9A-Za-z.-]*|${IMAGE_PATH}|g" ${CURRENT}/deployment.yaml
kubectl apply -f ${CURRENT}/deployment.yaml

declare exit_flag=-1
declare count=0
while [ $exit_flag -eq -1 ];
do
    sleep 30
    results=`kubectl get pod --selector=app=app -o jsonpath="{.items[*].spec.containers[*].image}@{..phase}" | tr -s '[[:space:]]' '\n'`
    for result in $results
    do
        echo $result
        images=(${result/@/ })
        if [ "${images[0]}" == "${IMAGE_PATH}" ]; then
            if [ "${images[1]}" == "Running" ]; then
                exit_flag=0
                break
            fi
        fi
    done
    count=$(( $count+1 ))

    if [ $count -eq 60 ]; then
        exit_flag=1
    fi
done

exit $exit_flag