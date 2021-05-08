# The `hitstree` package

Suppose you have some API with a set of paths like this:

    /api/users/
    /api/users/{UserID}/posts
    /api/users/{UserID}/posts/{PostID}

But full list of paths is unknown and you want to get it from webserver logs.

If you will count each path separately, you could get too much records, because you will have one separate counter for each post of each user.

This is where `hitstree` package comes into play. This package builds a tree of paths and automatically merges all subtrees at some level when number of such subtrees reaches some limit.

See [example/example.go](example/example.go) for quick start:

    $ make run-example
    go build -o build/bin/example example/example.go
    ./build/bin/example
    1	/
    1	/content/bar
    1	/content/baz
    1	/content/foo
    1	/users
    2	/users/{}
    181	/users/{}/posts
    400	/users/{}/posts/{}
