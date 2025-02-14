# github.com/carabiner-dev/yamplate

A drop in replacement of `gopkg.in/yaml` with support for templated variables
when unmarshaling and decoding. 

## What is this?

This module emulates the stock 	gopkg.in/yaml.vX yaml packages to decode
YAML data but adds support for variable substitution on the fly. The module wraps
the native yaml modules and has minimal dependencies.

## Templating

Yaml files can embed bash-like variables in curly braces, for example `${VAR1}`
that will be replaced as the file is read from predefined values in a substitution
table. Here is an example of a templated YAML file:

```yaml
---
name: ${USER}
hostname: ${HOSTNAME} 
```

Here's an example parsing the templated example above, defining the required
variables. As you can see, the only difference from stock the yaml modules is
defininf the substitution table in the options:

```golang
package main

import yaml "github.com/carabiner-dev/yamplate"

func main() {
    // Templated YAML:
    reader := strings.NewReader(`---
name: ${USER}
hostname: ${HOSTNAME}
    `)

    // Create a decoder as usual
    dec := yaml.NewDecoder(reader)
    
    // Set the variable substitutions:
    dec.Options.Variables = map[string]string{
        "USER": "John Doe",
        "HOSTNAME": "localhost"
	}

    parsed := struct {
        Name string
        Hostname string
    }{}

    err := dec.Decode(&parsed)
}
```

If a variable is found in the template and no substitution is found in the 
symbols table, both `decoder.Decode` and `Unmarshal()` will return an error.
