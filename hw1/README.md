# Homework 1: A*

### Running precompiled binaries
For usage:
```shell
./astar -h
````

Example of running:
```shell
./astar -problem=3 -tree=<tree-file> -path=<path-file> <path-file>
```

### Compiling and running (assuming Ubuntu Linux)
1. Install go (aka. golang), if not already installed.
    ```shell
    snap install --classic go
    ````
    After reloading the shell you should be able to get somehting like:
    ```shell
    $ go version
    go version go1.10 linux/amd64
    ```
2. In the `hw`-folder:
    ```shell
    go build . -o astar
    ````
3. Run as described above.