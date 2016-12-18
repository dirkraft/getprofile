# getprofile

Synchronize user home files between machines.

<img src="https://travis-ci.org/dirkraft/getprofile.svg?branch=master">

Example: I have `~/.bash_prompt` that I would like to be the same
on all my environments. Once configured (explained in detail below)
I follow any changes to `~/.bash_prompt` with

    getprofile sync

Then the next time I get on any other machine

    getprofile sync

It is called getprofile because I originally wanted to synchronize
.bash_profile between machines. I think of that and all sourced scripts
as my *profile files*. Now I don't actually sync `.bash_profile` itself,
but rather that includes other sourced files. See
[How I use getprofile](#how-i-use-getprofile) for more details.

### Install getprofile

Dev builds: https://github.com/dirkraft/getprofile/releases/tag/dev

For Linux, installation might look like this:

    sudo curl -o /usr/local/bin/getprofile -L https://github.com/dirkraft/getprofile/releases/download/dev/getprofile.linux.amd64
    sudo chmod +x /usr/local/bin/getprofile

### Configure getprofile

getprofile currently stores tracked files via any git repo that can
be addressed in the form`<user>@<host>:<repo>`,
e.g. `git@github.com:dirkraft/profile`.

Currently, only an SSH-accessible git repos are supported (it doesn't
have to be on GitHub!). If you plan to synchronize files with sensitive
information like keys, be sure that the repository is **private**.

    $ getprofile config --help
    NAME:
       getprofile config - Set up getprofile

    USAGE:
       getprofile config REPOSITORY_URL

    DESCRIPTION:

    Set where the profile repository is located. Supported formats:

        git-compatible names: e.g. git@github.com:user/repo

### Use getprofile

Files must first be tracked for getprofile to care about them.
Note that getprofile only operates on the assumption that tracked files
live within the current user's home directory.

    $ getprofile track --help
    NAME:
       getprofile track - Track or untrack a file

    USAGE:
       getprofile track [command options] FILE

    OPTIONS:
       --untrack, -u  Stop tracking a tracked file. The file is not deleted from the local machine.

Then they can be synchronized to (and from) the configured repo.

    $ getprofile sync --help
    NAME:
       getprofile sync - Synchronize profile

    USAGE:
       getprofile sync [command options] [arguments...]

    OPTIONS:
       --watch, -w  Continuously watch and synchronize changes
       --force, -f  Copy from repo to local whether or not there is an update

### What getprofile does not do well

  - Concurrent changes. Conflicting updates will simply overwrite each
    other. If you are working from multiple environments and modifying
    tracked files on each, be careful about the order of `getprofile
    sync`s.

### How I use getprofile

I track .vimrc because I want it identical on all environments without
exception.

I do not track .bashrc or .bash_profile since these aren't consistently
integrated between OSX and different Linux distros. Instead I have the
following tracked by getprofile in my own scripts folder.

    $ tree ~/.dirkraftrc/
    /home/dirkraft/.dirkraftrc/
    ├── source_all.sh
    └── sources
        ├── bash_prompt.sh
        ├── github.sh
        ├── go.sh
        └── ... and so on ...

Here is source_all.sh

    #/usr/bin/env bash

    for f in ~/.dirkraftrc/sources/*; do source $f; done

Then on each machine I add to whatever .bash_* that makes sense

    . ~/.dirkraftrc/source_all.sh

If certains scripts can't be sourced on a machine, I'm still able to
take advantage of `getprofile sync` and selectively source whichever
scripts I like.

### License

```
The MIT License (MIT)
Copyright (c) 2016 Jason Dunkelberger (a.k.a. "dirkraft")

Permission is hereby granted, free of charge, to any person obtaining a
copy of this software and associated documentation files (the
"Software"), to deal in the Software without restriction, including
without limitation the rights to use, copy, modify, merge, publish,
distribute, sublicense, and/or sell copies of the Software, and to
permit persons to whom the Software is furnished to do so, subject to
the following conditions:

The above copyright notice and this permission notice shall be included
in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS
OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.
IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY
CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT,
TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE
SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
```