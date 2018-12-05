// Copyright © 2018 Doug Johnson <doug4j@gmail.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"log"
	"regexp"
)

func main2() {
	regex := regexp.MustCompile(`(?ms)(.*spring.application.name=)[\$|\{|\}|_|A-Z|:|a-z|\.|\-|0-9]*(.*)`)
	match := regex.FindStringSubmatch(theStr)
	log.Printf("%v%v%v", match[1], "7.0.0-beta3", match[2])
}

var theStr = `spring.application.name=${ACT_RB_APP_NAME:rb-my-app}
activiti.cloud.application.name=default-app

spring.cloud.stream.bindings.auditProducer.destination=${ACT_RB_AUDIT_PRODUCER_DEST:engineEvents}
spring.cloud.stream.bindings.auditProducer.contentType=${ACT_RB_AUDIT_PRODUCER_CONTENT_TYPE:application/json}
spring.cloud.stream.bindings.myCmdResults.destination=${ACT_RB_COMMAND_RESULTS_DEST:commandResults}
spring.cloud.stream.bindings.myCmdResults.group=${ACT_RB_COMMAND_RESULTS_GROUP:myCmdGroup}
spring.cloud.stream.bindings.myCmdResults.contentType=${ACT_RB_COMMAND_RESULTS_CONTENT_TYPE:application/json}
spring.cloud.stream.bindings.myCmdProducer.destination=${ACT_RB_COMMAND_RESULTS_DEST:commandConsumer}
spring.cloud.stream.bindings.myCmdProducer.contentType=${ACT_RB_COMMAND_RESULTS_CONTENT_TYPE:application/json}
spring.cloud.stream.bindings.signalProducer.destination=${ACT_RB_SIGNAL_PRODUCER_DEST:signalEvent}
spring.cloud.stream.bindings.signalProducer.contentType=${ACT_RB_SIGNAL_PRODUCER_CONTENT_TYPE:application/json}
spring.cloud.stream.bindings.signalConsumer.destination=${ACT_RB_SIGNAL_CONSUMER_DEST:signalEvent}
spring.cloud.stream.bindings.signalConsumer.group=${ACT_RB_SIGNAL_CONSUMER_GROUP:mySignalConsumerGroup}
spring.cloud.stream.bindings.signalConsumer.contentType=${ACT_RB_SIGNAL_CONSUMER_CONTENT_TYPE:application/json}
spring.jackson.serialization.fail-on-unwrapped-type-identifiers=${ACT_RB_JACKSON_FAIL_ON_UNWRAPPED_IDS:false}

keycloak.auth-server-url=${ACT_KEYCLOAK_URL:http://activiti-keycloak:8180/auth}
keycloak.realm=${ACT_KEYCLOAK_REALM:activiti}
keycloak.resource=${ACT_KEYCLOAK_RESOURCE:activiti}
keycloak.ssl-required=${ACT_KEYCLOAK_SSL_REQUIRED:none}
keycloak.public-client=${ACT_KEYCLOAK_CLIENT:true}

keycloak.security-constraints[0].authRoles[0]=${ACT_KEYCLOAK_USER_ROLE:ACTIVITI_USER}
keycloak.security-constraints[0].securityCollections[0].patterns[0]=${ACT_KEYCLOAK_PATTERNS:/v1/*}
keycloak.security-constraints[1].authRoles[0]=${ACT_KEYCLOAK_ADMIN_ROLE:ACTIVITI_ADMIN}
keycloak.security-constraints[1].securityCollections[0].patterns[0]=/admin/*

keycloak.principal-attribute=${ACT_KEYCLOAK_PRINCIPAL_ATTRIBUTE:preferred-username}
# see https://issues.jboss.org/browse/KEYCLOAK-810 for configuration options

activiti.keycloak.admin-client-app=${ACT_KEYCLOAK_CLIENT_APP:admin-cli}
activiti.keycloak.client-user=${ACT_KEYCLOAK_CLIENT_USER:client}
activiti.keycloak.client-password=${ACT_KEYCLOAK_CLIENT_PASSWORD:client}
# this user needs to have the realm management roles assigned
spring.rabbitmq.host=${ACT_RABBITMQ_HOST:rabbitmq}

spring.activiti.useStrongUuids=true`
