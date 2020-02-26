# goroutine graceful shutdown

```bash
> go run main.go
start: 9
start: 7
start: 4
start: 2
start: 5
start: 6
start: 3
start: 8
start: 1
end: 1
end: 2
end: 3
end: 4
start: 1
start: 9
start: 2
start: 3
start: 7
start: 4
start: 8
start: 6
start: 5
end: 5
end: 1
end: 6
^C2020/02/27 01:27:11 received signal: interrupt # Ctrl + C
2020/02/27 01:27:11 context canceled
end: 2
end: 7
end: 8
end: 3
end: 9
end: 4
end: 5
end: 6
end: 7
end: 8
end: 9
```
