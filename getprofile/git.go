package getprofile

import (
    "regexp"
    "os/exec"
    "path"
    "path/filepath"
    "os"
    "errors"
    "fmt"
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
        return exec.Command("git", "clone", repoUrl, path.Join(basePath, "repo")).Run()
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

func (sync *gitSyncer) Out() error {
    inf("Sending updates out")
    if err := filepath.Walk(repoPath, makeWalkFunc(copyToRepo)); err != nil {
        return err
    } else {
        repoPath := path.Join(basePath, "repo")

        dbg("Command:", "git", "diff", "HEAD", "--quiet")
        cmd := exec.Command("git", "diff", "HEAD", "--quiet")
        cmd.Dir = repoPath
        if err := cmd.Run(); err != nil {
            if execExitStatus(err, 128) {
                dbg("New git repo. Will try to continue.")
            } else {
                return err
            }
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

func (sync *gitSyncer) In() error {
    inf("Copying updates in")
    if prevSha, err := gitSha(); err != nil {
        return err
    } else if _, err := gitRepoExec("git", "pull"); err != nil {
        return err
    } else if nowSha, err := gitSha(); err != nil {
        return err
    } else if prevSha == nowSha {
        return nil
    } else {
        return filepath.Walk(path.Join(basePath, "repo"), makeWalkFunc(copyToLocal))
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
        } else if !filepath.HasPrefix(relPath, ".git/") {
            return processor(absPath, relPath)
        } else {
            return nil
        }
    }
}

func gitRepoExec(cmdToks ...string) (string, error) {
    dbgf("Command: %v", cmdToks)
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
    inf("Copy to repo:", relPath)
    src := path.Join(homePath, relPath)
    dest := path.Join(repoPath, relPath)
    if err := os.MkdirAll(path.Dir(dest), 0700); err != nil {
        return err
    }

    cmd := fmt.Sprintf("[ -e '%v' ] && cp '%v' '%v' || exit 0", src, src, dest)
    dbg("Command:", cmd)
    return exec.Command("bash", "-c", cmd).Run()
}

func copyToLocal(absPath, relPath string) error {
    inf("Copy to local:", relPath)
    return errors.New("TODO") // TODO
}

