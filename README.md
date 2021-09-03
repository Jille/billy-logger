# billy-logger

[![GoDoc](https://godoc.org/github.com/Jille/billy-logger?status.svg)](https://godoc.org/github.com/Jille/billy-logger)

This library wraps your billy.Filesystem and logs all calls made and their results. Simplest usage:

```diff
-fs := memfs.New()
+fs := logger.Wrap(memfs.New(), log.Println)
```

If you're debugging multiple Billy filesystems, you can prefix your log lines with this:

```golang
el := log.New(log.Writer(), "emptyfs: ", log.Flags())
rl := log.New(log.Writer(), "router : ", log.Flags())
ml := log.New(log.Writer(), "mergeFS: ", log.Flags())
fs := router.New(logger.Wrap(emptyfs.New(), el.Println))
fs.Mount("/all", logger.Wrap(mergeFS{}, ml.Println))
fs := logger.Wrap(fs, rl.Println)
````
