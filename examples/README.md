# gost Examples

This directory holds examples of how the gost library can be used.

The `sample_component` demonstrates tha bare minimum needed to build a gost "component" and register your own
configuration as well as a trivial HTTP handler.
Build the sample from the root of the repo using the following command:

```
docker build -t gost-sample:latest -f examples/sample_component/Dockerfile .
```

And run it, providing env that is reflected in the component:

```
docker run -e CUSTOM_VALUE=foobar -p 8080:8080 gost-sample:latest
```

And confirm the component is functional with an HTTP request:

```
% curl -i localhost:8080/sample
HTTP/1.1 200 OK

{"CUSTOM_VALUE": "foobar"}
```
