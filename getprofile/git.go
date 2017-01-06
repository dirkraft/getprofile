package getprofile

import (
    "regexp"
    "os/exec"
    "path"
    "path/filepath"
    "os"
)

type gitSyncer struct{}

var gitSyncerInstance *gitSyncer = &gitSyncer{}

func (sync *gitSyncer) Supports(repoUrl string) bool {
    if matched, err := regexp.MatchString(".+@.+:.+", repoUrl); err != nil {
        dbgf("Error checking gitSyncer.Supports(%v)", repoUrl)
        return false
    } else {
        return matched
    }
}

func (sync *gitSyncer) Init() error {
    if repoUrl, err := getConfig(); err != nil {
        return err
    } else if _, err := os.Stat(path.Join(repoPath, ".git")); os.IsNotExist(err) {
        dest := path.Join(basePath, "repo")
        dbgf("Cloning %v to %v", repoUrl, dest)
        return execWithDebug("git", "clone", repoUrl, dest)
    } else {
        return nil // Already cloned
    }
}

func (sync *gitSyncer) Track(absPath, relPath string) error {
    copyToRepo(absPath, relPath)

    cmd := exec.Command("git", "add", relPath)
    cmd.Dir = repoPath
    return cmd.Run()
}

func (sync *gitSyncer) Untrack(relPath string) error {
    cmd := exec.Command("git", "rm", "-f", relPath)
    cmd.Dir = repoPath
    return cmd.Run()
}

func (sync *gitSyncer) Push() error {
    inf("Sending updates out")
    if err := filepath.Walk(repoPath, makeWalkFunc(copyToRepo)); err != nil {
        return err
    } else {
        repoPath := path.Join(basePath, "repo")

        dbg("Command:", "git", "diff", "HEAD", "--quiet")
        cmd := exec.Command("git", "diff", "HEAD", "--quiet")
        cmd.Dir = repoPath
        if err := cmd.Run(); err == nil {
            inf("No changes to push.")
            return nil
        } else if execExitStatus(err, 1) {
            dbg("Changes to push. Continuing.")
        } else if execExitStatus(err, 128) {
            dbg("New git repo. Will try to continue.")
        } else {
            return err
        }

        if cmd.ProcessState.Success() {
            return nil // No changes
        } else if _, err := gitRepoExec("git", "add", "-u"); err != nil {
            return err
        } else if _, err := gitRepoExec("git", "commit", "-m", "auto-commit"); err != nil {
            return err
        } else {
            _, err = gitRepoExec("git", "push")
            return err
        }
    }
}

func (sync *gitSyncer) Pull(force bool) error {
    inf("Copying updates in")
    if prevSha, err := gitSha(); err != nil {
        return err
    } else if _, err := gitRepoExec("git", "pull"); err != nil {
        return err
    } else if nowSha, err := gitSha(); err != nil {
        return err
    } else if force || prevSha != nowSha {
        return filepath.Walk(path.Join(basePath, "repo"), makeWalkFunc(copyToLocal))
    } else {
        inf("No updates. Use --force to overwrite local files anyways.")
        return nil
    }
}

func makeWalkFunc(processor func(absPath, relPath string) error) filepath.WalkFunc {
    return func(absPath string, info os.FileInfo, err error) error {
        if err != nil {
            return err
        } else if info.IsDir() {
            return nil
        } else if relPath, err := filepath.Rel(repoPath, absPath); err != nil {
            return err
        } else if !filepath.HasPrefix(relPath, ".git") {
            return processor(absPath, relPath)
        } else {
            return nil
        }
    }
}

func gitRepoExec(cmdToks ...string) (string, error) {
    dbgf("Checking: %v", cmdToks)
    cmd := exec.Command(cmdToks[0], cmdToks[1:]...)
    cmd.Dir = repoPath
    if bytes, err := cmd.Output(); err != nil {
        return string(bytes), err
    } else {
        return string(bytes), nil
    }
}

func gitSha() (string, error) {
    return gitRepoExec("git", "rev-parse", "HEAD")
}

func copyToRepo(absPath, relPath string) error {
    dbgf("Checking:", relPath)
    src := path.Join(homePath, relPath)
    dest := path.Join(repoPath, relPath)
    if err := os.MkdirAll(path.Dir(dest), 0700); err != nil {
        return err
    } else if _, err := os.Stat(src); os.IsNotExist(err) {
        return nil
    } else {
        dbgf("Copying %v to %v", src, dest)
        return err2(cp.Single(dest, src))
    }
}

func copyToLocal(absPath, relPath string) error {
    inf("Copying:", relPath)
    src := absPath
    dest := path.Join(homePath, relPath)
    if err := os.MkdirAll(path.Dir(dest), 0700); err != nil {
        return err
    } else {
        return err2(cp.Single(dest, src))
    }
}

