# yamlwalker

Get, set and walk over generic YAML without prior knowlage of the structure of YAML.

Get/Set/Walk methods take dot separated path to the parameter.
The following document fragment

```yaml
object:
  name:
    param: value
another:
  something:
    interresting: thing
```

should be accessed by path as:

```"object.name.param"```
and 
```"another.something.interresting"```

# TODO:

Provide:
- generic methods accepting/returning interface{}
- methods to get the type of returning objects.

Nice to have:
- convinient methods to get values e.g. ToArray(), ToString(), ToInt() and others.
- convinient methods to set values e.g. AssignArray([]slice), AssignString(str string), AssignInt(num int) and others.
