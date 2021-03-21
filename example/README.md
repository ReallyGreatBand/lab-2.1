The simplest bood example
=========================

From this directory, try the following commands:

#### Install bood
```
$ go get -u github.com/ReallyGreatBand/lab-2.1/build/cmd/bood
```

#### Build the program
```
$ bood
INFO 2021/03/21 13:39:49 Ninja build file is generated at out/build.ninja
INFO 2021/03/21 13:39:49 Starting the build now
[3/3] go binary archivation by module archiver
  adding: five (deflated 46%)
```

#### Run the program
```
$ out/bin/hello
Hello World
```

#### Run build again 
```
INFO 2021/03/21 13:39:53 Ninja build file is generated at out/build.ninja
INFO 2021/03/21 13:39:53 Starting the build now
[1/1] go binary archivation by module archiver
updating: five (deflated 46%)