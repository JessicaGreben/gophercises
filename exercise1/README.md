# Exercise 1 Quiz

https://github.com/gophercises/quiz

### Steps to run the program

Build the docker image:
```
$ docker build -t exercise1 .
```

Run program with default flags. `-i` allows reading from stdin:
```
$ docker run -i exercise1
```

Run program with options:
```
$ docker run -i exercise1 -timer=3
$ docker run -i exercise1 -file=problems.csv
```

### Steps to run tests

From the root of the `exercise1/` directory:
```
$ go test ./...
```
