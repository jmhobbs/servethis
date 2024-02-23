# Serve This

Read input from `stdin` then open a browser to view it.

## Usage

```
Usage of servethis:
  -file
    	open from filesystem instead of http
  -p int
    	port to serve on (default: a random open port)
  -v	verbose output
```

### Examples

```
$ echo "<b>Hello World</b>" | servethis
```

```
$ mustache data.yaml template.mustache | servethis
```

```
$ pbpaste | servethis
```
