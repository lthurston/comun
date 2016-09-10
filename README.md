# Comun

Comun is a little tool for finding the first common line between two or more files. Once it finds one, it's done.

## How to Use Comun

Do this:

`comun file1.txt file2.txt file3.txt`

If the files have a line in common, comun will output the line, and exit with status 0. Otherwise it'll exit with status 1.

## Why Use Comun?

Dunno, really. I'm using it to assist in finding common ancestors between various branches in git repos, like this:

`comun <(git rev-list master) <(git rev-list development) <git rev-list whatever)`

Knowing shared ancestor commits is useful in seeing a show-branch beyond a merge commit.