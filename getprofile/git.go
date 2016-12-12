package getprofile

import (
    "regexp"
    "log"
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
        log.Printf("Error checking gitSyncer.Supports(%v)", repoUrl)
        return false
    } else {
        return matched
    }
}

func (sync *gitSyncer) Init() error {
    if repoUrl, err := getConfig(); err != nil {
        return err
    } else if basePath, _, err := getConfigPath(); err != nil {
        return err
    } else if _, err := os.Stat(path.Join(basePath, "repo/.git")); os.IsNotExist(err) {
        return exec.Command("git", "clone", repoUrl, path.Join(basePath, "repo")).Run()
    } else {
        return nil // Already cloned
    }
}

func (sync *gitSyncer) Track(absPath string, relPath string) error {
    if basePath, _, err := getConfigPath(); err != nil {
        return err
    } else {
        repoPath := path.Join(basePath, "repo")
        dest := path.Join(repoPath, relPath)
        if err := os.MkdirAll(path.Dir(dest), 0700); err != nil {
            return err
        } else {
            return exec.Command("cp", absPath, dest).Run()
        }
    }
}

func (sync *gitSyncer) Untrack(relPath string) error {
    if basePath, _, err := getConfigPath(); err != nil {
        return err
    } else {
        cmd := exec.Command("git", "rm", "-f", relPath)
        cmd.Dir = path.Join(basePath, "repo")
        return cmd.Run()
    }
}

func (sync *gitSyncer) Out() error {
    if basePath, _, err := getConfigPath(); err != nil {
        return err
    } else if err := filepath.Walk(path.Join(basePath, "repo"), copyToRepo); err != nil {
        return err
    } else {
        repoPath := path.Join(basePath, "repo")

        cmd := exec.Command("git", "diff", "HEAD", "--quiet")
        cmd.Dir = repoPath
        if err := cmd.Run(); err != nil {
            return err
        } else if cmd.ProcessState.Success() {
            return nil // No changes
        }

        cmd = exec.Command("git", "add", ".")
        cmd.Dir = repoPath
        if err := cmd.Run(); err != nil {
            return err
        }

        cmd = exec.Command("git", "commit", "-m", "auto-commit")
        cmd.Dir = repoPath
        if err := cmd.Run(); err != nil {
            return err
        }

        cmd = exec.Command("git", "push")
        cmd.Dir = repoPath
        return cmd.Run()
    }
}

func copyToRepo(path string, info os.FileInfo, err error) error {
    if err != nil {
        return err
    }
    if (info.IsDir()) {
        return nil
    }
    fmt.Println(path)
    fmt.Println(info.IsDir())
    return errors.New("TODO") // TODO
}

func (sync *gitSyncer) In() error {
    if basePath, _, err := getConfigPath(); err != nil {
        return err
    } else {
        repoPath := path.Join(basePath, "repo")
        prevSha, err := gitSha(repoPath)
        if err != nil {
            return err
        }

        cmd := exec.Command("git", "pull")
        cmd.Dir = repoPath
        if err := cmd.Run(); err != nil {
            return err
        }

        nowSha, err := gitSha(repoPath)
        if err != nil {
            return err
        }

        if prevSha == nowSha {
            return nil
        }

        return filepath.Walk(path.Join(basePath, "repo"), copyToLocal)
    }
}

func gitSha(repoPath string) (string, error) {
    cmd := exec.Command("git", "rev-parse", "HEAD")
    cmd.Dir = repoPath
    if err := cmd.Run(); err != nil {
        return "", err
    }
    bytes, err := cmd.Output()
    if err != nil {
        return "", err
    }
    return string(bytes), nil
}

func copyToLocal(path string, info os.FileInfo, err error) error {
    return errors.New("TODO") // TODO
}

