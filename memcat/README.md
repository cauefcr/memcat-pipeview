# Memcat - print your program's entire address space

This little program, for windows mac and linux, will print your program's entire address space (your program's memory), can be useful for searching something or in conjunction with pipeview for visualization.

# Usage

```shell
git clone https://github.com/cauefcr/memscroll.git
cd memscroll/memcat
go build . # or go install . to make it available to your PATH
# print some memory!
memcat memcat
# combine with pipeview for a more interesting visualization
memcat firefox | pipeview > /dev/null
```
