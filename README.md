# go-file-hash-utility
   
   The main goal is to build a Command line utility to get hash value for a file (based on data in file). This Command line utility turns a UTF-8 input stream into a 64bit hash.
   The process reads an arbitrary long byte stream on the standard input and partition them into equally sized chunks.
   The chunk size is specified by the user and must be a multiple of 8.
   
   It is capable of using threads for processing hash values. Default value is 1.
   
   Use two files minimal.txt and sherlock.txt to test the application.
   
   Please review code documentation which explains all the logic and how the application is written.

# How to run the application:

To run application directly:
```
$ cd go-file-hash-utility
$ go build
$ cat filename | ./go-file-hash-utility --window 4096 

optional:
$ cat filename | ./go-file-hash-utility --window 4096 --threads 5

```

 
Using Docker:

To run the application you would need docker installed.

```
$ cd go-file-hash-utility
$ docker build -t filehash .
$ cat filename | docker run -i filehash --window 4096 

optional:
$ cat filename | docker run -i test --window 4096 --threads 5

```
