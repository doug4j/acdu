## acdu generate process project

Creates or updates Process Runtime Bundle and Connector. [NOT IMPLEMENTED]

### Synopsis

Creates or updates Process Runtime Bundle and Connector.

```
acdu generate process project [flags]
```

### Options

```
  -b, --bundlename string       MANDATORY: Name of the runtime bundle (friendly for kubernetes and jars).
  -c, --channel string          MANDATORY: Name of implementation (starting lower case alpha and all alphanum).
  -d, --destdir string          optional: Destination directory for writing the runtime bundle template. This directory will be appended with the BundleName. Example: a destdir of '/Users/john/projects' and a bundlename 'my-bundle' will results in the runtime bundle being created in a final directory '/Users/john/projects/my-bundle' (default "./")
  -t, --downloadtag string      optional: Tag name to pull the zip file github. (default "7.0.0.Beta3")
  -h, --help                    help for project
  -i, --implementation string   MANDATORY: Name of implementation (starting lower case alpha and all alphanum).
  -p, --packagename string      MANDATORY: Name of java package (friendly for java packages).
  -a, --projectname string      MANDATORY: The default Project to use for the process bundle.
```

### Options inherited from parent commands

```
      --config string   config file (default is $HOME/.acdu.yaml)
  -v, --verbose         optional verbosity setting (default:'false'), if true temporary side effects (for instance temp directories/files are not cleaned up) for debugging
```

### SEE ALSO

* [acdu generate process](acdu_generate_process.md)	 - Creates a process runtime bundle for development.

###### Auto generated by spf13/cobra on 5-Dec-2018