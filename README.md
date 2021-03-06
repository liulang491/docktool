# docktool

```
NAME:
   docktool

USAGE:
   docktool [global options] COMMAND [command options] [arguments ...]

VERSION:
   v1.2.x

AUTHORS:
   应卓 <yingzhor@gmail.com>

COMMANDS:
   filegen   generate file using template
   filedel   delete files/dirs using wildcard
   wait      wait until tcp reachable or timeout
   sleep     make current thread sleep
   secret    encode/decode a string
   uuid      create random uuid
   test      test env/file/dir/tcp

GLOBALS OPTIONS:
   -q, --quiet     quiet mode (default: false)
   -h, --help      print this usage
   -v, --version   print version information

SEE ALSO:
   https://github.com/yingzhuo/docktool

Run 'docktool COMMAND --help' for more information on a command.

```

### Installing on docker image

also, you can install it on your docker image. Two examples:

```dockerfile
FROM busybox

COPY --from=yingzhuo/docktool:latest /bin/docktool /bin/docktool
```

#### Build it on your computer

```bash
git clone git@github.com:yingzhuo/docktool.git
cd ./docktool/
make install
```

### Sub Command

* [filegen](./.github/filegen.md)
* [filedel](./.github/filedel.md)
* [wait](./.github/wait.md)
* [sleep](./.github/sleep.md)
* [secret](./.github/secret.md)
* [uuid](./.github/uuid.md)
* [test](./.github/test.md)

### Contributing

* Fork it
* Create your feature branch (git checkout -b my-new-feature)
* Commit your changes (git commit -am 'add some feature')
* Push to the branch (git push origin my-new-feature)
* Create new Pull Request

### License

Apache 2.0 license. See [LICENSE](./LICENSE)

### Chang Log

See [chang log](./CHANGELOG.md)

### Authors

* 应卓 - [github](https://github.com/yingzhuo)

See also the list of [contributors](https://github.com/yingzhuo/docktool/graphs/contributors) who participated in this project.

### Acknowledgments

* [subchen](https://github.com/subchen)
	* [go-cli](https://github.com/subchen/go-cli)
