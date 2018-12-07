# Activiti Cloud Development Utilities

See [Command Line Doc](./doc/acdu.md) for more details.

## Pre-reqs
* [Maven](https://maven.apache.org/download.cgi) Installed
* [Java](https://www.oracle.com/technetwork/java/javase/downloads/jdk8-downloads-2133151.html) 8 Installed
* [Docker for Desktop](https://www.docker.com/products/docker-desktop) with Kubernetes enabled and >=10 gigs of ram and >=4 cores
* [Helm](https://docs.helm.sh/using_helm/#installing-helm) installed
* [Postman](https://www.getpostman.com/apps) installed
* Install acdu into your path:
- [MacOS](bin/darwin_amd64/acdu)
- [Linux 32 bit](bin/linux_386/acdu)
- [Linux arm](bin/linux_arm/acdu)
- [Linux 64 bit](bin/linux_arm/acdu)
- [Windows 32 bit](bin/windows_386/acdu.exe)
- [Windows 64 bit](bin/windows_amd64/acdu.exe)

## Pre-req Config 

### Import Postman Collection and Environment

Import [this Activiti 7 Postman Collection](postman/Activiti-v7-REST-API.postman_collection.json) into Postman.

Import [this Activiti 7 Postman Enviornment](postman/activiti-local.postman_environment.json) environment into Postman.

- [ ] Update pull request https://github.com/Activiti/activiti-cloud-examples/pull/95 with the latest from above.

### Create a local project directory

mkdir ```${projectDir}```

### Install Docker Dashboard

``
kubectl create -f https://raw.githubusercontent.com/kubernetes/dashboard/master/src/deploy/recommended/kubernetes-dashboard.yaml
``

```
kubectl get pods --namespace=kube-system
NAME                                         READY     STATUS              RESTARTS   AGE
etcd-docker-for-desktop                      1/1       Running             0          1h
kube-apiserver-docker-for-desktop            1/1       Running             0          1h
kube-controller-manager-docker-for-desktop   1/1       Running             0          1h
kube-dns-86f4d74b45-4jw6f                    3/3       Running             0          1h
kube-proxy-zwwpx                             1/1       Running             0          1h
kube-scheduler-docker-for-desktop            1/1       Running             0          1h
kubernetes-dashboard-7b9c7bc8c9-d7dcz        0/1       ContainerCreating   0          39s
```

Wait for the Dashboard POD to load and you see STATUS as Running. It could take some time to change from ContainerCreating to Running like so:

```
kubectl get pods --namespace=kube-system
NAME                                         READY     STATUS    RESTARTS   AGE
etcd-docker-for-desktop                      1/1       Running   0          1h
kube-apiserver-docker-for-desktop            1/1       Running   0          1h
kube-controller-manager-docker-for-desktop   1/1       Running   0          1h
kube-dns-86f4d74b45-lw4bg                    3/3       Running   0          1h
kube-proxy-sm5ps                             1/1       Running   0          1h
kube-scheduler-docker-for-desktop            1/1       Running   0          1h
kubernetes-dashboard-7b9c7bc8c9-5rr4g        1/1       Running   0          13s
tiller-deploy-6fd8d857bc-kjsxz               1/1       Running   0          30m
```

Create a file called in the root of ${projectDir} called k8s-dashboard-nodeport-service.yaml with the following content:

```
apiVersion: v1
kind: Service
metadata:
  labels:
    k8s-app: kubernetes-dashboard
  name: kubernetes-dashboard-nodeport
  namespace: kube-system
spec:
  ports:
  - port: 8443
    protocol: TCP
    targetPort: 8443
    nodePort: 31234
  selector:
    k8s-app: kubernetes-dashboard
  sessionAffinity: None
  type: NodePort
```

CD to ${projectDir} and run the following command:

```
kubectl create -f k8s-dashboard-nodeport-service.yaml
service "kubernetes-dashboard-nodeport" created
```

Open open up a browser at https://localhost:31234

You should be able to by-pass security warnings and navigate to the kubernetes dashboard.

### Setup Activiti Cloud Helm Repo

Run:

```
helm repo add activiti-cloud-charts https://activiti.github.io/activiti-cloud-charts/
```

## Install Activiti Postman Templates 


## Install Activiti 7 Infrastructure

Find you IP ```${ipAddress}``` such as 192.168.7.24

192.168.7.24

```
acdu install infra -i ${ipAddress}.nip.io -n activiti7
```
Note it is normal to wait for ~20 minutes after issuing the above commmand.

```
acdu install infra -i 192.168.7.24.nip.io -n activiti7
2018/12/07 11:31:35 ⏳ [WORKING]    Deploy project via helm (and verify) [helm install stable/nginx-ingress --namespace activiti7 --timeout 720 --wait]...
2018/12/07 11:31:39 👍 [OK]         Code compiled and packaged  [helm install stable/nginx-ingress --namespace activiti7 --timeout 720 --wait]
2018/12/07 11:31:39 ⏱ [TIME]       Deploy project via helm (and verify) elapsed time: 3.93s
2018/12/07 11:31:39 ⏳ [WORKING]    Waiting until all pods are ready, expiring at [2018-12-07 11:43:39.190995 -0500 EST m=+723.970970062], pods being checked: [youngling-ragdoll-nginx-ingress-controller-79644f7b56-7wsbj youngling-ragdoll-nginx-ingress-default-backend-78d6bd5998d5zr7]
2018/12/07 11:31:59 👍 [OK]         All Containers in Namespace [activiti7] ready. Pods [youngling-ragdoll-nginx-ingress-controller-79644f7b56-7wsbj youngling-ragdoll-nginx-ingress-default-backend-78d6bd5998d5zr7] with 2 containers in total
2018/12/07 11:31:59 ⏱ [TIME]       Install Ingress elapsed time: 24.001s
2018/12/07 11:31:59 ⏱ [TIME]       Elapsed time thus far: 24.001s
2018/12/07 11:31:59 ⏳ [WORKING]    Deploy project via helm (and verify) [helm install activiti-cloud-charts/activiti-cloud-full-example --namespace activiti7 --timeout 720 --wait --set global.keycloak.url=http://activiti-keycloak.192.168.7.24.nip.io/auth --set global.gateway.host=activiti-cloud-gateway.192.168.7.24.nip.io --set infrastructure.activiti-keycloak.keycloak.keycloak.ingress.hosts[0]=activiti-keycloak.192.168.7.24.nip.io --set infrastructure.activiti-cloud-gateway.ingress.hostName=activiti-cloud-gateway.192.168.7.24.nip.io --set application.activiti-cloud-connector.enabled=false --set application.runtime-bundle.enabled=false --set activiti-cloud-modeling.enabled=true]...
2018/12/07 11:33:43 👍 [OK]         Code compiled and packaged  [helm install activiti-cloud-charts/activiti-cloud-full-example --namespace activiti7 --timeout 720 --wait --set global.keycloak.url=http://activiti-keycloak.192.168.7.24.nip.io/auth --set global.gateway.host=activiti-cloud-gateway.192.168.7.24.nip.io --set infrastructure.activiti-keycloak.keycloak.keycloak.ingress.hosts[0]=activiti-keycloak.192.168.7.24.nip.io --set infrastructure.activiti-cloud-gateway.ingress.hostName=activiti-cloud-gateway.192.168.7.24.nip.io --set application.activiti-cloud-connector.enabled=false --set application.runtime-bundle.enabled=false --set activiti-cloud-modeling.enabled=true]
2018/12/07 11:33:43 ⏱ [TIME]       Deploy project via helm (and verify) elapsed time: 1m44.214s
2018/12/07 11:33:43 ⏳ [WORKING]    Waiting until all pods are ready, expiring at [2018-12-07 11:45:43.451963 -0500 EST m=+848.235824737], pods being checked: [fallacious-turkey-activiti-cloud-audit-94d87cf79-pm4nf fallacious-turkey-activiti-cloud-gateway-79bffd74d-87mlz fallacious-turkey-activiti-cloud-modeling-76cfc9684-45wh2 fallacious-turkey-activiti-cloud-query-5df745f5b4-nr5rs fallacious-turkey-ke-0 fallacious-turkey-rabbitmq-0 youngling-ragdoll-nginx-ingress-controller-79644f7b56-7wsbj youngling-ragdoll-nginx-ingress-default-backend-78d6bd5998d5zr7]
2018/12/07 11:36:40 👍 [OK]         All Containers in Namespace [activiti7] ready. Pods [fallacious-turkey-activiti-cloud-audit-94d87cf79-pm4nf fallacious-turkey-activiti-cloud-gateway-79bffd74d-87mlz fallacious-turkey-activiti-cloud-modeling-76cfc9684-45wh2 fallacious-turkey-activiti-cloud-query-5df745f5b4-nr5rs fallacious-turkey-ke-0 fallacious-turkey-rabbitmq-0 youngling-ragdoll-nginx-ingress-controller-79644f7b56-7wsbj youngling-ragdoll-nginx-ingress-default-backend-78d6bd5998d5zr7] with 9 containers in total
2018/12/07 11:36:40 ⏱ [TIME]       Install Activiti Full Example elapsed time: 4m40.969s
2018/12/07 11:36:40 ⏱ [TIME]       Elapsed time thus far: 5m4.97s
2018/12/07 11:36:40 ℹ️ [INFO]       identityURL url is available at
http://activiti-keycloak.192.168.7.24.nip.io/auth/admin/master/console
default user/name: admin/admin
2018/12/07 11:36:40 ℹ️ [INFO]       modelerURL url is available at
http://activiti-cloud-gateway.192.168.7.24.nip.io/activiti-cloud-modeling
default user/name: modeler/password
2018/12/07 11:36:40 ℹ️ [INFO]       modelingSwaggerURL url is available at
http://activiti-cloud-gateway.192.168.7.24.nip.io/activiti-cloud-modeling-backend/swagger-ui.html
2018/12/07 11:36:40 ⏱ [TIME]       Total Elapsed time: 5m4.97s
```

Notice the Helm deployment prefix for the Activiti cloud infrastructure ```famous-sasquatch-activiti-cloud-audit-5d44f5c69c-j79ff``` as in the pattern ```${activitiInfraDeployName}-activiti-cloud-audit-...``` 

Per the output above, login to the modeling app, ensuring it is up-and-running via ```http://activiti-cloud-gateway.192.168.7.24.nip.io/activiti-cloud-modeling``` or ```http://activiti-cloud-gateway.${ipAddress}.nip.io/activiti-cloud-modeling```

## Generate Runtime Bundle Code

From ```${projectDir}``` run 

```
acdu generate process bundle -b my-rb-1 -p com.example -a project1
```

The output looks like:
```
2018/12/05 19:44:11 ⏳ TagName '7.0.0.Beta3' requested for download (using downloader=stereotype-github-quickstart)
2018/12/05 19:44:11 👍 Created temp directory ./.acdu-tmp-prj-project1-pkg-com.example-bun-my-rb-1-tag-7.0.0.Beta3/
2018/12/05 19:44:11 👍 Read 23626 bytes from https://github.com/Activiti/activiti-cloud-runtime-bundle-quickstart/archive/master.zip
2018/12/05 19:44:11 👍 Wrote temp 23626 bytes of runtime bundle template to ./.acdu-tmp-prj-project1-pkg-com.example-bun-my-rb-1-tag-7.0.0.Beta3/Tmp.zip
2018/12/05 19:44:11 👍 Unzipped runtime bundle template to ./.acdu-tmp-prj-project1-pkg-com.example-bun-my-rb-1-tag-7.0.0.Beta3/.acdu-tmp-prj-project1-pkg-com.example-bun-my-rb-1-tag-7.0.0.Beta3/activiti-cloud-runtime-bundle-quickstart-master
2018/12/05 19:44:11 👍 Removed temp zip file ./.acdu-tmp-prj-project1-pkg-com.example-bun-my-rb-1-tag-7.0.0.Beta3/Tmp.zip
2018/12/05 19:44:11 👍 RT Bundle Transform 1 of 7: Rule rename template directory. Renamed from .acdu-tmp-prj-project1-pkg-com.example-bun-my-rb-1-tag-7.0.0.Beta3/activiti-cloud-runtime-bundle-quickstart-master ./.acdu-tmp-prj-project1-pkg-com.example-bun-my-rb-1-tag-7.0.0.Beta3/my-rb-1/
2018/12/05 19:44:11 👍 RT Bundle Transform 2 of 7: Rule adjust group id, artifact id, REPLACE_ME_APP_NAME, and activiti-cloud-dependencies  in pom.xml. Adjusted with com.example, my-rb-1, project1, and 7.0.0.Beta3 writing 3587 bytes at ./.acdu-tmp-prj-project1-pkg-com.example-bun-my-rb-1-tag-7.0.0.Beta3/my-rb-1/pom.xml
2018/12/05 19:44:11 👍 RT Bundle Transform 3 of 7: Rule change the spring.application.name and activiti.cloud.application.name. Adjusted to my-rb-1 and project1 with 2611 bytes at ./.acdu-tmp-prj-project1-pkg-com.example-bun-my-rb-1-tag-7.0.0.Beta3/my-rb-1/src/main/resources/application.properties
2018/12/05 19:44:11 👍 RT Bundle Transform 4 of 7: Rule rename the charts folder. Adjusted to ./.acdu-tmp-prj-project1-pkg-com.example-bun-my-rb-1-tag-7.0.0.Beta3/my-rb-1/charts/my-rb-1
2018/12/05 19:44:11 👍 RT Bundle Transform 5 of 7: Rule change helm service name and repository to bundle name and tag to 'latest'. Adjusted to my-rb-1 with 1418 bytes at ./.acdu-tmp-prj-project1-pkg-com.example-bun-my-rb-1-tag-7.0.0.Beta3/my-rb-1/charts/my-rb-1/values.yaml
2018/12/05 19:44:11 👍 RT Bundle Transform 6 of 7: Rule change helm name. Adjusted to my-rb-1 with 186 bytes at ./.acdu-tmp-prj-project1-pkg-com.example-bun-my-rb-1-tag-7.0.0.Beta3/my-rb-1/charts/my-rb-1/Chart.yaml
2018/12/05 19:44:11 👍 RT Bundle Transform 7 of 7: Rule move the folder to the desired destination. Adjusted to ./my-rb-1
2018/12/05 19:44:11 👍 Removed temp directory from ./.acdu-tmp-prj-project1-pkg-com.example-bun-my-rb-1-tag-7.0.0.Beta3/
2018/12/05 19:44:11 ℹ️ Ready to use the Activiti Cloud Runtime Bundle
```


## Generate Cloud Connector Code

From ```${projectDir}``` run 

```
acdu generate process connector -b my-con-1 -p com.example -a project1 -c myChannel1 -i MyFirstConn
```

The output looks like:
```
acdu generate process connector -b my-con-1 -p com.example -a project1 -c myChannel1 -i MyFirstConn
2018/12/05 19:46:56 ⏳ TagName '7.0.0.Beta3' requested for download
2018/12/05 19:46:56 👍 Created temp directory ./.acdu-tmp-prj-project1-pkg-com.example-bun-my-con-1-tag-7.0.0.Beta3/
2018/12/05 19:46:57 👍 Read 24370 bytes from https://github.com/Activiti/activiti-cloud-connector-quickstart/archive/master.zip
2018/12/05 19:46:57 👍 Wrote temp 24370 bytes of runtime bundle template to ./.acdu-tmp-prj-project1-pkg-com.example-bun-my-con-1-tag-7.0.0.Beta3/Tmp.zip
2018/12/05 19:46:57 👍 Unzipped runtime bundle template to ./.acdu-tmp-prj-project1-pkg-com.example-bun-my-con-1-tag-7.0.0.Beta3/.acdu-tmp-prj-project1-pkg-com.example-bun-my-con-1-tag-7.0.0.Beta3/activiti-cloud-connector-quickstart-master
2018/12/05 19:46:57 👍 Removed temp zip file ./.acdu-tmp-prj-project1-pkg-com.example-bun-my-con-1-tag-7.0.0.Beta3/Tmp.zip
2018/12/05 19:46:57 👍 RT Bundle Transform 1 of 9: Rule rename template directory. Renamed from .acdu-tmp-prj-project1-pkg-com.example-bun-my-con-1-tag-7.0.0.Beta3/activiti-cloud-connector-quickstart-master ./.acdu-tmp-prj-project1-pkg-com.example-bun-my-con-1-tag-7.0.0.Beta3/my-con-1/
2018/12/05 19:46:57 👍 RT Bundle Transform 2 of 9: Rule adjust group id, artifact id, REPLACE_ME_APP_NAME, and activiti-cloud-dependencies  in pom.xml. Adjusted with com.example, my-con-1, project1, and 7.0.0.Beta3 writing 3357 bytes at ./.acdu-tmp-prj-project1-pkg-com.example-bun-my-con-1-tag-7.0.0.Beta3/my-con-1/pom.xml
2018/12/05 19:46:57 👍 RT Bundle Transform 3 of 9: Rule change the spring.application.name and activiti.cloud.application.name. Adjusted to my-con-1 and project1 with 1027 bytes at ./.acdu-tmp-prj-project1-pkg-com.example-bun-my-con-1-tag-7.0.0.Beta3/my-con-1/src/main/resources/application.properties
2018/12/05 19:46:57 👍 RT Bundle Transform 4 of 9: Rule rename the charts folder. Adjusted to ./.acdu-tmp-prj-project1-pkg-com.example-bun-my-con-1-tag-7.0.0.Beta3/my-con-1/charts/my-con-1
2018/12/05 19:46:57 👍 RT Bundle Transform 5 of 9: Rule change helm service name and repository to bundle name and tag to 'latest'. Adjusted to my-con-1 with 1420 bytes at ./.acdu-tmp-prj-project1-pkg-com.example-bun-my-con-1-tag-7.0.0.Beta3/my-con-1/charts/my-con-1/values.yaml
2018/12/05 19:46:57 👍 RT Bundle Transform 6 of 9: Rule change helm name. Adjusted to my-con-1 with 187 bytes at ./.acdu-tmp-prj-project1-pkg-com.example-bun-my-con-1-tag-7.0.0.Beta3/my-con-1/charts/my-con-1/Chart.yaml
2018/12/05 19:46:57 👍 RT Bundle Transform 7 of 9: Rule move the folder to the desired destination. Adjusted to ./my-con-1
2018/12/05 19:46:57 👍 Removed temp directory from ./.acdu-tmp-prj-project1-pkg-com.example-bun-my-con-1-tag-7.0.0.Beta3/
2018/12/05 19:46:57 👍 Transform 8 of 9: Rule adjust chanel in ExampleConnectorChannels.java. Adjusted as myChannel1 with 965 bytes at ./my-con-1/src/main/java/org/activiti/cloud/connector/impl/ExampleConnectorChannels.java
2018/12/05 19:46:57 👍 Transform 8 of 9: Rule application.properties file for channel name and . Adjusted as myChannel1 with 957 bytes at ./my-con-1/src/main/resources/application.properties
2018/12/05 19:46:57 ℹ️ Ready to use the Activiti Cloud Connector
```

## Create simple Hello World Process and Export to Process Runtime Bundle

In the modeling app:
1. Create an Application (call it 'TestApp')
2. Click on the newly created app
3. Create a process in the App (call it 'TestProcess')
4. In the newly created process, create an Activity after the Start Event (call it 'Automation Task')
5. With the 'Automation Task' highlighted, click the gears icon and choose 'Service Task'
6. In the properties area under Implementation type 'MyFirstConn' (note: this is the same name used in the '-i' or implementation flag in the creation of the cloud connector above)
7. Close the process with an End Event.
8. Save the process
9. Export the process as a file.
10. Move the exported file to your process runtime bundle as in ```${projectDir}/my-rb-1/src/main/resources/processes/TestProcess.bpmn20.xml```
11. Near the top after 'Dashboard' click 'TestApp'
12. On the far right near the version click the download symbol.
13. Move the saved app zip file to ```${projectDir}```. (This is simply saving the app for when you shutdown the Kubernetes cluster. Currently, the Modeling app doesn't use a persistent connection.)  

## Complile, Package, and Deploy Runtime Bundle

Look up the rabbitmq instance from navigating in the Kubernetes Dashboard to the **famous-sasquatch**-rabbitmq

From ```${projectDir}```/my-rb-1 run
```
acdu install process quickstart -k http://${ipAddress}.nip.io/auth -i ${ipAddress}.nip.io -m ${activitiInfraDeployName}-rabbitmq -n activiti7 
```

as in

```
acdu install process quickstart -k http://192.168.7.24.nip.io/auth -i 192.168.7.24.nip.io -m famous-sasquatch-rabbitmq -n activiti7
2018/12/05 20:08:33 ℹ️ Using source directory [./]
2018/12/05 20:08:33 ⏳ Compile and package code [mvn package]...
2018/12/05 20:12:21 👍 Code compiled and packaged  [mvn package]
2018/12/05 20:12:21 ⏱ Compile and package code elapsed time: 3m47.862s
2018/12/05 20:12:21 ⏳ Build docker image into local registry [docker build -t my-rb-1 .]...
2018/12/05 20:12:39 👍 Code compiled and packaged  [docker build -t my-rb-1 .]
2018/12/05 20:12:39 ⏱ Build docker image into local registry elapsed time: 18.042s
2018/12/05 20:12:39 ⏳ Update helm dependencies [helm dep update ./charts/my-rb-1]...
2018/12/05 20:12:50 👍 Code compiled and packaged  [helm dep update ./charts/my-rb-1]
2018/12/05 20:12:50 ⏱ Update helm dependencies elapsed time: 10.961s
2018/12/05 20:12:50 ⏳ Deploy project via helm (and verify) [helm install ./charts/my-rb-1 --namespace activiti7 --timeout 720 --wait --set global.rabbitmq.host.value=famous-sasquatch-rabbitmq --set global.keycloak.url=http://192.168.7.24.nip.io/auth]...
2018/12/05 20:12:52 👍 Code compiled and packaged  [helm install ./charts/my-rb-1 --namespace activiti7 --timeout 720 --wait --set global.rabbitmq.host.value=famous-sasquatch-rabbitmq --set global.keycloak.url=http://192.168.7.24.nip.io/auth]
2018/12/05 20:12:52 ⏱ Deploy project via helm (and verify) elapsed time: 2.55s
2018/12/05 20:12:52 ⏳ Waiting until all pods are ready, expiring at [2018-12-05 20:24:52.653274 -0500 EST m=+979.597314914], pods being checked: [exhaling-alligator-my-rb-1-759dbdd6dc-rzkbq famous-sasquatch-activiti-cloud-audit-5d44f5c69c-j79ff famous-sasquatch-activiti-cloud-gateway-b45f9d4b5-46tj8 famous-sasquatch-activiti-cloud-modeling-84df585675-dt8rw famous-sasquatch-activiti-cloud-query-7d9d465f5b-9jcbh famous-sasquatch-key-0 famous-sasquatch-rabbitmq-0 killjoy-bear-nginx-ingress-controller-59c986f84-f7cnn killjoy-bear-nginx-ingress-default-backend-95dd5fc44-njtn4]
2018/12/05 20:13:39 👍 All Containers in Namespace [activiti7] ready. Pods [exhaling-alligator-my-rb-1-759dbdd6dc-rzkbq famous-sasquatch-activiti-cloud-audit-5d44f5c69c-j79ff famous-sasquatch-activiti-cloud-gateway-b45f9d4b5-46tj8 famous-sasquatch-activiti-cloud-modeling-84df585675-dt8rw famous-sasquatch-activiti-cloud-query-7d9d465f5b-9jcbh famous-sasquatch-key-0 famous-sasquatch-rabbitmq-0 killjoy-bear-nginx-ingress-controller-59c986f84-f7cnn killjoy-bear-nginx-ingress-default-backend-95dd5fc44-njtn4] with 10 containers in total
2018/12/05 20:13:39 ℹ️ swaggerURL (if available) url is available at
http://activiti-cloud-gateway..192.168.7.24.nip.io/project1-my-rb-1/swagger-ui.html
2018/12/05 20:13:39 ⏱ Total Elapsed time: 5m6.167s
```

Note: this may take a long time the first time.

In a browser nativate to the swaggerURL as noted in your output.

## Verify the Deployed process in Postman

1. Open Postman.

2. Ensure the 'activiti-local' environment is selected.

3. Edit 'activiti-local' in the [Postman Enviornment](images/postman-environment-config.jpg) to match yours.

In quick succession, do the following:

4. Under the 'Activiti v7 REST API' collection, navigate to 'keycloak > getKeycloakToken'

5. Click ```send```. You should 

6. Under the 'Activiti v7 REST API' collection, navigate to 'rb-my-app > getProcessDefinitions'

7. Click ```send```.

8. 

## Deploy Cloud Connector

Look up the rabbitmq instance from navigating in the Kubernetes Dashboard to the famous-sasquatch-rabbitmq

From ```${projectDir}```/my-rb-1 run
```
acdu install process quickstart -k http://${ipAddress}.nip.io/auth -i ${ipAddress}.nip.io -m ${activitiInfraDeployName}-rabbitmq -n activiti7 
```

as in

```
acdu install process quickstart -k http://192.168.7.24.nip.io/auth -i 192.168.7.24.nip.io -m famous-sasquatch-rabbitmq -n activiti7 
```
