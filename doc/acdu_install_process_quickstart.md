## acdu install process quickstart

Installs a Process Runtime Bundle or Process Connector Quick Start.

### Synopsis

Installs a Process Runtime Bundle or Process Connector Quick Start.

```
acdu install process quickstart [flags]
```

### Options

```
  -h, --help                                help for quickstart
  -k, --identityHost string                 MANDATORY: Hostname of the identity service connection (Keycloak).
  -i, --ingressIP string                    MANDATORY: Kubernetes ingress IP address. Tip: for a local install, when connected to the internet this can suffixed with '.nip.io' to map external ips to internal ones.
  -m, --mqhost string                       MANDATORY: Hostname of the message and queuing connection (RabbitMQ).
  -n, --namespace string                    MANDATORY: Kubernetes namespace to install into.
  -q, --queryForAllPodsRunningSeconds int   optional: Number of seconds to wait until querying to see if all pods are running. (default 2)
  -s, --sourceDir string                    optional: The directory where the source code exists. (default "./")
  -t, --timeoutSeconds int                  optional: Number of seconds to wait until the kubernetes commands give up. (default 720)
```

### Options inherited from parent commands

```
      --config string   config file (default is $HOME/.acdu.yaml)
  -v, --verbose         optional verbosity setting (default:'false'), if true temporary side effects (for instance temp directories/files are not cleaned up) for debugging
```

### SEE ALSO

* [acdu install process](acdu_install_process.md)	 - Installs Process related objects. [NOT IMPLEMENTED]

###### Auto generated by spf13/cobra on 5-Dec-2018